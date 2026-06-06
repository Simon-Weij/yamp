package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Theme string `json:"theme"`
}

type ThemeService struct{}

// TODO: tests & cleanup
func (t *ThemeService) Theme() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("could not get config directory", err)
		return "dark"
	}

	dir := filepath.Join(configDir, "yamp")
	path := filepath.Join(dir, "config.json")

	_ = os.MkdirAll(dir, 0755)

	if data, err := os.ReadFile(path); err == nil {
		var cfg Config
		if jsonErr := json.Unmarshal(data, &cfg); jsonErr == nil {
			return cfg.Theme
		}
	}

	cfg := Config{Theme: "dark"}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fmt.Println("error marshaling:", err)
		return "dark"
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		fmt.Println("error writing file:", err)
		return "dark"
	}

	return "dark"
}
