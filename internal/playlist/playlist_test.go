package playlist

import (
	"fmt"
	"path/filepath"
	"testing"
	"yamp/internal/musicdiscovery"

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

func TestListPlaylistItems(t *testing.T) {
	tests := []struct {
		name         string
		playlistName string
		isInternal   bool
		setup        func(t *testing.T)
		want         []musicdiscovery.Metadata
		wantErr      bool
	}{
		{
			name:         "missing playlist returns error",
			playlistName: "missing",
			wantErr:      true,
		},
		{
			name:         "empty playlist returns error",
			playlistName: "empty",
			setup: func(t *testing.T) {
				err := CreatePlaylist("empty", false)
				require.NoError(t, err)
			},
			wantErr: true,
		},
		{
			name:         "list single item",
			playlistName: "single",
			setup: func(t *testing.T) {
				err := CreatePlaylist("single", false)
				require.NoError(t, err)
				err = AddItemToPlaylist("single", "artist", "title", "/path/one", false)
				require.NoError(t, err)
			},
			want: []musicdiscovery.Metadata{{Artist: "artist", Title: "title"}},
		},
		{
			name:         "list multiple items",
			playlistName: "multi",
			setup: func(t *testing.T) {
				err := CreatePlaylist("multi", false)
				require.NoError(t, err)
				err = AddItemToPlaylist("multi", "artist1", "title1", "/path/one", false)
				require.NoError(t, err)
				err = AddItemToPlaylist("multi", "artist2", "title2", "/path/two", false)
				require.NoError(t, err)
			},
			want: []musicdiscovery.Metadata{{Artist: "artist1", Title: "title1"}, {Artist: "artist2", Title: "title2"}},
		},
		{
			name:         "internal playlist",
			playlistName: "internaltest",
			isInternal:   true,
			setup: func(t *testing.T) {
				err := CreatePlaylist("internaltest", true)
				require.NoError(t, err)
				err = AddItemToPlaylist("internaltest", "artist", "title", "/path/internal", true)
				require.NoError(t, err)
			},
			want: []musicdiscovery.Metadata{{Artist: "artist", Title: "title"}},
		},
		{
			name:         "ignores non-song lines",
			playlistName: "mixed",
			setup: func(t *testing.T) {
				err := CreatePlaylist("mixed", false)
				require.NoError(t, err)
				baseDir := filepath.Join(xdg.UserDirs.Music, "playlists")
				path := filepath.Join(baseDir, "mixed.m3u")
				content := "#EXTM3U\n# comment\n#EXTINF:-1,artist - title\n/path/one\n#EXTINF:-1,artist2 - title2\n/path/two\n"
				err = afero.WriteFile(playlistFs, path, []byte(content), 0644)
				require.NoError(t, err)
			},
			want: []musicdiscovery.Metadata{{Artist: "artist", Title: "title"}, {Artist: "artist2", Title: "title2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldFs := playlistFs
			playlistFs = afero.NewMemMapFs()
			t.Cleanup(func() {
				playlistFs = oldFs
			})

			if tt.setup != nil {
				tt.setup(t)
			}

			items, err := ListPlaylistItems(tt.playlistName, tt.isInternal)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, items)
		})
	}
}

func TestRemoveItemFromPlaylist(t *testing.T) {
	type playlistItem struct {
		artist   string
		title    string
		location string
	}

	tests := []struct {
		name         string
		playlistName string
		setup        func(t *testing.T)
		removeArtist string
		removeTitle  string
		wantItems    []playlistItem
		wantErr      bool
	}{
		{
			name:         "remove single item",
			playlistName: "single",
			setup: func(t *testing.T) {
				err := CreatePlaylist("single", false)
				require.NoError(t, err)
				err = AddItemToPlaylist("single", "artist", "title", "/path/one", false)
				require.NoError(t, err)
			},
			removeArtist: "artist",
			removeTitle:  "title",
			wantItems:    []playlistItem{},
		},
		{
			name:         "remove one of multiple items",
			playlistName: "multi",
			setup: func(t *testing.T) {
				err := CreatePlaylist("multi", false)
				require.NoError(t, err)
				err = AddItemToPlaylist("multi", "artist1", "title1", "/path/one", false)
				require.NoError(t, err)
				err = AddItemToPlaylist("multi", "artist2", "title2", "/path/two", false)
				require.NoError(t, err)
			},
			removeArtist: "artist1",
			removeTitle:  "title1",
			wantItems:    []playlistItem{{artist: "artist2", title: "title2", location: "/path/two"}},
		},
		{
			name:         "case-insensitive match",
			playlistName: "case",
			setup: func(t *testing.T) {
				err := CreatePlaylist("case", false)
				require.NoError(t, err)
				err = AddItemToPlaylist("case", "Artist", "Title", "/path/one", false)
				require.NoError(t, err)
			},
			removeArtist: "artist",
			removeTitle:  "title",
			wantItems:    []playlistItem{},
		},
		{
			name:         "song not found",
			playlistName: "notfound",
			setup: func(t *testing.T) {
				err := CreatePlaylist("notfound", false)
				require.NoError(t, err)
				err = AddItemToPlaylist("notfound", "artist", "title", "/path/one", false)
				require.NoError(t, err)
			},
			removeArtist: "other",
			removeTitle:  "song",
			wantErr:      true,
		},
		{
			name:         "missing playlist",
			playlistName: "missing",
			removeArtist: "artist",
			removeTitle:  "title",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldFs := playlistFs
			playlistFs = afero.NewMemMapFs()
			t.Cleanup(func() {
				playlistFs = oldFs
			})

			if tt.setup != nil {
				tt.setup(t)
			}

			err := RemoveItemFromPlaylist(tt.playlistName, tt.removeArtist, tt.removeTitle)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			baseDir := filepath.Join(xdg.UserDirs.Music, "playlists")
			path := filepath.Join(baseDir, tt.playlistName+".m3u")

			data, err := afero.ReadFile(playlistFs, path)
			require.NoError(t, err)

			expected := "#EXTM3U\n"
			for _, item := range tt.wantItems {
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
