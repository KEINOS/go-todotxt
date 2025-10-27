/*
Tests for todo.txt format rules based on: https://github.com/todotxt/todo.txt

Some of these tests may overlap with other unit tests, but it's intentional to
clarify the rules.
*/
package todo_test

import (
	"testing"

	"github.com/KEINOS/go-todotxt/todo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Rule 1: If priority exists, it ALWAYS appears first.
func Test_rule1(t *testing.T) {
	t.Parallel()

	t.Log("The priority is an uppercase character from A-Z " +
		"enclosed in parentheses and followed by a space.")

	t.Run(
		"Golden Case",
		func(t *testing.T) {
			t.Parallel()

			tasks, err := todo.LoadFromString(`
				(A) Call Mom
			`)
			require.NoError(t, err, "failed to load tasks")

			filteredTasks := tasks.Filter(todo.FilterHasPriority)

			require.Len(t, filteredTasks, 1,
				"this task should be treated as having priority")
			assert.Equal(t, "(A) Call Mom", filteredTasks[0].String(),
				"the task shoulbe be the same as the original")
		},
	)

	t.Run(
		"Dark Case",
		func(t *testing.T) {
			t.Parallel()

			tasks, err := todo.LoadFromString(`
				Really gotta call Mom (A) @phone @someday
				(b) Get back to the boss
				(B)->Submit TPS report
			`)
			require.NoError(t, err, "failed to load tasks")

			filteredTasks := tasks.Filter(todo.FilterHasPriority)

			require.Empty(t, filteredTasks,
				"non of the tasks should be treated as having priority")
		},
	)
}

// Rule 2: A task's creation date may optionally appear directly after priority
// and a space.
func Test_rule2(t *testing.T) {
	t.Parallel()

	t.Log("If there is no priority, the creation date appears first. " +
		"If the creation date exists, it should be in the format YYYY-MM-DD.")

	t.Run(
		"Tasks that have creation dates",
		func(t *testing.T) {
			t.Parallel()

			// These tasks have creation dates.
			tasks, err := todo.LoadFromString(`
				2011-03-02 Document +TodoTxt task format
				(A) 2011-03-02 Call Mom
			`)
			require.NoError(t, err, "failed to load tasks")

			for index, task := range tasks {
				require.NotEmptyf(t, task.CreatedDate,
					"#%v: this task should have a creation date", index)
			}
		},
	)

	t.Run(
		"Tasks that doesn't have a creation date",
		func(t *testing.T) {
			t.Parallel()

			// This task doesn't have a creation date.
			tasks, err := todo.LoadFromString(`
				(A) Call Mom 2011-03-02
				Review article of 2011-03-02
			`)
			require.NoError(t, err, "failed to load tasks")

			for index, task := range tasks {
				require.Emptyf(t, task.CreatedDate,
					"#%v: the task should not have a creation date", index+1)
			}
		},
	)
}

// Rule 3: Contexts and Projects may appear anywhere in the line after priority
// or prepended creation date.
func Test_rule3(t *testing.T) {
	t.Parallel()

	t.Log("Contexts and Projects are:\n" +
		"1. A context is preceded by a single space and an at-sign (@).\n" +
		"2. A project is preceded by a single space and a plus-sign (+).\n" +
		"3. A project or context contains any non-whitespace character.\n" +
		"4. A task may have zero, one, or more than one projects and contexts included in it.")

	t.Run(
		"Tasks that have contexts and projects",
		func(t *testing.T) {
			t.Parallel()

			// this task is part of the +Family and +PeaceLoveAndHappiness projects
			// as well as the @iphone and @phone contexts.
			task, err := todo.LoadFromString(`
				(A) Call Mom +Family +PeaceLoveAndHappiness @iphone @phone
			`)
			require.NoError(t, err, "failed to load tasks")

			// Test projects
			require.Len(t, task[0].Projects, 2,
				"the task should have 2 projects")
			require.Contains(t, task[0].Projects, "Family",
				"the task should have the project 'Family'")
			require.Contains(t, task[0].Projects, "PeaceLoveAndHappiness",
				"the task should have the project 'PeaceLoveAndHappiness'")

			// Test contexts
			require.Len(t, task[0].Contexts, 2,
				"the task should have 2 contexts")
			require.Contains(t, task[0].Contexts, "iphone",
				"the task should have the context 'iphone'")
			require.Contains(t, task[0].Contexts, "phone",
				"the task should have the context 'phone'")
		},
	)

	t.Run(
		"Tasks that has no contexts nor projects in it",
		func(t *testing.T) {
			t.Parallel()

			task, err := todo.LoadFromString(`
				Email SoAndSo at soandso@example.com
				Learn how to add 2+2
				Email SoAndSo at soandso@example.com to ask what 2+2 is
			`)
			require.NoError(t, err, "failed to load tasks")

			for index, task := range task {
				require.Emptyf(t, task.Contexts,
					"#%v: the task should not have any contexts", index+1)
				require.Emptyf(t, task.Projects,
					"#%v: the task should not have any projects", index+1)
			}
		},
	)
}

// Rule 4: Contexts and Projects may appear anywhere in the line after priority
// or prepended creation date.
