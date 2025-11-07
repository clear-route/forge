package coding

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/entrhq/forge/pkg/agent/tools"
	"github.com/entrhq/forge/pkg/security/workspace"
)

// WriteFileTool creates or overwrites files with workspace validation.
type WriteFileTool struct {
	guard *workspace.Guard
}

// NewWriteFileTool creates a new WriteFileTool with workspace security.
func NewWriteFileTool(guard *workspace.Guard) *WriteFileTool {
	return &WriteFileTool{
		guard: guard,
	}
}

// Name returns the tool name.
func (t *WriteFileTool) Name() string {
	return "write_file"
}

// Description returns the tool description.
func (t *WriteFileTool) Description() string {
	return "Write content to a file, creating it if it doesn't exist or overwriting if it does. Automatically creates parent directories as needed."
}

// Schema returns the JSON schema for the tool's input parameters.
func (t *WriteFileTool) Schema() map[string]interface{} {
	return tools.BaseToolSchema(
		map[string]interface{}{
			"path": map[string]interface{}{
				"type":        "string",
				"description": "Path to the file to write (relative to workspace)",
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "Content to write to the file",
			},
		},
		[]string{"path", "content"},
	)
}

// Execute writes content to the specified file.
func (t *WriteFileTool) Execute(ctx context.Context, arguments json.RawMessage) (string, error) {
	// Parse arguments
	var input struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}

	if err := json.Unmarshal(arguments, &input); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	if input.Path == "" {
		return "", fmt.Errorf("missing required parameter: path")
	}

	// Validate path with workspace guard
	if err := t.guard.ValidatePath(input.Path); err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}

	// Resolve to absolute path
	absPath, err := t.guard.ResolvePath(input.Path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path: %w", err)
	}

	// Create parent directories if they don't exist
	dir := filepath.Dir(absPath)
	if mkdirErr := os.MkdirAll(dir, 0755); mkdirErr != nil {
		return "", fmt.Errorf("failed to create directories: %w", mkdirErr)
	}

	// Check if file exists
	fileExists := false
	if _, statErr := os.Stat(absPath); statErr == nil {
		fileExists = true
	}

	// Write file atomically using a temporary file
	tmpPath := absPath + ".tmp"
	if writeErr := os.WriteFile(tmpPath, []byte(input.Content), 0600); writeErr != nil {
		return "", fmt.Errorf("failed to write temporary file: %w", writeErr)
	}

	// Rename temporary file to target file (atomic operation)
	if renameErr := os.Rename(tmpPath, absPath); renameErr != nil {
		// Clean up temporary file on error
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to rename temporary file: %w", renameErr)
	}

	// Get relative path for output message
	relPath, err := t.guard.MakeRelative(absPath)
	if err != nil {
		relPath = input.Path // Fallback to original path
	}

	var message string
	if fileExists {
		message = fmt.Sprintf("File '%s' overwritten successfully", relPath)
	} else {
		message = fmt.Sprintf("File '%s' created successfully", relPath)
	}

	return message, nil
}

// IsLoopBreaking returns false as this tool doesn't break the agent loop.
func (t *WriteFileTool) IsLoopBreaking() bool {
	return false
}
