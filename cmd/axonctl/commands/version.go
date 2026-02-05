package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "v0.1.0"

// NewVersionCmd creates a new version command
func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Long:  `Print the version number of axonctl.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("axonctl version %s\n", version)
		},
	}

	return cmd
}
