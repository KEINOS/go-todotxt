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
		// i18n: Unicode characters are valid in project tags
		{"+日本語プロジェクト", true},
		{"+项目", true},     // Chinese
		{"+проект", true}, // Russian
		{"+émoji🚀", true}, // Mixed with emoji
		{"+café", true},   // Accented characters

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
		// i18n: Unicode characters are valid in context tags
		{"@自宅", true},      // Japanese "home"
		{"@办公室", true},     // Chinese "office"
		{"@дом", true},     // Russian "home"
		{"@café", true},    // Accented characters
		{"@работа🏢", true}, // Mixed with emoji

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
		// Edge case: space in value (valid - space is part of the value)
		// Note: In practice, segments are split by whitespace, so this
		// would not occur. But if it did, it's technically valid.
		{"key: value", true},
		{"note: remember to call", true},
		// i18n: Unicode in keys and values
		{"場所:東京", true},        // Japanese key:value
		{"地点:北京", true},        // Chinese key:value
		{"место:дом", true},    // Russian key:value
		{"emoji:🎉party", true}, // Emoji in value
		{"café:latte", true},   // Accented key
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
		actual := test.input.IsTagKeyValue()
		require.Equal(t, test.expected, actual,
			"IsTagKeyValue(%q) should return %v", test.input, test.expected)
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
