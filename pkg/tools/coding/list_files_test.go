package coding

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListFilesTool_Name(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewListFilesTool(guard)
	if got := tool.Name(); got != "list_files" {
		t.Errorf("Name() = %v, want %v", got, "list_files")
	}
}

func TestListFilesTool_Description(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewListFilesTool(guard)
	desc := tool.Description()
	if desc == "" {
		t.Error("Description() returned empty string")
	}
}

func TestListFilesTool_Schema(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewListFilesTool(guard)
	schema := tool.Schema()

	if schema["type"] != "object" {
		t.Error("Schema type should be 'object'")
	}

	props, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Schema should have properties")
	}

	if _, ok := props["path"]; !ok {
		t.Error("Schema should have 'path' property")
	}
	if _, ok := props["recursive"]; !ok {
		t.Error("Schema should have 'recursive' property")
	}
	if _, ok := props["pattern"]; !ok {
		t.Error("Schema should have 'pattern' property")
	}
}

func TestListFilesTool_IsLoopBreaking(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewListFilesTool(guard)
	if tool.IsLoopBreaking() {
		t.Error("ListFilesTool should not be loop-breaking")
	}
}

func TestListFilesTool_Execute_NonRecursive(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	// Create test files
	os.WriteFile(filepath.Join(guard.WorkspaceDir(), "file1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(guard.WorkspaceDir(), "file2.go"), []byte("test"), 0644)
	os.Mkdir(filepath.Join(guard.WorkspaceDir(), "subdir"), 0755)

	tool := NewListFilesTool(guard)
	args, _ := json.Marshal(map[string]interface{}{})

	result, err := tool.Execute(context.Background(), args)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Check output contains files
	if !strings.Contains(result, "file1.txt") {
		t.Error("Result should contain file1.txt")
	}
	if !strings.Contains(result, "file2.go") {
		t.Error("Result should contain file2.go")
	}
	if !strings.Contains(result, "subdir") {
		t.Error("Result should contain subdir")
	}
}

func TestListFilesTool_Execute_Recursive(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	// Create nested structure
	subdir := filepath.Join(guard.WorkspaceDir(), "subdir")
	os.Mkdir(subdir, 0755)
	os.WriteFile(filepath.Join(guard.WorkspaceDir(), "root.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(subdir, "nested.txt"), []byte("test"), 0644)

	tool := NewListFilesTool(guard)
	args, _ := json.Marshal(map[string]interface{}{
		"recursive": true,
	})

	result, err := tool.Execute(context.Background(), args)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(result, "root.txt") {
		t.Error("Result should contain root.txt")
	}
	if !strings.Contains(result, "nested.txt") {
		t.Error("Result should contain nested.txt")
	}
}

func TestListFilesTool_Execute_WithPattern(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	// Create test files with different extensions
	os.WriteFile(filepath.Join(guard.WorkspaceDir(), "file1.go"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(guard.WorkspaceDir(), "file2.go"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(guard.WorkspaceDir(), "file3.txt"), []byte("test"), 0644)

	tool := NewListFilesTool(guard)
	args, _ := json.Marshal(map[string]interface{}{
		"pattern": "*.go",
	})

	result, err := tool.Execute(context.Background(), args)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(result, "file1.go") {
		t.Error("Result should contain file1.go")
	}
	if !strings.Contains(result, "file2.go") {
		t.Error("Result should contain file2.go")
	}
	if strings.Contains(result, "file3.txt") {
		t.Error("Result should not contain file3.txt")
	}
}

func TestListFilesTool_Execute_EmptyDirectory(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewListFilesTool(guard)
	args, _ := json.Marshal(map[string]interface{}{})

	result, err := tool.Execute(context.Background(), args)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(result, "No files found") {
		t.Error("Result should indicate no files found")
	}
}

func TestListFilesTool_Execute_InvalidPath(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewListFilesTool(guard)
	args, _ := json.Marshal(map[string]interface{}{
		"path": "../outside",
	})

	_, err := tool.Execute(context.Background(), args)
	if err == nil {
		t.Error("Execute() should fail for invalid path")
	}
}

func TestListFilesTool_Execute_NonExistentPath(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewListFilesTool(guard)
	args, _ := json.Marshal(map[string]interface{}{
		"path": "nonexistent",
	})

	_, err := tool.Execute(context.Background(), args)
	if err == nil {
		t.Error("Execute() should fail for non-existent path")
	}
}
