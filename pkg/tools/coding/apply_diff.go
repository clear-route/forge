package coding

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/entrhq/forge/pkg/agent/tools"
	"github.com/entrhq/forge/pkg/security/workspace"
)

// ApplyDiffTool applies search/replace operations to files for precise code editing.
type ApplyDiffTool struct {
	guard *workspace.Guard
}

// NewApplyDiffTool creates a new ApplyDiffTool with workspace security.
func NewApplyDiffTool(guard *workspace.Guard) *ApplyDiffTool {
	return &ApplyDiffTool{
		guard: guard,
	}
}

// Name returns the tool name.
func (t *ApplyDiffTool) Name() string {
	return "apply_diff"
}

// Description returns the tool description.
func (t *ApplyDiffTool) Description() string {
	return "Apply precise search/replace operations to files. Supports multiple edits in a single operation for surgical code changes."
}

// Schema returns the JSON schema for the tool's input parameters.
func (t *ApplyDiffTool) Schema() map[string]interface{} {
	return tools.BaseToolSchema(
		map[string]interface{}{
			"path": map[string]interface{}{
				"type":        "string",
				"description": "Path to the file to edit (relative to workspace)",
			},
			"edits": map[string]interface{}{
				"type":        "array",
				"description": "List of search/replace operations to apply",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"search": map[string]interface{}{
							"type":        "string",
							"description": "Exact text to search for (must match exactly including whitespace)",
						},
						"replace": map[string]interface{}{
							"type":        "string",
							"description": "Text to replace the search text with",
						},
					},
					"required": []string{"search", "replace"},
				},
			},
		},
		[]string{"path", "edits"},
	)
}

// Execute performs the search/replace operations on the file.
func (t *ApplyDiffTool) Execute(ctx context.Context, args json.RawMessage) (string, error) {
	var input struct {
		Path  string `json:"path"`
		Edits []struct {
			Search  string `json:"search"`
			Replace string `json:"replace"`
		} `json:"edits"`
	}

	if err := json.Unmarshal(args, &input); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	if input.Path == "" {
		return "", fmt.Errorf("path is required")
	}

	if len(input.Edits) == 0 {
		return "", fmt.Errorf("at least one edit is required")
	}

	// Resolve path to absolute path
	absPath, err := t.guard.ResolvePath(input.Path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path: %w", err)
	}

	// Validate path is within workspace
	if validateErr := t.guard.ValidatePath(input.Path); validateErr != nil {
		return "", fmt.Errorf("invalid path: %w", validateErr)
	}

	// Read current file content
	content, err := os.ReadFile(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	fileContent := string(content)
	originalContent := fileContent

	// Apply each edit in sequence
	appliedEdits := 0
	for i, edit := range input.Edits {
		if edit.Search == "" {
			return "", fmt.Errorf("edit %d: search text cannot be empty", i+1)
		}

		// Check if search text exists
		if !strings.Contains(fileContent, edit.Search) {
			return "", fmt.Errorf("edit %d: search text not found in file:\n%s", i+1, edit.Search)
		}

		// Count occurrences to warn about multiple matches
		count := strings.Count(fileContent, edit.Search)
		if count > 1 {
			return "", fmt.Errorf("edit %d: search text appears %d times in file, must be unique", i+1, count)
		}

		// Apply the replacement
		fileContent = strings.Replace(fileContent, edit.Search, edit.Replace, 1)
		appliedEdits++
	}

	// Only write if changes were made
	if fileContent == originalContent {
		return "No changes made to file", nil
	}

	// Write the modified content atomically
	tmpPath := absPath + ".tmp"
	if writeErr := os.WriteFile(tmpPath, []byte(fileContent), 0600); writeErr != nil {
		return "", fmt.Errorf("failed to write temporary file: %w", writeErr)
	}

	if renameErr := os.Rename(tmpPath, absPath); renameErr != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to rename temporary file: %w", renameErr)
	}

	// Get relative path for response
	relPath, err := t.guard.MakeRelative(absPath)
	if err != nil {
		relPath = input.Path
	}

	return fmt.Sprintf("Successfully applied %d edit(s) to %s", appliedEdits, relPath), nil
}

// IsLoopBreaking returns whether this tool should break the agent loop.
func (t *ApplyDiffTool) IsLoopBreaking() bool {
	return false
}
