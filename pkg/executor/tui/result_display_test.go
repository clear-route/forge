package tui

import (
	"strings"
	"testing"
)

func TestExtractFilename(t *testing.T) {
	s := NewToolResultSummarizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "quoted path with directory",
			input:    `Read file "pkg/executor/tui/result_display.go"`,
			expected: "result_display.go",
		},
		{
			name:     "single quoted path",
			input:    `Modified 'internal/agent/handler.py'`,
			expected: "handler.py",
		},
		{
			name:     "backtick quoted path",
			input:    "Processing file `src/main.ts`",
			expected: "main.ts",
		},
		{
			name:     "path with slash separators",
			input:    "Reading from pkg/config/settings.go line 42",
			expected: "settings.go",
		},
		{
			name:     "relative path with ./",
			input:    "Found in ./cmd/forge/main.go",
			expected: "main.go",
		},
		{
			name:     "absolute path",
			input:    "Loading /etc/nginx/nginx.conf",
			expected: "nginx.conf",
		},
		{
			name:     "filename with extension only",
			input:    "Processing config.yaml",
			expected: "config.yaml",
		},
		{
			name:     "go file in text",
			input:    "The main.go file contains the entry point",
			expected: "main.go",
		},
		{
			name:     "python file",
			input:    "Running test_handler.py",
			expected: "test_handler.py",
		},
		{
			name:     "typescript file",
			input:    "Compiling components/Button.tsx",
			expected: "Button.tsx",
		},
		{
			name:     "config file - json",
			input:    "Loading package.json",
			expected: "package.json",
		},
		{
			name:     "config file - toml",
			input:    "Reading Cargo.toml",
			expected: "Cargo.toml",
		},
		{
			name:     "markdown file",
			input:    "Rendering README.md",
			expected: "README.md",
		},
		{
			name:     "dockerfile",
			input:    "Building from Dockerfile",
			expected: "Dockerfile",
		},
		{
			name:     "makefile",
			input:    "Executing Makefile",
			expected: "Makefile",
		},
		{
			name:     "shell script",
			input:    "Running deploy.sh",
			expected: "deploy.sh",
		},
		{
			name:     "no filename present",
			input:    "This is just some text without any file reference",
			expected: "",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "multiline with filename in second line",
			input:    "Successfully completed\nModified src/utils/helper.js\nAll done",
			expected: "helper.js",
		},
		{
			name:     "filename with hyphens",
			input:    "Reading my-config-file.yaml",
			expected: "my-config-file.yaml",
		},
		{
			name:     "filename with underscores",
			input:    "Processing test_utils_helper.py",
			expected: "test_utils_helper.py",
		},
		{
			name:     "c++ header file",
			input:    "Including algorithm.hpp",
			expected: "algorithm.hpp",
		},
		{
			name:     "rust file",
			input:    "Compiling main.rs",
			expected: "main.rs",
		},
		{
			name:     "java file",
			input:    "Running Application.java",
			expected: "Application.java",
		},
		{
			name:     "scala file",
			input:    "Processing Handler.scala",
			expected: "Handler.scala",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.extractFilename(tt.input)
			if result != tt.expected {
				t.Errorf("extractFilename(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractFromQuotedPath(t *testing.T) {
	s := NewToolResultSummarizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "double quotes",
			input:    `"pkg/executor/main.go"`,
			expected: "main.go",
		},
		{
			name:     "single quotes",
			input:    `'internal/config.yaml'`,
			expected: "config.yaml",
		},
		{
			name:     "backticks",
			input:    "`src/handler.ts`",
			expected: "handler.ts",
		},
		{
			name:     "quoted path in sentence",
			input:    `The file "pkg/utils/helper.py" was modified`,
			expected: "helper.py",
		},
		{
			name:     "no quotes",
			input:    "pkg/main.go",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.extractFromQuotedPath(tt.input)
			if result != tt.expected {
				t.Errorf("extractFromQuotedPath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractFromPathPattern(t *testing.T) {
	s := NewToolResultSummarizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "relative path with pkg",
			input:    "pkg/executor/tui/handler.go:42",
			expected: "handler.go",
		},
		{
			name:     "relative path with ./",
			input:    "./cmd/forge/main.go",
			expected: "main.go",
		},
		{
			name:     "nested path",
			input:    "internal/agent/tools/executor.py",
			expected: "executor.py",
		},
		{
			name:     "path in sentence",
			input:    "Found in src/components/Button.tsx at line 10",
			expected: "Button.tsx",
		},
		{
			name:     "no path separators",
			input:    "config.yaml",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.extractFromPathPattern(tt.input)
			if result != tt.expected {
				t.Errorf("extractFromPathPattern(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractFromExtensionPattern(t *testing.T) {
	s := NewToolResultSummarizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "go file",
			input:    "Processing main.go",
			expected: "main.go",
		},
		{
			name:     "python file",
			input:    "Running test.py",
			expected: "test.py",
		},
		{
			name:     "typescript file",
			input:    "Compiling index.ts",
			expected: "index.ts",
		},
		{
			name:     "config file",
			input:    "Loading config.yaml",
			expected: "config.yaml",
		},
		{
			name:     "multiple files, gets first",
			input:    "Processing main.go and test.py",
			expected: "main.go",
		},
		{
			name:     "file with hyphens",
			input:    "Reading my-config.toml",
			expected: "my-config.toml",
		},
		{
			name:     "file with underscores",
			input:    "Running test_utils.py",
			expected: "test_utils.py",
		},
		{
			name:     "no recognizable extension",
			input:    "Processing data.xyz",
			expected: "",
		},
		{
			name:     "no filename",
			input:    "Just some text",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.extractFromExtensionPattern(tt.input)
			if result != tt.expected {
				t.Errorf("extractFromExtensionPattern(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGenerateSummary(t *testing.T) {
	s := NewToolResultSummarizer()

	tests := []struct {
		name     string
		toolName string
		result   string
		contains []string // strings that should be in the summary
	}{
		{
			name:     "read_file with filename",
			toolName: "read_file",
			result:   "Reading pkg/config/settings.go\nline 1\nline 2\nline 3",
			contains: []string{"Read", "lines", "settings.go", "Ctrl+R"},
		},
		{
			name:     "search_files",
			toolName: "search_files",
			result:   "file1.go:10: match\nfile2.py:20: match\nfile3.js:30: match",
			contains: []string{"Found", "matches", "files", "Ctrl+R"},
		},
		{
			name:     "list_files",
			toolName: "list_files",
			result:   "📁 dir1\n📄 file1.go\n📄 file2.py\n📁 dir2",
			contains: []string{"Listed", "files", "directories", "Ctrl+R"},
		},
		{
			name:     "write_file with filename",
			toolName: "write_file",
			result:   "Wrote to output.txt\nContent here",
			contains: []string{"Wrote", "lines", "output.txt"},
		},
		{
			name:     "apply_diff with filename",
			toolName: "apply_diff",
			result:   "Applied edits to main.go\nedit 1\nedit 2",
			contains: []string{"Applied", "edits", "main.go"},
		},
		{
			name:     "generic tool",
			toolName: "unknown_tool",
			result:   "Some result\nwith multiple\nlines",
			contains: []string{"unknown_tool", "completed", "lines", "Ctrl+R"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := s.GenerateSummary(tt.toolName, tt.result)
			for _, expected := range tt.contains {
				if !contains(summary, expected) {
					t.Errorf("GenerateSummary(%q, ...) = %q, expected to contain %q",
						tt.toolName, summary, expected)
				}
			}
		})
	}
}

// contains checks if a string contains a substring (case-sensitive)
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
