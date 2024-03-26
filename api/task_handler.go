package api

import (
	"context"
	"encoding/json"
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
	GetAllTasks(ctx context.Context) ([]entity.Task, error)
	UpdateTask(ctx context.Context, id int64, task entity.UpdateTask) error
}

func HandlerError(w http.ResponseWriter, err error) {
	log.Println(err)
	w.Write([]byte("The problem is in program"))
}

func (h *TaskHandler) HandlerAnswerEncode(w http.ResponseWriter, body any) error {
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
	//idR := r.URL.Query().Get("id")
	idR := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idR)
	if err != nil {
		HandlerError(w, err)
		return
	}

	task, err := h.service.GetTask(r.Context(), int64(id))
	if err != nil {
		HandlerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = h.HandlerAnswerEncode(w, task)
	if err != nil {
		HandlerError(w, err)
		return
	}
}
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {

	tasks, err := h.service.GetAllTasks(r.Context()) //передаем контекст, полученный из запроса
	if err != nil {
		HandlerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = h.HandlerAnswerEncode(w, tasks)
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
