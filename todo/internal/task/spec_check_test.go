package task_test

import (
	"testing"

	"github.com/KEINOS/go-todotxt/todo/internal/todoparse"
	"github.com/KEINOS/go-todotxt/todo/internal/task"
	"github.com/stretchr/testify/require"
)

// ============================================================================
//  Spec Checks
// ============================================================================
//  These tests ensuure that the package follows the rules of todo.txt spec by
//  using the official todo.txt example.
//
//  Some tests and descriptions may seem redundant or a duplicate of other package
//  tests, but they are kept here to ensure compliance with the spec.
//
//	* Reference:
//    * As of 2025-11-18 @ Commit e79d866
//    * https://github.com/todotxt/todo.txt/tree/e79d866af927013dc4fcfbe27ed51254a2cac394
//    * For the latest spec see: https://github.com/todotxt/todo.txt

// ----------------------------------------------------------------------------
//  Incomplete Tasks: 3 Format Rules
// ----------------------------------------------------------------------------

// Rule 1: If priority exists, it ALWAYS appears first.
//
// The priority is an uppercase character from A-Z enclosed in parentheses and
// followed by a space.
func Test_Incomplete_Rule1(t *testing.T) {
	t.Parallel()

	t.Run("Has priority", func(t *testing.T) {
		t.Parallel()

		tsk, err := task.New("(A) Call Mom")
		require.NoError(t, err)

		require.Equal(t, "A", tsk.Priority())
		require.Equal(t, "Call Mom", tsk.Description())
	})

	t.Run("No priority", func(t *testing.T) {
		t.Parallel()

		for index, taskTxt := range []string{
			"Really gotta call Mom (A) @phone @someday",
			"(b) Get back to the boss",
			"(B)->Submit TPS report",
		} {
			tsk, err := task.New(taskTxt)
			require.NoError(t, err,
				"Case #%d should not error: %s", index+1, taskTxt)

			require.Empty(t, tsk.Priority(),
				"Case #%d should have no priority: %s", index+1, taskTxt)
		}
	})
}

// Rule 2: A task's creation date may optionally appear directly after priority
// and a space.
//
// If there is no priority, the creation date appears first. If the creation date
// exists, it should be in the format `YYYY-MM-DD`.
func Test_Incomplete_Rule2(t *testing.T) {
	t.Parallel()

	t.Run("Have creation date", func(t *testing.T) {
		t.Parallel()

		for index, taskTxt := range []string{
			"2011-03-02 Document +TodoTxt task format",
			"(A) 2011-03-02 Call Mom",
		} {
			tsk, err := task.New(taskTxt)
			require.NoError(t, err,
				"Case #%d should not error: %s", index+1, taskTxt)

			require.Equal(t, "2011-03-02", tsk.DateCreated(),
				"Case #%d should have creation date: %s", index+1, taskTxt)
		}
	})

	t.Run("No creation date", func(t *testing.T) {
		t.Parallel()

		tsk, err := task.New("(A) Call Mom 2011-03-02")
		require.NoError(t, err)

		require.Empty(t, tsk.DateCreated())
	})
}

// Rule 3: Contexts and Projects may appear anywhere in the line after priority/
// prepended date.
//
//   - A context is preceded by a single space and an at-sign (`@`).
//   - A project is preceded by a single space and a plus-sign (`+`).
//   - A project or context contains any non-whitespace character.
//   - A task may have zero, one, or more than one projects and contexts included
//     in it.
//
// For example, this task is part of the +Family and +PeaceLoveAndHappiness
// projects as well as the `@iphone` and `@phone` contexts:
//
//	"((A) Call Mom +Family +PeaceLoveAndHappiness @iphone @phone"
func Test_Incomplete_Rule3(t *testing.T) {
	t.Parallel()

	t.Run("No context in it", func(t *testing.T) {
		t.Parallel()

		tsk, err := task.New("Email SoAndSo at soandso@example.com")
		require.NoError(t, err)

		require.Empty(t, tsk.Contexts())
	})

	t.Run("No project in it", func(t *testing.T) {
		t.Parallel()

		tsk, err := task.New("Learn how to add 2+2")
		require.NoError(t, err)

		require.Empty(t, tsk.Projects())
	})
}

// ----------------------------------------------------------------------------
//  Complete Tasks: 2 Format Rules
// ----------------------------------------------------------------------------

// Rule 1: A completed task starts with a lowercase x character (x).
//
// If a task starts with an `x“ (case-sensitive and lowercase) followed directly
// by a space, it is marked as complete.
func Test_Completed_Rule1(t *testing.T) {
	t.Parallel()

	t.Run("Is complete", func(t *testing.T) {
		t.Parallel()

		tsk, err := task.New("x 2011-03-03 Call Mom")
		require.NoError(t, err)
		require.True(t, tsk.IsCompleted())
	})

	t.Run("Not complete", func(t *testing.T) {
		t.Parallel()

		for index, taskTxt := range []string{
			"xylophone lesson",
			"X 2012-01-01 Make resolutions",
			"(A) x Find ticket prices",
		} {
			tsk, err := task.New(taskTxt)
			require.NoError(t, err,
				"Case #%d should not error: %s", index+1, taskTxt)

			require.False(t, tsk.IsCompleted(),
				"Case #%d should not be marked as complete: %s", index+1, taskTxt)
		}
	})
}

// Rule 2: The date of completion appears directly after the x, separated by a
// space.
//
// If you’ve prepended the creation date to your task, on completion it will
// appear directly after the completion date. This is so your completed tasks
// sort by date using standard sort tools. Many Todo.txt clients discard priority
// on task completion.
// To preserve it, use the `key:value` format (e.g. `pri:A`).
//
// With the completed date (required), if you've used the prepended date (optional),
// you can calculate how many days it took to complete a task.
func Test_Completed_Rule2(t *testing.T) {
	t.Parallel()

	const taskStr = "x 2011-03-02 2011-03-01 Review Tim's pull request +TodoTxtTouch @github"

	tsk, err := task.New(taskStr)
	require.NoError(t, err)

	require.Equal(t, "2011-03-02", tsk.DateCompleted())
	require.Equal(t, "2011-03-01", tsk.DateCreated())
}

// ----------------------------------------------------------------------------
//  Additional File Format Definitions
// ----------------------------------------------------------------------------

// Tool developers may define additional formatting rules for extra metadata.
//
// Developers should use the format `key:value` to define additional metadata
// (e.g. `due:2010-01-02` as a due date).
//
// Both `key` and `value` must consist of non-whitespace characters, which are
// not colons. Only one colon separates the `key` and `value`.
func Test_Additional_KeyValue(t *testing.T) {
	t.Parallel()

	const taskStr = "x (A) 2016-05-20 2016-04-30 measure space for " +
		"+chapelShelving @chapel due:2016-05-30"

	tsk, err := task.New(taskStr)
	require.NoError(t, err)

	expect := []todoparse.KeyValue{
		{Key: "due", Value: "2016-05-30"},
	}
	actual := tsk.KeyValues()
	require.Equal(t, expect, actual)
}
