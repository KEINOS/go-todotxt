//nolint:varnamelen // allow short name for simplicity
package parse

import (
	"errors"
	"fmt"
)

var (
	// ErrTaskNotCompleted is returned when trying to retrieve a completion date of done task.
	ErrTaskNotCompleted = newError("task is not marked as completed")
	// ErrNoCompletionDate is returned when no completion date is found in the task.
	ErrNoCompletionDate = newError("no completion date found")
	// ErrTaskWithControlChars is returned when a task string contains control characters.
	ErrTaskWithControlChars = newError("task contains control characters")
)

// wrapError wraps an existing error with a new message.
//
// If the provided error is nil, it returns nil.
func wrapError(err error, message string, a ...any) error {
	if err == nil {
		return nil
	}

	if len(a) > 0 {
		message = fmt.Sprintf(message, a...)
	}

	return fmt.Errorf("%s: %w", message, err)
}

// newError creates a new error with the given formatted message.
//
// If 'a' arguments are provided, it formats the message accordingly.
// Note that even the message is empty, it returns a valid error.
func newError(message string, a ...any) error {
	if len(a) > 0 {
		message = fmt.Sprintf(message, a...)
	}

	//nolint:err113 // using errors.New is acceptable here
	return errors.New(message)
}
