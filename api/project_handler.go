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
	ProjectByID(ctx context.Context, id int64) (entity.Project, error)
	Projects(ctx context.Context) ([]entity.Project, error)
	UpdateProject(ctx context.Context, id int64, project entity.Project) error
	DeleteProject(ctx context.Context, id int64) error
	AddProjectMembers(ctx context.Context, projectID int64, userID int64) error
	UserProjects(ctx context.Context) ([]entity.Project, error)
}

func (p *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var project entity.Project

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

	err = p.service.AddProject(r.Context(), project)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
}

func (p *ProjectHandler) ProjectByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idR := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idR, 10, 64)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

	project, err := p.service.ProjectByID(ctx, id)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
	err = sendJSON(w, project)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
}

func (p *ProjectHandler) Projects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	projects, err := p.service.Projects(ctx)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

	err = sendJSON(w, projects)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
}

func (p *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idR := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idR, 10, 64)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
	var project entity.Project

	err = json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

	err = p.service.UpdateProject(ctx, id, project)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
}
func (p *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idR := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idR)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
	err = p.service.DeleteProject(ctx, int64(id))
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
}

type AddProjectMemberRequest struct {
	ProjectID int64 `json:"project_id"`
	UserID    int64 `json:"user_id"`
}

func (p *ProjectHandler) AddProjectMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var data AddProjectMemberRequest

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

	err = p.service.AddProjectMembers(ctx, data.ProjectID, data.UserID)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

}

func (p *ProjectHandler) UserProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	projects, err := p.service.UserProjects(ctx)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

	err = sendJSON(w, projects)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

}
