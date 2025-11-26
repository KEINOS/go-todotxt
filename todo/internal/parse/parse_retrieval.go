package parse

import (
	"slices"
	"strings"

	"github.com/KEINOS/go-todotxt/todo/internal/spec"
)

// ============================================================================
//  Methods of Parsed to retrieve task components (GETTERS)
// ============================================================================

// ----------------------------------------------------------------------------
// Official spec getter methods
// ----------------------------------------------------------------------------

// String returns the original task string.
func (p *Parsed) String() string {
	return p.originalText
}

// Priority returns the priority letter if the task has a priority segment.
func (p *Parsed) Priority() string {
	if p.IsCommentLine() {
		return ""
	}

	maxIndex := 2 // Priority can only be at index 1 if completion mark is at index 0

	for index, seg := range p.Segments {
		if index > maxIndex {
			break
		}

		if seg.IsMarkPriority() {
			return string(seg)[1:2] // Extract the priority letter without parentheses
		}
	}

	return ""
}

// DateCompleted returns the completion date as a string in "YYYY-MM-DD" format
// if the task is marked as completed and has a completion date segment.
func (p *Parsed) DateCompleted() string {
	if p.IsCommentLine() || !p.IsDone() {
		return ""
	}

	index := p.consumeHeadParts() // Skip done mark and priority

	if index < len(p.Segments) && p.Segments[index].IsDate() {
		return p.Segments[index].String()
	}

	return ""
}

// DateCreated returns the creation date as a string in "YYYY-MM-DD" format.
//
//   - If the task is completed, it returns the second date segment found (after
//     the completion date).
//   - If the task is not completed, it returns the first date segment found before
//     the description text.
func (p *Parsed) DateCreated() string {
	if p.IsCommentLine() {
		return ""
	}

	isDoneTask := p.IsDone()
	index := p.consumeHeadParts() // Skip done mark and priority

	// First date segment (only if not done)
	if index < len(p.Segments) && p.Segments[index].IsDate() {
		if !isDoneTask {
			return p.Segments[index].String()
		}

		index++
	}

	// Second date segment (only if done)
	if index < len(p.Segments) && p.Segments[index].IsDate() {
		if isDoneTask {
			return p.Segments[index].String()
		}
	}

	return ""
}

// Contexts returns a slice of context tags found in the task.
//
// Contexts are segments that start with '@'. Repeated contexts are treated as
// one.
func (p *Parsed) Contexts() []string {
	if p.IsCommentLine() {
		return nil
	}

	var contexts []string

	index := p.consumeHeadParts() // Skip done mark and priority

	for _, seg := range p.Segments[index:] {
		if seg.IsTagContext() {
			contexts = append(contexts, seg.String()[1:]) // Exclude '@' prefix
		}
	}

	return slices.Compact(contexts)
}

// Projects returns a slice of project tags found in the task.
//
// Projects are segments that start with '+'. Repeated projects are treated as
// one.
func (p *Parsed) Projects() []string {
	if p.IsCommentLine() {
		return nil
	}

	var projects []string

	index := p.consumeHeadParts() // Skip done mark and priority

	for _, seg := range p.Segments[index:] {
		if seg.IsTagProject() {
			projects = append(projects, seg.String()[1:]) // Exclude '+' prefix
		}
	}

	return slices.Compact(projects)
}

// Description returns the task description text without:
//
//   - Completion mark
//   - Priority
//   - Date headers at the beginning
//   - Comments (both inline comments and comment lines)
//
// Due to the free-form nature of the todo.txt format, this method does not
// reconstruct the description from parsed segments.
//
//	Example: "x   (A)   Do something   +project   # comment" --> "Do something   +project"
//	Example: "# This is a comment line" --> ""
func (p *Parsed) Description() string {
	// Comment lines have no description, only comment
	if p.IsCommentLine() {
		return ""
	}

	index := p.consumeHeadParts() // segments to skip

	// Skip completion date
	for index < len(p.Segments) && p.Segments[index].IsDate() {
		index++
	}

	// No segments to skip, return original text
	if index == 0 {
		result := strings.TrimSpace(p.originalText)
		// Remove inline comment if present
		if p.HasInlineComment() {
			result = strings.TrimSpace(p.trimInlineComment(result))
		}

		return result
	}

	// Use offsetAfterSegments to accurately skip header tokens,
	// even when identical strings appear (e.g., same completion/creation dates)
	offset := p.offsetAfterSegments(index)
	remaining := p.originalText[offset:]
	result := strings.TrimSpace(remaining)

	// Remove inline comment if present
	if p.HasInlineComment() {
		result = strings.TrimSpace(p.trimInlineComment(result))
	}

	return result
}

// KeyValues returns a slice of key-value pairs found in the task.
//
// Splits segments at the first ':' to form key-value pairs. The value part may
// contain additional colons (e.g., "url:https://example.com").
// Returns nil if no key-value pairs are found.
//
// Note: Results are cached after the first call. The slice preserves the order
// of appearance in the original task string and may contain duplicate keys.
func (p *Parsed) KeyValues() []KeyValue {
	if p.IsCommentLine() {
		return nil
	}

	if p.keyValueInit {
		return p.keyValueCache // decrease mem alloc 11 --> 2, 280 --> 103 ns/op
	}

	const maxNumParts = 2 // Key and value

	separator := spec.SepKeyValue.String()

	index := p.consumeHeadParts() // Skip done mark and priority

	for _, seg := range p.Segments[index:] {
		if !seg.IsKeyValue() {
			continue
		}

		parts := strings.SplitN(seg.String(), separator, maxNumParts)

		// ignore 'key: ' or ':value ' cases
		if len(parts) == maxNumParts {
			p.keyValueCache = append(p.keyValueCache, KeyValue{
				Key:   parts[0],
				Value: parts[1],
			})
		}
	}

	p.keyValueInit = true

	return p.keyValueCache
}

// ----------------------------------------------------------------------------
//  Extended getter methods
// ----------------------------------------------------------------------------
//  The official spec does not define about comments. In this package, we treat
//  comments as follows:
//   - A comment line is a line that begins with '#' character.
//   - An inline comment is a segment that begins with '#' character, but not at
//     the beginning of the line.
// ----------------------------------------------------------------------------

// Comment returns the comment text found in the task.
//
// It returns the entire task string with leading/trailing whitespace trimmed
// if it is a comment line.
// For inline comments, it returns the text starting from the '#' character
// to the end of the line.
func (p *Parsed) Comment() string {
	if p.IsCommentLine() {
		return strings.TrimSpace(p.originalText)
	}

	if p.HasInlineComment() {
		return *p.inlineComment
	}

	return ""
}

// Components returns all parsed components as a Components struct.
//
// This method aggregates all individual getter results into a single structured
// object, useful for JSON serialization or comprehensive task inspection.
func (p *Parsed) Components() Components {
	return Components{
		Priority:         p.Priority(),
		DateCompleted:    p.DateCompleted(),
		DateCreated:      p.DateCreated(),
		Description:      p.Description(),
		Contexts:         p.Contexts(),
		Projects:         p.Projects(),
		KeyValues:        p.KeyValues(),
		Comment:          p.Comment(),
		IsDone:           p.IsDone(),
		IsCommentLine:    p.IsCommentLine(),
		HasInlineComment: p.HasInlineComment(),
		PosInlineComment: p.InlineCommentPos(),
	}
}
