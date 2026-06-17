package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrg/xdg"
	"github.com/spf13/afero"
)

type BrowserRepository struct {
	baseURL    string
	httpClient *http.Client
	Fs         afero.Fs
}

func NewBrowserRepository(fs afero.Fs) *BrowserRepository {
	return &BrowserRepository{
		baseURL:    "https://itunes.apple.com",
		httpClient: &http.Client{Timeout: 10 * time.Second},
		Fs:         fs,
	}
}

type Song struct {
	TrackName      string `json:"trackName"`
	Artist         string `json:"artistName"`
	CollectionName string `json:"collectionName"`
	Cover          string `json:"artworkUrl100"`
	DurationMillis int    `json:"trackTimeMillis"`
}

type searchResponse struct {
	Results []Song `json:"results"`
}

func (br *BrowserRepository) SearchSong(ctx context.Context, query string) ([]Song, error) {
	httpClient := &http.Client{Timeout: 10 * time.Second}

	term := strings.ReplaceAll(query, " ", "+")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/search?term=%s&media=music&entity=song&sort=popular", br.baseURL, term), nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var body searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	if len(body.Results) == 0 {
		return nil, fmt.Errorf("no results for %q", query)
	}

	return body.Results, nil
}

func (br *BrowserRepository) AddSongToPlaylist(song Song, playlist string) error {
	playlistLocation := filepath.Join(xdg.UserDirs.Music, "playlists", playlist+".m3u")
	songLocation := filepath.Join(xdg.UserDirs.Music, song.TrackName)
	file, err := os.OpenFile(playlistLocation, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("could not open playlist: %w", err)
	}
	defer file.Close()

	entry := fmt.Sprintf("#EXTINF:%d,%s - %s - %s\n%s\n",
		song.DurationMillis,
		song.Artist,
		song.CollectionName,
		song.TrackName,
		songLocation,
	)

	if _, err := file.WriteString(entry); err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}
	return nil
}

func (br *BrowserRepository) RemoveSongFromPlaylist(title, album, artist, playlist string) error {
	playlistLocation := filepath.Join(xdg.UserDirs.Music, "playlists", playlist+".m3u")

	data, err := os.ReadFile(playlistLocation)
	if err != nil {
		return fmt.Errorf("could not read playlist: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	out := make([]string, 0, len(lines))

	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], artist+" - "+album+" - "+title) {
			i++
			continue
		}
		out = append(out, lines[i])
	}

	return os.WriteFile(playlistLocation, []byte(strings.Join(out, "\n")), 0644)
}

// This exists because webkitgtk can't behave and throws tls errors when you interact with network
func (br *BrowserRepository) FetchImageAsBase64(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	b64 := base64.StdEncoding.EncodeToString(bytes)
	return "data:image/jpeg;base64," + b64, nil
}
