package segment_test

import (
	"fmt"

	"github.com/KEINOS/go-todotxt/todo/internal/segment"
)

// ----------------------------------------------------------------------------
//  Type: Type
// ----------------------------------------------------------------------------

func ExampleType() {
	for _, segType := range []segment.Type{
		segment.Undefined,
		segment.MarkCompletion,
		segment.MarkPriority,
		segment.Date,
		segment.TagProject,
		segment.TagContext,
		segment.TagKeyValue,
		segment.MarkComment,
		segment.PlainText,
	} {
		fmt.Printf("%012b\n", segType)
	}
	// Output:
	// 000000000000
	// 000000000010
	// 000000000100
	// 000000001000
	// 000000010000
	// 000000100000
	// 000001000000
	// 000010000000
	// 000100000000
}

// ----------------------------------------------------------------------------
//	Type: Segment
// ----------------------------------------------------------------------------

func ExampleSegment() {
	// Example of parsed/chanked segments from a todo.txt task string
	segments := []segment.Segment{
		"x", "(A)", "2016-05-20", "Thank", "mom", "for", "the", "meatballs",
		"+dinner", "@phone", "due:2016-05-25", "#", "comment",
	}

	for index, seg := range segments {
		fmt.Printf("%02d: %-15s (%s)\n", index+1, seg.String(), seg.Type().String())
	}
	// Output:
	// 01: x               (completion mark)
	// 02: (A)             (priority mark)
	// 03: 2016-05-20      (date)
	// 04: Thank           (plain text)
	// 05: mom             (plain text)
	// 06: for             (plain text)
	// 07: the             (plain text)
	// 08: meatballs       (plain text)
	// 09: +dinner         (project tag)
	// 10: @phone          (context tag)
	// 11: due:2016-05-25  (key-value pair tag)
	// 12: #               (comment prefix)
	// 13: comment         (plain text)
}

func ExampleSegment_Is() {
	segExample := segment.Segment("2023-12-31")

	fmt.Println(segExample.Is(segment.MarkCompletion)) // false
	fmt.Println(segExample.Is(segment.MarkPriority))   // false
	fmt.Println(segExample.Is(segment.Date))           // true
	fmt.Println(segExample.Is(segment.TagProject))     // false
	fmt.Println(segExample.Is(segment.TagContext))     // false
	fmt.Println(segExample.Is(segment.TagKeyValue))    // false
	fmt.Println(segExample.Is(segment.MarkComment))    // false
	fmt.Println(segExample.Is(segment.PlainText))      // false
	fmt.Println(segExample.Is(segment.Undefined))      // false
	// Output:
	// false
	// false
	// true
	// false
	// false
	// false
	// false
	// false
	// false
}

func ExampleSegment_IsDate() {
	seg1 := segment.Segment("2024-06-15")
	seg2 := segment.Segment("2024-02-30") // invalid date but valid format
	seg3 := segment.Segment("not-a-date")

	fmt.Printf("Is '%s' a date? --> %v\n", seg1, seg1.IsDate())
	fmt.Printf("Is '%s' a date? --> %v\n", seg2, seg2.IsDate())
	fmt.Printf("Is '%s' a date? --> %v\n", seg3, seg3.IsDate())
	// Output:
	// Is '2024-06-15' a date? --> true
	// Is '2024-02-30' a date? --> true
	// Is 'not-a-date' a date? --> false
}

func ExampleSegment_IsDateValid() {
	seg1 := segment.Segment("2024-02-29") // valid leap year date
	seg2 := segment.Segment("2024-06-31") // invalid date
	seg3 := segment.Segment("not-a-date") // invalid format

	fmt.Printf("Is '%s' a valid date? --> %v\n", seg1, seg1.IsDateValid())
	fmt.Printf("Is '%s' a valid date? --> %v\n", seg2, seg2.IsDateValid())
	fmt.Printf("Is '%s' a valid date? --> %v\n", seg3, seg3.IsDateValid())
	// Output:
	// Is '2024-02-29' a valid date? --> true
	// Is '2024-06-31' a valid date? --> false
	// Is 'not-a-date' a valid date? --> false
}

// ----------------------------------------------------------------------------
//  Type: Segments
// ----------------------------------------------------------------------------

func ExampleSegments_Types() {
	segments := segment.Segments{
		segment.Segment("x"),
		segment.Segment("(A)"),
		segment.Segment("2023-12-31"),
		segment.Segment("Call"),
		segment.Segment("mom"),
		segment.Segment("+family"),
		// Omitted to show absence
		// segment.Segment("@phone"),
		// segment.Segment("due:2024-01-05"),
	}

	types := segments.Types()
	fmt.Printf("Combined Types: 0b%012b (0x%03X)\n", types, uint16(types))

	// Check for specific types using bitwise AND
	hasCompletion := types&segment.MarkCompletion != 0
	hasPriority := types&segment.MarkPriority != 0
	hasDate := types&segment.Date != 0
	hasProject := types&segment.TagProject != 0
	hasContext := types&segment.TagContext != 0
	hasKeyValue := types&segment.TagKeyValue != 0
	hasComment := types&segment.MarkComment != 0

	fmt.Printf("Has completion mark: %v\n", hasCompletion)
	fmt.Printf("Has priority mark: %v\n", hasPriority)
	fmt.Printf("Has date: %v\n", hasDate)
	fmt.Printf("Has project tag: %v\n", hasProject)
	fmt.Printf("Has context tag: %v\n", hasContext)
	fmt.Printf("Has key-value tag: %v\n", hasKeyValue)
	fmt.Printf("Has comment: %v\n", hasComment)
	// Output:
	// Combined Types: 0b000100011110 (0x11E)
	// Has completion mark: true
	// Has priority mark: true
	// Has date: true
	// Has project tag: true
	// Has context tag: false
	// Has key-value tag: false
	// Has comment: false
}
