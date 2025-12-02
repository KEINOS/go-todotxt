package spec

import (
	"slices"
	"strings"
	"unicode"
)

// ============================================================================
//  Helper Functions To Check Specifications
// ============================================================================

// ----------------------------------------------------------------------------
//  Any text after "<space char>#" are treated as inline comment.
// ----------------------------------------------------------------------------

// ContainsInlineComment returns the index and true if the text contains an
// inline comment indicator (any Unicode space character followed by "#". e.g.
// " #", "\t#").
//
// If not found, returns -1 and false.
// Note that a comment line (starting with "#") is not considered as an inline
// comment.
func ContainsInlineComment(text string) (int, bool) {
	idxCommentMark := strings.Index(text, PrefixComment.String())

	// no comment mark found or is a comment line (starts with '#')
	if idxCommentMark == -1 || idxCommentMark == 0 {
		return -1, false
	}

	preChar := rune(text[idxCommentMark-1])
	if !unicode.IsSpace(preChar) {
		return -1, false
	}

	return idxCommentMark, true

	// inlineCommentIndicator := DelimSegments.String() + PrefixComment.String()
	// return strings.Contains(text, inlineCommentIndicator)
}

// ----------------------------------------------------------------------------
//  Unless explicitly allowed, control characters are not permitted.
// ----------------------------------------------------------------------------

// ContainsCtlChars returns true if any control character is found in the text
// to prevent parsing malformed or potentially harmful input.
//
// To exclude certain control characters from this check, provide them in the
// 'allow' slice.
func ContainsCtlChars(text string, allow []rune) bool {
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

// ----------------------------------------------------------------------------
//  Task length should not be longer than 64KB (default bufio.MaxScanTokenSize)
// ----------------------------------------------------------------------------

// IsValidTaskLength checks if the length of the task text is within the allowed
// limit.The maximum allowed length is 64 * 1024 bytes (64KB).
func IsValidTaskLength(text string) bool {
	return len(text) <= MaxTaskLength
}
