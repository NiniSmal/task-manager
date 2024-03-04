package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
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

func (r *UserRepository) GetUserIDByLoginAndPassword(ctx context.Context, user entity.User) (int64, error) {
	query := "SELECT id FROM users WHERE login = $1 AND password = $2"

	var id int64

	err := r.db.QueryRowContext(ctx, query, user.Login, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserRepository) SaveSession(ctx context.Context, sessionID uuid.UUID, userID int64) error {
	query := "INSERT INTO sessions (id, user_id) VALUES ($1, $2)"

	_, err := r.db.ExecContext(ctx, query, sessionID, userID)
	if err != nil {
		return err
	}

	return nil
}
