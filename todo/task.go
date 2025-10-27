package todo

import (
	"strings"
	"time"
)

// Clock interface for time operations, allowing dependency injection for testing.
type Clock interface {
	Now() time.Time
}

// realClock implements Clock using the real time.
type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now()
}

// ----------------------------------------------------------------------------
//  Type: Task
// ----------------------------------------------------------------------------

// Task represents a todo.txt task entry.
//
// The 'Contexts' and 'Projects' are both used to categorize tasks.
// The difference is that 'Contexts' are used to categorize tasks by location or
// situation where you'll work on the job, while 'Projects' are used to categorize
// tasks by project.
//
// For the "todo.txt" format specification see:
// https://github.com/todotxt/todo.txt#todotxt-format-rules
//
//nolint:godox,recvcheck // False positive TODO in the comment. Stringer requires non-pointer receiver for String()
type Task struct {
	DueDate        time.Time         // DueDate is the due date calculated from the 'due:' tag.
	CompletedDate  time.Time         // CompletedDate is the date the task was completed.
	CreatedDate    time.Time         // CreatedDate is the date the task was created.
	AdditionalTags map[string]string // AdditionalTags of the task in a key:value format (e.g. "due:2012-12-12")
	Original       string            // Original raw task text.
	Priority       string            // Priority of the task in (A)-(Z) range.
	Todo           string            // Todo part of task text.
	Contexts       []string          // Contexts of the task (e.g. @MyContext).
	Projects       []string          // Projects of the task (e.g. +MyProject).
	ID             int               // ID of the task internaly.
	Completed      bool              // Completed flag. If true, the task has been completed.
	clock          Clock             // clock for time operations, defaults to realClock.
}

// ----------------------------------------------------------------------------
//  Constructors
// ----------------------------------------------------------------------------

// NewTask creates a new empty Task with default values. (CreatedDate is set to Now()).
func NewTask() Task {
	task := new(Task)
	task.clock = realClock{}
	task.CreatedDate = task.clock.Now()

	return *task
}

// NewTaskWithClock creates a new empty Task with a custom clock for testing.
func NewTaskWithClock(clock Clock) Task {
	task := new(Task)
	task.clock = clock
	task.CreatedDate = task.clock.Now()

	return *task
}

// ParseTask parses the input text string into a Task struct.
func ParseTask(text string) (*Task, error) {
	parser := newTaskParser(text)

	return parser.parse()
}

// ----------------------------------------------------------------------------
//  Methods
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
//  Status Methods
// ----------------------------------------------------------------------------

// Complete sets Task.Completed to 'true' if the task was not already completed.
// Also sets Task.CompletedDate to time.Now().
func (task *Task) Complete() {
	if !task.Completed {
		task.Completed = true
		if task.clock == nil {
			task.clock = realClock{}
		}

		task.CompletedDate = task.clock.Now()
	}
}

// IsCompleted returns true if the task has already been completed.
func (task *Task) IsCompleted() bool {
	return task.Completed
}

// Reopen sets Task.Completed to 'false' if the task was completed.
// Also resets Task.CompletedDate.
func (task *Task) Reopen() {
	if task.Completed {
		task.Completed = false
		task.CompletedDate = time.Time{} // time.IsZero() value
	}
}

// ----------------------------------------------------------------------------
//  Date Methods
// ----------------------------------------------------------------------------

// Due returns the duration left until due date from now. The duration is negative
// if the task is overdue.
//
// Just as with IsOverdue(), this function does also not take the Completed flag
// into consideration. You should check Task.Completed first if needed.
func (task *Task) Due() time.Duration {
	return time.Until(task.DueDate.AddDate(0, 0, 1))
}

// HasCompletedDate returns true if the task has a completed date.
func (task *Task) HasCompletedDate() bool {
	return !task.CompletedDate.IsZero() && task.Completed
}

// HasCreatedDate returns true if the task has a created date.
func (task *Task) HasCreatedDate() bool {
	return !task.CreatedDate.IsZero()
}

// HasDueDate returns true if the task has a due date.
func (task *Task) HasDueDate() bool {
	return !task.DueDate.IsZero()
}

// IsDueToday returns true if the task is due today.
func (task *Task) IsDueToday() bool {
	if task.HasDueDate() {
		due := task.Due()

		return 0 < due && due <= oneDay
	}

	return false
}

// IsOverdue returns true if due date is in the past.
//
// This function does not take the Completed flag into consideration.
// You should check Task.Completed first if needed.
func (task *Task) IsOverdue() bool {
	if task.HasDueDate() {
		return task.Due() < 0
	}

	return false
}

// ----------------------------------------------------------------------------
//  Attribute Methods
// ----------------------------------------------------------------------------

// HasAdditionalTags returns true if the task has any additional tags.
func (task *Task) HasAdditionalTags() bool {
	return len(task.AdditionalTags) > 0
}

// HasContexts returns true if the task has any contexts.
func (task *Task) HasContexts() bool {
	return len(task.Contexts) > 0
}

// HasPriority returns true if the task has a priority.
func (task *Task) HasPriority() bool {
	return isNotEmpty(task.Priority)
}

// HasProjects returns true if the task has any projects.
func (task *Task) HasProjects() bool {
	return len(task.Projects) > 0
}

// ----------------------------------------------------------------------------
//  String Methods
// ----------------------------------------------------------------------------

// String returns a complete task string in todo.txt format.
//
// Contexts, Projects and additional tags are alphabetically sorted,
// and appended at the end in the following order:
// Contexts, Projects, Tags
//
// For example:
//
//	"(A) 2013-07-23 Call Dad @Home @Phone +Family due:2013-07-31 customTag1:Important!"
func (task Task) String() string {
	segs := task.Segments()

	displays := make([]string, len(segs))
	for i, seg := range segs {
		displays[i] = seg.Display
	}

	return strings.Join(displays, " ")
}

// Task returns a complete task string in todo.txt format.
//
// It is an alias of String(). See *Task.String() for further information.
func (task *Task) Task() string {
	return task.String()
}
