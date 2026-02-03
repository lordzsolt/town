package cmd

import (
	"fmt"
	"os"

	"github.com/lordzsolt/town/internal"

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
		// Check if config exists before loading
		configExists := internal.ConfigExists()

		// Load config
		var err error
		cfg, err = internal.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// If org was provided via flag and no config exists, create one
		if org != "" && !configExists {
			cfg.DefaultOrg = org
			if err := internal.SaveConfig(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not save config: %v\n", err)
			} else {
				fmt.Printf("Created config with default org: %s\n", org)
			}
		}

		// Apply config defaults if flags not provided
		if org == "" {
			org = cfg.DefaultOrg
		}

		return nil
	},
}

// Execute runs the root command (for standalone usage)
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Command returns the root cobra command for embedding in other CLIs.
// Example usage in parent CLI:
//
//	import "github.com/yourusername/town/cmd"
//	parentCmd.AddCommand(cmd.Command())
func Command() *cobra.Command {
	return rootCmd
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&org, "org", "o", "", "GitHub organization name")
}
