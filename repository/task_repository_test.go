package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"gitlab.com/nina8884807/task-manager/entity"
	"strconv"
	"testing"
	"time"
)

func TestApplyTaskFilter(t *testing.T) {
	wantQuery := "SELECT id FROM tasks WHERE user_id = $1 AND project_id = $2"
	wantArgs := []any{"3", "13"}

	query := "SELECT id FROM tasks"
	f := entity.TaskFilter{
		UserID:    "3",
		ProjectID: "13",
	}
	query, args := applyTaskFilter(query, f)

	require.Equal(t, wantQuery, query)
	require.Equal(t, wantArgs, args)
}

func TaskConnection(t *testing.T) (*sql.DB, *redis.Client) {
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

func TestTaskRepository_Create(t *testing.T) {
	db, rds := TaskConnection(t)
	tr := NewTaskRepository(db, rds)
	ur := NewUserRepository(db, rds)
	pr := NewProjectRepository(db, rds)
	ctx := context.Background()
	createdAt := time.Now()
	user := entity.User{Email: uuid.NewString()}
	userID, err := ur.CreateUser(ctx, user)
	require.NoError(t, err)
	user.ID = userID

	project := entity.Project{UserID: user.ID}
	projectID, err := pr.SaveProject(ctx, project)
	require.NoError(t, err)
	project.ID = projectID

	task := entity.Task{
		Name:      uuid.New().String(),
		Status:    "not done",
		CreatedAt: createdAt,
		UserID:    user.ID,
		ProjectID: project.ID,
	}
	taskID, err := tr.Create(ctx, task)
	require.NoError(t, err)

	dbTask, err := tr.ByID(ctx, taskID)
	require.NoError(t, err)
	require.Equal(t, dbTask.Name, task.Name)
	require.Equal(t, dbTask.Status, task.Status)
	require.Equal(t, dbTask.CreatedAt.Unix(), dbTask.CreatedAt.Unix())
	require.Equal(t, dbTask.UserID, dbTask.UserID)
	require.Equal(t, dbTask.ProjectID, task.ProjectID)
}

func TestTaskRepository_Tasks(t *testing.T) {
	db, rds := TaskConnection(t)
	tr := NewTaskRepository(db, rds)
	ctx := context.Background()

	userID := int64(1)
	tasks := []entity.Task{
		{
			Name:      uuid.New().String(),
			Status:    "done",
			CreatedAt: time.Now().Round(time.Millisecond).UTC(),
			UserID:    userID,
			ProjectID: 2,
		}, {
			Name:      uuid.New().String(),
			Status:    "done",
			CreatedAt: time.Now().Round(time.Millisecond).UTC(),
			UserID:    userID,
			ProjectID: 2,
		},
	}
	for i, task := range tasks {
		id, err := tr.Create(ctx, task)
		require.NoError(t, err)
		tasks[i].ID = id
	}
	filter := entity.TaskFilter{
		UserID:    strconv.FormatInt(userID, 10),
		ProjectID: "2",
	}
	dbTasks, err := tr.Tasks(ctx, filter)
	require.NoError(t, err)

	for _, task := range tasks {
		require.Contains(t, dbTasks, task)
	}

}

func TestTaskRepository_ByID_Error(t *testing.T) {
	db, rds := TaskConnection(t)
	tr := NewTaskRepository(db, rds)
	ctx := context.Background()
	_, err := tr.ByID(ctx, 1234)
	require.Error(t, err)
}

func TestTaskRepository_Update(t *testing.T) {
	db, rds := TaskConnection(t)
	tr := NewTaskRepository(db, rds)
	ctx := context.Background()
	task := entity.Task{
		Name:      uuid.New().String(),
		Status:    "not done",
		CreatedAt: time.Now(),
		UserID:    1,
		ProjectID: 2,
	}
	id, err := tr.Create(ctx, task)
	require.NoError(t, err)

	taskUp := entity.UpdateTask{
		Name:   task.Name,
		Status: "done",
		UserID: 1,
	}

	idDB, err := tr.Update(ctx, id, taskUp)
	require.NoError(t, err)

	taskDB, err := tr.ByID(ctx, idDB)
	require.NoError(t, err)
	require.Equal(t, taskDB, taskUp)

}
func TestTaskRepository_Delete(t *testing.T) {
	db, rds := TaskConnection(t)
	tr := NewTaskRepository(db, rds)
	ctx := context.Background()

	task := entity.Task{
		Name:      uuid.New().String(),
		ProjectID: 2,
	}

	id, err := tr.Create(ctx, task)
	require.NoError(t, err)

	err = tr.Delete(ctx, id)
	require.NoError(t, err)

	_, err = tr.ByID(ctx, id)
	require.Error(t, err)

}
