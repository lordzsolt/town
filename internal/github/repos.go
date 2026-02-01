package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v58/github"
)

// Common CODEOWNERS file locations
var codeownersLocations = []string{
	".github/CODEOWNERS",
	"CODEOWNERS",
}

func FetchAllRepos(ctx context.Context, client *github.Client, org string) ([]*github.Repository, error) {
	var allRepos []*github.Repository

	opts := &github.RepositoryListByOrgOptions{
		Type:        "all",
		ListOptions: github.ListOptions{PerPage: 100},
		Sort:        "updated",
		Direction:   "desc",
	}

	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, org, opts)
		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repos...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allRepos, nil
}

func GetCodeownersContent(ctx context.Context, client *github.Client, owner, repo string) (string, error) {
	for _, path := range codeownersLocations {
		content, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
		if err != nil {
			continue // Try next location
		}

		if content != nil {
			decoded, err := content.GetContent()
			if err != nil {
				return "", err
			}
			return decoded, nil
		}
	}

	return "", nil // No CODEOWNERS found
}

func FetchReposWithTeamInCodeowners(ctx context.Context, client *github.Client, org string, team string) ([]*github.Repository, error) {
	repos, err := FetchAllRepos(ctx, client, org)
	if err != nil {
		return nil, fmt.Errorf("fetching repos: %w", err)
	}

	fmt.Printf("Scanning %d repositories for team '%s'...\n\n", len(repos), team)

	var results []*github.Repository

	for _, repo := range repos {
		if repo.GetArchived() {
			continue // Skip archived repos
		}

		content, err := GetCodeownersContent(ctx, client, org, repo.GetName())
		if err != nil {
			continue // Skip repos we can't access
		}

		if content == "" {
			continue // No CODEOWNERS file
		}

		if !strings.Contains(strings.ToLower(content), strings.ToLower(team)) {
			continue
		}

		fmt.Printf("%s\n%s\n\n", repo.GetFullName(), repo.GetHTMLURL())
		results = append(results, repo)
	}

	fmt.Println()
	return results, nil
}

func PrintReposWithTeam(results []*github.Repository, team string) {
	fmt.Printf("\nTotal: %d repositories\n", len(results))
}
