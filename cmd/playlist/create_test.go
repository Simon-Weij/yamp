package playlistcmd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCommand(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		setup     func(t *testing.T)
		wantErr   bool
		wantErrIn string
	}{
		{
			name: "create error",
			args: []string{"playlist"},
			setup: func(t *testing.T) {
				createPlaylistFn = func(string, bool) error { return errors.New("boom") }
			},
			wantErr:   true,
			wantErrIn: "could not create playlist",
		},
		{
			name: "success",
			args: []string{"playlist"},
			setup: func(t *testing.T) {
				createPlaylistFn = func(string, bool) error { return nil }
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldCreate := createPlaylistFn
			t.Cleanup(func() {
				createPlaylistFn = oldCreate
			})

			if tt.setup != nil {
				tt.setup(t)
			}

			err := createCmd.RunE(createCmd, tt.args)
			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrIn != "" {
					assert.Contains(t, err.Error(), tt.wantErrIn)
				}
				return
			}
			require.NoError(t, err)
		})
	}
}
