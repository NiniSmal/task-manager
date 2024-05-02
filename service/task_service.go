package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"gitlab.com/nina8884807/task-manager/entity"
	"time"
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
	Create(ctx context.Context, task entity.Task) error
	ByID(ctx context.Context, id int64) (entity.Task, error)
	Tasks(ctx context.Context, f entity.TaskFilter) ([]entity.Task, error)
	Update(ctx context.Context, id int64, task entity.UpdateTask) error
}

func (s *TaskService) AddTask(ctx context.Context, task entity.Task) error {
	err := task.Validate()
	if err != nil {
		return err
	}

	user := ctx.Value("user").(entity.User)

	users, err := s.projects.ProjectUsers(ctx, task.ProjectID)
	if err != nil {
		return err
	}

	err = isUserInProject(users, user.ID)
	if err != nil {
		return err
	}

	task.Status = entity.StatusNotDone
	task.CreatedAt = time.Now()
	task.UserID = user.ID

	err = s.tasks.Create(ctx, task)
	if err != nil {
		return fmt.Errorf("save task: %w", err)
	}
	//при создании задачи в проекте слать уведомление об этом всем участникам проекта
	usersEmail, err := s.projects.ProjectUsers(ctx, task.ProjectID)
	if err != nil {
		return err
	}
	var emails []SendEmail
	for _, userTo := range usersEmail {
		email := SendEmail{
			Text: fmt.Sprintf("New task created in project  %d", task.ProjectID),
			To:   userTo.Email,
		}
		emails = append(emails, email)
	}
	msg, err := json.Marshal(&emails)
	if err != nil {
		return fmt.Errorf("failed to marshal message: , %w", err)
	}

	err = s.kafka.WriteMessages(ctx, kafka.Message{Value: msg})
	if err != nil {
		return fmt.Errorf("failed to write messages %w", err)
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
	tasks, err := s.tasks.Tasks(ctx, f)
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

	taskOld, err := s.tasks.ByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get task by id: %w", err)
	}

	if user.Role == entity.RoleAdmin {
		err = s.tasks.Update(ctx, id, task)
		if err != nil {
			return fmt.Errorf("update task: %w", err)
		}
	}
	if user.ID == taskOld.UserID {
		err = s.tasks.Update(ctx, id, task)
		if err != nil {
			return fmt.Errorf("update task: %w", err)
		}
	}
	return nil
}
