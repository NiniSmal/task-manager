package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"

	"gitlab.com/nina8884807/task-manager/entity"
)

type TaskService struct {
	tasks    TaskRepository
	projects ProjectRepository
	kafka    *kafka.Writer
}

func NewTaskService(r TaskRepository, pr ProjectRepository, w *kafka.Writer) *TaskService {
	return &TaskService{
		tasks:    r,
		projects: pr,
		kafka:    w,
	}
}

type TaskRepository interface {
	Create(ctx context.Context, task entity.Task) (entity.Task, error)
	ByID(ctx context.Context, id int64) (entity.Task, error)
	Tasks(ctx context.Context, f entity.TaskFilter) ([]entity.Task, error)
	Update(ctx context.Context, id int64, task entity.UpdateTask) error
	Delete(ctx context.Context, id int64) error
}

func (s *TaskService) AddTask(ctx context.Context, task entity.Task) (entity.Task, error) {
	err := task.Validate()
	if err != nil {
		return entity.Task{}, err
	}

	user := ctx.Value("user").(entity.User)

	users, err := s.projects.ProjectUsers(ctx, task.ProjectID)
	if err != nil {
		return entity.Task{}, err
	}

	err = isUserInProject(users, user.ID)
	if err != nil {
		return entity.Task{}, err
	}

	task.Status = entity.StatusNotDone
	task.CreatedAt = time.Now()
	task.UserID = user.ID

	taskDB, err := s.tasks.Create(ctx, task)
	if err != nil {
		return entity.Task{}, err
	}

	err = s.sendCreateTaskNotification(ctx, task.ProjectID)
	if err != nil {
		l := ctx.Value("logger").(*slog.Logger)
		l.Error("sendCreateTaskNotification", "err", err)
		err = nil
	}

	return taskDB, nil
}

func (s *TaskService) sendCreateTaskNotification(ctx context.Context, projectID int64) error {
	// при создании задачи в проекте слать уведомление об этом всем участникам проекта
	userEmails, err := s.projects.ProjectUsers(ctx, projectID)
	if err != nil {
		return err
	}

	user := ctx.Value("user").(entity.User)

	for _, userTo := range userEmails {
		if userTo.ID == user.ID {
			continue
		}

		email := Email{
			Text:    fmt.Sprintf("New task created in project  %d", projectID),
			To:      userTo.Email,
			Subject: "New task",
		}

		msg, err := json.Marshal(&email)
		if err != nil {
			return fmt.Errorf("failed to marshal message: , %w", err)
		}

		err = s.kafka.WriteMessages(ctx, kafka.Message{Value: msg})
		if err != nil {
			return fmt.Errorf("failed to write messages %w", err)
		}
	}

	return nil
}

func isUserInProject(users []entity.User, id int64) error {
	for _, userM := range users {
		if userM.ID == id {
			return nil
		}
	}
	return entity.ErrForbidden
}

func (s *TaskService) GetTask(ctx context.Context, id int64) (entity.Task, error) {
	user := ctx.Value("user").(entity.User)

	task, err := s.tasks.ByID(ctx, id)
	if err != nil {
		return entity.Task{}, fmt.Errorf("get task by %d: %w", id, err)
	}
	if user.Role == entity.RoleAdmin {
		return task, nil
	}
	if task.UserID == user.ID {
		return task, nil
	} else {
		return entity.Task{}, err
	}
}

func (s *TaskService) GetAllTasks(ctx context.Context, f entity.TaskFilter) ([]entity.Task, error) {
	user, ok := ctx.Value("user").(entity.User)
	if !ok {
		return nil, entity.ErrNotAuthenticated
	}

	f.UserID = strconv.FormatInt(user.ID, 10)

	tasks, err := s.tasks.Tasks(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("get all tasks: %w", err)
	}
	return tasks, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id int64, task entity.UpdateTask) (entity.Task, error) {
	user := ctx.Value("user").(entity.User)

	if task.Status != entity.StatusDone && task.Status != entity.StatusNotDone {
		return entity.Task{}, fmt.Errorf("update task: %w", entity.ErrForbidden)
	}

	taskOld, err := s.tasks.ByID(ctx, id)
	if err != nil {
		return entity.Task{}, err
	}

	if user.Role == entity.RoleAdmin {
		err = s.tasks.Update(ctx, id, task)
		if err != nil {
			return entity.Task{}, err
		}
	}
	if user.ID == taskOld.UserID {

		err = s.tasks.Update(ctx, id, task)
		if err != nil {
			return entity.Task{}, err
		}
	}
	taskUp, err := s.tasks.ByID(ctx, id)
	return taskUp, nil
}

func (s *TaskService) Delete(ctx context.Context, id int64) error {
	user := ctx.Value("user").(entity.User)
	taskOld, err := s.tasks.ByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get task by id %d: %w", id, err)
	}
	if user.ID == taskOld.UserID {
		err = s.tasks.Delete(ctx, id)
		if err != nil {
			return fmt.Errorf("delete task %d: %w", id, err)
		}
	}
	return nil
}
