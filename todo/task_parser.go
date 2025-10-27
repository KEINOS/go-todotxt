package todo

import (
	"strings"

	"github.com/pkg/errors"
)

// taskParser handles parsing of a todo.txt task string into a Task struct.
type taskParser struct {
	text string
	task *Task
}

// newTaskParser creates a new taskParser instance.
func newTaskParser(text string) *taskParser {
	oriText := strings.Trim(text, whitespaces)
	task := new(Task)

	task.Original = oriText
	task.Todo = oriText
	task.clock = realClock{}

	return &taskParser{
		text: oriText,
		task: task,
	}
}

// parse performs the full parsing of the task.
func (p *taskParser) parse() (*Task, error) {
	err := p.parseCompleted()
	if err != nil {
		return nil, err
	}

	p.parsePriority()

	err = p.parseCreatedDate()
	if err != nil {
		return nil, err
	}

	p.parseContexts()
	p.parseProjects()

	err = p.parseAdditionalTags()
	if err != nil {
		return nil, err
	}

	p.finalizeTodo()

	return p.task, nil
}

// parseCompleted checks and parses the completed status.
func (p *taskParser) parseCompleted() error {
	if completedRx.MatchString(p.text) {
		err := parseCompleted(p.text, p.task)
		if err != nil {
			return errors.Wrap(err, "failed to parse completed status")
		}
	}

	return nil
}

// parsePriority checks and parses the priority.
func (p *taskParser) parsePriority() {
	if priorityRx.MatchString(p.text) {
		parsePriority(p.text, p.task)
	}
}

// parseCreatedDate checks and parses the created date.
func (p *taskParser) parseCreatedDate() error {
	if createdDateRx.MatchString(p.text) {
		err := parseCreatedDate(p.text, p.task)
		if err != nil {
			return errors.Wrap(err, "failed to parse created date")
		}
	}

	return nil
}

// parseContexts checks and parses contexts.
func (p *taskParser) parseContexts() {
	if contextRx.MatchString(p.text) {
		p.task.Contexts = getSlice(p.text, contextRx)
		p.task.Todo = contextRx.ReplaceAllString(p.task.Todo, emptyStr)
	}
}

// parseProjects checks and parses projects.
func (p *taskParser) parseProjects() {
	if projectRx.MatchString(p.text) {
		p.task.Projects = getSlice(p.text, projectRx)
		p.task.Todo = projectRx.ReplaceAllString(p.task.Todo, emptyStr)
	}
}

// parseAdditionalTags checks and parses additional tags.
func (p *taskParser) parseAdditionalTags() error {
	if addonTagRx.MatchString(p.text) {
		err := parseAdditionalTags(p.text, p.task)
		if err != nil {
			return errors.Wrap(err, "failed to parse additional tags")
		}
	}

	return nil
}

// finalizeTodo trims whitespaces from the Todo text.
func (p *taskParser) finalizeTodo() {
	p.task.Todo = strings.Trim(p.task.Todo, "\t\n\r\f ")
}
