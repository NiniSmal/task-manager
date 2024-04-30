package service

import (
	"context"
	"fmt"
	"gitlab.com/nina8884807/task-manager/entity"
	"time"
)

type ProjectService struct {
	repo ProjectRepository
}

func NewProjectService(r ProjectRepository) *ProjectService {
	return &ProjectService{
		repo: r}
}

type ProjectRepository interface {
	SaveProject(ctx context.Context, project entity.Project) (int64, error)
	ProjectByID(ctx context.Context, id int64) (entity.Project, error)
	Projects(ctx context.Context, filter entity.ProjectFilter) ([]entity.Project, error)
	AddProjectMembersByID(ctx context.Context, userID int64, projectID int64) error
	UpdateProject(ctx context.Context, id int64, project entity.Project) error
	DeleteProject(ctx context.Context, id int64) error
	UserProjects(ctx context.Context, filter entity.ProjectFilter) ([]entity.Project, error)
}

func (p *ProjectService) AddProject(ctx context.Context, project entity.Project) error {
	err := project.Validate()
	if err != nil {
		return fmt.Errorf("validation: %w", entity.ErrIncorrectName)
	}
	project.CreatedAt = time.Now()
	project.UpdatedAt = project.CreatedAt

	user := ctx.Value("user").(entity.User)
	project.UserID = user.ID

	_, err = p.repo.SaveProject(ctx, project)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProjectService) ProjectByID(ctx context.Context, id int64) (entity.Project, error) {
	project, err := p.repo.ProjectByID(ctx, id)
	if err != nil {
		return entity.Project{}, fmt.Errorf("get project by id %d :%w", id, err)
	}

	user := ctx.Value("user").(entity.User)
	if user.Role != entity.RoleAdmin && project.UserID != user.ID {
		return entity.Project{}, fmt.Errorf("get project by id %d :%w", id, entity.ErrForbidden)
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

func (p *ProjectService) AddProjectMembers(ctx context.Context, projectID int64, userID int64) error {
	user := ctx.Value("user").(entity.User)

	project, err := p.repo.ProjectByID(ctx, projectID)
	if err != nil {
		return err
	}

	if project.UserID != user.ID {
		return fmt.Errorf("add project member: %w", entity.ErrForbidden)
	}

	err = p.repo.AddProjectMembersByID(ctx, userID, projectID)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProjectService) UpdateProject(ctx context.Context, id int64, project entity.Project) error {
	user := ctx.Value("user").(entity.User)

	projectOld, err := p.repo.ProjectByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get project by id: %w", err)
	}

	project.UpdatedAt = time.Now()

	if user.Role != entity.RoleAdmin && user.ID != projectOld.UserID {
		return fmt.Errorf("update project: %w", entity.ErrForbidden)
	}
	err = p.repo.UpdateProject(ctx, id, project)
	if err != nil {
		return fmt.Errorf("update project: %w", err)
	}

	return nil
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
