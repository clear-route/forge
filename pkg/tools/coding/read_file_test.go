package coding

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/entrhq/forge/pkg/security/workspace"
)

func TestReadFileTool_Name(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewReadFileTool(guard)
	if got := tool.Name(); got != "read_file" {
		t.Errorf("Name() = %v, want %v", got, "read_file")
	}
}

func TestReadFileTool_Description(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewReadFileTool(guard)
	desc := tool.Description()
	if desc == "" {
		t.Error("Description() returned empty string")
	}
	if !strings.Contains(desc, "Read") {
		t.Error("Description() should mention reading files")
	}
}

func TestReadFileTool_Schema(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewReadFileTool(guard)
	schema := tool.Schema()

	// Verify schema structure
	if schema["type"] != "object" {
		t.Error("Schema type should be 'object'")
	}

	props, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Schema should have properties")
	}

	// Verify path property
	if _, pathOk := props["path"]; !pathOk {
		t.Error("Schema should have 'path' property")
	}

	// Verify required fields
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("Schema should have required fields")
	}
	if len(required) != 1 || required[0] != "path" {
		t.Error("Schema should require 'path' field")
	}
}

func TestReadFileTool_IsLoopBreaking(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewReadFileTool(guard)
	if tool.IsLoopBreaking() {
		t.Error("ReadFileTool should not be loop-breaking")
	}
}

func TestReadFileTool_Execute_WholeFile(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	// Create test file
	testFile := filepath.Join(guard.WorkspaceDir(), "test.txt")
	content := "Line 1\nLine 2\nLine 3\n"
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tool := NewReadFileTool(guard)
	args, _ := json.Marshal(map[string]interface{}{
		"path": "test.txt",
	})

	result, err := tool.Execute(context.Background(), args)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	expected := "1 | Line 1\n2 | Line 2\n3 | Line 3"
	if result != expected {
		t.Errorf("Execute() output = %q, want %q", result, expected)
	}
}

func TestReadFileTool_Execute_LineRange(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	// Create test file with multiple lines
	testFile := filepath.Join(guard.WorkspaceDir(), "test.txt")
	content := "Line 1\nLine 2\nLine 3\nLine 4\nLine 5\n"
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tool := NewReadFileTool(guard)
	args, _ := json.Marshal(map[string]interface{}{
		"path":       "test.txt",
		"start_line": 2,
		"end_line":   4,
	})

	result, err := tool.Execute(context.Background(), args)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	expected := "2 | Line 2\n3 | Line 3\n4 | Line 4"
	if result != expected {
		t.Errorf("Execute() output = %q, want %q", result, expected)
	}
}

func TestReadFileTool_Execute_InvalidPath(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewReadFileTool(guard)

	tests := []struct {
		name string
		path string
	}{
		{"parent directory", "../outside.txt"},
		{"absolute path", "/etc/passwd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, _ := json.Marshal(map[string]interface{}{
				"path": tt.path,
			})
			_, err := tool.Execute(context.Background(), args)
			if err == nil {
				t.Error("Execute() should fail for invalid path")
			}
		})
	}
}

// setupTestWorkspace creates a temporary workspace for testing.
func setupTestWorkspace(t *testing.T) (*workspace.Guard, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "forge-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	guard, err := workspace.NewGuard(tmpDir)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to create guard: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return guard, cleanup
}
