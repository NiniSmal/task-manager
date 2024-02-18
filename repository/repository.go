package repository

import (
	"database/sql"
	"gitlab.com/nina8884807/task-manager/entity"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}
func (r *TaskRepository) SaveTask(task entity.Task) error {
	query := "INSERT INTO tasks (name, status, created_at) VALUES ($1, $2, $3) "

	_, err := r.db.Exec(query, task.Name, task.Status, task.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) GetTaskByID(id int64) (entity.Task, error) {
	query := "SELECT id, name, status, created_at FROM tasks WHERE id=$1"

	var task entity.Task

	err := r.db.QueryRow(query, id).Scan(&task.ID, &task.Name, &task.Status, &task.CreatedAt)
	if err != nil {
		return entity.Task{}, err
	}

	return task, nil
}
