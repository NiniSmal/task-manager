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
		return fmt.Errorf("the login must be minimum %d symbol and not more %d symbols", minLogin, maxLogin)
	}

	if user.Password == "" {
		return errors.New("the password can't be empty")
	}
	return nil
}
