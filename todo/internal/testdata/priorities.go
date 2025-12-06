package testdata

// Priorities contains test cases for priority validation. Both for priority
// letters (A-Z) and priority marks (e.g., (A)).
var Priorities = []struct {
	Title            string
	Input            string
	IsPriorityLetter bool // valid priority letter A-Z
	IsPriorityMark   bool // valid priority mark (e.g., (A))
}{
	// Golden test cases
	{"letter A (no parentheses)", "A", true, false},
	{"letter Z (no parentheses)", "Z", true, false},
	{"priority mark (A)", "(A)", false, true},
	{"priority mark (Z)", "(Z)", false, true},

	// Invalid cases
	{"empty string", "", false, false},
	{"lowercase 'a' (no parentheses)", "a", false, false},
	{"digit '1' (no parentheses)", "1", false, false},
	{"too long 'AA'", "AA", false, false},
	{"unicode letter 'Ä'", "Ä", false, false},

	{"priority mark empty ()", "()", false, false},
	{"priority mark lowercase (a)", "(a)", false, false},
	{"priority mark digit (1)", "(1)", false, false},
	{"priority mark too long (AA)", "(AA)", false, false},
	{"priority mark unicode (Ä)", "(Ä)", false, false},

	{"missing opening parenthesis A)", "A)", false, false},
	{"missing closing parenthesis (A", "(A", false, false},
	{"extra space '(A) '", "(A) ", false, false},
	{"contains control characters", "(\x01)", false, false},
}
