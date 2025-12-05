package spec_test

import (
	"fmt"

	"github.com/KEINOS/go-todotxt/todo/internal/spec"
)

func Example() {
	fmt.Printf("Project Prefix: \"%s\"\n", spec.PrefixProject)
	fmt.Printf("Context Prefix: \"%s\"\n", spec.PrefixContext)
	fmt.Printf("Comment Prefix: \"%s\"\n", spec.PrefixComment)
	fmt.Printf("Done Marker: \"%s\"\n", spec.MarkerDone)
	fmt.Printf("Key-Value Separator: \"%s\"\n", spec.SepKeyValue)
	fmt.Printf("Segment Delimiter: \"%s\"\n", spec.DelimSegments)
	// Output:
	// Project Prefix: "+"
	// Context Prefix: "@"
	// Comment Prefix: "#"
	// Done Marker: "x"
	// Key-Value Separator: ":"
	// Segment Delimiter: " "
}

// ============================================================================
//  Method Examples
// ============================================================================

func ExampleMark_Byte() {
	fmt.Printf("Project Prefix: %x\n", spec.PrefixProject.Byte())
	fmt.Printf("Context Prefix: %x\n", spec.PrefixContext.Byte())
	fmt.Printf("Comment Prefix: %x\n", spec.PrefixComment.Byte())
	fmt.Printf("Done Marker: %x\n", spec.MarkerDone.Byte())
	fmt.Printf("Key-Value Separator: %x\n", spec.SepKeyValue.Byte())
	fmt.Printf("Segment Delimiter: %x\n", spec.DelimSegments.Byte())
	// Output:
	// Project Prefix: 2b
	// Context Prefix: 40
	// Comment Prefix: 23
	// Done Marker: 78
	// Key-Value Separator: 3a
	// Segment Delimiter: 20
}

func ExampleMark_Rune() {
	fmt.Printf("Project Prefix: %q\n", spec.PrefixProject.Rune())
	fmt.Printf("Context Prefix: %q\n", spec.PrefixContext.Rune())
	fmt.Printf("Comment Prefix: %q\n", spec.PrefixComment.Rune())
	fmt.Printf("Done Marker: %q\n", spec.MarkerDone.Rune())
	fmt.Printf("Key-Value Separator: %q\n", spec.SepKeyValue.Rune())
	fmt.Printf("Segment Delimiter: %q\n", spec.DelimSegments.Rune())
	// Output:
	// Project Prefix: '+'
	// Context Prefix: '@'
	// Comment Prefix: '#'
	// Done Marker: 'x'
	// Key-Value Separator: ':'
	// Segment Delimiter: ' '
}

func ExampleMark_String() {
	fmt.Printf("Project Prefix: %q\n", spec.PrefixProject.String())
	fmt.Printf("Context Prefix: %q\n", spec.PrefixContext.String())
	fmt.Printf("Comment Prefix: %q\n", spec.PrefixComment.String())
	fmt.Printf("Done Marker: %q\n", spec.MarkerDone.String())
	fmt.Printf("Key-Value Separator: %q\n", spec.SepKeyValue.String())
	fmt.Printf("Segment Delimiter: %q\n", spec.DelimSegments.String())
	// Output:
	// Project Prefix: "+"
	// Context Prefix: "@"
	// Comment Prefix: "#"
	// Done Marker: "x"
	// Key-Value Separator: ":"
	// Segment Delimiter: " "
}

func ExampleMark_IsPrintable() {
	// All defined spec markers are printable ASCII
	fmt.Printf("Project Prefix (+): %v\n", spec.PrefixProject.IsPrintable())
	fmt.Printf("Context Prefix (@): %v\n", spec.PrefixContext.IsPrintable())
	fmt.Printf("Segment Delimiter (space): %v\n", spec.DelimSegments.IsPrintable())

	// Control characters are not printable
	tab := spec.Mark('\t')
	nul := spec.Mark(0x00)
	del := spec.Mark(0x7F)

	fmt.Printf("Tab (0x09): %v\n", tab.IsPrintable())
	fmt.Printf("NUL (0x00): %v\n", nul.IsPrintable())
	fmt.Printf("DEL (0x7F): %v\n", del.IsPrintable())
	// Output:
	// Project Prefix (+): true
	// Context Prefix (@): true
	// Segment Delimiter (space): true
	// Tab (0x09): false
	// NUL (0x00): false
	// DEL (0x7F): false
}

// ExampleMark_validation demonstrates that all Mark methods validate printability.
// Non-printable characters return zero/empty values for safety.
func ExampleMark_validation() {
	// Valid printable mark
	validMark := spec.Mark('A')
	fmt.Printf("Valid 'A' - Byte: 0x%02X, Rune: %q, String: %q\n",
		validMark.Byte(), validMark.Rune(), validMark.String())

	// Invalid non-printable mark (control character)
	invalidMark := spec.Mark(0x1F)
	fmt.Printf("Invalid 0x1F - Byte: 0x%02X, Rune: %q, String: %q\n",
		invalidMark.Byte(), invalidMark.Rune(), invalidMark.String())

	// Check before use pattern
	mark := spec.Mark(0x80)
	if mark.IsPrintable() {
		fmt.Printf("Mark: %s\n", mark.String())
	} else {
		fmt.Println("Mark is not printable, skipping")
	}
	// Output:
	// Valid 'A' - Byte: 0x41, Rune: 'A', String: "A"
	// Invalid 0x1F - Byte: 0x00, Rune: '\x00', String: ""
	// Mark is not printable, skipping
}

// ============================================================================
//  Checker Function Examples
// ============================================================================

// ----------------------------------------------------------------------------
//  IsDate()
// ----------------------------------------------------------------------------

func ExampleIsDate() {
	sampleDate := "2024-12-31"    // valid date
	notADate := "milk"            // not a date
	invalidMonth := "2024-13-01"  // valid format but bad month
	invalidFormat := "31-12-2024" // bad format (not YYYY-MM-DD)

	// validation disabled
	validate := false

	fmt.Printf("IsDate(%q, %v): %v\n", sampleDate, validate,
		spec.IsDate(sampleDate, validate))
	fmt.Printf("IsDate(%q, %v): %v\n", notADate, validate,
		spec.IsDate(notADate, validate))
	fmt.Printf("IsDate(%q, %v): %v\n", invalidMonth, validate,
		spec.IsDate(invalidMonth, validate))
	fmt.Printf("IsDate(%q, %v): %v\n", invalidFormat, validate,
		spec.IsDate(invalidFormat, validate))

	// validation enabled
	validate = true

	fmt.Printf("IsDate(%q, %v): %v\n", sampleDate, validate,
		spec.IsDate(sampleDate, validate))
	fmt.Printf("IsDate(%q, %v): %v\n", notADate, validate,
		spec.IsDate(notADate, validate))
	fmt.Printf("IsDate(%q, %v): %v\n", invalidMonth, validate,
		spec.IsDate(invalidMonth, validate))
	fmt.Printf("IsDate(%q, %v): %v\n", invalidFormat, validate,
		spec.IsDate(invalidFormat, validate))
	// Output:
	// IsDate("2024-12-31", false): true
	// IsDate("milk", false): false
	// IsDate("2024-13-01", false): true
	// IsDate("31-12-2024", false): false
	// IsDate("2024-12-31", true): true
	// IsDate("milk", true): false
	// IsDate("2024-13-01", true): false
	// IsDate("31-12-2024", true): false
}

// ----------------------------------------------------------------------------
//  ContainsInlineComment()
// ----------------------------------------------------------------------------

func ExampleContainsInlineComment() {
	for index, task := range []struct {
		text        string
		expectIndex int
		expectFound bool
	}{
		// True cases (treat as inline comment)
		{
			text:        "Task w/ +Project @Context # This is an inline comment",
			expectIndex: 34,
			expectFound: true,
		},
		{
			text:        "  # Leading spaces before comment",
			expectIndex: 2,
			expectFound: true,
		},
		{
			text: "Task w/\t# inline comment after tab",
			// The index is byte-based; tab is a single byte
			expectIndex: 8,
			expectFound: true,
		},
		// False cases (no inline comment)
		{text: "Task w/ +Project @Context with no comment",
			expectIndex: -1,
			expectFound: false,
		},
		{
			text:        "# This entire line is a comment",
			expectIndex: -1,
			expectFound: false,
		},
		{
			text:        "Task w/ key:value of url:http://example.com/page#fragment",
			expectIndex: -1,
			expectFound: false,
		},
	} {
		pos, found := spec.ContainsInlineComment(task.text)

		fmt.Printf("#%d: Input: %q, Pos: %d, Found: %v\n",
			index+1, task.text, pos, found)
	}
	// Output:
	// #1: Input: "Task w/ +Project @Context # This is an inline comment", Pos: 26, Found: true
	// #2: Input: "  # Leading spaces before comment", Pos: 2, Found: true
	// #3: Input: "Task w/\t# inline comment after tab", Pos: 8, Found: true
	// #4: Input: "Task w/ +Project @Context with no comment", Pos: -1, Found: false
	// #5: Input: "# This entire line is a comment", Pos: -1, Found: false
	// #6: Input: "Task w/ key:value of url:http://example.com/page#fragment", Pos: -1, Found: false
}
