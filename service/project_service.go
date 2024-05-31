package service

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"gitlab.com/nina8884807/task-manager/entity"
)

type ProjectService struct {
	repo   ProjectRepository
	sender *SenderService
	appURL string
	user   UserRepository
}

func NewProjectService(r ProjectRepository, s *SenderService, appURL string, user UserRepository) *ProjectService {
	return &ProjectService{
		repo:   r,
		sender: s,
		appURL: appURL,
		user:   user,
	}
}

type ProjectRepository interface {
	SaveProject(ctx context.Context, project entity.Project) (entity.Project, error)
	ProjectByID(ctx context.Context, id int64) (entity.Project, error)
	Projects(ctx context.Context, filter entity.ProjectFilter) ([]entity.Project, error)
	AddProjectMembers(ctx context.Context, code string) error
	UpdateProject(ctx context.Context, id int64, project entity.Project) error
	DeleteProject(ctx context.Context, id int64) error
	UserProjects(ctx context.Context, filter entity.ProjectFilter) ([]entity.Project, error)
	ProjectUsers(ctx context.Context, projectID int64) ([]entity.User, error)
	JoiningUsers(ctx context.Context, projectID int64, userID int64, code string) error
}

func (p *ProjectService) AddProject(ctx context.Context, project entity.Project) (entity.Project, error) {
	err := project.Validate()
	if err != nil {
		return entity.Project{}, err
	}
	project.CreatedAt = time.Now()
	project.UpdatedAt = project.CreatedAt

	user := ctx.Value("user").(entity.User)
	project.UserID = user.ID

	projectDB, err := p.repo.SaveProject(ctx, project)
	if err != nil {
		return entity.Project{}, err
	}

	return projectDB, nil
}

func (p *ProjectService) ProjectByID(ctx context.Context, projectID int64) (entity.Project, error) {
	project, err := p.repo.ProjectByID(ctx, projectID)
	if err != nil {
		return entity.Project{}, fmt.Errorf("get project by id %d :%w", projectID, err)
	}

	projectUsers, err := p.repo.ProjectUsers(ctx, projectID)
	if err != nil {
		return entity.Project{}, fmt.Errorf("get project %d members: %w", projectID, err)
	}

	project.Members = projectUsers

	user := ctx.Value("user").(entity.User)

	if user.Role != entity.RoleAdmin && !slices.ContainsFunc(projectUsers, func(u entity.User) bool {
		return u.ID == user.ID
	}) {
		return entity.Project{}, fmt.Errorf("get project by id %d :%w", projectID, entity.ErrForbidden)
	}

	return project, nil
}

func (p *ProjectService) Projects(ctx context.Context) ([]entity.Project, error) {
	user := ctx.Value("user").(entity.User)
	var filter entity.ProjectFilter

	if user.Role != entity.RoleAdmin {
		filter.UserID = user.ID
	}

	projects, err := p.repo.Projects(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("get all projects: %w", err)
	}

	return projects, nil
}

func (p *ProjectService) AddProjectMembers(ctx context.Context, code string) error {
	err := p.repo.AddProjectMembers(ctx, code)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProjectService) JoiningUsers(ctx context.Context, projectID int64, userEmail string) error {
	user := ctx.Value("user").(entity.User)
	project, err := p.repo.ProjectByID(ctx, projectID)
	if err != nil {
		return err
	}

	if project.UserID != user.ID {
		return fmt.Errorf("add project member: %w", entity.ErrForbidden)
	}
	code := uuid.NewString()

	userToInvite, err := p.user.UserByEmail(ctx, userEmail)
	if err != nil {
		return fmt.Errorf("get user by email: %w", err)
	}

	err = p.repo.JoiningUsers(ctx, projectID, userToInvite.ID, code)
	if err != nil {
		return fmt.Errorf("joining user: %w", err)
	}

	email := Email{
		Text:    p.appURL + "/projects/joining?code=" + code,
		To:      userEmail,
		Subject: "Joining the project",
	}
	err = p.sender.SendEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	return nil
}

func (p *ProjectService) UpdateProject(ctx context.Context, projectID int64, project entity.Project) (entity.Project, error) {
	user := ctx.Value("user").(entity.User)

	projectOld, err := p.repo.ProjectByID(ctx, projectID)
	if err != nil {
		return entity.Project{}, fmt.Errorf("get project by projectID: %w", err)
	}

	project.UpdatedAt = time.Now()

	if user.Role != entity.RoleAdmin && user.ID != projectOld.UserID {
		return entity.Project{}, fmt.Errorf("update project: %w", entity.ErrForbidden)
	}

	err = p.repo.UpdateProject(ctx, projectID, project)
	if err != nil {
		return entity.Project{}, fmt.Errorf("update project: %w", err)
	}
	project, err = p.repo.ProjectByID(ctx, projectID)
	if err != nil {
		return entity.Project{}, fmt.Errorf("get project by projectID: %w", err)
	}
	return project, nil
}

func (p *ProjectService) DeleteProject(ctx context.Context, id int64) error {
	user := ctx.Value("user").(entity.User)

	projectOld, err := p.repo.ProjectByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get project by id: %w", err)
	}

	if user.ID != projectOld.UserID && user.Role != entity.RoleAdmin {
		return fmt.Errorf("delete project by id %d: %w", id, entity.ErrForbidden)
	}

	err = p.repo.DeleteProject(ctx, id)
	if err != nil {
		return fmt.Errorf("delete project by id %d: %w", id, err)
	}

	return nil
}

func (p *ProjectService) UserProjects(ctx context.Context) ([]entity.Project, error) {
	user := ctx.Value("user").(entity.User)

	var filter entity.ProjectFilter
	if user.Role != entity.RoleAdmin {
		filter.UserID = user.ID
	}
	projects, err := p.repo.UserProjects(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("get user projects: %w", err)
	}
	return projects, nil
}
