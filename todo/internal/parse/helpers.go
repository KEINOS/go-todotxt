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

// getPosComment returns the position of inline comment in the original text.
//
// It will return -1 if no inline comment is found.
func getPosComment(p *Parsed) int {
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

// hasControlChars returns true if any control character is found in the text.
//
// This prevents parsing of malformed or potentially harmful input. To exclude
// certain control characters from this check, provide them in the 'allow' slice.
func hasControlChars(text string, allow []rune) bool {
	isInAllowed := func(r rune) bool {
		if len(allow) == 0 {
			return false
		}

		return slices.Contains(allow, r)
	}

	// traverse each character in the text
	for _, char := range text {
		// Check for latin and other Unicode control characters
		isControl := unicode.IsControl(char) || unicode.Is(unicode.C, char)

		if !isControl {
			continue
		}

		if isInAllowed(char) {
			continue
		}

		return true // found disallowed control character
	}

	return false
}
