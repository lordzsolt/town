package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	org string
)

var rootCmd = &cobra.Command{
	Use:   "town",
	Short: "A CLI tool for exploring GitHub teams and repositories",
	Long: `town is a CLI tool that helps you explore GitHub organizations,
their teams, and find repositories owned by specific teams via CODEOWNERS.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&org, "org", "o", "", "GitHub organization name (required)")
	rootCmd.MarkPersistentFlagRequired("org")
}
