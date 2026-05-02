package playlist

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

type Song struct {
	Title  string
	Artist string
	Path   string
}

func CreatePlaylist(playlistName string) error {
	playlistFile := filepath.Join(xdg.UserDirs.Music, "playlists", playlistName)

	if _, err := os.Stat(playlistFile); err == nil {
		return fmt.Errorf("playlist already exists")
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to check playlist: %w", err)
	}

	dir := filepath.Dir(playlistFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(playlistFile, []byte("#EXTM3U\n"), 0644); err != nil {
		return fmt.Errorf("failed to create playlist: %w", err)
	}

	return nil
}
