package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"

	"gitlab.com/nina8884807/task-manager/entity"
)

type TaskService struct {
	tasks    TaskRepository
	projects ProjectRepository
	users    UserRepository
	kafka    *kafka.Writer
}

func NewTaskService(r TaskRepository, pr ProjectRepository, us UserRepository, w *kafka.Writer) *TaskService {
	return &TaskService{
		tasks:    r,
		projects: pr,
		users:    us,
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

	if !slices.ContainsFunc(users, func(u entity.User) bool {
		return u.ID == user.ID
	}) {
		return entity.Task{}, fmt.Errorf("user %d is not a member of project %d", user.ID, task.ProjectID)
	}

	task.Status = entity.StatusNotDone
	task.CreatedAt = time.Now()
	task.UserID = user.ID

	if task.AssignerID == 0 {
		task.AssignerID = user.ID
	}

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

	if f.ProjectID != "" {
		projectID, err := strconv.ParseInt(f.ProjectID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse project id: %w", err)
		}

		users, err := s.projects.ProjectUsers(ctx, projectID)
		if err != nil {
			return nil, fmt.Errorf("get project users: %w", err)
		}

		if !slices.ContainsFunc(users, func(u entity.User) bool {
			return u.ID == user.ID
		}) {
			return nil, fmt.Errorf("user %d is not a member of project %d", user.ID, projectID)
		}
	}

	tasks, err := s.tasks.Tasks(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("get all tasks: %w", err)
	}
	return tasks, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id int64, task entity.UpdateTask) (entity.Task, error) {
	user := ctx.Value("user").(entity.User)

	_, ok := entity.Statuses[task.Status]
	if !ok {
		return entity.Task{}, fmt.Errorf("update task: %w", entity.ErrValidate)
	}

	if task.AssignerID == 0 {
		task.AssignerID = user.ID
	}

	members, err := s.projects.ProjectUsers(ctx, task.ProjectID)
	if err != nil {
		return entity.Task{}, fmt.Errorf("get project users: %w", entity.ErrNotFound)
	}

	if !slices.ContainsFunc(members, func(u entity.User) bool {
		return u.ID == user.ID
	}) {
		return entity.Task{}, fmt.Errorf("user %d is not a member of project %d", user.ID, task.ProjectID)
	}

	err = s.tasks.Update(ctx, id, task)
	if err != nil {
		return entity.Task{}, fmt.Errorf("update task: %w", err)
	}

	taskUp, err := s.tasks.ByID(ctx, id)
	if err != nil {
		return entity.Task{}, err
	}

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
