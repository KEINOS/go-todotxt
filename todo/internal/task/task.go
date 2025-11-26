/*
Package task provides methods to edit todo.txt strings while preserving their
original formatting.

It uses the "parse" package for reading task data and provides methods to
modify priorities, tags, dates, and comments. Changes are applied only when
Apply() is called.

For read-only operations, and no modifications are needed, use the "parse"
package instead.
*/
package task

import (
	"strings"
	"time"
	"unicode"

	"github.com/KEINOS/go-todotxt/todo/internal/parse"
	"github.com/KEINOS/go-todotxt/todo/internal/spec"
)

// Task holds the parsed data of a todo.txt task.
// It embeds the parse.Parsed struct to delegate read-only operations.
//
// The design of this package is to treat the original task string as the
// source of truth. Modifications are made directly to the string to preserve
// user-defined formatting, such as inline tags.
//
// Changes are deferred until Apply() is called. This method re-parses the
// modified string to update the task's structured data.
type Task struct {
	*parse.Parsed // Delegates read-only methods to Parsed.

	parseOpts []parse.Option // Options for re-parsing via parse.Parsed.Parse()
	isDirty   bool           // True if the task string has been modified.
}

// TimeNow is a variable holding time.Now for testing purposes.
//
// Replace this function in tests (monkey patch) to control time-dependent behavior:
//
//	task.TimeNow = func() time.Time { return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) }
var TimeNow = time.Now

// ============================================================================
//  Constructor
// ============================================================================

// New creates a new Task from a task string.
// It accepts "With*" options to customize parsing and manipulation.
//
// The creation process is:
//  1. Extracts parse-related options.
//  2. Parses the initial task string.
//  3. Applies manipulation options to the raw task string.
//  4. Calls Apply() to re-parse the modified string and finalize the state.
func New(taskString string, opts ...Option) (*Task, error) {
	newTask := new(Task)

	// 1. Extract parse-related options.
	for _, opt := range opts {
		err := opt(newTask)
		if err != nil {
			return nil, wrapError(err, "failed to extract parse option")
		}
	}

	// 2. Parse the initial task string.
	parsed, err := parse.FromTaskString(taskString, newTask.parseOpts...)
	if err != nil {
		return nil, wrapError(err, "failed to create new task")
	}

	newTask.Parsed = parsed

	// 3. Apply manipulation options.
	for _, opt := range opts {
		err := opt(newTask)
		if err != nil {
			return nil, wrapError(err, "failed to apply option")
		}
	}

	// 4. Finalize the task by applying changes.
	err = newTask.Apply()
	if err != nil {
		return nil, wrapError(err, "failed to apply initial options")
	}

	return newTask, nil
}

// ============================================================================
//  Public Methods
// ============================================================================

// ----------------------------------------------------------------------------
//  Read (get components from the task)
// ----------------------------------------------------------------------------
// Most read methods are inherited from the embedded parse.Parsed struct.
// These methods are read-only and do not modify the task.
// See the example in task_example_test.go for more details.

// String returns the full task string.
func (t *Task) String() string {
	return t.Parsed.String()
}

// IsDirty returns true if the task has been modified since the last Apply call.
func (t *Task) IsDirty() bool {
	return t.isDirty
}

// ----------------------------------------------------------------------------
//  Create (add new components to the task)
// ----------------------------------------------------------------------------

// AppendSegment adds a segment to the end of the task string.
//
// If there is an inline comment, the segment is inserted before it. If the task
// is a comment line, the segment is prepended.
// The task is marked as dirty. Call Apply() to save the changes.
func (t *Task) AppendSegment(segment string) error {
	if segment == "" {
		return nil
	}

	currentText := t.String()
	if currentText == "" {
		t.SetText(segment)

		return nil
	}

	err := t.ensureParsed()
	if err != nil {
		return wrapError(err, "failed to re-parse dirty task before appending segment")
	}

	var newText string

	// For comment lines (starting with #), prepend the segment.
	if t.IsCommentLine() {
		newText = segment + spec.DelimSegments.String() + currentText
		t.SetText(newText)

		return nil
	}

	// If there is an inline comment, insert the segment before it.
	if t.HasInlineComment() {
		commentIdx := t.InlineCommentPos()
		if commentIdx != -1 {
			sep := string(currentText[commentIdx-1])
			beforeComment := currentText[:commentIdx-1]
			comment := currentText[commentIdx:]
			newText = beforeComment + spec.DelimSegments.String() + segment + sep + comment
		}
	} else {
		newText = currentText + spec.DelimSegments.String() + segment
	}

	t.SetText(newText)

	return nil
}

// InsertAfter adds a segment after the first occurrence of a target segment.
//
// It returns true if the segment was inserted, false otherwise.
// The task is marked as dirty. Call Apply() to save the changes.
func (t *Task) InsertAfter(target, segment string) bool {
	currentText := t.String()

	newText, inserted := insertSegmentAfterWithBoundaries(currentText, target, segment)
	if inserted {
		t.SetText(newText)
	}

	return inserted
}

// AddContext adds a context tag to the end of the task. The given context may
// include or omit the '@' prefix.
//
// This method will not:
// - Add contexts to comment lines.
// - Add duplicate contexts.
// - Add contexts with empty or invalid values.
//
// The task is marked as dirty. Call Apply() to save the changes.
func (t *Task) AddContext(context string) error {
	if t.isDirty {
		err := t.ensureParsed()
		if err != nil {
			return wrapError(err, "failed to re-parse dirty task before adding context")
		}
	}

	if t.IsCommentLine() {
		return nil
	}

	normalized, err := normalizeTag(context, spec.PrefixContext)
	if err != nil {
		return wrapError(err, ErrEmptyContextValue.Error())
	}

	tag := spec.PrefixContext.String() + normalized

	// Ignore duplicates.
	if _, _, exists := findSegmentBounds(t.String(), tag); exists {
		return nil
	}

	return t.AppendSegment(tag)
}

// AddProject adds a project tag to the task.
//
// It normalizes the project, prevents duplicates, and handles inline comments.
// It does not add projects to comment lines.
// The task is marked as dirty. Call Apply() to save the changes.
func (t *Task) AddProject(project string) error {
	if t.isDirty {
		err := t.ensureParsed()
		if err != nil {
			return wrapError(err, "failed to re-parse dirty task before adding project")
		}
	}

	if t.IsCommentLine() {
		return nil
	}

	normalized, err := normalizeTag(project, spec.PrefixProject)
	if err != nil {
		return wrapError(err, ErrEmptyProjectValue.Error())
	}

	tag := spec.PrefixProject.String() + normalized

	if _, _, exists := findSegmentBounds(t.String(), tag); exists {
		return nil
	}

	return t.AppendSegment(tag)
}

// ----------------------------------------------------------------------------
//  Modify (update existing components)
// ----------------------------------------------------------------------------

// SetText updates the raw task string and marks it dirty.
// Call Apply() to re-parse and sync structured data.
func (t *Task) SetText(newText string) {
	t.Parsed.SetText(newText)
	t.isDirty = true
}

// SetPriority sets the priority of the task.
//
// It validates the priority, replaces any existing priority, and inserts it
// correctly in the task string.
// The task is marked as dirty. Call Apply() to save the changes.
func (t *Task) SetPriority(priority string) error {
	// Validate given priority.
	if len(priority) != 1 || priority[0] < 'A' || priority[0] > 'Z' {
		return newError("invalid priority: must be a single uppercase letter A-Z, got: %s", priority)
	}

	priorityMark := "(" + priority + ")"

	if t.isDirty {
		err := t.ensureParsed()
		if err != nil {
			return wrapError(err, "failed to re-parse dirty task before setting priority")
		}
	}

	currentText := t.String()
	currentPriority := t.Priority()
	isDone := t.IsCompleted()

	doneMarker := spec.MarkerDone.String()
	sep := spec.DelimSegments.String()

	var newText string

	switch {
	case currentPriority != "":
		// Replace existing priority.
		oldPriority := "(" + currentPriority + ")"
		newText = replaceFirst(currentText, oldPriority, priorityMark)
	case isDone:
		// Insert priority after completion mark.
		donePrefix := doneMarker + sep
		newText = replaceFirst(currentText, donePrefix, donePrefix+priorityMark+sep)
	default:
		// Insert priority at the beginning.
		newText = priorityMark + sep + currentText
	}

	t.SetText(newText)
	t.isDirty = true

	return nil
}

// CompleteWithDate marks the task as completed with a specified date.
//
// Unlike Complete(), this method ALWAYS updates the completion date, even if
// the task is already completed.
// The task is marked as dirty. Call Apply() to save the changes.
func (t *Task) CompleteWithDate(date string) error {
	// Re-parse to ensure we have the latest data. It prevents inconsistencies
	// when this method is called multiple times before Apply().
	err := t.ensureParsed()
	if err != nil {
		return wrapError(err, "failed to re-parse dirty task before completing")
	}

	currentText := t.String()
	sep := spec.DelimSegments.String() // " "

	var newText string

	switch {
	case t.IsCompleted():
		// Already completed: replace completion date
		oldDate := t.DateCompleted()
		if oldDate != "" {
			// Has completion date: replace it
			newText = replaceFirst(currentText, oldDate, date)
		} else {
			// No completion date: insert after 'x'
			marker := spec.MarkerDone.String() + sep
			newText = replaceFirst(currentText, marker, marker+date+sep)
		}
	default:
		// Not completed: add completion marker and date
		// Check if task has priority in raw text - it needs to be moved after 'x'
		// Use raw text inspection instead of t.Priority() to handle dirty state
		// Format: x [(priority)] [date] remaining_text
		trimmed := strings.TrimSpace(currentText)

		var priority string

		textToMark := currentText

		if len(trimmed) > 2 && trimmed[0] == '(' && trimmed[2] == ')' {
			// Extract priority from raw text (e.g., "(A)" -> "A")
			priority = string(trimmed[1])
			// Remove priority from current position
			priorityMark := "(" + priority + ")" + sep
			textToMark = replaceFirst(currentText, priorityMark, "")
		}

		// Build new text: x [(priority)] [date] remaining
		newText = spec.MarkerDone.String() + sep

		if priority != "" {
			newText += "(" + priority + ")" + sep
		}

		if date != "" {
			newText += date + sep
		}

		newText += textToMark
	}

	t.SetText(newText)

	return nil
}

// Complete marks the task as completed.
//
// It adds 'x' marker at the beginning and sets the completion date to the
// current time. Use `CompleteWithDate` to specify a custom completion date.
//
// If the task is already completed, it does nothing (idempotent).
// To update the date of an already-completed task, use CompleteWithDate.
func (t *Task) Complete() error {
	if t.IsCompleted() {
		return nil // Idempotent: preserve original date
	}

	now := TimeNow().Format(spec.DateFormat)

	return t.CompleteWithDate(now)
}

// MarkDone marks the task as completed.
//
// It is an alias/shorthand for Complete. See Complete for details.
func (t *Task) MarkDone() error {
	return t.Complete()
}

// Reopen marks the task as incomplete.
//
// It removes 'x' marker and completion date. If the task is already incomplete,
// it does nothing.
// The task is marked as dirty. Call Apply() to save the changes.
func (t *Task) Reopen() error {
	err := t.ensureParsed()
	if err != nil {
		return wrapError(err, "failed to re-parse dirty task before reopening")
	}

	if !t.IsCompleted() {
		return nil // Already incomplete
	}

	currentText := t.String()

	// Remove completion marker
	marker := spec.MarkerDone.String() + spec.DelimSegments.String()
	newText := replaceFirst(currentText, marker, "")

	// Remove completion date if present
	dateCompleted := t.DateCompleted()
	if dateCompleted != "" {
		dateWithSpace := dateCompleted + spec.DelimSegments.String()
		newText = replaceFirst(newText, dateWithSpace, "")
	}

	t.SetText(newText)

	return nil
}

// MarkDirty forces the task state to dirty.
//
// This method is useful when external modifications are made to the task
// string via SetText(). It ensures that the next Apply() call will re-parse
// the task to update its structured data.
func (t *Task) MarkDirty() {
	t.isDirty = true
}

// ----------------------------------------------------------------------------
//  Delete (remove components from the task)
// ----------------------------------------------------------------------------

// RemoveSegment removes the first occurrence of a segment from the task.
//
// The task is marked as dirty. Call Apply() to save the changes.
func (t *Task) RemoveSegment(segment string) error {
	currentText := t.String()

	newText, removed := removeSegmentWithBoundaries(currentText, segment)
	if removed {
		t.SetText(newText)
	}

	return nil
}

// RemoveContext removes a context tag from the task.
//
// It is a shorthand for RemoveSegment with context normalization.
func (t *Task) RemoveContext(context string) error {
	if t.isDirty {
		err := t.ensureParsed()
		if err != nil {
			return wrapError(err, "failed to re-parse dirty task before removing context")
		}
	}

	if t.IsCommentLine() {
		return nil
	}

	normalized, err := normalizeTag(context, spec.PrefixContext)
	if err != nil {
		return wrapError(err, ErrEmptyContextValue.Error())
	}

	tag := spec.PrefixContext.String() + normalized

	return t.RemoveSegment(tag)
}

// RemovePriority removes the priority from the task.
//
// The task is marked as dirty. Call Apply() to save the changes.
func (t *Task) RemovePriority() error {
	if t.isDirty {
		err := t.ensureParsed()
		if err != nil {
			return wrapError(err, "failed to re-parse dirty task before removing priority")
		}
	}

	currentPriority := t.Priority()
	if currentPriority == "" {
		return nil // No priority to remove.
	}

	priorityMark := "(" + currentPriority + ")"

	return t.RemoveSegment(priorityMark)
}

// RemoveProject removes a project tag from the task.
//
// It is a shorthand for RemoveSegment with project normalization.
func (t *Task) RemoveProject(project string) error {
	if t.isDirty {
		err := t.ensureParsed()
		if err != nil {
			return wrapError(err, "failed to re-parse dirty task before removing project")
		}
	}

	if t.IsCommentLine() {
		return nil
	}

	normalized, err := normalizeTag(project, spec.PrefixProject)
	if err != nil {
		return wrapError(err, ErrEmptyProjectValue.Error())
	}

	tag := spec.PrefixProject.String() + normalized

	return t.RemoveSegment(tag)
}

// SetTag sets a key-value pair tag (e.g., "due:2024-12-25").
// Updates existing key or appends new one.
//
// Note: If the same key appears multiple times in the original text,
// SetTag updates only the first occurrence. The [Task.KeyValues] method
// returns all values including duplicates, preserving their original order.
// To avoid ambiguity, use a single key with comma-separated values
// (e.g., "tag:value1,value2") for multiple values.
func (t *Task) SetTag(key, value string) error {
	if t.isDirty {
		err := t.ensureParsed()
		if err != nil {
			return wrapError(err, "failed to re-parse dirty task before setting tag")
		}
	}

	if t.IsCommentLine() {
		return nil
	}

	// Validate key and value
	key = strings.TrimSpace(key)
	if key == "" {
		return newError("invalid key-value: key cannot be empty")
	}

	if strings.ContainsFunc(key, unicode.IsSpace) {
		return newError("invalid key-value: key cannot contain whitespace")
	}

	value = strings.TrimSpace(value)
	if value == "" {
		return newError("invalid key-value: value cannot be empty")
	}

	if strings.ContainsFunc(value, unicode.IsSpace) {
		return newError("invalid key-value: value cannot contain whitespace")
	}

	newTag := key + spec.SepKeyValue.String() + value

	// Check if the key already exists and update it (first occurrence)
	existingValues := t.KeyValues()
	for _, kv := range existingValues {
		if kv.Key == key {
			oldTag := key + spec.SepKeyValue.String() + kv.Value
			currentText := t.String()
			newText := replaceFirst(currentText, oldTag, newTag)
			t.SetText(newText)

			return nil
		}
	}

	// Key doesn't exist, append new tag
	return t.AppendSegment(newTag)
}

// RemoveTag removes a key-value pair tag by key name.
// Does nothing if key does not exist.
func (t *Task) RemoveTag(key string) error {
	if t.isDirty {
		err := t.ensureParsed()
		if err != nil {
			return wrapError(err, "failed to re-parse dirty task before removing tag")
		}
	}

	if t.IsCommentLine() {
		return nil
	}

	// Validate key
	key = strings.TrimSpace(key)
	if key == "" {
		return newError("invalid key-value: key cannot be empty")
	}

	if strings.ContainsFunc(key, unicode.IsSpace) {
		return newError("invalid key-value: key cannot contain whitespace")
	}

	// Check if the key exists (first occurrence)
	existingValues := t.KeyValues()
	for _, kv := range existingValues {
		if kv.Key == key {
			// Remove the key-value tag
			oldTag := key + spec.SepKeyValue.String() + kv.Value

			return t.RemoveSegment(oldTag)
		}
	}

	return nil // Key not found
}

// ----------------------------------------------------------------------------
//  Apply Changes (finalize modifications)
// ----------------------------------------------------------------------------

// Apply re-parses the task string if it has been modified (isDirty is true).
//
// Note that the given Options are applied linearly before re-parsing.
//
// It is safe to call this method multiple times. It only re-parses if there are
// pending changes (isDirty is true).
func (t *Task) Apply(opts ...Option) error {
	// Apply given options.
	if len(opts) > 0 {
		for _, opt := range opts {
			err := opt(t)
			if err != nil {
				return wrapError(err, "failed to apply option")
			}
		}
	}

	// Only re-parse if the task was modified.
	if t.isDirty {
		err := t.Parse(t.String())
		if err != nil {
			return wrapError(err, "failed to re-parse task after applying changes")
		}

		t.isDirty = false // Reset dirty flag after successful parse.
	}

	return nil
}

// ============================================================================
//  Private Methods
// ============================================================================

// ensureParsed re-parses the task if it is marked as dirty.
//
// This method is used internally to guarantee that the task's structured data
// is up-to-date before performing operations that depends on it.
func (t *Task) ensureParsed() error {
	if t.isDirty {
		err := t.Parse(t.String())
		if err == nil {
			t.isDirty = false
		}

		return wrapError(err, "failed to re-parse dirty task")
	}

	return nil
}

// ============================================================================
//  Helper Functions
// ============================================================================

// findSegmentBounds finds the start and end index of a segment.
//
// If the segment is not found, it returns false with -1 indices.
func findSegmentBounds(text, segment string) (int, int, bool) {
	const notFound = -1

	if segment == "" {
		return notFound, notFound, false
	}

	offset := 0

	for {
		idx := strings.Index(text[offset:], segment)
		if idx == -1 {
			return -1, -1, false
		}

		idx += offset
		start := idx
		end := idx + len(segment)

		beforeOK := start == 0 || text[start-1] == spec.DelimSegments.Byte()
		afterOK := end == len(text) || text[end] == spec.DelimSegments.Byte()

		if beforeOK && afterOK {
			return start, end, true
		}

		// offset = start + 1
		offset = end
	}
}

// insertSegmentAfterWithBoundaries inserts a segment after the first occurrence
// of a target segment. It respects word boundaries.
func insertSegmentAfterWithBoundaries(text, target, segment string) (string, bool) {
	_, end, ok := findSegmentBounds(text, target)
	if !ok {
		return text, false
	}

	before := text[:end]
	after := text[end:]

	if len(after) == 0 {
		return before + spec.DelimSegments.String() + segment, true
	}

	return before + spec.DelimSegments.String() + segment + after, true
}

// normalizeTag trims whitespace and an optional prefix from a tag.
//
// It returns an error if the tag is empty or contains spaces.
func normalizeTag(value string, mark spec.Mark) (string, error) {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return "", ErrValIsEmpty
	}

	if rune(normalized[0]) == mark.Rune() {
		normalized = strings.TrimSpace(normalized[1:])
	}

	if normalized == "" {
		return "", ErrValIsEmpty
	}

	if strings.IndexFunc(normalized, unicode.IsSpace) != -1 {
		return "", ErrValWithSpace
	}

	return normalized, nil
}

// removeSegmentWithBoundaries removes the first occurrence of a segment.
//
// It respects word boundaries to avoid removing parts of other words.
func removeSegmentWithBoundaries(text, segment string) (string, bool) {
	start, end, ok := findSegmentBounds(text, segment)
	if !ok {
		return text, false
	}

	before := text[:start]
	after := text[end:]

	if start == 0 {
		after = strings.TrimPrefix(after, spec.DelimSegments.String())

		return after, true
	}

	if end == len(text) {
		before = strings.TrimSuffix(before, spec.DelimSegments.String())

		return before, true
	}

	after = strings.TrimPrefix(after, spec.DelimSegments.String())

	return before + after, true
}

// replaceFirst replaces the first occurrence of a substring.
//
// It operates on the raw text to preserve formatting.
// It is a thin wrapper around strings.Replace with count=1.
func replaceFirst(s, old, replacement string) string {
	const firstFound = 1

	return strings.Replace(s, old, replacement, firstFound)
}
