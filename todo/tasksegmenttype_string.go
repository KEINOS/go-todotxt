// Code generated by "stringer -type TaskSegmentType -trimprefix Segment -output tasksegmenttype_string.go"; DO NOT EDIT.

package todo

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SegmentIsCompleted-2]
	_ = x[SegmentCompletedDate-3]
	_ = x[SegmentPriority-4]
	_ = x[SegmentCreatedDate-5]
	_ = x[SegmentTodoText-6]
	_ = x[SegmentContext-7]
	_ = x[SegmentProject-8]
	_ = x[SegmentTag-9]
	_ = x[SegmentDueDate-10]
}

const _TaskSegmentType_name = "IsCompletedCompletedDatePriorityCreatedDateTodoTextContextProjectTagDueDate"

var _TaskSegmentType_index = [...]uint8{0, 11, 24, 32, 43, 51, 58, 65, 68, 75}

func (i TaskSegmentType) String() string {
	i -= 2
	if i >= TaskSegmentType(len(_TaskSegmentType_index)-1) {
		return "TaskSegmentType(" + strconv.FormatInt(int64(i+2), 10) + ")"
	}
	return _TaskSegmentType_name[_TaskSegmentType_index[i]:_TaskSegmentType_index[i+1]]
}
