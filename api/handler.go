package api

import (
	"encoding/json"
	"gitlab.com/nina8884807/task-manager/entity"
	"net/http"
	"strconv"
)

type Handler struct {
	service TaskService //зависимость, чтобы выполнить задачу нужен метод другого типа
}

func NewHandler(s TaskService) *Handler {
	return &Handler{
		service: s,
	}
}

type TaskService interface {
	AddTask(task entity.Task) error
	GetTaskByID(id int64) (entity.Task, error)
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task entity.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.service.AddTask(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	idR := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idR)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	task, err := h.service.GetTaskByID(int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
