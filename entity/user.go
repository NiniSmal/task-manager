package entity

import (
	"errors"
	"fmt"
	"time"
	"unicode/utf8"
)

type Role string

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
	Password         string    `json:"password"`
	CreatedAt        time.Time `json:"created_at"`
	Role             Role      `json:"role"`
	Verification     bool      `json:"verification"`
	VerificationCode string    `json:"verification_code"`
}

func (user *User) Validate() error {
	rl := utf8.RuneCountInString(user.Email)
	if rl < minLogin {
		return fmt.Errorf("the login must be minimum %d symbols", minLogin)
	}
	if rl > maxLogin {
		return fmt.Errorf("the login can be max %d symbols", maxLogin)
	}
	if user.Password == "" {
		return errors.New("the password can't be empty")
	}
	return nil
}
