package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gitlab.com/nina8884807/task-manager/entity"
)

type UserHandler struct {
	service UserService
	appHost string
}

func NewUserHandler(u UserService, appHost string) *UserHandler {
	return &UserHandler{
		service: u,
		appHost: appHost,
	}
}

type UserService interface {
	CreateUser(ctx context.Context, login, password string) error
	Login(ctx context.Context, login, password string) (uuid.UUID, error)
	Verification(ctx context.Context, verificationCode string, verification bool) error
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user entity.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

	err = u.service.CreateUser(ctx, user.Email, user.Password)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user entity.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

	sessionID, err := u.service.Login(ctx, user.Email, user.Password)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

	l := r.Context().Value("logger").(*slog.Logger)

	l.Info(fmt.Sprintf("login OK session_id: %s", sessionID))

	cookie := http.Cookie{
		Name:  "session_id",
		Value: sessionID.String(),
		Path:  "/",
		// Domain:     u.appHost,
		Expires:    time.Now().Add(time.Hour * 24 * 7),
		RawExpires: "",
		MaxAge:     3600,
		Secure:     true,
		HttpOnly:   true,
		SameSite:   http.SameSiteNoneMode,
		Raw:        "",
		Unparsed:   nil,
	}
	http.SetCookie(w, &cookie)
}
func (u *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID := uuid.UUID{}

	cookie := http.Cookie{
		Name:    "session_id",
		Value:   sessionID.String(),
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Now().Add(-time.Hour),
	}
	http.SetCookie(w, &cookie)

}

func (u *UserHandler) Verification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user entity.User

	code := r.URL.Query().Get("code")
	if code == "" {
		HandlerError(ctx, w, errors.New("code is empty"))
		return
	}

	user.VerificationCode = code
	user.Verification = true

	err := u.service.Verification(ctx, user.VerificationCode, user.Verification)
	if err != nil {
		HandlerError(ctx, w, err)
	}

	_, _ = fmt.Fprintln(w, "verification successful, now you can login")
}
