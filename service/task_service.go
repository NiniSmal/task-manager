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
	UpdateTask(ctx context.Context, task entity.Task) error
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

func (s *TaskService) UpdateTask(ctx context.Context, task entity.Task) error {
	if task.Status != entity.StatusDone && task.Status != entity.StatusNotDone {
		return errors.New("status  is not correct")
	}
	_, err := s.repo.GetTaskByID(ctx, task.ID)
	if err != nil {
		return fmt.Errorf("get task by id: %w", err)
	}

	err = s.repo.UpdateTask(ctx, task)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}
	return nil
}
