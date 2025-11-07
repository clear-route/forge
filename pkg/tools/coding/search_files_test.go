package coding

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSearchFilesTool_Name(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewSearchFilesTool(guard)
	if got := tool.Name(); got != "search_files" {
		t.Errorf("Name() = %v, want %v", got, "search_files")
	}
}

func TestSearchFilesTool_Description(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewSearchFilesTool(guard)
	desc := tool.Description()
	if desc == "" {
		t.Error("Description() returned empty string")
	}
}

func TestSearchFilesTool_Schema(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewSearchFilesTool(guard)
	schema := tool.Schema()

	if schema["type"] != "object" {
		t.Error("Schema type should be 'object'")
	}

	props, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Schema should have properties")
	}

	if _, patternOk := props["pattern"]; !patternOk {
		t.Error("Schema should have 'pattern' property")
	}

	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("Schema should have required fields")
	}
	if len(required) != 1 || required[0] != "pattern" {
		t.Error("Schema should require 'pattern' field")
	}
}

func TestSearchFilesTool_IsLoopBreaking(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewSearchFilesTool(guard)
	if tool.IsLoopBreaking() {
		t.Error("SearchFilesTool should not be loop-breaking")
	}
}

func TestSearchFilesTool_Execute_SimpleMatch(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	// Create test file with searchable content
	testFile := filepath.Join(guard.WorkspaceDir(), "test.txt")
	content := "Line 1\nTODO: fix this\nLine 3\n"
	os.WriteFile(testFile, []byte(content), 0644)

	tool := NewSearchFilesTool(guard)
	args, _ := json.Marshal(map[string]interface{}{
		"pattern": "TODO",
	})

	result, err := tool.Execute(context.Background(), args)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(result, "TODO: fix this") {
		t.Error("Result should contain the matching line")
	}
	if !strings.Contains(result, "test.txt") {
		t.Error("Result should contain the filename")
	}
}

func TestSearchFilesTool_Execute_RegexMatch(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	// Create test file
	testFile := filepath.Join(guard.WorkspaceDir(), "code.go")
	content := "func main() {\n\tfmt.Println(\"hello\")\n}\n"
	os.WriteFile(testFile, []byte(content), 0644)

	tool := NewSearchFilesTool(guard)
	args, _ := json.Marshal(map[string]interface{}{
		"pattern": "func \\w+\\(",
	})

	result, err := tool.Execute(context.Background(), args)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(result, "func main()") {
		t.Error("Result should contain function declaration")
	}
}

func TestSearchFilesTool_Execute_WithFilePattern(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	// Create files with different extensions
	os.WriteFile(filepath.Join(guard.WorkspaceDir(), "file1.go"), []byte("package main\n"), 0644)
	os.WriteFile(filepath.Join(guard.WorkspaceDir(), "file2.txt"), []byte("package main\n"), 0644)

	tool := NewSearchFilesTool(guard)
	args, _ := json.Marshal(map[string]interface{}{
		"pattern":      "package",
		"file_pattern": "*.go",
	})

	result, err := tool.Execute(context.Background(), args)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(result, "file1.go") {
		t.Error("Result should contain file1.go")
	}
	if strings.Contains(result, "file2.txt") {
		t.Error("Result should not contain file2.txt")
	}
}

func TestSearchFilesTool_Execute_WithContext(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	// Create test file with context lines
	testFile := filepath.Join(guard.WorkspaceDir(), "test.txt")
	content := "Line 1\nLine 2\nMATCH HERE\nLine 4\nLine 5\n"
	os.WriteFile(testFile, []byte(content), 0644)

	tool := NewSearchFilesTool(guard)
	args, _ := json.Marshal(map[string]interface{}{
		"pattern":       "MATCH",
		"context_lines": 1,
	})

	result, err := tool.Execute(context.Background(), args)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Should contain context lines
	if !strings.Contains(result, "Line 2") {
		t.Error("Result should contain context before match")
	}
	if !strings.Contains(result, "Line 4") {
		t.Error("Result should contain context after match")
	}
}

func TestSearchFilesTool_Execute_NoMatches(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	// Create test file
	testFile := filepath.Join(guard.WorkspaceDir(), "test.txt")
	os.WriteFile(testFile, []byte("No matches here\n"), 0644)

	tool := NewSearchFilesTool(guard)
	args, _ := json.Marshal(map[string]interface{}{
		"pattern": "NOTFOUND",
	})

	result, err := tool.Execute(context.Background(), args)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(result, "No matches found") {
		t.Error("Result should indicate no matches found")
	}
}

func TestSearchFilesTool_Execute_InvalidPattern(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewSearchFilesTool(guard)
	args, _ := json.Marshal(map[string]interface{}{
		"pattern": "[invalid",
	})

	_, err := tool.Execute(context.Background(), args)
	if err == nil {
		t.Error("Execute() should fail for invalid regex pattern")
	}
}

func TestSearchFilesTool_Execute_MissingPattern(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewSearchFilesTool(guard)
	args, _ := json.Marshal(map[string]interface{}{})

	_, err := tool.Execute(context.Background(), args)
	if err == nil {
		t.Error("Execute() should fail when pattern is missing")
	}
}

func TestSearchFilesTool_Execute_Recursive(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	// Create nested structure
	subdir := filepath.Join(guard.WorkspaceDir(), "subdir")
	os.Mkdir(subdir, 0755)
	os.WriteFile(filepath.Join(guard.WorkspaceDir(), "root.txt"), []byte("FOUND in root\n"), 0644)
	os.WriteFile(filepath.Join(subdir, "nested.txt"), []byte("FOUND in nested\n"), 0644)

	tool := NewSearchFilesTool(guard)
	args, _ := json.Marshal(map[string]interface{}{
		"pattern": "FOUND",
	})

	result, err := tool.Execute(context.Background(), args)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Should find matches in both files
	if !strings.Contains(result, "root.txt") {
		t.Error("Result should contain root.txt")
	}
	if !strings.Contains(result, "nested.txt") {
		t.Error("Result should contain nested.txt")
	}
}

func TestSearchFilesTool_Execute_InvalidPath(t *testing.T) {
	guard, cleanup := setupTestWorkspace(t)
	defer cleanup()

	tool := NewSearchFilesTool(guard)
	args, _ := json.Marshal(map[string]interface{}{
		"path":    "../outside",
		"pattern": "test",
	})

	_, err := tool.Execute(context.Background(), args)
	if err == nil {
		t.Error("Execute() should fail for invalid path")
	}
}
