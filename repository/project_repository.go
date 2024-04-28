package repository

import (
	"context"
	"database/sql"
	"github.com/redis/go-redis/v9"
	"gitlab.com/nina8884807/task-manager/entity"
)

type ProjectRepository struct {
	db  *sql.DB
	rds *redis.Client
}

func NewProjectRepository(db *sql.DB, rds *redis.Client) *ProjectRepository {
	return &ProjectRepository{
		db:  db,
		rds: rds,
	}
}
func (p *ProjectRepository) SaveProject(ctx context.Context, project entity.Project) error {
	query := "INSERT INTO projects (name, created_at, updated_at, user_id) VALUES ($1, $2, $3, $4)"

	_, err := p.db.ExecContext(ctx, query, project.Name, project.CreatedAt, project.UpdatedAt, project.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProjectRepository) GetProject(ctx context.Context, id int64) (entity.Project, error) {
	query := "SELECT id, name, created_at, updated_at, user_id FROM projects WHERE id  = $1"

	var project entity.Project

	err := p.db.QueryRowContext(ctx, query, id).Scan(&project.ID, &project.Name, &project.CreatedAt, &project.UpdatedAt, &project.UserID)
	if err != nil {
		return entity.Project{}, err
	}
	return project, nil
}

func (p *ProjectRepository) GetUserProjects(ctx context.Context, userID int64) ([]entity.Project, error) {
	query := "SELECT id, name,  created_at, updated_at, user_id FROM projects WHERE user_id = $1"

	rows, err := p.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []entity.Project
	for rows.Next() {
		var project entity.Project
		err = rows.Scan(&project.ID, &project.Name, &project.CreatedAt, &project.UpdatedAt, &project.UserID)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func (p *ProjectRepository) GetProjects(ctx context.Context) ([]entity.Project, error) {
	query := "SELECT id, name, created_at, updated_at, user_id FROM projects"
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []entity.Project
	for rows.Next() {
		var project entity.Project
		err = rows.Scan(&project.ID, &project.Name, &project.CreatedAt, &project.UpdatedAt, &project.UserID)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (p *ProjectRepository) UpdateProject(ctx context.Context, id int64, project entity.Project) error {
	query := "UPDATE projects SET name = $1 WHERE id = $2"

	_, err := p.db.ExecContext(ctx, query, project.Name, id)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProjectRepository) DeleteProject(ctx context.Context, id int64) error {
	query := "DELETE FROM projects WHERE id = $1 "

	_, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProjectRepository) AddProjectMembers() {

}
