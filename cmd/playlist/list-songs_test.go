package playlistcmd

import (
	"errors"
	"testing"

	"yamp/internal/musicdiscovery"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListSongsCommand(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T)
		args      []string
		wantErr   bool
		wantErrIn string
	}{
		{
			name: "list error",
			args: []string{"playlist"},
			setup: func(t *testing.T) {
				listPlaylistItemsFn = func(string, bool) ([]musicdiscovery.Metadata, error) {
					return nil, errors.New("boom")
				}
			},
			wantErr:   true,
			wantErrIn: "could not list playlist items",
		},
		{
			name: "success",
			args: []string{"playlist"},
			setup: func(t *testing.T) {
				listPlaylistItemsFn = func(string, bool) ([]musicdiscovery.Metadata, error) {
					return []musicdiscovery.Metadata{{Artist: "a", Title: "t"}}, nil
				}
			},
		},
		{
			name: "empty list",
			args: []string{"playlist"},
			setup: func(t *testing.T) {
				listPlaylistItemsFn = func(string, bool) ([]musicdiscovery.Metadata, error) {
					return []musicdiscovery.Metadata{}, nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldList := listPlaylistItemsFn
			t.Cleanup(func() {
				listPlaylistItemsFn = oldList
			})

			if tt.setup != nil {
				tt.setup(t)
			}

			err := listCmd.RunE(listCmd, tt.args)
			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrIn != "" {
					assert.Contains(t, err.Error(), tt.wantErrIn)
				}
				return
			}
			require.NoError(t, err)
		})
	}
}
