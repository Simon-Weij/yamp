package cmd

import (
	"errors"
	"path/filepath"
	"testing"
	"os"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlayCommand(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T)
		args      []string
		wantErr   bool
		wantErrIn string
	}{
		{
			name: "mkdir error",
			args: []string{"song"},
			setup: func(t *testing.T) {
				tempDir = func() string { return "/tmp" }
				mkdirAll = func(string, os.FileMode) error { return errors.New("boom") }
			},
			wantErr:   true,
			wantErrIn: "couldn't create directory",
		},
		{
			name: "download error",
			args: []string{"song"},
			setup: func(t *testing.T) {
				tempDir = func() string { return "/tmp" }
				mkdirAll = func(string, os.FileMode) error { return nil }
				downloadSong = func(string, string) error { return errors.New("boom") }
			},
			wantErr:   true,
			wantErrIn: "could not download song",
		},
		{
			name: "play error",
			args: []string{"song"},
			setup: func(t *testing.T) {
				tempDir = func() string { return "/tmp" }
				mkdirAll = func(string, os.FileMode) error { return nil }
				downloadSong = func(string, string) error { return nil }
				playSong = func(string) error { return errors.New("boom") }
			},
			wantErr:   true,
			wantErrIn: "could not play song",
		},
		{
			name: "cleanup error",
			args: []string{"song"},
			setup: func(t *testing.T) {
				tempDir = func() string { return "/tmp" }
				mkdirAll = func(string, os.FileMode) error { return nil }
				downloadSong = func(string, string) error { return nil }
				playSong = func(string) error { return nil }
				removeAll = func(string) error { return errors.New("boom") }
			},
			wantErr:   true,
			wantErrIn: "could not clean up properly",
		},
		{
			name: "success",
			args: []string{"song"},
			setup: func(t *testing.T) {
				tempDir = func() string { return "/tmp" }
				mkdirAll = func(path string, _ os.FileMode) error {
					assert.Equal(t, filepath.Join("/tmp", "yamp"), path)
					return nil
				}
				downloadSong = func(name, path string) error {
					assert.Equal(t, "song", name)
					assert.Equal(t, filepath.Join("/tmp", "yamp", "song.mp3"), path)
					return nil
				}
				playSong = func(path string) error {
					assert.Equal(t, filepath.Join("/tmp", "yamp", "song.mp3"), path)
					return nil
				}
				removeAll = func(path string) error {
					assert.Equal(t, filepath.Join("/tmp", "yamp"), path)
					return nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldTempDir := tempDir
			oldMkdir := mkdirAll
			oldRemoveAll := removeAll
			oldDownload := downloadSong
			oldPlay := playSong
			t.Cleanup(func() {
				tempDir = oldTempDir
				mkdirAll = oldMkdir
				removeAll = oldRemoveAll
				downloadSong = oldDownload
				playSong = oldPlay
			})

			if tt.setup != nil {
				tt.setup(t)
			}

			err := playCmd.RunE(playCmd, tt.args)
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
