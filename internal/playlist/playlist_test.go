package playlist

import (
	"path/filepath"
	"testing"

	"github.com/adrg/xdg"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePlaylist(t *testing.T) {
	tests := []struct {
		name                string
		wantPlaylistToExist bool
		playlistName        string
		isInternal          bool
		wantErr             bool
	}{
		{"creates playlist correctly", true, "test", false, false},
		{"require content of file to be #EXTM3U", true, "test", false, false},
		{"disallowed name", false, ".", false, true},
		{"isInternal makes playlist in internal", true, "test", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldFs := playlistFs
			playlistFs = afero.NewMemMapFs()
			t.Cleanup(func() {
				playlistFs = oldFs
			})

			err := CreatePlaylist(tt.playlistName, tt.isInternal)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			baseDir := filepath.Join(xdg.UserDirs.Music, "playlists")
			if tt.isInternal {
				baseDir = filepath.Join(baseDir, "internal")
			}

			path := filepath.Join(baseDir, tt.playlistName+".m3u")

			exists, err := afero.Exists(playlistFs, path)
			require.NoError(t, err)
			assert.Equal(t, tt.wantPlaylistToExist, exists)

			data, err := afero.ReadFile(playlistFs, path)
			require.NoError(t, err)

			assert.Equal(t, "#EXTM3U\n", string(data))
		})
	}
}
