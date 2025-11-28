package parse

import (
	"slices"
	"unicode"

	"github.com/KEINOS/go-todotxt/todo/internal/spec"
)

// ----------------------------------------------------------------------------
//  Helper Functions (private)
// ----------------------------------------------------------------------------
// This file contains helper functions for parsing todo.txt lines to reduce
// complexity in the main parsing logic.

// findInlineCommentPos returns the position of inline comment in the original text.
// Returns -1 if not found.
func findInlineCommentPos(p *Parsed) int {
	const notFound = -1

	for index, char := range p.originalText {
		if char != spec.PrefixComment.Rune() {
			continue
		}

		if index == 0 {
			return 0 // whole line is comment
		}

		// Check if previous char is space (to avoid '#' in words)
		prevChar := rune(p.originalText[index-1])
		if !unicode.IsSpace(prevChar) {
			continue // '#' is part of a word, skip
		}

		if prevChar == spec.DelimSegments.Rune() {
			return index
		}

		// Check if prevChar is in allowed control chars and treat as space
		if slices.Contains(p.allowedCtrlChars, prevChar) {
			return index
		}
	}

	return notFound
}
