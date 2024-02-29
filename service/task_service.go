package service

import (
	"context"
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
	SaveTask(ctx context.Context, task entity.Task) error
	GetTaskByID(ctx context.Context, id int64) (entity.Task, error)
	GetAllTasks(ctx context.Context) ([]entity.Task, error)
}

func (s *TaskService) AddTask(ctx context.Context, task entity.Task) error {
	if task.Name == "" {
		return errors.New("name is empty")
	}

	task.Status = entity.StatusNotDone
	task.CreatedAt = time.Now()

	err := s.repo.SaveTask(ctx, task)
	if err != nil {
		return fmt.Errorf("save task: %w", err)
	}
	return nil
}

func (s *TaskService) GetTask(ctx context.Context, id int64) (entity.Task, error) {
	task, err := s.repo.GetTaskByID(ctx, id)
	if err != nil {
		return entity.Task{}, fmt.Errorf("get task: %w", err)
	}

	return task, nil
}

func (s *TaskService) GetAllTasks(ctx context.Context) ([]entity.Task, error) {
	tasks, err := s.repo.GetAllTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all tasks: %w", err)
	}
	return tasks, nil
}
