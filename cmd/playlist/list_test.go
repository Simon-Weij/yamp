package playlistcmd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlaylistListCommand(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T)
		wantErr   bool
		wantErrIn string
	}{
		{
			name: "list error",
			setup: func(t *testing.T) {
				listPlaylistsFn = func() ([]string, error) { return nil, errors.New("boom") }
			},
			wantErr:   true,
			wantErrIn: "boom",
		},
		{
			name: "success",
			setup: func(t *testing.T) {
				listPlaylistsFn = func() ([]string, error) { return []string{"one", "two"}, nil }
			},
		},
		{
			name: "empty",
			setup: func(t *testing.T) {
				listPlaylistsFn = func() ([]string, error) { return []string{}, nil }
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldList := listPlaylistsFn
			t.Cleanup(func() {
				listPlaylistsFn = oldList
			})

			if tt.setup != nil {
				tt.setup(t)
			}

			err := playlistListCmd.RunE(playlistListCmd, []string{})
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
