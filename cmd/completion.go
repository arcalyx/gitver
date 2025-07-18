package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

  $ source <(gitver completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ gitver completion bash > /etc/bash_completion.d/gitver
  # macOS:
  $ gitver completion bash > $(brew --prefix)/etc/bash_completion.d/gitver

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ gitver completion zsh > "${fpath[1]}/_gitver"

  # You will need to start a new shell for this setup to take effect.

Fish:

  $ gitver completion fish | source

  # To load completions for each session, execute once:
  $ gitver completion fish > ~/.config/fish/completions/gitver.fish

PowerShell:

  PS> gitver completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> gitver completion powershell > gitver.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			if err := cmd.Root().GenBashCompletion(os.Stdout); err != nil {
				log.Fatalf("Fehler beim Generieren der Bash-Completion: %v", err)
			}
		case "zsh":
			if err := cmd.Root().GenZshCompletion(os.Stdout); err != nil {
				log.Fatalf("Fehler beim Generieren der Zsh-Completion: %v", err)
			}
		case "fish":
			if err := cmd.Root().GenFishCompletion(os.Stdout, true); err != nil {
				log.Fatalf("Fehler beim Generieren der Fish-Completion: %v", err)
			}
		case "powershell":
			if err := cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout); err != nil {
				log.Fatalf("Fehler beim Generieren der PowerShell-Completion: %v", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
