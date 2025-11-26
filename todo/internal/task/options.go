package task

import (
	"github.com/KEINOS/go-todotxt/todo/internal/parse"
)

// Option is a functional option for configuring or modifying a Task.
// Use with New() or Apply().
type Option func(*Task) error

// ============================================================================
//  Parse Options
// ============================================================================

// WithAllowedCtrlChars allows specified control characters in the task string.
// By default, all control characters (including tabs) are disallowed.
// Use with caution; caller must ensure input safety.
func WithAllowedCtrlChars(allowedCtrlChars []rune) Option {
	return func(t *Task) error {
		if t.Parsed == nil { // pre-parse phase
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

// WithCompleted marks the task as done.
//   - WithCompleted() uses current date
//   - WithCompleted("") marks done without date (non-standard)
//   - WithCompleted("2024-01-01") uses specified date
func WithCompleted(date ...string) Option {
	return func(t *Task) error {
		if t.Parsed == nil { // pre-parse phase
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
		if t.Parsed == nil { // pre-parse phase
			return nil
		}

		return t.AddContext(context)
	}
}

// WithIncomplete marks the task as not done.
func WithIncomplete() Option {
	return func(t *Task) error {
		if t.Parsed == nil { // pre-parse phase
			return nil
		}

		return t.Reopen()
	}
}

// WithoutContext removes a context from the task.
func WithoutContext(context string) Option {
	return func(t *Task) error {
		if t.Parsed == nil { // pre-parse phase
			return nil
		}

		return t.RemoveContext(context)
	}
}

// WithoutPriority removes the task priority.
func WithoutPriority() Option {
	return func(t *Task) error {
		if t.Parsed == nil { // pre-parse phase
			return nil
		}

		return t.RemovePriority()
	}
}

// WithoutProject removes a project from the task.
func WithoutProject(project string) Option {
	return func(t *Task) error {
		if t.Parsed == nil { // pre-parse phase
			return nil
		}

		return t.RemoveProject(project)
	}
}

// WithPriority sets or updates the task priority (A-Z).
// Leading spaces are preserved: " Buy milk" -> "(A)  Buy milk".
func WithPriority(priority string) Option {
	return func(t *Task) error {
		if t.Parsed == nil { // pre-parse phase
			return nil
		}

		return t.SetPriority(priority)
	}
}

// WithProject adds a project tag to the task.
func WithProject(project string) Option {
	return func(t *Task) error {
		if t.Parsed == nil { // pre-parse phase
			return nil
		}

		return t.AddProject(project)
	}
}

// WithKeyValue sets a key-value pair tag (e.g., "due:2024-12-25").
// Updates existing key or appends new one.
func WithKeyValue(key, value string) Option {
	return func(t *Task) error {
		if t.Parsed == nil { // pre-parse phase
			return nil
		}

		return t.SetTag(key, value)
	}
}

// WithoutKeyValue removes a key-value pair tag by key name.
func WithoutKeyValue(key string) Option {
	return func(t *Task) error {
		if t.Parsed == nil { // pre-parse phase
			return nil
		}

		return t.RemoveTag(key)
	}
}

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
