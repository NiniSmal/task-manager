package entity

import (
	"fmt"
	"time"
	"unicode/utf8"
)

const (
	minNameProject = 5
	maxNameProject = 200
)

type Project struct {
	ID        int64     ` json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    int64     `json:"user_id"`
	Members   []User    `json:"members"`
}

type ProjectFilter struct {
	UserID int64
}

func (project *Project) Validate() error {
	rp := utf8.RuneCountInString(project.Name)
	if rp < minNameProject {
		return fmt.Errorf("%w: the name project must be minimum %d symbols", ErrValidate, minNameProject)
	}
	if rp > maxNameProject {
		return fmt.Errorf("%w: the name project can be maximum %d symbols", ErrValidate, maxNameProject)
	}
	return nil
}
