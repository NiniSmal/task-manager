package entity

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestTask_Validate(t *testing.T) {
	testCasesTask := []struct {
		name      string
		task      Task
		wantError bool
	}{
		{
			name:      "ok",
			task:      Task{Name: "learn go"},
			wantError: false,
		},
		{
			name:      "court name",
			task:      Task{Name: "read"},
			wantError: true,
		},
		{
			name:      "long name",
			task:      Task{Name: strings.Repeat("n", maxNameTask+1)},
			wantError: true,
		},
	}

	for _, tt := range testCasesTask {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
