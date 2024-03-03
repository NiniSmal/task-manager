package repository

import (
	"context"
	"database/sql"
	"gitlab.com/nina8884807/task-manager/entity"
	"time"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}
func (r *TaskRepository) SaveTask(ctx context.Context, task entity.Task) error {
	query := "INSERT INTO tasks (name, status, created_at) VALUES ($1, $2, $3) "

	time.Sleep(3 * time.Second)

	_, err := r.db.ExecContext(ctx, query, task.Name, task.Status, task.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) GetTaskByID(ctx context.Context, id int64) (entity.Task, error) {
	query := "SELECT id, name, status, created_at FROM tasks WHERE id=$1"

	var task entity.Task

	err := r.db.QueryRowContext(ctx, query, id).Scan(&task.ID, &task.Name, &task.Status, &task.CreatedAt)
	if err != nil {
		return entity.Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) GetAllTasks(ctx context.Context) ([]entity.Task, error) {
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

func (r *TaskRepository) UpdateTask(ctx context.Context, task entity.Task) error {
	query := "UPDATE tasks SET name = $1, status = $2 WHERE id = $3"

	_, err := r.db.ExecContext(ctx, query, task.Name, task.Status, task.ID)
	if err != nil {
		return err
	}
	return nil
}
