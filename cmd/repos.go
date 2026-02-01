package cmd

import (
	"context"
	"fmt"
	"os"

	gh "team-repos/internal/github"

	"github.com/spf13/cobra"
)

var team string

var reposCmd = &cobra.Command{
	Use:   "repos",
	Short: "Find repositories where a team is mentioned in CODEOWNERS",
	Long: `Searches all repositories in the organization and returns those
where the specified team is mentioned in the CODEOWNERS file.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Apply config default for team if not provided
		if team == "" && cfg != nil {
			team = cfg.DefaultTeam
		}

		if org == "" {
			return fmt.Errorf("organization is required: use --org flag or set default_org in config")
		}
		if team == "" {
			return fmt.Errorf("team is required: use --team flag or set default_team in config")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		client, err := gh.NewClient()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		ctx := context.Background()
		repos, err := gh.FetchReposWithTeamInCodeowners(ctx, client, org, team)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error fetching repos:", err)
			os.Exit(1)
		}

		gh.PrintReposWithTeam(repos, team)
	},
}

func init() {
	rootCmd.AddCommand(reposCmd)
	reposCmd.Flags().StringVarP(&team, "team", "t", "", "Team name to search for in CODEOWNERS")
}
