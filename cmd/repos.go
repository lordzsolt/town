package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"team-repos/internal"
	gh "team-repos/internal/github"

	"github.com/spf13/cobra"
)

var (
	team    string
	noOwner bool
)

var reposCmd = &cobra.Command{
	Use:   "repos",
	Short: "Find repositories based on CODEOWNERS",
	Long: `Searches all repositories in the organization. By default, returns those
where the specified team is mentioned in the CODEOWNERS file.

Use --no-owner to find repositories without a CODEOWNERS file.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if org == "" {
			return fmt.Errorf("organization is required: use --org flag or set default_org in config")
		}

		// --no-owner mode doesn't need a team
		if noOwner {
			return nil
		}

		// Apply config default for team if not provided
		if team == "" && cfg != nil {
			team = cfg.DefaultTeam
		}

		if team == "" {
			return fmt.Errorf("team is required: use --team flag or set default_team in config (or use --no-owner)")
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

		if noOwner {
			repos, err := gh.FetchReposWithoutCodeowners(ctx, client, org)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error fetching repos:", err)
				os.Exit(1)
			}
			gh.PrintRepoCount(repos)
			return
		}

		repos, err := gh.FetchReposWithTeamInCodeowners(ctx, client, org, team)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error fetching repos:", err)
			os.Exit(1)
		}

		gh.PrintRepoCount(repos)
	},
}

func init() {
	rootCmd.AddCommand(reposCmd)
	reposCmd.Flags().StringVarP(&team, "team", "t", "", "Team name to search for in CODEOWNERS")
	reposCmd.Flags().BoolVar(&noOwner, "no-owner", false, "List repositories without a CODEOWNERS file")

	// Register completion for --team flag using cached teams
	reposCmd.RegisterFlagCompletionFunc("team", completeTeamFlag)
}

// completeTeamFlag provides autocomplete suggestions for the --team flag
func completeTeamFlag(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Determine org: check flag first, then config
	completionOrg, _ := cmd.Flags().GetString("org")
	if completionOrg == "" {
		// Try loading config for default org
		cfg, err := internal.LoadConfig()
		if err == nil && cfg.DefaultOrg != "" {
			completionOrg = cfg.DefaultOrg
		}
	}

	if completionOrg == "" {
		// Can't complete without knowing the org
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// Load cached teams for this org
	teams, err := internal.LoadCachedTeams(completionOrg)
	if err != nil || teams == nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// Filter teams by prefix if user has started typing
	if toComplete != "" {
		var filtered []string
		for _, t := range teams {
			if strings.HasPrefix(t, toComplete) {
				filtered = append(filtered, t)
			}
		}
		return filtered, cobra.ShellCompDirectiveNoFileComp
	}

	return teams, cobra.ShellCompDirectiveNoFileComp
}
