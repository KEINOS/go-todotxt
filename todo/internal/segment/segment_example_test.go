package segment_test

import (
	"fmt"

	"github.com/KEINOS/go-todotxt/todo/internal/segment"
)

func ExampleSegment() {
	// Example of parsed/chanked segments from a todo.txt task string
	segments := []segment.Segment{
		"x", "(A)", "2016-05-20", "Thank", "Mom", "for", "the", "meatballs",
		"+dinner", "@phone", "due:2016-05-25", "#", "comment",
	}

	for index, seg := range segments {
		isSegType := "unknown segment type"

		switch {
		case seg.IsMarkCompletion():
			isSegType = "completion mark"
		case seg.IsMarkPriority():
			isSegType = "priority"
		case seg.IsDate():
			isSegType = "date"
		case seg.IsTagProject():
			isSegType = "project tag"
		case seg.IsTagContext():
			isSegType = "context tag"
		case seg.IsKeyValue():
			isSegType = "key-value pair"
		case seg.IsComment():
			isSegType = "comment prefix"
		case seg.IsPlainText():
			isSegType = "plain text"
		}

		fmt.Printf("%02d: %-15s (%s)\n", index+1, seg.String(), isSegType)
	}
	// Output:
	// 01: x               (completion mark)
	// 02: (A)             (priority)
	// 03: 2016-05-20      (date)
	// 04: Thank           (plain text)
	// 05: Mom             (plain text)
	// 06: for             (plain text)
	// 07: the             (plain text)
	// 08: meatballs       (plain text)
	// 09: +dinner         (project tag)
	// 10: @phone          (context tag)
	// 11: due:2016-05-25  (key-value pair)
	// 12: #               (comment prefix)
	// 13: comment         (plain text)
}
