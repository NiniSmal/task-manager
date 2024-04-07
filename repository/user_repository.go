package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gitlab.com/nina8884807/task-manager/entity"
)

type UserRepository struct {
	db  *sql.DB
	rds *redis.Client
}

func NewUserRepository(db *sql.DB, rds *redis.Client) *UserRepository {
	return &UserRepository{
		db:  db,
		rds: rds,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user entity.User) error {
	query := "INSERT INTO users ( login, password, created_at, role, verification, verification_code ) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := r.db.ExecContext(ctx, query, user.Login, user.Password, user.CreatedAt, user.Role, user.Verification, user.VerificationCode)
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

//type UserSession struct {
//	UserID    int64     `redis:"user_id"`
//	SessionID uuid.UUID `redis:"session_id"`
//}

func (r *UserRepository) SetUserSession(ctx context.Context, sessionID uuid.UUID, user entity.User) error {
	us, err := json.Marshal(user)
	if err != nil {
		return err
	}

	_, err = r.rds.Set(ctx, sessionID.String(), us, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UserByLogin(ctx context.Context, login string) (entity.User, error) {
	query := "SELECT id, login, password, verification FROM users WHERE login = $1"

	var user entity.User

	err := r.db.QueryRowContext(ctx, query, login).Scan(&user.ID, &user.Login, &user.Password, &user.Verification)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, entity.ErrNotFound
		}
		return entity.User{}, err
	}
	return user, nil
}

func (r *UserRepository) Verification(ctx context.Context, verificationCode string, verification bool) error {
	query := "UPDATE users SET verification = $1 WHERE verification_code = $2 "
	_, err := r.db.ExecContext(ctx, query, verification, verificationCode)
	if err != nil {
		return err
	}
	return nil
}
