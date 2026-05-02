package playlist

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

func CreatePlaylist(playlistName string) error {
	playlistFile := filepath.Join(xdg.UserDirs.Music, "playlists", playlistName)

	playlistExists, err := PlaylistExists(playlistName)
	if err != nil {
		return err
	}

	if playlistExists {
		return fmt.Errorf("playlist %s already exists", playlistName)
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

func AddItemToPlaylist(playlistName, artist, title, location string) error {
	playlistFile := filepath.Join(xdg.UserDirs.Music, "playlists", playlistName)

	playlistExists, err := PlaylistExists(playlistName)
	if err != nil {
		return err
	}

	if !playlistExists {
		return fmt.Errorf("playlist %s doesn't exist", playlistName)
	}

	file, err := os.OpenFile(playlistFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("could not open file %s: %w", playlistFile, err)
	}
	defer func() {
		_ = file.Close()
	}()

	if _, err = fmt.Fprintf(file, "#EXTINF:-1,%s - %s\n%s\n", artist, title, location); err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}

	return nil
}

func PlaylistExists(playlistName string) (bool, error) {
	playlistFile := filepath.Join(xdg.UserDirs.Music, "playlists", playlistName)

	if _, err := os.Stat(playlistFile); err == nil {
		return true, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return false, fmt.Errorf("failed to check playlist: %w", err)
	}
	return false, nil
}
