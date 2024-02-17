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
	query := "INSERT INTO tasks (name, status) VALUES ($1, $2) "

	_, err := r.db.Exec(query, task.Name, task.Status)
	if err != nil {
		return err
	}

	return nil
}
