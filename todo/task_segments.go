package todo

import (
	"fmt"
	"sort"
)

// segmentBuilder helps build task segments.
type segmentBuilder struct {
	segs []*TaskSegment
}

// addBasic adds a basic segment.
func (sb *segmentBuilder) addBasic(t TaskSegmentType, s string) {
	sb.segs = append(sb.segs, &TaskSegment{
		Type:      t,
		Originals: []string{s},
		Display:   s,
	})
}

// addWithDisplay adds a segment with original and display.
func (sb *segmentBuilder) addWithDisplay(t TaskSegmentType, orig, display string) {
	sb.segs = append(sb.segs, &TaskSegment{
		Type:      t,
		Originals: []string{orig},
		Display:   display,
	})
}

// addTag adds a tag segment.
func (sb *segmentBuilder) addTag(key, val string) {
	sb.segs = append(sb.segs, &TaskSegment{
		Type:      SegmentTag,
		Originals: []string{key, val},
		Display:   fmt.Sprintf("%s:%s", key, val),
	})
}

// addCompletedSegments adds completed status and date segments.
func (sb *segmentBuilder) addCompletedSegments(task *Task) {
	if task.Completed {
		sb.addBasic(SegmentIsCompleted, "x")

		if task.HasCompletedDate() {
			sb.addBasic(SegmentCompletedDate, task.CompletedDate.Format(DateLayout))
		}
	}
}

// addPrioritySegment adds priority segment.
func (sb *segmentBuilder) addPrioritySegment(task *Task) {
	if task.HasPriority() && (!task.Completed || !RemoveCompletedPriority) {
		sb.addWithDisplay(SegmentPriority, task.Priority, fmt.Sprintf("(%s)", task.Priority))
	}
}

// addCreatedDateSegment adds created date segment.
func (sb *segmentBuilder) addCreatedDateSegment(task *Task) {
	if task.HasCreatedDate() {
		sb.addBasic(SegmentCreatedDate, task.CreatedDate.Format(DateLayout))
	}
}

// addTodoTextSegment adds todo text segment.
func (sb *segmentBuilder) addTodoTextSegment(task *Task) {
	sb.addBasic(SegmentTodoText, task.Todo)
}

// addContextSegments adds context segments.
func (sb *segmentBuilder) addContextSegments(task *Task) {
	if task.HasContexts() {
		sortedContexts := make([]string, len(task.Contexts))
		copy(sortedContexts, task.Contexts)
		sort.Strings(sortedContexts)

		for _, context := range sortedContexts {
			sb.addWithDisplay(SegmentContext, context, contextPrefix+context)
		}
	}
}

// addProjectSegments adds project segments.
func (sb *segmentBuilder) addProjectSegments(task *Task) {
	if task.HasProjects() {
		sortedProjects := make([]string, len(task.Projects))
		copy(sortedProjects, task.Projects)
		sort.Strings(sortedProjects)

		for _, project := range sortedProjects {
			sb.addWithDisplay(SegmentProject, project, projectPrefix+project)
		}
	}
}

// addTagSegments adds additional tag segments.
func (sb *segmentBuilder) addTagSegments(task *Task) {
	if task.HasAdditionalTags() {
		keys := make([]string, 0, len(task.AdditionalTags))
		for key := range task.AdditionalTags {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		for _, key := range keys {
			sb.addTag(key, task.AdditionalTags[key])
		}
	}
}

// addDueDateSegment adds due date segment.
func (sb *segmentBuilder) addDueDateSegment(task *Task) {
	if task.HasDueDate() {
		sb.addBasic(SegmentDueDate, duePrefix+task.DueDate.Format(DateLayout))
	}
}

// Segments returns a segmented task string in todo.txt format. The order of
// segments is the same as String.
func (task *Task) Segments() []*TaskSegment {
	//nolint:exhaustruct // segs field is initialized in add methods
	segmentBuilder := &segmentBuilder{}

	segmentBuilder.addCompletedSegments(task)
	segmentBuilder.addPrioritySegment(task)
	segmentBuilder.addCreatedDateSegment(task)
	segmentBuilder.addTodoTextSegment(task)
	segmentBuilder.addContextSegments(task)
	segmentBuilder.addProjectSegments(task)
	segmentBuilder.addTagSegments(task)
	segmentBuilder.addDueDateSegment(task)

	return segmentBuilder.segs
}
