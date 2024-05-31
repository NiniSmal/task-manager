package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

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

func (r *UserRepository) CreateUser(ctx context.Context, user entity.User) (int64, error) {
	query := "INSERT INTO users ( email, password, created_at, role, verification, verification_code ) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	var id int64

	err := r.db.QueryRowContext(ctx, query, user.Email, user.Password, user.CreatedAt, user.Role, user.Verification,
		user.VerificationCode).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (entity.User, error) {
	query := "SELECT email, password, created_at, role, verification, verification_code FROM users WHERE id = $1"
	var user entity.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.Email, &user.Password, &user.CreatedAt, &user.Role, &user.Verification,
		&user.VerificationCode)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *UserRepository) SaveSession(ctx context.Context, sessionID uuid.UUID, user entity.User, createdAtSession time.Time) error {
	query := "INSERT INTO sessions (id, user_id, created_at) VALUES ($1, $2, $3)"

	_, err := r.db.ExecContext(ctx, query, sessionID, user.ID, createdAtSession)
	if err != nil {
		return err
	}
	err = r.saveSessionToCache(ctx, sessionID, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) saveSessionToCache(ctx context.Context, sessionID uuid.UUID, user entity.User) error {
	us, err := json.Marshal(user)
	if err != nil {
		return err
	}

	_, err = r.rds.Set(ctx, sessionID.String(), us, time.Hour).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UserByEmail(ctx context.Context, email string) (entity.User, error) {
	query := "SELECT id, email, password, created_at, role, verification, verification_code FROM users WHERE email = $1"

	var user entity.User

	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.Role,
		&user.Verification, &user.VerificationCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, entity.ErrNotFound
		}
		return entity.User{}, err
	}
	return user, nil
}

func (r *UserRepository) Verification(ctx context.Context, verificationCode string, verification bool) (int64, error) {
	query := "UPDATE users SET verification = $1 WHERE verification_code = $2 RETURNING id"
	var id int64
	err := r.db.QueryRowContext(ctx, query, verification, verificationCode).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (r *UserRepository) UpdateVerificationCode(ctx context.Context, id int64, verificationCode string) error {
	query := "UPDATE users SET verification_code =$1 WHERE id = $2"
	_, err := r.db.ExecContext(ctx, query, verificationCode, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetSession(ctx context.Context, sessionID uuid.UUID) (entity.User, error) {
	user, err := r.getSessionFromCache(ctx, sessionID)
	if err == nil {
		return user, nil
	}

	query := "SELECT u.id, u.email, u.created_at, u.verification, u.role FROM users u JOIN sessions ON u.id = sessions.user_id WHERE sessions.id = $1"

	err = r.db.QueryRowContext(ctx, query, sessionID).Scan(&user.ID, &user.Email, &user.CreatedAt, &user.Verification, &user.Role)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *UserRepository) getSessionFromCache(ctx context.Context, sessionID uuid.UUID) (entity.User, error) {
	var s string

	err := r.rds.Get(ctx, sessionID.String()).Scan(&s)
	if err != nil {
		return entity.User{}, err
	}

	var user entity.User

	err = json.Unmarshal([]byte(s), &user)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *UserRepository) UsersToSendVIP(ctx context.Context) ([]entity.User, error) {
	query := `SELECT u.id, u.email, u.created_at, u.role, u.verification, u.verification_code 
FROM users u LEFT JOIN messages m ON m.user_id = u.id 
WHERE u.created_at < NOW() - INTERVAL '1 month' AND m.created_at IS NULL AND m.message_type =$1`

	rows, err := r.db.QueryContext(ctx, query, "VIP message")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User

		err = rows.Scan(&user.ID, &user.Email, &user.CreatedAt, &user.Role, &user.Verification, &user.VerificationCode)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) SaveVIPMessage(ctx context.Context, userID int64, createdAt time.Time) error {
	query := "INSERT INTO messages(user_id, message_type, created_at) VALUES($1, $2, $3)"
	_, err := r.db.ExecContext(ctx, query, userID, "VIP message", createdAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UsersToSendAuth(ctx context.Context) ([]entity.User, error) {
	query := `SELECT u.id, u.email, u.created_at, u.role, u.verification, u.verification_code
FROM users AS u
    LEFT JOIN messages AS m ON u.id = m.user_id
    JOIN sessions AS ss ON u.id = ss.user_id
WHERE ss.created_at < NOW() - INTERVAL '1 month'
  AND m.message_type != $1 OR m.message_type IS NULL`

	rows, err := r.db.QueryContext(ctx, query, "absence message")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User

		err = rows.Scan(&user.ID, &user.Email, &user.CreatedAt, &user.Role, &user.Verification, &user.VerificationCode)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) SaveSendAbsenceReminder(ctx context.Context, userID int64, createdAt time.Time) error {
	query := "INSERT INTO messages(user_id, message_type, created_at) VALUES ($1, $2, $3)"
	_, err := r.db.ExecContext(ctx, query, userID, "absence message", createdAt)
	if err != nil {
		return err
	}
	return nil
}
