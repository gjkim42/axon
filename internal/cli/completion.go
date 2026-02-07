package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/client"

	axonv1alpha1 "github.com/gjkim42/axon/api/v1alpha1"
)

func completeTaskNames(cfg *ClientConfig) cobra.CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cl, ns, err := cfg.NewClient()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		taskList := &axonv1alpha1.TaskList{}
		if err := cl.List(ctx, taskList, client.InNamespace(ns)); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		var names []string
		for _, t := range taskList.Items {
			names = append(names, t.Name)
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	}
}

func completeTaskSpawnerNames(cfg *ClientConfig) cobra.CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cl, ns, err := cfg.NewClient()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		tsList := &axonv1alpha1.TaskSpawnerList{}
		if err := cl.List(ctx, tsList, client.InNamespace(ns)); err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		var names []string
		for _, ts := range tsList.Items {
			names = append(names, ts.Name)
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	}
}

func newCompletionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for axon.

To load completions:

Bash:
  $ source <(axon completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ axon completion bash > /etc/bash_completion.d/axon
  # macOS:
  $ axon completion bash > $(brew --prefix)/etc/bash_completion.d/axon

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ axon completion zsh > "${fpath[1]}/_axon"

  # You will need to start a new shell for this setup to take effect.

Fish:
  $ axon completion fish | source

  # To load completions for each session, execute once:
  $ axon completion fish > ~/.config/fish/completions/axon.fish

PowerShell:
  PS> axon completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> axon completion powershell > axon.ps1
  # and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletionV2(out, true)
			case "zsh":
				return cmd.Root().GenZshCompletion(out)
			case "fish":
				return cmd.Root().GenFishCompletion(out, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(out)
			}
			return fmt.Errorf("unsupported shell: %s", args[0])
		},
	}

	return cmd
}
