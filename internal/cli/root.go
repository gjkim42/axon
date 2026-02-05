package cli

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cfg := &ClientConfig{}

	cmd := &cobra.Command{
		Use:           "axon",
		Short:         "CLI for managing AI agent Tasks on Kubernetes",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().StringVar(&cfg.Kubeconfig, "kubeconfig", "", "path to kubeconfig file")
	cmd.PersistentFlags().StringVarP(&cfg.Namespace, "namespace", "n", "", "Kubernetes namespace")

	cmd.AddCommand(
		newRunCommand(cfg),
		newGetCommand(cfg),
		newLogsCommand(cfg),
		newDeleteCommand(cfg),
	)

	return cmd
}
