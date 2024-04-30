package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"gitlab.com/nina8884807/task-manager/entity"
	"log"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	service TaskService //зависимость, чтобы выполнить задачу нужен метод другого типа
}

func NewTaskHandler(s TaskService) *TaskHandler {
	return &TaskHandler{
		service: s,
	}
}

type TaskService interface {
	AddTask(ctx context.Context, task entity.Task) error
	GetTask(ctx context.Context, id int64) (entity.Task, error)
	GetAllTasks(ctx context.Context, f entity.TaskFilter) ([]entity.Task, error)
	UpdateTask(ctx context.Context, id int64, task entity.UpdateTask) error
}

func HandlerError(w http.ResponseWriter, err error) {
	log.Println("API error:", err)

	switch {
	case errors.Is(err, entity.ErrNotVerification):
		http.Error(w, err.Error(), http.StatusUnauthorized)
	case errors.Is(err, entity.ErrNotAuthenticated):
		http.Error(w, err.Error(), http.StatusUnauthorized)
	case errors.Is(err, entity.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, entity.ErrEmailExists):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, entity.ErrIncorrectName):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, entity.ErrIncorrectEmail):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, entity.ErrForbidden):
		http.Error(w, err.Error(), http.StatusForbidden)
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func sendJSON(w http.ResponseWriter, body any) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		return err
	}
	return nil

}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {

	var task entity.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		HandlerError(w, err)
		return
	}

	err = h.service.AddTask(r.Context(), task) //передаем контекст, полученный из запроса.
	if err != nil {
		HandlerError(w, err)
		return
	}
}

func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	idR := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idR, 10, 64)
	if err != nil {
		HandlerError(w, err)
		return
	}

	task, err := h.service.GetTask(r.Context(), id)
	if err != nil {
		HandlerError(w, err)
		return
	}

	err = sendJSON(w, task)
	if err != nil {
		HandlerError(w, err)
		return
	}
}
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	filter := entity.TaskFilter{
		UserID:    r.URL.Query().Get("user_id"),
		ProjectID: r.URL.Query().Get("project_id"),
	}

	tasks, err := h.service.GetAllTasks(r.Context(), filter)
	if err != nil {
		HandlerError(w, err)
		return
	}
	err = sendJSON(w, tasks)
	if err != nil {
		HandlerError(w, err)
		return
	}
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idR := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idR)
	if err != nil {
		HandlerError(w, err)
		return
	}
	var task entity.UpdateTask

	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		HandlerError(w, err)
		return
	}

	err = h.service.UpdateTask(ctx, int64(id), task)
	if err != nil {
		HandlerError(w, err)
		return
	}
}
