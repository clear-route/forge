package coding

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/entrhq/forge/pkg/security/workspace"
)

func TestExecuteCommandTool_Name(t *testing.T) {
	guard, err := workspace.NewGuard(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to create workspace guard: %v", err)
	}

	tool := NewExecuteCommandTool(guard)
	if tool.Name() != "execute_command" {
		t.Errorf("Expected name 'execute_command', got '%s'", tool.Name())
	}
}

func TestExecuteCommandTool_Description(t *testing.T) {
	guard, err := workspace.NewGuard(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to create workspace guard: %v", err)
	}

	tool := NewExecuteCommandTool(guard)
	desc := tool.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}

func TestExecuteCommandTool_Schema(t *testing.T) {
	guard, err := workspace.NewGuard(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to create workspace guard: %v", err)
	}

	tool := NewExecuteCommandTool(guard)
	schema := tool.Schema()

	if schema == nil {
		t.Fatal("Schema should not be nil")
	}

	// Verify schema has required properties
	props, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Schema should have properties")
	}

	if _, hasCommand := props["command"]; !hasCommand {
		t.Error("Schema should have 'command' property")
	}
}

func TestExecuteCommandTool_IsLoopBreaking(t *testing.T) {
	guard, err := workspace.NewGuard(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to create workspace guard: %v", err)
	}

	tool := NewExecuteCommandTool(guard)
	if tool.IsLoopBreaking() {
		t.Error("ExecuteCommandTool should not be loop breaking")
	}
}

func TestExecuteCommandTool_Execute_SimpleCommand(t *testing.T) {
	workspaceDir := t.TempDir()
	guard, err := workspace.NewGuard(workspaceDir)
	if err != nil {
		t.Fatalf("Failed to create workspace guard: %v", err)
	}

	tool := NewExecuteCommandTool(guard)
	ctx := context.Background()

	args, _ := json.Marshal(map[string]interface{}{
		"command": "echo 'Hello, World!'",
	})

	result, err := tool.Execute(ctx, args)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !strings.Contains(result, "Hello, World!") {
		t.Errorf("Expected output to contain 'Hello, World!', got: %s", result)
	}

	if !strings.Contains(result, "Exit code: 0") {
		t.Errorf("Expected exit code 0, got: %s", result)
	}
}

func TestExecuteCommandTool_Execute_CommandWithStderr(t *testing.T) {
	workspaceDir := t.TempDir()
	guard, err := workspace.NewGuard(workspaceDir)
	if err != nil {
		t.Fatalf("Failed to create workspace guard: %v", err)
	}

	tool := NewExecuteCommandTool(guard)
	ctx := context.Background()

	// Use a command that fails to ensure stderr is captured
	args, _ := json.Marshal(map[string]interface{}{
		"command": "echo 'error message' >&2; exit 1",
	})

	result, err := tool.Execute(ctx, args)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !strings.Contains(result, "Stderr:") {
		t.Errorf("Expected stderr section in output, got: %s", result)
	}

	if !strings.Contains(result, "error message") {
		t.Errorf("Expected 'error message' in output, got: %s", result)
	}
}

func TestExecuteCommandTool_Execute_FailingCommand(t *testing.T) {
	workspaceDir := t.TempDir()
	guard, err := workspace.NewGuard(workspaceDir)
	if err != nil {
		t.Fatalf("Failed to create workspace guard: %v", err)
	}

	tool := NewExecuteCommandTool(guard)
	ctx := context.Background()

	args, _ := json.Marshal(map[string]interface{}{
		"command": "exit 42",
	})

	result, err := tool.Execute(ctx, args)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !strings.Contains(result, "Exit code: 42") {
		t.Errorf("Expected exit code 42, got: %s", result)
	}

	if !strings.Contains(result, "failed") {
		t.Errorf("Expected failure message, got: %s", result)
	}
}

func TestExecuteCommandTool_Execute_WithTimeout(t *testing.T) {
	workspaceDir := t.TempDir()
	guard, err := workspace.NewGuard(workspaceDir)
	if err != nil {
		t.Fatalf("Failed to create workspace guard: %v", err)
	}

	tool := NewExecuteCommandTool(guard)
	ctx := context.Background()

	args, _ := json.Marshal(map[string]interface{}{
		"command": "sleep 5",
		"timeout": 0.5, // 500ms timeout
	})

	result, err := tool.Execute(ctx, args)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !strings.Contains(result, "timed out") {
		t.Errorf("Expected timeout message, got: %s", result)
	}
}

func TestExecuteCommandTool_Execute_WithWorkingDirectory(t *testing.T) {
	workspaceDir := t.TempDir()
	guard, err := workspace.NewGuard(workspaceDir)
	if err != nil {
		t.Fatalf("Failed to create workspace guard: %v", err)
	}

	// Create a subdirectory
	subdir := filepath.Join(workspaceDir, "subdir")
	if mkdirErr := os.Mkdir(subdir, 0755); mkdirErr != nil {
		t.Fatalf("Failed to create subdirectory: %v", mkdirErr)
	}

	tool := NewExecuteCommandTool(guard)
	ctx := context.Background()

	args, _ := json.Marshal(map[string]interface{}{
		"command":     "pwd",
		"working_dir": "subdir",
	})

	result, err := tool.Execute(ctx, args)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !strings.Contains(result, "subdir") {
		t.Errorf("Expected output to contain 'subdir', got: %s", result)
	}
}

func TestExecuteCommandTool_Execute_InvalidWorkingDirectory(t *testing.T) {
	workspaceDir := t.TempDir()
	guard, err := workspace.NewGuard(workspaceDir)
	if err != nil {
		t.Fatalf("Failed to create workspace guard: %v", err)
	}

	tool := NewExecuteCommandTool(guard)
	ctx := context.Background()

	args, _ := json.Marshal(map[string]interface{}{
		"command":     "pwd",
		"working_dir": "../outside",
	})

	_, err = tool.Execute(ctx, args)
	if err == nil {
		t.Error("Expected error for invalid working directory")
	}

	if !strings.Contains(err.Error(), "invalid working directory") {
		t.Errorf("Expected 'invalid working directory' error, got: %v", err)
	}
}

func TestExecuteCommandTool_Execute_EmptyCommand(t *testing.T) {
	workspaceDir := t.TempDir()
	guard, err := workspace.NewGuard(workspaceDir)
	if err != nil {
		t.Fatalf("Failed to create workspace guard: %v", err)
	}

	tool := NewExecuteCommandTool(guard)
	ctx := context.Background()

	args, _ := json.Marshal(map[string]interface{}{
		"command": "",
	})

	_, err = tool.Execute(ctx, args)
	if err == nil {
		t.Error("Expected error for empty command")
	}

	if !strings.Contains(err.Error(), "command cannot be empty") {
		t.Errorf("Expected 'command cannot be empty' error, got: %v", err)
	}
}

func TestExecuteCommandTool_Execute_InvalidJSON(t *testing.T) {
	workspaceDir := t.TempDir()
	guard, err := workspace.NewGuard(workspaceDir)
	if err != nil {
		t.Fatalf("Failed to create workspace guard: %v", err)
	}

	tool := NewExecuteCommandTool(guard)
	ctx := context.Background()

	_, err = tool.Execute(ctx, json.RawMessage("invalid json"))
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestExecuteCommandTool_Execute_ContextCancellation(t *testing.T) {
	workspaceDir := t.TempDir()
	guard, err := workspace.NewGuard(workspaceDir)
	if err != nil {
		t.Fatalf("Failed to create workspace guard: %v", err)
	}

	tool := NewExecuteCommandTool(guard)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	args, _ := json.Marshal(map[string]interface{}{
		"command": "sleep 10",
	})

	result, err := tool.Execute(ctx, args)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Should timeout due to context cancellation
	if !strings.Contains(result, "timed out") {
		t.Errorf("Expected timeout due to context cancellation, got: %s", result)
	}
}

func TestExecuteCommandTool_Execute_MultilineOutput(t *testing.T) {
	workspaceDir := t.TempDir()
	guard, err := workspace.NewGuard(workspaceDir)
	if err != nil {
		t.Fatalf("Failed to create workspace guard: %v", err)
	}

	tool := NewExecuteCommandTool(guard)
	ctx := context.Background()

	args, _ := json.Marshal(map[string]interface{}{
		"command": "echo 'line1'; echo 'line2'; echo 'line3'",
	})

	result, err := tool.Execute(ctx, args)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !strings.Contains(result, "line1") || !strings.Contains(result, "line2") || !strings.Contains(result, "line3") {
		t.Errorf("Expected all three lines in output, got: %s", result)
	}
}

func TestExecuteCommandTool_Execute_FileCreation(t *testing.T) {
	workspaceDir := t.TempDir()
	guard, err := workspace.NewGuard(workspaceDir)
	if err != nil {
		t.Fatalf("Failed to create workspace guard: %v", err)
	}

	tool := NewExecuteCommandTool(guard)
	ctx := context.Background()

	testFile := "test.txt"
	args, _ := json.Marshal(map[string]interface{}{
		"command": "echo 'test content' > " + testFile,
	})

	result, err := tool.Execute(ctx, args)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !strings.Contains(result, "Exit code: 0") {
		t.Errorf("Expected success, got: %s", result)
	}

	// Verify file was created
	filePath := filepath.Join(workspaceDir, testFile)
	if _, statErr := os.Stat(filePath); os.IsNotExist(statErr) {
		t.Error("Expected file to be created")
	}
}
