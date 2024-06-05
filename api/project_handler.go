package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gitlab.com/nina8884807/task-manager/entity"
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
	AddProject(ctx context.Context, project entity.Project) (entity.Project, error)
	ProjectByID(ctx context.Context, id int64) (entity.Project, error)
	Projects(ctx context.Context) ([]entity.Project, error)
	UpdateProject(ctx context.Context, id int64, project entity.Project) (entity.Project, error)
	SoftDeleteProject(ctx context.Context, id int64) error
	HardDeleteProjects(ctx context.Context) error
	AddProjectMembers(ctx context.Context, code string) error
	UserProjects(ctx context.Context) ([]entity.Project, error)
	JoiningUsers(ctx context.Context, projectID int64, userEmail string) error
}

func (p *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var project entity.Project

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

	projectDB, err := p.service.AddProject(r.Context(), project)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
	err = sendJSON(w, projectDB)
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

	if len(projects) == 0 {
		projects = make([]entity.Project, 0)
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

	projectDB, err := p.service.UpdateProject(ctx, id, project)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
	err = sendJSON(w, projectDB)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
}
func (p *ProjectHandler) SoftDeleteProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idR := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idR, 10, 64)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
	err = p.service.SoftDeleteProject(ctx, id)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
}

func (p *ProjectHandler) HardDeleteProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := p.service.HardDeleteProjects(ctx)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
}

type AddProjectMemberRequest struct {
	ProjectID int64  `json:"project_id"`
	UserEmail string `json:"user_email"`
}

func (p *ProjectHandler) AddProjectMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := r.URL.Query().Get("code")
	if code == "" {
		HandlerError(ctx, w, errors.New("code is empty"))
		return
	}

	err := p.service.AddProjectMembers(ctx, code)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

	http.Redirect(w, r, "https://tm.anaxita.ru/projects", http.StatusSeeOther)
}

func (p *ProjectHandler) UserProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	projects, err := p.service.UserProjects(ctx)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}

	if len(projects) == 0 {
		projects = make([]entity.Project, 0)
	}

	err = sendJSON(w, projects)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
}

func (p *ProjectHandler) JoiningUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var data AddProjectMemberRequest

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		HandlerError(ctx, w, err)
		return
	}
	err = p.service.JoiningUsers(ctx, data.ProjectID, data.UserEmail)
	if err != nil {
		HandlerError(ctx, w, err)
	}

}
