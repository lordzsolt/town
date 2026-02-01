package internal

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

const teamsFileName = "teams"

// GetCacheDir returns the cache directory following XDG Base Directory Specification.
// Priority:
//  1. $XDG_CACHE_HOME/town
//  2. ~/.town/cache/ (fallback)
func GetCacheDir() (string, error) {
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome != "" {
		return filepath.Join(cacheHome, appName), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Check if ~/.town exists (legacy), use it if so
	legacyDir := filepath.Join(home, "."+appName)
	return filepath.Join(legacyDir, "cache"), nil
}

// CacheTeams stores team names to the cache file, one per line.
// The file is stored as <cache_dir>/<org>/teams
func CacheTeams(org string, teamNames []string) error {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return err
	}

	orgCacheDir := filepath.Join(cacheDir, org)
	if err := os.MkdirAll(orgCacheDir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(orgCacheDir, teamsFileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, name := range teamNames {
		if _, err := file.WriteString(name + "\n"); err != nil {
			return err
		}
	}

	return nil
}

// LoadCachedTeams reads team names from the cache file.
// Returns nil, nil if the cache file doesn't exist.
func LoadCachedTeams(org string) ([]string, error) {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(cacheDir, org, teamsFileName)
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	var teams []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		if name != "" {
			teams = append(teams, name)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}

// GetTeamsCachePath returns the path to the teams cache file for an org.
func GetTeamsCachePath(org string) (string, error) {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cacheDir, org, teamsFileName), nil
}
