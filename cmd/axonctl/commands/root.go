package commands

import (
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"

	axonv1alpha1 "github.com/gjkim/axon/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	kubeconfig string
	namespace  string
)

// NewRootCmd creates the root command for axonctl
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "axonctl",
		Short: "CLI for managing Axon tasks",
		Long:  `axonctl is a command-line tool for managing AI agent tasks in Kubernetes using Axon.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	// Global flags
	if home := homedir.HomeDir(); home != "" {
		rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "path to kubeconfig file")
	} else {
		rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "path to kubeconfig file")
	}
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes namespace")

	// Add subcommands
	rootCmd.AddCommand(NewCreateCmd())
	rootCmd.AddCommand(NewListCmd())
	rootCmd.AddCommand(NewGetCmd())
	rootCmd.AddCommand(NewDeleteCmd())
	rootCmd.AddCommand(NewLogsCmd())
	rootCmd.AddCommand(NewVersionCmd())

	return rootCmd
}

// getKubeConfig returns the Kubernetes configuration
func getKubeConfig() (*rest.Config, error) {
	// Try in-cluster config first
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	// Fall back to kubeconfig
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

// getKubeClient returns a Kubernetes client
func getKubeClient() (kubernetes.Interface, error) {
	config, err := getKubeConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

// getAxonClient returns a client for Axon resources
func getAxonClient() (client.Client, error) {
	config, err := getKubeConfig()
	if err != nil {
		return nil, err
	}

	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = axonv1alpha1.AddToScheme(scheme)

	return client.New(config, client.Options{Scheme: scheme})
}
