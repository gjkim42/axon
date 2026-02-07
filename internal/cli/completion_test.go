package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestValidArgsFunctionWired(t *testing.T) {
	root := NewRootCommand()

	tests := []struct {
		name string
		path []string
	}{
		{"get task", []string{"get", "task"}},
		{"get taskspawner", []string{"get", "taskspawner"}},
		{"delete task", []string{"delete", "task"}},
		{"logs", []string{"logs"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := findSubcommand(t, root, tt.path)
			if cmd.ValidArgsFunction == nil {
				t.Errorf("expected ValidArgsFunction to be set on %q", tt.name)
			}
		})
	}
}

func TestCompletionWithInvalidKubeconfig(t *testing.T) {
	cfg := &ClientConfig{Kubeconfig: "/nonexistent/kubeconfig"}

	fns := []struct {
		name string
		fn   cobra.CompletionFunc
	}{
		{"completeTaskNames", completeTaskNames(cfg)},
		{"completeTaskSpawnerNames", completeTaskSpawnerNames(cfg)},
	}

	for _, tt := range fns {
		t.Run(tt.name, func(t *testing.T) {
			results, directive := tt.fn(nil, nil, "")
			if len(results) != 0 {
				t.Errorf("expected no completions, got %v", results)
			}
			if directive != cobra.ShellCompDirectiveNoFileComp {
				t.Errorf("expected ShellCompDirectiveNoFileComp, got %d", directive)
			}
		})
	}
}

func TestCompletionSkipsAfterFirstArg(t *testing.T) {
	cfg := &ClientConfig{Kubeconfig: "/nonexistent/kubeconfig"}

	fns := []struct {
		name string
		fn   cobra.CompletionFunc
	}{
		{"completeTaskNames", completeTaskNames(cfg)},
		{"completeTaskSpawnerNames", completeTaskSpawnerNames(cfg)},
	}

	for _, tt := range fns {
		t.Run(tt.name, func(t *testing.T) {
			results, directive := tt.fn(nil, []string{"already-provided"}, "")
			if len(results) != 0 {
				t.Errorf("expected no completions when arg already provided, got %v", results)
			}
			if directive != cobra.ShellCompDirectiveNoFileComp {
				t.Errorf("expected ShellCompDirectiveNoFileComp, got %d", directive)
			}
		})
	}
}

func TestFlagCompletionOutput(t *testing.T) {
	root := NewRootCommand()

	root.SetArgs([]string{"__complete", "get", "task", "--output", ""})
	out := &strings.Builder{}
	root.SetOut(out)
	root.Execute()

	output := out.String()
	if !strings.Contains(output, "yaml") {
		t.Errorf("expected yaml in output flag completions, got %q", output)
	}
	if !strings.Contains(output, "json") {
		t.Errorf("expected json in output flag completions, got %q", output)
	}
	if !strings.Contains(output, ":4") {
		t.Errorf("expected ShellCompDirectiveNoFileComp (:4) in output, got %q", output)
	}
}

func TestFlagCompletionCredentialType(t *testing.T) {
	root := NewRootCommand()

	root.SetArgs([]string{"__complete", "run", "--credential-type", ""})
	out := &strings.Builder{}
	root.SetOut(out)
	root.Execute()

	output := out.String()
	if !strings.Contains(output, "api-key") {
		t.Errorf("expected api-key in credential-type completions, got %q", output)
	}
	if !strings.Contains(output, "oauth") {
		t.Errorf("expected oauth in credential-type completions, got %q", output)
	}
	if !strings.Contains(output, ":4") {
		t.Errorf("expected ShellCompDirectiveNoFileComp (:4) in output, got %q", output)
	}
}

func TestCompletionCommand_Bash(t *testing.T) {
	cmd := NewRootCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"completion", "bash"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "bash completion") {
		t.Error("expected bash completion output")
	}
}

func TestCompletionCommand_Zsh(t *testing.T) {
	cmd := NewRootCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"completion", "zsh"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Len() == 0 {
		t.Error("expected non-empty zsh completion output")
	}
}

func TestCompletionCommand_Fish(t *testing.T) {
	cmd := NewRootCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"completion", "fish"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Len() == 0 {
		t.Error("expected non-empty fish completion output")
	}
}

func TestCompletionCommand_Powershell(t *testing.T) {
	cmd := NewRootCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"completion", "powershell"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Len() == 0 {
		t.Error("expected non-empty powershell completion output")
	}
}

func TestCompletionCommand_InvalidShell(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"completion", "invalid"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for invalid shell")
	}
}

func TestCompletionCommand_NoArgs(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"completion"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when no shell specified")
	}
}

func TestCompletionCommand_ValidArgs(t *testing.T) {
	cmd := NewRootCommand()
	for _, sub := range cmd.Commands() {
		if sub.Name() == "completion" {
			if len(sub.ValidArgs) != 4 {
				t.Fatalf("expected 4 valid args, got %d", len(sub.ValidArgs))
			}
			expected := map[string]bool{"bash": true, "zsh": true, "fish": true, "powershell": true}
			for _, arg := range sub.ValidArgs {
				if !expected[arg] {
					t.Errorf("unexpected valid arg: %s", arg)
				}
			}
			return
		}
	}
	t.Fatal("completion subcommand not found")
}

func findSubcommand(t *testing.T, root *cobra.Command, path []string) *cobra.Command {
	t.Helper()
	cmd := root
	for _, name := range path {
		found := false
		for _, sub := range cmd.Commands() {
			if sub.Name() == name {
				cmd = sub
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("subcommand %q not found under %q", name, cmd.Name())
		}
	}
	return cmd
}
