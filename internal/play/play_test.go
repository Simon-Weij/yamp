package play

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/lrstanley/go-ytdlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownloadSong(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T, output string)
		output    func(t *testing.T) string
		wantErr   bool
		wantErrIn string
	}{
		{
			name: "existing file skips download",
			output: func(t *testing.T) string {
				dir := t.TempDir()
				path := filepath.Join(dir, "song.mp3")
				err := os.WriteFile(path, []byte("data"), 0644)
				require.NoError(t, err)
				return path
			},
			setup: func(t *testing.T, _ string) {
				oldInstall := ytdlpInstall
				oldInstallFFmpeg := ytdlpInstallFFmpeg
				oldDownload := ytdlpDownload
				t.Cleanup(func() {
					ytdlpInstall = oldInstall
					ytdlpInstallFFmpeg = oldInstallFFmpeg
					ytdlpDownload = oldDownload
				})
				ytdlpInstall = func() {
					t.Fatal("unexpected ytdlp install")
				}
				ytdlpInstallFFmpeg = func() {
					t.Fatal("unexpected ffmpeg install")
				}
				ytdlpDownload = func(string, string) (*ytdlp.Result, error) {
					return nil, errors.New("unexpected download")
				}
			},
		},
		{
			name: "output path is directory",
			output: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr:   true,
			wantErrIn: "directory",
		},
		{
			name: "download error returns wrapped error",
			output: func(t *testing.T) string {
				dir := t.TempDir()
				return filepath.Join(dir, "missing.mp3")
			},
			setup: func(t *testing.T, _ string) {
				oldInstall := ytdlpInstall
				oldInstallFFmpeg := ytdlpInstallFFmpeg
				oldDownload := ytdlpDownload
				t.Cleanup(func() {
					ytdlpInstall = oldInstall
					ytdlpInstallFFmpeg = oldInstallFFmpeg
					ytdlpDownload = oldDownload
				})
				ytdlpInstall = func() {}
				ytdlpInstallFFmpeg = func() {}
				ytdlpDownload = func(string, string) (*ytdlp.Result, error) {
					return &ytdlp.Result{Stdout: "log"}, errors.New("boom")
				}
			},
			wantErr:   true,
			wantErrIn: "something went wrong",
		},
		{
			name: "download success",
			output: func(t *testing.T) string {
				dir := t.TempDir()
				return filepath.Join(dir, "song.mp3")
			},
			setup: func(t *testing.T, _ string) {
				oldInstall := ytdlpInstall
				oldInstallFFmpeg := ytdlpInstallFFmpeg
				oldDownload := ytdlpDownload
				t.Cleanup(func() {
					ytdlpInstall = oldInstall
					ytdlpInstallFFmpeg = oldInstallFFmpeg
					ytdlpDownload = oldDownload
				})
				ytdlpInstall = func() {}
				ytdlpInstallFFmpeg = func() {}
				ytdlpDownload = func(string, string) (*ytdlp.Result, error) {
					return &ytdlp.Result{Stdout: ""}, nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.output(t)
			if tt.setup != nil {
				tt.setup(t, output)
			}

			err := DownloadSong("song name", output)
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

func TestPlaySong(t *testing.T) {
	tests := []struct {
		name      string
		execCmd   func(name string, args ...string) *exec.Cmd
		wantErr  bool
		wantName string
		wantArgs []string
	}{
		{
			name: "success",
			execCmd: func(name string, args ...string) *exec.Cmd {
				return exec.Command("true")
			},
			wantName: "mpv",
			wantArgs: []string{"/path/song.mp3"},
		},
		{
			name: "failure",
			execCmd: func(name string, args ...string) *exec.Cmd {
				return exec.Command("false")
			},
			wantErr:  true,
			wantName: "mpv",
			wantArgs: []string{"/path/song.mp3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldExec := execCommand
			var gotName string
			var gotArgs []string
			execCommand = func(name string, args ...string) *exec.Cmd {
				gotName = name
				gotArgs = append([]string{}, args...)
				return tt.execCmd(name, args...)
			}
			t.Cleanup(func() {
				execCommand = oldExec
			})

			err := PlaySong("/path/song.mp3")
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.wantName, gotName)
			assert.Equal(t, tt.wantArgs, gotArgs)
		})
	}
}
