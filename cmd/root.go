package cmd

import (
	"fmt"
	"os"

	"team-repos/internal"

	"github.com/spf13/cobra"
)

var (
	org string
	cfg *internal.Config
)

var rootCmd = &cobra.Command{
	Use:   "town",
	Short: "A CLI tool for exploring GitHub teams and repositories",
	Long: `town is a CLI tool that helps you explore GitHub organizations,
their teams, and find repositories owned by specific teams via CODEOWNERS.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Load config
		var err error
		cfg, err = internal.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Apply config defaults if flags not provided
		if org == "" {
			org = cfg.DefaultOrg
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&org, "org", "o", "", "GitHub organization name")
}
