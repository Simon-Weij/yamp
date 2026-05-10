package playlistcmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"yamp/internal/musicdiscovery"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func useMemFs(t *testing.T) {
	oldFs := fs
	oldMkdir := mkdirAll
	fs = afero.NewMemMapFs()
	mkdirAll = fs.MkdirAll
	t.Cleanup(func() {
		fs = oldFs
		mkdirAll = oldMkdir
	})
}

func TestAddCommand(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T)
		args      []string
		wantErr   bool
		wantErrIn string
	}{
		{
			name: "playlist missing",
			args: []string{"missing", "song"},
			setup: func(t *testing.T) {
				playlistExistsFn = func(string, bool) (bool, error) { return false, nil }
			},
			wantErr:   true,
			wantErrIn: "playlist doesn't exist",
		},
		{
			name: "playlist exists check error",
			args: []string{"playlist", "song"},
			setup: func(t *testing.T) {
				playlistExistsFn = func(string, bool) (bool, error) { return false, errors.New("boom") }
			},
			wantErr:   true,
			wantErrIn: "checking is playlist exists",
		},
		{
			name: "download error",
			args: []string{"playlist", "song"},
			setup: func(t *testing.T) {
				playlistExistsFn = func(string, bool) (bool, error) { return true, nil }
				uuidNewString = func() string { return "id" }
				tempDir = func() string { return "/tmp" }
				downloadSong = func(string, string) error { return errors.New("download failed") }
			},
			wantErr:   true,
			wantErrIn: "downloading song",
		},
		{
			name: "extract metadata error",
			args: []string{"playlist", "song"},
			setup: func(t *testing.T) {
				playlistExistsFn = func(string, bool) (bool, error) { return true, nil }
				uuidNewString = func() string { return "id" }
				tempDir = func() string { return "/tmp" }
				downloadSong = func(string, string) error { return nil }
				extractMetadata = func(string) (*musicdiscovery.Metadata, error) { return nil, errors.New("bad metadata") }
			},
			wantErr:   true,
			wantErrIn: "could not initialise",
		},
		{
			name: "add item error",
			args: []string{"playlist", "song"},
			setup: func(t *testing.T) {
				playlistExistsFn = func(string, bool) (bool, error) { return true, nil }
				uuidNewString = func() string { return "id" }
				tempDir = func() string { return "/tmp" }
				downloadSong = func(string, string) error { return nil }
				extractMetadata = func(string) (*musicdiscovery.Metadata, error) {
					return &musicdiscovery.Metadata{Artist: "artist", Album: "album", Title: "title"}, nil
				}
				userMusicDir = func() string { return "/music" }
				mkdirAll = func(string, os.FileMode) error { return nil }
				moveFileFn = func(string, string) error { return nil }
				addItemToPlaylist = func(string, string, string, string, bool) error { return errors.New("add failed") }
			},
			wantErr:   true,
			wantErrIn: "could not add song",
		},
		{
			name: "success",
			args: []string{"playlist", "song"},
			setup: func(t *testing.T) {
				playlistExistsFn = func(string, bool) (bool, error) { return true, nil }
				uuidNewString = func() string { return "id" }
				tempDir = func() string { return "/tmp" }
				downloadSong = func(song, path string) error {
					expected := filepath.Join("/tmp", "yamp", "id", "song.mp3")
					assert.Equal(t, expected, path)
					return nil
				}
				extractMetadata = func(string) (*musicdiscovery.Metadata, error) {
					return &musicdiscovery.Metadata{Artist: "artist", Album: "album", Title: "title"}, nil
				}
				userMusicDir = func() string { return "/music" }
				mkdirAll = func(path string, _ os.FileMode) error {
					expected := filepath.Join("/music", "yamp", "artist", "album")
					assert.Equal(t, expected, path)
					return nil
				}
				moveFileFn = func(src, dst string) error {
					expectedSrc := filepath.Join("/tmp", "yamp", "id", "song.mp3")
					expectedDst := filepath.Join("/music", "yamp", "artist", "album", "title.mp3")
					assert.Equal(t, expectedSrc, src)
					assert.Equal(t, expectedDst, dst)
					return nil
				}
				addItemToPlaylist = func(name, artist, title, location string, _ bool) error {
					expected := filepath.Join("/music", "yamp", "artist", "album", "title.mp3")
					assert.Equal(t, "playlist", name)
					assert.Equal(t, "artist", artist)
					assert.Equal(t, "title", title)
					assert.Equal(t, expected, location)
					return nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useMemFs(t)
			oldPlaylistExists := playlistExistsFn
			oldDownload := downloadSong
			oldExtract := extractMetadata
			oldAddItem := addItemToPlaylist
			oldMove := moveFileFn
			oldUserMusicDir := userMusicDir
			oldTempDir := tempDir
			oldUUID := uuidNewString
			t.Cleanup(func() {
				playlistExistsFn = oldPlaylistExists
				downloadSong = oldDownload
				extractMetadata = oldExtract
				addItemToPlaylist = oldAddItem
				moveFileFn = oldMove
				userMusicDir = oldUserMusicDir
				tempDir = oldTempDir
				uuidNewString = oldUUID
			})

			if tt.setup != nil {
				tt.setup(t)
			}

			err := addCmd.RunE(addCmd, tt.args)
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

func TestInitialiseSong(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T)
		wantErr   bool
		wantErrIn string
	}{
		{
			name: "extract metadata error",
			setup: func(t *testing.T) {
				extractMetadata = func(string) (*musicdiscovery.Metadata, error) { return nil, errors.New("bad metadata") }
			},
			wantErr:   true,
			wantErrIn: "failed to extract metadata",
		},
		{
			name: "mkdir error",
			setup: func(t *testing.T) {
				extractMetadata = func(string) (*musicdiscovery.Metadata, error) {
					return &musicdiscovery.Metadata{Artist: "artist", Album: "album", Title: "title"}, nil
				}
				userMusicDir = func() string { return "/music" }
				mkdirAll = func(string, os.FileMode) error { return errors.New("mkdir failed") }
			},
			wantErr:   true,
			wantErrIn: "creating directory",
		},
		{
			name: "move error",
			setup: func(t *testing.T) {
				extractMetadata = func(string) (*musicdiscovery.Metadata, error) {
					return &musicdiscovery.Metadata{Artist: "artist", Album: "album", Title: "title"}, nil
				}
				userMusicDir = func() string { return "/music" }
				mkdirAll = func(string, os.FileMode) error { return nil }
				moveFileFn = func(string, string) error { return errors.New("move failed") }
			},
			wantErr:   true,
			wantErrIn: "moving file",
		},
		{
			name: "success",
			setup: func(t *testing.T) {
				extractMetadata = func(string) (*musicdiscovery.Metadata, error) {
					return &musicdiscovery.Metadata{Artist: "artist", Album: "album", Title: "title"}, nil
				}
				userMusicDir = func() string { return "/music" }
				mkdirAll = func(path string, _ os.FileMode) error {
					expected := filepath.Join("/music", "yamp", "artist", "album")
					assert.Equal(t, expected, path)
					return nil
				}
				moveFileFn = func(src, dst string) error {
					expectedDst := filepath.Join("/music", "yamp", "artist", "album", "title.mp3")
					assert.Equal(t, "/tmp/song.mp3", src)
					assert.Equal(t, expectedDst, dst)
					return nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useMemFs(t)
			oldExtract := extractMetadata
			oldUserMusicDir := userMusicDir
			oldMove := moveFileFn
			t.Cleanup(func() {
				extractMetadata = oldExtract
				userMusicDir = oldUserMusicDir
				moveFileFn = oldMove
			})

			if tt.setup != nil {
				tt.setup(t)
			}

			_, _, err := initaliseSong("song", "/tmp/song.mp3")
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

func TestMoveFile(t *testing.T) {
	useMemFs(t)
	tmp := "/tmp"
	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "dst.txt")

	content := []byte("hello")
	err := afero.WriteFile(fs, src, content, 0644)
	require.NoError(t, err)

	err = moveFile(src, dst)
	require.NoError(t, err)

	_, err = fs.Stat(src)
	require.Error(t, err)

	data, err := afero.ReadFile(fs, dst)
	require.NoError(t, err)
	assert.Equal(t, content, data)
}

func TestInitialiseSongReturnsDestination(t *testing.T) {
	useMemFs(t)
	oldExtract := extractMetadata
	oldUserMusicDir := userMusicDir
	oldMove := moveFileFn
	t.Cleanup(func() {
		extractMetadata = oldExtract
		userMusicDir = oldUserMusicDir
		moveFileFn = oldMove
	})

	extractMetadata = func(string) (*musicdiscovery.Metadata, error) {
		return &musicdiscovery.Metadata{Artist: "artist", Album: "album", Title: "title"}, nil
	}
	userMusicDir = func() string { return "/music" }
	mkdirAll = func(string, os.FileMode) error { return nil }
	moveFileFn = func(string, string) error { return nil }

	path, meta, err := initaliseSong("song", "/tmp/song.mp3")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join("/music", "yamp", "artist", "album", "title.mp3"), path)
	require.NotNil(t, meta)
}

func TestAddCommandErrorMessage(t *testing.T) {
	useMemFs(t)
	oldPlaylistExists := playlistExistsFn
	oldDownload := downloadSong
	oldExtract := extractMetadata
	oldAddItem := addItemToPlaylist
	oldMove := moveFileFn
	oldUserMusicDir := userMusicDir
	oldTempDir := tempDir
	oldUUID := uuidNewString
	t.Cleanup(func() {
		playlistExistsFn = oldPlaylistExists
		downloadSong = oldDownload
		extractMetadata = oldExtract
		addItemToPlaylist = oldAddItem
		moveFileFn = oldMove
		userMusicDir = oldUserMusicDir
		tempDir = oldTempDir
		uuidNewString = oldUUID
	})

	playlistExistsFn = func(string, bool) (bool, error) { return true, nil }
	uuidNewString = func() string { return "id" }
	tempDir = func() string { return "/tmp" }
	downloadSong = func(string, string) error { return nil }
	extractMetadata = func(string) (*musicdiscovery.Metadata, error) {
		return &musicdiscovery.Metadata{Artist: "artist", Album: "album", Title: "title"}, nil
	}
	userMusicDir = func() string { return "/music" }
	mkdirAll = func(string, os.FileMode) error { return nil }
	moveFileFn = func(string, string) error { return nil }
	addItemToPlaylist = func(string, string, string, string, bool) error { return nil }

	err := addCmd.RunE(addCmd, []string{"playlist", "song"})
	require.NoError(t, err)
}

func TestInitialiseSongMoveFileErrorIncludesPaths(t *testing.T) {
	useMemFs(t)
	oldExtract := extractMetadata
	oldUserMusicDir := userMusicDir
	oldMove := moveFileFn
	t.Cleanup(func() {
		extractMetadata = oldExtract
		userMusicDir = oldUserMusicDir
		moveFileFn = oldMove
	})

	extractMetadata = func(string) (*musicdiscovery.Metadata, error) {
		return &musicdiscovery.Metadata{Artist: "artist", Album: "album", Title: "title"}, nil
	}
	userMusicDir = func() string { return "/music" }
	mkdirAll = func(string, os.FileMode) error { return nil }
	moveFileFn = func(string, string) error { return fmt.Errorf("boom") }

	_, _, err := initaliseSong("song", "/tmp/song.mp3")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "/tmp/song.mp3")
	assert.Contains(t, err.Error(), filepath.Join("/music", "yamp", "artist", "album", "title.mp3"))
}
