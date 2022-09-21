//go:build !windows
// +build !windows

package todo

import (
	"fmt"
	"log"
	"strings"
)

func ExampleLoadFromPath() {
	tasklist, err := LoadFromPath("testdata/todo.txt")
	if err != nil {
		log.Fatal(err)
	}

	// TaskList object implements Stringer interface
	fmt.Println(tasklist)
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

func ExampleTaskList_CustomSort() {
	//nolint:exhaustruct // fields of Task are missing but they are not used in this example
	tasks := TaskList{
		Task{Todo: "Task 3"},
		Task{Todo: "Task 1"},
		Task{Todo: "Task 4"},
		Task{Todo: "Task 2"},
	}

	customFunc := func(a, b Task) bool {
		return strings.Compare(a.Todo, b.Todo) < 0
	}

	tasks.CustomSort(customFunc)

	fmt.Println(tasks[0].Todo)
	fmt.Println(tasks[1].Todo)
	fmt.Println(tasks[2].Todo)
	fmt.Println(tasks[3].Todo)
	// Output:
	// Task 1
	// Task 2
	// Task 3
	// Task 4
}

func ExampleTaskList_Sort() {
	var tasklist TaskList
	if err := tasklist.LoadFromPath("testdata/sort_todo.txt"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Before #1:", tasklist[0].Projects)
	fmt.Println("Before #2:", tasklist[1].Projects)
	fmt.Println("Before #3:", tasklist[2].Projects)

	// sort tasks by project and then priority in ascending order
	if err := tasklist.Sort(SortProjectAsc, SortPriorityAsc); err != nil {
		log.Fatal(err)
	}

	fmt.Println("After  #1:", tasklist[0].Projects)
	fmt.Println("After  #2:", tasklist[1].Projects)
	fmt.Println("After  #3:", tasklist[2].Projects)
	// Output:
	// Before #1: []
	// Before #2: [go-todotxt]
	// Before #3: [go-todotxt]
	// After  #1: [Apple]
	// After  #2: [Apple]
	// After  #3: [Apple]
}

func ExampleTaskList_Filter() {
	var tasklist TaskList
	if err := tasklist.LoadFromPath("testdata/filter_todo.txt"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Before:", tasklist[0].String())

	// Get tasks that are not overdue and are priority A or B.
	tasklist = tasklist.Filter(FilterNot(FilterOverdue)).Filter(
		FilterByPriority("A"), FilterByPriority("B"),
	)

	fmt.Println("After :", tasklist[0].String())

	// Output:
	// Before: This is a task should be due yesterday due:2020-11-15
	// After : (A) Call Mom @Call @Phone +Family
}