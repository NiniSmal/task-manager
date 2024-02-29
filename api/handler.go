package api

import (
	"context"
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
	AddTask(ctx context.Context, task entity.Task) error
	GetTask(ctx context.Context, id int64) (entity.Task, error)
	GetAllTasks(ctx context.Context) ([]entity.Task, error)
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task entity.Task
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.service.AddTask(ctx, task)
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

	ctx := r.Context()

	task, err := h.service.GetTask(ctx, int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tasks, err := h.service.GetAllTasks(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
