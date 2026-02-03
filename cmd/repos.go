package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"team-repos/internal"
	"team-repos/internal/cache"
	gh "team-repos/internal/github"

	"github.com/google/go-github/v58/github"
	"github.com/spf13/cobra"
)

var (
	team     string
	noOwner  bool
	clone    bool
	cloneDir string
)

var reposCmd = &cobra.Command{
	Use:   "repos",
	Short: "Find repositories based on CODEOWNERS",
	Long: `Searches all repositories in the organization. By default, returns those
where the specified team is mentioned in the CODEOWNERS file.

Use --no-owner to find repositories without a CODEOWNERS file.
Use --clone to clone all matching repositories.

Results are cached for 1 hour to avoid unnecessary API calls.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if org == "" {
			return fmt.Errorf("organization is required: use --org flag or set defaultOrg in config")
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
			return fmt.Errorf("team is required: use --team flag or set defaultTeam in config (or use --no-owner)")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Check if we have a valid cached result
		if cached := cache.GetValidCache(org, team, noOwner); cached != nil {
			printCachedResult(cached)
			if clone {
				internal.CloneReposFromCache(cached.Repos, cloneDir)
			}
			return
		}

		client, err := gh.NewClient()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		ctx := context.Background()

		var repos []*github.Repository
		if noOwner {
			repos, err = gh.FetchReposWithoutCodeowners(ctx, client, org)
		} else {
			repos, err = gh.FetchReposWithTeamInCodeowners(ctx, client, org, team)
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error fetching repos:", err)
			os.Exit(1)
		}

		cache.CacheResult(org, team, noOwner, repos)

		if clone {
			internal.CloneRepos(repos, cloneDir)
		}
	},
}

func init() {
	rootCmd.AddCommand(reposCmd)
	reposCmd.Flags().StringVarP(&team, "team", "t", "", "Team name to search for in CODEOWNERS")
	reposCmd.Flags().BoolVar(&noOwner, "no-owner", false, "List repositories without a CODEOWNERS file")
	reposCmd.Flags().BoolVar(&clone, "clone", false, "Clone all matching repositories")
	reposCmd.Flags().StringVar(&cloneDir, "clone-dir", ".", "Directory to clone repositories into")

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
	teams, err := cache.LoadCachedTeams(completionOrg)
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

// printCachedResult prints the cached repos result
func printCachedResult(cached *cache.ReposResult) {
	runAt, _ := time.Parse(time.RFC3339, cached.RunAt)
	age := time.Since(runAt).Round(time.Second)

	if cached.NoOwner {
		fmt.Println("Repositories without CODEOWNERS:")
	} else {
		fmt.Printf("Repositories where '%s' is mentioned in CODEOWNERS:\n", cached.Team)
	}
	fmt.Println()

	for _, repo := range cached.Repos {
		fmt.Printf("%s\n%s\n\n", repo.Name, repo.HTMLURL)
	}

	fmt.Printf("\nTotal: %d repositories\n", len(cached.Repos))
	fmt.Printf("Used cached result from %s ago\n", age)
	fmt.Printf("You can delete the cache file at %s to force a new search.\n", cached.CachePath)
}
