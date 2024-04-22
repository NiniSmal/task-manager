package entity

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestUser_Validation(t *testing.T) {
	testCases := []struct {
		name      string
		user      User
		wantError bool
	}{
		{
			name:      "ok",
			user:      User{Login: "ninamusatova90+1@gmail.com", Password: "Querty123!"},
			wantError: false,
		},
		{
			name:      "empty password",
			user:      User{Login: "ninamusatova90+1@gmail.com", Password: ""},
			wantError: true,
		},
		{
			name:      "empty login",
			user:      User{Login: "", Password: "Querty123!"},
			wantError: true,
		},
		{
			name:      "empty login and empty password",
			user:      User{Login: "", Password: ""},
			wantError: true,
		},
		{
			name:      "long login",
			user:      User{Login: strings.Repeat("n", maxLogin+1), Password: "Querty123!"},
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.user.Validate()
			if tc.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
