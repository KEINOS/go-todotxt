package todo_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/KEINOS/go-todotxt/todo"
)

// ============================================================================
//  Basic Usage
// ============================================================================

func Example() {
	// Load tasks from a string.
	// You can also load from a file by using LoadFromFile().
	tasks, err := todo.LoadFromString(`
        (A) Call Mom @Phone +Family
        x (A) Schedule annual checkup +Health
        (C) Add cover sheets @Office +TPSReports
        Plan backyard herb garden @Home
        Pick up milk @GroceryStore
        Research self-publishing services +Novel @Computer
        x Download Todo.txt mobile app @Phone
		(A) This is a task should be due before yesterday due:2020-11-15
    `)
	if err != nil {
		log.Fatal(err)
	}

	// Add a new task.
	// Note the order of project ("+Novel") and context ("@Computer") later
	// in the output.
	newTask, err := todo.ParseTask("(B) Outline chapter 5 +Novel @Computer")
	if err != nil {
		log.Fatal(err)
	}

	tasks.AddTask(newTask) // append to the end of the list

	// Sort tasks by priority and then by context in ascending order.
	err = tasks.Sort(todo.SortPriorityAsc, todo.SortContextAsc)
	if err != nil {
		log.Fatal(err)
	}

	// AND filter.
	// Get tasks that have any priority AND are not completed AND are not overdue.
	filteredTasks := tasks.
		Filter(todo.FilterHasPriority).
		Filter(todo.FilterNotCompleted).
		Filter(todo.FilterNot(todo.FilterOverdue)) // NOT overdue

	// OR filter.
	// Filter the above tasks with priority "A" OR has project "Novel" OR has context "Office".
	filteredTasks = filteredTasks.Filter(
		todo.FilterByPriority("A"),     // has (A)
		todo.FilterByProject("Novel"),  // has +Novel
		todo.FilterByContext("Office"), // has @Office
	)

	// Print the filtered tasks.
	for _, task := range filteredTasks {
		fmt.Println(task.String())
	}
	// Output:
	// (A) Call Mom @Phone +Family
	// (B) Outline chapter 5 @Computer +Novel
	// (C) Add cover sheets @Office +TPSReports
}

func Example_with_comments() {
	// Load tasks from a string.
	// You can also load from a file by using LoadFromFile().
	tasks, err := todo.LoadFromString(`
		# This is a comment line.
        (A) Call Mom @Phone +Family
        x (C) Schedule annual checkup +Health
        Pick up milk @GroceryStore

		# This is another comment line.
        Research self-publishing services +Novel @Computer
        x Download Todo.txt mobile app @Phone
		(A) This is a task should be due before yesterday due:2020-11-15
    `)
	if err != nil {
		log.Fatal(err)
	}

	for _, task := range tasks {
		fmt.Println(task)
	}
	// Output:
	// (A) Call Mom @Phone +Family
	// x (C) Schedule annual checkup +Health
	// Pick up milk @GroceryStore
	// Research self-publishing services @Computer +Novel
	// x Download Todo.txt mobile app @Phone
	// (A) This is a task should be due before yesterday due:2020-11-15
}

func ExampleLoadFromFile() {
	// Create a reader from a string. Usually, for strings, you can use
	// LoadFromString() instead.
	strReader := strings.NewReader(`
		(A) Call Mom @Phone +Family
		x (A) Schedule annual checkup +Health
		(B) Outline chapter 5 +Novel @Computer
		(C) Add cover sheets @Office +TPSReports
		Plan backyard herb garden @Home
		Pick up milk @GroceryStore
		Research self-publishing services +Novel @Computer
		x Download Todo.txt mobile app @Phone
		(A) This is a task should be due before yesterday due:2020-11-15
	`)

	// Load tasks from any io.Reader (e.g. os.File, os.Stdin, etc.)
	tasks, err := todo.LoadFromFile(strReader)
	if err != nil {
		log.Fatal(err)
	}

	// Get tasks that are overdue (out of due date)
	filteredTasks := tasks.Filter(todo.FilterOverdue)

	for _, task := range filteredTasks {
		fmt.Println(task)
	}
	// Output:
	// (A) This is a task should be due before yesterday due:2020-11-15
}

func ExampleLoadFromPath() {
	tasks, err := todo.LoadFromPath("testdata/filter_todo.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Get tasks that are not completed AND are overdue AND has priority.
	filteredTasks := tasks.
		Filter(todo.FilterNotCompleted).
		Filter(todo.FilterOverdue).
		Filter(todo.FilterHasPriority)

	for _, task := range filteredTasks {
		fmt.Println(task)
	}
	// Output:
	// (B) 2013-12-01 Outline chapter 5 @Computer +Novel Level:5 private:false due:2014-02-17
}

func ExampleLoadFromString() {
	tasks, err := todo.LoadFromString(`
		(A) Call Mom @Phone +Family
		x (A) Schedule annual checkup +Health
		(B) Outline chapter 5 +Novel @Computer
		(C) Add cover sheets @Office +TPSReports
		Plan backyard herb garden @Home
		Pick up milk @GroceryStore
		Research self-publishing services +Novel @Computer
		x Download Todo.txt mobile app @Phone
		(A) This is a task should be due before yesterday due:2020-11-15
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Get tasks that are overdue (out of due date)
	filteredTasks := tasks.Filter(todo.FilterOverdue)

	for _, task := range filteredTasks {
		fmt.Println(task)
	}
	// Output:
	// (A) This is a task should be due before yesterday due:2020-11-15
}

// ============================================================================
//  TaskList
// ============================================================================

// ----------------------------------------------------------------------------
//  TaskList.CustomSort
// ----------------------------------------------------------------------------

func ExampleTaskList_CustomSort() {
	tasks, err := todo.LoadFromString(`
		Task 3
		Task 1
		Task 4
		Task 2
	`)
	if err != nil {
		log.Fatal(err)
	}

	// customFunc returns true if taskA is less than taskB.
	customFunc := func(taskA, taskB todo.Task) bool {
		// Compare strings of the text part of the task.
		return taskA.Todo < taskB.Todo
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

// ----------------------------------------------------------------------------
//  TaskList.Sort
// ----------------------------------------------------------------------------

func ExampleTaskList_Sort() {
	tasks, err := todo.LoadFromPath("testdata/sort_todo.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Before #1:", tasks[0].Projects)
	fmt.Println("Before #2:", tasks[1].Projects)
	fmt.Println("Before #3:", tasks[2].Projects)

	// Sort tasks by project and then priority in ascending order
	err = tasks.Sort(todo.SortProjectAsc, todo.SortPriorityAsc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("After  #1:", tasks[0].Projects)
	fmt.Println("After  #2:", tasks[1].Projects)
	fmt.Println("After  #3:", tasks[2].Projects)
	// Output:
	// Before #1: []
	// Before #2: [go-todotxt]
	// Before #3: [go-todotxt]
	// After  #1: [Apple]
	// After  #2: [Apple]
	// After  #3: [Apple]
}
