package todoparse

import (
	"strings"

	"github.com/KEINOS/go-todotxt/todo/internal/spec"
)

// ============================================================================
//  Methods of Parsed for status checks
// ============================================================================
//  Methods to obtain information about the current status of the task.

// HasPriority returns true if the task has a priority segment.
func (p *Parsed) HasPriority() bool {
	return p.Priority() != ""
}

// HasInlineComment returns true if the task has an inline comment segment or
// it is a comment line.
func (p *Parsed) HasInlineComment() bool {
	// Use cache
	if p.inlineComment != nil {
		return *p.inlineComment != ""
	}

	p.indexComment = new(int)
	*p.indexComment = -1

	// It begins with '#'. Set the whole line as comment
	if p.IsCommentLine() {
		p.setComment(p.originalText)
		*p.indexComment = strings.Index(p.originalText, spec.PrefixComment.String())

		return true
	}

	// Search for inline comment segment
	commentFound := false

	for _, seg := range p.Segments {
		if seg.IsComment() {
			commentFound = true // not comment line but inline comment

			break
		}
	}

	if !commentFound {
		p.setComment("") // empty cache and return

		return false
	}

	// Search rune by rune to find the index of inline comment
	foundPos := findInlineCommentPos(p)

	*p.indexComment = foundPos // cache the position

	if foundPos == -1 {
		p.setComment("") // should not happen, but just in case

		return false
	}

	foundComment := p.originalText[foundPos:] // include '#'

	p.setComment(foundComment)

	return *p.inlineComment != ""
}

// InlineCommentPos returns the position index of the inline comment in the
// original task string.
//
// It returns -1 if no inline comment is found.
func (p *Parsed) InlineCommentPos() int {
	if p.indexComment != nil {
		return *p.indexComment
	}

	if !p.HasInlineComment() {
		return -1 // not found
	}

	return *p.indexComment
}

// IsCommentLine returns true if the task is a comment line (begins with '#').
func (p *Parsed) IsCommentLine() bool {
	if p.isCommentLine != nil {
		return *p.isCommentLine
	}

	p.isCommentLine = new(bool)
	*p.isCommentLine = len(p.Segments) > 0 && p.Segments[0].IsComment()

	return *p.isCommentLine
}

// IsCompleted returns true if the task is marked as completed (i.e., starts with 'x' segment).
func (p *Parsed) IsCompleted() bool {
	if p.isDone != nil {
		return *p.isDone
	}

	p.isDone = new(bool)
	*p.isDone = len(p.Segments) > 0 && p.Segments[0].IsMarkCompletion()

	return *p.isDone
}

// IsDone returns true if the task is marked as completed (i.e., starts with 'x' segment).
//
// It is an alias of IsCompleted for backward compatibility.
func (p *Parsed) IsDone() bool {
	return p.IsCompleted()
}
