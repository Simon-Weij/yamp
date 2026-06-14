package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBrowserRepository_SearchSong(t *testing.T) {
	type args struct {
		ctx   context.Context
		query string
	}

	songs := []Song{
		{TrackName: "Title 1", Artist: "Artist 1", CollectionName: "Album 1", Cover: "https://example.com"},
		{TrackName: "Title 2", Artist: "Artist 2", CollectionName: "Album 2", Cover: "https://example.com"},
		{TrackName: "Title 3", Artist: "Artist 3", CollectionName: "Album 3", Cover: "https://example.com"},
	}

	tests := []struct {
		name           string
		args           args
		want           []Song
		wantErr        bool
		songs          []Song
		httpStatusCode int
	}{
		{"should return without errors", args{
			context.Background(),
			"query",
		}, songs, false, songs, 200,
		},
		{"should return error when no results", args{
			context.Background(),
			"query",
		}, nil, true, []Song{}, 200},
		{"should return error when server errors", args{
			context.Background(),
			"query",
		}, nil, true, songs, 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				w.WriteHeader(tt.httpStatusCode)
				json.NewEncoder(w).Encode(searchResponse{Results: tt.songs})
			}))
			defer srv.Close()

			br := &BrowserRepository{baseURL: srv.URL}
			got, err := br.SearchSong(tt.args.ctx, tt.args.query)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
