package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type BrowserRepository struct {
	baseURL string
}

func NewBrowserRepository() *BrowserRepository {
	return &BrowserRepository{
		baseURL: "https://itunes.apple.com",
	}
}

type Song struct {
	TrackName      string `json:"trackName"`
	Artist         string `json:"artistName"`
	CollectionName string `json:"collectionName"`
	Cover          string `json:"artworkUrl100"`
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
