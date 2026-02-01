package cmd

import (
	"context"
	"fmt"
	"os"

	gh "team-repos/internal/github"

	"github.com/spf13/cobra"
)

var teamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "List all teams in an organization",
	Long:  `Fetches and displays all teams in the specified GitHub organization.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if org == "" {
			return fmt.Errorf("organization is required: use --org flag or set default_org in config")
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
		teams, err := gh.FetchAllTeams(ctx, client, org)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error fetching teams:", err)
			os.Exit(1)
		}

		gh.PrintTeams(teams, org)
	},
}

func init() {
	rootCmd.AddCommand(teamsCmd)
}
