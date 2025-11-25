package spec_test

import (
	"fmt"

	"github.com/KEINOS/go-todotxt/todo/internal/spec"
)

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
