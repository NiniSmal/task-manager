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

	err = p.addProjectMembersByID(ctx, tx, project.UserID, project.ID)
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

func (p *ProjectRepository) Projects(ctx context.Context, filter entity.ProjectFilter) ([]entity.Project, error) {
	query := "SELECT id, name, created_at, updated_at, user_id FROM projects"

	var projects []entity.Project

	if filter.UserID != 0 {
		query += fmt.Sprintf(" WHERE user_id = %d", filter.UserID)
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

func (p *ProjectRepository) UpdateProject(ctx context.Context, projectID int64, project entity.Project) error {
	query := "UPDATE projects SET name = $1, updated_at = $2 WHERE id = $3"

	_, err := p.db.ExecContext(ctx, query, project.Name, project.UpdatedAt, projectID)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProjectRepository) DeleteProject(ctx context.Context, id int64) error {
	query := "DELETE FROM projects WHERE id = $1"

	_, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProjectRepository) AddProjectMembers(ctx context.Context, code string) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var userID, projectID int64

	query := "SELECT user_id, project_id FROM codes_projects_users WHERE code = $1"
	err = tx.QueryRowContext(ctx, query, code).Scan(&userID, &projectID)
	if err != nil {
		return err
	}

	err = p.addProjectMembersByID(ctx, tx, userID, projectID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (p *ProjectRepository) addProjectMembersByID(ctx context.Context, tx *sql.Tx, userID int64, projectID int64) error {
	query := "INSERT INTO user_projects (user_id, project_id) VALUES ($1, $2) ON CONFLICT (user_id, project_id) DO NOTHING"

	_, err := tx.ExecContext(ctx, query, userID, projectID)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProjectRepository) UserProjects(ctx context.Context, filter entity.ProjectFilter) ([]entity.Project, error) {
	query := "SELECT p.id, p.name, p.created_at, p.updated_at, p.user_id FROM projects p  JOIN user_projects up  ON p.id = up.project_id"

	if filter.UserID != 0 {
		query += fmt.Sprintf(" WHERE up.user_id = %d", filter.UserID)
	}

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

func (p *ProjectRepository) ProjectUsers(ctx context.Context, projectID int64) ([]entity.User, error) {
	query := "SELECT id, email,  created_at, role  FROM users  JOIN user_projects ON users.id = user_projects.user_id WHERE project_id = $1"

	rows, err := p.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User

	for rows.Next() {
		var user entity.User

		err = rows.Scan(&user.ID, &user.Email, &user.CreatedAt, &user.Role)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (p *ProjectRepository) JoiningUsers(ctx context.Context, projectID int64, userID int64, code string) error {
	query := "INSERT INTO codes_projects_users (code, project_id, user_id) VALUES($1, $2, $3)"

	_, err := p.db.ExecContext(ctx, query, code, projectID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProjectRepository) GetCodeProjectUser(ctx context.Context, projectID int64, userEmail string) (string, error) {
	query := "SELECT code FROM codes_projects_users WHERE projects_id = $1 AND user_email = $2"

	var code string
	err := p.db.QueryRowContext(ctx, query, projectID, userEmail).Scan(&code)
	if err != nil {
		return "", err
	}
	return code, nil
}
