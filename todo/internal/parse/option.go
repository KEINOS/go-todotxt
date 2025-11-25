package parse

import "slices"

// Option defines a functional option for configuring behavior of Parsed.
type Option func(*Parsed)

// WithAllowedCtrlChars is an option to allow control characters in the task string.
//
// By default, any control characters are considered disallowed including tab
// characters. Enabling this option permits control characters in the input. Use
// with caution.
func WithAllowedCtrlChars(allowedCtrlChars []rune) Option {
	return func(p *Parsed) {
		if p.allowedCtrlChars == nil {
			p.allowedCtrlChars = []rune{}
		}

		if allowedCtrlChars != nil {
			p.allowedCtrlChars = append(p.allowedCtrlChars, allowedCtrlChars...)
			p.allowedCtrlChars = slices.Compact(p.allowedCtrlChars) // remove dups
		}
	}
}
