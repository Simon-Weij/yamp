package main

// TODO: tests

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type MusicbrainzReleaseGroupResponse struct {
	ReleaseGroups []struct {
		ID           string `json:"id"`
		Title        string `json:"title"`
		ArtistCredit []struct {
			Artist struct {
				Name string `json:"name"`
			} `json:"artist"`
		} `json:"artist-credit"`
	} `json:"release-groups"`
}

type MusicService struct{}

func NewMusicService() *MusicService {
	return &MusicService{}
}

func (ms *MusicService) fetchReleaseGroups(artist, album string) ([]struct {
	ID string `json:"id"`
}, error) {
	query := url.QueryEscape(artist+" "+album) + "+AND+primarytype:Album+AND+secondarytype:(-Live)+AND+secondarytype:(-Compilation)+AND+secondarytype:(-Soundtrack)+AND+status:Official"
	endpoint := "https://musicbrainz.org/ws/2/release-group?query=" + query + "&fmt=json&limit=1"

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "yamp/0.1")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}
	var result struct {
		ReleaseGroups []struct {
			ID string `json:"id"`
		} `json:"release-groups"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.ReleaseGroups, nil
}

func (ms *MusicService) fetchAlbumID(artist, album string) (string, error) {
	groups, err := ms.fetchReleaseGroups(artist, album)
	if err != nil {
		return "", err
	}
	if len(groups) == 0 {
		return "", fmt.Errorf("no results found for %s - %s", artist, album)
	}
	return groups[0].ID, nil
}

func (ms *MusicService) GetAlbumCover(artist, album string) (string, error) {
	cacheBase, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	coverDir := filepath.Join(cacheBase, "yamp")
	if err := os.MkdirAll(coverDir, 0755); err != nil {
		return "", err
	}

	cacheName := artist + " - " + album
	exactPath := filepath.Join(coverDir, cacheName)
	if _, err := os.Stat(exactPath); err == nil {
		return exactPath, nil
	}

	matches, _ := filepath.Glob(exactPath + ".*")
	if len(matches) > 0 {
		return matches[0], nil
	}

	id, err := ms.fetchAlbumID(artist, album)
	if err != nil {
		return "", err
	}

	destPath := filepath.Join(coverDir, cacheName)
	if err := ms.downloadAlbumCover(destPath, id); err != nil {
		return "", err
	}
	return destPath, nil
}

func (ms *MusicService) GetAlbumCoverBase64(artist, album string) (string, error) {
	path, err := ms.GetAlbumCover(artist, album)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	mime := http.DetectContentType(data)
	if mime == "application/octet-stream" {
		mime = "image/jpeg"
	}

	b64 := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mime, b64), nil
}

func (ms *MusicService) downloadAlbumCover(location, albumID string) error {
	client := &http.Client{Timeout: 10 * time.Second}
	url := "https://coverartarchive.org/release-group/" + albumID + "/front"

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("could not download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	file, err := os.Create(location)
	if err != nil {
		return fmt.Errorf("could not make file %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
