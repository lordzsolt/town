package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/lordzsolt/town/internal/cache"

	"github.com/google/go-github/v58/github"
)

// CloneRepos clones the given repositories
func CloneRepos(repos []*github.Repository, cloneDir string) {
	if len(repos) == 0 {
		return
	}

	fmt.Printf("\nCloning %d repositories to %s...\n\n", len(repos), cloneDir)

	// Ensure clone directory exists
	if err := os.MkdirAll(cloneDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating clone directory: %v\n", err)
		return
	}

	var cloned, failed int
	for _, repo := range repos {
		var repoName = repo.GetName()
		err := cloneRepo(repoName, repo.GetCloneURL(), cloneDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to clone %s: %v\n", repoName, err)
			failed++
		} else {
			cloned++
		}
	}

	fmt.Printf("\nClone complete: %d cloned, %d failed\n", cloned, failed)
}

// CloneReposFromCache clones repositories from cached results
func CloneReposFromCache(repos []*cache.CachedRepo, cloneDir string) {
	if len(repos) == 0 {
		return
	}

	fmt.Printf("\nCloning %d repositories to %s...\n\n", len(repos), cloneDir)

	// Ensure clone directory exists
	if err := os.MkdirAll(cloneDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating clone directory: %v\n", err)
		return
	}

	var cloned, skipped, failed int
	for _, repo := range repos {
		err := cloneRepo(repo.Name, repo.CloneURL, cloneDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to clone %s: %v\n", repo.Name, err)
			failed++
		} else {
			cloned++
		}
	}

	fmt.Printf("\nClone complete: %d cloned, %d skipped, %d failed\n", cloned, skipped, failed)
}

func cloneRepo(name string, url string, cloneDir string) error {
	targetDir := filepath.Join(cloneDir, name)

	// Skip if directory already exists
	if _, err := os.Stat(targetDir); err == nil {
		fmt.Printf("Skipping %s (already exists)\n", name)
		return nil
	}

	fmt.Printf("Cloning %s...\n", name)
	cmd := exec.Command("git", "clone", url, targetDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
