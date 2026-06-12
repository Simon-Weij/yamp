package main

import (
	"fmt"
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

func Test_addSongToPlaylist(t *testing.T) {
	type fields struct {
		Fs afero.Fs
	}

	type args struct {
		playlistLocation string
		playlistItem     PlaylistItem
		songLocation     string
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		wantErr      bool
		playlistName string
	}{
		{
			"function should run successfully",
			fields{
				Fs: afero.NewMemMapFs(),
			},
			args{
				playlistLocation: "/home/someuser/playlist.mp3",

				playlistItem: PlaylistItem{
					Artist: "test artist",
					Album:  "test album",
					Title:  "test title",
				},
				songLocation: "/home/someuser/song.mp3",
			},
			false, "test_playlist",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := NewPlaylistRepository(tt.fields.Fs)
			path, err := pr.createPlaylist(tt.playlistName)
			require.NoError(t, err)

			err = pr.addSongToPlaylist(path, tt.args.playlistItem, tt.args.songLocation)
			if !tt.wantErr {
				require.NoError(t, err)
			}
			data, err := afero.ReadFile(tt.fields.Fs, path)
			require.NoError(t, err)
			contents := string(data)
			expectedFile := fmt.Sprintf("#EXTM3U\n#EXTINF:-1,%s - %s - %s\n%s\n",
				tt.args.playlistItem.Artist,
				tt.args.playlistItem.Album,
				tt.args.playlistItem.Title,
				tt.args.songLocation,
			)
			assert.Equal(t, expectedFile, contents)
		})
	}
}

func Test_ParsePlaylistFile(t *testing.T) {
	type fields struct {
		Fs afero.Fs
	}
	type args struct {
		path string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		playlistItem *[]PlaylistItem
		wantErr      bool
		songLocation string
	}{
		{
			"function should run without issues", fields{
				Fs: afero.NewMemMapFs(),
			}, args{
				path: "/home/someuser/testplaylist.m3u",
			}, &[]PlaylistItem{
				{Artist: "test artist 1", Album: "test album 1", Title: "test title 2"},
				{Artist: "test artist 2", Album: "test album 2", Title: "test title 3"},
			}, false, "/home/someuser/song.mp3",
		}, // TODO: more test cases
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := NewPlaylistRepository(afero.NewMemMapFs())
			path, err := pr.createPlaylist("test")
			require.NoError(t, err)
			for _, v := range *tt.playlistItem {
				err := pr.addSongToPlaylist(path, v, tt.songLocation)
				require.NoError(t, err)
			}
			playlistItems, err := pr.ParsePlaylistFile(path)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.playlistItem, playlistItems)
			}
		})
	}
}
