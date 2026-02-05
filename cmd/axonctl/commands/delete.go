package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/types"

	axonv1alpha1 "github.com/gjkim/axon/api/v1alpha1"
)

// NewDeleteCmd creates a new delete command
func NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete TASK_NAME",
		Short: "Delete a task",
		Long:  `Delete a task and its associated resources.`,
		Example: `  # Delete a task
  axonctl delete my-task

  # Delete a task in specific namespace
  axonctl delete my-task -n my-namespace`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskName := args[0]

			client, err := getAxonClient()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			ctx := context.Background()
			task := &axonv1alpha1.Task{}
			if err := client.Get(ctx, types.NamespacedName{
				Name:      taskName,
				Namespace: namespace,
			}, task); err != nil {
				return fmt.Errorf("failed to get task: %w", err)
			}

			if err := client.Delete(ctx, task); err != nil {
				return fmt.Errorf("failed to delete task: %w", err)
			}

			fmt.Printf("Task %s deleted successfully from namespace %s\n", taskName, namespace)
			return nil
		},
	}

	return cmd
}
