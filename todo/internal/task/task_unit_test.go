package task

import (
	"fmt"
	"strings"
	"testing"

	"github.com/KEINOS/go-todotxt/todo/internal/spec"
	"github.com/stretchr/testify/require"
)

// ============================================================================
//  Test Structure Overview
//    1. Tests for public functions
//    2. Tests for private functions
//    3. Tests for public methods
//    4. Tests for private methods
//    5. Tests for issue fixes
//    6. Helper functions for tests
// ============================================================================

// ============================================================================
//  Tests for public functions
// ============================================================================

//  Constructor
// ============================================================================

// ----------------------------------------------------------------------------
//  New()
// ----------------------------------------------------------------------------

// Test to fail in each internal step of New().
func TestNew_fail_cases(t *testing.T) {
	t.Parallel()

	t.Run("step 1: extract parse.Option from task.Option", func(t *testing.T) {
		t.Parallel()

		WithForcedError := func(tsk *Task) error {
			if tsk.Parsed == nil {
				return newError("forced error for testing step 1")
			}

			return nil
		}

		task, err := New("buy milk @shopA", WithForcedError)
		require.Error(t, err)

		require.ErrorContains(t, err, "forced error for testing step 1")
		require.Nil(t, task)
	})

	t.Run("step 2: parse task with control characters disallowed", func(t *testing.T) {
		t.Parallel()

		task, err := New("buy milk\t@shopB")
		require.Error(t, err)

		require.ErrorContains(t, err, "failed to create new task")
		require.ErrorContains(t, err, "invalid task string: task contains control characters")
		require.Nil(t, task)
	})

	t.Run("step 3: apply manipulation options", func(t *testing.T) {
		t.Parallel()

		WithForcedError := func(tsk *Task) error {
			if tsk.Parsed != nil {
				return newError("forced error for testing step 3")
			}

			return nil
		}

		task, err := New("buy milk @shopC", WithForcedError)
		require.Error(t, err)

		require.ErrorContains(t, err, "forced error for testing step 3")
		require.Nil(t, task)
	})

	t.Run("step 4: option re-writes task string to invalid format", func(t *testing.T) {
		t.Parallel()

		WithInvalidRewrite := func(tsk *Task) error {
			if tsk.Parsed != nil {
				invalidText := "invalid\ttask string !!!"
				tsk.SetText(invalidText)
			}

			return nil
		}

		task, err := New("buy milk @shopD", WithInvalidRewrite)
		require.Error(t, err)

		require.ErrorContains(t, err, "failed to apply initial options")
		require.ErrorContains(t, err, "invalid task string: task contains control characters")
		require.Nil(t, task)
	})
}

func TestNew_error_scenarios(t *testing.T) {
	t.Parallel()

	for _, test := range dataNewErrorScenarios {
		t.Run(test.title, func(t *testing.T) {
			t.Parallel()

			task, err := New(test.input)
			require.Error(t, err)

			require.ErrorContains(t, err, test.errorContains)
			require.Nil(t, task)
		})
	}
}

//  Setter functional options
// ============================================================================

// ----------------------------------------------------------------------------
//  WithPriority()
// ----------------------------------------------------------------------------

func TestWithPriority(t *testing.T) {
	t.Parallel()

	for index, test := range dataPriority {
		title := fmt.Sprintf("Test #%d: %s", index+1, test.title)

		t.Run(title, func(t *testing.T) {
			t.Parallel()

			task, err := New(test.task,
				WithPriority(test.priorityVal),
			)

			// Failure case
			if test.shouldFail {
				require.Error(t, err)
				require.ErrorContains(t, err, test.output)

				return
			}

			// Success case
			require.NoError(t, err)

			err = task.Apply()
			require.NoError(t, err)

			require.Equal(t, test.output, task.String())
			require.Equal(t, test.priorityVal, task.Priority())

			if test.additionalTest != nil {
				test.additionalTest(t, task)
			}
		})
	}
}

func TestWithPriority_in_apply(t *testing.T) {
	t.Parallel()

	t.Run("use WithPriority in Apply", func(t *testing.T) {
		t.Parallel()

		task, err := New("buy milk @shopA")
		require.NoError(t, err)

		err = task.Apply(WithPriority("A"))
		require.NoError(t, err)

		require.Equal(t, "(A) buy milk @shopA", task.String())
		require.Equal(t, "A", task.Priority())
	})

	t.Run("multiple priority changes via separate Apply calls", func(t *testing.T) {
		t.Parallel()

		task, err := New("buy milk @shopC")
		require.NoError(t, err)

		err = task.Apply(WithPriority("A"))
		require.NoError(t, err)
		require.Equal(t, "(A) buy milk @shopC", task.String())

		err = task.Apply(WithPriority("B"))
		require.NoError(t, err)

		require.Equal(t, "(B) buy milk @shopC", task.String())

		err = task.Apply(WithPriority("C"))
		require.NoError(t, err)

		require.Equal(t, "(C) buy milk @shopC", task.String())
		require.Equal(t, "C", task.Priority())
	})

	t.Run("multiple priority applications in single Apply call", func(t *testing.T) {
		t.Parallel()

		task, err := New("(A) buy milk @shopD")
		require.NoError(t, err)

		err = task.Apply(
			WithPriority("B"),
			WithPriority("C"), // final priority should be "C"
		)
		require.NoError(t, err)

		require.Equal(t, "(C) buy milk @shopD", task.String())
		require.Equal(t, "C", task.Priority())
	})

	t.Run("failure with bad option before WithPriority", func(t *testing.T) {
		t.Parallel()

		WithBadBehaviorOption := func(tsk2 *Task) error {
			tsk2.SetText(tsk2.String() + "\t") // not allowed control char
			tsk2.isDirty = true

			return nil
		}

		task, err := New("(A) buy milk @shopE")
		require.NoError(t, err)

		err = task.Apply(
			WithBadBehaviorOption,
			WithPriority("B"),
		)
		require.Error(t, err)

		require.ErrorContains(t, err,
			"failed to re-parse dirty task before setting priority")
	})
}

// ----------------------------------------------------------------------------
//  WithoutPriority()
// ----------------------------------------------------------------------------

func TestWithoutPriority(t *testing.T) {
	t.Parallel()

	t.Run("remove priority via Apply", func(t *testing.T) {
		t.Parallel()

		task, err := New("(A) buy milk @shopA")
		require.NoError(t, err)

		err = task.Apply(WithoutPriority())
		require.NoError(t, err)

		require.Equal(t, "buy milk @shopA", task.String())
		require.Empty(t, task.Priority())
	})

	t.Run("remove priority during New", func(t *testing.T) {
		t.Parallel()

		task, err := New("(B) buy milk @shopB", WithoutPriority())
		require.NoError(t, err)

		require.Equal(t, "buy milk @shopB", task.String())
		require.Empty(t, task.Priority())
	})

	t.Run("no-op when priority absent", func(t *testing.T) {
		t.Parallel()

		task, err := New("buy milk @shopC")
		require.NoError(t, err)

		err = task.Apply(WithoutPriority())
		require.NoError(t, err)

		require.Equal(t, "buy milk @shopC", task.String())
		require.Empty(t, task.Priority())
	})

	t.Run("remove priority from completed task", func(t *testing.T) {
		t.Parallel()

		task, err := New("x (C) 2024-01-15 buy milk @shopD")
		require.NoError(t, err)

		err = task.Apply(WithoutPriority())
		require.NoError(t, err)

		require.Equal(t, "x 2024-01-15 buy milk @shopD", task.String())
		require.Empty(t, task.Priority())
		require.True(t, task.IsDone())
	})

	t.Run("add then remove priority in single Apply call", func(t *testing.T) {
		t.Parallel()

		task, err := New("buy milk @shopE")
		require.NoError(t, err)

		err = task.Apply(
			WithPriority("A"), // add
			WithoutPriority(), // remove
		)
		require.NoError(t, err)

		require.Equal(t, "buy milk @shopE", task.String())
		require.Empty(t, task.Priority())
	})

	t.Run("failure with bad option before WithoutPriority", func(t *testing.T) {
		t.Parallel()

		WithBadBehaviorOption := func(tsk2 *Task) error {
			tsk2.SetText(tsk2.String() + "\t") // not allowed control char
			tsk2.isDirty = true

			return nil
		}

		task, err := New("(A) buy milk @shopE")
		require.NoError(t, err)

		err = task.Apply(
			WithBadBehaviorOption,
			WithoutPriority(),
		)
		require.Error(t, err)

		require.ErrorContains(t, err,
			"failed to re-parse dirty task before removing priority")
	})
}

// ----------------------------------------------------------------------------
//  WithContext()
// ----------------------------------------------------------------------------

func TestWithContext(t *testing.T) {
	t.Parallel()

	for index, test := range dataWithContext {
		title := fmt.Sprintf("Test #%d: %s", index+1, test.title)

		t.Run(title+" (during New)", func(t *testing.T) {
			t.Parallel()

			task, err := New(test.taskStr, WithContext(test.addCtx))

			if test.shouldErr {
				require.Error(t, err)
				require.ErrorContains(t, err, test.expectOut,
					"error message does not contain expected text")
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expectOut, task.String())
				require.Equal(t, test.expectCtx, task.Contexts(),
					"parsed contexts do not match expected")
				require.False(t, task.IsDirty(),
					"task should not be marked dirty in New")
			}
		})

		t.Run(title+" (during Apply)", func(t *testing.T) {
			t.Parallel()

			task, err := New(test.taskStr)
			require.NoError(t, err)

			err = task.Apply(WithContext(test.addCtx))

			if test.shouldErr {
				require.Error(t, err)
				require.ErrorContains(t, err, test.expectOut,
					"error message does not contain expected text")
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expectOut, task.String())
				require.Equal(t, test.expectCtx, task.Contexts(),
					"parsed contexts do not match expected")
				require.False(t, task.IsDirty(),
					"task should not be marked dirty after Apply")
			}
		})
	}

	t.Run("failure with bad option before WithContext", func(t *testing.T) {
		t.Parallel()

		WithBadBehaviorOption := func(tsk2 *Task) error {
			tsk2.SetText(tsk2.String() + "\t") // not allowed control char
			tsk2.isDirty = true

			return nil
		}

		task, err := New("buy milk")
		require.NoError(t, err)

		err = task.Apply(
			WithBadBehaviorOption,
			WithContext("shopE"),
		)
		require.Error(t, err)

		require.ErrorContains(t, err,
			"failed to re-parse dirty task before adding context")
	})
}

func TestWithContext_make_comment_line_then_try_to_add_context(t *testing.T) {
	t.Parallel()

	tsk, err := New("buy milk")
	require.NoError(t, err)

	// Custom option that makes the task a comment line and marks it dirty.
	WithUserCustomOption := func(tsk2 *Task) error {
		if tsk2.IsCommentLine() {
			return nil
		}

		tsk2.SetText(spec.PrefixComment.String() + tsk2.String())
		tsk2.isDirty = true

		return nil
	}

	err = tsk.Apply(
		WithUserCustomOption, // makes as comment line
		WithContext("shopX"),
	)
	require.NoError(t, err)

	require.Equal(t, "#buy milk", tsk.String(),
		"if the task is a comment line, it should do nothing when adding context")
}

// ----------------------------------------------------------------------------
//  WithoutContext()
// ----------------------------------------------------------------------------

//nolint:dupl // Similar to TestWithProject is intentional.
func TestWithoutContext(t *testing.T) {
	t.Parallel()

	for index, test := range dataWithoutContext {
		title := fmt.Sprintf("Test #%d: %s", index+1, test.title)

		t.Run(title+" (during New)", func(t *testing.T) {
			t.Parallel()

			tsk, err := New(test.taskStr, WithoutContext(test.removeCtx))

			if test.shouldError {
				require.Error(t, err)
				require.ErrorContains(t, err, test.errContains,
					"error message does not contain expected text")
			} else {
				require.NoError(t, err)

				require.Equal(t, test.expectOut, tsk.String())
				require.Equal(t, test.expectCtx, tsk.Contexts(),
					"parsed contexts do not match expected")
			}
		})

		t.Run(title+" (during Apply)", func(t *testing.T) {
			t.Parallel()

			tsk, err := New(test.taskStr)
			require.NoError(t, err)

			err = tsk.Apply(WithoutContext(test.removeCtx))

			if test.shouldError {
				require.Error(t, err)
				require.ErrorContains(t, err, test.errContains,
					"error message does not contain expected text")
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expectOut, tsk.String())
				require.Equal(t, test.expectCtx, tsk.Contexts(),
					"parsed contexts do not match expected")
			}
		})
	}

	t.Run("make comment line then try to remove context", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk @shopF")
		require.NoError(t, err)

		// Custom option that makes the task a comment line and marks it dirty.
		WithUserCustomOption := func(tsk2 *Task) error {
			if tsk2.IsCommentLine() {
				return nil
			}

			tsk2.SetText(spec.PrefixComment.String() + tsk2.String())

			return nil
		}

		err = tsk.Apply(
			WithUserCustomOption,    // makes as comment line
			WithoutContext("shopX"), // should do nothing
		)
		require.NoError(t, err)

		require.Equal(t, "#buy milk @shopF", tsk.String(),
			"if the task is a comment line, it should do nothing when removing context")
	})

	t.Run("failure with bad option before WithoutContext", func(t *testing.T) {
		t.Parallel()

		WithBadBehaviorOption := func(tsk2 *Task) error {
			tsk2.SetText(tsk2.String() + "\t") // not allowed control char
			tsk2.isDirty = true

			return nil
		}

		task, err := New("buy milk @shopE")
		require.NoError(t, err)

		err = task.Apply(
			WithBadBehaviorOption,
			WithoutContext("shopE"),
		)
		require.Error(t, err)

		require.ErrorContains(t, err,
			"failed to re-parse dirty task before removing context")
	})
}

// ----------------------------------------------------------------------------
//  WithProject()
// ----------------------------------------------------------------------------

//nolint:dupl // Similar to TestWithoutProject is intentional.
func TestWithProject(t *testing.T) {
	t.Parallel()

	for index, test := range dataWithProject {
		title := fmt.Sprintf("Test #%d: %s", index+1, test.title)

		t.Run(title+" (during New)", func(t *testing.T) {
			t.Parallel()

			tsk, err := New(test.taskStr, WithProject(test.addProject))

			if test.shouldError {
				require.Error(t, err)
				require.ErrorContains(t, err, test.expectOut,
					"error message does not contain expected text")
			} else {
				require.NoError(t, err)

				require.Equal(t, test.expectOut, tsk.String())
				require.Equal(t, test.expectPrj, tsk.Projects(),
					"parsed projects do not match expected")
			}
		})

		t.Run(title+" (during Apply)", func(t *testing.T) {
			t.Parallel()

			tsk, err := New(test.taskStr)
			require.NoError(t, err)

			err = tsk.Apply(WithProject(test.addProject))

			if test.shouldError {
				require.Error(t, err)
				require.ErrorContains(t, err, test.expectOut,
					"error message does not contain expected text")
			} else {
				require.NoError(t, err)

				require.Equal(t, test.expectOut, tsk.String())
				require.Equal(t, test.expectPrj, tsk.Projects(),
					"parsed projects do not match expected")
			}
		})
	}

	t.Run("make comment line then try to add project", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk")
		require.NoError(t, err)

		// Custom option that makes the task a comment line and marks it dirty.
		WithUserCustomOption := func(tsk2 *Task) error {
			if tsk2.IsCommentLine() {
				return nil
			}

			tsk2.SetText(spec.PrefixComment.String() + tsk2.String())

			return nil
		}

		err = tsk.Apply(
			WithUserCustomOption, // makes as comment line
			WithProject("shopF"), // should do nothing
		)
		require.NoError(t, err)

		require.Equal(t, "#buy milk", tsk.String(),
			"if the task is a comment line, it should do nothing when adding project")
	})

	t.Run("failure with bad option before WithProject", func(t *testing.T) {
		t.Parallel()

		WithBadBehaviorOption := func(tsk2 *Task) error {
			tsk2.SetText(tsk2.String() + "\t") // not allowed control char

			return nil
		}

		task, err := New("buy milk")
		require.NoError(t, err)

		err = task.Apply(
			WithBadBehaviorOption,
			WithProject("shopG"),
		)
		require.Error(t, err)

		require.ErrorContains(t, err,
			"failed to re-parse dirty task before adding project")
	})
}

// ----------------------------------------------------------------------------
//  WithoutProject()
// ----------------------------------------------------------------------------

//nolint:dupl // Similar to TestWithProject is intentional.
func TestWithoutProject(t *testing.T) {
	t.Parallel()

	for index, test := range dataWithoutProject {
		title := fmt.Sprintf("Test #%d: %s", index+1, test.title)

		t.Run(title+" (during New)", func(t *testing.T) {
			t.Parallel()

			tsk, err := New(test.taskStr, WithoutProject(test.removePrj))

			if test.shouldError {
				require.Error(t, err)
				require.ErrorContains(t, err, test.expectOut,
					"error message does not contain expected text")
			} else {
				require.NoError(t, err)

				require.Equal(t, test.expectOut, tsk.String())
				require.Equal(t, test.expectPrj, tsk.Projects(),
					"parsed Projects do not match expected")
			}
		})

		t.Run(title+" (during Apply)", func(t *testing.T) {
			t.Parallel()

			tsk, err := New(test.taskStr)
			require.NoError(t, err)

			err = tsk.Apply(WithoutProject(test.removePrj))

			if test.shouldError {
				require.Error(t, err)
				require.ErrorContains(t, err, test.expectOut,
					"error message does not contain expected text")
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expectOut, tsk.String())
				require.Equal(t, test.expectPrj, tsk.Projects(),
					"parsed projects do not match expected")
			}
		})
	}

	t.Run("make comment line then try to remove project", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk +shopH")
		require.NoError(t, err)

		// Custom option that makes the task a comment line and marks it dirty.
		WithUserCustomOption := func(tsk2 *Task) error {
			if tsk2.IsCommentLine() {
				return nil
			}

			tsk2.SetText(spec.PrefixComment.String() + tsk2.String())

			return nil
		}

		err = tsk.Apply(
			WithUserCustomOption,    // makes as comment line
			WithoutProject("shopH"), // should do nothing
		)
		require.NoError(t, err)

		require.Equal(t, "#buy milk +shopH", tsk.String(),
			"if the task is a comment line, it should do nothing when removing project")
	})

	t.Run("failure with bad option before WithoutProject", func(t *testing.T) {
		t.Parallel()

		WithBadBehaviorOption := func(tsk2 *Task) error {
			tsk2.SetText(tsk2.String() + "\t") // not allowed control char

			return nil
		}

		task, err := New("buy milk +shopI")
		require.NoError(t, err)

		err = task.Apply(
			WithBadBehaviorOption,
			WithoutProject("shopI"),
		)
		require.Error(t, err)

		require.ErrorContains(t, err, "failed to re-parse dirty task before removing project")
	})
}

// ----------------------------------------------------------------------------
//  WithCompleted()
// ----------------------------------------------------------------------------

func TestWithCompleted(t *testing.T) {
	t.Parallel()

	t.Run("no args - use current date via New", func(t *testing.T) {
		t.Parallel()

		// Test WithCompleted() without args during New()
		// This will use TimeNow() to get current date
		tsk, err := New("buy milk", WithCompleted())
		require.NoError(t, err)

		// Verify task is completed with today's date
		require.True(t, tsk.IsCompleted())
		require.NotEmpty(t, tsk.DateCompleted())
		require.True(t, strings.HasPrefix(tsk.String(), "x "))
	})

	t.Run("no args - use current date via Complete", func(t *testing.T) {
		t.Parallel()

		// Test with Complete() method which also uses TimeNow internally
		tsk, err := New("buy milk")
		require.NoError(t, err)

		err = tsk.Complete() // Uses TimeNow internally
		require.NoError(t, err)

		err = tsk.Apply()
		require.NoError(t, err)

		// Verify task is completed (date will be today's date)
		require.True(t, tsk.IsCompleted())
		require.NotEmpty(t, tsk.DateCompleted())
		require.True(t, strings.HasPrefix(tsk.String(), "x "))
	})

	t.Run("empty string - no date", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk", WithCompleted(""))
		require.NoError(t, err)

		require.Equal(t, "x buy milk", tsk.String())
		require.True(t, tsk.IsCompleted())
		require.Empty(t, tsk.DateCompleted())
	})

	t.Run("specific date provided", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk", WithCompleted("2024-12-25"))
		require.NoError(t, err)

		require.Equal(t, "x 2024-12-25 buy milk", tsk.String())
		require.Equal(t, "2024-12-25", tsk.DateCompleted())
	})

	t.Run("via Apply", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk")
		require.NoError(t, err)

		// Use WithCompleted with specific date instead of relying on TimeNow
		err = tsk.Apply(WithCompleted("2024-07-15"))
		require.NoError(t, err)

		require.Equal(t, "x 2024-07-15 buy milk", tsk.String())
		require.True(t, tsk.IsCompleted())
	})

	t.Run("chained WithCompleted options should avoid duplicate markers", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk")
		require.NoError(t, err)

		err = tsk.Apply(
			WithCompleted("2024-01-01"),
			WithCompleted("2024-02-02"),
		)
		require.NoError(t, err)

		require.Equal(t, "x 2024-02-02 buy milk", tsk.String())
	})
}

func TestWithCompleted_error_cases(t *testing.T) {
	t.Parallel()

	t.Run("error with multiple arguments", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk", WithCompleted("2024-01-01", "2024-12-31"))

		require.Error(t, err)
		require.ErrorContains(t, err, "expects 0 or 1 argument")
		require.Nil(t, tsk)
	})

	t.Run("failure with bad option before WithCompleted", func(t *testing.T) {
		t.Parallel()

		tsk1, err := New("buy milk")
		require.NoError(t, err)

		WithBadBehaviorOption := func(tsk2 *Task) error {
			tsk2.SetText(tsk2.String() + "\t") // not allowed control char
			tsk2.isDirty = true

			return nil
		}

		err = tsk1.Apply(
			WithBadBehaviorOption, // bad option before WithCompleted
			WithCompleted("2024-02-02"),
		)
		require.Error(t, err)

		require.ErrorContains(t, err, "failed to re-parse dirty task")
	})
}

// ----------------------------------------------------------------------------
//  WithKeyValue()
// ----------------------------------------------------------------------------

func TestWithKeyValue(t *testing.T) {
	t.Parallel()

	for index, test := range dataWithKeyValue {
		title := fmt.Sprintf("Test #%d: %s", index+1, test.title)

		t.Run(title+" (during New)", func(t *testing.T) {
			t.Parallel()

			tsk, err := New(test.taskStr, WithKeyValue(test.key, test.value))

			if test.shouldError {
				require.Error(t, err)
				require.ErrorContains(t, err, test.errContains,
					"error message does not contain expected text")
				require.Nil(t, tsk)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expectOut, tsk.String())
				require.Equal(t, test.expectKV, tsk.KeyValues(),
					"parsed key-values do not match expected")
				require.False(t, tsk.IsDirty(),
					"task should not be marked dirty after New")
			}
		})

		t.Run(title+" (during Apply)", func(t *testing.T) {
			t.Parallel()

			tsk, err := New(test.taskStr)
			require.NoError(t, err)

			err = tsk.Apply(WithKeyValue(test.key, test.value))

			if test.shouldError {
				require.Error(t, err)
				require.ErrorContains(t, err, test.errContains,
					"error message does not contain expected text")
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expectOut, tsk.String())
				require.Equal(t, test.expectKV, tsk.KeyValues(),
					"parsed key-values do not match expected")
				require.False(t, tsk.IsDirty(),
					"task should not be marked dirty after Apply")
			}
		})
	}

	t.Run("failure with bad option before WithKeyValue", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk")
		require.NoError(t, err)

		WithBadBehaviorOption := func(tsk2 *Task) error {
			tsk2.SetText(tsk2.String() + "\t")
			tsk2.isDirty = true

			return nil
		}

		err = tsk.Apply(
			WithBadBehaviorOption,
			WithKeyValue("due", "2024-12-25"),
		)
		require.Error(t, err)

		require.ErrorContains(t, err, "failed to re-parse dirty task before setting tag")
	})
}

// ----------------------------------------------------------------------------
//  WithoutKeyValue()
// ----------------------------------------------------------------------------

func TestWithoutKeyValue(t *testing.T) {
	t.Parallel()

	for index, test := range dataWithoutKeyValue {
		title := fmt.Sprintf("Test #%d: %s", index+1, test.title)

		t.Run(title+" (during New)", func(t *testing.T) {
			t.Parallel()

			tsk, err := New(test.taskStr, WithoutKeyValue(test.key))

			if test.shouldError {
				require.Error(t, err)
				require.ErrorContains(t, err, test.errContains,
					"error message does not contain expected text")
				require.Nil(t, tsk)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expectOut, tsk.String())
				require.Equal(t, test.expectKV, tsk.KeyValues(),
					"parsed key-values do not match expected")
				require.False(t, tsk.IsDirty(),
					"task should not be marked dirty after New")
			}
		})

		t.Run(title+" (during Apply)", func(t *testing.T) {
			t.Parallel()

			tsk, err := New(test.taskStr)
			require.NoError(t, err)

			err = tsk.Apply(WithoutKeyValue(test.key))

			if test.shouldError {
				require.Error(t, err)
				require.ErrorContains(t, err, test.errContains,
					"error message does not contain expected text")
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expectOut, tsk.String())
				require.Equal(t, test.expectKV, tsk.KeyValues(),
					"parsed key-values do not match expected")
				require.False(t, tsk.IsDirty(),
					"task should not be marked dirty after Apply")
			}
		})
	}

	t.Run("failure with bad option before WithoutKeyValue", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk due:2024-12-25")
		require.NoError(t, err)

		WithBadBehaviorOption := func(tsk2 *Task) error {
			tsk2.SetText(tsk2.String() + "\t")
			tsk2.isDirty = true

			return nil
		}

		err = tsk.Apply(
			WithBadBehaviorOption,
			WithoutKeyValue("due"),
		)
		require.Error(t, err)

		require.ErrorContains(t, err, "failed to re-parse dirty task before removing tag")
	})
}

// ----------------------------------------------------------------------------
//  WithIncomplete()
// ----------------------------------------------------------------------------

func TestWithIncomplete(t *testing.T) {
	t.Parallel()

	t.Run("reopen completed task during New", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("x 2024-01-01 buy milk", WithIncomplete())
		require.NoError(t, err)

		require.Equal(t, "buy milk", tsk.String())
		require.False(t, tsk.IsCompleted())
	})

	t.Run("via Apply", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("x 2024-06-01 buy milk")
		require.NoError(t, err)

		err = tsk.Apply(WithIncomplete())
		require.NoError(t, err)

		require.Equal(t, "buy milk", tsk.String())
		require.False(t, tsk.IsCompleted())
	})

	t.Run("idempotent - already incomplete", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk", WithIncomplete())
		require.NoError(t, err)

		require.Equal(t, "buy milk", tsk.String())
		require.False(t, tsk.IsCompleted())
	})

	t.Run("Apply with WithCompleted then WithIncomplete should reopen", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk")
		require.NoError(t, err)

		err = tsk.Apply(
			WithCompleted("2024-03-10"),
			WithIncomplete(),
		)
		require.NoError(t, err)

		require.Equal(t, "buy milk", tsk.String())
	})

	t.Run("failure with bad option before WithIncomplete", func(t *testing.T) {
		t.Parallel()

		tsk1, err := New("x buy milk")
		require.NoError(t, err)

		WithBadBehaviorOption := func(tsk2 *Task) error {
			tsk2.SetText(tsk2.String() + "\t") // not allowed control char
			tsk2.isDirty = true

			return nil
		}

		err = tsk1.Apply(
			WithBadBehaviorOption, // bad option before WithIncomplete
			WithIncomplete(),
		)
		require.Error(t, err)

		require.ErrorContains(t, err, "failed to re-parse dirty task")
		require.True(t, tsk1.isDirty,
			"task should remain dirty after failed Apply") // sanity check
	})
}

// ----------------------------------------------------------------------------
//  WithDueDate()
// ----------------------------------------------------------------------------

// TestWithDueDate_multiple_calls verifies that multiple WithDueDate calls
// result in the last call taking precedence.
func TestWithDueDate_multiple_calls(t *testing.T) {
	t.Parallel()

	t.Run("multiple WithDueDate in single Apply - last wins", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk")
		require.NoError(t, err)

		err = tsk.Apply(
			WithDueDate("2024-01-01"),
			WithDueDate("2024-06-15"),
			WithDueDate("2024-12-25"), // last one should win
		)
		require.NoError(t, err)

		require.Equal(t, "buy milk due:2024-12-25", tsk.String())
		require.Equal(t, map[string]string{"due": "2024-12-25"}, tsk.KeyValues())
	})

	t.Run("multiple WithDueDate across separate Apply calls", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk")
		require.NoError(t, err)

		err = tsk.Apply(WithDueDate("2024-01-01"))
		require.NoError(t, err)
		require.Equal(t, "buy milk due:2024-01-01", tsk.String())

		err = tsk.Apply(WithDueDate("2024-06-15"))
		require.NoError(t, err)
		require.Equal(t, "buy milk due:2024-06-15", tsk.String())

		err = tsk.Apply(WithDueDate("2024-12-25"))
		require.NoError(t, err)
		require.Equal(t, "buy milk due:2024-12-25", tsk.String())

		require.Equal(t, map[string]string{"due": "2024-12-25"}, tsk.KeyValues())
	})

	t.Run("WithDueDate during New then update via Apply", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk", WithDueDate("2024-01-01"))
		require.NoError(t, err)
		require.Equal(t, "buy milk due:2024-01-01", tsk.String())

		err = tsk.Apply(WithDueDate("2024-12-25"))
		require.NoError(t, err)
		require.Equal(t, "buy milk due:2024-12-25", tsk.String())
	})

	t.Run("WithDueDate then WithoutDueDate removes due date", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk")
		require.NoError(t, err)

		err = tsk.Apply(
			WithDueDate("2024-12-25"),
			WithoutDueDate(), // should remove it
		)
		require.NoError(t, err)

		require.Equal(t, "buy milk", tsk.String())
		require.Nil(t, tsk.KeyValues())
	})

	t.Run("WithoutDueDate then WithDueDate adds due date", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk due:2024-01-01")
		require.NoError(t, err)

		err = tsk.Apply(
			WithoutDueDate(),          // remove existing
			WithDueDate("2024-12-25"), // add new
		)
		require.NoError(t, err)

		require.Equal(t, "buy milk due:2024-12-25", tsk.String())
		require.Equal(t, map[string]string{"due": "2024-12-25"}, tsk.KeyValues())
	})
}

// ============================================================================
//  Tests for private functions
// ============================================================================

// ----------------------------------------------------------------------------
//  newError()
// ----------------------------------------------------------------------------

func Test_newError(t *testing.T) {
	t.Parallel()

	t.Run("basic error creation", func(t *testing.T) {
		t.Parallel()

		msg := "sample error message"
		err := newError(msg)

		require.Error(t, err,
			"newError() should not return nil")
		require.ErrorContains(t, err, msg,
			"Error message should match input")
	})

	t.Run("formatted error creation", func(t *testing.T) {
		t.Parallel()

		errCode := 42
		err := newError("sample error message with code: %d", errCode)

		require.Error(t, err,
			"newError() should not return nil")
		require.ErrorContains(t, err, "sample error message with code: 42",
			"Error message should match input")
	})
}

// ----------------------------------------------------------------------------
//  wrapError()
// ----------------------------------------------------------------------------

func Test_wrapError(t *testing.T) {
	t.Parallel()

	t.Run("wrap nil error", func(t *testing.T) {
		t.Parallel()

		//nolint:revive // zero-value but explicit for clarity
		var innerErr error = nil

		wrappedMsg := "additional context"
		err := wrapError(innerErr, wrappedMsg)

		require.NoError(t, err,
			"nil error should return nil when wrapped")
	})

	t.Run("wrap with args (formatted)", func(t *testing.T) {
		t.Parallel()

		innerErr := newError("inner error")

		errCode := 42
		err := wrapError(innerErr, "failed with error code: %d", errCode)

		require.Error(t, err,
			"Wrapped error should not be nil")
		require.ErrorContains(t, err, "failed with error code: 42",
			"Wrapped message should contain formatted context")
		require.ErrorIs(t, err, innerErr,
			"Wrapped error should contain the inner error")
	})
}

// ----------------------------------------------------------------------------
//  normalizeTag()
// ----------------------------------------------------------------------------

func Test_normalizeTag(t *testing.T) {
	t.Parallel()

	for index, test := range dataNormalizeTag {
		title := fmt.Sprintf("Test #%d: %s", index+1, test.name)

		t.Run(title, func(t *testing.T) {
			t.Parallel()

			normalized, err := normalizeTag(test.value, test.mark)

			if test.shouldErr {
				require.Error(t, err)
				require.ErrorContains(t, err, test.expectOut,
					"error message does not contain expected text")
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expectOut, normalized)
			}
		})
	}
}

// ----------------------------------------------------------------------------
//  replaceFirst()
// ----------------------------------------------------------------------------

func Test_replaceFirst_DirectUnitTest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		title       string
		input       string
		old         string
		replacement string
		expected    string
	}{
		{
			title:       "basic replacement",
			input:       "(A) buy milk",
			old:         "(A)",
			replacement: "(B)",
			expected:    "(B) buy milk",
		},
		{
			title:       "pattern appears multiple times - only first replaced",
			input:       "(A) check (A) levels",
			old:         "(A)",
			replacement: "(B)",
			expected:    "(B) check (A) levels",
		},
		{
			title:       "pattern not found - no change",
			input:       "buy milk",
			old:         "(A)",
			replacement: "(B)",
			expected:    "buy milk",
		},
		{
			title:       "x space pattern with multiple occurrences",
			input:       "x 2024-01-15 fix x coordinate",
			old:         "x ",
			replacement: "x (A) ",
			expected:    "x (A) 2024-01-15 fix x coordinate",
		},
		{
			title:       "empty old string - strings.Replace behavior",
			input:       "(A) buy milk",
			old:         "",
			replacement: "(B)",
			expected:    "(B)(A) buy milk",
		},
	}

	for _, test := range tests {
		actual := replaceFirst(test.input, test.old, test.replacement)
		require.Equal(t, test.expected, actual)
	}
}

func Test_findSegmentBounds(t *testing.T) {
	t.Parallel()

	//nolint:dupword // duplicate words are intentional
	const targetText = "a a a @context"

	for index, test := range []struct {
		title         string
		searchSegment string
		expectedStart int
		expectedEnd   int
		expectedOk    bool
	}{
		{
			title:         "segment found at start",
			searchSegment: "a",
			expectedStart: 0,
			expectedEnd:   1,
			expectedOk:    true,
		},
		{
			title:         "segment found at end",
			searchSegment: "@context",
			expectedStart: 6,
			expectedEnd:   14,
			expectedOk:    true,
		},
		{
			title:         "segment not found",
			searchSegment: "missing",
			expectedStart: -1,
			expectedEnd:   -1,
			expectedOk:    false,
		},
	} {
		t.Run(fmt.Sprintf("Test #%d: %s", index+1, test.title), func(t *testing.T) {
			t.Parallel()

			start, end, ok := findSegmentBounds(targetText, test.searchSegment)
			require.Equal(t, test.expectedOk, ok)
			require.Equal(t, test.expectedStart, start)
			require.Equal(t, test.expectedEnd, end)
		})
	}
}

// ============================================================================
//  Tests for public methods
// ============================================================================

// ----------------------------------------------------------------------------
//  Task.AppendSegment()
// ----------------------------------------------------------------------------

func TestTask_AppendSegment(t *testing.T) {
	t.Parallel()

	for index, test := range dataAppendSegment {
		title := fmt.Sprintf("Test #%d: %s", index+1, test.title)

		t.Run(title, func(t *testing.T) {
			t.Parallel()

			task, err := New(test.taskStr, test.options...)
			require.NoError(t, err)

			err = task.AppendSegment(test.segment)
			require.NoError(t, err)

			err = task.Apply()
			require.NoError(t, err)

			require.Equal(t, test.expected, task.String())
		})
	}

	t.Run("make comment line then try to append segment", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk")
		require.NoError(t, err)

		// Force the task to be a comment line.
		tsk.SetText(spec.PrefixComment.String() + tsk.String())

		err = tsk.AppendSegment("eggs")
		require.NoError(t, err)

		err = tsk.Apply()
		require.NoError(t, err)

		require.Equal(t, "eggs #buy milk", tsk.String(),
			"if the task is a comment line, it should append segment before the comment")
	})

	t.Run("set bad task then try to append segment", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk")
		require.NoError(t, err)

		// Force the task to be a bad formatted task.
		tsk.SetText(tsk.String() + string([]byte{0x01}))

		err = tsk.AppendSegment("eggs")
		require.Error(t, err)

		require.ErrorContains(t, err, "failed to re-parse dirty task before appending segment")
	})
}

// ----------------------------------------------------------------------------
//  Task.RemovePriority()
// ----------------------------------------------------------------------------

func TestTask_RemovePriority(t *testing.T) {
	t.Parallel()

	t.Run("remove priority from task with priority", func(t *testing.T) {
		t.Parallel()

		task, err := New("(A) buy milk @shopA")
		require.NoError(t, err)

		err = task.RemovePriority()
		require.NoError(t, err)

		err = task.Apply()
		require.NoError(t, err)

		require.Equal(t, "buy milk @shopA", task.String())
		require.Empty(t, task.Priority())
	})

	t.Run("remove priority from task without priority", func(t *testing.T) {
		t.Parallel()

		task, err := New("buy milk @shopB")
		require.NoError(t, err)

		err = task.RemovePriority()
		require.NoError(t, err)

		err = task.Apply()
		require.NoError(t, err)

		require.Equal(t, "buy milk @shopB", task.String())
		require.Empty(t, task.Priority())
	})

	t.Run("remove priority from completed task", func(t *testing.T) {
		t.Parallel()

		task, err := New("x (A) 2024-01-15 buy milk")
		require.NoError(t, err)

		err = task.RemovePriority()
		require.NoError(t, err)

		err = task.Apply()
		require.NoError(t, err)

		require.Equal(t, "x 2024-01-15 buy milk", task.String())
		require.Empty(t, task.Priority())
		require.True(t, task.IsDone())
	})

	t.Run("remove priority preserves other components", func(t *testing.T) {
		t.Parallel()

		task, err := New("(B) buy @grocery +shopping due:2024-12-25")
		require.NoError(t, err)

		err = task.RemovePriority()
		require.NoError(t, err)

		err = task.Apply()
		require.NoError(t, err)

		require.Equal(t, "buy @grocery +shopping due:2024-12-25", task.String())
		require.Empty(t, task.Priority())
		require.Equal(t, []string{"grocery"}, task.Contexts())
		require.Equal(t, []string{"shopping"}, task.Projects())
		require.Equal(t, "2024-12-25", task.KeyValues()["due"])
	})
}

// ----------------------------------------------------------------------------
//  Task.SetPriority()
// ----------------------------------------------------------------------------

func TestTask_SetPriority(t *testing.T) {
	t.Parallel()

	for index, test := range dataPriority {
		title := fmt.Sprintf("Test #%d: %s", index+1, test.title)

		t.Run(title, func(t *testing.T) {
			t.Parallel()

			task, err := New(test.task)
			require.NoError(t, err)

			err = task.SetPriority(test.priorityVal)

			// Failure case
			if test.shouldFail {
				require.Error(t, err)
				require.ErrorContains(t, err, test.output)

				return
			}

			// Success case
			require.NoError(t, err)
			err = task.Apply()
			require.NoError(t, err)

			require.Equal(t, test.output, task.String())
			require.Equal(t, test.priorityVal, task.Priority())

			if test.additionalTest != nil {
				test.additionalTest(t, task)
			}
		})
	}
}

// ----------------------------------------------------------------------------
//  Task.InsertAfter()
// ----------------------------------------------------------------------------

func TestTask_InsertAfter(t *testing.T) {
	t.Parallel()

	for index, test := range dataInsertAfter {
		title := fmt.Sprintf("Test #%d: %s", index+1, test.title)

		t.Run(title, func(t *testing.T) {
			t.Parallel()

			task, err := New(test.taskStr)
			require.NoError(t, err)

			ok := task.InsertAfter(test.targetSeg, test.insertSeg)

			if test.expectOK {
				require.True(t, ok, "expected insertion to succeed")
			} else {
				require.False(t, ok, "expected insertion to fail")
			}

			err = task.Apply()
			require.NoError(t, err)

			require.Equal(t, test.expectOut, task.String())
		})
	}
}

// ----------------------------------------------------------------------------
//  Task.SetTag()
// ----------------------------------------------------------------------------

// TestSetTag_duplicate_keys verifies SetTag behavior when the original task
// contains duplicate keys.
//
// Current behavior: SetTag updates the LAST occurrence of a key. If duplicate
// keys exist, earlier occurrences remain unchanged. This is because KeyValues()
// returns only the last value for each key (Go map behavior).
func TestSetTag_duplicate_keys(t *testing.T) {
	t.Parallel()

	t.Run("KeyValues returns last value when duplicates exist", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk due:2024-11-01 due:2024-12-01")
		require.NoError(t, err)

		// KeyValues() returns the last value due to Go map behavior
		require.Equal(t, map[string]string{"due": "2024-12-01"}, tsk.KeyValues())
	})

	t.Run("SetTag updates last occurrence only", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk due:2024-11-01 due:2024-12-01")
		require.NoError(t, err)

		err = tsk.Apply(WithDueDate("2025-01-01"))
		require.NoError(t, err)

		// Current behavior: only the last "due" is updated
		require.Equal(t, "buy milk due:2024-11-01 due:2025-01-01", tsk.String())
		require.Equal(t, map[string]string{"due": "2025-01-01"}, tsk.KeyValues())
	})

	t.Run("RemoveTag removes last occurrence only", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk due:2024-11-01 due:2024-12-01")
		require.NoError(t, err)

		err = tsk.Apply(WithoutDueDate())
		require.NoError(t, err)

		// Current behavior: only the last "due" is removed
		require.Equal(t, "buy milk due:2024-11-01", tsk.String())
		require.Equal(t, map[string]string{"due": "2024-11-01"}, tsk.KeyValues())
	})

	t.Run("multiple different keys with duplicates", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("task url:http://a.com url:http://b.com tag:foo tag:bar")
		require.NoError(t, err)

		// KeyValues returns last value for each key
		require.Equal(t, map[string]string{
			"url": "http://b.com",
			"tag": "bar",
		}, tsk.KeyValues())
	})
}

// ============================================================================
//  Private Methods
// ============================================================================

// ----------------------------------------------------------------------------
//  Task.isDirty()
// ----------------------------------------------------------------------------

func TestTask_DirtyFlag(t *testing.T) {
	t.Parallel()

	const testString = "x (A) 2016-05-20 2016-04-30 measure space for +chapelShelving @chapel"

	t.Run("new task should not be dirty", func(t *testing.T) {
		t.Parallel()

		task, err := New(testString)
		require.NoError(t, err)
		require.False(t, task.isDirty)
	})

	t.Run("updateText should set dirty flag", func(t *testing.T) {
		t.Parallel()

		task, err := New(testString)
		require.NoError(t, err)

		task.SetText("(A) Updated task text")
		require.True(t, task.isDirty)
	})

	t.Run("Apply with no options should not parse if not dirty", func(t *testing.T) {
		t.Parallel()

		task, err := New(testString)
		require.NoError(t, err)

		initialSegments := len(task.Segments)

		err = task.Apply()
		require.NoError(t, err)

		require.Len(t, task.Segments, initialSegments)
		require.False(t, task.isDirty)
	})

	t.Run("Apply should reparse when dirty and reset flag", func(t *testing.T) {
		t.Parallel()

		task, err := New(testString)
		require.NoError(t, err)

		newText := "(B) Different task @context +project"
		task.SetText(newText)

		require.True(t, task.isDirty)

		err = task.Apply()
		require.NoError(t, err)

		require.Equal(t, newText, task.String())
		require.Equal(t, "B", task.Priority())
		require.False(t, task.isDirty)
	})

	t.Run("multiple text updates before Apply should only parse once", func(t *testing.T) {
		t.Parallel()

		task, err := New(testString)
		require.NoError(t, err)

		task.SetText("(A) First change")
		task.SetText("(B) Second change")
		task.SetText("(C) Final change")

		require.True(t, task.isDirty)

		err = task.Apply()
		require.NoError(t, err)

		require.Equal(t, "(C) Final change", task.String())
		require.Equal(t, "C", task.Priority())
		require.False(t, task.isDirty)
	})
}

// ----------------------------------------------------------------------------
//  Task.Complete()
// ----------------------------------------------------------------------------

func TestTask_Complete(t *testing.T) {
	t.Parallel()

	t.Run("mark incomplete task as complete", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk")
		require.NoError(t, err)

		// Use CompleteWithDate to specify exact date for testing
		err = tsk.CompleteWithDate("2024-06-01")
		require.NoError(t, err)

		err = tsk.Apply()
		require.NoError(t, err)

		require.Equal(t, "x 2024-06-01 buy milk", tsk.String())
		require.True(t, tsk.IsCompleted())
		require.Equal(t, "2024-06-01", tsk.DateCompleted())
	})

	t.Run("idempotent - preserve date when already complete", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("x 2024-06-01 buy milk")
		require.NoError(t, err)

		// Complete() should be idempotent and preserve existing date
		err = tsk.Complete()
		require.NoError(t, err)

		err = tsk.Apply()
		require.NoError(t, err)

		// Should preserve original date
		require.Equal(t, "x 2024-06-01 buy milk", tsk.String())
		require.Equal(t, "2024-06-01", tsk.DateCompleted())
	})
}

// ----------------------------------------------------------------------------
//  Task.CompleteWithDate()
// ----------------------------------------------------------------------------

func TestTask_CompleteWithDate(t *testing.T) {
	t.Parallel()

	t.Run("mark incomplete task with specific date", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk")
		require.NoError(t, err)

		err = tsk.CompleteWithDate("2024-01-01")
		require.NoError(t, err)

		err = tsk.Apply()
		require.NoError(t, err)

		require.Equal(t, "x 2024-01-01 buy milk", tsk.String())
		require.True(t, tsk.IsCompleted())
		require.Equal(t, "2024-01-01", tsk.DateCompleted())
	})

	t.Run("update date when already complete", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("x 2024-01-01 buy milk")
		require.NoError(t, err)

		err = tsk.CompleteWithDate("2024-12-31")
		require.NoError(t, err)

		err = tsk.Apply()
		require.NoError(t, err)

		require.Equal(t, "x 2024-12-31 buy milk", tsk.String())
		require.Equal(t, "2024-12-31", tsk.DateCompleted())
	})

	t.Run("mark complete without date (non-standard)", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk")
		require.NoError(t, err)

		err = tsk.CompleteWithDate("")
		require.NoError(t, err)

		err = tsk.Apply()
		require.NoError(t, err)

		require.Equal(t, "x buy milk", tsk.String())
		require.True(t, tsk.IsCompleted())
		require.Empty(t, tsk.DateCompleted())
	})

	t.Run("add date to completed task without date", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("x buy milk")
		require.NoError(t, err)

		err = tsk.CompleteWithDate("2024-06-15")
		require.NoError(t, err)

		err = tsk.Apply()
		require.NoError(t, err)

		require.Equal(t, "x 2024-06-15 buy milk", tsk.String())
		require.Equal(t, "2024-06-15", tsk.DateCompleted())
	})

	t.Run("complete task with priority - priority should move after x", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("(A) buy milk")
		require.NoError(t, err)

		err = tsk.CompleteWithDate("2024-06-20")
		require.NoError(t, err)

		err = tsk.Apply()
		require.NoError(t, err)

		require.Equal(t, "x (A) 2024-06-20 buy milk", tsk.String())
		require.Equal(t, "A", tsk.Priority())
		require.Equal(t, "2024-06-20", tsk.DateCompleted())
	})

	t.Run("complete task with priority without date", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("(B) buy milk")
		require.NoError(t, err)

		err = tsk.CompleteWithDate("")
		require.NoError(t, err)

		err = tsk.Apply()
		require.NoError(t, err)

		require.Equal(t, "x (B) buy milk", tsk.String())
		require.Equal(t, "B", tsk.Priority())
		require.True(t, tsk.IsCompleted())
	})
}

// ----------------------------------------------------------------------------
//  Task.Reopen()
// ----------------------------------------------------------------------------

func TestTask_Reopen(t *testing.T) {
	t.Parallel()

	t.Run("reopen completed task with date", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("x 2024-01-01 buy milk")
		require.NoError(t, err)

		err = tsk.Reopen()
		require.NoError(t, err)

		err = tsk.Apply()
		require.NoError(t, err)

		require.Equal(t, "buy milk", tsk.String())
		require.False(t, tsk.IsCompleted())
		require.Empty(t, tsk.DateCompleted())
	})

	t.Run("reopen completed task without date", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("x buy milk")
		require.NoError(t, err)

		err = tsk.Reopen()
		require.NoError(t, err)

		err = tsk.Apply()
		require.NoError(t, err)

		require.Equal(t, "buy milk", tsk.String())
		require.False(t, tsk.IsCompleted())
	})

	t.Run("idempotent - already incomplete", func(t *testing.T) {
		t.Parallel()

		tsk, err := New("buy milk")
		require.NoError(t, err)

		err = tsk.Reopen()
		require.NoError(t, err)

		err = tsk.Apply()
		require.NoError(t, err)

		require.Equal(t, "buy milk", tsk.String())
		require.False(t, tsk.IsCompleted())
	})
}

// ============================================================================
//  Tests for issue fixes (Reproduction and fix verification)
// ============================================================================

func Test_Issue21_todo_task_formatting(t *testing.T) {
	t.Parallel()

	const taskTxt = "2025-01-02 testing @inline +project definition for a task @testing time:0s"

	tsk, err := New(taskTxt)
	require.NoError(t, err)

	expect := taskTxt
	actual := tsk.String()

	require.Equal(t, expect, actual,
		"Task string after parsing should match original input")
}

// ============================================================================
//  Helper functions for tests
// ============================================================================
