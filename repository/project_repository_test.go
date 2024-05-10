package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"gitlab.com/nina8884807/task-manager/entity"
	"testing"
	"time"
)

func ProjectConnection(t *testing.T) (*sql.DB, *redis.Client) {
	t.Helper()
	db, err := sql.Open("postgres", "postgres://postgres:dev@localhost:9000/postgres?sslmode=disable")
	require.NoError(t, err)

	t.Cleanup(func() {
		err := db.Close()
		require.NoError(t, err)
	})
	err = db.Ping()
	require.NoError(t, err)
	ctx := context.Background()
	rds := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	t.Cleanup(func() {
		err := rds.Close()
		require.NoError(t, err)
	})
	_, err = rds.Ping(ctx).Result()
	require.NoError(t, err)

	return db, rds
}

func TestProjectRepository_SaveProject(t *testing.T) {
	db, rds := ProjectConnection(t)
	pr := NewProjectRepository(db, rds)
	ur := NewUserRepository(db, rds)

	ctx := context.Background()

	user := entity.User{Email: uuid.NewString()}
	userID, err := ur.CreateUser(ctx, user)
	require.NoError(t, err)
	user.ID = userID
	createdAt := time.Now()
	updatedAt := createdAt
	project := entity.Project{
		Name:      uuid.New().String(),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		UserID:    user.ID,
		Members:   nil,
	}

	idPr, err := pr.SaveProject(ctx, project)
	require.NoError(t, err)

	dbProject, err := pr.ProjectByID(ctx, idPr)
	require.NoError(t, err)
	require.Equal(t, project.Name, dbProject.Name)
	require.Equal(t, project.CreatedAt.Unix(), dbProject.CreatedAt.Unix())
	require.Equal(t, project.UpdatedAt.Unix(), dbProject.UpdatedAt.Unix())
	require.Equal(t, project.UserID, dbProject.UserID)
	require.Equal(t, project.Members, dbProject.Members)
}

func TestProjectRepository_ByID_Error(t *testing.T) {
	db, rds := ProjectConnection(t)
	pr := NewTaskRepository(db, rds)
	ctx := context.Background()
	_, err := pr.ByID(ctx, 1234)
	require.Error(t, err)
}

// ошибка, не могу найти
func TestProjectRepository_UpdateProject(t *testing.T) {
	db, rds := ProjectConnection(t)
	pr := NewProjectRepository(db, rds)
	ur := NewUserRepository(db, rds)
	ctx := context.Background()

	user := entity.User{Email: uuid.NewString()}
	userID, err := ur.CreateUser(ctx, user)
	require.NoError(t, err)
	user.ID = userID

	createdAt := time.Now()
	updatedAt := createdAt
	project := entity.Project{
		Name:      uuid.New().String(),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		UserID:    user.ID,
	}

	idPr, err := pr.SaveProject(ctx, project)
	require.NoError(t, err)

	project2 := entity.Project{
		ID:        idPr,
		Name:      uuid.New().String(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
	}
	err = pr.UpdateProject(ctx, idPr, project2)
	require.NoError(t, err)

	dbProject, err := pr.ProjectByID(ctx, idPr)
	require.NoError(t, err)
	require.Equal(t, project2.Name, dbProject.Name)
	require.Equal(t, project2.UpdatedAt.Unix(), dbProject.UpdatedAt.Unix())
}

func TestProjectRepository_AddProjectMembersByID(t *testing.T) {
	db, rds := ProjectConnection(t)
	pr := NewProjectRepository(db, rds)
	ur := NewUserRepository(db, rds)
	ctx := context.Background()

	user := entity.User{Email: uuid.NewString()}
	userID, err := ur.CreateUser(ctx, user)
	require.NoError(t, err)
	user.ID = userID

	project := entity.Project{UserID: userID}
	projectID, err := pr.SaveProject(ctx, project)
	require.NoError(t, err)
	project.ID = projectID

	err = pr.AddProjectMembersByID(ctx, userID, projectID)
	require.NoError(t, err)

	users, err := pr.ProjectUsers(ctx, projectID)
	require.NoError(t, err)
	require.Contains(t, users, user)
}

func TestProjectRepository_UserProjects(t *testing.T) {
	db, rds := ProjectConnection(t)
	pr := NewProjectRepository(db, rds)
	ur := NewUserRepository(db, rds)
	ctx := context.Background()

	user := entity.User{Email: uuid.NewString()}
	userID, err := ur.CreateUser(ctx, user)
	require.NoError(t, err)
	user.ID = userID
	createdAt := time.Now().Round(time.Millisecond).UTC()
	projects := []entity.Project{
		{
			Name:      uuid.New().String(),
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
			UserID:    user.ID,
		}, {
			Name:      uuid.New().String(),
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
			UserID:    user.ID,
		},
	}
	for i, project := range projects {
		projectID, err := pr.SaveProject(ctx, project)
		require.NoError(t, err)
		projects[i].ID = projectID
	}
	filter := entity.ProjectFilter{UserID: user.ID}
	projectsDB, err := pr.UserProjects(ctx, filter)
	require.NoError(t, err)

	for _, project := range projects {
		require.Contains(t, projectsDB, project)
	}
}

func TestProjectRepository_Projects(t *testing.T) {
	db, rds := ProjectConnection(t)
	pr := NewProjectRepository(db, rds)
	ur := NewUserRepository(db, rds)
	ctx := context.Background()

	user := entity.User{Email: uuid.NewString()}
	userID, err := ur.CreateUser(ctx, user)
	require.NoError(t, err)
	user.ID = userID

	createdAt := time.Now().Round(time.Millisecond).UTC()
	projects := []entity.Project{
		{
			Name:      uuid.New().String(),
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
			UserID:    user.ID,
		}, {
			Name:      uuid.New().String(),
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
			UserID:    user.ID,
		},
	}
	for i, project := range projects {
		projectID, err := pr.SaveProject(ctx, project)
		require.NoError(t, err)
		projects[i].ID = projectID
	}
	filter := entity.ProjectFilter{UserID: user.ID}
	projectsDB, err := pr.Projects(ctx, filter)
	require.NoError(t, err)

	for _, project := range projects {
		require.Contains(t, projectsDB, project)
	}
}
