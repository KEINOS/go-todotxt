package segment

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// ============================================================================
//  Tests for Segment methods
// ============================================================================

// ----------------------------------------------------------------------------
//  Tests for Segment.IsMarkCompletion()
// ----------------------------------------------------------------------------

func TestSegment_IsMarkCompletion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    Segment
		expected bool
	}{
		{"x", true},

		{"X", false},
		{"", false},
		{"y", false},
		{"xx", false},
	}

	for _, test := range tests {
		actual := test.input.IsMarkCompletion()
		require.Equal(t, test.expected, actual,
			"IsMarkCompletion(%q) should return %v", test.input, test.expected)
	}
}

// ----------------------------------------------------------------------------
//  Tests for segment.IsMarkPriority()
// ----------------------------------------------------------------------------

func TestSegment_IsMarkPriority(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    Segment
		expected bool
	}{
		{"(A)", true},
		{"(Z)", true},

		{"(a)", false},
		{"(1)", false},
		{"(A", false},
		{"A)", false},
		{"(AA)", false},
		{"", false},
		{"(A) ", false},
	}

	for _, test := range tests {
		actual := test.input.IsMarkPriority()
		require.Equal(t, test.expected, actual,
			"IsMarkPriority(%q) should return %v", test.input, test.expected)
	}
}

// ----------------------------------------------------------------------------
//  Tests for segment.IsDate()
// ----------------------------------------------------------------------------

func TestSegment_IsDate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    Segment
		expected bool
	}{
		{"2016-05-20", true},
		{"0000-01-01", true},
		{"9999-12-31", true},

		{"", false},
		{"2016-05-2", false},   // too short
		{"2016-05-200", false}, // too long
		{"2016/05/20", false},  // wrong separator
		{"2016-05-2A", false},  // non-digit

		{"2016-13-01", true}, // invalid month, but format ok
		{"2016-05-32", true}, // invalid day, but format ok
	}

	for _, test := range tests {
		actual := test.input.IsDate()
		require.Equal(t, test.expected, actual,
			"IsDate(%q) should return %v", test.input, test.expected)
	}
}

// ----------------------------------------------------------------------------
//  Tests for segment.IsTagProject()
// ----------------------------------------------------------------------------

func TestSegment_IsTagProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    Segment
		expected bool
	}{
		{"+project", true},
		{"+Project", true},
		{"+123", true},

		{"+", false},
		{"project", false},
		{"@project", false},
		{"", false},
	}

	for _, test := range tests {
		actual := test.input.IsTagProject()
		require.Equal(t, test.expected, actual,
			"IsTagProject(%q) should return %v", test.input, test.expected)
	}
}

// ----------------------------------------------------------------------------
//  Tests for segment.IsTagContext()
// ----------------------------------------------------------------------------

func TestSegment_IsTagContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    Segment
		expected bool
	}{
		{"@context", true},
		{"@Context", true},
		{"@123", true},

		{"@", false},
		{"context", false},
		{"+context", false},
		{"", false},
	}

	for _, test := range tests {
		actual := test.input.IsTagContext()
		require.Equal(t, test.expected, actual,
			"IsTagContext(%q) should return %v", test.input, test.expected)
	}
}

// ----------------------------------------------------------------------------
//  Tests for segment.IsKeyValue()
// ----------------------------------------------------------------------------

func TestSegment_IsKeyValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    Segment
		expected bool
	}{
		// Valid key:value pairs
		{"key:value", true},
		{"due:2016-05-20", true},
		// Multiple colons allowed in value (real-world use cases)
		{"url:https://github.com/KEINOS/go-todotxt", true},
		{"host:192.168.1.100:8080", true},
		{"endpoint:http://api.example.com:3000/v1/users", true},
		{"time:10:00-11:30", true},
		{"key:value:extra", true},
		// Invalid formats
		{"key:", false},     // trailing colon
		{":value", false},   // leading colon
		{"keyvalue", false}, // no colon
		{":", false},        // colon only
		{"::", false},       // multiple colons only
		{"::value", false},  // leading colons
		{"key::", false},    // trailing colons
		{"", false},         // empty
	}

	for _, test := range tests {
		actual := test.input.IsKeyValue()
		require.Equal(t, test.expected, actual,
			"IsKeyValue(%q) should return %v", test.input, test.expected)
	}
}

// ----------------------------------------------------------------------------
//  Tests for segment.IsComment()
// ----------------------------------------------------------------------------

func TestSegment_IsComment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    Segment
		expected bool
	}{
		{"# comment", true},
		{"#", true},
		{"#comment", true},

		{"comment", false},
		{"buy#milk", false},
		{"", false},
	}

	for _, test := range tests {
		actual := test.input.IsComment()
		require.Equal(t, test.expected, actual,
			"IsComment(%q) should return %v", test.input, test.expected)
	}
}

// ----------------------------------------------------------------------------
//  Tests for segment.IsPlainText()
// ----------------------------------------------------------------------------

func TestSegment_IsPlainText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    Segment
		expected bool
	}{
		{"plain text", true},
		{"task", true},
		{"", true}, // empty is plain text

		{"x", false},          // completion mark
		{"(A)", false},        // priority
		{"2016-05-20", false}, // date
		{"+project", false},   // project
		{"@context", false},   // context
		{"key:value", false},  // key-value
		{"# comment", false},  // comment
	}

	for _, test := range tests {
		actual := test.input.IsPlainText()
		require.Equal(t, test.expected, actual,
			"IsPlainText(%q) should return %v", test.input, test.expected)
	}
}
