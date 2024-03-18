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
	SaveSession(ctx context.Context, sessionID uuid.UUID, user entity.User) error
	CheckUserByLogin(ctx context.Context, login string) (entity.User, error)
}

func (u *UserService) CreateUser(ctx context.Context, user entity.User) error {
	err := user.Validate()
	if err != nil {
		return fmt.Errorf("validation: %w", err)
	}

	userDB, err := u.repo.CheckUserByLogin(ctx, user.Login)
	if err != nil {
		return fmt.Errorf("check user by login: %w", err)
	}

	if userDB.Login == user.Login {
		return fmt.Errorf("this login already exists")
	}

	user.CreatedAt = time.Now()
	user.Role = entity.RoleUser
	err = u.repo.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("create user %w", err)
	}
	return nil
}

func (u *UserService) Login(ctx context.Context, user entity.User) (uuid.UUID, error) {
	sessionID := uuid.New()

	err := u.repo.SaveSession(ctx, sessionID, user)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("save session: %w", err)
	}
	return sessionID, nil
}
