package coding

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFileTool_Name(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewWriteFileTool(guard)
	if got := tool.Name(); got != "write_file" {
		t.Errorf("Name() = %v, want %v", got, "write_file")
	}
}

func TestWriteFileTool_Description(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewWriteFileTool(guard)
	desc := tool.Description()
	if desc == "" {
		t.Error("Description() returned empty string")
	}
}

func TestWriteFileTool_Schema(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewWriteFileTool(guard)
	schema := tool.Schema()

	// Verify schema structure
	if schema["type"] != "object" {
		t.Error("Schema type should be 'object'")
	}

	props, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Schema should have properties")
	}

	// Verify required properties
	if _, pathOk := props["path"]; !pathOk {
		t.Error("Schema should have 'path' property")
	}
	if _, contentOk := props["content"]; !contentOk {
		t.Error("Schema should have 'content' property")
	}

	// Verify required fields
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("Schema should have required fields")
	}
	if len(required) != 2 {
		t.Errorf("Schema should require 2 fields, got %d", len(required))
	}
}

func TestWriteFileTool_IsLoopBreaking(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewWriteFileTool(guard)
	if tool.IsLoopBreaking() {
		t.Error("WriteFileTool should not be loop-breaking")
	}
}

func TestWriteFileTool_Execute_CreateFile(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewWriteFileTool(guard)
	content := "Hello, World!"

	args, _ := json.Marshal(map[string]interface{}{
		"path":    "test.txt",
		"content": content,
	})

	result, err := tool.Execute(context.Background(), args)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result == "" {
		t.Error("Execute() should return a message")
	}

	// Verify file was created
	filePath := filepath.Join(guard.WorkspaceDir(), "test.txt")
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read created file: %v", err)
	}

	if string(data) != content {
		t.Errorf("File content = %q, want %q", string(data), content)
	}
}

func TestWriteFileTool_Execute_OverwriteFile(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	// Create initial file
	filePath := filepath.Join(guard.WorkspaceDir(), "test.txt")
	initialContent := "Initial content"
	if err := os.WriteFile(filePath, []byte(initialContent), 0644); err != nil {
		t.Fatalf("Failed to create initial file: %v", err)
	}

	// Overwrite file
	tool := NewWriteFileTool(guard)
	newContent := "New content"

	args, _ := json.Marshal(map[string]interface{}{
		"path":    "test.txt",
		"content": newContent,
	})

	result, err := tool.Execute(context.Background(), args)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result == "" {
		t.Error("Execute() should return a message")
	}

	// Verify file was overwritten
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(data) != newContent {
		t.Errorf("File content = %q, want %q", string(data), newContent)
	}
}

func TestWriteFileTool_Execute_CreateDirectories(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewWriteFileTool(guard)
	content := "Content in nested directory"

	args, _ := json.Marshal(map[string]interface{}{
		"path":    "dir1/dir2/test.txt",
		"content": content,
	})

	result, err := tool.Execute(context.Background(), args)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result == "" {
		t.Error("Execute() should return a message")
	}

	// Verify directories and file were created
	filePath := filepath.Join(guard.WorkspaceDir(), "dir1", "dir2", "test.txt")
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read created file: %v", err)
	}

	if string(data) != content {
		t.Errorf("File content = %q, want %q", string(data), content)
	}
}

func TestWriteFileTool_Execute_InvalidPath(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewWriteFileTool(guard)

	tests := []struct {
		name string
		path string
	}{
		{"parent directory", "../outside.txt"},
		{"absolute path", "/tmp/file.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, _ := json.Marshal(map[string]interface{}{
				"path":    tt.path,
				"content": "test",
			})
			_, err := tool.Execute(context.Background(), args)
			if err == nil {
				t.Error("Execute() should fail for invalid path")
			}
		})
	}
}

func TestWriteFileTool_Execute_MissingPath(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewWriteFileTool(guard)
	args, _ := json.Marshal(map[string]interface{}{
		"content": "test",
	})
	_, err := tool.Execute(context.Background(), args)
	if err == nil {
		t.Error("Execute() should fail when path is missing")
	}
}
