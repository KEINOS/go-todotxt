package task

import (
	"testing"

	"github.com/KEINOS/go-todotxt/todo/internal/spec"
	"github.com/stretchr/testify/require"
)

// ============================================================================
//  Data Providers
// ============================================================================
//  Place test data here if it's long or re-usable. Useful to reduce lint errors
//  of 'dupl'.
//  "Common-first, edge-last" ordering is recommended.

// ----------------------------------------------------------------------------
//  New()
// ----------------------------------------------------------------------------

// Error scenarios for New() function specifically.
var dataNewErrorScenarios = []struct {
	title         string
	input         string
	errorContains string
}{
	{
		title:         "parse task with control chars (CRLF)",
		input:         "buy milk\r\n",
		errorContains: "control characters",
	},
	{
		title:         "parse task with control chars (LF)",
		input:         "buy milk\n",
		errorContains: "control characters",
	},
	{
		title:         "parse task with control chars (tab)",
		input:         "buy\tmilk",
		errorContains: "control characters",
	},
	{
		title:         "parse task with control chars (null)",
		input:         "buy\x00milk",
		errorContains: "control characters",
	},
}

// ----------------------------------------------------------------------------
//  Priority Setting Test Data
// ----------------------------------------------------------------------------

// Comprehensive test data for priority setter methods.
var dataPriority = []struct {
	title          string
	task           string
	priorityVal    string
	output         string // task.String() on success, err msg to contain on failure
	additionalTest func(t *testing.T, task *Task)
	shouldFail     bool // if true requires error on setting
}{
	// Basic/regular cases
	{
		title:          "add priority to task without priority",
		task:           "buy milk @shopA",
		priorityVal:    "A",
		output:         "(A) buy milk @shopA",
		additionalTest: nil,
		shouldFail:     false,
	},
	{
		title:          "update existing priority to different",
		task:           "(B) buy milk @shopB",
		priorityVal:    "B",
		output:         "(B) buy milk @shopB",
		additionalTest: nil,
		shouldFail:     false,
	},
	{
		title:          "update existing priority to same",
		task:           "(C) buy milk @shopC",
		priorityVal:    "C",
		output:         "(C) buy milk @shopC",
		additionalTest: nil,
		shouldFail:     false,
	},
	{
		title:       "completed task (add)",
		task:        "x buy milk @shopD",
		priorityVal: "D",
		output:      "x (D) buy milk @shopD",
		additionalTest: func(t *testing.T, tsk *Task) {
			t.Helper()

			require.True(t, tsk.IsDone(),
				"task should remain completed after setting priority")
		},
		shouldFail: false,
	},
	{
		title:       "completed task with completed date (add)",
		task:        "x 2024-01-15 buy milk @shopD",
		priorityVal: "E",
		output:      "x (E) 2024-01-15 buy milk @shopD",
		additionalTest: func(t *testing.T, tsk *Task) {
			t.Helper()

			require.Equal(t, "2024-01-15", tsk.DateCompleted(),
				"completed date should remain unchanged after setting priority")
		},
		shouldFail: false,
	},
	{
		title:       "completed task with completed and created date (add)",
		task:        "x 2024-01-15 2024-01-02 buy milk @shopD",
		priorityVal: "F",
		output:      "x (F) 2024-01-15 2024-01-02 buy milk @shopD",
		additionalTest: func(t *testing.T, tsk *Task) {
			t.Helper()

			require.Equal(t, "2024-01-15", tsk.DateCompleted(),
				"completed date should remain unchanged after setting priority")
			require.Equal(t, "2024-01-02", tsk.DateCreated(),
				"created date should remain unchanged after setting priority")
		},
		shouldFail: false,
	},
	{
		title:          "completed task with existing priority (update)",
		task:           "x (A) 2024-01-15 buy milk @shopE",
		priorityVal:    "G",
		output:         "x (G) 2024-01-15 buy milk @shopE",
		additionalTest: nil,
		shouldFail:     false,
	},
	{
		title:       "dirty flag after Apply",
		task:        "buy milk @shopF",
		priorityVal: "H",
		output:      "(H) buy milk @shopF",
		additionalTest: func(t *testing.T, tsk *Task) {
			t.Helper()

			require.False(t, tsk.isDirty,
				"isDirty should be false after Apply")
		},
		shouldFail: false,
	},
	{
		title:       "preserve structure with dates and tags",
		task:        "buy @grocery +shopping due:2024-12-25",
		priorityVal: "I",
		output:      "(I) buy @grocery +shopping due:2024-12-25",
		additionalTest: func(t *testing.T, tsk *Task) {
			t.Helper()

			require.Equal(t, []string{"grocery"}, tsk.Contexts())
			require.Equal(t, []string{"shopping"}, tsk.Projects())
			require.Equal(t, "2024-12-25", tsk.KeyValues()["due"])
		},
		shouldFail: false,
	},
	{
		title:          "task with inline comment (add)",
		task:           "buy milk # remember the milk",
		priorityVal:    "J",
		output:         "(J) buy milk # remember the milk",
		additionalTest: nil,
		shouldFail:     false,
	},
	// i18n cases (Chinese → Japanese → Korean alphabetical order)
	{
		title:          "Chinese task",
		task:           "买牛奶 @商店",
		priorityVal:    "A",
		output:         "(A) 买牛奶 @商店",
		additionalTest: nil,
		shouldFail:     false,
	},
	{
		title:          "Japanese task",
		task:           "牛乳を買う @店",
		priorityVal:    "B",
		output:         "(B) 牛乳を買う @店",
		additionalTest: nil,
		shouldFail:     false,
	},
	{
		title:          "Korean task",
		task:           "우유 사기 @가게",
		priorityVal:    "C",
		output:         "(C) 우유 사기 @가게",
		additionalTest: nil,
		shouldFail:     false,
	},
	// Failure cases
	{
		title:          "invalid priority - empty string",
		task:           "buy milk @shopG",
		priorityVal:    "",
		output:         "invalid priority",
		additionalTest: nil,
		shouldFail:     true,
	},
	{
		title:          "invalid priority - too long",
		task:           "buy milk @shopAB",
		priorityVal:    "AB",
		output:         "invalid priority",
		additionalTest: nil,
		shouldFail:     true,
	},
	{
		title:          "invalid priority - lowercase",
		task:           "buy milk @shop_g",
		priorityVal:    "g",
		output:         "invalid priority",
		additionalTest: nil,
		shouldFail:     true,
	},
	{
		title:          "invalid priority - number",
		task:           "buy milk @shop1",
		priorityVal:    "1",
		output:         "invalid priority",
		additionalTest: nil,
		shouldFail:     true,
	},
	{
		title:          "invalid priority - special char",
		task:           "buy milk @shop!",
		priorityVal:    "!",
		output:         "invalid priority",
		additionalTest: nil,
		shouldFail:     true,
	},
	// Edge/niche/irregular/dangerous cases
	{
		title:          "task with leading spaces (add)",
		task:           "   buy milk @shopH   ",
		priorityVal:    "A",
		output:         "(A)    buy milk @shopH   ",
		additionalTest: nil,
		shouldFail:     false,
	},
	{
		title:          "task with multiple contexts and projects (add)",
		task:           "meeting @office @team +project1 +project2",
		priorityVal:    "A",
		output:         "(A) meeting @office @team +project1 +project2",
		additionalTest: nil,
		shouldFail:     false,
	},
	{
		title:          "task with date in description",
		task:           "report incident of 2024-01-13",
		priorityVal:    "A",
		output:         "(A) report incident of 2024-01-13",
		additionalTest: nil,
		shouldFail:     false,
	},
	{
		title:          "priority marks in description (uncompleted)",
		task:           "(A) Check priority levels (A) and (B) for review",
		priorityVal:    "B",
		output:         "(B) Check priority levels (A) and (B) for review",
		additionalTest: nil,
		shouldFail:     false,
	},
	{
		title:          "priority marks in description (completed)",
		task:           "x (C) 2024-01-10 Check priority levels (A) and (B) for review",
		priorityVal:    "A",
		output:         "x (A) 2024-01-10 Check priority levels (A) and (B) for review",
		additionalTest: nil,
		shouldFail:     false,
	},
	{
		title:          "priority marks in description (no current priority)",
		task:           "Check priority levels (A) and (B) for review",
		priorityVal:    "C",
		output:         "(C) Check priority levels (A) and (B) for review",
		additionalTest: nil,
		shouldFail:     false,
	},
	{
		title:          "completed mark in description (uncompleted)",
		task:           "fix x coordinate bug @dev",
		priorityVal:    "A",
		output:         "(A) fix x coordinate bug @dev",
		additionalTest: nil,
		shouldFail:     false,
	},
	{
		title:          "completed mark in description (completed)",
		task:           "x 2024-01-12 fix x coordinate bug @dev",
		priorityVal:    "B",
		output:         "x (B) 2024-01-12 fix x coordinate bug @dev",
		additionalTest: nil,
		shouldFail:     false,
	},
	{
		title:          "extra spaces around priority",
		task:           "  (A)   buy milk  ",
		priorityVal:    "B",
		output:         "  (B)   buy milk  ",
		additionalTest: nil,
		shouldFail:     false,
	},
}

// ----------------------------------------------------------------------------
//  AppendSegment()
// ----------------------------------------------------------------------------

var dataAppendSegment = []struct {
	title    string
	taskStr  string
	segment  string // target segment to append
	expected string
	options  []Option
}{
	{
		title:    "append to simple task",
		taskStr:  "buy milk",
		segment:  "@home",
		expected: "buy milk @home",
		options:  nil,
	},
	{
		title:    "append to task with existing tags",
		taskStr:  "buy milk @store",
		segment:  "+shopping",
		expected: "buy milk @store +shopping",
		options:  nil,
	},
	{
		title:    "append to empty string",
		taskStr:  "",
		segment:  "@home",
		expected: "@home",
		options:  nil,
	},
	{
		title:    "append empty segment to task",
		taskStr:  "buy milk",
		segment:  "",
		expected: "buy milk",
		options:  nil,
	},
	{
		title:    "append to task with priority",
		taskStr:  "(A) buy milk",
		segment:  "due:2024-12-31",
		expected: "(A) buy milk due:2024-12-31",
		options:  nil,
	},
	{
		title:    "append to completed task",
		taskStr:  "x 2024-01-15 buy milk",
		segment:  "@store",
		expected: "x 2024-01-15 buy milk @store",
		options:  nil,
	},
	{
		title:    "append to task with inline comment",
		taskStr:  "buy milk # remember",
		segment:  "+shopping",
		expected: "buy milk +shopping # remember",
		options:  nil,
	},
	{
		title:    "append to task with no space before comment (included hash)",
		taskStr:  "buy#comment",
		segment:  "@tag",
		expected: "buy#comment @tag",
		options:  nil,
	},
	{
		title:    "prepend to comment line",
		taskStr:  "# this is a comment",
		segment:  "@tag",
		expected: "@tag # this is a comment",
		options:  nil,
	},
	{
		title:    "append before tab-delimited inline comment",
		taskStr:  "buy milk\t# remember",
		segment:  "+shopping",
		expected: "buy milk +shopping\t# remember",
		options:  []Option{WithAllowedCtrlChars([]rune{'\t'})}, // allow tab chars
	},
	{
		title:    "append to task with tab characters",
		taskStr:  "buy\tmilk\tcofee # remember",
		segment:  "+no-sugar",
		expected: "buy\tmilk\tcofee +no-sugar # remember",
		options:  []Option{WithAllowedCtrlChars([]rune{'\t'})}, // allow tab chars
	},
}

// ----------------------------------------------------------------------------
//  InsertAfter()
// ----------------------------------------------------------------------------

var dataInsertAfter = []struct {
	title     string
	taskStr   string
	targetSeg string
	insertSeg string
	expectOut string
	expectOK  bool // expected return value from InsertAfter
}{
	{
		title:     "insert after completion mark",
		taskStr:   "x buy milk",
		targetSeg: "x",
		insertSeg: "(A)",
		expectOut: "x (A) buy milk",
		expectOK:  true,
	},
	{
		title:     "insert after priority",
		taskStr:   "(A) buy milk",
		targetSeg: "(A)",
		insertSeg: "2024-01-01",
		expectOut: "(A) 2024-01-01 buy milk",
		expectOK:  true,
	},
	{
		title:     "insert after word in description",
		taskStr:   "buy milk @store",
		targetSeg: "buy",
		insertSeg: "some",
		expectOut: "buy some milk @store",
		expectOK:  true,
	},
	{
		title:     "insert after word at end",
		taskStr:   "buy milk",
		targetSeg: "milk",
		insertSeg: "@home",
		expectOut: "buy milk @home",
		expectOK:  true,
	},
	{
		title:     "insert after date",
		taskStr:   "x 2024-01-15 buy milk",
		targetSeg: "2024-01-15",
		insertSeg: "2024-01-01",
		expectOut: "x 2024-01-15 2024-01-01 buy milk",
		expectOK:  true,
	},
	// I18n cases (Chinese → Japanese → Korean alphabetical order)
	{
		title:     "insert in Chinese task",
		taskStr:   "买牛奶 @商店",
		targetSeg: "@商店",
		insertSeg: "@低脂肪",
		expectOut: "买牛奶 @商店 @低脂肪",
		expectOK:  true,
	},
	{
		title:     "insert in Japanese task",
		taskStr:   "牛乳を買う @店",
		targetSeg: "@店",
		insertSeg: "@低脂肪",
		expectOut: "牛乳を買う @店 @低脂肪",
		expectOK:  true,
	},
	{
		title:     "insert in Korean task",
		taskStr:   "우유 사기 @가게",
		targetSeg: "@가게",
		insertSeg: "@저지방",
		expectOut: "우유 사기 @가게 @저지방",
		expectOK:  true,
	},
	// Edge/niche cases
	{
		title:     "insert in empty string",
		taskStr:   "",
		targetSeg: "x",
		insertSeg: "(A)",
		expectOut: "",
		expectOK:  false,
	},
	{
		title:     "empty target - no change",
		taskStr:   "buy milk",
		targetSeg: "",
		insertSeg: "(A)",
		expectOut: "buy milk",
		expectOK:  false,
	},
	{
		title:     "target not found - no change",
		taskStr:   "buy milk",
		targetSeg: "@store",
		insertSeg: "@home",
		expectOut: "buy milk",
		expectOK:  false,
	},
	{
		title:     "insert after first occurrence only",
		taskStr:   "x buy milk x coordinate",
		targetSeg: "x",
		insertSeg: "(A)",
		expectOut: "x (A) buy milk x coordinate",
		expectOK:  true,
	},
	{
		title:     "avoid inserting into context sharing prefix",
		taskStr:   "buy milk @homework",
		targetSeg: "@home",
		insertSeg: "@urgent",
		expectOut: "buy milk @homework",
		expectOK:  false,
	},
	{
		title:     "avoid inserting into key-value sharing prefix",
		taskStr:   "watch overdue:2024-01-01 movies",
		targetSeg: "due:2024-01-01",
		insertSeg: "@followup",
		expectOut: "watch overdue:2024-01-01 movies",
		expectOK:  false,
	},
}

// ----------------------------------------------------------------------------
//  WithContext()
// ----------------------------------------------------------------------------

//nolint:dupl // Context and project data sets intentionally mirror each other for parity.
var dataWithContext = []struct {
	title     string
	taskStr   string // valid task format
	addCtx    string
	expectOut string // task.String() on success, err msg to contain on failure
	expectCtx []string
	shouldErr bool
}{
	// Golden path
	{
		title:     "add new context",
		taskStr:   "buy milk +shopping",
		addCtx:    "store",
		expectOut: "buy milk +shopping @store",
		expectCtx: []string{"store"},
		shouldErr: false,
	},
	{
		title:     "allow prefixed context",
		taskStr:   "buy milk +shopping",
		addCtx:    "@store",
		expectOut: "buy milk +shopping @store",
		expectCtx: []string{"store"},
		shouldErr: false,
	},
	// Edge cases
	{
		title:     "avoid duplicate contexts",
		taskStr:   "buy milk @home +chores",
		addCtx:    "home", // already exists
		expectOut: "buy milk @home +chores",
		expectCtx: []string{"home"},
		shouldErr: false,
	},
	{
		title:     "insert context before inline comment",
		taskStr:   "buy milk # reminder",
		addCtx:    "errand",
		expectOut: "buy milk @errand # reminder",
		expectCtx: []string{"errand"},
		shouldErr: false,
	},
	{
		title:     "ignore comment lines",
		taskStr:   "# this is a comment line",
		addCtx:    "archive",
		expectOut: "# this is a comment line",
		expectCtx: nil,
		shouldErr: false,
	},
	// Error cases
	{
		title:     "reject empty context",
		taskStr:   "buy milk",
		addCtx:    "",
		expectOut: "invalid context",
		expectCtx: nil,
		shouldErr: true,
	},
	{
		title:     "reject whitespace only",
		taskStr:   "buy milk",
		addCtx:    "   ",
		expectOut: "invalid context",
		expectCtx: nil,
		shouldErr: true,
	},
	{
		title:     "reject context containing whitespace",
		taskStr:   "buy milk",
		addCtx:    "work office",
		expectOut: "invalid context",
		expectCtx: nil,
		shouldErr: true,
	},
	{
		title:     "reject context with prefix only",
		taskStr:   "buy milk",
		addCtx:    "@",
		expectOut: "invalid context",
		expectCtx: nil,
		shouldErr: true,
	},
}

// ----------------------------------------------------------------------------
//  WithoutContext()
// ----------------------------------------------------------------------------

// Test data for WithoutContext() option while New() or Apply().
var dataWithoutContext = []struct {
	title       string
	taskStr     string
	removeCtx   string   // value passed to WithoutContext()
	expectOut   string   // task.String() on success
	errContains string   // error message to contain on failure
	expectCtx   []string // expected contexts after removal
	shouldError bool
}{
	// Basic cases
	{
		title:       "remove existing context",
		taskStr:     "buy milk @home @storeA",
		removeCtx:   "home",
		expectOut:   "buy milk @storeA",
		errContains: "",
		expectCtx:   []string{"storeA"},
		shouldError: false,
	},
	{
		title:       "no change when context missing",
		taskStr:     "buy milk @storeB",
		removeCtx:   "home",
		expectOut:   "buy milk @storeB",
		errContains: "",
		expectCtx:   []string{"storeB"},
		shouldError: false,
	},
	{
		title:       "trim whitespace in input",
		taskStr:     "buy milk @storeC",
		removeCtx:   " storeC ",
		expectOut:   "buy milk",
		errContains: "",
		expectCtx:   nil,
		shouldError: false,
	},
	{
		title:       "ignore comment lines",
		taskStr:     "# buy milk @storeD",
		removeCtx:   "storeD",
		expectOut:   "# buy milk @storeD",
		errContains: "",
		expectCtx:   nil,
		shouldError: false,
	},
	{
		title:       "avoid removing context sharing prefix",
		taskStr:     "buy milk @homework",
		removeCtx:   "home",
		expectOut:   "buy milk @homework",
		errContains: "",
		expectCtx:   []string{"homework"},
		shouldError: false,
	},
	// Error cases
	{
		title:       "reject empty value",
		taskStr:     "buy milk @storeE",
		removeCtx:   "",
		expectOut:   "",
		errContains: "invalid context",
		expectCtx:   nil,
		shouldError: true,
	},
	{
		title:       "reject whitespace only",
		taskStr:     "buy milk @storeF",
		removeCtx:   "   ",
		expectOut:   "",
		errContains: "invalid context",
		expectCtx:   nil,
		shouldError: true,
	},
	{
		title:       "reject prefix only",
		taskStr:     "buy milk @homeG",
		removeCtx:   "@",
		expectOut:   "",
		errContains: "invalid context",
		expectCtx:   nil,
		shouldError: true,
	},
	{
		title:       "reject context containing whitespace",
		taskStr:     "buy milk @storeH",
		removeCtx:   "store H",
		expectOut:   "",
		errContains: "invalid context",
		expectCtx:   nil,
		shouldError: true,
	},
}

// ----------------------------------------------------------------------------
//  WithProject()
// ----------------------------------------------------------------------------

//nolint:dupl // Context and project data sets intentionally mirror each other for parity.
var dataWithProject = []struct {
	title       string
	taskStr     string
	addProject  string
	expectOut   string // task.String() on success, err msg to contain on failure
	expectPrj   []string
	shouldError bool
}{
	// Golden path
	{
		title:       "add new project",
		taskStr:     "buy milk @store",
		addProject:  "shopping",
		expectOut:   "buy milk @store +shopping",
		expectPrj:   []string{"shopping"},
		shouldError: false,
	},
	{
		title:       "allow prefixed project",
		taskStr:     "buy milk @store",
		addProject:  "+shopping",
		expectOut:   "buy milk @store +shopping",
		expectPrj:   []string{"shopping"},
		shouldError: false,
	},
	// Edge cases
	{
		title:       "avoid duplicate projects",
		taskStr:     "buy milk @home +chores",
		addProject:  "chores", // already exists
		expectOut:   "buy milk @home +chores",
		expectPrj:   []string{"chores"},
		shouldError: false,
	},
	{
		title:       "insert project before inline comment",
		taskStr:     "buy milk # reminder",
		addProject:  "errand",
		expectOut:   "buy milk +errand # reminder",
		expectPrj:   []string{"errand"},
		shouldError: false,
	},
	{
		title:       "ignore comment lines",
		taskStr:     "# this is a comment line",
		addProject:  "archive",
		expectOut:   "# this is a comment line",
		expectPrj:   nil,
		shouldError: false,
	},
	// Error cases
	{
		title:       "reject empty project",
		taskStr:     "buy milk",
		addProject:  "",
		expectOut:   "invalid project",
		expectPrj:   nil,
		shouldError: true,
	},
	{
		title:       "reject whitespace only",
		taskStr:     "buy milk",
		addProject:  "   ",
		expectOut:   "invalid project",
		expectPrj:   nil,
		shouldError: true,
	},
	{
		title:       "reject project containing whitespace",
		taskStr:     "buy milk",
		addProject:  "work project",
		expectOut:   "invalid project",
		expectPrj:   nil,
		shouldError: true,
	},
	{
		title:       "reject project with prefix only",
		taskStr:     "buy milk",
		addProject:  "+",
		expectOut:   "invalid project",
		expectPrj:   nil,
		shouldError: true,
	},
}

// ----------------------------------------------------------------------------
//  WithoutProject()
// ----------------------------------------------------------------------------

// Test data for WithoutProject() option while New() or Apply().
var dataWithoutProject = []struct {
	title       string
	taskStr     string
	removePrj   string   // value passed to WithoutProject()
	expectOut   string   // task.String() on success, err msg to contain on failure
	expectPrj   []string // expected projects after removal
	shouldError bool
}{
	// Basic cases
	{
		title:       "remove existing project",
		taskStr:     "buy milk +shopping +storeA",
		removePrj:   "shopping",
		expectOut:   "buy milk +storeA",
		expectPrj:   []string{"storeA"},
		shouldError: false,
	},
	{
		title:       "no change when project missing",
		taskStr:     "buy milk +storeB",
		removePrj:   "home",
		expectOut:   "buy milk +storeB",
		expectPrj:   []string{"storeB"},
		shouldError: false,
	},
	{
		title:       "trim whitespace in input",
		taskStr:     "buy milk +storeC",
		removePrj:   " storeC ",
		expectOut:   "buy milk",
		expectPrj:   nil,
		shouldError: false,
	},
	{
		title:       "ignore comment lines",
		taskStr:     "# buy milk +storeD",
		removePrj:   "storeD",
		expectOut:   "# buy milk +storeD",
		expectPrj:   nil,
		shouldError: false,
	},
	{
		title:       "allow prefixed project for removal",
		taskStr:     "buy milk +shopping +storeE",
		removePrj:   "+shopping",
		expectOut:   "buy milk +storeE",
		expectPrj:   []string{"storeE"},
		shouldError: false,
	},
	{
		title:       "avoid removing project sharing prefix",
		taskStr:     "buy milk +homework",
		removePrj:   "home",
		expectOut:   "buy milk +homework",
		expectPrj:   []string{"homework"},
		shouldError: false,
	},
	// Error cases
	{
		title:       "reject empty value",
		taskStr:     "buy milk +storeF",
		removePrj:   "",
		expectOut:   "invalid project",
		expectPrj:   nil,
		shouldError: true,
	},
	{
		title:       "reject whitespace only",
		taskStr:     "buy milk +storeG",
		removePrj:   "   ",
		expectOut:   "invalid project",
		expectPrj:   nil,
		shouldError: true,
	},
	{
		title:       "reject prefix only",
		taskStr:     "buy milk +homeH",
		removePrj:   "+",
		expectOut:   "invalid project",
		expectPrj:   nil,
		shouldError: true,
	},
	{
		title:       "reject project containing whitespace",
		taskStr:     "buy milk +storeI",
		removePrj:   "store I",
		expectOut:   "invalid project",
		expectPrj:   nil,
		shouldError: true,
	},
}

// ----------------------------------------------------------------------------
//  normalizeTag()
// ----------------------------------------------------------------------------

var dataNormalizeTag = []struct {
	name      string
	value     string
	mark      spec.Mark
	expectOut string // normalized value on success, err msg on failure
	shouldErr bool
}{
	// Valid cases (golden path)
	{
		name:      "context basic",
		value:     "store",
		mark:      spec.PrefixContext,
		expectOut: "store",
		shouldErr: false,
	},
	{
		name:      "project basic",
		value:     "shopping",
		mark:      spec.PrefixProject,
		expectOut: "shopping",
		shouldErr: false,
	},
	{
		name:      "context with prefix",
		value:     "@store",
		mark:      spec.PrefixContext,
		expectOut: "store",
		shouldErr: false,
	},
	{
		name:      "project with prefix",
		value:     "+shopping",
		mark:      spec.PrefixProject,
		expectOut: "shopping",
		shouldErr: false,
	},
	// Invalid cases (error expected)
	{
		name:      "context prefix only",
		value:     "@",
		mark:      spec.PrefixContext,
		expectOut: ErrValIsEmpty.Error(),
		shouldErr: true,
	},
	{
		name:      "project prefix only",
		value:     "+",
		mark:      spec.PrefixProject,
		expectOut: ErrValIsEmpty.Error(),
		shouldErr: true,
	},
	{
		name:      "context empty",
		value:     "",
		mark:      spec.PrefixContext,
		expectOut: ErrValIsEmpty.Error(),
		shouldErr: true,
	},
	{
		name:      "project empty",
		value:     "",
		mark:      spec.PrefixProject,
		expectOut: ErrValIsEmpty.Error(),
		shouldErr: true,
	},
	{
		name:      "context whitespace within value",
		value:     "home office",
		mark:      spec.PrefixContext,
		expectOut: ErrValWithSpace.Error(),
		shouldErr: true,
	},
	{
		name:      "project whitespace within value",
		value:     "work project",
		mark:      spec.PrefixProject,
		expectOut: ErrValWithSpace.Error(),
		shouldErr: true,
	},
}

// ----------------------------------------------------------------------------
//  WithKeyValue()
// ----------------------------------------------------------------------------

// Test data for WithKeyValue() option while New() or Apply().
var dataWithKeyValue = []struct {
	title       string
	taskStr     string
	key         string
	value       string
	expectOut   string            // task.String() on success
	expectKV    map[string]string // expected key-values after operation
	errContains string            // error message to contain on failure
	shouldError bool
}{
	// Golden path - basic cases
	{
		title:       "add new key-value to simple task",
		taskStr:     "buy milk",
		key:         "due",
		value:       "2024-12-25",
		expectOut:   "buy milk due:2024-12-25",
		expectKV:    map[string]string{"due": "2024-12-25"},
		errContains: "",
		shouldError: false,
	},
	{
		title:       "add key-value to task with existing tags",
		taskStr:     "buy milk @store +shopping",
		key:         "priority",
		value:       "high",
		expectOut:   "buy milk @store +shopping priority:high",
		expectKV:    map[string]string{"priority": "high"},
		errContains: "",
		shouldError: false,
	},
	{
		title:       "update existing key-value",
		taskStr:     "buy milk due:2024-01-01",
		key:         "due",
		value:       "2024-12-25",
		expectOut:   "buy milk due:2024-12-25",
		expectKV:    map[string]string{"due": "2024-12-25"},
		errContains: "",
		shouldError: false,
	},
	{
		title:       "add multiple key-values",
		taskStr:     "buy milk due:2024-12-25",
		key:         "priority",
		value:       "high",
		expectOut:   "buy milk due:2024-12-25 priority:high",
		expectKV:    map[string]string{"due": "2024-12-25", "priority": "high"},
		errContains: "",
		shouldError: false,
	},
	// Edge cases - special formats
	{
		title:       "key-value with URL-like value (contains colons)",
		taskStr:     "check website",
		key:         "url",
		value:       "https://example.com",
		expectOut:   "check website url:https://example.com",
		expectKV:    map[string]string{"url": "https://example.com"},
		errContains: "",
		shouldError: false,
	},
	{
		title:       "key-value with numeric value",
		taskStr:     "exercise routine",
		key:         "reps",
		value:       "30",
		expectOut:   "exercise routine reps:30",
		expectKV:    map[string]string{"reps": "30"},
		errContains: "",
		shouldError: false,
	},
	{
		title:       "key-value with time format",
		taskStr:     "meeting",
		key:         "time",
		value:       "14:30",
		expectOut:   "meeting time:14:30",
		expectKV:    map[string]string{"time": "14:30"},
		errContains: "",
		shouldError: false,
	},
	// Edge cases - task structure
	{
		title:       "add to completed task",
		taskStr:     "x 2024-01-15 buy milk",
		key:         "due",
		value:       "2024-01-10",
		expectOut:   "x 2024-01-15 buy milk due:2024-01-10",
		expectKV:    map[string]string{"due": "2024-01-10"},
		errContains: "",
		shouldError: false,
	},
	{
		title:       "add to task with priority",
		taskStr:     "(A) buy milk",
		key:         "due",
		value:       "2024-12-25",
		expectOut:   "(A) buy milk due:2024-12-25",
		expectKV:    map[string]string{"due": "2024-12-25"},
		errContains: "",
		shouldError: false,
	},
	{
		title:       "preserve inline comment position",
		taskStr:     "buy milk # shopping list",
		key:         "due",
		value:       "2024-12-25",
		expectOut:   "buy milk due:2024-12-25 # shopping list",
		expectKV:    map[string]string{"due": "2024-12-25"},
		errContains: "",
		shouldError: false,
	},
	{
		title:       "update key-value preserving other segments",
		taskStr:     "buy milk @store due:2024-01-01 +shopping",
		key:         "due",
		value:       "2024-12-25",
		expectOut:   "buy milk @store due:2024-12-25 +shopping",
		expectKV:    map[string]string{"due": "2024-12-25"},
		errContains: "",
		shouldError: false,
	},
	// Edge cases - whitespace handling
	{
		title:       "trim key whitespace",
		taskStr:     "buy milk",
		key:         "  due  ",
		value:       "2024-12-25",
		expectOut:   "buy milk due:2024-12-25",
		expectKV:    map[string]string{"due": "2024-12-25"},
		errContains: "",
		shouldError: false,
	},
	{
		title:       "trim value whitespace",
		taskStr:     "buy milk",
		key:         "due",
		value:       "  2024-12-25  ",
		expectOut:   "buy milk due:2024-12-25",
		expectKV:    map[string]string{"due": "2024-12-25"},
		errContains: "",
		shouldError: false,
	},
	// Comment line handling
	{
		title:       "ignore on comment line",
		taskStr:     "# this is a comment",
		key:         "due",
		value:       "2024-12-25",
		expectOut:   "# this is a comment",
		expectKV:    nil,
		errContains: "",
		shouldError: false,
	},
	// I18n cases
	{
		title:       "Chinese key-value",
		taskStr:     "买牛奶",
		key:         "截止日期",
		value:       "2024-12-25",
		expectOut:   "买牛奶 截止日期:2024-12-25",
		expectKV:    map[string]string{"截止日期": "2024-12-25"},
		errContains: "",
		shouldError: false,
	},
	{
		title:       "Japanese key-value",
		taskStr:     "牛乳を買う",
		key:         "期限",
		value:       "2024-12-25",
		expectOut:   "牛乳を買う 期限:2024-12-25",
		expectKV:    map[string]string{"期限": "2024-12-25"},
		errContains: "",
		shouldError: false,
	},
	{
		title:       "Korean key-value",
		taskStr:     "우유 사기",
		key:         "마감일",
		value:       "2024-12-25",
		expectOut:   "우유 사기 마감일:2024-12-25",
		expectKV:    map[string]string{"마감일": "2024-12-25"},
		errContains: "",
		shouldError: false,
	},
	// Error cases - key validation
	{
		title:       "reject empty key",
		taskStr:     "buy milk",
		key:         "",
		value:       "2024-12-25",
		expectOut:   "",
		expectKV:    nil,
		errContains: "key cannot be empty",
		shouldError: true,
	},
	{
		title:       "reject whitespace only key",
		taskStr:     "buy milk",
		key:         "   ",
		value:       "2024-12-25",
		expectOut:   "",
		expectKV:    nil,
		errContains: "key cannot be empty",
		shouldError: true,
	},
	{
		title:       "reject key with whitespace",
		taskStr:     "buy milk",
		key:         "due date",
		value:       "2024-12-25",
		expectOut:   "",
		expectKV:    nil,
		errContains: "key cannot contain whitespace",
		shouldError: true,
	},
	{
		title:       "reject key with tab",
		taskStr:     "buy milk",
		key:         "due\tdate",
		value:       "2024-12-25",
		expectOut:   "",
		expectKV:    nil,
		errContains: "key cannot contain whitespace",
		shouldError: true,
	},
	// Error cases - value validation
	{
		title:       "reject empty value",
		taskStr:     "buy milk",
		key:         "due",
		value:       "",
		expectOut:   "",
		expectKV:    nil,
		errContains: "value cannot be empty",
		shouldError: true,
	},
	{
		title:       "reject whitespace only value",
		taskStr:     "buy milk",
		key:         "due",
		value:       "   ",
		expectOut:   "",
		expectKV:    nil,
		errContains: "value cannot be empty",
		shouldError: true,
	},
	{
		title:       "reject value with whitespace",
		taskStr:     "buy milk",
		key:         "due",
		value:       "2024 12 25",
		expectOut:   "",
		expectKV:    nil,
		errContains: "value cannot contain whitespace",
		shouldError: true,
	},
	{
		title:       "reject value with newline",
		taskStr:     "buy milk",
		key:         "note",
		value:       "line1\nline2",
		expectOut:   "",
		expectKV:    nil,
		errContains: "value cannot contain whitespace",
		shouldError: true,
	},
}

// ----------------------------------------------------------------------------
//  WithoutKeyValue()
// ----------------------------------------------------------------------------

// Test data for WithoutKeyValue() option while New() or Apply().
var dataWithoutKeyValue = []struct {
	title       string
	taskStr     string
	key         string
	expectOut   string            // task.String() on success
	expectKV    map[string]string // expected key-values after operation
	errContains string            // error message to contain on failure
	shouldError bool
}{
	// Golden path - basic cases
	{
		title:       "remove existing key-value",
		taskStr:     "buy milk due:2024-12-25",
		key:         "due",
		expectOut:   "buy milk",
		expectKV:    nil,
		errContains: "",
		shouldError: false,
	},
	{
		title:       "remove one of multiple key-values",
		taskStr:     "buy milk due:2024-12-25 priority:high",
		key:         "due",
		expectOut:   "buy milk priority:high",
		expectKV:    map[string]string{"priority": "high"},
		errContains: "",
		shouldError: false,
	},
	{
		title:       "remove key-value preserving other segments",
		taskStr:     "buy milk @store due:2024-12-25 +shopping",
		key:         "due",
		expectOut:   "buy milk @store +shopping",
		expectKV:    nil,
		errContains: "",
		shouldError: false,
	},
	// Edge cases - non-existent key
	{
		title:       "no change when key not found",
		taskStr:     "buy milk due:2024-12-25",
		key:         "priority",
		expectOut:   "buy milk due:2024-12-25",
		expectKV:    map[string]string{"due": "2024-12-25"},
		errContains: "",
		shouldError: false,
	},
	{
		title:       "no change when task has no key-values",
		taskStr:     "buy milk @store +shopping",
		key:         "due",
		expectOut:   "buy milk @store +shopping",
		expectKV:    nil,
		errContains: "",
		shouldError: false,
	},
	// Edge cases - special formats
	{
		title:       "remove key-value with URL-like value",
		taskStr:     "check website url:https://example.com",
		key:         "url",
		expectOut:   "check website",
		expectKV:    nil,
		errContains: "",
		shouldError: false,
	},
	{
		title:       "remove key-value with colons in value",
		taskStr:     "meeting time:14:30:00",
		key:         "time",
		expectOut:   "meeting",
		expectKV:    nil,
		errContains: "",
		shouldError: false,
	},
	// Edge cases - task structure
	{
		title:       "remove from completed task",
		taskStr:     "x 2024-01-15 buy milk due:2024-01-10",
		key:         "due",
		expectOut:   "x 2024-01-15 buy milk",
		expectKV:    nil,
		errContains: "",
		shouldError: false,
	},
	{
		title:       "remove from task with priority",
		taskStr:     "(A) buy milk due:2024-12-25",
		key:         "due",
		expectOut:   "(A) buy milk",
		expectKV:    nil,
		errContains: "",
		shouldError: false,
	},
	{
		title:       "remove preserving inline comment",
		taskStr:     "buy milk due:2024-12-25 # shopping list",
		key:         "due",
		expectOut:   "buy milk # shopping list",
		expectKV:    nil,
		errContains: "",
		shouldError: false,
	},
	// Edge cases - whitespace handling
	{
		title:       "trim key whitespace",
		taskStr:     "buy milk due:2024-12-25",
		key:         "  due  ",
		expectOut:   "buy milk",
		expectKV:    nil,
		errContains: "",
		shouldError: false,
	},
	// Edge cases - similar key names
	{
		title:       "avoid removing key sharing prefix",
		taskStr:     "watch overdue:2024-01-01 movies",
		key:         "due",
		expectOut:   "watch overdue:2024-01-01 movies",
		expectKV:    map[string]string{"overdue": "2024-01-01"},
		errContains: "",
		shouldError: false,
	},
	{
		title:       "remove exact key match only",
		taskStr:     "task due:2024-01-01 dueDate:2024-02-01",
		key:         "due",
		expectOut:   "task dueDate:2024-02-01",
		expectKV:    map[string]string{"dueDate": "2024-02-01"},
		errContains: "",
		shouldError: false,
	},
	// Comment line handling
	{
		title:       "ignore on comment line",
		taskStr:     "# comment due:2024-12-25",
		key:         "due",
		expectOut:   "# comment due:2024-12-25",
		expectKV:    nil,
		errContains: "",
		shouldError: false,
	},
	// I18n cases
	{
		title:       "remove Chinese key-value",
		taskStr:     "买牛奶 截止日期:2024-12-25",
		key:         "截止日期",
		expectOut:   "买牛奶",
		expectKV:    nil,
		errContains: "",
		shouldError: false,
	},
	{
		title:       "remove Japanese key-value",
		taskStr:     "牛乳を買う 期限:2024-12-25",
		key:         "期限",
		expectOut:   "牛乳を買う",
		expectKV:    nil,
		errContains: "",
		shouldError: false,
	},
	{
		title:       "remove Korean key-value",
		taskStr:     "우유 사기 마감일:2024-12-25",
		key:         "마감일",
		expectOut:   "우유 사기",
		expectKV:    nil,
		errContains: "",
		shouldError: false,
	},
	// Error cases - key validation
	{
		title:       "reject empty key",
		taskStr:     "buy milk due:2024-12-25",
		key:         "",
		expectOut:   "",
		expectKV:    nil,
		errContains: "key cannot be empty",
		shouldError: true,
	},
	{
		title:       "reject whitespace only key",
		taskStr:     "buy milk due:2024-12-25",
		key:         "   ",
		expectOut:   "",
		expectKV:    nil,
		errContains: "key cannot be empty",
		shouldError: true,
	},
	{
		title:       "reject key with whitespace",
		taskStr:     "buy milk due:2024-12-25",
		key:         "due date",
		expectOut:   "",
		expectKV:    nil,
		errContains: "key cannot contain whitespace",
		shouldError: true,
	},
	{
		title:       "reject key with tab",
		taskStr:     "buy milk due:2024-12-25",
		key:         "due\tdate",
		expectOut:   "",
		expectKV:    nil,
		errContains: "key cannot contain whitespace",
		shouldError: true,
	},
}
