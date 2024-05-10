package entity

import (
	"fmt"
	"time"
	"unicode/utf8"
)

const (
	minNameTask = 5
	maxNameTask = 200
)

const (
	StatusNotDone = "not_done"
	StatusDone    = "done"
)

type Task struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      int64     `json:"user_id"`
	ProjectID   int64     `json:"project_id"`
}

type UpdateTask struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UserID      int64  `json:"user_id"`
	ProjectID   int64  `json:"project_id"`
}

type TaskFilter struct {
	UserID    string
	ProjectID string
}

func (task *Task) Validate() error {
	rt := utf8.RuneCountInString(task.Name)
	if rt < minNameTask {
		return fmt.Errorf("%w: the name task must be minimum %d symbols", ErrValidate, minNameTask)
	}
	if rt > maxNameTask {
		return fmt.Errorf("%w: the name task can be maximum %d symbols", ErrValidate, maxNameTask)
	}
	if task.ProjectID == 0 {
		return fmt.Errorf("%w: project ID is empty", ErrValidate)
	}

	return nil
}
