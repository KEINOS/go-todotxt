/*
Package todo provides a library for parsing and manipulating tasks and to-do
lists in todo.txt format.

The todo.txt format is a simple text-based format for managing tasks. It allows
users to define tasks with priorities, due dates, contexts, projects, and custom
tags. This package provides structs and functions to parse, create, and manipulate
such tasks.

Key features:
- Parse todo.txt formatted strings into Task structs
- Create and modify tasks programmatically
- Filter and sort tasks based on various criteria
- Load and save task lists from/to files
- Support for all standard todo.txt elements: priority, completion, dates, contexts, projects, tags

Example usage:

	import "github.com/KEINOS/go-todotxt/todo"

	// Parse a task
	task, err := todo.ParseTask("(A) Call Mom @Phone +Family due:2023-12-25")
	if err != nil {
		log.Fatal(err)
	}

	// Create a task list
	tasks := todo.NewTaskList()
	tasks.AddTask(task)

	// Filter completed tasks
	completed := tasks.Filter(todo.FilterCompleted)

For more examples, see the example functions in this package.
*/
package todo

// Go generate directives.
//
// These will generate the stringer implementations for TaskSortByType and
// TaskSegmentType types.
// Note that to call `go generate ./...` you need `stringer` command installed.
// You can use `docker compose run go_generate` for convenience.
//
//go:generate stringer -type TaskSortByType -trimprefix Sort -output tasksortbytype_string.go
//go:generate stringer -type TaskSegmentType -trimprefix Segment -output tasksegmenttype_string.go
