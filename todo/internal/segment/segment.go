/*
Package segment provides types and functions for working with todo.txt segments.
A segment is a part of a task string, separated by spaces.
*/
package segment

import (
	"strings"

	"github.com/KEINOS/go-todotxt/todo/internal/spec"
)

// ----------------------------------------------------------------------------
//  Type: Segment
// ----------------------------------------------------------------------------

// Segment is a token from a todo.txt task string.
//
// Example: "x", "(A)", "2023-12-25", "call", "mom", "+project", "@context", "due:2023-12-31".
type Segment string

// ----------------------------------------------------------------------------
//  Methods for Segment
// ----------------------------------------------------------------------------

// IsComment returns true if the segment is a comment.
//
// Comments start with '#' and can contain any text.
func (s Segment) IsComment() bool {
	return len(s) > 0 && rune(s[0]) == spec.PrefixComment.Rune()
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

// IsDate returns true if the segment is a date in YYYY-MM-DD format.
//
// This method only checks the format, not if the date is valid.
func (s Segment) IsDate() bool {
	// We do not use time.Parse here for performance reasons.
	// This implementation is 7x faster than time.Parse in benchmarks.
	const dateLen = 10 // "YYYY-MM-DD"

	if len(s) != dateLen {
		return false
	}

	// Check format: YYYY-MM-DD
	if s[4] != '-' || s[7] != '-' {
		return false
	}

	// Check that all other characters are digits.
	for i, ch := range s {
		if i == 4 || i == 7 {
			continue
		}

		if ch < '0' || ch > '9' {
			return false
		}
	}

	return true
}

// IsTagProject returns true if the segment is a project tag.
//
// A project tag starts with '+' and has at least one more character.
func (s Segment) IsTagProject() bool {
	return len(s) > 1 && rune(s[0]) == spec.PrefixProject.Rune()
}

// IsTagContext returns true if the segment is a context tag.
//
// A context tag starts with '@' and has at least one more character.
func (s Segment) IsTagContext() bool {
	return len(s) > 1 && rune(s[0]) == spec.PrefixContext.Rune()
}

// IsKeyValue returns true if the segment is a key-value pair like "key:value".
// Value may contain colons (e.g., "url:https://example.com").
// Returns false for empty key, empty value, or missing separator.
func (s Segment) IsKeyValue() bool {
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

// IsPlainText returns true if the segment is plain text.
//
// Plain text is a regular word that does not match any special format.
func (s Segment) IsPlainText() bool {
	return !s.IsMarkCompletion() &&
		!s.IsMarkPriority() &&
		!s.IsDate() &&
		!s.IsTagProject() &&
		!s.IsTagContext() &&
		!s.IsKeyValue() &&
		!s.IsComment()
}

// String returns the segment as a string.
func (s Segment) String() string {
	return string(s)
}
