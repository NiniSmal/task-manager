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
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UserID    int64     `json:"user_id"`
	ProjectID int64     `json:"project_id"`
}

type UpdateTask struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	UserID int64  `json:"user_id"`
}

type TaskFilter struct {
	UserID    string `json:"user_id"`
	ProjectID string `json:"project_id"`
}

func (task *Task) Validate() error {
	rt := utf8.RuneCountInString(task.Name)
	if rt < minNameTask {
		return fmt.Errorf("the name task must be minimum %d symbols", minNameTask)
	}
	if rt > maxNameTask {
		return fmt.Errorf("the name task can be maximum %d symbols", maxNameTask)
	}
	return nil
}
