package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"gitlab.com/nina8884807/task-manager/entity"
)

type UserService struct {
	repo   UserRepository
	kafka  *kafka.Writer
	appURL string
}

func NewUserService(r UserRepository, w *kafka.Writer, appURL string) *UserService {
	return &UserService{
		repo:   r,
		kafka:  w,
		appURL: appURL,
	}
}

type UserRepository interface {
	CreateUser(ctx context.Context, user entity.User) (int64, error)
	SaveSession(ctx context.Context, sessionID uuid.UUID, user entity.User) error
	UserByEmail(ctx context.Context, login string) (entity.User, error)
	Verification(ctx context.Context, verificationCode string, verification bool) (int64, error)
	UpdateVerificationCode(ctx context.Context, id int64, verificationCode string) error
	UsersToSendVIP(ctx context.Context) ([]entity.User, error)
	SaveVIPMessage(ctx context.Context, userID int64, createdAt time.Time) error
}

type SendEmail struct {
	Text    string `json:"text"`
	To      string `json:"to"`
	Subject string `json:"subject"`
}

func (u *UserService) sendEmail(ctx context.Context, user entity.User) error {
	email := SendEmail{
		Text:    u.appURL + "/verification?code=" + user.VerificationCode,
		To:      user.Email,
		Subject: "Account verification",
	}

	msg, err := json.Marshal(&email)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	err = u.kafka.WriteMessages(ctx, kafka.Message{Value: msg})
	if err != nil {
		return fmt.Errorf("write messages: %w", err)
	}

	return nil
}

func (u *UserService) CreateUser(ctx context.Context, login, password string) error {
	user := entity.User{
		Email:            login,
		Password:         password,
		CreatedAt:        time.Now(),
		Role:             entity.RoleUser,
		Verification:     false,
		VerificationCode: uuid.NewString(),
	}

	err := user.Validate()
	if err != nil {
		return err
	}

	_, err = u.repo.UserByEmail(ctx, login)
	if err == nil {
		return entity.ErrEmailExists
	}

	if !errors.Is(err, entity.ErrNotFound) {
		return fmt.Errorf("get user by login: %w", err)
	}

	_, err = u.repo.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	err = u.sendEmail(ctx, user)
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}

func (u *UserService) Login(ctx context.Context, login, password string) (uuid.UUID, error) {
	user, err := u.repo.UserByEmail(ctx, login)
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

	return sessionID, nil
}

func (u *UserService) Verification(ctx context.Context, verificationCode string, verification bool) error {
	_, err := u.repo.Verification(ctx, verificationCode, verification)
	if err != nil {
		return fmt.Errorf("%w, follow the link in the email to verify", entity.ErrNotVerification)
	}
	return nil
}

func (u *UserService) SendVIPStatus(ctx context.Context, intervalTime string) error {
	users, err := u.repo.UsersToSendVIP(ctx)
	if err != nil {
		return fmt.Errorf("get users for VIP status: %w", err)
	}

	for _, user := range users {
		email := SendEmail{
			Text:    "You have been assigned VIP status",
			To:      user.Email,
			Subject: "VIP Status",
		}
		msg2, err := json.Marshal(&email)
		if err != nil {
			return fmt.Errorf("marshal message: %w", err)
		}

		err = u.kafka.WriteMessages(ctx, kafka.Message{Value: msg2})
		if err != nil {
			return fmt.Errorf("write messages: %w", err)
		}
		err = u.repo.SaveVIPMessage(ctx, user.ID, time.Now())
		if err != nil {
			return fmt.Errorf("save VIP message: %w", err)
		}
	}
	return nil
}

func (u *UserService) ResendVerificationCode(ctx context.Context, email string) error {
	user, err := u.repo.UserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("get user by email: %w", err)
	}

	user.VerificationCode = uuid.NewString()
	err = u.repo.UpdateVerificationCode(ctx, user.ID, user.VerificationCode)
	if err != nil {
		return fmt.Errorf("update verification code: %w", err)
	}

	err = u.sendEmail(ctx, user)
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}
