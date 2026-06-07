package main

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_createPlaylist(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		wantPlaylist bool
	}{
		{"create playlist", args{
			name: "testplaylist",
		}, false, true},
		{"create playlist with invalid name", args{
			name: "../../playlist",
		}, true, false},
		{"create playlist with an empty name", args{
			name: "",
		}, true, false},
		{"create playlist with just dots", args{
			name: "...",
		}, true, false},
		{"create a name that's too long", args{
			name: strings.Repeat("a", 999),
		}, true, false},
		{"create a playlist with spaces", args{
			name: " test_playlist ",
		}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			playlistRepo := &PlaylistRepository{
				Fs: afero.NewMemMapFs(),
			}
			path, err := playlistRepo.createPlaylist(tt.args.name)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			if tt.wantPlaylist {
				exists, err := afero.Exists(playlistRepo.Fs, path)
				assert.NoError(t, err)
				require.True(t, exists)
			} else {
				exists, err := afero.Exists(playlistRepo.Fs, path)
				assert.NoError(t, err)
				assert.False(t, exists)
			}
			_, err = playlistRepo.createPlaylist(tt.args.name)
			require.Error(t, err)
		})
	}
}
func TestPlaylistRepository_ListPlaylists(t *testing.T) {
	type fields struct {
		Fs afero.Fs
	}

	tests := []struct {
		name             string
		fields           fields
		createdPlaylists []string
		want             []string
		wantErr          bool
	}{
		{
			name:   "create playlists, then should appear in ListPlaylists",
			fields: fields{Fs: afero.NewMemMapFs()},
			createdPlaylists: []string{
				"playlist 1",
				"playlist 2",
				"playlist 3",
			},
			want:    []string{"playlist 1", "playlist 2", "playlist 3"},
			wantErr: false,
		},
		{
			name:             "when no playlists, should not error",
			fields:           fields{Fs: afero.NewMemMapFs()},
			createdPlaylists: []string{},
			want:             []string{},
			wantErr:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := &PlaylistRepository{
				Fs: tt.fields.Fs,
			}

			for _, playlist := range tt.createdPlaylists {
				_, err := pr.createPlaylist(playlist)
				require.NoError(t, err)
			}

			got, err := pr.ListPlaylists()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.ElementsMatch(t, tt.want, got)
		})
	}
}
