package spec

import (
	"slices"
	"unicode"
)

// ============================================================================
//  Helper Functions To Check Specifications
// ============================================================================

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
