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
