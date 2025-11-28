package parse

// Option defines a functional option for configuring behavior of Parsed.
type Option func(*Parsed)

// WithAllowedCtrlChars is an option to allow control characters in the task string.
// Such as tab (`\t`) character.
//
// By default, any control characters are considered disallowed. Enabling this option
// permits control characters in the input.
//
// Control characters such as line breaks may affect parsing. Use with caution.
func WithAllowedCtrlChars(allowedCtrlChars []rune) Option {
	return func(p *Parsed) {
		p.AllowCtrlChars(allowedCtrlChars)
	}
}
