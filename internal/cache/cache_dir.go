package cache

import (
	"os"
	"path/filepath"
)

const appName = "town"

// GetCacheDir returns the cache directory following XDG Base Directory Specification.
// Priority:
//  1. $XDG_CACHE_HOME/town
//  2. ~/.town/cache/ (fallback)
func getCacheDir() (string, error) {
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
