package spec

import (
	"slices"
	"strings"
	"time"
	"unicode"
)

// ============================================================================
//  Helper Functions To Check Specifications
// ============================================================================

// ----------------------------------------------------------------------------
//  Priority should be an uppercase letter A-Z in parentheses (Rule 1)
// ----------------------------------------------------------------------------

// IsPriorityLetter returns true if the letter is an uppercase letter A-Z.
func IsPriorityLetter(letter string) bool {
	return len(letter) == 1 && letter[0] >= 'A' && letter[0] <= 'Z'
}

// IsPriorityMark returns true if the string is a valid priority mark.
//
// A valid priority mark is exactly 3 characters: '(', an uppercase letter A-Z,
// and ')'. For example: "(A)", "(B)", "(Z)".
//
// Examples:
//   - "(A)" → true
//   - "(Z)" → true
//   - "(a)" → false (lowercase)
//   - "(1)" → false (not a letter)
//   - "(AA)" → false (too long)
//   - "A" → false (missing parentheses)
func IsPriorityMark(seg string) bool {
	if len(seg) != PriorityLen {
		return false
	}

	return seg[0] == WrapPriorityOpen.Byte() &&
		seg[2] == WrapPriorityClose.Byte() &&
		seg[1] >= 'A' && seg[1] <= 'Z'
}

// ----------------------------------------------------------------------------
//  Date should be in the format of "YYYY-MM-DD" (Rule 2)
// ----------------------------------------------------------------------------

// IsDate returns true if the string is a date in YYYY-MM-DD format.
//
// If 'validate' is true, it also checks if the date is a valid calendar date,
// but it is 12x slower and more alloc than format-only check (validate=false).
func IsDate(seg string, validate bool) bool {
	if validate {
		_, err := time.Parse(DateFormat, seg)

		return err == nil
	}

	// We do not use time.Parse here for performance reasons.
	// This implementation is 12x faster than time.Parse in benchmarks.
	if len(seg) != DateLen {
		return false
	}

	// Check format: YYYY-MM-DD
	if seg[4] != '-' || seg[7] != '-' {
		return false
	}

	// Check that all other characters are digits.
	for i, ch := range seg {
		if i == 4 || i == 7 {
			continue // skip hyphens
		}

		if ch < '0' || ch > '9' {
			return false
		}
	}

	return true
}

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
	searchAllow := len(allow) > 0

	// traverse each character in the text
	for _, char := range text {
		// Check for latin and other Unicode control characters
		isControl := unicode.IsControl(char) || unicode.Is(unicode.C, char)

		if !isControl {
			continue
		}

		if searchAllow && slices.Contains(allow, char) {
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
