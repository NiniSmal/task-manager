package repository

import (
	"context"
	"database/sql"
	"fmt"
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
func (r *TaskRepository) Create(ctx context.Context, task entity.Task) (entity.Task, error) {
	query := "INSERT INTO tasks ( name, description, status, created_at, user_id, project_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	err := r.db.QueryRowContext(ctx, query, task.Name, task.Description, task.Status, task.CreatedAt, task.UserID, task.ProjectID).Scan(&task.ID)
	if err != nil {
		return entity.Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) ByID(ctx context.Context, id int64) (entity.Task, error) {
	query := "SELECT id, name, description, status, created_at, user_id, project_id FROM tasks WHERE id = $1"

	var task entity.Task

	err := r.db.QueryRowContext(ctx, query, id).Scan(&task.ID, &task.Name, &task.Description, &task.Status, &task.CreatedAt, &task.UserID, &task.ProjectID)
	if err != nil {
		return entity.Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) Tasks(ctx context.Context, f entity.TaskFilter) ([]entity.Task, error) {
	query := "SELECT id, name, description, status, created_at, project_id, user_id FROM tasks"

	query, args := applyTaskFilter(query, f)
	query += " ORDER BY created_at DESC"
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []entity.Task

	for rows.Next() {
		var task entity.Task
		err = rows.Scan(&task.ID, &task.Name, &task.Description, &task.Status, &task.CreatedAt, &task.ProjectID, &task.UserID)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
func applyTaskFilter(query string, f entity.TaskFilter) (string, []any) {
	var args []any
	where := ""

	if f.UserID != "" {
		args = append(args, f.UserID)
		where += fmt.Sprintf("user_id = $%d", len(args))
	}
	if f.ProjectID != "" {
		args = append(args, f.ProjectID)
		if where != "" {
			where += " AND "
		}
		where += fmt.Sprintf("project_id = $%d", len(args))
	}
	if where != "" {
		query += " WHERE " + where
	}

	return query, args
}

func (r *TaskRepository) Update(ctx context.Context, taskID int64, task entity.UpdateTask) error {
	query := "UPDATE tasks SET name = $1, description = $2, status = $3, user_id = $4, project_id = $5 WHERE id = $6"

	_, err := r.db.ExecContext(ctx, query, task.Name, task.Description, task.Status, task.UserID, task.ProjectID, taskID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM tasks WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
