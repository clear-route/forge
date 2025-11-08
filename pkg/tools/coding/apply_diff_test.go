package coding

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestApplyDiffTool_Name(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewApplyDiffTool(guard)
	if got := tool.Name(); got != "apply_diff" {
		t.Errorf("Name() = %v, want %v", got, "apply_diff")
	}
}

func TestApplyDiffTool_Description(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewApplyDiffTool(guard)
	desc := tool.Description()
	if desc == "" {
		t.Error("Description() returned empty string")
	}
}

func TestApplyDiffTool_Schema(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewApplyDiffTool(guard)
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
	if _, editsOk := props["edits"]; !editsOk {
		t.Error("Schema should have 'edits' property")
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

func TestApplyDiffTool_IsLoopBreaking(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewApplyDiffTool(guard)
	if tool.IsLoopBreaking() {
		t.Error("ApplyDiffTool should not be loop-breaking")
	}
}

func TestApplyDiffTool_Execute_SingleEdit(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewApplyDiffTool(guard)

	// Create test file
	testFile := filepath.Join(guard.WorkspaceDir(), "test.txt")
	originalContent := "Hello, World!\nThis is a test."
	if err := os.WriteFile(testFile, []byte(originalContent), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	args, _ := json.Marshal(map[string]interface{}{
		"path": "test.txt",
		"edits": []map[string]string{
			{
				"search":  "World",
				"replace": "Universe",
			},
		},
	})

	result, err := tool.Execute(context.Background(), args)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result == "" {
		t.Error("Execute() should return non-empty result")
	}

	// Verify file content
	content, _ := os.ReadFile(testFile)
	expected := "Hello, Universe!\nThis is a test."
	if string(content) != expected {
		t.Errorf("File content = %q, want %q", string(content), expected)
	}
}

func TestApplyDiffTool_Execute_MultipleEdits(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewApplyDiffTool(guard)

	// Create test file
	testFile := filepath.Join(guard.WorkspaceDir(), "test.go")
	originalContent := `package main

func hello() {
	fmt.Println("Hello")
}

func world() {
	fmt.Println("World")
}`
	if err := os.WriteFile(testFile, []byte(originalContent), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	args, _ := json.Marshal(map[string]interface{}{
		"path": "test.go",
		"edits": []map[string]string{
			{
				"search":  `fmt.Println("Hello")`,
				"replace": `fmt.Println("Hi")`,
			},
			{
				"search":  `fmt.Println("World")`,
				"replace": `fmt.Println("Universe")`,
			},
		},
	})

	result, err := tool.Execute(context.Background(), args)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result == "" {
		t.Error("Execute() should return non-empty result")
	}

	// Verify file content
	content, _ := os.ReadFile(testFile)
	expectedContent := `package main

func hello() {
	fmt.Println("Hi")
}

func world() {
	fmt.Println("Universe")
}`
	if string(content) != expectedContent {
		t.Errorf("File content = %q, want %q", string(content), expectedContent)
	}
}

func TestApplyDiffTool_Execute_SearchNotFound(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewApplyDiffTool(guard)

	// Create test file
	testFile := filepath.Join(guard.WorkspaceDir(), "test.txt")
	if err := os.WriteFile(testFile, []byte("Hello, World!"), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	args, _ := json.Marshal(map[string]interface{}{
		"path": "test.txt",
		"edits": []map[string]string{
			{
				"search":  "NotFound",
				"replace": "Something",
			},
		},
	})

	_, err := tool.Execute(context.Background(), args)
	if err == nil {
		t.Error("Execute() should return error when search text not found")
	}
}

func TestApplyDiffTool_Execute_DuplicateSearchText(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewApplyDiffTool(guard)

	// Create test file with duplicate text
	testFile := filepath.Join(guard.WorkspaceDir(), "test.txt")
	content := "Hello, World!\nHello, Universe!"
	if err := os.WriteFile(testFile, []byte(content), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	args, _ := json.Marshal(map[string]interface{}{
		"path": "test.txt",
		"edits": []map[string]string{
			{
				"search":  "Hello",
				"replace": "Hi",
			},
		},
	})

	_, err := tool.Execute(context.Background(), args)
	if err == nil {
		t.Error("Execute() should return error when search text appears multiple times")
	}
}

func TestApplyDiffTool_Execute_EmptySearch(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewApplyDiffTool(guard)

	// Create test file
	testFile := filepath.Join(guard.WorkspaceDir(), "test.txt")
	if err := os.WriteFile(testFile, []byte("Hello, World!"), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	args, _ := json.Marshal(map[string]interface{}{
		"path": "test.txt",
		"edits": []map[string]string{
			{
				"search":  "",
				"replace": "Something",
			},
		},
	})

	_, err := tool.Execute(context.Background(), args)
	if err == nil {
		t.Error("Execute() should return error for empty search text")
	}
}

func TestApplyDiffTool_Execute_NoEdits(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewApplyDiffTool(guard)

	args, _ := json.Marshal(map[string]interface{}{
		"path":  "test.txt",
		"edits": []map[string]string{},
	})

	_, err := tool.Execute(context.Background(), args)
	if err == nil {
		t.Error("Execute() should return error when no edits provided")
	}
}

func TestApplyDiffTool_Execute_InvalidPath(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewApplyDiffTool(guard)

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
				"edits": []map[string]string{
					{
						"search":  "test",
						"replace": "test",
					},
				},
			})

			_, err := tool.Execute(context.Background(), args)
			if err == nil {
				t.Errorf("Execute() should return error for %s", tt.name)
			}
		})
	}
}

func TestApplyDiffTool_Execute_MissingPath(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewApplyDiffTool(guard)

	args, _ := json.Marshal(map[string]interface{}{
		"edits": []map[string]string{
			{
				"search":  "test",
				"replace": "test",
			},
		},
	})

	_, err := tool.Execute(context.Background(), args)
	if err == nil {
		t.Error("Execute() should return error for missing path")
	}
}

func TestApplyDiffTool_Execute_FileNotFound(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewApplyDiffTool(guard)

	args, _ := json.Marshal(map[string]interface{}{
		"path": "nonexistent.txt",
		"edits": []map[string]string{
			{
				"search":  "test",
				"replace": "test",
			},
		},
	})

	_, err := tool.Execute(context.Background(), args)
	if err == nil {
		t.Error("Execute() should return error for nonexistent file")
	}
}

func TestApplyDiffTool_Execute_WhitespacePreserved(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewApplyDiffTool(guard)

	// Create test file with specific indentation
	testFile := filepath.Join(guard.WorkspaceDir(), "test.txt")
	originalContent := "    indented line\n  less indented"
	if err := os.WriteFile(testFile, []byte(originalContent), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	args, _ := json.Marshal(map[string]interface{}{
		"path": "test.txt",
		"edits": []map[string]string{
			{
				"search":  "    indented line",
				"replace": "    modified line",
			},
		},
	})

	_, err := tool.Execute(context.Background(), args)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Verify whitespace is preserved
	content, _ := os.ReadFile(testFile)
	expected := "    modified line\n  less indented"
	if string(content) != expected {
		t.Errorf("Whitespace not preserved: got %q, want %q", string(content), expected)
	}
}
