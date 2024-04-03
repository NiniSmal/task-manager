package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/nina8884807/task-manager/entity"
	"net/http"
	"time"
)

type UserHandler struct {
	service UserService
}

func NewUserHandler(u UserService) *UserHandler {
	return &UserHandler{
		service: u,
	}
}

type UserService interface {
	CreateUser(ctx context.Context, login, password string) error
	Login(ctx context.Context, login, password string) (uuid.UUID, error)
	Verification(ctx context.Context, verificationCode string, verification bool) error
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user entity.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		HandlerError(w, err)
		return
	}

	err = u.service.CreateUser(r.Context(), user.Login, user.Password)
	if err != nil {
		HandlerError(w, err)
		return
	}

}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user entity.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		HandlerError(w, err)
		return
	}

	sessionID, err := u.service.Login(r.Context(), user.Login, user.Password)

	if err != nil {
		HandlerError(w, err)
		return
	}

	cookie := http.Cookie{
		Name:       "session_id",
		Value:      sessionID.String(),
		Path:       "/",
		Domain:     "localhost",
		Expires:    time.Now().Add(time.Hour),
		RawExpires: "",
		MaxAge:     3600,
		Secure:     false,
		HttpOnly:   false,
		SameSite:   http.SameSiteNoneMode,
		Raw:        "",
		Unparsed:   nil,
	}
	http.SetCookie(w, &cookie)
}

func (u *UserHandler) Verification(w http.ResponseWriter, r *http.Request) {
	var user entity.User

	code := r.URL.Query().Get("code")
	if code == "" {
		HandlerError(w, errors.New("code is empty"))
		return
	}

	user.VerificationCode = code
	user.Verification = true

	err := u.service.Verification(r.Context(), user.VerificationCode, user.Verification)
	if err != nil {
		HandlerError(w, err)
	}

	_, _ = fmt.Fprintln(w, "verification successful, now you can login")
}
