package repository

import (
	"context"
	"database/sql"
	"gitlab.com/nina8884807/task-manager/entity"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user entity.User) error {
	query := "INSERT INTO users ( login, password, created_at ) VALUES ($1, $2, $3)"

	_, err := r.db.ExecContext(ctx, query, user.Login, user.Password, user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
