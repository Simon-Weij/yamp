package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeTestConfig(ts *ThemeService, theme string) error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("could not get config dir %w", err)
	}
	path := filepath.Join(configDir, "yamp", "config.json")
	data, err := json.MarshalIndent(Config{Theme: theme}, "", "  ")
	if err != nil {
		return err
	}
	return afero.WriteFile(ts.Fs, path, data, 0644)
}

func TestThemeService_loadTheme(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
		setup   func(ts *ThemeService) error
	}{
		{
			name:    "by default, it should return dark",
			want:    "dark",
			wantErr: false,
			setup:   func(ts *ThemeService) error { return nil },
		},
		{
			name:    "if the theme is light, it should return light",
			want:    "light",
			wantErr: false,
			setup: func(ts *ThemeService) error {
				return writeTestConfig(ts, "light")
			},
		},
		{
			name:    "if the theme is empty, return dark",
			want:    "dark",
			wantErr: false,
			setup: func(ts *ThemeService) error {
				return writeTestConfig(ts, "")
			},
		},
		{
			name:    "if the theme is invalid, return dark",
			want:    "dark",
			wantErr: false,
			setup: func(ts *ThemeService) error {
				return writeTestConfig(ts, "invalidTheme")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &ThemeService{
				Fs: afero.NewMemMapFs(),
			}
			err := tt.setup(ts)
			require.NoError(t, err)
			got, err := ts.loadTheme()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
