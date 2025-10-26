//go:build !windows
// +build !windows

package todo_test

import (
	"fmt"
	"log"

	//nolint:depguard // external test package needs to import the module
	"github.com/KEINOS/go-todotxt/todo"
)

// ============================================================================
//  TaskList
// ============================================================================

// ----------------------------------------------------------------------------
//  TaskList.String
// ----------------------------------------------------------------------------

//nolint:lll // Allow long line length for output example.
func ExampleTaskList_String() {
	tasks, err := todo.LoadFromPath("testdata/todo.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Note that the end of line character is '\n'. It differs from OS to OS.
	// '\n' for Unix-like OS and '\r\n' for Windows.
	fmt.Printf("%#v", tasks.String())
	// Output:
	// "(A) Call Mom @Phone +Family\n(A) Schedule annual checkup +Health\n(B) Outline chapter 5 @Computer +Novel\n(C) Add cover sheets @Office +TPSReports\nPlan backyard herb garden @Home\nPick up milk @GroceryStore\nResearch self-publishing services @Computer +Novel\nx Download Todo.txt mobile app @Phone\n"
}
