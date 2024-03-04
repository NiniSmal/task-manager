package api

import (
	"context"
	"encoding/json"
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
	CreateUser(ctx context.Context, user entity.User) error
	Login(ctx context.Context, user entity.User) (uuid.UUID, error)
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user entity.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		HandlerError(w, err)
		return
	}

	err = u.service.CreateUser(r.Context(), user)
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
	sessionID, err := u.service.Login(r.Context(), user)
	if err != nil {
		HandlerError(w, err)
		return
	}
	cookie := http.Cookie{
		Name:    "session_id",
		Value:   sessionID.String(),
		Path:    "/",
		Expires: time.Now().Add(time.Hour),
		MaxAge:  3600,
	}
	http.SetCookie(w, &cookie)
}
