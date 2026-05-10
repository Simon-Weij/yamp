package musicdiscovery

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lrstanley/go-ytdlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestExtractMetadata(t *testing.T) {
	tests := []struct {
		name      string
		response  string
		status    int
		want      *Metadata
		wantErr   bool
		wantErrIn string
		validate  func(t *testing.T, r *http.Request)
	}{
		{
			name:     "success",
			response: `{"recordings":[{"title":"Song","artist-credit":[{"artist":{"name":"Artist"}}],"releases":[{"release-group":{"title":"Album"}}]}]}`,
			status:   200,
			want:     &Metadata{Artist: "Artist", Album: "Album", Title: "Song"},
			validate: func(t *testing.T, r *http.Request) {
				assert.Equal(t, "yamp/0.1", r.Header.Get("User-Agent"))
			},
		},
		{
			name:      "no results",
			response:  `{"recordings":[]}`,
			status:    200,
			wantErr:   true,
			wantErrIn: "no results found",
		},
		{
			name:      "incomplete metadata",
			response:  `{"recordings":[{"title":"Song","artist-credit":[],"releases":[]}]}`,
			status:    200,
			wantErr:   true,
			wantErrIn: "incomplete metadata",
		},
		{
			name:      "invalid json",
			response:  `not-json`,
			status:    200,
			wantErr:   true,
			wantErrIn: "invalid character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.validate != nil {
					tt.validate(t, r)
				}
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer srv.Close()

			oldBaseURL := musicbrainzBaseURL
			oldClient := httpClient
			t.Cleanup(func() {
				musicbrainzBaseURL = oldBaseURL
				httpClient = oldClient
			})
			musicbrainzBaseURL = srv.URL
			httpClient = srv.Client()

			got, err := ExtractMetadata("Song Name")
			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrIn != "" {
					assert.Contains(t, err.Error(), tt.wantErrIn)
				}
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExtractMetadataHTTPError(t *testing.T) {
	oldClient := httpClient
	oldBaseURL := musicbrainzBaseURL
	t.Cleanup(func() {
		httpClient = oldClient
		musicbrainzBaseURL = oldBaseURL
	})

	httpClient = &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("network down")
	})}
	musicbrainzBaseURL = "http://example.invalid"

	_, err := ExtractMetadata("Song Name")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "network down")
}

func TestGetSimilarSongs(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T)
		wantErr   bool
		wantErrIn string
	}{
		{
			name: "search error",
			setup: func(t *testing.T) {
				oldInstall := ytdlpInstall
				oldInstallBun := ytdlpInstallBun
				oldInstallFFmpeg := ytdlpInstallFFmpeg
				oldSearch := ytdlpSearchID
				oldGetSimilar := ytdlpGetSimilar
				t.Cleanup(func() {
					ytdlpInstall = oldInstall
					ytdlpInstallBun = oldInstallBun
					ytdlpInstallFFmpeg = oldInstallFFmpeg
					ytdlpSearchID = oldSearch
					ytdlpGetSimilar = oldGetSimilar
				})
				ytdlpInstall = func() {}
				ytdlpInstallBun = func() {}
				ytdlpInstallFFmpeg = func() {}
				ytdlpSearchID = func(string) (*ytdlp.Result, error) {
					return nil, errors.New("search failed")
				}
				ytdlpGetSimilar = func(string) (*ytdlp.Result, error) {
					return nil, errors.New("unexpected get")
				}
			},
			wantErr:   true,
			wantErrIn: "could not search",
		},
		{
			name: "empty id",
			setup: func(t *testing.T) {
				oldInstall := ytdlpInstall
				oldInstallBun := ytdlpInstallBun
				oldInstallFFmpeg := ytdlpInstallFFmpeg
				oldSearch := ytdlpSearchID
				oldGetSimilar := ytdlpGetSimilar
				t.Cleanup(func() {
					ytdlpInstall = oldInstall
					ytdlpInstallBun = oldInstallBun
					ytdlpInstallFFmpeg = oldInstallFFmpeg
					ytdlpSearchID = oldSearch
					ytdlpGetSimilar = oldGetSimilar
				})
				ytdlpInstall = func() {}
				ytdlpInstallBun = func() {}
				ytdlpInstallFFmpeg = func() {}
				ytdlpSearchID = func(string) (*ytdlp.Result, error) {
					return &ytdlp.Result{Stdout: "\n"}, nil
				}
				ytdlpGetSimilar = func(string) (*ytdlp.Result, error) {
					return nil, errors.New("unexpected get")
				}
			},
			wantErr:   true,
			wantErrIn: "empty video id",
		},
		{
			name: "similar error",
			setup: func(t *testing.T) {
				oldInstall := ytdlpInstall
				oldInstallBun := ytdlpInstallBun
				oldInstallFFmpeg := ytdlpInstallFFmpeg
				oldSearch := ytdlpSearchID
				oldGetSimilar := ytdlpGetSimilar
				t.Cleanup(func() {
					ytdlpInstall = oldInstall
					ytdlpInstallBun = oldInstallBun
					ytdlpInstallFFmpeg = oldInstallFFmpeg
					ytdlpSearchID = oldSearch
					ytdlpGetSimilar = oldGetSimilar
				})
				ytdlpInstall = func() {}
				ytdlpInstallBun = func() {}
				ytdlpInstallFFmpeg = func() {}
				ytdlpSearchID = func(string) (*ytdlp.Result, error) {
					return &ytdlp.Result{Stdout: "abc123"}, nil
				}
				ytdlpGetSimilar = func(string) (*ytdlp.Result, error) {
					return nil, errors.New("get failed")
				}
			},
			wantErr:   true,
			wantErrIn: "could not fetch similar",
		},
		{
			name: "success",
			setup: func(t *testing.T) {
				oldInstall := ytdlpInstall
				oldInstallBun := ytdlpInstallBun
				oldInstallFFmpeg := ytdlpInstallFFmpeg
				oldSearch := ytdlpSearchID
				oldGetSimilar := ytdlpGetSimilar
				t.Cleanup(func() {
					ytdlpInstall = oldInstall
					ytdlpInstallBun = oldInstallBun
					ytdlpInstallFFmpeg = oldInstallFFmpeg
					ytdlpSearchID = oldSearch
					ytdlpGetSimilar = oldGetSimilar
				})
				ytdlpInstall = func() {}
				ytdlpInstallBun = func() {}
				ytdlpInstallFFmpeg = func() {}
				ytdlpSearchID = func(string) (*ytdlp.Result, error) {
					return &ytdlp.Result{Stdout: "abc123"}, nil
				}
				ytdlpGetSimilar = func(string) (*ytdlp.Result, error) {
					return &ytdlp.Result{Stdout: "song list"}, nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t)
			}

			result, err := GetSimilarSongs("name")
			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrIn != "" {
					assert.Contains(t, err.Error(), tt.wantErrIn)
				}
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, result)
		})
	}
}

func TestExtractMetadataHTTPStatusNonOK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "{}")
	}))
	defer srv.Close()

	oldBaseURL := musicbrainzBaseURL
	oldClient := httpClient
	t.Cleanup(func() {
		musicbrainzBaseURL = oldBaseURL
		httpClient = oldClient
	})
	musicbrainzBaseURL = srv.URL
	httpClient = srv.Client()

	_, err := ExtractMetadata("Song Name")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected response status")
}
