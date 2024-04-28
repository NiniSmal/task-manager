package entity

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestProject_Validate(t *testing.T) {
	testProject := []struct {
		name      string
		project   Project
		wantError bool
	}{
		{
			name:      "ok",
			project:   Project{Name: "launch the application"},
			wantError: false,
		},
		{
			name:      "short name",
			project:   Project{Name: "step"},
			wantError: true,
		},
		{
			name:      "long name",
			project:   Project{Name: strings.Repeat("n", maxNameProject+1)},
			wantError: true,
		},
	}

	for _, tp := range testProject {
		t.Run(tp.name, func(t *testing.T) {
			err := tp.project.Validate()
			if tp.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}

}
