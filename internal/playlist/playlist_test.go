package playlist

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/adrg/xdg"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddItemToPlaylist(t *testing.T) {
	tests := []struct {
		name         string
		playlistName string
		isInternal   bool
		preCreate    bool
		items        []struct {
			artist   string
			title    string
			location string
		}
		wantErr bool
	}{
		{
			name:         "add item to playlist",
			playlistName: "testplaylist",
			preCreate:    true,
			items: []struct {
				artist   string
				title    string
				location string
			}{{"testartist", "testtitle", "testlocation"}},
		},
		{
			name:         "append multiple items",
			playlistName: "testplaylist",
			preCreate:    true,
			items: []struct {
				artist   string
				title    string
				location string
			}{{"artist1", "title1", "location1"}, {"artist2", "title2", "location2"}},
		},
		{
			name:         "error when playlist missing",
			playlistName: "missingplaylist",
			preCreate:    false,
			items: []struct {
				artist   string
				title    string
				location string
			}{{"artist", "title", "location"}},
			wantErr: true,
		},
		{
			name:         "add item to internal playlist",
			playlistName: "internalplaylist",
			isInternal:   true,
			preCreate:    true,
			items: []struct {
				artist   string
				title    string
				location string
			}{{"artist", "title", "location"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldFs := playlistFs
			playlistFs = afero.NewMemMapFs()
			t.Cleanup(func() {
				playlistFs = oldFs
			})

			if tt.preCreate {
				err := CreatePlaylist(tt.playlistName, tt.isInternal)
				require.NoError(t, err)
			}

			for _, item := range tt.items {
				err := AddItemToPlaylist(tt.playlistName, item.artist, item.title, item.location, tt.isInternal)
				if tt.wantErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
			}

			baseDir := filepath.Join(xdg.UserDirs.Music, "playlists")
			if tt.isInternal {
				baseDir = filepath.Join(baseDir, "internal")
			}
			path := filepath.Join(baseDir, tt.playlistName+".m3u")

			data, err := afero.ReadFile(playlistFs, path)
			require.NoError(t, err)

			expected := "#EXTM3U\n"
			for _, item := range tt.items {
				expected += fmt.Sprintf("#EXTINF:-1,%s - %s\n%s\n", item.artist, item.title, item.location)
			}
			assert.Equal(t, expected, string(data))
		})
	}
}

func TestCreatePlaylist(t *testing.T) {
	tests := []struct {
		name                string
		wantPlaylistToExist bool
		playlistName        string
		isInternal          bool
		wantErr             bool
		preCreate           bool
	}{
		{"creates playlist correctly", true, "test", false, false, false},
		{"require content of file to be #EXTM3U", true, "test", false, false, false},
		{"disallowed name", false, ".", false, true, false},
		{"empty name", false, "", false, true, false},
		{"parent dir name", false, "..", false, true, false},
		{"internal name", false, "internal", false, true, false},
		{"isInternal makes playlist in internal", true, "test", true, false, false},
		{"duplicate playlist", false, "dup", false, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldFs := playlistFs
			playlistFs = afero.NewMemMapFs()
			t.Cleanup(func() {
				playlistFs = oldFs
			})

			if tt.preCreate {
				err := CreatePlaylist(tt.playlistName, tt.isInternal)
				require.NoError(t, err)
			}

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
