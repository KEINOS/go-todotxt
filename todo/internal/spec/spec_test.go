package spec

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMark_methods(t *testing.T) {
	t.Parallel()

	// Test data for Mark methods
	// All methods return zero/empty values for non-printable characters.
	testCases := []struct {
		title        string
		mark         byte
		expectByte   byte
		expectRune   rune
		expectString string
	}{
		// Printable ASCII characters (0x20-0x7E)
		{
			title:        "Space (0x20) - lower bound of printable",
			mark:         0x20,
			expectByte:   0x20,
			expectRune:   ' ',
			expectString: " ",
		},
		{
			title:        "Exclamation mark (0x21) - lower bound + 1",
			mark:         0x21,
			expectByte:   0x21,
			expectRune:   '!',
			expectString: "!",
		},
		{
			title:        "Right brace (0x7D) - upper bound - 1",
			mark:         0x7D,
			expectByte:   0x7D,
			expectRune:   '}',
			expectString: "}",
		},
		{
			title:        "Tilde (0x7E) - upper bound of printable",
			mark:         0x7E,
			expectByte:   0x7E,
			expectRune:   '~',
			expectString: "~",
		},
		{
			title:        "Plus sign (project prefix)",
			mark:         '+',
			expectByte:   0x2B,
			expectRune:   '+',
			expectString: "+",
		},
		{
			title:        "At sign (context prefix)",
			mark:         '@',
			expectByte:   0x40,
			expectRune:   '@',
			expectString: "@",
		},
		// Edge cases - outside printable range
		// All methods return zero/empty for non-printable characters
		{
			title:        "Zero value (0x00) - control character",
			mark:         0x00,
			expectByte:   0x00, // InvalidMark
			expectRune:   0,    // InvalidRune
			expectString: "",   // empty string
		},
		{
			title:        "Tab (0x09) - control character",
			mark:         0x09,
			expectByte:   0x00, // InvalidMark
			expectRune:   0,    // InvalidRune
			expectString: "",   // empty string
		},
		{
			title:        "Unit Separator (0x1F) - just below printable",
			mark:         0x1F,
			expectByte:   0x00, // InvalidMark
			expectRune:   0,    // InvalidRune
			expectString: "",   // empty string
		},
		{
			title:        "DEL (0x7F) - control character",
			mark:         0x7F,
			expectByte:   0x00, // InvalidMark
			expectRune:   0,    // InvalidRune
			expectString: "",   // empty string
		},
		{
			title:        "Extended ASCII start (0x80)",
			mark:         0x80,
			expectByte:   0x00, // InvalidMark
			expectRune:   0,    // InvalidRune
			expectString: "",   // empty string
		},
		{
			title:        "Max byte value (0xFF)",
			mark:         0xFF,
			expectByte:   0x00, // InvalidMark
			expectRune:   0,    // InvalidRune
			expectString: "",   // empty string
		},
	}

	for index, test := range testCases {
		mark := Mark(test.mark)
		title := fmt.Sprintf("Test #%d: %s", index+1, test.title)

		t.Run(title+" (Byte)", func(t *testing.T) {
			t.Parallel()

			require.Equal(t, test.expectByte, mark.Byte())
		})

		t.Run(title+" (Rune)", func(t *testing.T) {
			t.Parallel()

			require.Equal(t, test.expectRune, mark.Rune())
		})

		t.Run(title+" (String)", func(t *testing.T) {
			t.Parallel()

			require.Equal(t, test.expectString, mark.String())
		})
	}
}

func TestMark_IsPrintable(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		title    string
		mark     byte
		expected bool
	}{
		// Printable ASCII (0x20-0x7E) - should return true
		{
			title:    "Space (0x20) - lower bound",
			mark:     0x20,
			expected: true,
		},
		{
			title:    "Exclamation (0x21) - lower bound + 1",
			mark:     0x21,
			expected: true,
		},
		{
			title:    "Right brace (0x7D) - upper bound - 1",
			mark:     0x7D,
			expected: true,
		},
		{
			title:    "Tilde (0x7E) - upper bound",
			mark:     0x7E,
			expected: true,
		},
		{
			title:    "Plus sign (+)",
			mark:     '+',
			expected: true,
		},
		{
			title:    "At sign (@)",
			mark:     '@',
			expected: true,
		},
		{
			title:    "Hash (#)",
			mark:     '#',
			expected: true,
		},
		{
			title:    "Lowercase x",
			mark:     'x',
			expected: true,
		},
		{
			title:    "Colon (:)",
			mark:     ':',
			expected: true,
		},
		// Non-printable (control characters) - should return false
		{
			title:    "NUL (0x00)",
			mark:     0x00,
			expected: false,
		},
		{
			title:    "Tab (0x09)",
			mark:     0x09,
			expected: false,
		},
		{
			title:    "Line Feed (0x0A)",
			mark:     0x0A,
			expected: false,
		},
		{
			title:    "Carriage Return (0x0D)",
			mark:     0x0D,
			expected: false,
		},
		{
			title:    "Unit Separator (0x1F) - just below printable",
			mark:     0x1F,
			expected: false,
		},
		{
			title:    "DEL (0x7F) - just above printable",
			mark:     0x7F,
			expected: false,
		},
		{
			title:    "Extended ASCII start (0x80)",
			mark:     0x80,
			expected: false,
		},
		{
			title:    "High byte (0xFF)",
			mark:     0xFF,
			expected: false,
		},
	}

	for index, tc := range testCases {
		t.Run(fmt.Sprintf("Test #%d: %s", index+1, tc.title), func(t *testing.T) {
			t.Parallel()

			m := Mark(tc.mark)
			require.Equal(t, tc.expected, m.IsPrintable())
		})
	}
}

// TestMark_constants verifies that all defined constants are printable ASCII.
func TestMark_constants(t *testing.T) {
	t.Parallel()

	constants := []struct {
		name string
		mark Mark
	}{
		{"PrefixProject", PrefixProject},
		{"PrefixContext", PrefixContext},
		{"PrefixComment", PrefixComment},
		{"MarkerDone", MarkerDone},
		{"SepKeyValue", SepKeyValue},
		{"DelimSegments", DelimSegments},
	}

	for _, c := range constants {
		t.Run(c.name+" is printable", func(t *testing.T) {
			t.Parallel()

			require.True(t, c.mark.IsPrintable(),
				"%s (0x%02X) should be a printable ASCII character", c.name, c.mark)
		})
	}
}
