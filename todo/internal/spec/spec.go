/*
Package spec defines shared constants for the todo.txt format.
*/
package spec

// Mark is a special character in the todo.txt format.
type Mark rune

// Byte returns the byte representation of the Mark.
func (m Mark) Byte() byte {
	return byte(m)
}

// Rune returns the rune representation of the Mark.
func (m Mark) Rune() rune {
	return rune(m)
}

// String returns the string representation of the mark.
func (m Mark) String() string {
	return string(m)
}

const (
	// PrefixProject is the marker for project tags.
	PrefixProject Mark = '+'
	// PrefixContext is the marker for context tags.
	PrefixContext Mark = '@'
	// PrefixComment is the marker for comments (extended syntax).
	PrefixComment Mark = '#'
	// MarkerDone is the marker indicating a completed task.
	MarkerDone Mark = 'x'
	// SepKeyValue is the separator for key-value pairs.
	SepKeyValue Mark = ':'
	// DelimSegments is the delimiter for segments in a task string.
	DelimSegments Mark = ' '
)

const (
	// DateFormat is the standard date format in todo.txt (YYYY-MM-DD).
	DateFormat = "2006-01-02"
)
