//nolint:varnamelen // allow short name for simplicity
package task

import (
	"errors"
	"fmt"
)

var (
	// ErrNotImplemented is returned when a feature is not yet implemented.
	ErrNotImplemented = newError("not implemented yet")
	// ErrTaskIsNil is returned when a nil task object is encountered.
	ErrTaskIsNil = newError("task object is nil")
	// ErrValIsEmpty is returned when a value is required but is empty.
	ErrValIsEmpty = newError("value is empty")
	// ErrValWithSpace is returned when a value contains white spaces but shouldn't.
	ErrValWithSpace = newError("value cannot contain white spaces")
	// ErrEmptyContextValue is returned when a context value is invalid or empty.
	ErrEmptyContextValue = newError("invalid context: value cannot be empty")
	// ErrEmptyProjectValue is returned when a project value is invalid or empty.
	ErrEmptyProjectValue = newError("invalid project: value cannot be empty")
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
