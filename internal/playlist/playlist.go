package playlist

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"yamp/internal/musicbrainz"

	"github.com/adrg/xdg"
)

func CreatePlaylist(playlistName string) error {
	wantPlaylist := false
	playlistFile, err := playlistSetup(playlistName, wantPlaylist)
	if err != nil {
		return err
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

func playlistSetup(playlistName string, wantPlaylist bool) (string, error) {
	playlistFile := filepath.Join(xdg.UserDirs.Music, "playlists", playlistName)

	playlistExists, err := PlaylistExists(playlistName)
	if err != nil {
		return "", err
	}

	if !wantPlaylist && playlistExists {
		return "", fmt.Errorf("could not find playlist %s: %w", playlistName, err)
	}
	if wantPlaylist && !playlistExists {
		return "", fmt.Errorf("playlist %s doesn't exist", playlistName)
	}

	return playlistFile, nil
}

func ListPlaylistItems(playlistName string) ([]musicbrainz.Metadata, error) {
	wantPlaylist := false
	playlistFile, err := playlistSetup(playlistName, wantPlaylist)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(playlistFile)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s: %w", playlistFile, err)
	}

	scanner := bufio.NewScanner(file)

	songsMetadata := []musicbrainz.Metadata{}
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#EXTINF:-1,") {
			rest := strings.TrimPrefix(scanner.Text(), "#EXTINF:-1,")
			parts := strings.Split(rest, " - ")
			artist := parts[0]
			title := parts[1]
			metadata := musicbrainz.Metadata{
				Artist: artist,
				Title:  title,
			}
			songsMetadata = append(songsMetadata, metadata)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error occurred in scanner: %w", err)
	}

	if len(songsMetadata) == 0 {
		return nil, fmt.Errorf("no songs found in playlist %s", playlistName)
	}

	return songsMetadata, nil
}

func AddItemToPlaylist(playlistName, artist, title, location string) error {
	wantPlaylist := true
	playlistFile, err := playlistSetup(playlistName, wantPlaylist)
	if err != nil {
		return err
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
