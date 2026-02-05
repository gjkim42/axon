package commands

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	axonv1alpha1 "github.com/gjkim/axon/api/v1alpha1"
)

// NewLogsCmd creates a new logs command
func NewLogsCmd() *cobra.Command {
	var follow bool
	var tail int64

	cmd := &cobra.Command{
		Use:   "logs TASK_NAME",
		Short: "Get logs from a task's pod",
		Long:  `Fetch logs from the pod running the specified task.`,
		Example: `  # Get logs from a task
  axonctl logs my-task

  # Follow logs from a task
  axonctl logs my-task -f

  # Get last 50 lines of logs
  axonctl logs my-task --tail=50`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskName := args[0]

			axonClient, err := getAxonClient()
			if err != nil {
				return fmt.Errorf("failed to create axon client: %w", err)
			}

			// Get the task to find the pod name
			ctx := context.Background()
			task := &axonv1alpha1.Task{}
			if err := axonClient.Get(ctx, types.NamespacedName{
				Name:      taskName,
				Namespace: namespace,
			}, task); err != nil {
				return fmt.Errorf("failed to get task: %w", err)
			}

			if task.Status.PodName == "" {
				return fmt.Errorf("task %s has no pod yet", taskName)
			}

			// Get Kubernetes client for pod logs
			kubeClient, err := getKubeClient()
			if err != nil {
				return fmt.Errorf("failed to create kubernetes client: %w", err)
			}

			// Prepare log options
			logOpts := &corev1.PodLogOptions{
				Follow: follow,
			}
			if tail > 0 {
				logOpts.TailLines = &tail
			}

			// Get logs
			req := kubeClient.CoreV1().Pods(namespace).GetLogs(task.Status.PodName, logOpts)
			stream, err := req.Stream(ctx)
			if err != nil {
				return fmt.Errorf("failed to get logs: %w", err)
			}
			defer stream.Close()

			// Stream logs to stdout
			reader := bufio.NewReader(stream)
			for {
				line, err := reader.ReadBytes('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					return fmt.Errorf("error reading logs: %w", err)
				}
				fmt.Print(string(line))
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output")
	cmd.Flags().Int64Var(&tail, "tail", -1, "Number of lines to show from the end of the logs")

	return cmd
}
