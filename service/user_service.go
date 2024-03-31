package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	gen "gitlab.com/nina8884807/mail/proto"
	"gitlab.com/nina8884807/task-manager/entity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	Verification(ctx context.Context, user entity.User) error
}

func (u *UserService) CreateUser(ctx context.Context, user entity.User) error {
	err := user.Validate()
	if err != nil {
		return fmt.Errorf("validation: %w", err)
	}

	_, err = u.repo.CheckUserByLogin(ctx, user.Login)
	if err == nil {
		return fmt.Errorf("this login already exists")
	}
	if !errors.Is(err, entity.ErrNotFound) {
		return fmt.Errorf("get user by login: %w", err)
	}

	user.CreatedAt = time.Now()
	user.Role = entity.RoleUser
	user.Verification = false
	code := uuid.NewString()

	user.VerificationCode = code

	err = u.repo.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	con, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	mailClient := gen.NewMailClient(con)
	_, err = mailClient.SendEmail(ctx, &gen.SendEmailRequest{
		Text: "http://localhost:8021/verification?code=" + code,
		To:   user.Login,
	})
	if err != nil {
		return fmt.Errorf("send email: %w", err)
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

func (u *UserService) Verification(ctx context.Context, user entity.User) error {
	err := u.repo.Verification(ctx, user)
	if err != nil {
		fmt.Errorf("verification: %w", err)
	}
	return nil
}
