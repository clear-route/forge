package coding

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/entrhq/forge/pkg/agent/tools"
	"github.com/entrhq/forge/pkg/security/workspace"
)

func TestApplyDiffTool_GeneratePreview(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "preview-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "test.go")
	originalContent := "package main\n\nfunc hello() {\n\tprintln(\"world\")\n}\n"
	if err := os.WriteFile(testFile, []byte(originalContent), 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	guard, err := workspace.NewGuard(tempDir)
	if err != nil {
		t.Fatalf("failed to create guard: %v", err)
	}
	tool := NewApplyDiffTool(guard)

	args := map[string]interface{}{
		"path": "test.go",
		"edits": []map[string]interface{}{
			{
				"search":  "println(\"world\")",
				"replace": "fmt.Println(\"hello, world!\")",
			},
		},
	}
	argsJSON, _ := json.Marshal(args)

	ctx := context.Background()
	preview, err := tool.GeneratePreview(ctx, argsJSON)
	if err != nil {
		t.Fatalf("GeneratePreview failed: %v", err)
	}

	if preview.Type != tools.PreviewTypeDiff {
		t.Errorf("preview type = %v, want %v", preview.Type, tools.PreviewTypeDiff)
	}

	if preview.Title == "" {
		t.Error("preview title is empty")
	}

	if preview.Content == "" {
		t.Error("preview content is empty")
	}
}

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		filename string
		want     string
	}{
		{"test.go", "go"},
		{"test.py", "python"},
		{"test.js", "javascript"},
		{"test.md", "markdown"},
		{"test.txt", "text"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			got := detectLanguage(tt.filename)
			if got != tt.want {
				t.Errorf("detectLanguage(%q) = %q, want %q", tt.filename, got, tt.want)
			}
		})
	}
}
