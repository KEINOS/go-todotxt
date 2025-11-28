/*
Package spec defines shared constants for the todo.txt format.

The Mark type represents single-byte ASCII markers used in todo.txt syntax.
Only printable ASCII characters (0x20-0x7E) are valid markers.
*/
package spec

const (
	// DelimSegments is the delimiter for segments in a task string.
	DelimSegments Mark = ' '
	// MarkerDone is the marker indicating a completed task.
	MarkerDone Mark = 'x'
	// PrefixComment is the marker for comments (extended syntax).
	PrefixComment Mark = '#'
	// PrefixContext is the marker for context tags.
	PrefixContext Mark = '@'
	// PrefixProject is the marker for project tags.
	PrefixProject Mark = '+'
	// SepKeyValue is the separator for key-value pairs.
	SepKeyValue Mark = ':'
)

const (
	// DateFormat is the standard date format in todo.txt (YYYY-MM-DD).
	DateFormat = "2006-01-02"
	// InvalidMark is the zero value returned by Byte() for non-printable marks.
	InvalidMark byte = 0x00
	// InvalidRune is the zero value returned by Rune() for non-printable marks.
	InvalidRune rune = 0
	// MaxTaskLength is the maximum allowed length for a task string. Equivalent
	// to bufio.MaxScanTokenSize (64K).
	MaxTaskLength int = 64 * 1024 // 65536
)

// Mark is a single-byte marker character in the todo.txt format.
//
// Valid marks are printable ASCII characters in the range 0x20 (space) to 0x7E
// (tilde). Control characters and extended ASCII are not valid markers.
type Mark byte

// Byte returns the byte value, or InvalidMark (0x00) if not printable.
func (m Mark) Byte() byte {
	if !m.IsPrintable() {
		return InvalidMark
	}

	return byte(m)
}

// Rune returns the rune value, or InvalidRune (0) if not printable.
func (m Mark) Rune() rune {
	if !m.IsPrintable() {
		return InvalidRune
	}

	return rune(m)
}

// String returns the string value, or empty string if not printable.
func (m Mark) String() string {
	if !m.IsPrintable() {
		return ""
	}

	return string(m)
}

// IsPrintable reports whether the Mark is a printable ASCII character.
//
// Printable ASCII characters are in the range 0x20 (space) to 0x7E (tilde).
// Control characters (0x00-0x1F, 0x7F) and extended bytes (0x80-0xFF) return
// false.
func (m Mark) IsPrintable() bool {
	return m >= 0x20 && m <= 0x7E
}
