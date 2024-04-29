package repository

import (
	"context"
	"database/sql"
	"fmt"
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
func (p *ProjectRepository) SaveProject(ctx context.Context, project entity.Project) (int64, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	query := "INSERT INTO projects (name, created_at, updated_at, user_id) VALUES ($1, $2, $3, $4) RETURNING id"

	err = tx.QueryRowContext(ctx, query, project.Name, project.CreatedAt, project.UpdatedAt, project.UserID).Scan(&project.ID)
	if err != nil {
		return 0, err
	}

	err = p.addProjectMembersByID(ctx, project.UserID, project.ID, tx)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return project.ID, nil
}

func (p *ProjectRepository) ProjectByID(ctx context.Context, id int64) (entity.Project, error) {
	query := "SELECT id, name, created_at, updated_at, user_id FROM projects WHERE id  = $1"

	var project entity.Project

	err := p.db.QueryRowContext(ctx, query, id).Scan(&project.ID, &project.Name, &project.CreatedAt, &project.UpdatedAt, &project.UserID)
	if err != nil {
		return entity.Project{}, err
	}

	return project, nil
}

func (p *ProjectRepository) Projects(ctx context.Context, filter entity.ProjectFilter) (projects []entity.Project, err error) {
	query := "SELECT id, name, created_at, updated_at, user_id FROM projects"
	if filter.UserID != 0 {
		query += fmt.Sprintf("WHERE user_id = %d", filter.UserID)
	}

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

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
	query := "UPDATE projects SET name = $1, updated_at = $2 WHERE id = $3"

	_, err := p.db.ExecContext(ctx, query, project.Name, project.UpdatedAt, id)
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

func (p *ProjectRepository) AddProjectMembersByID(ctx context.Context, userID int64, projectID int64) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = p.addProjectMembersByID(ctx, userID, projectID, tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (p *ProjectRepository) addProjectMembersByID(ctx context.Context, userID int64, projectID int64, tx *sql.Tx) error {
	query := "INSERT INTO user_projects (user_id, project_id) VALUES ($1, $2) ON CONFLICT (user_id, project_id) DO NOTHING"

	_, err := tx.ExecContext(ctx, query, userID, projectID)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProjectRepository) UserProjects(ctx context.Context, userID int64) ([]entity.Project, error) {
	query := "SELECT p.id, p.name, p.created_at, p.updated_at, p.user_id FROM projects p JOIN user_projects up ON p.user_id = up.user_id WHERE up.user_id = $1 "

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
