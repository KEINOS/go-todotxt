package parse_test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/KEINOS/go-todotxt/todo/internal/parse"
)

// ----------------------------------------------------------------------------
//  FromTaskString()
// ----------------------------------------------------------------------------

//nolint:lll // allow long lines in examples
func ExampleFromTaskString() {
	taskStr := "x (A) 2016-05-20 2016-05-18 Thank Mom for the meatballs " +
		"+dinner @phone due:2016-05-25 # this is an inline-comment"

	parsed, err := parse.FromTaskString(taskStr)
	if err != nil {
		log.Fatalf("failed to parse task string: %v", err)
	}

	fmt.Println("Stringer:", parsed) // Equivalent to parsed.String()
	fmt.Println("Is completed:", parsed.IsDone())
	fmt.Println("Priority:", parsed.Priority())
	fmt.Println("Has priority:", parsed.HasPriority())
	fmt.Println("Date completed:", parsed.DateCompleted())
	fmt.Println("Date created:", parsed.DateCreated())
	fmt.Println("Contexts:", parsed.Contexts())
	fmt.Println("Projects:", parsed.Projects())
	fmt.Println("Description:", parsed.Description())
	fmt.Println("Key-Values:", parsed.KeyValues())

	fmt.Println("Comment:", parsed.Comment())
	fmt.Println("Is comment line:", parsed.IsCommentLine())
	fmt.Println("Has inline comment:", parsed.HasInlineComment())
	fmt.Println("Pos inline comment:", parsed.InlineCommentPos())
	// Unordered output:
	// Stringer: x (A) 2016-05-20 2016-05-18 Thank Mom for the meatballs +dinner @phone due:2016-05-25 # this is an inline-comment
	// Is completed: true
	// Priority: A
	// Has priority: true
	// Date completed: 2016-05-20
	// Date created: 2016-05-18
	// Contexts: [phone]
	// Projects: [dinner]
	// Description: Thank Mom for the meatballs +dinner @phone due:2016-05-25
	// Key-Values: [{due 2016-05-25}]
	// Comment: # this is an inline-comment
	// Is comment line: false
	// Has inline comment: true
	// Pos inline comment: 86
}

func ExampleFromTaskString_with_options() {
	// Task string with tab characters ('\t')
	taskStr := "x\t(A)\t2016-05-20\t2016-05-18\tThank Mom for the meatballs +dinner @phone due:2016-05-25 # comment"

	// Allow tab control characters in the task string
	parsed, err := parse.FromTaskString(taskStr,
		parse.WithAllowedCtrlChars([]rune{'\t'}),
	)
	if err != nil {
		log.Fatalf("failed to parse task string: %v", err)
	}

	fmt.Println("Is completed:", parsed.IsDone())
	fmt.Println("Priority:", parsed.Priority())
	fmt.Println("Date completed:", parsed.DateCompleted())
	fmt.Println("Date created:", parsed.DateCreated())
	// Output:
	// Is completed: true
	// Priority: A
	// Date completed: 2016-05-20
	// Date created: 2016-05-18
}

func ExampleFromTaskString_with_key_value_tags() {
	// Task string with key-value tags
	taskStr := "code review site url:https://example.com/#fragment due:2024-07-01 # important"

	parsed, err := parse.FromTaskString(taskStr)
	if err != nil {
		log.Fatalf("failed to parse task string: %v", err)
	}

	fmt.Println("Description:", parsed.Description())
	fmt.Println("Comment:", parsed.Comment())

	for _, kv := range parsed.KeyValues() {
		fmt.Printf("Key-Value: %s = %s\n", kv.Key, kv.Value)
	}
	// Unordered output:
	// Description: code review site url:https://example.com/#fragment due:2024-07-01
	// Comment: # important
	// Key-Value: url = https://example.com/#fragment
	// Key-Value: due = 2024-07-01
}

// ----------------------------------------------------------------------------
//  Parsed.Parse()
// ----------------------------------------------------------------------------

func ExampleParsed_Parse() {
	// Initial task string
	taskStr := "Buy milk @store"

	parsed, err := parse.FromTaskString(taskStr)
	if err != nil {
		log.Fatalf("failed to parse task string: %v", err)
	}

	fmt.Println("Original description:", parsed.Description())

	// Update the task string including tab control character
	updatedTaskStr := "(B)\tBuy milk @store # updated"

	err = parsed.Parse(updatedTaskStr,
		// Allow tab control character in the updated task string
		parse.WithAllowedCtrlChars([]rune{'\t'}),
	)
	if err != nil {
		log.Fatalf("failed to re-parse updated task string: %v", err)
	}

	fmt.Println("Updated priority:", parsed.Priority())
	fmt.Println("Updated description:", parsed.Description())
	fmt.Println("Updated comment:", parsed.Comment())
	fmt.Printf("Allowed control characters: %q\n", parsed.AllowedCtrlChars())
	// Output:
	// Original description: Buy milk @store
	// Updated priority: B
	// Updated description: Buy milk @store
	// Updated comment: # updated
	// Allowed control characters: ['\t']
}

// ----------------------------------------------------------------------------
//  Parsed.Segments
// ----------------------------------------------------------------------------

//nolint:cyclop // complexity 11 is acceptable for an example
func ExampleParsed_segments_type_listing() {
	taskStr := "x (A) 2016-05-20 Thank Mom for the meatballs +dinner @phone due:2016-05-25 # comment"

	parsed, err := parse.FromTaskString(taskStr)
	if err != nil {
		log.Fatalf("failed to parse task string: %v", err)
	}

	for index, seg := range parsed.Segments {
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
			isSegType = "comment"
		case seg.IsPlainText():
			isSegType = "plain text"
		}

		fmt.Printf("%02d: %-15s (%s)\n", index+1, seg, isSegType)
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
	// 12: #               (comment)
	// 13: comment         (plain text)
}

// ----------------------------------------------------------------------------
//  Parsed.Components()
// ----------------------------------------------------------------------------

func ExampleParsed_Components() {
	// Parse a completed task with all components
	taskStr := "x (A) 2016-05-20 2016-04-30 measure space for +chapelShelving @chapel due:2016-05-30 # comment"

	parsed, err := parse.FromTaskString(taskStr)
	if err != nil {
		log.Fatalf("failed to parse task: %v", err)
	}

	// Get all components at once
	comp := parsed.Components()

	// Access individual fields
	fmt.Printf("Done: %v\n", comp.IsDone)
	fmt.Printf("Priority: %s\n", comp.Priority)
	fmt.Printf("Completed: %s\n", comp.DateCompleted)
	fmt.Printf("Created: %s\n", comp.DateCreated)
	fmt.Printf("Description: %s\n", comp.Description)
	fmt.Printf("Contexts: %v\n", comp.Contexts)
	fmt.Printf("Projects: %v\n", comp.Projects)
	fmt.Printf("Tags: %v\n", comp.KeyValues)
	fmt.Printf("Comment: %s\n", comp.Comment)

	// Unordered output:
	// Done: true
	// Priority: A
	// Completed: 2016-05-20
	// Created: 2016-04-30
	// Description: measure space for +chapelShelving @chapel due:2016-05-30
	// Contexts: [chapel]
	// Projects: [chapelShelving]
	// Tags: [{due 2016-05-30}]
	// Comment: # comment
}

func ExampleParsed_Components_simple() {
	// Parse a simple task with minimal components
	taskStr := "Buy milk @store"

	parsed, err := parse.FromTaskString(taskStr)
	if err != nil {
		log.Fatalf("failed to parse task: %v", err)
	}

	comp := parsed.Components()

	fmt.Printf("Done: %v\n", comp.IsDone)
	fmt.Printf("Description: %s\n", comp.Description)
	fmt.Printf("Contexts: %v\n", comp.Contexts)

	// Output:
	// Done: false
	// Description: Buy milk @store
	// Contexts: [store]
}

func ExampleParsed_Components_commentLine() {
	// Parse a comment line
	taskStr := "# This is a comment"

	parsed, err := parse.FromTaskString(taskStr)
	if err != nil {
		log.Fatalf("failed to parse task: %v", err)
	}

	comp := parsed.Components()

	fmt.Printf("Comment line: %v\n", comp.IsCommentLine)
	fmt.Printf("Comment: %s\n", comp.Comment)
	fmt.Printf("Description: %s\n", comp.Description)

	// Output:
	// Comment line: true
	// Comment: # This is a comment
	// Description:
}

// ----------------------------------------------------------------------------
//  Components as JSON
// ----------------------------------------------------------------------------

func ExampleComponents_json() {
	// Create a Components struct manually for demonstration
	comp := parse.Components{
		IsDone:           true,
		IsCommentLine:    false,
		HasInlineComment: true,
		PosInlineComment: 68,
		Priority:         "A",
		DateCompleted:    "2016-05-20",
		DateCreated:      "2016-04-30",
		Description:      "measure space for +chapelShelving @chapel due:2016-05-30",
		Comment:          "# comment",
		Contexts:         []string{"chapel"},
		Projects:         []string{"chapelShelving"},
		KeyValues:        []parse.KeyValue{{Key: "due", Value: "2016-05-30"}},
	}

	// Marshal to JSON with indentation
	jsonBytes, err := json.MarshalIndent(comp, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal to JSON: %v", err)
	}

	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "priority": "A",
	//   "dateCompleted": "2016-05-20",
	//   "dateCreated": "2016-04-30",
	//   "description": "measure space for +chapelShelving @chapel due:2016-05-30",
	//   "comment": "# comment",
	//   "contexts": [
	//     "chapel"
	//   ],
	//   "projects": [
	//     "chapelShelving"
	//   ],
	//   "keyValues": [
	//     {
	//       "key": "due",
	//       "value": "2016-05-30"
	//     }
	//   ],
	//   "posInlineComment": 68,
	//   "isDone": true,
	//   "isCommentLine": false,
	//   "hasInlineComment": true
	// }
}

func ExampleComponents_json_minimal() {
	// Minimal task with only required fields
	//
	//nolint:exhaustruct // allow missing fields for example
	comp := parse.Components{
		IsDone:        false,
		IsCommentLine: false,
		Description:   "Buy milk",
	}

	// Marshal to JSON with indentation
	jsonBytes, err := json.MarshalIndent(comp, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal to JSON: %v", err)
	}

	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "description": "Buy milk",
	//   "isDone": false,
	//   "isCommentLine": false
	// }
}
