package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"gitlab.com/nina8884807/task-manager/entity"
	"os"
	"testing"
	"time"
)

func DBConnection(t *testing.T) (*sql.DB, *redis.Client) {
	t.Helper()
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:dev@localhost:9000/postgres?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := db.Close()
		require.NoError(t, err)
	})
	err = db.Ping()
	require.NoError(t, err)

	goose.SetLogger(goose.NopLogger())
	err = goose.Up(db, "../migrations")
	require.NoError(t, err)

	ctx := context.Background()

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rds := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
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
	db, rds := DBConnection(t)
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

	projectDB, err := pr.SaveProject(ctx, project)
	require.NoError(t, err)

	dbProject, err := pr.ProjectByID(ctx, projectDB.ID)
	require.NoError(t, err)
	require.Equal(t, project.Name, dbProject.Name)
	require.Equal(t, project.CreatedAt.Unix(), dbProject.CreatedAt.Unix())
	require.Equal(t, project.UpdatedAt.Unix(), dbProject.UpdatedAt.Unix())
	require.Equal(t, project.UserID, dbProject.UserID)
	require.Equal(t, project.Members, dbProject.Members)
}

func TestProjectRepository_ByID_Error(t *testing.T) {
	db, rds := DBConnection(t)
	pr := NewTaskRepository(db, rds)
	ctx := context.Background()
	_, err := pr.ByID(ctx, 1234)
	require.Error(t, err)
}

// ошибка, не могу найти
func TestProjectRepository_UpdateProject(t *testing.T) {
	db, rds := DBConnection(t)
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

	projectDB, err := pr.SaveProject(ctx, project)
	require.NoError(t, err)

	project2 := entity.Project{
		ID:        projectDB.ID,
		Name:      uuid.New().String(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
	}
	err = pr.UpdateProject(ctx, projectDB.ID, project2)
	require.NoError(t, err)

	dbProject, err := pr.ProjectByID(ctx, projectDB.ID)
	require.NoError(t, err)
	require.Equal(t, project2.Name, dbProject.Name)
	require.Equal(t, project2.UpdatedAt.Unix(), dbProject.UpdatedAt.Unix())
}

func TestProjectRepository_AddProjectMembersByID(t *testing.T) {
	db, rds := DBConnection(t)
	pr := NewProjectRepository(db, rds)
	ur := NewUserRepository(db, rds)
	ctx := context.Background()

	user := entity.User{Email: uuid.NewString()}
	userID, err := ur.CreateUser(ctx, user)
	require.NoError(t, err)
	user.ID = userID

	project := entity.Project{UserID: userID}
	projectDB, err := pr.SaveProject(ctx, project)
	require.NoError(t, err)
	project.ID = projectDB.ID

	code := uuid.NewString()

	err = pr.JoiningUsers(ctx, project.ID, userID, code)
	require.NoError(t, err)

	err = pr.AddProjectMembers(ctx, code)
	require.NoError(t, err)

	users, err := pr.ProjectUsers(ctx, project.ID)
	require.NoError(t, err)
	require.Contains(t, users, user)
}

func TestProjectRepository_UserProjects(t *testing.T) {
	db, rds := DBConnection(t)
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
		project, err := pr.SaveProject(ctx, project)
		require.NoError(t, err)
		projects[i].ID = project.ID
	}
	filter := entity.ProjectFilter{UserID: user.ID}
	projectsDB, err := pr.UserProjects(ctx, filter)
	require.NoError(t, err)

	for _, project := range projects {
		require.Contains(t, projectsDB, project)
	}
}

func TestProjectRepository_Projects(t *testing.T) {
	db, rds := DBConnection(t)
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
		projectDB, err := pr.SaveProject(ctx, project)
		require.NoError(t, err)
		projects[i].ID = projectDB.ID
	}
	filter := entity.ProjectFilter{UserID: user.ID}
	projectsDB, err := pr.Projects(ctx, filter)
	require.NoError(t, err)

	for _, project := range projects {
		require.Contains(t, projectsDB, project)
	}
}

func TestProjectRepository_SoftDeleteProject(t *testing.T) {
	db, rds := DBConnection(t)
	ur := NewUserRepository(db, rds)
	pr := NewProjectRepository(db, rds)
	ctx := context.Background()

	user := entity.User{Email: uuid.NewString()}
	userID, err := ur.CreateUser(ctx, user)
	require.NoError(t, err)
	user.ID = userID

	project := entity.Project{
		Name:   uuid.NewString(),
		UserID: user.ID,
	}
	projectDB, err := pr.SaveProject(ctx, project)
	require.NoError(t, err)

	err = pr.SoftDeleteProject(ctx, projectDB.ID)
	require.NoError(t, err)

	_, err = pr.ProjectByID(ctx, projectDB.ID)
	require.Error(t, err)
}
