/*
Package parse provides read-only parsing of todo.txt task strings.

It extracts structured information from a task string.
For modifying tasks, use the "task" package.
*/
package parse

import (
	"strings"
	"unicode"

	"github.com/KEINOS/go-todotxt/todo/internal/segment"
)

// Parsed holds the segments of a todo.txt task string.
type Parsed struct {
	segment.Segments // Segments of the task string.

	// The original task string.
	originalText string
	// Caches
	inlineComment *string           // Cached inline comment.
	isCommentLine *bool             // Cached comment line status.
	isDone        *bool             // Cached done status.
	indexSkip     *int              // Cached index to skip header parts.
	indexComment  *int              // Cached index of the inline comment.
	keyValueCache map[string]string // Cached key-value pairs.
	// Options
	options          []Option // Options for parsing.
	allowedCtrlChars []rune   // Allowed control characters.
}

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// FromTaskString creates a new Parsed object from a todo.txt task string.
// It splits the string by whitespace into segments.
// It returns an empty slice if the input is empty or only whitespace.
func FromTaskString(todoTxt string, opts ...Option) (*Parsed, error) {
	parsed := new(Parsed)

	parsed.options = opts

	// Set the original task string and parse it.
	err := parsed.Parse(todoTxt)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}

// ----------------------------------------------------------------------------
//  Methods
// ----------------------------------------------------------------------------
//  - For retrieving parsed components, see parse_retrieval.go
//  - For status checks, see parse_status.go
// ----------------------------------------------------------------------------

// SetText updates the original task string.
// This method is for the "task" package to modify the task string directly.
// After calling this, you should call Parse() to update the segments and caches.
func (p *Parsed) SetText(newText string) {
	p.originalText = newText
}

// Parse updates the Parsed object by parsing a new task string.
// If 'opts' are provided, they override the original options.
func (p *Parsed) Parse(taskTxtUpdate string, opts ...Option) error {
	p.resetCaches()                // Clear all caches.
	p.originalText = taskTxtUpdate // Set/update the original text.

	if len(opts) > 0 {
		p.options = opts // Override options.
	}

	// Apply options.
	for _, opt := range p.options {
		opt(p)
	}

	if hasControlChars(p.originalText, p.allowedCtrlChars) {
		return wrapError(ErrTaskWithControlChars, "invalid task string")
	}

	// Split the string into segments.
	fields := strings.Fields(p.originalText)
	p.Segments = make([]segment.Segment, len(fields))

	for i, field := range fields {
		p.Segments[i] = segment.Segment(field)
	}

	// Cached values are computed on demand (lazy evaluation).
	// Use Components() to get all fields pre-computed.

	return nil
}

// ----------------------------------------------------------------------------
//  Helper Methods
// ----------------------------------------------------------------------------

// consumeHeadParts returns the index after the completion marker and priority.
// Use this to skip optional header parts of the task.
// It does not skip date tokens.
func (p *Parsed) consumeHeadParts() int {
	// Return cached value if available.
	if p.indexSkip != nil {
		return *p.indexSkip
	}

	p.indexSkip = new(int)
	index := 0

	// Skip completion marker "x" if it is the first segment.
	if index < len(p.Segments) && p.Segments[index].IsMarkCompletion() {
		index++
	}

	// Skip priority marker "(A)" if it is next.
	if index < len(p.Segments) && p.Segments[index].IsMarkPriority() {
		index++
	}

	*p.indexSkip = index

	return *p.indexSkip
}

// offsetAfterSegments returns the byte offset in the original text
// after consuming the specified number of segments from the start.
//
// It returns -1 if the count exceeds the number of segments.
func (p *Parsed) offsetAfterSegments(count int) int {
	offset := 0

	const notFound = -1

	if count > len(p.Segments) {
		return notFound
	}

	for i := 0; i < count && i < len(p.Segments); i++ {
		segText := p.Segments[i].String()

		idx := strings.Index(p.originalText[offset:], segText)
		if idx == -1 {
			return notFound // Segment not found; should not happen.
		}

		offset += idx + len(segText)

		// drop whitespaces until next segment
		for offset < len(p.originalText) && unicode.IsSpace(rune(p.originalText[offset])) {
			offset++
		}
	}

	return offset
}

// resetCaches clears all cached values.
// Call this method if the Parsed object is modified directly.
func (p *Parsed) resetCaches() {
	p.inlineComment = nil
	p.isCommentLine = nil
	p.indexComment = nil
	p.isDone = nil
	p.indexSkip = nil
	p.keyValueCache = nil
}

// setComment caches the inline comment text.
func (p *Parsed) setComment(cmt string) {
	if p.inlineComment == nil {
		p.inlineComment = new(string)
	}

	*p.inlineComment = cmt
}

// trimInlineComment removes the cached inline comment from the given text.
//
// It relies on HasInlineComment() having populated p.inlineComment and trims
// only when the cached comment fragment is found inside the provided text.
func (p *Parsed) trimInlineComment(text string) string {
	if text == "" {
		return text
	}

	if !p.HasInlineComment() || p.inlineComment == nil {
		return text
	}

	comment := *p.inlineComment

	idx := strings.LastIndex(text, comment)
	if idx == -1 {
		return text
	}

	return text[:idx]
}
