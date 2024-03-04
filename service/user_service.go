package service

import (
	"context"
	"fmt"
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
}

func (u *UserService) AddUser(ctx context.Context, user entity.User) error {
	user.CreatedAt = time.Now()

	err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("add user %w", err)
	}
	return nil
}
