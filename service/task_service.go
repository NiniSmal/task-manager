package service

import (
	"errors"
	"fmt"
	"gitlab.com/nina8884807/task-manager/entity"
	"time"
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
	GetTaskByID(id int64) (entity.Task, error)
}

func (s *TaskService) AddTask(task entity.Task) error {
	if task.Name == "" {
		return errors.New("name is empty")
	}

	task.Status = entity.StatusNotDone
	task.CreatedAt = time.Now()

	err := s.repo.SaveTask(task)
	if err != nil {
		return fmt.Errorf("save task: %w", err)
	}

	return nil
}

func (s *TaskService) GetTaskByID(id int64) error {

	_, err := s.repo.GetTaskByID(id)
	if err != nil {
		return fmt.Errorf("get task: %w", err)
	}
	return nil
}
