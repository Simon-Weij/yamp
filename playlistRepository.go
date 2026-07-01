package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/adrg/xdg"
	"github.com/spf13/afero"
)

type PlaylistRepository struct {
	Fs afero.Fs
}

type PlaylistItem struct {
	Artist   string
	Album    string
	Title    string
	Cover    string
	Duration int
}

func NewPlaylistRepository(fs afero.Fs) *PlaylistRepository {
	return &PlaylistRepository{
		Fs: fs,
	}
}

func validateFilename(name string) error {
	invalidFilename := regexp.MustCompile(`[/\\:*?"<>|]|\.\.|^\.`)
	if invalidFilename.MatchString(name) {
		return fmt.Errorf("invalid playlist name: %q", name)
	}
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("filename can't be empty")
	}
	if len(name) > 254 {
		return fmt.Errorf("name is too long")
	}
	if strings.HasPrefix(name, " ") || strings.HasSuffix(name, " ") {
		return fmt.Errorf("name can't start or end with a space")
	}
	return nil
}

func getPlaylistPath(name string) string {
	return filepath.Join(xdg.DataHome, "yamp", "playlists", name+".json")
}

func getPlaylistDir() string {
	return filepath.Join(xdg.DataHome, "yamp", "playlists")
}

func (pr *PlaylistRepository) CreatePlaylist(name string) (string, error) {
	if err := validateFilename(name); err != nil {
		return "", err
	}
	path := getPlaylistPath(name)

	if _, err := pr.Fs.Stat(path); err == nil {
		return "", fmt.Errorf("file already exists")
	}

	dir := filepath.Dir(path)
	if err := pr.Fs.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	if err := afero.WriteFile(pr.Fs, path, []byte("[]"), 0o644); err != nil {
		return "", err
	}

	return path, nil
}

func (pr *PlaylistRepository) getSongsInPlaylist(playlistName string) ([]Song, error) {
	var songs []Song
	path := getPlaylistPath(playlistName)
	if _, err := pr.Fs.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("%s does not exist", playlistName)
	}
	data, err := afero.ReadFile(pr.Fs, path)
	if err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &songs)
	}
	return songs, nil
}

func (pr *PlaylistRepository) AddSongToPlaylist(song Song, playlistName string) error {
	var songs []Song
	path := getPlaylistPath(playlistName)
	if _, err := pr.Fs.Stat(path); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%s does not exist", playlistName)
	}

	data, err := afero.ReadFile(pr.Fs, path)
	if err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &songs)
	}
	songs = append(songs, song)

	out, err := json.MarshalIndent(songs, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := pr.Fs.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	return afero.WriteFile(pr.Fs, path, out, 0o644)
}

func (pr *PlaylistRepository) RemoveSongFromPlaylist(title, artist, playlistName string) error {
	var songs []Song
	path := getPlaylistPath(playlistName)
	if _, err := pr.Fs.Stat(path); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%s does not exist", playlistName)
	}

	data, err := afero.ReadFile(pr.Fs, path)
	if err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, &songs)
	}

	filtered := songs[:0]
	for _, s := range songs {
		if s.TrackName != title || s.Artist != artist {
			filtered = append(filtered, s)
		}
	}

	if len(filtered) == len(songs) {
		return fmt.Errorf("song not found in %s", playlistName)
	}

	out, err := json.MarshalIndent(filtered, "", "  ")
	if err != nil {
		return err
	}

	return afero.WriteFile(pr.Fs, path, out, 0o644)
}

func (pr *PlaylistRepository) ListSongsInPlaylist(playlistName string) ([]PlaylistItem, error) {
	songs, err := pr.getSongsInPlaylist(playlistName)
	if err != nil {
		return nil, err
	}
	var playlistItems []PlaylistItem
	for _, song := range songs {
		playlistItems = append(playlistItems, PlaylistItem{
			Artist:   song.Artist,
			Album:    song.CollectionName,
			Title:    song.TrackName,
			Cover:    song.Cover,
			Duration: song.DurationMillis,
		})
	}
	return playlistItems, nil
}

func (pr *PlaylistRepository) ListPlaylists() []string {
	dir := getPlaylistDir()
	entries, err := afero.ReadDir(pr.Fs, dir)
	if err != nil {
		return []string{}
	}

	var playlists []string
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".json") {
			playlists = append(playlists, strings.TrimSuffix(entry.Name(), ".json"))
		}
	}
	return playlists
}
