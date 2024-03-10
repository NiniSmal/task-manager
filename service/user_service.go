package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/nina8884807/task-manager/entity"
	"time"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(r UserRepository) *UserService {
	return &UserService{
		repo: r,
	}
}

type UserRepository interface {
	CreateUser(ctx context.Context, user entity.User) error
	GetUserIDByLoginAndPassword(ctx context.Context, user entity.User) (int64, string, error)
	SaveSession(ctx context.Context, sessionID uuid.UUID, userID int64, role string) error
}

func (u *UserService) CreateUser(ctx context.Context, user entity.User) error {
	user.CreatedAt = time.Now()
	user.Role = entity.RoleUser
	err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("create user %w", err)
	}
	return nil
}

func (u *UserService) Login(ctx context.Context, user entity.User) (uuid.UUID, error) {
	userID, role, err := u.repo.GetUserIDByLoginAndPassword(ctx, user)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("login: %w", err)
	}

	sessionID := uuid.New()

	err = u.repo.SaveSession(ctx, sessionID, userID, role)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("save session: %w", err)
	}
	return sessionID, nil
}
