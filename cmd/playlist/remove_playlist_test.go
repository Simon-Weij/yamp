package playlistcmd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemovePlaylistCommand(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		setup     func(t *testing.T)
		wantErr   bool
		wantErrIn string
	}{
		{
			name: "delete error",
			args: []string{"playlist"},
			setup: func(t *testing.T) {
				deletePlaylistFn = func(string) error { return errors.New("boom") }
			},
			wantErr:   true,
			wantErrIn: "could not remove playlist",
		},
		{
			name: "success",
			args: []string{"playlist"},
			setup: func(t *testing.T) {
				deletePlaylistFn = func(string) error { return nil }
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldDelete := deletePlaylistFn
			t.Cleanup(func() {
				deletePlaylistFn = oldDelete
			})

			if tt.setup != nil {
				tt.setup(t)
			}

			err := removePlaylistCmd.RunE(removePlaylistCmd, tt.args)
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
