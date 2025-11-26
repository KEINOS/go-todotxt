package parse

// Components represents all parsed components from a todo.txt task string.
//
// This struct is used internally to store parsed segments to ease testing
// and further processing.
type Components struct {
	Priority         string     `json:"priority,omitempty"`
	DateCompleted    string     `json:"dateCompleted,omitempty"`
	DateCreated      string     `json:"dateCreated,omitempty"`
	Description      string     `json:"description"`       // without completed marker, priority and dates
	Comment          string     `json:"comment,omitempty"` // including '#'
	Contexts         []string   `json:"contexts,omitempty"`
	Projects         []string   `json:"projects,omitempty"`
	KeyValues        []KeyValue `json:"keyValues,omitempty"`
	PosInlineComment int        `json:"posInlineComment,omitempty"`
	IsDone           bool       `json:"isDone"`
	IsCommentLine    bool       `json:"isCommentLine"`
	HasInlineComment bool       `json:"hasInlineComment,omitempty"`
}

// KeyValue represents a `key:value` tag found in a todo.txt task string.
type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
