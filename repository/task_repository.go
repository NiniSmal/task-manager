package repository

import (
	"context"
	"database/sql"
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
	query := "INSERT INTO tasks (name, status, created_at, user_id, project_id) VALUES ($1, $2, $3, $4, $5) "

	_, err := r.db.ExecContext(ctx, query, task.Name, task.Status, task.CreatedAt, task.UserID, task.ProjectID)
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) GetTaskByID(ctx context.Context, id int64) (entity.Task, error) {
	query := "SELECT id, name, status, created_at, user_id, project_id FROM tasks WHERE id=$1"

	var task entity.Task

	err := r.db.QueryRowContext(ctx, query, id).Scan(&task.ID, &task.Name, &task.Status, &task.CreatedAt, &task.UserID, &task.ProjectID)
	if err != nil {
		return entity.Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) GetTasks(ctx context.Context) ([]entity.Task, error) {
	query := "SELECT id, name, status, created_at, project_id FROM tasks"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []entity.Task

	for rows.Next() {
		var task entity.Task
		err = rows.Scan(&task.ID, &task.Name, &task.Status, &task.CreatedAt, &task.ProjectID)
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
