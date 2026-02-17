package todoparse

import (
	"fmt"
	"testing"

	"github.com/KEINOS/go-todotxt/todo/internal/segment"
	"github.com/stretchr/testify/require"
)

// ============================================================================
//  Test Structure Overview
//    1. Tests for public functions
//    2. Tests for private functions
//    3. Tests for public methods
//	  4. Tests for private methods
// ============================================================================

// ============================================================================
//  Tests for public functions
// ============================================================================

// ----------------------------------------------------------------------------
//  FromTaskString()
// ----------------------------------------------------------------------------

func TestFromTaskString_parsing_failure(t *testing.T) {
	t.Parallel()

	for index, test := range dataControlChars {
		title := fmt.Sprintf("TestCase #%d", index+1)

		t.Run(title, func(t *testing.T) {
			t.Parallel()

			allowedCtrlChars := test.allowedChars

			parsed, err := FromTaskString(test.input,
				WithAllowedCtrlChars(allowedCtrlChars))

			if test.containsCtlChar {
				require.Error(t, err,
					"Task string with control characters should return an error. Input: '%q'", test.input)
				require.Nil(t, parsed,
					"Parsed result should be nil on error for input: '%q'", test.input)
			} else {
				require.NoError(t, err,
					"Task string without control characters should not return an error. Input: '%q'", test.input)
				require.NotNil(t, parsed,
					"Parsed result should not be nil for valid input: '%q'", test.input)
				require.ElementsMatch(t, test.expectSegments, parsed.Segments,
					"Parsed segments should match expected for input: '%q'", test.input)
			}
		})
	}
}

func TestFromTaskString_caching(t *testing.T) {
	t.Parallel()

	testString := "x (A) 2016-05-20 2016-04-30 measure space for " +
		"+chapelShelving @chapel due:2016-05-30 location:mainOffice"

	parsed, err := FromTaskString(testString)
	require.NoErrorf(t, err,
		"Failed to parse input: %s", testString)

	// First calls to populate cache
	isDoneFirst := parsed.IsDone()
	isCommentFirst := parsed.IsCommentLine()
	headerIdxFirst := parsed.consumeHeadParts()
	keyValuesFirst := parsed.KeyValues()

	// Subsequent calls to verify cached results
	isDoneSecond := parsed.IsDone()
	isCommentSecond := parsed.IsCommentLine()
	headerIdxSecond := parsed.consumeHeadParts()
	keyValuesSecond := parsed.KeyValues()

	require.Equal(t, isDoneFirst, isDoneSecond,
		"IsDone() should return consistent results with caching")
	require.Equal(t, isCommentFirst, isCommentSecond,
		"IsCommentLine() should return consistent results with caching")
	require.Equal(t, headerIdxFirst, headerIdxSecond,
		"consumeHeadParts() should return consistent results with caching")
	require.Equal(t, keyValuesFirst, keyValuesSecond,
		"KeyValues() should return consistent results with caching")
}

// ============================================================================
//  Tests for private functions
// ============================================================================

// ----------------------------------------------------------------------------
//  findInlineCommentPos()
// ----------------------------------------------------------------------------

func Test_findInlineCommentPos(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		input        string
		allowedRunes []rune
		expected     int
	}{
		{
			name:         "comment at line start",
			input:        "#note",
			allowedRunes: nil,
			expected:     0,
		},
		{
			name:         "comment line but with leading spaces",
			input:        "    # comment line",
			allowedRunes: nil,
			expected:     4,
		},
		{
			name:         "comment separated by space",
			input:        "milk #note",
			allowedRunes: nil,
			expected:     len("milk "),
		},
		{
			name:         "hash embedded without whitespace",
			input:        "milk#note",
			allowedRunes: nil,
			expected:     -1,
		},
		{
			name:         "comment separated by tab when allowed",
			input:        "milk\t#note",
			allowedRunes: []rune{'\t'},
			expected:     len("milk\t"),
		},
		{
			name:         "tab not allowed keeps hash inside word",
			input:        "milk\t#note",
			allowedRunes: nil,
			expected:     -1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			parsed := new(Parsed)
			parsed.originalText = test.input
			parsed.allowedCtrlChars = test.allowedRunes

			actual := findInlineCommentPos(parsed)

			require.Equalf(t, test.expected, actual,
				"findInlineCommentPos() should return %d for input %q",
				test.expected, test.input)
		})
	}
}

// ----------------------------------------------------------------------------
//  newError()
// ----------------------------------------------------------------------------

func Test_newError(t *testing.T) {
	t.Parallel()

	t.Run("basic error creation", func(t *testing.T) {
		t.Parallel()

		msg := "sample error message"
		err := newError(msg)

		require.Error(t, err,
			"newError() should not return nil")
		require.ErrorContains(t, err, msg,
			"Error message should match input")
	})

	t.Run("formatted error creation", func(t *testing.T) {
		t.Parallel()

		errCode := 42
		err := newError("sample error message with code: %d", errCode)

		require.Error(t, err,
			"newError() should not return nil")
		require.ErrorContains(t, err, "sample error message with code: 42",
			"Error message should match input")
	})
}

// ----------------------------------------------------------------------------
//  wrapError()
// ----------------------------------------------------------------------------

func Test_wrapError(t *testing.T) {
	t.Parallel()

	t.Run("wrap nil error", func(t *testing.T) {
		t.Parallel()

		//nolint:revive // zero-value but explicit for clarity
		var innerErr error = nil

		wrappedMsg := "additional context"
		err := wrapError(innerErr, wrappedMsg)

		require.NoError(t, err,
			"nil error should return nil when wrapped")
	})

	t.Run("wrap with args (formatted)", func(t *testing.T) {
		t.Parallel()

		innerErr := newError("inner error")

		errCode := 42
		err := wrapError(innerErr, "failed with error code: %d", errCode)

		require.Error(t, err,
			"Wrapped error should not be nil")
		require.ErrorContains(t, err, "failed with error code: 42",
			"Wrapped message should contain formatted context")
		require.ErrorIs(t, err, innerErr,
			"Wrapped error should contain the inner error")
	})
}

// ============================================================================
//  Tests for public methods
// ============================================================================

// ----------------------------------------------------------------------------
//  Parsed.Parse()
// ----------------------------------------------------------------------------

func TestParsed_Parse_clears_caches(t *testing.T) {
	t.Parallel()

	taskBefore := "buy a milk @grocery +errands due:2023-10-01"
	taskAfter := "x 2023-10-02 buy a milk @grocery +errands due:2023-10-01"

	parsed, err := FromTaskString(taskBefore)
	require.NoErrorf(t, err,
		"Failed to parse input: %s", taskBefore)

	require.False(t, parsed.IsDone(),
		"IsDone() should be false for uncompleted task")
	require.Equal(t, taskBefore, parsed.String(),
		"String() should return the original task string before re-parsing")

	// Re-parse
	err = parsed.Parse(taskAfter)
	require.NoErrorf(t, err,
		"Failed to re-parse input: %s", taskAfter)

	require.True(t, parsed.IsDone(),
		"IsDone() should be true after re-parsing completed task")

	// Verify originalText
	require.Equal(t, taskAfter, parsed.String(),
		"String() should return the updated task string after re-parsing")
	// Verify DateCompleted
	expectDoneDate := "2023-10-02"
	actualDoneDate := parsed.DateCompleted()
	require.Equalf(t, expectDoneDate, actualDoneDate,
		"DateCompleted() should return %q after re-parsing, got %q",
		expectDoneDate, actualDoneDate)
}

// ----------------------------------------------------------------------------
//  Parsed.HasInlineComment()
/// ----------------------------------------------------------------------------

func TestParsed_HasInlineComment(t *testing.T) {
	t.Parallel()

	for index, test := range dataValidVarious {
		parsed, err := FromTaskString(test.input)
		require.NoErrorf(t, err,
			"Failed to parse input for test case %d: %s", index+1, test.title)

		expect := test.expected.HasInlineComment
		actual := parsed.HasInlineComment()

		require.Equalf(t, expect, actual,
			"Test case #%d (%s): HasInlineComment() should return %v for input '%s'",
			index+1, test.title, expect, test.input)
	}
}

func TestParsed_HasInlineComment_hashSpacing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		input        string
		allowedRunes []rune
		expected     bool
	}{
		{
			name:         "hash attached to word",
			input:        "milk#note",
			allowedRunes: []rune{},
			expected:     false,
		},
		{
			name:         "hash separated by space",
			input:        "milk #note",
			allowedRunes: []rune{},
			expected:     true,
		},
		{
			name:         "hash with leading space",
			input:        " #note milk",
			allowedRunes: []rune{},
			expected:     true,
		},
		{
			name:         "hash at beginning",
			input:        "#note milk",
			allowedRunes: []rune{},
			expected:     true,
		},
		{
			name:         "hash separated by tab",
			input:        "milk\t#note",
			allowedRunes: []rune{'\t'},
			expected:     true,
		},
	}

	for index, test := range tests {
		parsed, err := FromTaskString(test.input,
			WithAllowedCtrlChars(test.allowedRunes),
		)
		t.Logf("Test case #%d (%s): Segments: %v, Parsed: %#v", index+1, test.name, parsed.Segments, parsed)

		require.NoErrorf(t, err,
			"Failed to parse input: %s", test.input)

		if test.expected {
			require.Truef(t, parsed.HasInlineComment(),
				"HasInlineComment() should return true for input '%s'", test.input)
		} else {
			require.Falsef(t, parsed.HasInlineComment(),
				"HasInlineComment() should return false for input '%s'", test.input)
		}
	}
}

func TestParsed_HasInlineComment_cache(t *testing.T) {
	t.Parallel()

	input := "Buy milk #remember"

	parsed, err := FromTaskString(input)
	require.NoErrorf(t, err, "Failed to parse input: %s", input)

	require.Nil(t, parsed.inlineComment, "inlineComment cache should be nil before first call")

	first := parsed.HasInlineComment()
	require.True(t, first, "HasInlineComment() should detect inline comment on first call")

	ptrFirst := parsed.inlineComment
	require.NotNil(t, ptrFirst, "inlineComment cache pointer should not be nil after first call")
	require.Equal(t, "#remember", *ptrFirst, "inline comment cache should store comment text")

	second := parsed.HasInlineComment()
	require.True(t, second, "HasInlineComment() should remain true on subsequent call")

	require.Equal(t, ptrFirst, parsed.inlineComment,
		"inlineComment cache pointer should remain stable between calls")
}

func TestParsed_HasInlineComment_nonSeparatedWhitespace(t *testing.T) {
	t.Parallel()

	// Not ordinary way to segmentize
	parsed := new(Parsed)
	parsed.Segments = segment.Segments{
		segment.Segment("milk"),
		segment.Segment("#note"),
	}
	parsed.originalText = "milk\n#note"

	result := parsed.HasInlineComment()
	require.False(t, result, "HasInlineComment() should return false when comment is not separated")

	require.NotNil(t, parsed.inlineComment, "inlineComment cache should be initialized")
	require.Empty(t, *parsed.inlineComment, "inlineComment cache should store empty string when comment not found")
}

func TestParsed_InlineCommentPos(t *testing.T) {
	t.Parallel()

	for index, test := range dataInlineCommentPos {
		t.Run(test.title, func(t *testing.T) {
			t.Parallel()

			parsed, err := FromTaskString(test.input, WithAllowedCtrlChars(test.allowedRunes))
			require.NoErrorf(t, err,
				"Failed to parse input for test case %d: %s", index+1, test.title)

			actual := parsed.InlineCommentPos()
			require.Equalf(t, test.expectedPos, actual,
				"Test case #%d (%s): InlineCommentPos() should return %d for input '%s'",
				index+1, test.title, test.expectedPos, test.input)

			expectHas := test.expectedPos >= 0
			require.Equalf(t, expectHas, parsed.HasInlineComment(),
				"Test case #%d (%s): HasInlineComment() should match InlineCommentPos() expectation for input '%s'",
				index+1, test.title, test.input)
		})
	}
}

func TestParsed_InlineCommentPos_cacheReuse(t *testing.T) {
	t.Parallel()

	input := "Archive logs # keep 30 days"

	parsed, err := FromTaskString(input)
	require.NoErrorf(t, err, "Failed to parse input: %s", input)

	require.Nil(t, parsed.indexComment,
		"indexComment cache should be nil before InlineCommentPos() is invoked")

	first := parsed.InlineCommentPos()
	require.Equal(t, len("Archive logs "), first,
		"InlineCommentPos() should detect inline comment position on first call")

	ptrFirst := parsed.indexComment
	require.NotNil(t, ptrFirst,
		"indexComment cache pointer should not be nil after InlineCommentPos() call")

	second := parsed.InlineCommentPos()
	require.Equal(t, first, second,
		"InlineCommentPos() should return consistent position on subsequent calls")

	require.Equal(t, ptrFirst, parsed.indexComment,
		"indexComment cache pointer should remain stable between calls")
}

// ----------------------------------------------------------------------------
//  Parsed.IsCommentLine()
/// ----------------------------------------------------------------------------

func TestParsed_IsCommentLine(t *testing.T) {
	t.Parallel()

	for index, test := range dataValidVarious {
		parsed, err := FromTaskString(test.input)
		require.NoErrorf(t, err,
			"Failed to parse input for test case %d: %s", index+1, test.title)

		expect := test.expected.IsCommentLine
		actual := parsed.IsCommentLine()

		require.Equalf(t, expect, actual,
			"Test case #%d (%s): IsCommentLine() should return %v for input '%s'",
			index+1, test.title, expect, test.input)
	}
}

func TestParsed_IsCommentLine_cache(t *testing.T) {
	t.Parallel()

	testString := "x (A) 2016-05-20 2016-04-30 measure space for " +
		"+chapelShelving @chapel due:2016-05-30 location:mainOffice"

	parsed, err := FromTaskString(testString)
	require.NoErrorf(t, err,
		"Failed to parse input: %s", testString)

	require.Nil(t, parsed.isCommentLine, "isCommentLine cache should be nil before first call")

	// First calls to populate cache
	isCommentFirst := parsed.IsCommentLine()
	ptrFirst := parsed.isCommentLine
	require.NotNil(t, parsed.isCommentLine, "isCommentLine cache should not be nil after first call")

	// Subsequent calls to verify cached results
	isCommentSecond := parsed.IsCommentLine()
	ptrSecond := parsed.isCommentLine

	require.Equal(t, ptrFirst, ptrSecond,
		"isCommentLine cache pointer should remain the same after subsequent calls")

	require.Equal(t, isCommentFirst, isCommentSecond,
		"IsCommentLine() should return consistent results with caching")
}

// ----------------------------------------------------------------------------
//  Parsed.IsDone()
// ----------------------------------------------------------------------------

func TestParsed_IsDone(t *testing.T) {
	t.Parallel()

	for index, test := range dataValidVarious {
		parsed, err := FromTaskString(test.input)
		require.NoErrorf(t, err,
			"Failed to parse input for test case %d: %s", index+1, test.title)

		expect := test.expected.IsDone
		actual := parsed.IsDone()

		require.Equalf(t, expect, actual,
			"Test case #%d (%s): IsDone() should return %v for input '%s'",
			index+1, test.title, expect, test.input)
	}
}

// ----------------------------------------------------------------------------
//  Parsed.Priority()
// ----------------------------------------------------------------------------

func TestParsed_Priority(t *testing.T) {
	t.Parallel()

	for index, test := range dataValidVarious {
		parsed, err := FromTaskString(test.input)
		require.NoErrorf(t, err,
			"Failed to parse input for test case %d: %s", index+1, test.title)

		expect := test.expected.Priority
		actual := parsed.Priority()

		require.Equalf(t, expect, actual,
			"Test case #%d (%s): Priority() should return %q for input '%s'",
			index+1, test.title, expect, test.input)
	}
}

func TestParsed_Priority_edge_and_invalid_cases(t *testing.T) {
	t.Parallel()

	for index, test := range dataPriority {
		parsed, err := FromTaskString(test.input)
		require.NoErrorf(t, err,
			"Failed to parse input for test case %d: %s", index+1, test.title)

		expect := test.expected
		actual := parsed.Priority()

		require.Equalf(t, expect, actual,
			"Test case #%d (%s): Priority() should return %q for input '%s'",
			index+1, test.title, expect, test.input)
	}
}

// ----------------------------------------------------------------------------
//  Parsed.DateCompleted()
// ----------------------------------------------------------------------------

func TestParsed_DateCompleted(t *testing.T) {
	t.Parallel()

	for index, test := range dataValidVarious {
		parsed, err := FromTaskString(test.input)
		require.NoErrorf(t, err,
			"Failed to parse input for test case %d: %s", index+1, test.title)

		expect := test.expected.DateCompleted
		actual := parsed.DateCompleted()

		require.Equalf(t, expect, actual,
			"Test case #%d (%s): DateCompleted() should return %q for input '%s'",
			index+1, test.title, expect, test.input)
	}
}

// ----------------------------------------------------------------------------
//  Parsed.DateCreated()
// ----------------------------------------------------------------------------

func TestParsed_DateCreated(t *testing.T) {
	t.Parallel()

	for index, test := range dataValidVarious {
		parsed, err := FromTaskString(test.input)
		require.NoErrorf(t, err,
			"Failed to parse input for test case %d: %s", index+1, test.title)

		expect := test.expected.DateCreated
		actual := parsed.DateCreated()

		require.Equalf(t, expect, actual,
			"Test case #%d (%s): DateCreated() should return %q for input '%s'",
			index+1, test.title, expect, test.input)
	}
}

// ----------------------------------------------------------------------------
//  Parsed.Contexts()
// ----------------------------------------------------------------------------

func TestParsed_Contexts(t *testing.T) {
	t.Parallel()

	for index, test := range dataValidVarious {
		parsed, err := FromTaskString(test.input)
		require.NoErrorf(t, err,
			"Failed to parse input for test case %d: %s", index+1, test.title)

		expect := test.expected.Contexts
		actual := parsed.Contexts()

		require.Equalf(t, expect, actual,
			"Test case #%d (%s): Contexts() should return %q for input '%s'",
			index+1, test.title, expect, test.input)
	}
}

// ----------------------------------------------------------------------------
//  Parsed.Projects()
// ----------------------------------------------------------------------------

func TestParsed_Projects(t *testing.T) {
	t.Parallel()

	for index, test := range dataValidVarious {
		parsed, err := FromTaskString(test.input)
		require.NoErrorf(t, err,
			"Failed to parse input for test case %d: %s", index+1, test.title)

		expect := test.expected.Projects
		actual := parsed.Projects()

		require.Equalf(t, expect, actual,
			"Test case #%d (%s): Projects() should return %q for input '%s'",
			index+1, test.title, expect, test.input)
	}
}

// ----------------------------------------------------------------------------
//  Parsed.Description()
// ----------------------------------------------------------------------------

func TestParsed_Description(t *testing.T) {
	t.Parallel()

	for index, test := range dataValidVarious {
		parsed, err := FromTaskString(test.input)
		require.NoErrorf(t, err,
			"Failed to parse input for test case %d: %s", index+1, test.title)

		expect := test.expected.Description
		actual := parsed.Description()

		require.Equalf(t, expect, actual,
			"Test case #%d (%s): Description() should return %q for input '%s'",
			index+1, test.title, expect, test.input)
	}
}

// ----------------------------------------------------------------------------
//  Parsed.KeyValues()
// ----------------------------------------------------------------------------

func TestParsed_KeyValues(t *testing.T) {
	t.Parallel()

	for index, test := range dataValidVarious {
		parsed, err := FromTaskString(test.input)
		require.NoErrorf(t, err,
			"Failed to parse input for test case %d: %s", index+1, test.title)

		expect := test.expected.KeyValues
		actual := parsed.KeyValues()

		require.Equalf(t, expect, actual,
			"Test case #%d (%s): KeyValues() should return %v for input '%s'",
			index+1, test.title, expect, test.input)
	}
}

// ----------------------------------------------------------------------------
//  Parsed.Comment()
/// ----------------------------------------------------------------------------

func TestParsed_Comment(t *testing.T) {
	t.Parallel()

	for index, test := range dataValidVarious {
		parsed, err := FromTaskString(test.input)
		require.NoErrorf(t, err,
			"Failed to parse input for test case %d: %s", index+1, test.title)

		expect := test.expected.Comment
		actual := parsed.Comment()

		require.Equalf(t, expect, actual,
			"Test case #%d (%s): Comment() should return %q for input '%s'",
			index+1, test.title, expect, test.input)
	}
}

// ----------------------------------------------------------------------------
//  Parsed.Components()
// ----------------------------------------------------------------------------

func TestParsed_Components(t *testing.T) {
	t.Parallel()

	for index, test := range dataValidVarious {
		parsed, err := FromTaskString(test.input)
		require.NoErrorf(t, err,
			"Failed to parse input for test case %d: %s", index+1, test.title)

		expect := test.expected
		actual := parsed.Components()

		require.Equalf(t, expect, actual,
			"Test case #%d (%s): Components() should return expected struct for input '%s'",
			index+1, test.title, test.input)
	}
}

// ----------------------------------------------------------------------------
//  Parsed.SetText()
// ----------------------------------------------------------------------------

func TestParsed_SetText(t *testing.T) {
	t.Parallel()

	t.Run("update original text", func(t *testing.T) {
		t.Parallel()

		// Create parsed task with initial text
		initialText := "(A) Initial task @context"

		p, err := FromTaskString(initialText)
		require.NoError(t, err)
		require.Equal(t, initialText, p.String(),
			"String() should return initial text")

		// Update text using SetText
		newText := "(B) Updated task @newcontext"
		p.SetText(newText)

		// Verify text was updated
		require.Equal(t, newText, p.String(),
			"String() should return updated text after SetText()")

		// Verify that parsing results are still from old text
		// (SetText doesn't re-parse, just updates the text field)
		require.Equal(t, "A", p.Priority(),
			"Priority should still be from original parse")
		require.Contains(t, p.Contexts(), "context",
			"Context should still be from original parse")
	})

	t.Run("empty text", func(t *testing.T) {
		t.Parallel()

		p, err := FromTaskString("(A) Task")
		require.NoError(t, err)

		p.SetText("")
		require.Empty(t, p.String(),
			"SetText should accept empty string")
	})
}

// ============================================================================
//  Tests for private methods
// ============================================================================

// ----------------------------------------------------------------------------
//  Parsed.offsetAfterSegments()
// ----------------------------------------------------------------------------

func TestParsed_offsetAfterSegments(t *testing.T) {
	t.Parallel()

	const (
		taskString = "x  2024-01-01 2024-01-01 one +two #three @four due:five #six"
		notFound   = -1
	)

	parsed, err := FromTaskString(taskString)
	require.NoErrorf(t, err,
		"Failed to parse input: %s", taskString)

	tests := []struct {
		name           string
		numSegments    int
		expectedOffset int
	}{
		{
			name:           "after 0 segments",
			numSegments:    0,
			expectedOffset: 0,
		},
		{
			name:           "after 1 segment",
			numSegments:    1,
			expectedOffset: len("x  "),
		},
		{
			name:           "after 2 segments",
			numSegments:    2,
			expectedOffset: len("x  2024-01-01 "),
		},
		{
			name:           "after 3 segments",
			numSegments:    3,
			expectedOffset: len("x  2024-01-01 2024-01-01 "),
		},
		{
			name:           "after all segments",
			numSegments:    len(parsed.Segments),
			expectedOffset: len(taskString),
		},
		{
			name:           "after more than total segments",
			numSegments:    len(parsed.Segments) + 5,
			expectedOffset: notFound,
		},
	}

	for index, test := range tests {
		title := fmt.Sprintf("Test #%d: %s", index+1, test.name)

		t.Run(title, func(t *testing.T) {
			t.Parallel()

			actualOffset := parsed.offsetAfterSegments(test.numSegments)
			require.Equalf(t, test.expectedOffset, actualOffset,
				"offsetAfterSegments(%d) should return %d for input '%s'",
				test.numSegments, test.expectedOffset, taskString)
		})
	}

	t.Run("segment not found in original text", func(t *testing.T) {
		t.Parallel()

		// Create a Parsed object with mismatched originalText and Segments
		parsed := Parsed{
			originalText:     "x 2024-01-01 task",
			Segments:         []segment.Segment{"x", "2024-01-01", "NONEXISTENT", "task"},
			inlineComment:    nil,
			isCommentLine:    nil,
			isDone:           nil,
			headerEndIndex:   nil,
			indexComment:     nil,
			keyValueCache:    nil,
			keyValueInit:     false,
			options:          nil,
			allowedCtrlChars: nil,
		}

		// Should stop at the unfound segment and return -1
		offset := parsed.offsetAfterSegments(3)

		require.Equalf(t, notFound, offset,
			"offsetAfterSegments should return offset up to unfound segment")
	})
}

// ----------------------------------------------------------------------------
//  Parsed.trimInlineComment()
// ----------------------------------------------------------------------------

func TestParsed_trimInlineComment(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		task     string
		argument string
		expect   string
	}{
		{
			name:     "inline comment is trimmed",
			task:     "Task description # inline comment",
			argument: "",
			expect:   "Task description ",
		},
		{
			name:     "no comment returns original text",
			task:     "Task description",
			argument: "",
			expect:   "Task description",
		},
		{
			name:     "only comment becomes empty",
			task:     "# comment",
			argument: "",
			expect:   "",
		},
		{
			name:     "empty task stays empty",
			task:     "",
			argument: "",
			expect:   "",
		},
		{
			name:     "first comment wins when multiple",
			task:     "Task # first # second",
			argument: "",
			expect:   "Task ",
		},
		{
			name:     "hash within word is preserved",
			task:     "visit site url:https://example.com/#fragment # remember to check",
			argument: "",
			expect:   "visit site url:https://example.com/#fragment ",
		},
		{
			name:     "partial text without cached fragment",
			task:     "Task description # inline comment",
			argument: "Task description ",
			expect:   "Task description ",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			tsk, err := FromTaskString(testcase.task)
			require.NoErrorf(t, err,
				"Failed to parse test input %q", testcase.task)

			arg := testcase.argument
			if arg == "" {
				arg = testcase.task
			}

			actual := tsk.trimInlineComment(arg)
			require.Equalf(t, testcase.expect, actual,
				"trimInlineComment(%q)", arg)
		})
	}
}
