// Package workspace provides security mechanisms for enforcing workspace boundaries
// on file system operations. It prevents path traversal attacks and ensures all
// operations stay within the designated working directory.
package workspace

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Guard enforces workspace boundary restrictions on file paths.
// It validates that all file operations remain within the workspace directory,
// preventing path traversal attacks and unauthorized file access.
type Guard struct {
	workspaceDir string // Absolute path to workspace root
}

// NewGuard creates a new workspace guard for the given directory.
// The directory path is converted to an absolute path, cleaned, and symlinks are evaluated.
func NewGuard(workspaceDir string) (*Guard, error) {
	if workspaceDir == "" {
		return nil, fmt.Errorf("workspace directory cannot be empty")
	}

	// Convert to absolute path and clean it
	absPath, err := filepath.Abs(workspaceDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve workspace directory: %w", err)
	}

	// Evaluate any symlinks in the workspace path itself
	evalPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate workspace directory symlinks: %w", err)
	}

	return &Guard{
		workspaceDir: evalPath,
	}, nil
}

// ValidatePath checks if the given path is within the workspace boundaries.
// It resolves the path to an absolute path and ensures it's a child of the workspace.
//
// Returns an error if:
// - The path is empty
// - The path contains invalid characters or patterns
// - The resolved path is outside the workspace
// - The path attempts directory traversal
func (g *Guard) ValidatePath(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	// Resolve to absolute path
	resolvedPath, err := g.ResolvePath(path)
	if err != nil {
		return err
	}

	// Check if resolved path is within workspace
	if !g.IsWithinWorkspace(resolvedPath) {
		return fmt.Errorf("path '%s' is outside workspace boundaries", path)
	}

	return nil
}

// ResolvePath converts a relative or absolute path to an absolute path
// within the workspace context. It cleans the path and resolves any
// symbolic links.
func (g *Guard) ResolvePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	// Clean the path to remove any .. or . components
	cleanPath := filepath.Clean(path)

	// If path is already absolute, use it directly
	// Otherwise, join with workspace directory
	var absPath string
	if filepath.IsAbs(cleanPath) {
		absPath = cleanPath
	} else {
		absPath = filepath.Join(g.workspaceDir, cleanPath)
	}

	// Clean the absolute path
	absPath = filepath.Clean(absPath)

	// Evaluate any symbolic links
	evalPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		// If the file doesn't exist yet, that's okay for write operations
		// Just ensure the parent directory structure would be valid
		parentDir := filepath.Dir(absPath)
		if parentDir != absPath {
			evalParent, parentErr := filepath.EvalSymlinks(parentDir)
			if parentErr == nil {
				// Use the evaluated parent with the original filename
				evalPath = filepath.Join(evalParent, filepath.Base(absPath))
			} else {
				// Parent doesn't exist either, evaluate workspace and use relative from there
				evalWorkspace, wsErr := filepath.EvalSymlinks(g.workspaceDir)
				if wsErr == nil {
					// Get the relative path from workspace
					relPath, relErr := filepath.Rel(g.workspaceDir, absPath)
					if relErr == nil {
						evalPath = filepath.Join(evalWorkspace, relPath)
					} else {
						evalPath = absPath
					}
				} else {
					evalPath = absPath
				}
			}
		} else {
			evalPath = absPath
		}
	}

	return evalPath, nil
}

// IsWithinWorkspace checks if an absolute path is within the workspace boundaries.
// This is done by ensuring the path starts with the workspace directory path.
func (g *Guard) IsWithinWorkspace(absPath string) bool {
	// Ensure both paths end with separator for accurate comparison
	workspacePrefix := g.workspaceDir
	if !strings.HasSuffix(workspacePrefix, string(filepath.Separator)) {
		workspacePrefix += string(filepath.Separator)
	}

	testPath := absPath
	if !strings.HasSuffix(testPath, string(filepath.Separator)) && absPath != g.workspaceDir {
		testPath += string(filepath.Separator)
	}

	// Check if path is exactly the workspace or a child of it
	return absPath == g.workspaceDir || strings.HasPrefix(testPath, workspacePrefix)
}

// WorkspaceDir returns the absolute path of the workspace directory.
func (g *Guard) WorkspaceDir() string {
	return g.workspaceDir
}

// MakeRelative converts an absolute path to a path relative to the workspace.
// Returns an error if the path is not within the workspace.
func (g *Guard) MakeRelative(absPath string) (string, error) {
	if !g.IsWithinWorkspace(absPath) {
		return "", fmt.Errorf("path '%s' is not within workspace", absPath)
	}

	relPath, err := filepath.Rel(g.workspaceDir, absPath)
	if err != nil {
		return "", fmt.Errorf("failed to make path relative: %w", err)
	}

	return relPath, nil
}