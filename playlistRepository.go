package main

import (
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
