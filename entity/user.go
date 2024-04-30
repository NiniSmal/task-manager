package entity

import (
	"fmt"
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
	ID               int64     `json:"id"`
	Email            string    `json:"email"`
	Password         string    `json:"password,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	Role             Role      `json:"role"`
	Verification     bool      `json:"verification,omitempty"`
	VerificationCode string    `json:"verification_code,omitempty"`
}

func (user *User) Validate() error {
	rl := utf8.RuneCountInString(user.Email)
	if rl < minLogin {
		return fmt.Errorf("%w: the login must be minimum %d symbols", ErrValidate, minLogin)
	}
	if rl > maxLogin {
		return fmt.Errorf("%w: the login can be max %d symbols", ErrValidate, maxLogin)
	}
	if user.Password == "" {
		return fmt.Errorf("%w: the password can't be empty", ErrValidate)
	}
	return nil
}
