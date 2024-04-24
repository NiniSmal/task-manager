package api

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"gitlab.com/nina8884807/task-manager/entity"
	"net/http"
	"strconv"
)

type ProjectHandler struct {
	service ProjectService
}

func NewProjectHandler(s ProjectService) *ProjectHandler {
	return &ProjectHandler{
		service: s,
	}
}

type ProjectService interface {
	AddProject(ctx context.Context, project entity.Project) error
	GetProject(ctx context.Context, id int64) (entity.Project, error)
	GetAllProjects(ctx context.Context) ([]entity.Project, error)
	UpdateProject(ctx context.Context, id int64, project entity.Project) error
}

func (p *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var project entity.Project

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		HandlerError(w, err)
		return
	}

	err = p.service.AddProject(r.Context(), project)
	if err != nil {
		HandlerError(w, err)
		return
	}
}

func (p *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	idR := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idR, 10, 64)
	if err != nil {
		HandlerError(w, err)
		return
	}

	project, err := p.service.GetProject(r.Context(), id)
	if err != nil {
		HandlerError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(project)
	if err != nil {
		HandlerError(w, err)
		return
	}
}

func (p *ProjectHandler) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := p.service.GetAllProjects(r.Context()) //передаем контекст, полученный из запроса
	if err != nil {
		HandlerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(projects)
	if err != nil {
		HandlerError(w, err)
		return
	}
}

func (p *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idR := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idR)
	if err != nil {
		HandlerError(w, err)
		return
	}
	var project entity.Project

	err = json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		HandlerError(w, err)
		return
	}

	err = p.service.UpdateProject(ctx, int64(id), project)
	if err != nil {
		HandlerError(w, err)
		return
	}
}
