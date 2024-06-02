package repository

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gitlab.com/nina8884807/task-manager/entity"
)

func TestApplyTaskFilter(t *testing.T) {
	wantQuery := "SELECT id FROM tasks WHERE t.user_id = $1 AND t.project_id = $2"
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

func TestTaskRepository_Create(t *testing.T) {
	db, rds := DBConnection(t)
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
	projectDB, err := pr.SaveProject(ctx, project)
	require.NoError(t, err)
	project.ID = projectDB.ID

	task := entity.Task{
		Name:        uuid.NewString(),
		Description: uuid.NewString(),
		Status:      "not done",
		CreatedAt:   createdAt,
		UserID:      user.ID,
		ProjectID:   projectDB.ID,
	}
	taskDB, err := tr.Create(ctx, task)
	require.NoError(t, err)

	dbTask, err := tr.ByID(ctx, taskDB.ID)
	require.NoError(t, err)
	require.Equal(t, dbTask.Name, task.Name)
	require.Equal(t, dbTask.Description, task.Description)
	require.Equal(t, dbTask.Status, task.Status)
	require.Equal(t, dbTask.CreatedAt.Unix(), dbTask.CreatedAt.Unix())
	require.Equal(t, dbTask.UserID, dbTask.UserID)
	require.Equal(t, dbTask.ProjectID, task.ProjectID)
}

func TestTaskRepository_Tasks(t *testing.T) {
	db, rds := DBConnection(t)
	tr := NewTaskRepository(db, rds)
	ur := NewUserRepository(db, rds)
	pr := NewProjectRepository(db, rds)
	ctx := context.Background()

	user := entity.User{Email: uuid.NewString()}
	userID, err := ur.CreateUser(ctx, user)
	require.NoError(t, err)
	user.ID = userID

	project := entity.Project{UserID: userID}
	projectDB, err := pr.SaveProject(ctx, project)
	require.NoError(t, err)
	project.ID = projectDB.ID
	tasks := []entity.Task{
		{
			Name:      uuid.New().String(),
			Status:    "done",
			CreatedAt: time.Now().Round(time.Millisecond).UTC(),
			UserID:    user.ID,
			ProjectID: projectDB.ID,
		}, {
			Name:      uuid.New().String(),
			Status:    "done",
			CreatedAt: time.Now().Round(time.Millisecond).UTC(),
			UserID:    user.ID,
			ProjectID: project.ID,
		},
	}
	for i, task := range tasks {
		taskDB, err := tr.Create(ctx, task)
		require.NoError(t, err)
		tasks[i].ID = taskDB.ID
	}
	filter := entity.TaskFilter{
		UserID:    strconv.FormatInt(user.ID, 10),
		ProjectID: strconv.FormatInt(project.ID, 10),
	}
	dbTasks, err := tr.Tasks(ctx, filter)
	require.NoError(t, err)

	for _, task := range tasks {
		require.Contains(t, dbTasks, task)
	}

}

func TestTaskRepository_ByID_Error(t *testing.T) {
	db, rds := DBConnection(t)
	tr := NewTaskRepository(db, rds)
	ctx := context.Background()
	_, err := tr.ByID(ctx, 1234)
	require.Error(t, err)
}

func TestTaskRepository_Update(t *testing.T) {
	db, rds := DBConnection(t)
	tr := NewTaskRepository(db, rds)
	ur := NewUserRepository(db, rds)
	pr := NewProjectRepository(db, rds)
	ctx := context.Background()

	user := entity.User{Email: uuid.NewString()}
	userID, err := ur.CreateUser(ctx, user)
	require.NoError(t, err)
	user.ID = userID

	project := entity.Project{UserID: userID}
	projectDB, err := pr.SaveProject(ctx, project)
	require.NoError(t, err)
	project.ID = projectDB.ID
	createdAt := time.Now()
	task := entity.Task{
		Name:        uuid.New().String(),
		Description: uuid.NewString(),
		Status:      "not done",
		CreatedAt:   createdAt,
		UserID:      user.ID,
		ProjectID:   project.ID,
	}
	taskDB, err := tr.Create(ctx, task)
	require.NoError(t, err)

	taskUp := entity.UpdateTask{
		Name:        task.Name,
		Description: task.Description,
		Status:      "done",
		UserID:      user.ID,
		ProjectID:   projectDB.ID,
	}

	err = tr.Update(ctx, taskDB.ID, taskUp)
	require.NoError(t, err)

	taskDB2, err := tr.ByID(ctx, taskDB.ID)
	require.NoError(t, err)
	require.Equal(t, taskUp.Name, taskDB2.Name)
	require.Equal(t, taskUp.Description, taskDB2.Description)
	require.Equal(t, taskUp.Status, taskDB2.Status)

}
func TestTaskRepository_Delete(t *testing.T) {
	db, rds := DBConnection(t)
	tr := NewTaskRepository(db, rds)
	ur := NewUserRepository(db, rds)
	pr := NewProjectRepository(db, rds)
	ctx := context.Background()

	user := entity.User{Email: uuid.NewString()}
	userID, err := ur.CreateUser(ctx, user)
	require.NoError(t, err)
	user.ID = userID

	project := entity.Project{UserID: user.ID}
	projectDB, err := pr.SaveProject(ctx, project)
	require.NoError(t, err)
	project.ID = projectDB.ID

	task := entity.Task{
		Name:      uuid.New().String(),
		ProjectID: project.ID,
		UserID:    user.ID,
	}

	taskDB, err := tr.Create(ctx, task)
	require.NoError(t, err)

	err = tr.Delete(ctx, taskDB.ID)
	require.NoError(t, err)

	_, err = tr.ByID(ctx, taskDB.ID)
	require.Error(t, err)

}
