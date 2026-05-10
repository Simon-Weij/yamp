package musicdiscovery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/lrstanley/go-ytdlp"
)

type MusicbrainzResponse struct {
	Recordings []struct {
		Title        string `json:"title"`
		ArtistCredit []struct {
			Artist struct {
				Name string `json:"name"`
			} `json:"artist"`
		} `json:"artist-credit"`
		Releases []struct {
			ReleaseGroup struct {
				Title string `json:"title"`
			} `json:"release-group"`
		} `json:"releases"`
	} `json:"recordings"`
}

type Metadata struct {
	Artist string
	Album  string
	Title  string
}

func ExtractMetadata(name string) (*Metadata, error) {
	query := url.QueryEscape(name) + "+AND+primarytype:Album+AND+secondarytype:(-Live)+AND+secondarytype:(-Compilation)+AND+secondarytype:(-Soundtrack)+AND+status:Official"
	endpoint := musicbrainzBaseURL + "?query=" + query + "&fmt=json&limit=1"

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "yamp/0.1")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	var result MusicbrainzResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Recordings) == 0 {
		return nil, fmt.Errorf("no results found for %q", name)
	}

	rec := result.Recordings[0]

	if len(rec.ArtistCredit) == 0 || len(rec.Releases) == 0 {
		return nil, fmt.Errorf("incomplete metadata for %q", name)
	}

	return &Metadata{
		Artist: rec.ArtistCredit[0].Artist.Name,
		Title:  rec.Title,
		Album:  rec.Releases[0].ReleaseGroup.Title,
	}, nil
}

func GetSimilarSongs(name string) (*ytdlp.Result, error) {
	ytdlpInstall()
	ytdlpInstallBun()
	ytdlpInstallFFmpeg()

	result, err := ytdlpSearchID(name)
	if err != nil {
		return nil, fmt.Errorf("could not search for song %s: %w", name, err)
	}
	id := strings.TrimSpace(result.Stdout)
	if id == "" {
		return nil, fmt.Errorf("could not search for song %s: empty video id", name)
	}

	url := fmt.Sprintf("https://www.youtube.com/watch?v=%s&list=RD%s", id, id)
	songs, err := ytdlpGetSimilar(url)
	if err != nil {
		return nil, fmt.Errorf("could not fetch similar songs for %s: %w", name, err)
	}

	return songs, nil
}

var (
	musicbrainzBaseURL = "https://musicbrainz.org/ws/2/recording/"
	httpClient         = http.DefaultClient
	ytdlpInstall       = func() {
		ytdlp.MustInstall(context.TODO(), nil)
	}
	ytdlpInstallBun = func() {
		ytdlp.MustInstallBun(context.TODO(), nil)
	}
	ytdlpInstallFFmpeg = func() {
		ytdlp.MustInstallFFmpeg(context.TODO(), nil)
	}
	ytdlpSearchID = func(name string) (*ytdlp.Result, error) {
		return ytdlp.New().Print("id").Run(context.TODO(), "ytsearch1:"+name)
	}
	ytdlpGetSimilar = func(url string) (*ytdlp.Result, error) {
		return ytdlp.New().Print("title").FlatPlaylist().PlaylistItems("2-11").Run(context.TODO(), url)
	}
)
