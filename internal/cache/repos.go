package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/google/go-github/v58/github"
)

const cacheTTL = 1 * time.Hour

// ReposResult represents the cached result of a repos command run
type ReposResult struct {
	Org       string        `json:"org"`
	Team      string        `json:"team,omitempty"`
	NoOwner   bool          `json:"noOwner,omitempty"`
	Repos     []*CachedRepo `json:"repos"`
	RunAt     string        `json:"runAt"`
	CachePath string        `json:"cachePath"`
}

type CachedRepo struct {
	Name     string `json:"name"`
	HTMLURL  string `json:"html_url"`
	CloneURL string `json:"clone_url"`
}

// cacheReposResult saves the repos result to cache
func CacheResult(org, team string, noOwner bool, repos []*github.Repository) error {
	cachedRepos := make([]*CachedRepo, len(repos))
	for i, r := range repos {
		cachedRepos[i] = &CachedRepo{
			Name:     r.GetName(),
			HTMLURL:  r.GetHTMLURL(),
			CloneURL: r.GetCloneURL(),
		}
	}

	result := &ReposResult{
		Org:     org,
		Team:    team,
		NoOwner: noOwner,
		Repos:   cachedRepos,
		RunAt:   time.Now().Format(time.RFC3339),
	}

	return cacheReposResult(result)
}

// getValidCache returns the cached result if it matches the parameters and is less than 15 minutes old
func GetValidCache(org, team string, noOwner bool) *ReposResult {
	cached, err := loadCachedReposResult(org)
	if err != nil || cached == nil {
		return nil
	}

	// Check if parameters match
	if cached.Team != team || cached.NoOwner != noOwner {
		return nil
	}

	// Check if cache is still fresh
	runAt, err := time.Parse(time.RFC3339, cached.RunAt)
	if err != nil {
		return nil
	}

	if time.Since(runAt) > cacheTTL {
		return nil
	}

	return cached
}

// CacheReposResult stores the result of a repos command run.
// File is stored as <cache_dir>/<org>/repos-last.json
func cacheReposResult(result *ReposResult) error {
	cacheDir, err := getCacheDir()
	if err != nil {
		return err
	}

	orgCacheDir := filepath.Join(cacheDir, result.Org)
	if err := os.MkdirAll(orgCacheDir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(orgCacheDir, "repos-last.json")
	result.CachePath = filePath

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// LoadCachedReposResult reads the last repos command result from cache.
// Returns nil, nil if no cache exists.
func loadCachedReposResult(org string) (*ReposResult, error) {
	cacheDir, err := getCacheDir()
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(cacheDir, org, "repos-last.json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var result ReposResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
