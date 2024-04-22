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
	GetUserTasks(ctx context.Context, userID int64) ([]entity.Task, error)
	GetTasks(ctx context.Context) ([]entity.Task, error)
	UpdateTask(ctx context.Context, id int64, task entity.UpdateTask) error
}

func (s *TaskService) AddTask(ctx context.Context, task entity.Task) error {
	if task.Name == "" {
		return errors.New("name is empty")
	}

	task.Status = entity.StatusNotDone
	task.CreatedAt = time.Now()

	user := ctx.Value("user").(entity.User)

	task.UserID = user.ID

	err := s.repo.SaveTask(ctx, task)
	if err != nil {
		return fmt.Errorf("save task: %w", err)
	}
	return nil
}

func (s *TaskService) GetTask(ctx context.Context, id int64) (entity.Task, error) {
	user := ctx.Value("user").(entity.User)

	task, err := s.repo.GetTaskByID(ctx, id)
	if err != nil {
		return entity.Task{}, fmt.Errorf("get task by %d: %w", id, err)
	}
	if user.Role == entity.RoleAdmin {
		return task, nil
	}
	if task.UserID == user.ID {
		return task, nil
	} else {
		return entity.Task{}, fmt.Errorf("get task by %d: %w", id, err)
	}
}

func (s *TaskService) GetAllTasks(ctx context.Context) ([]entity.Task, error) {
	user := ctx.Value("user").(entity.User)

	if user.Role == entity.RoleAdmin {
		tasks, err := s.repo.GetTasks(ctx)
		if err != nil {
			return nil, fmt.Errorf("get all tasks: %w", err)
		}
		return tasks, nil
	}
	tasks, err := s.repo.GetUserTasks(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("get all tasks: %w", err)
	}
	return tasks, nil

}

func (s *TaskService) UpdateTask(ctx context.Context, id int64, task entity.UpdateTask) error {
	user := ctx.Value("user").(entity.User)

	if task.Status != entity.StatusDone && task.Status != entity.StatusNotDone {
		return errors.New("status  is not correct")
	}

	_, err := s.repo.GetTaskByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get task by id: %w", err)
	}

	if user.Role == entity.RoleAdmin {
		err = s.repo.UpdateTask(ctx, id, task)
		if err != nil {
			return fmt.Errorf("update task: %w", err)
		}
	}
	if user.ID == task.UserID {
		err = s.repo.UpdateTask(ctx, id, task)
		if err != nil {
			return fmt.Errorf("update task: %w", err)
		}
	}
	return nil
}
