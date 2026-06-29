package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type errorFs struct {
	afero.Fs
	shouldFail bool
}

func (efs *errorFs) Create(name string) (afero.File, error) {
	if efs.shouldFail {
		return nil, os.ErrPermission
	}
	return efs.Fs.Create(name)
}

func makeEmptyFile(t *testing.T, fs afero.Fs, path string) {
	t.Helper()
	err := afero.WriteFile(fs, path, []byte("[]"), 0o644)
	require.NoError(t, err)
}

func (efs *errorFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	if efs.shouldFail && (flag&os.O_WRONLY != 0 || flag&os.O_RDWR != 0 || flag&os.O_CREATE != 0) {
		return nil, os.ErrPermission
	}
	return efs.Fs.OpenFile(name, flag, perm)
}

func Test_validateFilename(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid name", "my_playlist", false},
		{"valid name with numbers", "playlist123", false},
		{"valid name with dash and underscore", "play-list_new", false},
		{"invalid char slash", "play/list", true},
		{"invalid char backslash", "play\\list", true},
		{"invalid char colon", "play:list", true},
		{"invalid char asterisk", "play*list", true},
		{"invalid char question mark", "play?list", true},
		{"invalid char double quote", "play\"list", true},
		{"invalid char less than", "play<list", true},
		{"invalid char greater than", "play>list", true},
		{"invalid char pipe", "play|list", true},
		{"invalid dot dot", "play..list", true},
		{"invalid starts with dot", ".playlist", true},
		{"invalid empty name", "", true},
		{"invalid whitespace only", "   ", true},
		{"invalid leading space", " playlist", true},
		{"invalid trailing space", "playlist ", true},
		{"valid boundary length 254", strings.Repeat("a", 254), false},
		{"invalid length 255", strings.Repeat("a", 255), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFilename(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_createPlaylist(t *testing.T) {
	tests := []struct {
		name     string
		plName   string
		wantErr  bool
		readOnly bool
	}{
		{
			name:   "create playlist",
			plName: "testplaylist",
		},
		{
			name:     "create playlist directory fails (read-only filesystem)",
			plName:   "testplaylist",
			wantErr:  true,
			readOnly: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fs afero.Fs = afero.NewMemMapFs()
			if tt.readOnly {
				fs = afero.NewReadOnlyFs(fs)
			}
			pr := &PlaylistRepository{Fs: fs}
			path, err := pr.CreatePlaylist(tt.plName)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			exists, err := afero.Exists(pr.Fs, path)
			require.NoError(t, err)
			require.True(t, exists)

			_, err = pr.CreatePlaylist(tt.plName)
			assert.Error(t, err)
		})
	}
}

func GetTestSong() Song {
	return Song{
		TrackID:        1,
		ArtistID:       10,
		CollectionID:   100,
		TrackName:      "test track",
		Artist:         "test artist",
		CollectionName: "test collection",
		Cover:          "http://example.com/art.jpg",
		DurationMillis: 180000,
		PreviewURL:     "http://example.com/preview.mp3",
		TrackURL:       "http://example.com/track",
	}
}

func assertPlaylistContents(t *testing.T, fs afero.Fs, playlistName string, expectedSongs []Song) {
	path := getPlaylistPath(playlistName)
	exists, err := afero.Exists(fs, path)
	require.NoError(t, err)
	require.True(t, exists, "playlist file should exist")

	data, err := afero.ReadFile(fs, path)
	require.NoError(t, err)

	var actualSongs []Song
	err = json.Unmarshal(data, &actualSongs)
	require.NoError(t, err, "playlist data should be valid JSON")

	assert.Equal(t, expectedSongs, actualSongs)
}

func TestPlaylistRepository_AddSongToPlaylist(t *testing.T) {
	songA := GetTestSong()
	songB := Song{
		TrackID:        2,
		ArtistID:       20,
		CollectionID:   200,
		TrackName:      "another track",
		Artist:         "another artist",
		CollectionName: "another collection",
		Cover:          "http://example.com/art2.jpg",
		DurationMillis: 200000,
		PreviewURL:     "http://example.com/preview2.mp3",
		TrackURL:       "http://example.com/track2",
	}

	tests := []struct {
		name         string
		Fs           afero.Fs
		song         Song
		playlistName string
		wantErr      bool
		setup        func(pr *PlaylistRepository) error
		verify       func(t *testing.T, pr *PlaylistRepository)
	}{
		{
			name:         "playlist does not exist",
			Fs:           afero.NewMemMapFs(),
			song:         songA,
			playlistName: "nonExistentPlaylist",
			wantErr:      true,
			setup:        func(pr *PlaylistRepository) error { return nil },
			verify: func(t *testing.T, pr *PlaylistRepository) {
				path := getPlaylistPath("nonExistentPlaylist")
				exists, err := afero.Exists(pr.Fs, path)
				assert.NoError(t, err)
				assert.False(t, exists)
			},
		},
		{
			name:         "playlist exists",
			Fs:           afero.NewMemMapFs(),
			song:         songA,
			playlistName: "testPlaylist",
			wantErr:      false,
			setup: func(pr *PlaylistRepository) error {
				_, err := pr.CreatePlaylist("testPlaylist")
				return err
			},
			verify: func(t *testing.T, pr *PlaylistRepository) {
				assertPlaylistContents(t, pr.Fs, "testPlaylist", []Song{songA})
			},
		},
		{
			name:         "playlist exists and has existing songs",
			Fs:           afero.NewMemMapFs(),
			song:         songB,
			playlistName: "multiSongPlaylist",
			wantErr:      false,
			setup: func(pr *PlaylistRepository) error {
				_, err := pr.CreatePlaylist("multiSongPlaylist")
				if err != nil {
					return err
				}
				return pr.AddSongToPlaylist(songA, "multiSongPlaylist")
			},
			verify: func(t *testing.T, pr *PlaylistRepository) {
				assertPlaylistContents(t, pr.Fs, "multiSongPlaylist", []Song{songA, songB})
			},
		},
		{
			name:         "playlist exists but has invalid JSON content",
			Fs:           afero.NewMemMapFs(),
			song:         songA,
			playlistName: "corruptPlaylist",
			wantErr:      false,
			setup: func(pr *PlaylistRepository) error {
				path, err := pr.CreatePlaylist("corruptPlaylist")
				if err != nil {
					return err
				}
				return afero.WriteFile(pr.Fs, path, []byte("invalid-json{"), 0o644)
			},
			verify: func(t *testing.T, pr *PlaylistRepository) {
				assertPlaylistContents(t, pr.Fs, "corruptPlaylist", []Song{songA})
			},
		},
		{
			name:         "playlist exists but file is empty",
			Fs:           afero.NewMemMapFs(),
			song:         songA,
			playlistName: "emptyPlaylistFile",
			wantErr:      false,
			setup: func(pr *PlaylistRepository) error {
				path, err := pr.CreatePlaylist("emptyPlaylistFile")
				if err != nil {
					return err
				}
				return afero.WriteFile(pr.Fs, path, []byte(""), 0o644)
			},
			verify: func(t *testing.T, pr *PlaylistRepository) {
				assertPlaylistContents(t, pr.Fs, "emptyPlaylistFile", []Song{songA})
			},
		},
		{
			name:         "write failure during AddSongToPlaylist",
			Fs:           &errorFs{Fs: afero.NewMemMapFs()},
			song:         songA,
			playlistName: "writeFailurePlaylist",
			wantErr:      true,
			setup: func(pr *PlaylistRepository) error {
				_, err := pr.CreatePlaylist("writeFailurePlaylist")
				if err != nil {
					return err
				}
				pr.Fs.(*errorFs).shouldFail = true
				return nil
			},
			verify: func(t *testing.T, pr *PlaylistRepository) {
				pr.Fs.(*errorFs).shouldFail = false
				assertPlaylistContents(t, pr.Fs, "writeFailurePlaylist", []Song{})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := &PlaylistRepository{Fs: tt.Fs}
			err := tt.setup(pr)
			require.NoError(t, err)

			err = pr.AddSongToPlaylist(tt.song, tt.playlistName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.verify != nil {
				tt.verify(t, pr)
			}
		})
	}
}

func TestPlaylistRepository_ListPlaylists(t *testing.T) {
	tests := []struct {
		name  string
		Fs    afero.Fs
		setup func(pr *PlaylistRepository) error
		want  []string
	}{
		{
			name: "should return correct playlists",
			Fs:   afero.NewMemMapFs(),
			setup: func(pr *PlaylistRepository) error {
				for _, name := range []string{"playlist1", "playlist2", "playlist3"} {
					makeEmptyFile(t, pr.Fs, getPlaylistPath(name))
				}
				return nil
			},
			want: []string{"playlist1", "playlist2", "playlist3"},
		},
		{
			name: "should return correct playlists (empty)",
			Fs:   afero.NewMemMapFs(),
			setup: func(pr *PlaylistRepository) error {
				for _, name := range []string{} {
					makeEmptyFile(t, pr.Fs, getPlaylistPath(name))
				}
				return nil
			},
			want: []string{},
		},
		{
			name: "should not list non-json files",
			Fs:   afero.NewMemMapFs(),
			setup: func(pr *PlaylistRepository) error {
				makeEmptyFile(t, pr.Fs, getPlaylistDir()+"playlist.txt")
				return nil
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := &PlaylistRepository{Fs: tt.Fs}
			err := tt.setup(pr)
			require.NoError(t, err)

			got := pr.ListPlaylists()
			assert.Equal(t, tt.want, got)
		})
	}
}
