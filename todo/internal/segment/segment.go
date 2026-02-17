/*
Package segment provides types and functions for working with todo.txt segments.
A segment is a part of a task string, separated by spaces.
*/
package segment

import (
	"strings"

	"github.com/KEINOS/go-todotxt/todo/internal/spec"
)

// ============================================================================
//  Type: Type (Segment type)
// ============================================================================

// Type represents the type of a Segment.
type Type uint16

// Enum flags for Type.
const (
	Undefined      Type = 0
	MarkCompletion Type = 1 << iota
	MarkPriority
	Date
	TagProject
	TagContext
	TagKeyValue
	MarkComment
	PlainText
)

// ----------------------------------------------------------------------------
//  Methods for Type
// ----------------------------------------------------------------------------

// String returns the string representation of the Type.
//
//nolint:cyclop // Simple enum-to-string mapping with no branching logic
func (t Type) String() string {
	switch t {
	case MarkCompletion:
		return "completion mark"
	case MarkPriority:
		return "priority mark"
	case Date:
		return "date"
	case TagProject:
		return "project tag"
	case TagContext:
		return "context tag"
	case TagKeyValue:
		return "key-value pair tag"
	case MarkComment:
		return "comment prefix"
	case PlainText:
		return "plain text"
	case Undefined:
		fallthrough
	default:
		return "undefined"
	}
}

// ============================================================================
//  Type: Segment
// ============================================================================
//  Segments type follows.

// Segment is a token from a todo.txt task string.
//
// For example each of the below are segments:
//
//	"x", "(A)", "2023-12-25", "call", "mom", "+project", "@context", "due:2023-12-31"
type Segment string

// ----------------------------------------------------------------------------
//  Methods for Segment
// ----------------------------------------------------------------------------

// Is returns true if the Segment is of the specified Type. Undefined Type is
// treated as PlainText.
func (s Segment) Is(segType Type) bool {
	if segType == Undefined {
		segType = PlainText
	}

	return s.Type() == segType
}

// IsComment returns true if the segment is a comment.
//
// Comments start with '#' and can contain any text.
func (s Segment) IsComment() bool {
	return len(s) > 0 && rune(s[0]) == spec.PrefixComment.Rune()
}

// IsDate returns true if the segment is a date in YYYY-MM-DD format.
//
// This method only checks the format, not if the date is valid. For full
// validation, use IsDateValid().
func (s Segment) IsDate() bool {
	return spec.IsDate(string(s), false)
}

// IsDateValid returns true if the segment is a valid date in YYYY-MM-DD format
// and is a valid calendar date.
//
// This is similar to IsDate but performs full validation and is 12x slower.
func (s Segment) IsDateValid() bool {
	return spec.IsDate(string(s), true)
}

// IsMarkCompletion returns true if the segment is a completion mark ("x").
func (s Segment) IsMarkCompletion() bool {
	return string(s) == spec.MarkerDone.String()
}

// IsMarkPriority returns true if the segment is a priority mark, like "(A)".
//
// A priority mark is a 3-character string: '(', an uppercase letter, and ')'.
func (s Segment) IsMarkPriority() bool {
	const priorityLen = 3 // "(A)"

	if len(s) != priorityLen {
		return false
	}

	return s[0] == '(' && s[2] == ')' && s[1] >= 'A' && s[1] <= 'Z'
}

// IsPlainText returns true if the segment is plain text.
//
// Plain text is a regular word that does not match any special format.
func (s Segment) IsPlainText() bool {
	return !s.IsMarkCompletion() &&
		!s.IsMarkPriority() &&
		!s.IsDate() &&
		!s.IsTagProject() &&
		!s.IsTagContext() &&
		!s.IsTagKeyValue() &&
		!s.IsComment()
}

// IsTagContext returns true if the segment is a context tag.
//
// A context tag starts with '@' and has at least one more character.
func (s Segment) IsTagContext() bool {
	return len(s) > 1 && rune(s[0]) == spec.PrefixContext.Rune()
}

// IsTagKeyValue returns true if the segment is a key-value pair like "key:value".
// Value may contain colons (e.g., "url:https://example.com").
// Returns false for empty key, empty value, or missing separator.
func (s Segment) IsTagKeyValue() bool {
	str := string(s)
	sep := spec.SepKeyValue.String()
	colonIndex := strings.Index(str, sep)

	// The colon must not be at the beginning or end.
	if colonIndex <= 0 || colonIndex >= len(str)-1 {
		return false
	}

	// The value part (after the first colon) must not be empty.
	// This handles cases like "key::".
	return str[colonIndex+1] != ':'
}

// IsTagProject returns true if the segment is a project tag.
//
// A project tag starts with '+' and has at least one more character.
func (s Segment) IsTagProject() bool {
	return len(s) > 1 && rune(s[0]) == spec.PrefixProject.Rune()
}

// String returns the segment as a string.
func (s Segment) String() string {
	return string(s)
}

// Type returns the Type of the Segment. Non-matching segments are classified
// as PlainText.
func (s Segment) Type() Type {
	switch {
	case s.IsMarkCompletion():
		return MarkCompletion
	case s.IsMarkPriority():
		return MarkPriority
	case s.IsDate():
		return Date
	case s.IsTagProject():
		return TagProject
	case s.IsTagContext():
		return TagContext
	case s.IsTagKeyValue():
		return TagKeyValue
	case s.IsComment():
		return MarkComment
	case s.IsPlainText():
		fallthrough
	default:
		return PlainText
	}
}

// ============================================================================
//  Type: Segments
// ============================================================================

// Segments represents a slice of Segment.
type Segments []Segment

// ----------------------------------------------------------------------------
//  Methods for Segments
// ----------------------------------------------------------------------------

// Types returns a bitwise OR of all segment types contained in this Segments.
// The result can be used with bitwise operations to check for type presence.
func (s Segments) Types() Type {
	result := Undefined

	for _, seg := range s {
		result |= seg.Type()
	}

	return result
}
