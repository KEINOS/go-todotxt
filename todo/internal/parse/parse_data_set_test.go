package parse

import "github.com/KEINOS/go-todotxt/todo/internal/segment"

// ============================================================================
//  Data Providers
// ============================================================================
//  Place here if test cases require large data sets or repeatable use of data.
//  "Common-first, edge-last" ordering is recommended.

// dataValidVarious provides various data for valid task patterns.
//
//nolint:lll // overlong lines but acceptable as test data
var dataValidVarious = []struct {
	title    string
	input    string
	expected Components
}{
	// Full and empty components
	{
		title: "Completed task with full components",
		input: "x (A) 2016-05-20 2016-04-30 measure space for +chapelShelving @chapel due:2016-05-30 location:mainOffice # got gratuity from chaplain",
		expected: Components{
			HasInlineComment: true,
			PosInlineComment: 105,
			IsCommentLine:    false,
			IsDone:           true,
			Priority:         "A",
			DateCompleted:    "2016-05-20",
			DateCreated:      "2016-04-30",
			Contexts:         []string{"chapel"},
			Projects:         []string{"chapelShelving"},
			Description:      "measure space for +chapelShelving @chapel due:2016-05-30 location:mainOffice",
			KeyValues:        map[string]string{"due": "2016-05-30", "location": "mainOffice"},
			Comment:          "# got gratuity from chaplain",
		},
	},
	{
		title: "Empty input",
		input: "",
		expected: Components{
			HasInlineComment: false,
			PosInlineComment: -1,
			IsCommentLine:    false,
			IsDone:           false,
			Priority:         "",
			DateCompleted:    "",
			DateCreated:      "",
			Contexts:         nil,
			Projects:         nil,
			Description:      "",
			KeyValues:        nil,
			Comment:          "",
		},
	},
	// Individual component tests
	{
		title: "Task with priority and context tag (uncompleted)",
		input: "(A) Thank Mom for the meatballs @phone",
		expected: Components{
			HasInlineComment: false,
			PosInlineComment: -1,
			IsCommentLine:    false,
			IsDone:           false,
			Priority:         "A",
			DateCompleted:    "",
			DateCreated:      "",
			Contexts:         []string{"phone"},
			Projects:         nil,
			Description:      "Thank Mom for the meatballs @phone",
			KeyValues:        nil,
			Comment:          "",
		},
	},
	{
		title: "Task with created date (uncompleted)",
		input: "(A) 2016-05-20 Thank Mom for the meatballs @phone",
		expected: Components{
			HasInlineComment: false,
			PosInlineComment: -1,
			IsCommentLine:    false,
			IsDone:           false,
			Priority:         "A",
			DateCompleted:    "",
			DateCreated:      "2016-05-20",
			Contexts:         []string{"phone"},
			Projects:         nil,
			Description:      "Thank Mom for the meatballs @phone",
			KeyValues:        nil,
			Comment:          "",
		},
	},
	{
		title: "Task with project tag only (uncompleted)",
		input: "Post signs around the neighborhood +GarageSale",
		expected: Components{
			HasInlineComment: false,
			PosInlineComment: -1,
			IsCommentLine:    false,
			IsDone:           false,
			Priority:         "",
			DateCompleted:    "",
			DateCreated:      "",
			Contexts:         nil,
			Projects:         []string{"GarageSale"},
			Description:      "Post signs around the neighborhood +GarageSale",
			KeyValues:        nil,
			Comment:          "",
		},
	},
	{
		title: "Task with date in description",
		input: "x (A) 2025-10-15 report incident of 2025-10-13 location:AWS due:2025-10-15",
		expected: Components{
			HasInlineComment: false,
			PosInlineComment: -1,
			IsCommentLine:    false,
			IsDone:           true,
			Priority:         "A",
			DateCompleted:    "2025-10-15",
			DateCreated:      "",
			Contexts:         nil,
			Projects:         nil,
			Description:      "report incident of 2025-10-13 location:AWS due:2025-10-15",
			KeyValues:        map[string]string{"due": "2025-10-15", "location": "AWS"},
			Comment:          "",
		},
	},
	{
		title: "Task with leading x in word",
		input: "xylophone practice session @home +music",
		expected: Components{
			HasInlineComment: false,
			PosInlineComment: -1,
			IsCommentLine:    false,
			IsDone:           false,
			Priority:         "",
			DateCompleted:    "",
			DateCreated:      "",
			Contexts:         []string{"home"},
			Projects:         []string{"music"},
			Description:      "xylophone practice session @home +music",
			KeyValues:        nil,
			Comment:          "",
		},
	},
	{
		title: "Completed task with completion date only",
		input: "x 2016-05-20 measure space",
		expected: Components{
			HasInlineComment: false,
			PosInlineComment: -1,
			IsCommentLine:    false,
			IsDone:           true,
			Priority:         "",
			DateCompleted:    "2016-05-20",
			DateCreated:      "",
			Contexts:         nil,
			Projects:         nil,
			Description:      "measure space",
			KeyValues:        nil,
			Comment:          "",
		},
	},
	{
		title: "Contains date but not as completion or created date (no priority found)",
		input: "x incident 2016-05-20 reporting needs review location:AWS due:2016-05-25",
		expected: Components{
			HasInlineComment: false,
			PosInlineComment: -1,
			IsCommentLine:    false,
			IsDone:           true,
			Priority:         "",
			DateCompleted:    "",
			DateCreated:      "",
			Contexts:         nil,
			Projects:         nil,
			Description:      "incident 2016-05-20 reporting needs review location:AWS due:2016-05-25",
			KeyValues:        map[string]string{"due": "2016-05-25", "location": "AWS"},
			Comment:          "",
		},
	},
	// Test #10
	{
		title: "Comment line (Commented out completed task with full components)",
		input: "#x (A) 2016-05-20 2016-04-30 measure space for " +
			"+chapelShelving @chapel due:2016-05-30 location:mainOffice",
		expected: Components{
			HasInlineComment: true,
			PosInlineComment: 0,
			IsCommentLine:    true,
			IsDone:           false,
			Priority:         "",
			DateCompleted:    "",
			DateCreated:      "",
			Contexts:         nil,
			Projects:         nil,
			Description:      "",
			KeyValues:        nil,
			Comment: "#x (A) 2016-05-20 2016-04-30 measure space for " +
				"+chapelShelving @chapel due:2016-05-30 location:mainOffice",
		},
	},
	{
		title: "Comment line (only '#')",
		input: "#",
		expected: Components{
			HasInlineComment: true,
			PosInlineComment: 0,
			IsCommentLine:    true,
			IsDone:           false,
			Priority:         "",
			DateCompleted:    "",
			DateCreated:      "",
			Contexts:         nil,
			Projects:         nil,
			Description:      "",
			KeyValues:        nil,
			Comment:          "#",
		},
	},
	{
		title: "Simple task with inline comment (no priority, no dates)",
		input: "Buy milk @store # don't forget",
		expected: Components{
			HasInlineComment: true,
			PosInlineComment: 16,
			IsCommentLine:    false,
			IsDone:           false,
			Priority:         "",
			DateCompleted:    "",
			DateCreated:      "",
			Contexts:         []string{"store"},
			Projects:         nil,
			Description:      "Buy milk @store",
			KeyValues:        nil,
			Comment:          "# don't forget",
		},
	},
	{
		title: "Completed task with extra spaces and inline comment",
		input: "x    2016-05-20 2016-05-19  measure space    # comment",
		expected: Components{
			HasInlineComment: true,
			PosInlineComment: 45,
			IsCommentLine:    false,
			IsDone:           true,
			Priority:         "",
			DateCompleted:    "2016-05-20",
			DateCreated:      "2016-05-19",
			Contexts:         nil,
			Projects:         nil,
			Description:      "measure space",
			KeyValues:        nil,
			Comment:          "# comment",
		},
	},
	{
		title: "Task with consecutive spaces before context tag",
		input: "task    @context",
		expected: Components{
			HasInlineComment: false,
			PosInlineComment: -1,
			IsCommentLine:    false,
			IsDone:           false,
			Priority:         "",
			DateCompleted:    "",
			DateCreated:      "",
			Contexts:         []string{"context"},
			Projects:         nil,
			Description:      "task    @context",
			KeyValues:        nil,
			Comment:          "",
		},
	},
	{
		title: "Task with multiple consecutive spaces between segments",
		input: "(A)   2016-05-20   buy    milk    @store   +shopping",
		expected: Components{
			HasInlineComment: false,
			PosInlineComment: -1,
			IsCommentLine:    false,
			IsDone:           false,
			Priority:         "A",
			DateCompleted:    "",
			DateCreated:      "2016-05-20",
			Contexts:         []string{"store"},
			Projects:         []string{"shopping"},
			Description:      "buy    milk    @store   +shopping",
			KeyValues:        nil,
			Comment:          "",
		},
	},
	{
		title: "Completed task with same-day completion and created dates",
		input: "x 2016-05-20 2016-05-20 submit report +work @office # done quickly",
		expected: Components{
			HasInlineComment: true,
			PosInlineComment: 52,
			IsCommentLine:    false,
			IsDone:           true,
			Priority:         "",
			DateCompleted:    "2016-05-20",
			DateCreated:      "2016-05-20",
			Contexts:         []string{"office"},
			Projects:         []string{"work"},
			Description:      "submit report +work @office",
			KeyValues:        nil,
			Comment:          "# done quickly",
		},
	},
	{
		title: "Issue #21: any tag can be placed anywhere in the description",
		input: "2025-01-02 testing @inline +project definition for a task @testing time:0s",
		expected: Components{
			HasInlineComment: false,
			PosInlineComment: -1,
			IsCommentLine:    false,
			IsDone:           false,
			Priority:         "",
			DateCompleted:    "",
			DateCreated:      "2025-01-02",
			Contexts:         []string{"inline", "testing"},
			Projects:         []string{"project"},
			Description:      "testing @inline +project definition for a task @testing time:0s",
			KeyValues:        map[string]string{"time": "0s"},
			Comment:          "",
		},
	},
	// i18n test cases (Chinese → Japanese → Korean alphabetical order)
	{
		title: "Task with Chinese description",
		input: "(A) 中文任务 +项目 @团队",
		expected: Components{
			HasInlineComment: false,
			PosInlineComment: -1,
			IsCommentLine:    false,
			IsDone:           false,
			Priority:         "A",
			DateCompleted:    "",
			DateCreated:      "",
			Contexts:         []string{"团队"},
			Projects:         []string{"项目"},
			Description:      "中文任务 +项目 @团队",
			KeyValues:        nil,
			Comment:          "",
		},
	},
	{
		title: "Task with Japanese description",
		input: "(B) 日本語のタスク +プロジェクト @office",
		expected: Components{
			HasInlineComment: false,
			PosInlineComment: -1,
			IsCommentLine:    false,
			IsDone:           false,
			Priority:         "B",
			DateCompleted:    "",
			DateCreated:      "",
			Contexts:         []string{"office"},
			Projects:         []string{"プロジェクト"},
			Description:      "日本語のタスク +プロジェクト @office",
			KeyValues:        nil,
			Comment:          "",
		},
	},
	{
		title: "Task with Korean description",
		input: "(C) 한국어 작업 +프로젝트 @팀",
		expected: Components{
			HasInlineComment: false,
			PosInlineComment: -1,
			IsCommentLine:    false,
			IsDone:           false,
			Priority:         "C",
			DateCompleted:    "",
			DateCreated:      "",
			Contexts:         []string{"팀"},
			Projects:         []string{"프로젝트"},
			Description:      "한국어 작업 +프로젝트 @팀",
			KeyValues:        nil,
			Comment:          "",
		},
	},
	{
		title: "Task with mixed language description and inline comment",
		input: "x 2016-05-20 2016-04-30 これは @日本語 @한국어 @中文 +プロジェクト +프로젝트 +项目 #超重要",
		expected: Components{
			HasInlineComment: true,
			PosInlineComment: 106,
			IsCommentLine:    false,
			IsDone:           true,
			Priority:         "",
			DateCompleted:    "2016-05-20",
			DateCreated:      "2016-04-30",
			Contexts:         []string{"日本語", "한국어", "中文"},
			Projects:         []string{"プロジェクト", "프로젝트", "项目"},
			Description:      "これは @日本語 @한국어 @中文 +プロジェクト +프로젝트 +项目",
			KeyValues:        nil,
			Comment:          "#超重要",
		},
	},
	// niche/edge cases
	{
		title: "Two date segments but not completed (treat first matched date as created)",
		input: "2016-05-20 2016-04-30 measure space for +chapelShelving @chapel",
		expected: Components{
			HasInlineComment: false,
			PosInlineComment: -1,
			IsCommentLine:    false,
			IsDone:           false,
			Priority:         "",
			DateCompleted:    "",
			DateCreated:      "2016-05-20",
			Contexts:         []string{"chapel"},
			Projects:         []string{"chapelShelving"},
			Description:      "measure space for +chapelShelving @chapel",
			KeyValues:        nil,
			Comment:          "",
		},
	},
	{
		title: "Repeated context and project tags",
		input: "Prepare for +event @home +event @home +anotherProject @anotherContext",
		expected: Components{
			HasInlineComment: false,
			PosInlineComment: -1,
			IsCommentLine:    false,
			IsDone:           false,
			Priority:         "",
			DateCompleted:    "",
			DateCreated:      "",
			Contexts:         []string{"home", "anotherContext"},
			Projects:         []string{"event", "anotherProject"},
			Description:      "Prepare for +event @home +event @home +anotherProject @anotherContext",
			KeyValues:        nil,
			Comment:          "",
		},
	},
	{
		title: "Comment line ('#' with leading spaces)",
		input: "    # Comment line for readability",
		expected: Components{
			HasInlineComment: true,
			PosInlineComment: 4,
			IsCommentLine:    true,
			IsDone:           false,
			Priority:         "",
			DateCompleted:    "",
			DateCreated:      "",
			Contexts:         nil,
			Projects:         nil,
			Description:      "",
			KeyValues:        nil,
			Comment:          "# Comment line for readability",
		},
	},
}

// dataInlineCommentPos provides scenarios covering inline comment detection positions.
var dataInlineCommentPos = []struct {
	title        string
	input        string
	allowedRunes []rune
	expectedPos  int
}{
	{
		title:        "Task without comment",
		input:        "Deploy release candidate",
		allowedRunes: nil,
		expectedPos:  -1,
	},
	{
		title:        "Comment line at start",
		input:        "# System maintenance window",
		allowedRunes: nil,
		expectedPos:  0,
	},
	{
		title:        "Comment line with leading spaces",
		input:        "   # Comment line for readability",
		allowedRunes: nil,
		expectedPos:  3,
	},
	{
		title:        "Inline comment separated by space",
		input:        "Review logs # keep last 30 days",
		allowedRunes: nil,
		expectedPos:  len("Review logs "),
	},
	{
		title:        "Inline comment separated by tab",
		input:        "Backup\t# nightly",
		allowedRunes: []rune{'\t'},
		expectedPos:  len("Backup\t"),
	},
	{
		title:        "Inline comment separated by multiple tabs when allowed",
		input:        "Backup\t\t\t# nightly",
		allowedRunes: []rune{'\t'},
		expectedPos:  len("Backup\t\t\t"),
	},
	{
		title:        "Inline comment separated by mixed spaces and tabs when allowed",
		input:        "Backup \t # nightly",
		allowedRunes: []rune{'\t'},
		expectedPos:  len("Backup \t "),
	},
	{
		title:        "Hash within word should not count",
		input:        "Rehash#plan",
		allowedRunes: nil,
		expectedPos:  -1,
	},
	{
		title:        "Hash in URL should not count",
		input:        "visit site url:https://example.com/#fragment # remember to check",
		allowedRunes: nil,
		expectedPos:  len("visit site url:https://example.com/#fragment "),
	},
}

// dataPriority provides various data for priority tests.
var dataPriority = []struct {
	title    string
	input    string
	expected string
}{
	{
		title:    "Task with priority A",
		input:    "(A) Thank Mom for the meatballs @phone",
		expected: "A",
	},
	{
		title:    "Task with priority Z",
		input:    "(Z) Schedule meeting with team @office",
		expected: "Z",
	},
	{
		title:    "Task with no priority",
		input:    "Call Mom +Family",
		expected: "",
	},
	{
		title:    "Task with wrong priority (ZZ)",
		input:    "(ZZ) Call Billy Gibbons +RockNRoll",
		expected: "",
	},
	{
		title:    "Task with wrong priority '（A）' (full-width parentheses)",
		input:    "（A） Call Billy Gibbons +RockNRoll",
		expected: "",
	},
	{
		title:    "Task with wrong priority '(Ａ)' (full-width alphabet)",
		input:    "(Ａ) Call Billy Gibbons +RockNRoll",
		expected: "",
	},
	{
		title:    "Completed task with priority in wrong position (invalid)",
		input:    "x 2016-05-20 2016-04-30 (B) measure space for +chapelShelving @chapel",
		expected: "",
	},
	{
		title:    "Task with lowercase priority (invalid)",
		input:    "(a) Work on project",
		expected: "",
	},
	{
		title:    "Task with priority at end (invalid)",
		input:    "Work on project (A)",
		expected: "",
	},
	{
		title:    "Empty input",
		input:    "",
		expected: "",
	},
	{
		title:    "Only whitespace",
		input:    "   ",
		expected: "",
	},
}

// Data set for control character tests.
var dataControlChars = []struct {
	input           string
	allowedChars    []rune
	expectSegments  []segment.Segment
	containsCtlChar bool
}{
	// Golden/regular cases
	{
		input:           "",
		allowedChars:    nil,
		containsCtlChar: false,
		expectSegments:  []segment.Segment{},
	},
	{
		input:           "normal text",
		allowedChars:    nil,
		containsCtlChar: false,
		expectSegments:  []segment.Segment{segment.Segment("normal"), segment.Segment("text")},
	},
	{
		input:           "text with space",
		allowedChars:    nil,
		containsCtlChar: false,
		expectSegments:  []segment.Segment{segment.Segment("text"), segment.Segment("with"), segment.Segment("space")},
	},
	// allowed cases
	{
		input:           "text with \ttabs\t+project\t@context\t#comment",
		allowedChars:    []rune{'\t'},
		containsCtlChar: false,
		expectSegments: []segment.Segment{
			segment.Segment("text"),
			segment.Segment("with"),
			segment.Segment("tabs"),
			segment.Segment("+project"),
			segment.Segment("@context"),
			segment.Segment("#comment"),
		},
	},
	// Cases with control characters (not allowed)
	{
		input:           "text with \x00 null",
		allowedChars:    nil,
		containsCtlChar: true,
		expectSegments:  nil,
	},
	{
		input:           "text with \t tab",
		allowedChars:    nil,
		containsCtlChar: true,
		expectSegments:  nil,
	},
	{
		input:           "text with \n newline",
		allowedChars:    nil,
		containsCtlChar: true,
		expectSegments:  nil,
	},
	{
		input:           "text with \r return",
		allowedChars:    nil,
		containsCtlChar: true,
		expectSegments:  nil,
	},
	{
		input:           "text with \x1b escape",
		allowedChars:    nil,
		containsCtlChar: true,
		expectSegments:  nil,
	},
	{
		input:           "text with \uFEFF BOM",
		allowedChars:    nil,
		containsCtlChar: true,
		expectSegments:  nil,
	},
	{
		input:           "text with \u202E RTL override",
		allowedChars:    nil,
		containsCtlChar: true,
		expectSegments:  nil,
	},
	// edge/niche cases
	// allowed control characters only
	{
		input:           "\t\t\t",
		allowedChars:    []rune{'\t'},
		containsCtlChar: false,
		expectSegments:  []segment.Segment{},
	},
	// mixture of allowed and disallowed control characters
	{
		input:           "\ttext with tab and \x00 newline",
		allowedChars:    []rune{'\t'},
		containsCtlChar: true,
		expectSegments:  nil,
	},
}
