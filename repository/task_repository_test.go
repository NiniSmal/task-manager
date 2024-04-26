package repository

import (
	"github.com/stretchr/testify/require"
	"gitlab.com/nina8884807/task-manager/entity"
	"testing"
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
