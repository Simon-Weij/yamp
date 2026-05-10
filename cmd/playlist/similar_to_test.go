package playlistcmd

import (
	"errors"
	"testing"

	"yamp/internal/musicdiscovery"

	"github.com/lrstanley/go-ytdlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimilarToCommand(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T)
		args      []string
		wantErr   bool
		wantErrIn string
	}{
		{
			name: "get similar error",
			args: []string{"song"},
			setup: func(t *testing.T) {
				getSimilarSongsFn = func(string) (*ytdlp.Result, error) { return nil, errors.New("boom") }
			},
			wantErr:   true,
			wantErrIn: "could not find similar songs",
		},
		{
			name: "add songs error",
			args: []string{"song"},
			setup: func(t *testing.T) {
				getSimilarSongsFn = func(string) (*ytdlp.Result, error) { return &ytdlp.Result{Stdout: "one\n"}, nil }
				addSongsToPlaylistFn = func([]string, string) error { return errors.New("boom") }
			},
			wantErr:   true,
			wantErrIn: "could not add songs to playlist",
		},
		{
			name: "success",
			args: []string{"song"},
			setup: func(t *testing.T) {
				getSimilarSongsFn = func(string) (*ytdlp.Result, error) {
					return &ytdlp.Result{Stdout: "one\n\n two\n"}, nil
				}
				addSongsToPlaylistFn = func(lines []string, playlistName string) error {
					assert.Equal(t, []string{"one", " two"}, lines)
					assert.Equal(t, "similar-to-song", playlistName)
					return nil
				}
			},
		},
		{
			name: "empty stdout",
			args: []string{"song"},
			setup: func(t *testing.T) {
				getSimilarSongsFn = func(string) (*ytdlp.Result, error) { return &ytdlp.Result{Stdout: "\n"}, nil }
				addSongsToPlaylistFn = func(lines []string, playlistName string) error {
					assert.Empty(t, lines)
					assert.Equal(t, "similar-to-song", playlistName)
					return nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldGetSimilar := getSimilarSongsFn
			oldAddSongs := addSongsToPlaylistFn
			t.Cleanup(func() {
				getSimilarSongsFn = oldGetSimilar
				addSongsToPlaylistFn = oldAddSongs
			})

			if tt.setup != nil {
				tt.setup(t)
			}

			err := playlistSimilarToCmd.RunE(playlistSimilarToCmd, tt.args)
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

func TestAddSongsToPlaylist(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T)
		songs     []string
		wantErr   bool
		wantErrIn string
	}{
		{
			name:  "create playlist error",
			songs: []string{"one"},
			setup: func(t *testing.T) {
				createPlaylistSimilarToFn = func(string, bool) error { return errors.New("boom") }
			},
			wantErr:   true,
			wantErrIn: "could not create playlist",
		},
		{
			name:  "extract metadata error",
			songs: []string{"one"},
			setup: func(t *testing.T) {
				createPlaylistSimilarToFn = func(string, bool) error { return nil }
				extractMetadataFn = func(string) (*musicdiscovery.Metadata, error) { return nil, errors.New("boom") }
			},
			wantErr:   true,
			wantErrIn: "could not get metadata",
		},
		{
			name:  "download error",
			songs: []string{"one"},
			setup: func(t *testing.T) {
				createPlaylistSimilarToFn = func(string, bool) error { return nil }
				extractMetadataFn = func(string) (*musicdiscovery.Metadata, error) {
					return &musicdiscovery.Metadata{Artist: "a", Album: "b", Title: "t"}, nil
				}
				convertSongMetadataToFilePathFn = func(string, string, string) string { return "/tmp/song.mp3" }
				downloadSongFn = func(string, string) error { return errors.New("boom") }
			},
			wantErr:   true,
			wantErrIn: "could not download",
		},
		{
			name:  "add item error",
			songs: []string{"one"},
			setup: func(t *testing.T) {
				createPlaylistSimilarToFn = func(string, bool) error { return nil }
				extractMetadataFn = func(string) (*musicdiscovery.Metadata, error) {
					return &musicdiscovery.Metadata{Artist: "a", Album: "b", Title: "t"}, nil
				}
				convertSongMetadataToFilePathFn = func(string, string, string) string { return "/tmp/song.mp3" }
				downloadSongFn = func(string, string) error { return nil }
				addItemToPlaylistFn = func(string, string, string, string, bool) error { return errors.New("boom") }
			},
			wantErr:   true,
			wantErrIn: "could not add",
		},
		{
			name:  "success",
			songs: []string{"one"},
			setup: func(t *testing.T) {
				createPlaylistSimilarToFn = func(string, bool) error { return nil }
				extractMetadataFn = func(string) (*musicdiscovery.Metadata, error) {
					return &musicdiscovery.Metadata{Artist: "a", Album: "b", Title: "t"}, nil
				}
				convertSongMetadataToFilePathFn = func(string, string, string) string { return "/tmp/song.mp3" }
				downloadSongFn = func(string, string) error { return nil }
				addItemToPlaylistFn = func(string, string, string, string, bool) error { return nil }
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldCreate := createPlaylistSimilarToFn
			oldExtract := extractMetadataFn
			oldConvert := convertSongMetadataToFilePathFn
			oldDownload := downloadSongFn
			oldAddItem := addItemToPlaylistFn
			t.Cleanup(func() {
				createPlaylistSimilarToFn = oldCreate
				extractMetadataFn = oldExtract
				convertSongMetadataToFilePathFn = oldConvert
				downloadSongFn = oldDownload
				addItemToPlaylistFn = oldAddItem
			})

			if tt.setup != nil {
				tt.setup(t)
			}

			err := addSongsToPlaylist(tt.songs, "similar-to-song")
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

func TestCleanLines(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "strip brackets and parens",
			input: "Artist - Song (Live) [Remastered]",
			want:  "Artist - Song",
		},
		{
			name:  "strip quotes",
			input: "\"Artist\" - \"Song\"",
			want:  "Artist - Song",
		},
		{
			name:  "trim whitespace",
			input: "  Artist - Song  ",
			want:  "Artist - Song",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, cleanLines(tt.input))
		})
	}
}
