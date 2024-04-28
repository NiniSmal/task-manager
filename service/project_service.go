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
	SaveProject(ctx context.Context, project entity.Project) error
	GetProject(ctx context.Context, id int64) (entity.Project, error)
	GetUserProjects(ctx context.Context, userID int64) ([]entity.Project, error)
	GetProjects(ctx context.Context) ([]entity.Project, error)
	UpdateProject(ctx context.Context, id int64, project entity.Project) error
	DeleteProject(ctx context.Context, id int64) error
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

	err = p.repo.SaveProject(ctx, project)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProjectService) GetProject(ctx context.Context, id int64) (entity.Project, error) {
	user := ctx.Value("user").(entity.User)

	project, err := p.repo.GetProject(ctx, id)
	if err != nil {
		return entity.Project{}, fmt.Errorf("get project by id %d :%w", id, err)
	}

	if user.Role == entity.RoleAdmin {
		return project, nil
	}
	if project.UserID == user.ID {
		return project, nil
	} else {
		return entity.Project{}, fmt.Errorf("get project by id %d :%w", id, err)
	}
}

func (p *ProjectService) GetAllProjects(ctx context.Context) ([]entity.Project, error) {
	user := ctx.Value("user").(entity.User)

	if user.Role == entity.RoleAdmin {
		projects, err := p.repo.GetProjects(ctx)
		if err != nil {
			return nil, fmt.Errorf("get all projects: %w", err)
		}
		return projects, nil
	}
	projects, err := p.repo.GetUserProjects(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("get all projects: %w", err)
	}
	return projects, nil
}

func (p *ProjectService) UpdateProject(ctx context.Context, id int64, project entity.Project) error {
	user := ctx.Value("user").(entity.User)

	projectOld, err := p.repo.GetProject(ctx, id)
	if err != nil {
		return fmt.Errorf("get project by id: %w", err)
	}

	project.UpdatedAt = time.Now()

	if user.Role == entity.RoleAdmin {
		err = p.repo.UpdateProject(ctx, id, project)
		if err != nil {
			return fmt.Errorf("update project: %w", err)
		}
	}
	if user.ID == projectOld.UserID {
		err = p.repo.UpdateProject(ctx, id, project)
		if err != nil {
			return fmt.Errorf("update project: %w", err)
		}
	}
	return nil
}

func (p *ProjectService) DeleteProject(ctx context.Context, id int64) error {
	user := ctx.Value("user").(entity.User)
	projectOld, err := p.repo.GetProject(ctx, id)
	if err != nil {
		return fmt.Errorf("get project by id: %w", err)
	}
	if user.ID == projectOld.UserID {
		err = p.repo.DeleteProject(ctx, id)
		if err != nil {
			return fmt.Errorf("delete progect by id %d: %w", id, err)
		}
	}
	return nil
}
