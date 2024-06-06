package entity

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

type Role string
type View string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

const (
	minLogin = 5
	maxLogin = 200
)

type User struct {
	ID               int64      `json:"id"`
	Email            string     `json:"email"`
	Password         string     `json:"password,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	Role             Role       `json:"role"`
	Verification     bool       `json:"verification,omitempty"`
	VerificationCode string     `json:"verification_code,omitempty"`
	Photo            *string    `json:"photo"`
	DeleteAt         *time.Time `json:"deleted_at,omitempty"`
}

func (user *User) Validate() error {
	userEmail := strings.TrimSpace(user.Email)
	userPassword := strings.TrimSpace(user.Password)
	rl := utf8.RuneCountInString(userEmail)

	if rl < minLogin {
		return fmt.Errorf("%w: the login must be minimum %d symbols", ErrValidate, minLogin)
	}
	if rl > maxLogin {
		return fmt.Errorf("%w: the login can be max %d symbols", ErrValidate, maxLogin)
	}
	if userPassword == "" {
		return fmt.Errorf("%w: the password can't be empty", ErrValidate)
	}
	return nil
}
