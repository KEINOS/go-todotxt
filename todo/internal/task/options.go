package task

import (
	"github.com/KEINOS/go-todotxt/todo/internal/parse"
)

// Option is a functional option for configuring or modifying a Task.
//
// Options can be used during task creation with New() or for modifications via Apply().
// They follow the functional options pattern for flexible configuration.
type Option func(*Task) error

// ============================================================================
//  Parse Options
// ============================================================================
//  Functional Options Pattern for task parsing.

// WithAllowedCtrlChars is a wrapper of parse.WithAllowedCtrlChars(). It is an
// option to allow control characters in the task string.
//
// By default, any control characters are considered disallowed including tab
// characters. Enabling this option permits control characters in the input.
//
// It is the caller's responsibility to ensure that the input is safe. Use
// with caution.
func WithAllowedCtrlChars(allowedCtrlChars []rune) Option {
	return func(t *Task) error {
		// If Parsed is not yet set, we're in pre-parse phase
		if t.Parsed == nil {
			t.parseOpts = append(t.parseOpts, parse.WithAllowedCtrlChars(allowedCtrlChars))

			return nil
		}

		// Post-parse: do nothing (already applied)
		return nil
	}
}

// ============================================================================
//  Modification Options
// ============================================================================
//  Functional Options Pattern for task modifications.

// WithCompleted marks the task as done with an optional completion date.
//
// Usage patterns:
//   - WithCompleted() - Uses current date
//   - WithCompleted("") - Marks done without completion date (non-standard)
//   - WithCompleted("2024-01-01") - Uses specified date
func WithCompleted(date ...string) Option {
	return func(t *Task) error {
		// When Parsed is nil, it is a pre-parse phase. Skip.
		if t.Parsed == nil {
			return nil
		}

		switch len(date) {
		case 0:
			// No arguments: use current date
			return t.Complete()
		case 1:
			// One argument: use provided date (may be empty)
			return t.CompleteWithDate(date[0])
		default:
			return newError("WithCompleted expects 0 or 1 argument, got %d", len(date))
		}
	}
}

// WithContext adds a context to the task.
func WithContext(context string) Option {
	return func(t *Task) error {
		// When Parsed is nil, it is a pre-parse phase. Skip.
		if t.Parsed == nil {
			return nil
		}

		return t.AddContext(context)
	}
}

// WithIncomplete marks the task as not done.
func WithIncomplete() Option {
	return func(t *Task) error {
		// When Parsed is nil, it is a pre-parse phase. Skip.
		if t.Parsed == nil {
			return nil
		}

		return t.Reopen()
	}
}

// WithoutContext removes a context from the task.
func WithoutContext(context string) Option {
	return func(t *Task) error {
		// When Parsed is nil, it is a pre-parse phase. Skip.
		if t.Parsed == nil {
			return nil
		}

		return t.RemoveContext(context)
	}
}

// WithoutPriority removes the task priority.
func WithoutPriority() Option {
	return func(t *Task) error {
		// When Parsed is nil, it is a pre-parse phase. Skip.
		if t.Parsed == nil {
			return nil
		}

		return t.RemovePriority()
	}
}

// WithoutProject removes a project from the task.
func WithoutProject(project string) Option {
	return func(t *Task) error {
		// When Parsed is nil, it is a pre-parse phase. Skip.
		if t.Parsed == nil {
			return nil
		}

		return t.RemoveProject(project)
	}
}

// WithPriority sets or updates the task priority.
//
// This is a convenience wrapper around SetPriority for use with New() or Apply().
// Priority must be a single uppercase letter A-Z.
//
// Example:
//
//	task.Apply(WithPriority("A"))
//
// Note that when adding a priority, leading spaces in the task text are not
// trimmed. The priority is simply inserted before any leading spaces.
//
// Example:
//
//	" Buy milk" --> WithPriority("A") --> "(A)  Buy milk"
func WithPriority(priority string) Option {
	return func(t *Task) error {
		// When Parsed is nil, it is a pre-parse phase. Skip.
		if t.Parsed == nil {
			return nil
		}

		return t.SetPriority(priority)
	}
}

// WithProject adds a project tag to the task.
func WithProject(project string) Option {
	return func(t *Task) error {
		// When Parsed is nil, it is a pre-parse phase. Skip.
		if t.Parsed == nil {
			return nil
		}

		return t.AddProject(project)
	}
}

// // WithKeyValue sets a key-value pair tag.
// func WithKeyValue(key, value string) Option {
// 	return func(_ *Task) error {
// 		// to be implemented
// 		return nil
// 		// return t.SetTag(key, value)
// 	}
// }

// // WithoutKeyValue removes the key-value pair tag.
// func WithoutKeyValue(key string) Option {
// 	return func(_ *Task) error {
// 		// to be implemented
// 		return nil
// 		// return t.RemoveTag(key)
// 	}
// }

// // WithDueDate sets the due date.
// func WithDueDate(date string) Option {
// 	return func(_ *Task) error {
// 		// to be implemented
// 		return nil
// 		// return t.SetDueDate(date)
// 	}
// }

// // WithoutDueDate removes the due date.
// func WithoutDueDate() Option {
// 	return func(_ *Task) error {
// 		// to be implemented
// 		return nil
// 		// return t.RemoveDueDate()
// 	}
// }

// // WithInlineComment adds or updates an inline comment.
// func WithInlineComment(comment string) Option {
// 	return func(_ *Task) error {
// 		// to be implemented
// 		return nil
// 		// return t.AddInlineComment(comment)
// 	}
// }

// // WithoutInlineComment removes the inline comment.
// func WithoutInlineComment() Option {
// 	return func(_ *Task) error {
// 		// to be implemented
// 		return nil
// 		// return t.RemoveInlineComment()
// 	}
// }

// // WithDescription sets the description.
// func WithDescription(description string) Option {
// 	return func(_ *Task) error {
// 		// to be implemented
// 		return nil
// 		// return t.SetDescription(description)
// 	}
// }

// // WithoutDescription removes the description.
// func WithoutDescription() Option {
// 	return func(_ *Task) error {
// 		// to be implemented
// 		return nil
// 		// return t.RemoveDescription()
// 	}
// }
