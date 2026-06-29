package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

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
	TrackID      int64 `json:"trackId"`
	ArtistID     int64 `json:"artistId"`
	CollectionID int64 `json:"collectionId"`

	TrackName      string `json:"trackName"`
	Artist         string `json:"artistName"`
	CollectionName string `json:"collectionName"`

	Cover          string `json:"artworkUrl100"`
	DurationMillis int    `json:"trackTimeMillis"`

	PreviewURL string `json:"previewUrl,omitempty"`
	TrackURL   string `json:"trackViewUrl,omitempty"`
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
