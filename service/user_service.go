package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	gen "gitlab.com/nina8884807/mail/proto"
	"gitlab.com/nina8884807/task-manager/entity"
	"time"
)

type UserService struct {
	repo   UserRepository
	client gen.MailClient
}

func NewUserService(r UserRepository, mc gen.MailClient) *UserService {
	return &UserService{
		repo:   r,
		client: mc,
	}
}

type UserRepository interface {
	CreateUser(ctx context.Context, user entity.User) error
	SaveSession(ctx context.Context, sessionID uuid.UUID, user entity.User) error
	UserByLogin(ctx context.Context, login string) (entity.User, error)
	Verification(ctx context.Context, verificationCode string, verification bool) error
	SetUserSession(ctx context.Context, sessionID uuid.UUID, user entity.User) error
}

func (u *UserService) CreateUser(ctx context.Context, login, password string) error {
	user := entity.User{
		Login:            login,
		Password:         password,
		CreatedAt:        time.Now(),
		Role:             entity.RoleUser,
		Verification:     false,
		VerificationCode: uuid.NewString(),
	}

	err := user.Validate()
	if err != nil {
		return fmt.Errorf("validation: %w", err)
	}

	_, err = u.repo.UserByLogin(ctx, login)
	if err == nil {
		return fmt.Errorf("this login already exists")
	}
	if !errors.Is(err, entity.ErrNotFound) {
		return fmt.Errorf("get user by login: %w", err)
	}

	err = u.repo.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	_, err = u.client.SendEmail(ctx, &gen.SendEmailRequest{
		Text: "http://localhost:8021/verification?code=" + user.VerificationCode,
		To:   user.Login,
	})
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}

func (u *UserService) Login(ctx context.Context, login, password string) (uuid.UUID, error) {
	user, err := u.repo.UserByLogin(ctx, login)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("get user by login: %w", err)
	}
	if user.Password != password {
		return uuid.UUID{}, entity.ErrNotAuthenticated
	}
	if !user.Verification {
		return uuid.UUID{}, entity.ErrNotVerification
	}
	sessionID := uuid.New()

	err = u.repo.SaveSession(ctx, sessionID, user)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("save session: %w", err)
	}

	err = u.repo.SetUserSession(ctx, sessionID, user)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("set session %w", err)
	}
	return sessionID, nil
}

func (u *UserService) Verification(ctx context.Context, verificationCode string, verification bool) error {
	err := u.repo.Verification(ctx, verificationCode, verification)
	if err != nil {
		return fmt.Errorf("%w, follow the link in the email to verify", entity.ErrNotVerification)
	}
	return nil
}
