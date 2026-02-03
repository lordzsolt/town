package internal

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const appName = "town"

// Config holds the application configuration
type Config struct {
	DefaultOrg  string `json:"default_org"`
	DefaultTeam string `json:"default_team"`
}

// LoadConfig reads the config file following XDG Base Directory Specification.
// It checks in order:
//  1. $XDG_CONFIG_HOME/town/config.json
//  2. ~/.town/config.json (fallback)
//
// Returns an empty Config (not an error) if no config file exists.
func LoadConfig() (*Config, error) {
	configPath, err := findConfigFile()
	if err != nil {
		// No config file found - return empty config
		if errors.Is(err, os.ErrNotExist) {
			return &Config{}, nil
		}
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// findConfigFile returns the path to the config file if it exists.
// Returns os.ErrNotExist if no config file is found.
func findConfigFile() (string, error) {
	candidates := getConfigPaths()

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", os.ErrNotExist
}

// getConfigPaths returns the list of config file paths to check, in priority order.
func getConfigPaths() []string {
	var paths []string

	// XDG_CONFIG_HOME or default ~/.config
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome != "" {
		paths = append(paths, filepath.Join(configHome, appName, "config.json"))
	}

	// Legacy fallback: ~/.town/config.json
	home, err := os.UserHomeDir()
	if err == nil {
		paths = append(paths, filepath.Join(home, "."+appName, "config.json"))
	}

	return paths
}

// ConfigExists returns true if a config file exists
func ConfigExists() bool {
	_, err := findConfigFile()
	return err == nil
}

// SaveConfig writes the config to the preferred config file location.
// Creates the config directory if it doesn't exist.
func SaveConfig(cfg *Config) error {
	paths := getConfigPaths()
	if len(paths) == 0 {
		return errors.New("no config paths found")
	}

	configPath := getConfigPaths()[0]

	// Ensure directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
