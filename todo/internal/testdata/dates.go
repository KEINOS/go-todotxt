package testdata

// Dates contains test cases for date validation.
var Dates = []struct {
	Title   string
	Input   string
	IsDate  bool // valid date format
	IsValid bool // valid date format and valid calendar date
}{
	// Golden test cases
	{"golden case", "2024-01-01", true, true},
	{"golden case", "2024-12-31", true, true},

	// Edge cases
	{"old date", "1234-06-15", true, true},
	{"leap year", "2024-02-29", true, true},
	{"earliest possible date", "0000-01-01", true, true},
	{"far future date", "9999-12-31", true, true},
	{"valid format but invalid date", "2024-99-99", true, false},

	// Invalid format
	{"empty string", "", false, false},
	{"common but invalid (not supported)", "20240101", false, false},
	{"common but invalid (short year not supported)", "24-01-01", false, false},
	{"XXXX-XX-XX format but not a date", "date-is-10", false, false},
	{"XXXX-XX-XX format but in full-width characters", "２０２４－０１－０１", false, false},
	{"completely invalid", "hello-world", false, false},
	{"wrong separator", "2024/01/01", false, false},
	{"missing leading zeros", "2024-1-1", false, false},
	{"contains control characters", "\x00\x01\x02", false, false},
}
