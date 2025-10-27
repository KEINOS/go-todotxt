package todo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// fakeClock is a test helper that implements the Clock interface.
type fakeClock struct {
	now time.Time
}

func (f *fakeClock) Now() time.Time {
	return f.now
}

// ----------------------------------------------------------------------------
//  Tests: Constructors
// ----------------------------------------------------------------------------

//nolint:paralleltest // do not parallel to avoid race conditions
func TestNewTask_default_state(t *testing.T) {
	task := NewTask()

	t.Run("ID", func(t *testing.T) {
		expectID := 0
		actualID := task.ID
		require.Equal(t, expectID, actualID, "field ID failed to return expected ID")
	})

	t.Run("Original", func(t *testing.T) {
		expectStrRaw := ""
		actualStrRaw := task.Original
		require.Equal(t, expectStrRaw, actualStrRaw, "field Original failed to return expected string")
	})

	t.Run("Todo", func(t *testing.T) {
		expectTodo := ""
		actualTodo := task.Todo
		require.Equal(t, expectTodo, actualTodo, "field Todo failed to return expected string")
	})

	t.Run("HasPriority", func(t *testing.T) {
		require.False(t, task.HasPriority(),
			"method HasPriority failed to return false. the given task does not have a priority")
	})

	t.Run("HasProjects", func(t *testing.T) {
		require.False(t, task.HasProjects(),
			"method HasProjects failed to return false. the given task does not have projects")
	})

	t.Run("Projects", func(t *testing.T) {
		expectLenProj := 0
		actualLenProj := len(task.Projects)
		require.Equal(t, expectLenProj, actualLenProj, "field Projects failed to return expected length")
	})

	t.Run("HasContexts", func(t *testing.T) {
		require.False(t, task.HasContexts(),
			"method HasContexts failed to return false. the given task does not have contexts")
	})

	t.Run("Contexts", func(t *testing.T) {
		expectLenContexts := 0
		actualLenContexts := len(task.Contexts)
		require.Equal(t, expectLenContexts, actualLenContexts, "field Contexts failed to return expected length")
	})

	t.Run("HasAdditionalTags", func(t *testing.T) {
		require.False(t, task.HasAdditionalTags(),
			"method HasAdditionalTags failed to return false. the given task does not have additional tags")
	})

	t.Run("AdditionalTags", func(t *testing.T) {
		expectLenTags := 0
		actualLenTags := len(task.AdditionalTags)
		require.Equal(t, expectLenTags, actualLenTags, "field AdditionalTags failed to return expected length")
	})

	t.Run("HasCreatedDate", func(t *testing.T) {
		require.True(t, task.HasCreatedDate(),
			"method HasCreatedDate failed to return true. newly created tasks are automatically assigned a creation date")
	})

	t.Run("HasCompletedDate", func(t *testing.T) {
		require.False(t, task.HasCompletedDate(),
			"method HasCompletedDate failed to return false. the given task does not have a completed date")
	})

	t.Run("HasDueDate", func(t *testing.T) {
		require.False(t, task.HasDueDate(),
			"method HasDueDate failed to return false. the given task does not have a due date")
	})

	t.Run("Completed", func(t *testing.T) {
		require.False(t, task.Completed,
			"field Completed failed to return false. the given task is not completed")
	})
}

func TestNewTaskWithClock(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(2023, 10, 27, 12, 0, 0, 0, time.UTC)
	clock := &fakeClock{now: fixedTime}

	task := NewTaskWithClock(clock)

	require.Equal(t, fixedTime, task.CreatedDate, "CreatedDate should be set to the clock's now time")
	require.True(t, task.HasCreatedDate(), "HasCreatedDate should return true")
}

func TestTask_Complete_with_nil_clock(t *testing.T) {
	t.Parallel()

	// Create a task without setting clock (should be nil)
	task := new(Task)

	// Complete should set clock to realClock if nil
	task.Complete()

	require.True(t, task.Completed, "task should be completed")
	require.True(t, task.HasCompletedDate(), "task should have completed date")
}

func Test_ParseTask(t *testing.T) {
	t.Parallel()

	// ParseTask()
	task, err := ParseTask(
		"x (C) 2014-01-01 @Go due:2014-01-12 Create golang library documentation +go-todotxt  hello:world not::tag  ",
	)
	require.NoError(t, err, "method ParseTask failed to parse task")

	t.Run("Task", func(t *testing.T) {
		t.Parallel()

		expectStr := "x (C) 2014-01-01 Create golang library documentation " +
			"not::tag @Go +go-todotxt hello:world due:2014-01-12"
		actualStr := task.Task()
		require.Equal(t, expectStr, actualStr, "method Task failed to return expected string")
	})

	t.Run("ID", func(t *testing.T) {
		t.Parallel()

		expectID := 0
		actualID := task.ID
		require.Equal(t, expectID, actualID, "field ID failed to return expected ID")
	})

	t.Run("Original", func(t *testing.T) {
		t.Parallel()

		expectRaw := "x (C) 2014-01-01 @Go due:2014-01-12 Create " +
			"golang library documentation +go-todotxt  hello:world not::tag"
		actualRaw := task.Original
		require.Equal(t, expectRaw, actualRaw, "field Original failed to return expected string")
	})

	t.Run("Todo", func(t *testing.T) {
		t.Parallel()

		expectTask := "Create golang library documentation not::tag"
		actualTask := task.Todo
		require.Equal(t, expectTask, actualTask, "field Todo failed to return expected string")
	})

	t.Run("HasPriority", func(t *testing.T) {
		t.Parallel()
		require.True(t, task.HasPriority(),
			"method HasPriority failed to return true. the given task has a priority")
	})

	t.Run("Priority", func(t *testing.T) {
		t.Parallel()

		expectPriority := "C"
		actualPriority := task.Priority
		require.Equal(t, expectPriority, actualPriority, "field Priority failed to return expected string")
	})

	t.Run("HasProjects", func(t *testing.T) {
		t.Parallel()
		require.True(t, task.HasProjects(), "method HasProjects failed to return true. the given task has projects")
	})

	t.Run("Projects", func(t *testing.T) {
		t.Parallel()

		expectLenProj := 1
		actualLenProj := len(task.Projects)
		require.Equal(t, expectLenProj, actualLenProj, "field Projects failed to return expected length")
	})

	t.Run("HasContexts", func(t *testing.T) {
		t.Parallel()
		require.True(t, task.HasContexts(), "method HasContexts failed to return true. the given task has contexts")
	})

	t.Run("Contexts", func(t *testing.T) {
		t.Parallel()

		expectLenContexts := 1
		actualLenContexts := len(task.Contexts)
		require.Equal(t, expectLenContexts, actualLenContexts, "field Contexts failed to return expected length")
	})

	t.Run("HasAdditionalTags", func(t *testing.T) {
		t.Parallel()
		require.True(t, task.HasAdditionalTags(),
			"method HasAdditionalTags failed to return true. the given task has additional tags")
	})

	t.Run("AdditionalTags", func(t *testing.T) {
		t.Parallel()

		expectLenTag := 1
		actualLenTag := len(task.AdditionalTags)
		require.Equal(t, expectLenTag, actualLenTag, "field AdditionalTags failed to return expected length")
	})

	t.Run("Completed", func(t *testing.T) {
		t.Parallel()
		require.True(t, task.Completed, "field Completed failed to return true. the given task is completed")
	})

	t.Run("HasCreatedDate", func(t *testing.T) {
		t.Parallel()
		require.True(t, task.HasCreatedDate(),
			"method HasCreatedDate failed to return true. the given task has a created date")
	})

	t.Run("HasDueDate", func(t *testing.T) {
		t.Parallel()
		require.True(t, task.HasDueDate(), "method HasDueDate failed to return true. the given task has a due date")
	})

	t.Run("HasCompletedDate", func(t *testing.T) {
		t.Parallel()
		require.False(t, task.HasCompletedDate(),
			"method HasCompletedDate failed to return false. the given task does not have a completed date")
	})
}

// ----------------------------------------------------------------------------
//  Tests for fields of Task type
// ----------------------------------------------------------------------------

func TestTask_ID(t *testing.T) {
	t.Parallel()

	// Load the test tasklist
	testTasklist := testLoadFromPath(t, testInputTask)

	for _, test := range []struct {
		taskID   int
		expectID int
	}{
		{taskID: 1, expectID: 1},
		{taskID: 5, expectID: 5},
		{taskID: 27, expectID: 27},
	} {
		taskID := test.taskID
		task := testTasklist[taskID-1]

		expect := test.expectID
		actual := task.ID

		require.Equal(t, expect, actual, "field ID of task[%d] failed to return expected ID", taskID)
	}
}

func TestTask_Priority(t *testing.T) {
	t.Parallel()

	// Load the test tasklist
	testTasklist := testLoadFromPath(t, testInputTask)

	// Golden cases
	for _, test := range []struct {
		expect string
		taskID int
	}{
		{taskID: 6, expect: "B"},
		{taskID: 7, expect: "C"},
		{taskID: 8, expect: "B"},
	} {
		taskID := test.taskID
		task := testTasklist[taskID-1]

		expectPriority := test.expect
		actualPriority := task.Priority
		require.Equal(t, expectPriority, actualPriority,
			"field Priority of task[%d] did not return the expected priority: %s", taskID, task.String())
	}

	// Test cases with no priority
	{
		taskID := 9
		task := testTasklist[taskID-1]

		require.Empty(t, task.Priority, "field Priority of task[%d] should be empty: %s", taskID, task.String())
		require.False(t, task.HasPriority(), "method HasPriority of task[%d] should return false: %s", taskID, task.String())
	}
}

func TestTask_CreatedDate(t *testing.T) {
	t.Parallel()

	// Load the test tasklist
	testTasklist := testLoadFromPath(t, testInputTask)

	// Golden cases
	for _, test := range []struct {
		expect string
		taskID int
	}{
		{taskID: 10, expect: "2012-01-30"},
		{taskID: 11, expect: "2013-02-22"},
		{taskID: 12, expect: "2014-01-01"},
		{taskID: 13, expect: "2013-12-30"},
		{taskID: 14, expect: "2014-01-01"},
	} {
		taskID := test.taskID
		task := testTasklist[taskID-1]

		expectTime, err := parseTime(test.expect)
		require.NoError(t, err, "failed to parse time for testing")

		actualTime := task.CreatedDate
		require.Equal(t, expectTime, actualTime,
			"field CreatedDate of task[%d] did not return as expected: %s", taskID, task.String())
	}

	// Missing created date
	{
		taskID := 15
		task := testTasklist[taskID-1]

		require.Empty(t, task.CreatedDate,
			"field CreatedDate of task[%d] should be empty: %s", taskID, task.String())
		require.False(t, task.HasCreatedDate(),
			"method HasCreatedDate of task[%d] should return false: %s", taskID, task.String())
	}
}

func TestTask_Contexts(t *testing.T) {
	t.Parallel()

	// Load the test tasklist
	testTasklist := testLoadFromPath(t, testInputTask)

	for _, test := range []struct {
		taskID         int
		expectContexts []string
	}{
		{taskID: 16, expectContexts: []string{"Call", "Phone"}},
		{taskID: 17, expectContexts: []string{"Office"}},
		{taskID: 18, expectContexts: []string{"Electricity", "Home", "Of_Super-Importance", "Television"}},
		{taskID: 19, expectContexts: nil}, // No contexts
	} {
		task := testTasklist[test.taskID-1]

		require.Equal(t, test.expectContexts, task.Contexts,
			"task[%d] did not have the expected contexts: %s", test.taskID, task.String())
	}
}

func TestTask_Projects(t *testing.T) {
	t.Parallel()

	// Load the test tasklist
	testTasklist := testLoadFromPath(t, testInputTask)

	for _, test := range []struct {
		taskID         int
		expectProjects []string
	}{
		{taskID: 20, expectProjects: []string{"Gardening", "Improving", "Planning", "Relaxing-Work"}},
		{taskID: 21, expectProjects: []string{"Novel"}},
		{taskID: 22, expectProjects: nil}, // No projects
	} {
		task := testTasklist[test.taskID-1]

		require.Equal(t, test.expectProjects, task.Projects,
			"Task[%d] did not have the expected projects: %s", test.taskID, task.String())
	}
}

func TestTask_DueDate(t *testing.T) {
	t.Parallel()

	// Load the test tasklist
	testTasklist := testLoadFromPath(t, testInputTask)

	t.Run("HasDueDate", func(t *testing.T) {
		t.Parallel()

		taskID := 23
		task := testTasklist[taskID-1]

		expectTime, err := parseTime("2014-02-17")
		require.NoError(t, err, "failed to parse expected time for testing")

		actualTime := task.DueDate

		require.Equal(t, expectTime, actualTime, "task[%d] should have a due date: %s", taskID, task.String())
	})

	t.Run("NoDueDate", func(t *testing.T) {
		t.Parallel()

		taskID := 24
		task := testTasklist[taskID-1]

		// HasDueDate()
		require.False(t, task.HasDueDate(),
			"task[%d] does not have a due date but HasDueDate() returned true: %s", taskID, task.String())
		// IsDueToday()
		require.False(t, task.IsDueToday(),
			"task[%d] should not havebe due today but IsDueToday() returned true: %s", taskID, task.String())
	})

	t.Run("Overdue", func(t *testing.T) {
		t.Parallel()

		task, err := ParseTask("Hello Yesterday Task due:" + time.Now().AddDate(0, 0, -1).Format(DateLayout))

		require.NoError(t, err,
			"failed to parse task during testing")
		require.Less(t, task.Due(), time.Duration(0),
			"on overdue the duration time should be negative: %s", task.String())
		require.True(t, task.IsOverdue(),
			"on overdue tasks IsOverdue should return true: %s", task.String())
		require.False(t, task.IsDueToday(),
			"on overdue tasks IsDueToday should return false: %s", task.String())
	})

	t.Run("DueToday", func(t *testing.T) {
		t.Parallel()

		task, err := ParseTask("Hello Today Task due:" + time.Now().Format(DateLayout))

		require.NoError(t, err,
			"failed to parse task during testing")
		require.Less(t, task.Due(), 24*time.Hour,
			"on due today tasks duration time should not be greater than one day: %s", task.String())
		require.False(t, task.IsOverdue(),
			"on due today tasks IsOverdue should return false: %s", task.String())
		require.True(t, task.IsDueToday(),
			"on due today tasks IsDueToday should return true: %s", task.String())
	})

	t.Run("DueTomorrow", func(t *testing.T) {
		t.Parallel()

		task, err := ParseTask("Hello Tomorrow Task due:" + time.Now().AddDate(0, 0, 1).Format(DateLayout))

		require.NoError(t, err,
			"failed to parse task during testing")
		require.Greater(t, task.Due(), 24*time.Hour,
			"on due tomorrow tasks duration time should be greater than one day: %s", task.String())
		require.LessOrEqual(t, task.Due(), 48*time.Hour,
			"on due tomorrow tasks duration time should not be greater than two days: %s", task.String())
		require.False(t, task.IsOverdue(),
			"on due tomorrow tasks IsOverdue should return false: %s", task.String())
		require.False(t, task.IsDueToday(),
			"on due tomorrow tasks IsDueToday should return false: %s", task.String())
	})
}

func TestTask_AddonTags(t *testing.T) {
	t.Parallel()

	// Load the test tasklist
	testTasklist := testLoadFromPath(t, testInputTask)

	for _, test := range []struct {
		taskID               int
		expectAdditionalTags map[string]string
	}{
		{taskID: 25, expectAdditionalTags: map[string]string{"Level": "5", "private": "false"}},
		{taskID: 26, expectAdditionalTags: map[string]string{"Importance": "Very!"}},
		{taskID: 27, expectAdditionalTags: nil}, // No additional tags
		{taskID: 28, expectAdditionalTags: nil}, // No additional tags
	} {
		task := testTasklist[test.taskID-1]

		require.Equal(t, test.expectAdditionalTags, task.AdditionalTags,
			"task[%d] did not contain the expected additional tags: %s", test.taskID, task.String())
	}
}

// ----------------------------------------------------------------------------
//  Tests for Methods of Task type
// ----------------------------------------------------------------------------

func TestTask_IsCompleted(t *testing.T) {
	t.Parallel()

	// Load the test tasklist
	testTasklist := testLoadFromPath(t, testInputTask)

	for _, test := range []struct {
		taskID int
		expect bool
	}{
		{taskID: 31, expect: true},
		{taskID: 32, expect: false},
	} {
		task := testTasklist[test.taskID-1]

		testGot1 := task.Completed
		testGot2 := task.IsCompleted()

		require.Equal(t, test.expect, testGot2, "task[%d] should be as completed: %s", test.taskID, task.String())
		require.Equal(t, testGot1, testGot2,
			"task[%d] should be completed and IsCompleted() should return true: %s", test.taskID, task.String())
	}
}

func TestTask_Completed(t *testing.T) {
	t.Parallel()

	// Load the test tasklist
	testTasklist := testLoadFromPath(t, testInputTask)

	for _, test := range []struct {
		taskID int
		expect bool
	}{
		{29, true},
		{30, true},
		{31, true},
		{32, false},
		{33, false},
	} {
		task := testTasklist[test.taskID-1]

		if test.expect {
			require.True(t, task.Completed, "task[%d] should be completed: %s", test.taskID, task.String())
		} else {
			require.False(t, task.Completed, "task[%d] should be not completed: %s", test.taskID, task.String())
		}
	}
}

func TestTask_CompletedDate(t *testing.T) {
	t.Parallel()

	// Load the test tasklist
	testTasklist := testLoadFromPath(t, testInputTask)

	for _, test := range []struct {
		taskID           int
		hasCompletedDate bool
		expectedDate     string
	}{
		{taskID: 34, hasCompletedDate: false, expectedDate: ""},
		{taskID: 35, hasCompletedDate: true, expectedDate: "2014-01-03"},
		{taskID: 36, hasCompletedDate: false, expectedDate: ""},
		{taskID: 37, hasCompletedDate: true, expectedDate: "2014-01-02"},
		{taskID: 38, hasCompletedDate: true, expectedDate: "2014-01-03"},
		{taskID: 39, hasCompletedDate: false, expectedDate: ""},
	} {
		task := testTasklist[test.taskID-1]

		if test.hasCompletedDate {
			expectCompletedDate, err := parseTime(test.expectedDate)
			require.NoError(t, err, "failed to parse time for task[%d]", test.taskID)

			actualCompletedDate := task.CompletedDate
			require.Equal(t, expectCompletedDate, actualCompletedDate,
				"task[%d] should have a completed date of %s, but got %s", test.taskID, expectCompletedDate, actualCompletedDate)
		} else {
			require.False(t, task.HasCompletedDate(),
				"task[%d] should not have a completed date: %s", test.taskID, task.String())
		}
	}
}

//nolint:paralleltest // do not parallel to avoid race conditions
func TestTask_IsOverdue(t *testing.T) {
	// Load the test tasklist
	testTasklist := testLoadFromPath(t, testInputTask)

	for _, test := range []struct {
		taskID        int
		expectOverdue bool
		modifyDueDate func(*Task) // optional function to modify due date
		checkDueHours func(*Task) // optional function to check due hours
	}{
		{
			taskID:        40,
			expectOverdue: true,
			modifyDueDate: nil,
			checkDueHours: nil,
		},
		{
			taskID:        41,
			expectOverdue: false,
			modifyDueDate: func(task *Task) {
				task.DueDate = time.Now()
			},
			checkDueHours: func(task *Task) {
				dueHours := task.Due().Hours()
				require.True(t, dueHours > 23.0 && dueHours < 25.0,
					"task[%d] should be due in 24 hours: %s", 41, task.String())
			},
		},
		{
			taskID:        42,
			expectOverdue: true,
			modifyDueDate: func(task *Task) {
				task.DueDate = time.Now().AddDate(0, 0, -4)
			},
			checkDueHours: func(task *Task) {
				dueHours := task.Due().Hours()
				require.True(t, dueHours < 71 || dueHours > 73,
					"task[%d] should be due since 72 hours: due hours: %v, task: %s",
					42, dueHours, task.String())
			},
		},
		{
			taskID:        42,
			expectOverdue: false,
			modifyDueDate: func(task *Task) {
				task.DueDate = time.Now().AddDate(0, 0, 2)
			},
			checkDueHours: func(task *Task) {
				dueHours := task.Due().Hours()
				require.True(t, dueHours > 71 || dueHours < 73,
					"task[%d] should be due in 72 hours: due hours: %v, task: %s",
					42, dueHours, task.String())
			},
		},
		{
			taskID:        43,
			expectOverdue: false,
			modifyDueDate: nil,
			checkDueHours: nil,
		},
	} {
		task := testTasklist[test.taskID-1]

		if test.modifyDueDate != nil {
			test.modifyDueDate(&task)
		}

		require.Equal(t, test.expectOverdue, task.IsOverdue(),
			"task[%d] overdue status mismatch: %s", test.taskID, task.String())

		if test.checkDueHours != nil {
			test.checkDueHours(&task)
		}
	}
}

func TestTask_Complete(t *testing.T) {
	t.Parallel()

	// Load the test tasklist
	testTasklist := testLoadFromPath(t, testInputTask)

	for _, test := range []struct {
		taskID           int
		initialCompleted bool
		initialHasDate   bool
		expectCompleted  bool
		expectHasDate    bool
		expectDateFormat string // optional, for checking completed date
	}{
		{
			taskID:           44,
			initialCompleted: false,
			initialHasDate:   false,
			expectCompleted:  true,
			expectHasDate:    true,
			expectDateFormat: time.Now().Format(DateLayout),
		},
		{
			taskID:           45,
			initialCompleted: false,
			initialHasDate:   false,
			expectCompleted:  true,
			expectHasDate:    true,
			expectDateFormat: time.Now().Format(DateLayout),
		},
		{
			taskID:           46,
			initialCompleted: false,
			initialHasDate:   false,
			expectCompleted:  true,
			expectHasDate:    true,
			expectDateFormat: time.Now().Format(DateLayout),
		},
		{
			taskID:           47,
			initialCompleted: false,
			initialHasDate:   false,
			expectCompleted:  true,
			expectHasDate:    true,
			expectDateFormat: time.Now().Format(DateLayout),
		},
		{
			taskID:           48,
			initialCompleted: true,
			initialHasDate:   true,
			expectCompleted:  true,
			expectHasDate:    true,
			expectDateFormat: "2012-01-01", // already completed, date should not change
		},
	} {
		tmpTask := testTasklist[test.taskID-1]

		// Check initial state
		require.Equal(t, test.initialCompleted, tmpTask.Completed,
			"task[%d] initial completed state mismatch: %s", test.taskID, tmpTask.String())
		require.Equal(t, test.initialHasDate, tmpTask.HasCompletedDate(),
			"task[%d] initial has date state mismatch: %s", test.taskID, tmpTask.String())

		tmpTask.Complete() // close the task right now

		// Check after Complete()
		require.Equal(t, test.expectCompleted, tmpTask.Completed,
			"task[%d] after Complete() completed state mismatch: %s", test.taskID, tmpTask.String())
		require.Equal(t, test.expectHasDate, tmpTask.HasCompletedDate(),
			"task[%d] after Complete() has date state mismatch: %s", test.taskID, tmpTask.String())

		if test.expectDateFormat != "" {
			actualTime := tmpTask.CompletedDate.Format(DateLayout)
			require.Equal(t, test.expectDateFormat, actualTime,
				"task[%d] completed date mismatch: %s", test.taskID, tmpTask.String())
		}
	}
}

func TestTask_Reopen(t *testing.T) {
	t.Parallel()

	// Load the test tasklist
	testTasklist := testLoadFromPath(t, testInputTask)

	for _, test := range []struct {
		taskID               int
		initialCompleted     bool
		initialHasDate       bool
		expectCompletedAfter bool
		expectHasDateAfter   bool
	}{
		{
			taskID:               49,
			initialCompleted:     true,
			initialHasDate:       false,
			expectCompletedAfter: false,
			expectHasDateAfter:   false,
		},
		{
			taskID:               50,
			initialCompleted:     true,
			initialHasDate:       false,
			expectCompletedAfter: false,
			expectHasDateAfter:   false,
		},
		{
			taskID:               51,
			initialCompleted:     true,
			initialHasDate:       true,
			expectCompletedAfter: false,
			expectHasDateAfter:   false,
		},
		{
			taskID:               52,
			initialCompleted:     true,
			initialHasDate:       true,
			expectCompletedAfter: false,
			expectHasDateAfter:   false,
		},
		{
			taskID:               53,
			initialCompleted:     true,
			initialHasDate:       true,
			expectCompletedAfter: false,
			expectHasDateAfter:   false,
		},
		{
			taskID:               54,
			initialCompleted:     false,
			initialHasDate:       false,
			expectCompletedAfter: false,
			expectHasDateAfter:   false, // reopening uncompleted task
		},
	} {
		tmpTask := testTasklist[test.taskID-1]

		// Check initial state
		require.Equal(t, test.initialCompleted, tmpTask.Completed,
			"task[%d] initial completed state mismatch: %s", test.taskID, tmpTask.String())
		require.Equal(t, test.initialHasDate, tmpTask.HasCompletedDate(),
			"task[%d] initial has date state mismatch: %s", test.taskID, tmpTask.String())

		tmpTask.Reopen() // reopen the task

		// Check after Reopen()
		require.Equal(t, test.expectCompletedAfter, tmpTask.Completed,
			"task[%d] after Reopen() completed state mismatch: %s", test.taskID, tmpTask.String())
		require.Equal(t, test.expectHasDateAfter, tmpTask.HasCompletedDate(),
			"task[%d] after Reopen() has date state mismatch: %s", test.taskID, tmpTask.String())
	}
}

func TestTask_String(t *testing.T) {
	t.Parallel()

	// Load the test tasklist
	testTasklist := testLoadFromPath(t, testInputTask)

	for _, test := range []struct {
		expectStr string
		taskID    int
	}{
		{taskID: 1, expectStr: "2013-02-22 Pick up milk @GroceryStore"},
		{taskID: 2, expectStr: "x Download Todo.txt mobile app @Phone"},
		{taskID: 3, expectStr: "(B) 2013-12-01 Outline chapter 5 @Computer +Novel Level:5 private:false due:2014-02-17"},
		{taskID: 4, expectStr: "x 2014-01-02 (B) 2013-12-30 Create golang library test cases @Go +go-todotxt"},
		{taskID: 5, expectStr: "x 2014-01-03 2014-01-01 Create some more golang library test cases @Go +go-todotxt"},
	} {
		taskID := test.taskID
		task := testTasklist[taskID-1]

		expect := test.expectStr
		actual := task.String()
		require.Equal(t, expect, actual, "method String of task[%d] failed to return expected string", taskID)

		// Method Task() should return the same string
		require.Equal(t, expect, task.Task(), "method Task of task[%d] failed to return expected string", taskID)
	}
}
