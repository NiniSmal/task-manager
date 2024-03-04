package api

import (
	"context"
	"encoding/json"
	"gitlab.com/nina8884807/task-manager/entity"
	"net/http"
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
	AddUser(ctx context.Context, user entity.User) error
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user entity.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		HandlerError(w, err)
		return
	}

	err = u.service.AddUser(r.Context(), user)
	if err != nil {
		HandlerError(w, err)
		return
	}
}
