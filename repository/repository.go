package repository

import (
	"database/sql"
	"gitlab.com/nina8884807/task-manager/entity"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(r *sql.DB) *TaskRepository {
	return &TaskRepository{
		db: r,
	}
}
func (r *TaskRepository) SaveTask(task entity.Task) error {
	query := "INSERT INTO tasks (name, status) VALUES ($1, $2)"
	_, err := r.db.Exec(query, task.Name, task.Status)
	if err != nil {
		return err
	}
	return nil
}

//create tabl - task.sql
