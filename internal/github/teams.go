package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v58/github"
)

func FetchAllTeams(ctx context.Context, client *github.Client, org string) ([]*github.Team, error) {
	var allTeams []*github.Team

	opts := &github.ListOptions{PerPage: 100}

	for {
		teams, resp, err := client.Teams.ListTeams(ctx, org, opts)
		if err != nil {
			return nil, err
		}

		allTeams = append(allTeams, teams...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allTeams, nil
}

func PrintTeams(teams []*github.Team, org string) {
	fmt.Printf("Teams in organization '%s':\n\n", org)
	for _, team := range teams {
		fmt.Printf("%s", team.GetName())
		if desc := team.GetDescription(); desc != "" {
			fmt.Printf(": %s", desc)
		}
		fmt.Println()
	}
	fmt.Printf("Total: %d teams\n", len(teams))
}
