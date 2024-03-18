package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"gitlab.com/nina8884807/task-manager/entity"
	"unicode/utf8"
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

func (r *UserRepository) SaveSession(ctx context.Context, sessionID uuid.UUID, user entity.User) error {
	query := "INSERT INTO sessions (id, user_id, role) VALUES ($1, $2, $3)"

	_, err := r.db.ExecContext(ctx, query, sessionID, user.ID, user.Role)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Validate(user entity.User) error {
	rl := utf8.RuneCountInString(user.Login)
	if rl < 1 || rl > 15 {
		return errors.New("the login must be minimum 1 symbol and not more 15 symbols")
	}

	if utf8.RuneCountInString(user.Password) <= 0 {
		return errors.New("the password must be more than 0")
	}
	return nil
}
