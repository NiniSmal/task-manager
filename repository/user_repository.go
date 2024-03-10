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
	query := "INSERT INTO users ( login, password, created_at, role ) VALUES ($1, $2, $3, $4)"

	_, err := r.db.ExecContext(ctx, query, user.Login, user.Password, user.CreatedAt, user.Role)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUserIDByLoginAndPassword(ctx context.Context, user entity.User) (int64, string, error) {
	query := "SELECT id, role FROM users WHERE login = $1 AND password = $2"

	var id int64
	var role string

	err := r.db.QueryRowContext(ctx, query, user.Login, user.Password).Scan(&id, &role)
	if err != nil {
		return 0, "", err
	}

	return id, role, nil
}

func (r *UserRepository) SaveSession(ctx context.Context, sessionID uuid.UUID, userID int64, role string) error {
	query := "INSERT INTO sessions (id, user_id, role) VALUES ($1, $2, $3)"

	_, err := r.db.ExecContext(ctx, query, sessionID, userID, role)
	if err != nil {
		return err
	}

	return nil
}
