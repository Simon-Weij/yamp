package playlistcmd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoveCommand(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T)
		args      []string
		playlist  string
		wantErr   bool
		wantErrIn string
	}{
		{
			name:     "remove error",
			args:     []string{"artist", "title"},
			playlist: "playlist",
			setup: func(t *testing.T) {
				removeItemFromPlaylistFn = func(string, string, string) error { return errors.New("boom") }
			},
			wantErr:   true,
			wantErrIn: "could not remove song",
		},
		{
			name:     "success",
			args:     []string{"artist", "title"},
			playlist: "playlist",
			setup: func(t *testing.T) {
				removeItemFromPlaylistFn = func(playlistName, artist, title string) error {
					assert.Equal(t, "playlist", playlistName)
					assert.Equal(t, "artist", artist)
					assert.Equal(t, "title", title)
					return nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldRemove := removeItemFromPlaylistFn
			t.Cleanup(func() {
				removeItemFromPlaylistFn = oldRemove
			})

			if tt.setup != nil {
				tt.setup(t)
			}

			removeCmd.Flags().Set("playlist", tt.playlist)
			err := removeCmd.RunE(removeCmd, tt.args)
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
