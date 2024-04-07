package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gitlab.com/nina8884807/task-manager/entity"
)

type TaskRepository struct {
	db  *sql.DB
	rds *redis.Client
}

func NewTaskRepository(db *sql.DB, rds *redis.Client) *TaskRepository {
	return &TaskRepository{
		db:  db,
		rds: rds,
	}
}
func (r *TaskRepository) SaveTask(ctx context.Context, task entity.Task) error {
	query := "INSERT INTO tasks (name, status, created_at, user_id) VALUES ($1, $2, $3, $4) "

	_, err := r.db.ExecContext(ctx, query, task.Name, task.Status, task.CreatedAt, task.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) GetTaskByID(ctx context.Context, id int64) (entity.Task, error) {
	query := "SELECT id, name, status, created_at, user_id FROM tasks WHERE id=$1"

	var task entity.Task

	err := r.db.QueryRowContext(ctx, query, id).Scan(&task.ID, &task.Name, &task.Status, &task.CreatedAt, &task.UserID)
	if err != nil {
		return entity.Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) GetUserTasks(ctx context.Context, userID int64) ([]entity.Task, error) {
	query := "SELECT id, name, status, created_at, user_id FROM tasks WHERE user_id = $1"

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []entity.Task

	for rows.Next() {
		var task entity.Task
		err = rows.Scan(&task.ID, &task.Name, &task.Status, &task.CreatedAt, &task.UserID)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
func (r *TaskRepository) GetTasks(ctx context.Context) ([]entity.Task, error) {
	query := "SELECT id, name, status, created_at FROM tasks"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []entity.Task

	for rows.Next() {
		var task entity.Task
		err = rows.Scan(&task.ID, &task.Name, &task.Status, &task.CreatedAt)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskRepository) UpdateTask(ctx context.Context, id int64, task entity.UpdateTask) error {
	query := "UPDATE tasks SET name = $1, status = $2 WHERE id = $3"

	_, err := r.db.ExecContext(ctx, query, task.Name, task.Status, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepository) GetUserIDBySessionID(ctx context.Context, sessionID uuid.UUID) (entity.User, error) {
	query := "SELECT user_id, role FROM sessions WHERE id = $1"

	var user entity.User

	err := r.db.QueryRowContext(ctx, query, sessionID).Scan(&user.ID, &user.Role)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *TaskRepository) GetUserSession(ctx context.Context, sessionID uuid.UUID) (entity.User, error) {
	var user string

	err := r.rds.Get(ctx, sessionID.String()).Scan(&user)
	if err != nil {
		return entity.User{}, err
	}

	var usRep entity.User
	err = json.Unmarshal([]byte(user), &usRep)
	if err != nil {
		return entity.User{}, err
	}
	return usRep, nil
}
