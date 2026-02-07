package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

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

func TestGetTaskCommand_HasValidArgsFunction(t *testing.T) {
	cmd := NewRootCommand()
	getCmd, _, err := cmd.Find([]string{"get", "task"})
	if err != nil {
		t.Fatalf("finding get task command: %v", err)
	}
	if getCmd.ValidArgsFunction == nil {
		t.Error("expected ValidArgsFunction to be set on get task command")
	}
}

func TestGetTaskSpawnerCommand_HasValidArgsFunction(t *testing.T) {
	cmd := NewRootCommand()
	getCmd, _, err := cmd.Find([]string{"get", "taskspawner"})
	if err != nil {
		t.Fatalf("finding get taskspawner command: %v", err)
	}
	if getCmd.ValidArgsFunction == nil {
		t.Error("expected ValidArgsFunction to be set on get taskspawner command")
	}
}

func TestDeleteTaskCommand_HasValidArgsFunction(t *testing.T) {
	cmd := NewRootCommand()
	deleteCmd, _, err := cmd.Find([]string{"delete", "task"})
	if err != nil {
		t.Fatalf("finding delete task command: %v", err)
	}
	if deleteCmd.ValidArgsFunction == nil {
		t.Error("expected ValidArgsFunction to be set on delete task command")
	}
}

func TestLogsCommand_HasValidArgsFunction(t *testing.T) {
	cmd := NewRootCommand()
	logsCmd, _, err := cmd.Find([]string{"logs"})
	if err != nil {
		t.Fatalf("finding logs command: %v", err)
	}
	if logsCmd.ValidArgsFunction == nil {
		t.Error("expected ValidArgsFunction to be set on logs command")
	}
}

func TestTaskCompletionFunc_NoClientReturnsError(t *testing.T) {
	// With no kubeconfig, NewClient should fail and return error directive.
	cfg := &ClientConfig{Kubeconfig: "/nonexistent/kubeconfig"}
	fn := taskCompletionFunc(cfg)
	completions, directive := fn(&cobra.Command{}, nil, "")
	if completions != nil {
		t.Errorf("expected nil completions, got %v", completions)
	}
	if directive&cobra.ShellCompDirectiveError == 0 {
		t.Error("expected ShellCompDirectiveError")
	}
	if directive&cobra.ShellCompDirectiveNoFileComp == 0 {
		t.Error("expected ShellCompDirectiveNoFileComp")
	}
}

func TestTaskSpawnerCompletionFunc_NoClientReturnsError(t *testing.T) {
	cfg := &ClientConfig{Kubeconfig: "/nonexistent/kubeconfig"}
	fn := taskSpawnerCompletionFunc(cfg)
	completions, directive := fn(&cobra.Command{}, nil, "")
	if completions != nil {
		t.Errorf("expected nil completions, got %v", completions)
	}
	if directive&cobra.ShellCompDirectiveError == 0 {
		t.Error("expected ShellCompDirectiveError")
	}
	if directive&cobra.ShellCompDirectiveNoFileComp == 0 {
		t.Error("expected ShellCompDirectiveNoFileComp")
	}
}

func TestTaskCompletionFunc_NoCompletionAfterFirstArg(t *testing.T) {
	cfg := &ClientConfig{}
	fn := taskCompletionFunc(cfg)
	completions, directive := fn(&cobra.Command{}, []string{"already-provided"}, "")
	if completions != nil {
		t.Errorf("expected nil completions when arg already provided, got %v", completions)
	}
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("expected ShellCompDirectiveNoFileComp, got %d", directive)
	}
}

func TestTaskSpawnerCompletionFunc_NoCompletionAfterFirstArg(t *testing.T) {
	cfg := &ClientConfig{}
	fn := taskSpawnerCompletionFunc(cfg)
	completions, directive := fn(&cobra.Command{}, []string{"already-provided"}, "")
	if completions != nil {
		t.Errorf("expected nil completions when arg already provided, got %v", completions)
	}
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("expected ShellCompDirectiveNoFileComp, got %d", directive)
	}
}
