package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",
	Long: `Generate shell completion script for town.

To load completions:

Bash:
  $ source <(town completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ town completion bash > /etc/bash_completion.d/town
  # macOS:
  $ town completion bash > $(brew --prefix)/etc/bash_completion.d/town

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ town completion zsh > "${fpath[1]}/_town"

  # You may need to start a new shell for this setup to take effect.

Fish:
  $ town completion fish | source
  # To load completions for each session, execute once:
  $ town completion fish > ~/.config/fish/completions/town.fish

PowerShell:
  PS> town completion powershell | Out-String | Invoke-Expression
  # To load completions for every new session, run:
  PS> town completion powershell > town.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
