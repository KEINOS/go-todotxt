package todo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
//  Issue Fixes Tests
// ----------------------------------------------------------------------------

// Test_Issue21 tests the fix for issue #21 which once parsed the segments loses
// the original task string representation.
func Test_Issue21(t *testing.T) {
	t.Parallel()

	t.Skipf("[SKIPPED] %s: Temporarily skipped until implementation is completed", t.Name())

	const taskInput = "2025-01-02 testing @inline +project definition for a task @testing time:0s"

	task, err := ParseTask(taskInput)
	require.NoError(t, err,
		"failed to parse task for issue #21 test")

	expect := taskInput
	actual := task.String()
	require.Equal(t, expect, actual,
		"issue #21: task string mismatch")
}
