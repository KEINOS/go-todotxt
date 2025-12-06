package spec

import (
	"fmt"
	"strings"
	"testing"

	"github.com/KEINOS/go-todotxt/todo/internal/testdata"
	"github.com/stretchr/testify/require"
)

// ============================================================================
//  Mark methods
// ============================================================================

// ----------------------------------------------------------------------------
//  Mark.Byte(), Mark.Rune(), Mark.String()
// ----------------------------------------------------------------------------

func TestMark_methods(t *testing.T) {
	t.Parallel()

	for index, test := range dataTypeConversion {
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

// ----------------------------------------------------------------------------
//  Mark.IsPrintable()
// ----------------------------------------------------------------------------

func TestMark_IsPrintable(t *testing.T) {
	t.Parallel()

	for index, test := range dataPrintableDetection {
		t.Run(fmt.Sprintf("Test #%d: %s", index+1, test.title), func(t *testing.T) {
			t.Parallel()

			expect := test.expected
			actual := Mark(test.mark).IsPrintable()

			require.Equal(t, expect, actual)
		})
	}
}

// ----------------------------------------------------------------------------
//  Mark constants
// ----------------------------------------------------------------------------

// All predefined constants must be printable ASCII.
func TestMark_constants(t *testing.T) {
	t.Parallel()

	for index, markConst := range []struct {
		name string
		mark Mark
	}{
		// Predefined marks
		{"PrefixProject", PrefixProject},
		{"PrefixContext", PrefixContext},
		{"PrefixComment", PrefixComment},
		{"MarkerDone", MarkerDone},
		{"SepKeyValue", SepKeyValue},
		{"DelimSegments", DelimSegments},
		{"WrapPriorityOpen", WrapPriorityOpen},
		{"WrapPriorityClose", WrapPriorityClose},
	} {
		title := fmt.Sprintf("Test #%d: %s", index+1, markConst.name)

		t.Run(title, func(t *testing.T) {
			t.Parallel()

			require.True(t, markConst.mark.IsPrintable(),
				"%s (0x%02X) should be a printable ASCII character", markConst.name, markConst.mark)
		})
	}
}

// ============================================================================
//  Checker functions
// ============================================================================

// ----------------------------------------------------------------------------
//  IsDate()
// ----------------------------------------------------------------------------

func TestIsDate(t *testing.T) {
	t.Parallel()

	// use common test data
	for index, test := range testdata.Dates {
		title := fmt.Sprintf("Test #%d: %s", index+1, test.Title)

		t.Run(title+" (validate = false)", func(t *testing.T) {
			t.Parallel()

			validate := false
			expect := test.IsDate
			actual := IsDate(test.Input, validate)

			require.Equal(t, expect, actual,
				"IsDate('%s', %v) should return %v", test.Input, validate, expect)
		})

		t.Run(title+" (validate = true)", func(t *testing.T) {
			t.Parallel()

			validate := true
			expect := test.IsValid
			actual := IsDate(test.Input, validate)

			require.Equal(t, expect, actual,
				"IsDate('%s', %v) should return %v", test.Input, validate, expect)
		})
	}
}

// ----------------------------------------------------------------------------
//  IsPriorityLetter()
// ----------------------------------------------------------------------------

func TestIsPriorityLetter(t *testing.T) {
	t.Parallel()

	// use common test data
	for index, test := range testdata.Priorities {
		title := fmt.Sprintf("Test #%d: %s", index+1, test.Title)

		t.Run(title, func(t *testing.T) {
			t.Parallel()

			expect := test.IsPriorityLetter
			actual := IsPriorityLetter(test.Input)

			require.Equal(t, expect, actual,
				"IsPriorityLetter('%s') should return %v", test.Input, expect)
		})
	}
}

// ----------------------------------------------------------------------------
//  IsPriorityMark()
// ----------------------------------------------------------------------------

func TestIsPriorityMark(t *testing.T) {
	t.Parallel()

	// use common test data
	for index, test := range testdata.Priorities {
		title := fmt.Sprintf("Test #%d: %s", index+1, test.Title)

		t.Run(title, func(t *testing.T) {
			t.Parallel()

			expect := test.IsPriorityMark
			actual := IsPriorityMark(test.Input)

			require.Equal(t, expect, actual,
				"IsPriorityMark('%s') should return %v", test.Input, expect)
		})
	}
}

// ----------------------------------------------------------------------------
//  ContainsCtlChars()
// ----------------------------------------------------------------------------

func Test_ContainsCtlChars(t *testing.T) {
	t.Parallel()

	for _, test := range dataControlChars {
		allowedChars := test.allowedChars
		expect := test.containsCtlChar
		actual := ContainsCtlChars(test.input, allowedChars)

		require.Equal(t, expect, actual,
			"ContainsCtlChars('%q', %v) should return %v",
			test.input, allowedChars, expect)
	}
}

// ----------------------------------------------------------------------------
//  IsValidTaskLength()
// ----------------------------------------------------------------------------

func Test_IsValidTaskLength_invalid_length(t *testing.T) {
	t.Parallel()

	input := strings.Repeat("a", MaxTaskLength+1) // 64KB + 1 byte

	require.False(t, IsValidTaskLength(input),
		"it should return false for input length %d (> %d)",
		len(input), MaxTaskLength)
}
