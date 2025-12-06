package spec

// ----------------------------------------------------------------------------
//  Control characters detection in strings
// ----------------------------------------------------------------------------

// Data set for control character tests.
var dataControlChars = []struct {
	input           string
	allowedChars    []rune
	containsCtlChar bool
}{
	// Golden/regular cases
	{
		input:           "",
		allowedChars:    nil,
		containsCtlChar: false,
	},
	{
		input:           "normal text",
		allowedChars:    nil,
		containsCtlChar: false,
	},
	{
		input:           "text with space",
		allowedChars:    nil,
		containsCtlChar: false,
	},
	// allowed cases
	{
		input:           "text with \ttabs\t+project\t@context\t#comment",
		allowedChars:    []rune{'\t'},
		containsCtlChar: false,
	},
	// Cases with control characters (not allowed)
	{
		input:           "text with \x00 null",
		allowedChars:    nil,
		containsCtlChar: true,
	},
	{
		input:           "text with \t tab",
		allowedChars:    nil,
		containsCtlChar: true,
	},
	{
		input:           "text with \n newline",
		allowedChars:    nil,
		containsCtlChar: true,
	},
	{
		input:           "text with \r return",
		allowedChars:    nil,
		containsCtlChar: true,
	},
	{
		input:           "text with \x1b escape",
		allowedChars:    nil,
		containsCtlChar: true,
	},
	{
		input:           "text with \uFEFF BOM",
		allowedChars:    nil,
		containsCtlChar: true,
	},
	{
		input:           "text with \u202E RTL override",
		allowedChars:    nil,
		containsCtlChar: true,
	},
	// edge/niche cases
	// allowed control characters only
	{
		input:           "\t\t\t",
		allowedChars:    []rune{'\t'},
		containsCtlChar: false,
	},
	// mixture of allowed and disallowed control characters
	{
		input:           "text with allowed tab\t but with newline\x00",
		allowedChars:    []rune{'\t'},
		containsCtlChar: true,
	},
}

// ----------------------------------------------------------------------------
//  Printable character detection
// ----------------------------------------------------------------------------

var dataPrintableDetection = []struct {
	title    string
	mark     byte
	expected bool
}{
	// Printable ASCII (0x20-0x7E) - should return true
	{
		title:    "Space (0x20) - lower bound",
		mark:     0x20,
		expected: true,
	},
	{
		title:    "Exclamation (0x21) - lower bound + 1",
		mark:     0x21,
		expected: true,
	},
	{
		title:    "Right brace (0x7D) - upper bound - 1",
		mark:     0x7D,
		expected: true,
	},
	{
		title:    "Tilde (0x7E) - upper bound",
		mark:     0x7E,
		expected: true,
	},
	{
		title:    "Plus sign (+)",
		mark:     '+',
		expected: true,
	},
	{
		title:    "At sign (@)",
		mark:     '@',
		expected: true,
	},
	{
		title:    "Hash (#)",
		mark:     '#',
		expected: true,
	},
	{
		title:    "Lowercase x",
		mark:     'x',
		expected: true,
	},
	{
		title:    "Colon (:)",
		mark:     ':',
		expected: true,
	},
	// Non-printable (control characters) - should return false
	{
		title:    "NUL (0x00)",
		mark:     0x00,
		expected: false,
	},
	{
		title:    "Tab (0x09)",
		mark:     0x09,
		expected: false,
	},
	{
		title:    "Line Feed (0x0A)",
		mark:     0x0A,
		expected: false,
	},
	{
		title:    "Carriage Return (0x0D)",
		mark:     0x0D,
		expected: false,
	},
	{
		title:    "Unit Separator (0x1F) - just below printable",
		mark:     0x1F,
		expected: false,
	},
	{
		title:    "DEL (0x7F) - just above printable",
		mark:     0x7F,
		expected: false,
	},
	{
		title:    "Extended ASCII start (0x80)",
		mark:     0x80,
		expected: false,
	},
	{
		title:    "High byte (0xFF)",
		mark:     0xFF,
		expected: false,
	},
}

// ----------------------------------------------------------------------------
//  Type conversion methods for Mark
// ----------------------------------------------------------------------------

// Test data for Mark type conversion methods.
// All methods return zero/empty values for non-printable characters.
var dataTypeConversion = []struct {
	title        string
	mark         byte
	expectByte   byte
	expectRune   rune
	expectString string
}{
	// Printable ASCII characters (0x20-0x7E)
	{
		title:        "Space (0x20) - lower bound of printable",
		mark:         0x20,
		expectByte:   0x20,
		expectRune:   ' ',
		expectString: " ",
	},
	{
		title:        "Exclamation mark (0x21) - lower bound + 1",
		mark:         0x21,
		expectByte:   0x21,
		expectRune:   '!',
		expectString: "!",
	},
	{
		title:        "Right brace (0x7D) - upper bound - 1",
		mark:         0x7D,
		expectByte:   0x7D,
		expectRune:   '}',
		expectString: "}",
	},
	{
		title:        "Tilde (0x7E) - upper bound of printable",
		mark:         0x7E,
		expectByte:   0x7E,
		expectRune:   '~',
		expectString: "~",
	},
	{
		title:        "Plus sign (project prefix)",
		mark:         '+',
		expectByte:   0x2B,
		expectRune:   '+',
		expectString: "+",
	},
	{
		title:        "At sign (context prefix)",
		mark:         '@',
		expectByte:   0x40,
		expectRune:   '@',
		expectString: "@",
	},
	// Edge cases - outside printable range
	// All methods return zero/empty for non-printable characters
	{
		title:        "Zero value (0x00) - control character",
		mark:         0x00,
		expectByte:   0x00, // InvalidMark
		expectRune:   0,    // InvalidRune
		expectString: "",   // empty string
	},
	{
		title:        "Tab (0x09) - control character",
		mark:         0x09,
		expectByte:   0x00, // InvalidMark
		expectRune:   0,    // InvalidRune
		expectString: "",   // empty string
	},
	{
		title:        "Unit Separator (0x1F) - just below printable",
		mark:         0x1F,
		expectByte:   0x00, // InvalidMark
		expectRune:   0,    // InvalidRune
		expectString: "",   // empty string
	},
	{
		title:        "DEL (0x7F) - control character",
		mark:         0x7F,
		expectByte:   0x00, // InvalidMark
		expectRune:   0,    // InvalidRune
		expectString: "",   // empty string
	},
	{
		title:        "Extended ASCII start (0x80)",
		mark:         0x80,
		expectByte:   0x00, // InvalidMark
		expectRune:   0,    // InvalidRune
		expectString: "",   // empty string
	},
	{
		title:        "Max byte value (0xFF)",
		mark:         0xFF,
		expectByte:   0x00, // InvalidMark
		expectRune:   0,    // InvalidRune
		expectString: "",   // empty string
	},
}
