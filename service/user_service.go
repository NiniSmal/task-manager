package service

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/google/uuid"
	"gitlab.com/nina8884807/task-manager/entity"
)

type UserService struct {
	repo   UserRepository
	sender *SenderService
	appURL string
}

func NewUserService(r UserRepository, s *SenderService, appURL string) *UserService {
	return &UserService{
		repo:   r,
		sender: s,
		appURL: appURL,
	}
}

type UserRepository interface {
	CreateUser(ctx context.Context, user entity.User) (int64, error)
	SaveSession(ctx context.Context, sessionID uuid.UUID, user entity.User, createdAtSession time.Time) error
	UserByEmail(ctx context.Context, login string) (entity.User, error)
	Verification(ctx context.Context, verificationCode string, verification bool) (int64, error)
	UpdateVerificationCode(ctx context.Context, id int64, verificationCode string) error
	UsersToSendVIP(ctx context.Context) ([]entity.User, error)
	SaveVIPMessage(ctx context.Context, userID int64, createdAt time.Time) error
	UsersToSendAuth(ctx context.Context) ([]entity.User, error)
	SaveSendAbsenceReminder(ctx context.Context, userID int64, createdAT time.Time) error
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

func (u *UserService) CreateUser(ctx context.Context, login, password string) error {
	_, err := u.repo.UserByEmail(ctx, login)
	if err == nil {
		return entity.ErrEmailExists
	}

	if !errors.Is(err, entity.ErrNotFound) {
		return fmt.Errorf("get user by login: %w", err)
	}

	user := entity.User{
		Email:            login,
		Password:         password,
		CreatedAt:        time.Now(),
		Role:             entity.RoleUser,
		Verification:     false,
		VerificationCode: uuid.NewString(),
	}

	err = user.Validate()
	if err != nil {
		return err
	}
	user.Password, err = hashPassword(password)
	if err != nil {
		return err
	}

	_, err = u.repo.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	email := Email{
		Text:    u.appURL + "/verification?code=" + user.VerificationCode,
		To:      user.Email,
		Subject: "Account verification",
	}

	err = u.sender.SendEmail(ctx, email)
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

	err = checkPasswordHash(password, user.Password)
	if err != nil {
		return uuid.UUID{}, entity.ErrNotAuthenticated
	}

	if !user.Verification {
		return uuid.UUID{}, entity.ErrNotVerification
	}
	sessionID := uuid.New()
	createdAtSession := time.Now()
	err = u.repo.SaveSession(ctx, sessionID, user, createdAtSession)
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
		email := Email{
			Text:    "You have been assigned VIP status",
			To:      user.Email,
			Subject: "VIP Status",
		}
		err = u.sender.SendEmail(ctx, email)
		if err != nil {
			return fmt.Errorf("send email: %w", err)
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

	emailToSend := Email{
		Text:    u.appURL + "/verification?code=" + user.VerificationCode,
		To:      user.Email,
		Subject: "Account verification",
	}

	err = u.sender.SendEmail(ctx, emailToSend)
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}

// добавить функцию, кот. шлет письмо на почту, если пользователь не заходил больше чем N дней.
func (u *UserService) SendAnAbsenceLetter(ctx context.Context, intervalTime string) error {
	users, err := u.repo.UsersToSendAuth(ctx)
	if err != nil {
		return fmt.Errorf("get users: %w", err)
	}

	for _, user := range users {
		email := Email{
			Text:    fmt.Sprintf("you haven't logged into your account for %s ", intervalTime),
			To:      user.Email,
			Subject: "Absence reminder",
		}
		err = u.sender.SendEmail(ctx, email)
		if err != nil {
			return fmt.Errorf("send email: %w", err)
		}
		err = u.repo.SaveSendAbsenceReminder(ctx, user.ID, time.Now())
		if err != nil {
			return fmt.Errorf("save send message about absence remider; %w", err)
		}
	}
	return nil
}
