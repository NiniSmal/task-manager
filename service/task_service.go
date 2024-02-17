package service

import (
	"errors"
	"fmt"
	"gitlab.com/nina8884807/task-manager/entity"
)

type TaskService struct {
	repo Repository
}

func NewTaskService(r Repository) *TaskService {
	return &TaskService{
		repo: r,
	}
}

type Repository interface {
	SaveTask(task entity.Task) error
}

func (s *TaskService) AddTask(task entity.Task) error {
	if task.Name == "" {
		return errors.New("name is empty")
	}

	task.Status = "not done"

	err := s.repo.SaveTask(task)
	if err != nil {
		return fmt.Errorf("save task: %w", err)
	}

	return nil
}
