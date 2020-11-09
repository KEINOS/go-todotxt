// +build !windows

package todotxt

import (
	"fmt"
	"log"
)

func ExampleLoadFromPath() {
	if tasklist, err := LoadFromPath("testdata/todo.txt"); err != nil {
		log.Fatal(err)
	} else {
		fmt.Print(tasklist) // String representation of TaskList works as expected.
	}

	// Output:
	// (A) Call Mom @Phone +Family
	// (A) Schedule annual checkup +Health
	// (B) Outline chapter 5 @Computer +Novel
	// (C) Add cover sheets @Office +TPSReports
	// Plan backyard herb garden @Home
	// Pick up milk @GroceryStore
	// Research self-publishing services @Computer +Novel
	// x Download Todo.txt mobile app @Phone
}

func ExampleTaskList_LoadFromPath() {
	var tasklist TaskList

	// This will overwrite whatever was in the tasklist before.
	// Irrelevant here since the list is still empty.
	if err := tasklist.LoadFromPath("testdata/todo.txt"); err != nil {
		log.Fatal(err)
	}

	fmt.Println(tasklist[0].Todo)      // Text part of first task (Call Mom)
	fmt.Println(tasklist[2].Contexts)  // Slice of contexts from third task ([Computer])
	fmt.Println(tasklist[3].Priority)  // Priority of fourth task (C)
	fmt.Println(tasklist[7].Completed) // Completed flag of eighth task (true)
	// Output:
	// Call Mom
	// [Computer]
	// C
	// true
}