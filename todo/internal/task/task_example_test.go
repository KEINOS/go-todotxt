package task_test

import (
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/KEINOS/go-todotxt/todo/internal/spec"
	"github.com/KEINOS/go-todotxt/todo/internal/task"
)

// ============================================================================
//  Basic Examples
// ============================================================================

// Retrieval methods are inherited from embedded parse.Parsed (read-only).
//
//nolint:lll // acceptable in examples
func Example() {
	taskText := "x (A) 2016-05-20 2016-05-18 this is a completed task with " +
		"+project and @context due:2016-05-25 # with an inline-comment"

	// Create a task from text.
	tsk, err := task.New(taskText)
	if err != nil {
		panic(err)
	}

	// Retrieve the task info via parse.Parsed.
	fmt.Println("Priority:", tsk.Priority())
	fmt.Println("Date Completed:", tsk.DateCompleted())
	fmt.Println("Date Created:", tsk.DateCreated())
	fmt.Println("Description:", tsk.Description())
	fmt.Println("Projects:", tsk.Projects())
	fmt.Println("Contexts:", tsk.Contexts())
	fmt.Println("Key-Values:", tsk.KeyValues())

	fmt.Println("Is Completed:", tsk.IsDone())
	fmt.Println("Has Inline Comment:", tsk.HasInlineComment())
	fmt.Println("Is Comment Line:", tsk.IsCommentLine())
	fmt.Println("Comment:", tsk.Comment())

	fmt.Println("Task String:", tsk) // equivalent to tsk.String()
	fmt.Println("Segments:", tsk.Segments)
	// Output:
	// Priority: A
	// Date Completed: 2016-05-20
	// Date Created: 2016-05-18
	// Description: this is a completed task with +project and @context due:2016-05-25
	// Projects: [project]
	// Contexts: [context]
	// Key-Values: [{due 2016-05-25}]
	// Is Completed: true
	// Has Inline Comment: true
	// Is Comment Line: false
	// Comment: # with an inline-comment
	// Task String: x (A) 2016-05-20 2016-05-18 this is a completed task with +project and @context due:2016-05-25 # with an inline-comment
	// Segments: [x (A) 2016-05-20 2016-05-18 this is a completed task with +project and @context due:2016-05-25 # with an inline-comment]
}

func Example_output_as_JSON() {
	taskText := "x (A) 2025-05-20 2025-05-18 this is a completed task with " +
		"+project and @context due:2025-05-25 # with an inline-comment"

	// Create a task from text.
	tsk, err := task.New(taskText)
	if err != nil {
		panic(err)
	}

	// Retrieve the task components as JSON.
	jsonBytes, err := json.MarshalIndent(tsk.Components(), "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "priority": "A",
	//   "dateCompleted": "2025-05-20",
	//   "dateCreated": "2025-05-18",
	//   "description": "this is a completed task with +project and @context due:2025-05-25",
	//   "comment": "# with an inline-comment",
	//   "contexts": [
	//     "context"
	//   ],
	//   "projects": [
	//     "project"
	//   ],
	//   "keyValues": [
	//     {
	//       "key": "due",
	//       "value": "2025-05-25"
	//     }
	//   ],
	//   "posInlineComment": 95,
	//   "isDone": true,
	//   "isCommentLine": false,
	//   "hasInlineComment": true
	// }
}

// Creating a task with initial options.
func Example_create_new_task_with_options() {
	taskText := "buy mango and banana"

	// Create a task with setter options
	tsk, err := task.New(taskText,
		task.WithCompleted("2024-12-01"), // use `task.WithCompleted()` for current date
		task.WithPriority("B"),
		task.WithContext("unripe"),
		task.WithProject("shopping"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Task String:", tsk)
	fmt.Println("Is Completed:", tsk.IsDone())
	fmt.Println("Priority:", tsk.Priority())
	fmt.Println("Date Completed:", tsk.DateCompleted())
	fmt.Println("Contexts:", tsk.Contexts())
	fmt.Println("Projects:", tsk.Projects())
	// Output:
	// Task String: x (B) 2024-12-01 buy mango and banana @unripe +shopping
	// Is Completed: true
	// Priority: B
	// Date Completed: 2024-12-01
	// Contexts: [unripe]
	// Projects: [shopping]
}

func Example_use_user_custom_option() {
	tsk, err := task.New("take medicine @daily_dose")
	if err != nil {
		panic(err)
	}

	WithCommentOut := func(tsk2 *task.Task) error {
		if tsk2.IsCommentLine() {
			return nil
		}

		tsk2.SetText(
			spec.PrefixComment.String() + spec.DelimSegments.String() +
				tsk2.String(),
		)

		return nil
	}

	if slices.Contains(tsk.Contexts(), "daily_dose") {
		err = tsk.Apply(
			WithCommentOut,
		)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("task:", tsk.String())
	// Output:
	// task: # take medicine @daily_dose
}

// ============================================================================
//  Constructor Examples
// ============================================================================

// Creating a task and modifying it with setter methods.
// Call Apply() after modifications to finalize changes.
func ExampleNew_create_new_task_then_manage_components() {
	taskText := "buy Mac mini"

	panicOnErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	// Create a task
	tsk, err := task.New(taskText)
	panicOnErr(err)

	// Set options
	err = tsk.SetPriority("A")
	panicOnErr(err)

	err = tsk.AddContext("M4")
	panicOnErr(err)

	err = tsk.AddProject("shopping")
	panicOnErr(err)

	// Apply changes (**important**: call Update to reparse the task)
	err = tsk.Apply()
	panicOnErr(err)

	fmt.Println("Current task String:", tsk)
	// Output:
	// Current task String: (A) buy Mac mini @M4 +shopping
}

// Control characters (including tabs) are disallowed by default.
// Use WithAllowedCtrlChars to permit them.
//
//nolint:lll // acceptable in examples
func ExampleNew_allow_tabs_in_task_text_to_create_with_parse_option() {
	taskText := "x\t(A)\t2016-05-20 2016-05-18 this is a completed task with " +
		"+project and @context due:2016-05-25 # withan inline-comment"

	// Create and retrieve a task
	tsk, err := task.New(taskText, task.WithAllowedCtrlChars([]rune{'\t'}))
	if err != nil {
		panic(err)
	}

	fmt.Println("Task String:", tsk)
	fmt.Println("Is Completed:", tsk.IsDone())
	fmt.Println("Priority:", tsk.Priority())
	fmt.Println("Date Completed:", tsk.DateCompleted())
	fmt.Println("Date Created:", tsk.DateCreated())
	fmt.Println("Description:", tsk.Description())
	fmt.Println("Projects:", tsk.Projects())
	fmt.Println("Contexts:", tsk.Contexts())
	fmt.Println("Key-Values:", tsk.KeyValues())
	fmt.Println("Comment:", tsk.Comment())
	// Output:
	// Task String: x	(A)	2016-05-20 2016-05-18 this is a completed task with +project and @context due:2016-05-25 # withan inline-comment
	// Is Completed: true
	// Priority: A
	// Date Completed: 2016-05-20
	// Date Created: 2016-05-18
	// Description: this is a completed task with +project and @context due:2016-05-25
	// Projects: [project]
	// Contexts: [context]
	// Key-Values: [{due 2016-05-25}]
	// Comment: # withan inline-comment
}

// ============================================================================
//  Task Method Examples
// ============================================================================

// ----------------------------------------------------------------------------
//  Task.IsDirty
// ----------------------------------------------------------------------------

func ExampleTask_IsDirty() {
	taskText := "get dirty"

	tsk, err := task.New(taskText)
	if err != nil {
		panic(err)
	}

	fmt.Println("IsDirty (initial):", tsk.IsDirty())

	// Set priority
	err = tsk.SetPriority("A")
	if err != nil {
		panic(err)
	}

	fmt.Println("IsDirty (after SetPriority):", tsk.IsDirty())

	// Apply changes
	err = tsk.Apply()
	if err != nil {
		panic(err)
	}

	fmt.Println("IsDirty (after Apply):", tsk.IsDirty())

	// Force dirty
	tsk.MarkDirty()

	fmt.Println("IsDirty (after MarkDirty):", tsk.IsDirty())

	// Output:
	// IsDirty (initial): false
	// IsDirty (after SetPriority): true
	// IsDirty (after Apply): false
	// IsDirty (after MarkDirty): true
}

// ----------------------------------------------------------------------------
//  Task.MarkDone / Complete / CompleteWithDate
// ----------------------------------------------------------------------------

// MarkDone is an alias for Complete() using current date.
func ExampleTask_MarkDone() {
	tsk, err := task.New("buy milk")
	if err != nil {
		panic(err)
	}

	err = tsk.MarkDone() // use current date as completion date
	if err != nil {
		panic(err)
	}

	err = tsk.Apply()
	if err != nil {
		panic(err)
	}

	expectedDate := time.Now().Format(spec.DateFormat) // YYYY-MM-DD
	actualDate := tsk.DateCompleted()

	fmt.Println("Is completed:", tsk.IsDone())
	fmt.Println("Description:", tsk.Description())
	fmt.Println("Completed date matches today:", actualDate == expectedDate)
	// Output:
	// Is completed: true
	// Description: buy milk
	// Completed date matches today: true
}

func ExampleTask_Complete() {
	tsk, err := task.New("buy milk")
	if err != nil {
		panic(err)
	}

	// Complete uses current date as completion date.
	// Equivalent to:
	//   compDate := time.Now().Format(spec.DateFormat)
	//   err = tsk.CompleteWithDate(compDate)
	err = tsk.Complete()
	if err != nil {
		panic(err)
	}

	err = tsk.Apply()
	if err != nil {
		panic(err)
	}

	expectDate := time.Now().Format(spec.DateFormat) // YYYY-MM-DD
	actualDate := tsk.DateCompleted()

	fmt.Println("Is completed:", tsk.IsDone())
	fmt.Println("Description:", tsk.Description())
	fmt.Println("Completed date matches today:", actualDate == expectDate)
	// Output:
	// Is completed: true
	// Description: buy milk
	// Completed date matches today: true
}

func ExampleTask_CompleteWithDate_specific() {
	tsk, err := task.New("buy milk")
	if err != nil {
		panic(err)
	}

	err = tsk.CompleteWithDate("2024-01-15")
	if err != nil {
		panic(err)
	}

	err = tsk.Apply()
	if err != nil {
		panic(err)
	}

	fmt.Println("Is completed:", tsk.IsDone())
	fmt.Println("Description:", tsk.Description())
	fmt.Println("Completed date:", tsk.DateCompleted())
	// Output:
	// Is completed: true
	// Description: buy milk
	// Completed date: 2024-01-15
}

// This example shows completing a task without a date.
// It is non-standard but common in case users just want to mark tasks as done.
//
// E.g., "buy milk" --> "x buy milk".
func ExampleTask_CompleteWithDate_noDate() {
	tsk, err := task.New("buy milk")
	if err != nil {
		panic(err)
	}

	err = tsk.CompleteWithDate("")
	if err != nil {
		panic(err)
	}

	err = tsk.Apply()
	if err != nil {
		panic(err)
	}

	fmt.Println("Task String:", tsk)
	// Output:
	// Task String: x buy milk
}

// ----------------------------------------------------------------------------
//  Task.Reopen
// ----------------------------------------------------------------------------

func ExampleTask_Reopen() {
	// Completed task with completion date and creation date
	taskText := "x (A) 2024-01-15 2024-01-01 buy milk @grocery +shopping"

	tsk, err := task.New(taskText)
	if err != nil {
		panic(err)
	}

	err = tsk.Reopen()
	if err != nil {
		panic(err)
	}

	err = tsk.Apply()
	if err != nil {
		panic(err)
	}

	fmt.Println("Is completed:", tsk.IsDone())
	fmt.Println("Task String:", tsk.String())
	// Output:
	// Is completed: false
	// Task String: (A) 2024-01-01 buy milk @grocery +shopping
}

// ----------------------------------------------------------------------------
//  Task.SetPriority
// ----------------------------------------------------------------------------

func ExampleTask_SetPriority() {
	taskText := "buy milk @grocery +shopping"

	tsk, err := task.New(taskText)
	if err != nil {
		panic(err)
	}

	// Set priority
	err = tsk.SetPriority("A")
	if err != nil {
		panic(err)
	}

	// Apply changes
	err = tsk.Apply()
	if err != nil {
		panic(err)
	}

	fmt.Println(tsk)
	fmt.Println("Priority:", tsk.Priority())
	// Output:
	// (A) buy milk @grocery +shopping
	// Priority: A
}

func ExampleTask_SetPriority_updateExisting() {
	taskText := "(B) buy milk @grocery +shopping"

	tsk, err := task.New(taskText)
	if err != nil {
		panic(err)
	}

	fmt.Println("Before:", tsk)
	fmt.Println("Priority:", tsk.Priority())

	// Update priority from B to A
	err = tsk.SetPriority("A")
	if err != nil {
		panic(err)
	}

	err = tsk.Apply()
	if err != nil {
		panic(err)
	}

	fmt.Println("After:", tsk)
	fmt.Println("Priority:", tsk.Priority())
	// Output:
	// Before: (B) buy milk @grocery +shopping
	// Priority: B
	// After: (A) buy milk @grocery +shopping
	// Priority: A
}

func ExampleTask_SetPriority_withCompletedTask() {
	taskText := "x 2024-01-15 buy milk"

	tsk, err := task.New(taskText)
	if err != nil {
		panic(err)
	}

	// Set priority on completed task (priority goes after completion mark)
	err = tsk.SetPriority("A")
	if err != nil {
		panic(err)
	}

	err = tsk.Apply()
	if err != nil {
		panic(err)
	}

	fmt.Println(tsk)
	fmt.Println("Priority:", tsk.Priority())
	// Output:
	// x (A) 2024-01-15 buy milk
	// Priority: A
}

// ----------------------------------------------------------------------------
//  Task.RemovePriority
// ----------------------------------------------------------------------------

func ExampleTask_RemovePriority() {
	taskText := "(A) buy milk @grocery +shopping"

	tsk, err := task.New(taskText)
	if err != nil {
		panic(err)
	}

	fmt.Println("Before:", tsk)
	fmt.Println("Priority:", tsk.Priority())

	// Remove priority
	err = tsk.RemovePriority()
	if err != nil {
		panic(err)
	}

	err = tsk.Apply()
	if err != nil {
		panic(err)
	}

	fmt.Println("After:", tsk)
	fmt.Println("Priority:", tsk.Priority())
	// Output:
	// Before: (A) buy milk @grocery +shopping
	// Priority: A
	// After: buy milk @grocery +shopping
	// Priority:
}

// ----------------------------------------------------------------------------
//  Task.AppendSegment
// ----------------------------------------------------------------------------

func ExampleTask_AppendSegment() {
	tsk, err := task.New("buy milk")
	if err != nil {
		panic(err)
	}

	err = tsk.AppendSegment("@home")
	if err != nil {
		panic(err)
	}

	err = tsk.Apply()
	if err != nil {
		panic(err)
	}

	fmt.Println(tsk)
	// Output: buy milk @home
}

// ----------------------------------------------------------------------------
//  Task.RemoveSegment
// ----------------------------------------------------------------------------

func ExampleTask_RemoveSegment() {
	tsk, err := task.New("buy milk @home +shopping")
	if err != nil {
		panic(err)
	}

	err = tsk.RemoveSegment("@home")
	if err != nil {
		panic(err)
	}

	err = tsk.RemoveSegment("+shopping")
	if err != nil {
		panic(err)
	}

	err = tsk.Apply()
	if err != nil {
		panic(err)
	}

	fmt.Println(tsk)
	// Output: buy milk
}

// ----------------------------------------------------------------------------
//  Task.AddContext
// ----------------------------------------------------------------------------

func ExampleTask_AddContext() {
	tsk, err := task.New("buy milk +shopping")
	if err != nil {
		panic(err)
	}

	err = tsk.AddContext("home")
	if err != nil {
		panic(err)
	}

	err = tsk.Apply()
	if err != nil {
		panic(err)
	}

	fmt.Println(tsk)
	// Output: buy milk +shopping @home
}

// ----------------------------------------------------------------------------
//  Task.RemoveContext
// ----------------------------------------------------------------------------

func ExampleTask_RemoveContext() {
	tsk, err := task.New("buy milk @home @work +shopping")
	if err != nil {
		panic(err)
	}

	err = tsk.RemoveContext("home")
	if err != nil {
		panic(err)
	}

	err = tsk.Apply()
	if err != nil {
		panic(err)
	}

	fmt.Println(tsk)
	// Output: buy milk @work +shopping
}

// ----------------------------------------------------------------------------
//  Task.InsertAfter
// ----------------------------------------------------------------------------

// InsertAfter does not check for duplicates or validity.
func ExampleTask_InsertAfter() {
	tsk, err := task.New("x buy milk")
	if err != nil {
		panic(err)
	}

	ok := tsk.InsertAfter("x", "(A)")
	if !ok {
		panic("target segment not found")
	}

	err = tsk.Apply()
	if err != nil {
		panic(err)
	}

	fmt.Println(tsk)
	// Output: x (A) buy milk
}

// ============================================================================
//  Option Examples (Functional Options)
// ============================================================================

// ----------------------------------------------------------------------------
//  WithContext
// ----------------------------------------------------------------------------

func ExampleWithContext_while_creation() {
	tsk, err := task.New("buy milk",
		task.WithContext("home"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Contexts:", tsk.Contexts())
	// Output:
	// Current task: buy milk @home
	// Contexts: [home]
}

func ExampleWithContext_after_creation() {
	tsk, err := task.New("buy milk")
	if err != nil {
		panic(err)
	}

	// Equivalent to:
	//   err = tsk.AddContext("home")
	//   err = tsk.Apply()
	err = tsk.Apply(task.WithContext("home"))
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Contexts:", tsk.Contexts())
	// Output:
	// Current task: buy milk @home
	// Contexts: [home]
}

// ----------------------------------------------------------------------------
//  WithoutContext
// ----------------------------------------------------------------------------

func ExampleWithoutContext_while_creation() {
	tsk, err := task.New("buy milk @home @store",
		task.WithoutContext("home"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Contexts:", tsk.Contexts())
	// Output:
	// Current task: buy milk @store
	// Contexts: [store]
}

func ExampleWithoutContext_after_creation() {
	tsk, err := task.New("buy milk @home @store")
	if err != nil {
		panic(err)
	}

	// Equivalent to:
	//   err = tsk.RemoveContext("home")
	//   err = tsk.Apply()
	err = tsk.Apply(task.WithoutContext("home"))
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Contexts:", tsk.Contexts())
	// Output:
	// Current task: buy milk @store
	// Contexts: [store]
}

// ----------------------------------------------------------------------------
//  WithProject
// ----------------------------------------------------------------------------

func ExampleWithProject_while_creation() {
	tsk, err := task.New("buy milk",
		task.WithProject("shopping"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Projects:", tsk.Projects())
	// Output:
	// Current task: buy milk +shopping
	// Projects: [shopping]
}

func ExampleWithProject_after_creation() {
	tsk, err := task.New("buy milk")
	if err != nil {
		panic(err)
	}

	// Equivalent to:
	//   err = tsk.AddProject("shopping")
	//   err = tsk.Apply()
	err = tsk.Apply(task.WithProject("shopping"))
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Projects:", tsk.Projects())
	// Output:
	// Current task: buy milk +shopping
	// Projects: [shopping]
}

// ----------------------------------------------------------------------------
//  WithoutProject
// ----------------------------------------------------------------------------

func ExampleWithoutProject_while_creation() {
	tsk, err := task.New("buy milk +shopping +errands",
		task.WithoutProject("shopping"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Projects:", tsk.Projects())
	// Output:
	// Current task: buy milk +errands
	// Projects: [errands]
}

func ExampleWithoutProject_after_creation() {
	tsk, err := task.New("buy milk +shopping +errands")
	if err != nil {
		panic(err)
	}

	// Equivalent to:
	//   err = tsk.RemoveProject("shopping")
	//   err = tsk.Apply()
	err = tsk.Apply(task.WithoutProject("shopping"))
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Projects:", tsk.Projects())
	// Output:
	// Current task: buy milk +errands
	// Projects: [errands]
}

// ----------------------------------------------------------------------------
//  WithCompleted
// ----------------------------------------------------------------------------

func ExampleWithCompleted_default() {
	tsk, err := task.New("buy milk")
	if err != nil {
		panic(err)
	}

	// Default: use current date
	err = tsk.Apply(task.WithCompleted())
	if err != nil {
		panic(err)
	}

	fmt.Println("Is Done:", tsk.IsDone())
	fmt.Println("Has Completed Date:", tsk.DateCompleted() != "")
	// Output:
	// Is Done: true
	// Has Completed Date: true
}

func ExampleWithCompleted_withCustomDate() {
	tsk, err := task.New("buy milk")
	if err != nil {
		panic(err)
	}

	// Custom date
	err = tsk.Apply(task.WithCompleted("2024-01-01"))
	if err != nil {
		panic(err)
	}

	fmt.Println("Is Done:", tsk.IsDone())
	fmt.Println("Completed Date:", tsk.DateCompleted())
	// Output:
	// Is Done: true
	// Completed Date: 2024-01-01
}

func ExampleWithCompleted_withoutDate() {
	tsk, err := task.New("buy milk")
	if err != nil {
		panic(err)
	}

	// No date (non-standard)
	err = tsk.Apply(task.WithCompleted(""))
	if err != nil {
		panic(err)
	}

	fmt.Println("Is Done:", tsk.IsDone())
	fmt.Println("Completed Date:", tsk.DateCompleted())
	// Output:
	// Is Done: true
	// Completed Date:
}

// ----------------------------------------------------------------------------
//  WithKeyValue
// ----------------------------------------------------------------------------

func ExampleWithKeyValue_while_creation() {
	tsk, err := task.New("buy milk",
		task.WithKeyValue("due", "2024-12-25"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Key-Values:", tsk.KeyValues())
	// Output:
	// Current task: buy milk due:2024-12-25
	// Key-Values: [{due 2024-12-25}]
}

func ExampleWithKeyValue_after_creation() {
	tsk, err := task.New("buy milk")
	if err != nil {
		panic(err)
	}

	err = tsk.Apply(task.WithKeyValue("due", "2024-12-25"))
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Key-Values:", tsk.KeyValues())
	// Output:
	// Current task: buy milk due:2024-12-25
	// Key-Values: [{due 2024-12-25}]
}

func ExampleWithKeyValue_updateExisting() {
	tsk, err := task.New("buy milk due:2024-01-01")
	if err != nil {
		panic(err)
	}

	// Update existing key-value
	err = tsk.Apply(task.WithKeyValue("due", "2024-12-25"))
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Key-Values:", tsk.KeyValues())
	// Output:
	// Current task: buy milk due:2024-12-25
	// Key-Values: [{due 2024-12-25}]
}

// ----------------------------------------------------------------------------
//  WithoutKeyValue
// ----------------------------------------------------------------------------

func ExampleWithoutKeyValue_while_creation() {
	tsk, err := task.New("buy milk due:2024-12-25 priority:high",
		task.WithoutKeyValue("due"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Key-Values:", tsk.KeyValues())
	// Output:
	// Current task: buy milk priority:high
	// Key-Values: [{priority high}]
}

func ExampleWithoutKeyValue_after_creation() {
	tsk, err := task.New("buy milk due:2024-12-25 priority:high")
	if err != nil {
		panic(err)
	}

	err = tsk.Apply(task.WithoutKeyValue("priority"))
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Key-Values:", tsk.KeyValues())
	// Output:
	// Current task: buy milk due:2024-12-25
	// Key-Values: [{due 2024-12-25}]
}

// ----------------------------------------------------------------------------
//  WithDueDate
// ----------------------------------------------------------------------------

func ExampleWithDueDate_while_creation() {
	tsk, err := task.New("buy milk",
		task.WithDueDate("2024-12-25"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Key-Values:", tsk.KeyValues())
	// Output:
	// Current task: buy milk due:2024-12-25
	// Key-Values: [{due 2024-12-25}]
}

func ExampleWithDueDate_after_creation() {
	tsk, err := task.New("buy milk")
	if err != nil {
		panic(err)
	}

	err = tsk.Apply(task.WithDueDate("2024-12-25"))
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Key-Values:", tsk.KeyValues())
	// Output:
	// Current task: buy milk due:2024-12-25
	// Key-Values: [{due 2024-12-25}]
}

func ExampleWithDueDate_updateExisting() {
	tsk, err := task.New("buy milk due:2024-01-01")
	if err != nil {
		panic(err)
	}

	err = tsk.Apply(task.WithDueDate("2024-12-25"))
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Key-Values:", tsk.KeyValues())
	// Output:
	// Current task: buy milk due:2024-12-25
	// Key-Values: [{due 2024-12-25}]
}

// ----------------------------------------------------------------------------
//  WithoutDueDate
// ----------------------------------------------------------------------------

func ExampleWithoutDueDate_while_creation() {
	tsk, err := task.New("buy milk due:2024-12-25 priority:high",
		task.WithoutDueDate(),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Key-Values:", tsk.KeyValues())
	// Output:
	// Current task: buy milk priority:high
	// Key-Values: [{priority high}]
}

func ExampleWithoutDueDate_after_creation() {
	tsk, err := task.New("buy milk due:2024-12-25 priority:high")
	if err != nil {
		panic(err)
	}

	err = tsk.Apply(task.WithoutDueDate())
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Key-Values:", tsk.KeyValues())
	// Output:
	// Current task: buy milk priority:high
	// Key-Values: [{priority high}]
}

// ----------------------------------------------------------------------------
//  WithInlineComment
// ----------------------------------------------------------------------------

func ExampleWithInlineComment_while_creation() {
	tsk, err := task.New("buy milk",
		task.WithInlineComment("don't forget to buy organic"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Comment:", tsk.Comment())
	// Output:
	// Current task: buy milk # don't forget to buy organic
	// Comment: # don't forget to buy organic
}

func ExampleWithInlineComment_after_creation() {
	tsk, err := task.New("buy milk")
	if err != nil {
		panic(err)
	}

	err = tsk.Apply(
		task.WithInlineComment("don't forget to buy organic"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Comment:", tsk.Comment())
	// Output:
	// Current task: buy milk # don't forget to buy organic
	// Comment: # don't forget to buy organic
}

func ExampleWithInlineComment_updateExisting() {
	tsk, err := task.New("buy milk # buy low-fat")
	if err != nil {
		panic(err)
	}

	err = tsk.Apply(
		task.WithInlineComment("buy organic"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Comment:", tsk.Comment())
	// Output:
	// Current task: buy milk # buy organic
	// Comment: # buy organic
}

func ExampleWithInlineComment_allowedCtrlChars() {
	tsk, err := task.New("buy milk")
	if err != nil {
		panic(err)
	}

	// Add inline comment with tab character
	err = tsk.Apply(
		task.WithAllowedCtrlChars([]rune{'\t'}),
		task.WithInlineComment("don't forget to\tbuy organic"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Comment:", tsk.Comment())
	// Output:
	// Current task: buy milk # don't forget to	buy organic
	// Comment: # don't forget to	buy organic
}

// ----------------------------------------------------------------------------
//  WithoutInlineComment
// ----------------------------------------------------------------------------

func ExampleWithoutInlineComment_while_creation() {
	tsk, err := task.New("buy milk # don't forget to buy organic",
		task.WithoutInlineComment(),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Comment:", tsk.Comment())
	// Output:
	// Current task: buy milk
	// Comment:
}

func ExampleWithoutInlineComment_after_creation() {
	tsk, err := task.New("buy milk # don't forget to buy organic")
	if err != nil {
		panic(err)
	}

	err = tsk.Apply(
		task.WithoutInlineComment(),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Current task:", tsk)
	fmt.Println("Comment:", tsk.Comment())
	// Output:
	// Current task: buy milk
	// Comment:
}
