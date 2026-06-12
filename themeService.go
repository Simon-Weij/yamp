package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/spf13/afero"
)

type Config struct {
	Theme string `json:"theme"`
}

type ThemeService struct {
	Fs afero.Fs
}

func NewThemeService(fs afero.Fs) *ThemeService {
	return &ThemeService{
		Fs: fs,
	}
}

func (t *ThemeService) Theme() string {
	theme, err := t.loadTheme()
	if err != nil {
		fmt.Println("could not load theme:", err)
		return "dark"
	}
	return theme
}

var allowedThemes = []string{"light", "dark"}

func isValidTheme(theme string) bool {
	return slices.Contains(allowedThemes, theme)
}

func (t *ThemeService) loadTheme() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("could not get config directory: %w", err)
	}
	dir := filepath.Join(configDir, "yamp")
	path := filepath.Join(dir, "config.json")

	if err := t.Fs.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("could not create config directory: %w", err)
	}

	if data, err := afero.ReadFile(t.Fs, path); err == nil {
		var cfg Config
		if jsonErr := json.Unmarshal(data, &cfg); jsonErr == nil && isValidTheme(cfg.Theme) {
			return cfg.Theme, nil
		}
	}

	cfg := Config{Theme: "dark"}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshaling: %w", err)
	}
	if err = afero.WriteFile(t.Fs, path, data, 0644); err != nil {
		return "", fmt.Errorf("error writing file: %w", err)
	}
	return "dark", nil
}
