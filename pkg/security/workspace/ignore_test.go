package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultIgnorePatterns(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	matcher, err := NewIgnoreMatcher(tempDir)
	if err != nil {
		t.Fatalf("Failed to create matcher: %v", err)
	}

	tests := []struct {
		name     string
		path     string
		isDir    bool
		expected bool
	}{
		{"node_modules directory", "node_modules", true, true},
		{"node_modules file in root", "node_modules/package.json", false, true},
		{"nested node_modules", "src/node_modules/file.js", false, true},
		{".git directory", ".git", true, true},
		{".git file", ".git/config", false, true},
		{".env file", ".env", false, true},
		{".env.local file", ".env.local", false, true},
		{".env.production", ".env.production", false, true},
		{"log file", "app.log", false, true},
		{"nested log file", "logs/app.log", false, true},
		{".DS_Store", ".DS_Store", false, true},
		{"vendor directory", "vendor", true, true},
		{"vendor file", "vendor/package.go", false, true},
		{"__pycache__ directory", "__pycache__", true, true},
		{".pyc file", "module.pyc", false, true},
		{".vscode directory", ".vscode", true, true},
		{".idea directory", ".idea", true, true},
		{"dist directory", "dist", true, true},
		{"build directory", "build", true, true},
		{"tmp directory", "tmp", true, true},
		{"temp directory", "temp", true, true},
		{"coverage directory", "coverage", true, true},
		{".next directory", ".next", true, true},
		{".nuxt directory", ".nuxt", true, true},
		{"target directory", "target", true, true},
		{"swap file", "file.swp", false, true},
		{"normal go file", "main.go", false, false},
		{"normal directory", "src", true, false},
		{"normal nested file", "src/app/main.go", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matcher.ShouldIgnore(tt.path, tt.isDir)
			if result != tt.expected {
				t.Errorf("ShouldIgnore(%q, %v) = %v, want %v", tt.path, tt.isDir, result, tt.expected)
			}
		})
	}
}

func TestGitignorePatterns(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a .gitignore file
	gitignoreContent := `# Comment
*.test
/build/
debug.log
!important.log
`
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		t.Fatalf("Failed to create .gitignore: %v", err)
	}

	matcher, err := NewIgnoreMatcher(tempDir)
	if err != nil {
		t.Fatalf("Failed to create matcher: %v", err)
	}

	tests := []struct {
		name     string
		path     string
		isDir    bool
		expected bool
	}{
		{"test file", "app.test", false, true},
		{"debug log", "debug.log", false, true},
		{"important log (negated)", "important.log", false, false},
		{"build directory", "build", true, true},
		{"normal file", "app.go", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matcher.ShouldIgnore(tt.path, tt.isDir)
			if result != tt.expected {
				t.Errorf("ShouldIgnore(%q, %v) = %v, want %v", tt.path, tt.isDir, result, tt.expected)
			}
		})
	}
}

func TestForgeignorePatterns(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create both .gitignore and .forgeignore
	gitignoreContent := `*.test`
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		t.Fatalf("Failed to create .gitignore: %v", err)
	}

	forgeignoreContent := `!app.test
*.ignore
`
	forgeignorePath := filepath.Join(tempDir, ".forgeignore")
	if err := os.WriteFile(forgeignorePath, []byte(forgeignoreContent), 0644); err != nil {
		t.Fatalf("Failed to create .forgeignore: %v", err)
	}

	matcher, err := NewIgnoreMatcher(tempDir)
	if err != nil {
		t.Fatalf("Failed to create matcher: %v", err)
	}

	tests := []struct {
		name     string
		path     string
		isDir    bool
		expected bool
	}{
		{"app.test (negated in forgeignore)", "app.test", false, false},
		{"other.test (from gitignore)", "other.test", false, true},
		{"file.ignore (from forgeignore)", "file.ignore", false, true},
		{"normal file", "app.go", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matcher.ShouldIgnore(tt.path, tt.isDir)
			if result != tt.expected {
				t.Errorf("ShouldIgnore(%q, %v) = %v, want %v", tt.path, tt.isDir, result, tt.expected)
			}
		})
	}
}

func TestPatternPrecedence(t *testing.T) {
	// Test that .forgeignore > .gitignore > defaults
	tempDir := t.TempDir()

	// .gitignore ignores *.log
	gitignoreContent := `*.log`
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		t.Fatalf("Failed to create .gitignore: %v", err)
	}

	// .forgeignore negates app.log
	forgeignoreContent := `!app.log`
	forgeignorePath := filepath.Join(tempDir, ".forgeignore")
	if err := os.WriteFile(forgeignorePath, []byte(forgeignoreContent), 0644); err != nil {
		t.Fatalf("Failed to create .forgeignore: %v", err)
	}

	matcher, err := NewIgnoreMatcher(tempDir)
	if err != nil {
		t.Fatalf("Failed to create matcher: %v", err)
	}

	// app.log should NOT be ignored (forgeignore negation wins)
	if matcher.ShouldIgnore("app.log", false) {
		t.Error("Expected app.log to not be ignored due to .forgeignore negation")
	}

	// other.log should be ignored (gitignore pattern)
	if !matcher.ShouldIgnore("other.log", false) {
		t.Error("Expected other.log to be ignored due to .gitignore pattern")
	}
}

func TestDirectoryOnlyPatterns(t *testing.T) {
	tempDir := t.TempDir()

	// Create .forgeignore with directory-only pattern
	forgeignoreContent := `test/`
	forgeignorePath := filepath.Join(tempDir, ".forgeignore")
	if err := os.WriteFile(forgeignorePath, []byte(forgeignoreContent), 0644); err != nil {
		t.Fatalf("Failed to create .forgeignore: %v", err)
	}

	matcher, err := NewIgnoreMatcher(tempDir)
	if err != nil {
		t.Fatalf("Failed to create matcher: %v", err)
	}

	// Directory should be ignored
	if !matcher.ShouldIgnore("test", true) {
		t.Error("Expected 'test' directory to be ignored")
	}

	// Files inside the directory should also be ignored
	if !matcher.ShouldIgnore("test/file.txt", false) {
		t.Error("Expected files inside 'test/' directory to be ignored")
	}
}

func TestEmptyAndCommentLines(t *testing.T) {
	tempDir := t.TempDir()

	// Create .gitignore with empty lines and comments
	gitignoreContent := `
# This is a comment
*.test

# Another comment
*.log
`
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		t.Fatalf("Failed to create .gitignore: %v", err)
	}

	matcher, err := NewIgnoreMatcher(tempDir)
	if err != nil {
		t.Fatalf("Failed to create matcher: %v", err)
	}

	// Should ignore .test and .log files
	if !matcher.ShouldIgnore("app.test", false) {
		t.Error("Expected .test files to be ignored")
	}
	if !matcher.ShouldIgnore("app.log", false) {
		t.Error("Expected .log files to be ignored")
	}
}

func TestPatternCount(t *testing.T) {
	tempDir := t.TempDir()

	// Create .gitignore with 2 patterns
	gitignoreContent := `*.test
*.log
`
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		t.Fatalf("Failed to create .gitignore: %v", err)
	}

	matcher, err := NewIgnoreMatcher(tempDir)
	if err != nil {
		t.Fatalf("Failed to create matcher: %v", err)
	}

	// Should have default patterns + 2 from gitignore
	expectedMin := len(defaultIgnorePatterns) + 2
	actualCount := matcher.PatternCount()
	if actualCount < expectedMin {
		t.Errorf("Expected at least %d patterns, got %d", expectedMin, actualCount)
	}
}

func TestMissingIgnoreFiles(t *testing.T) {
	// Test that matcher works even without .gitignore or .forgeignore
	tempDir := t.TempDir()

	matcher, err := NewIgnoreMatcher(tempDir)
	if err != nil {
		t.Fatalf("Failed to create matcher: %v", err)
	}

	// Should still have default patterns
	if matcher.PatternCount() != len(defaultIgnorePatterns) {
		t.Errorf("Expected %d default patterns, got %d", len(defaultIgnorePatterns), matcher.PatternCount())
	}

	// Should ignore node_modules (default pattern)
	if !matcher.ShouldIgnore("node_modules", true) {
		t.Error("Expected node_modules to be ignored by default patterns")
	}
}

func TestGlobPatterns(t *testing.T) {
	tempDir := t.TempDir()

	// Create .gitignore with various glob patterns
	gitignoreContent := `*.test
test_*.go
*.tmp
`
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		t.Fatalf("Failed to create .gitignore: %v", err)
	}

	matcher, err := NewIgnoreMatcher(tempDir)
	if err != nil {
		t.Fatalf("Failed to create matcher: %v", err)
	}

	tests := []struct {
		name     string
		path     string
		isDir    bool
		expected bool
	}{
		{"simple glob", "app.test", false, true},
		{"prefix glob", "test_main.go", false, true},
		{"tmp file", "cache.tmp", false, true},
		{"nested tmp", "src/cache.tmp", false, true},
		{"no match", "main.go", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matcher.ShouldIgnore(tt.path, tt.isDir)
			if result != tt.expected {
				t.Errorf("ShouldIgnore(%q, %v) = %v, want %v", tt.path, tt.isDir, result, tt.expected)
			}
		})
	}
}
