package main

import (
	"bufio"
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
	Artist string
	Album  string
	Title  string
}

func NewPlaylistRepository(fs afero.Fs) *PlaylistRepository {
	return &PlaylistRepository{
		Fs: fs,
	}
}

func (pr *PlaylistRepository) ListPlaylists() ([]string, error) {
	dir := filepath.Join(xdg.UserDirs.Music, "playlists")

	items, err := afero.ReadDir(pr.Fs, dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("could not read playlists directory: %w", err)
	}

	if len(items) == 0 {
		return []string{}, nil
	}

	var names []string
	for _, file := range items {
		name := strings.TrimSuffix(file.Name(), ".m3u")
		names = append(names, name)
	}

	return names, nil
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

func (pr *PlaylistRepository) createPlaylist(name string) (string, error) {
	if err := validateFilename(name); err != nil {
		return "", err
	}
	path := filepath.Join(xdg.UserDirs.Music, "playlists", name+".m3u")

	if _, err := pr.Fs.Stat(path); err == nil {
		return "", fmt.Errorf("file already exists")
	}
	data := []byte("#EXTM3U\n")
	if err := afero.WriteFile(pr.Fs, path, data, 0644); err != nil {
		return "", fmt.Errorf("could not write to playlist file: %w", err)
	}
	return path, nil
}

func (pr *PlaylistRepository) addSongToPlaylist(playlistLocation string, playlistItem PlaylistItem, songLocation string) error {
	file, err := pr.Fs.OpenFile(playlistLocation, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("could not open playlist: %w", err)
	}
	defer file.Close()

	entry := fmt.Sprintf("#EXTINF:-1,%s - %s - %s\n%s\n",
		playlistItem.Artist,
		playlistItem.Album,
		playlistItem.Title,
		songLocation,
	)

	if _, err := file.WriteString(entry); err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}
	return nil
}

func (pr *PlaylistRepository) ParsePlaylistFile(name string) (*[]PlaylistItem, error) {
	path := filepath.Join(xdg.UserDirs.Music, "playlists", name+".m3u")
	file, err := pr.Fs.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file %w", err)
	}
	defer file.Close()
	re := regexp.MustCompile(`:-?\d+,`)

	var playlistItems []PlaylistItem

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		expectedPrefix := "#EXTINF"
		if !strings.HasPrefix(line, expectedPrefix) {
			continue
		}

		line = strings.TrimPrefix(line, expectedPrefix)
		line = re.ReplaceAllString(line, "")

		parts := strings.Split(line, " - ")

		if len(parts) != 3 {
			return nil, fmt.Errorf("unexpected format")
		}

		playlistItem := PlaylistItem{
			Artist: strings.TrimSpace(parts[0]),
			Album:  strings.TrimSpace(parts[1]),
			Title:  strings.TrimSpace(parts[2]),
		}

		playlistItems = append(playlistItems, playlistItem)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	return &playlistItems, nil
}
