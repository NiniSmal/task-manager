package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gitlab.com/nina8884807/task-manager/entity"
)

type TaskHandler struct {
	service TaskService // зависимость, чтобы выполнить задачу нужен метод другого типа
}

func NewTaskHandler(s TaskService) *TaskHandler {
	return &TaskHandler{
		service: s,
	}
}

type TaskService interface {
	AddTask(ctx context.Context, task entity.Task) (entity.Task, error)
	GetTask(ctx context.Context, id int64) (entity.Task, error)
	GetAllTasks(ctx context.Context, f entity.TaskFilter) ([]entity.Task, error)
	UpdateTask(ctx context.Context, id int64, task entity.UpdateTask) (entity.Task, error)
	Delete(ctx context.Context, id int64) error
}
type apiError struct {
	Error string `json:"error"`
}

func HandlerError(ctx context.Context, w http.ResponseWriter, err error) {
	l := ctx.Value("logger").(*slog.Logger)
	l.Error("Api response", "error", err.Error())

	errText := http.StatusText(http.StatusInternalServerError)
	errCode := http.StatusInternalServerError

	switch {
	case errors.Is(err, entity.ErrNotVerification):
		errText, errCode = err.Error(), http.StatusUnauthorized
	case errors.Is(err, entity.ErrNotAuthenticated):
		errText, errCode = err.Error(), http.StatusUnauthorized
	case errors.Is(err, entity.ErrNotFound):
		errText, errCode = err.Error(), http.StatusNotFound
	case errors.Is(err, entity.ErrEmailExists):
		errText, errCode = err.Error(), http.StatusConflict
	case errors.Is(err, entity.ErrValidate):
		errText, errCode = err.Error(), http.StatusBadRequest
	case errors.Is(err, entity.ErrForbidden):
		errText, errCode = err.Error(), http.StatusForbidden
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errCode)

	err = json.NewEncoder(w).Encode(apiError{Error: errText})
	if err != nil {
		l.Error("Encode", "error", err)
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
	ctx := r.Context()

	var task entity.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

	taskDB, err := h.service.AddTask(r.Context(), task)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
	err = sendJSON(w, taskDB)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
}

func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idR := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idR, 10, 64)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

	task, err := h.service.GetTask(ctx, id)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

	err = sendJSON(w, task)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
}
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter := entity.TaskFilter{
		ProjectID: r.URL.Query().Get("project_id"),
	}

	tasks, err := h.service.GetAllTasks(ctx, filter)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
	if len(tasks) == 0 {
		tasks = make([]entity.Task, 0)
	}
	err = sendJSON(w, tasks)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idR := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idR, 10, 64)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
	var task entity.UpdateTask

	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

	taskUp, err := h.service.UpdateTask(ctx, id, task)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
	err = sendJSON(w, taskUp)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idR := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idR, 10, 64)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
	err = h.service.Delete(ctx, id)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
}
