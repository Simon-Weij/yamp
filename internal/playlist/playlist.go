package playlist

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"yamp/internal/musicdiscovery"

	"github.com/adrg/xdg"
	"github.com/spf13/afero"
)

var playlistFs afero.Fs = afero.NewOsFs()

func CreatePlaylist(playlistName string, isInternal bool) error {
	var blacklistedNames = map[string]bool{
		"":         true,
		".":        true,
		"..":       true,
		"internal": true,
	}

	if blacklistedNames[playlistName] {
		return fmt.Errorf("playlist name is not allowed")
	}
	wantPlaylist := false
	playlistFile, err := playlistSetup(playlistName, wantPlaylist, isInternal)
	if err != nil {
		return err
	}

	dir := filepath.Dir(playlistFile)
	if err := playlistFs.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := afero.WriteFile(playlistFs, playlistFile, []byte("#EXTM3U\n"), 0644); err != nil {
		return fmt.Errorf("failed to create playlist: %w", err)
	}

	return nil
}

func playlistSetup(playlistName string, wantPlaylist, isInternal bool) (string, error) {
	playlistFile := ""
	if !isInternal {
		playlistFile = filepath.Join(xdg.UserDirs.Music, "playlists", playlistName+".m3u")
	} else {
		playlistFile = filepath.Join(xdg.UserDirs.Music, "playlists", "internal", playlistName+".m3u")
	}

	playlistExists, err := PlaylistExists(playlistName, isInternal)
	if err != nil {
		return "", err
	}

	if !wantPlaylist && playlistExists {
		return "", fmt.Errorf("playlist %s already exists", playlistName)
	}
	if wantPlaylist && !playlistExists {
		return "", fmt.Errorf("playlist %s doesn't exist", playlistName)
	}

	return playlistFile, nil
}

func ListPlaylistItems(playlistName string, isInternal bool) ([]musicdiscovery.Metadata, error) {
	wantPlaylist := true
	playlistFile, err := playlistSetup(playlistName, wantPlaylist, isInternal)
	if err != nil {
		return nil, err
	}

	file, err := playlistFs.Open(playlistFile)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s: %w", playlistFile, err)
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)

	songsMetadata := []musicdiscovery.Metadata{}
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#EXTINF:-1,") {
			rest := strings.TrimPrefix(scanner.Text(), "#EXTINF:-1,")
			parts := strings.Split(rest, " - ")
			artist := parts[0]
			title := parts[1]
			metadata := musicdiscovery.Metadata{
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

func AddItemToPlaylist(playlistName, artist, title, location string, isInternal bool) error {
	wantPlaylist := true
	playlistFile, err := playlistSetup(playlistName, wantPlaylist, isInternal)
	if err != nil {
		return err
	}

	file, err := playlistFs.OpenFile(playlistFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
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

func RemoveItemFromPlaylist(playlistName, artist, title string) error {
	wantPlaylist := true
	isInternal := false
	playlistFile, err := playlistSetup(playlistName, wantPlaylist, isInternal)
	if err != nil {
		return err
	}

	file, err := playlistFs.Open(playlistFile)
	if err != nil {
		return fmt.Errorf("could not open file %s: %w", playlistFile, err)
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	removed := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#EXTINF:-1,") {
			rest := strings.TrimPrefix(line, "#EXTINF:-1,")
			parts := strings.Split(rest, " - ")
			itemArtist := ""
			itemTitle := ""
			if len(parts) >= 2 {
				itemArtist = parts[0]
				itemTitle = parts[1]
			}

			if strings.EqualFold(itemArtist, artist) && strings.EqualFold(itemTitle, title) {
				if scanner.Scan() {
					removed = true
					continue
				}
				removed = true
				break
			}
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error occurred in scanner: %w", err)
	}

	if !removed {
		return fmt.Errorf("song not found in playlist %s", playlistName)
	}

	output := strings.Join(lines, "\n")
	if !strings.HasSuffix(output, "\n") {
		output += "\n"
	}
	if err := afero.WriteFile(playlistFs, playlistFile, []byte(output), 0644); err != nil {
		return fmt.Errorf("could not write to file %s: %w", playlistFile, err)
	}

	return nil
}

func PlaylistExists(playlistName string, isInternal bool) (bool, error) {
	playlistDir := filepath.Join(xdg.UserDirs.Music, "playlists")
	if isInternal {
		playlistDir = filepath.Join(playlistDir, "internal")
	}
	playlistFile := filepath.Join(playlistDir, playlistName+".m3u")

	if _, err := playlistFs.Stat(playlistFile); err == nil {
		return true, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return false, fmt.Errorf("failed to check playlist: %w", err)
	}
	return false, nil
}

func ListPlaylists() ([]string, error) {
	playlistsDir := filepath.Join(xdg.UserDirs.Music, "playlists")
	entries, err := afero.ReadDir(playlistFs, playlistsDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read playlists dir: %w", err)
	}

	playlists := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		playlists = append(playlists, entry.Name())
	}

	return playlists, nil
}

func DeletePlaylist(playlistName string) error {
	wantPlaylist := true
	isInternal := false
	playlistFile, err := playlistSetup(playlistName, wantPlaylist, isInternal)
	if err != nil {
		return err
	}

	if err := playlistFs.Remove(playlistFile); err != nil {
		return fmt.Errorf("failed to delete playlist %s: %w", playlistName, err)
	}

	return nil
}

func ConvertSongMetadataToFilePath(artist, album, songName string) string {
	return filepath.Join(xdg.UserDirs.Music, "yamp", artist, album, songName+".mp3")
}
