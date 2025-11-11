package workspace

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// defaultIgnorePatterns contains hardcoded patterns that are always ignored.
// These are common directories and files that should be excluded from file operations.
var defaultIgnorePatterns = []string{
	"node_modules/",
	".git/",
	".env",
	".env.*",
	"*.log",
	".DS_Store",
	"vendor/",
	"__pycache__/",
	"*.pyc",
	".vscode/",
	".idea/",
	"dist/",
	"build/",
	"tmp/",
	"temp/",
	"coverage/",
	".next/",
	".nuxt/",
	"target/",
	"*.swp",
	"*.swo",
	"*~",
}

// ignorePattern represents a single ignore pattern with metadata.
type ignorePattern struct {
	pattern   string // Original pattern string
	negation  bool   // True if this is a negation pattern (starts with !)
	dirOnly   bool   // True if pattern only matches directories (ends with /)
	isGlob    bool   // True if pattern contains glob characters
	source    string // Source of pattern: "default", "gitignore", "forgeignore"
}

// IgnoreMatcher handles pattern matching for file ignore rules.
// It supports layered patterns from multiple sources with defined precedence.
type IgnoreMatcher struct {
	patterns []ignorePattern
}

// NewIgnoreMatcher creates a new ignore matcher and loads patterns from all sources.
// Pattern loading order (all are merged, last match wins):
// 1. Default hardcoded patterns
// 2. .gitignore patterns (if file exists)
// 3. .forgeignore patterns (if file exists)
func NewIgnoreMatcher(workspaceDir string) (*IgnoreMatcher, error) {
	m := &IgnoreMatcher{
		patterns: make([]ignorePattern, 0),
	}

	// Load default patterns
	m.loadDefaultPatterns()

	// Load .gitignore if it exists
	gitignorePath := filepath.Join(workspaceDir, ".gitignore")
	if _, err := os.Stat(gitignorePath); err == nil {
		if err := m.loadPatternsFromFile(gitignorePath, "gitignore"); err != nil {
			// Log warning but continue - don't fail on parse errors
			fmt.Fprintf(os.Stderr, "Warning: failed to parse .gitignore: %v\n", err)
		}
	}

	// Load .forgeignore if it exists
	forgeignorePath := filepath.Join(workspaceDir, ".forgeignore")
	if _, err := os.Stat(forgeignorePath); err == nil {
		if err := m.loadPatternsFromFile(forgeignorePath, "forgeignore"); err != nil {
			// Log warning but continue - don't fail on parse errors
			fmt.Fprintf(os.Stderr, "Warning: failed to parse .forgeignore: %v\n", err)
		}
	}

	return m, nil
}

// loadDefaultPatterns loads the hardcoded default ignore patterns.
func (m *IgnoreMatcher) loadDefaultPatterns() {
	for _, pattern := range defaultIgnorePatterns {
		m.addPattern(pattern, "default")
	}
}

// loadPatternsFromFile loads patterns from a gitignore-style file.
func (m *IgnoreMatcher) loadPatternsFromFile(path, source string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Add pattern with source information
		m.addPattern(line, source)
	}

	return scanner.Err()
}

// addPattern adds a pattern to the matcher with metadata.
func (m *IgnoreMatcher) addPattern(pattern, source string) {
	// Check for negation
	negation := false
	if strings.HasPrefix(pattern, "!") {
		negation = true
		pattern = pattern[1:]
	}

	// Check for directory-only pattern
	dirOnly := strings.HasSuffix(pattern, "/")
	if dirOnly {
		pattern = strings.TrimSuffix(pattern, "/")
	}

	// Check if pattern contains glob characters
	isGlob := strings.ContainsAny(pattern, "*?[]")

	m.patterns = append(m.patterns, ignorePattern{
		pattern:  pattern,
		negation: negation,
		dirOnly:  dirOnly,
		isGlob:   isGlob,
		source:   source,
	})
}

// ShouldIgnore checks if a path should be ignored based on loaded patterns.
// The path should be relative to the workspace root.
// Returns true if the path matches an ignore pattern (last match wins).
func (m *IgnoreMatcher) ShouldIgnore(relPath string, isDir bool) bool {
	// Normalize path separators for matching
	relPath = filepath.ToSlash(relPath)
	
	// Track whether path is ignored (last match wins)
	ignored := false

	// Check each pattern in order
	for _, p := range m.patterns {
		var matches bool

		// For directory-only patterns (ending with /), we need to check if:
		// 1. The path IS that directory (if isDir is true)
		// 2. The path is INSIDE that directory (for both files and dirs)
		if p.dirOnly {
			// Check if path matches the directory name
			dirMatches := m.matchPattern(relPath, p.pattern, p.isGlob)
			// Check if path is inside this directory
			insideDir := strings.HasPrefix(relPath, p.pattern+"/")
			
			matches = dirMatches || insideDir
		} else {
			matches = m.matchPattern(relPath, p.pattern, p.isGlob)
		}

		if matches {
			// If this is a negation pattern, unignore the path
			// Otherwise, ignore it
			ignored = !p.negation
		}
	}

	return ignored
}

// matchPattern checks if a path matches a pattern.
func (m *IgnoreMatcher) matchPattern(path, pattern string, isGlob bool) bool {
	// Normalize pattern separators
	pattern = filepath.ToSlash(pattern)

	// Check for exact match first
	if path == pattern {
		return true
	}

	// Split path into components
	parts := strings.Split(path, "/")

	if isGlob {
		// Use filepath.Match for glob patterns
		// Try matching against the base name
		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err == nil && matched {
			return true
		}

		// Try matching full path
		matched, err = filepath.Match(pattern, path)
		if err == nil && matched {
			return true
		}

		// Try matching against each path segment
		for i := range parts {
			subpath := strings.Join(parts[:i+1], "/")
			matched, err := filepath.Match(pattern, subpath)
			if err == nil && matched {
				return true
			}
			// Also try matching just the segment
			matched, err = filepath.Match(pattern, parts[i])
			if err == nil && matched {
				return true
			}
		}
	} else {
		// For non-glob patterns (like directory names)
		// Check if pattern matches any component in the path
		for _, part := range parts {
			if part == pattern {
				return true
			}
		}

		// Check if pattern is a prefix of the path (for directory patterns)
		if strings.HasPrefix(path, pattern+"/") || path == pattern {
			return true
		}

		// Check if any segment of the path starts with the pattern
		for i := range parts {
			subpath := strings.Join(parts[:i+1], "/")
			if subpath == pattern || strings.HasPrefix(path, pattern+"/") {
				return true
			}
		}
	}

	return false
}

// PatternCount returns the total number of loaded patterns.
// Useful for debugging and testing.
func (m *IgnoreMatcher) PatternCount() int {
	return len(m.patterns)
}

// Patterns returns a copy of all loaded patterns for debugging.
func (m *IgnoreMatcher) Patterns() []string {
	result := make([]string, len(m.patterns))
	for i, p := range m.patterns {
		prefix := ""
		if p.negation {
			prefix = "!"
		}
		suffix := ""
		if p.dirOnly {
			suffix = "/"
		}
		result[i] = fmt.Sprintf("%s%s%s [%s]", prefix, p.pattern, suffix, p.source)
	}
	return result
}