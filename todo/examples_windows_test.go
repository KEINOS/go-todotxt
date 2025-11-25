//go:build windows

package todo_test

import (
	"fmt"
	"log"

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

	// Note that the end of line character is '\r\n'. It differs from OS to OS.
	// '\n' for Unix-like OS and '\r\n' for Windows.
	fmt.Printf("%#v", tasks.String())
	// Output:
	// "(A) Call Mom @Phone +Family\r\n(A) Schedule annual checkup +Health\r\n(B) Outline chapter 5 @Computer +Novel\r\n(C) Add cover sheets @Office +TPSReports\r\nPlan backyard herb garden @Home\r\nPick up milk @GroceryStore\r\nResearch self-publishing services @Computer +Novel\r\nx Download Todo.txt mobile app @Phone\r\n"
}
