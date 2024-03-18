package entity

import (
	"errors"
	"time"
	"unicode/utf8"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID        int64     `json:"id"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	Role      Role      `json:"role"`
}

func (user *User) Validate() error {
	const (
		minLogin = 1
		maxLogin = 15
	)

	rl := utf8.RuneCountInString(user.Login)
	if rl < minLogin || rl > maxLogin {
		return errors.New("the login must be minimum 1 symbol and not more 15 symbols")
	}

	if user.Password != "" {
		return errors.New("the password must be more than 0")
	}
	return nil
}
