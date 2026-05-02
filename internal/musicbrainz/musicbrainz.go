package musicbrainz

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	endpoint := "https://musicbrainz.org/ws/2/recording/?query=" + query + "&fmt=json&limit=1"

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "yamp/0.1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

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
