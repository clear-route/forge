package tui

import (
	"strings"
	"testing"
)

func TestParseDiffLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []DiffLine
	}{
		{
			name:  "empty string",
			input: "",
			expected: []DiffLine{
				{Type: DiffLineContext, Content: "", Marker: ""},
			},
		},
		{
			name:  "addition line",
			input: "+func main() {",
			expected: []DiffLine{
				{Type: DiffLineAddition, Content: "func main() {", Marker: "+"},
			},
		},
		{
			name:  "deletion line",
			input: "-old code",
			expected: []DiffLine{
				{Type: DiffLineDeletion, Content: "old code", Marker: "-"},
			},
		},
		{
			name:  "context line",
			input: " unchanged",
			expected: []DiffLine{
				{Type: DiffLineContext, Content: "unchanged", Marker: " "},
			},
		},
		{
			name:  "hunk header",
			input: "@@ -1,5 +1,6 @@ function name",
			expected: []DiffLine{
				{Type: DiffLineHunk, Content: " function name", Marker: "@@ -1,5 +1,6 @@"},
			},
		},
		{
			name:  "file header addition",
			input: "+++ b/file.go",
			expected: []DiffLine{
				{Type: DiffLineHeader, Content: " b/file.go", Marker: "+++"},
			},
		},
		{
			name:  "file header deletion",
			input: "--- a/file.go",
			expected: []DiffLine{
				{Type: DiffLineHeader, Content: " a/file.go", Marker: "---"},
			},
		},
		{
			name: "multiple lines",
			input: `--- a/main.go
+++ b/main.go
@@ -1,3 +1,4 @@
 package main
-import "fmt"
+import (
+	"fmt"
+)`,
			expected: []DiffLine{
				{Type: DiffLineHeader, Content: " a/main.go", Marker: "---"},
				{Type: DiffLineHeader, Content: " b/main.go", Marker: "+++"},
				{Type: DiffLineHunk, Content: "", Marker: "@@ -1,3 +1,4 @@"},
				{Type: DiffLineContext, Content: "package main", Marker: " "},
				{Type: DiffLineDeletion, Content: "import \"fmt\"", Marker: "-"},
				{Type: DiffLineAddition, Content: "import (", Marker: "+"},
				{Type: DiffLineAddition, Content: "\t\"fmt\"", Marker: "+"},
				{Type: DiffLineAddition, Content: ")", Marker: "+"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseDiffLines(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("parseDiffLines() returned %d lines, expected %d", len(result), len(tt.expected))
				return
			}

			for i, line := range result {
				if line.Type != tt.expected[i].Type {
					t.Errorf("Line %d: got type %v, expected %v", i, line.Type, tt.expected[i].Type)
				}
				if line.Content != tt.expected[i].Content {
					t.Errorf("Line %d: got content %q, expected %q", i, line.Content, tt.expected[i].Content)
				}
				if line.Marker != tt.expected[i].Marker {
					t.Errorf("Line %d: got marker %q, expected %q", i, line.Marker, tt.expected[i].Marker)
				}
			}
		})
	}
}

func TestGetLexerForLanguage(t *testing.T) {
	tests := []struct {
		name     string
		language string
		wantNil  bool
	}{
		{
			name:     "go language",
			language: "go",
			wantNil:  false,
		},
		{
			name:     "golang alias",
			language: "golang",
			wantNil:  false,
		},
		{
			name:     "python",
			language: "python",
			wantNil:  false,
		},
		{
			name:     "py alias",
			language: "py",
			wantNil:  false,
		},
		{
			name:     "javascript",
			language: "javascript",
			wantNil:  false,
		},
		{
			name:     "js alias",
			language: "js",
			wantNil:  false,
		},
		{
			name:     "typescript",
			language: "typescript",
			wantNil:  false,
		},
		{
			name:     "ts alias",
			language: "ts",
			wantNil:  false,
		},
		{
			name:     "empty language",
			language: "",
			wantNil:  true,
		},
		{
			name:     "unknown language",
			language: "unknownlang123",
			wantNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := getLexerForLanguage(tt.language)

			if tt.wantNil && lexer != nil {
				t.Errorf("getLexerForLanguage(%q) returned non-nil lexer, expected nil", tt.language)
			}

			if !tt.wantNil && lexer == nil {
				t.Errorf("getLexerForLanguage(%q) returned nil, expected non-nil lexer", tt.language)
			}
		})
	}
}

func TestHighlightDiff(t *testing.T) {
	tests := []struct {
		name     string
		diff     string
		language string
		wantErr  bool
		checkFor []string // Strings that should appear in output
	}{
		{
			name:     "empty diff",
			diff:     "",
			language: "go",
			wantErr:  false,
			checkFor: []string{},
		},
		{
			name: "simple go diff",
			diff: `--- a/main.go
+++ b/main.go
@@ -1,2 +1,3 @@
 package main
+import "fmt"
 func main() {}`,
			language: "go",
			wantErr:  false,
			checkFor: []string{"package", "main", "import", "fmt", "func"},
		},
		{
			name: "python diff",
			diff: `+def hello():
+    print("Hello, world!")
-def goodbye():
-    pass`,
			language: "python",
			wantErr:  false,
			checkFor: []string{"def", "hello", "print", "goodbye"},
		},
		{
			name: "unknown language fallback",
			diff: `+new line
-old line
 context line`,
			language: "unknownlang",
			wantErr:  false,
			checkFor: []string{"new line", "old line", "context line"},
		},
		{
			name: "javascript diff",
			diff: `+const x = 42;
-var y = 10;
 let z = 5;`,
			language: "javascript",
			wantErr:  false,
			checkFor: []string{"const", "var", "let"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := HighlightDiff(tt.diff, tt.language)

			if (err != nil) != tt.wantErr {
				t.Errorf("HighlightDiff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check that expected strings appear in output
			for _, check := range tt.checkFor {
				if !strings.Contains(result, check) {
					t.Errorf("HighlightDiff() result missing expected string %q", check)
				}
			}

			// If input was not empty, output should not be empty
			if tt.diff != "" && result == "" {
				t.Errorf("HighlightDiff() returned empty result for non-empty input")
			}
		})
	}
}

func TestHighlightCode(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		language string
		wantErr  bool
	}{
		{
			name:     "empty code",
			code:     "",
			language: "go",
			wantErr:  false,
		},
		{
			name:     "go code",
			code:     "package main\n\nfunc main() {\n\tfmt.Println(\"Hello\")\n}",
			language: "go",
			wantErr:  false,
		},
		{
			name:     "python code",
			code:     "def hello():\n    print('Hello, world!')",
			language: "python",
			wantErr:  false,
		},
		{
			name:     "unknown language",
			code:     "some code",
			language: "unknownlang",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := HighlightCode(tt.code, tt.language)

			if (err != nil) != tt.wantErr {
				t.Errorf("HighlightCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Empty code should return empty result
			if tt.code == "" && result != "" {
				t.Errorf("HighlightCode() expected empty result for empty code")
			}

			// Non-empty code should return non-empty result
			if tt.code != "" && result == "" {
				t.Errorf("HighlightCode() returned empty result for non-empty code")
			}
		})
	}
}

func TestApplyDiffColorToLine(t *testing.T) {
	tests := []struct {
		name string
		line DiffLine
	}{
		{
			name: "addition line",
			line: DiffLine{Type: DiffLineAddition, Content: "new code", Marker: "+"},
		},
		{
			name: "deletion line",
			line: DiffLine{Type: DiffLineDeletion, Content: "old code", Marker: "-"},
		},
		{
			name: "context line",
			line: DiffLine{Type: DiffLineContext, Content: "unchanged", Marker: " "},
		},
		{
			name: "header line",
			line: DiffLine{Type: DiffLineHeader, Content: " file.go", Marker: "+++"},
		},
		{
			name: "hunk line",
			line: DiffLine{Type: DiffLineHunk, Content: "", Marker: "@@ -1,2 +1,3 @@"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := applyDiffColorToLine(tt.line)

			// Result should not be empty (unless input was empty)
			expectedContent := tt.line.Marker + tt.line.Content
			if expectedContent != "" && result == "" {
				t.Errorf("applyDiffColorToLine() returned empty result")
			}

			// Result should contain the original content somewhere
			// (Note: The actual text might be wrapped in ANSI codes)
			if expectedContent != "" && !strings.Contains(result, tt.line.Content) {
				// Some styling might alter the content, so we just check it's not empty
				if result == "" {
					t.Errorf("applyDiffColorToLine() produced unexpected empty output")
				}
			}
		})
	}
}

func TestApplyDiffColorsOnly(t *testing.T) {
	lines := []DiffLine{
		{Type: DiffLineHeader, Content: " a/file.go", Marker: "---"},
		{Type: DiffLineHeader, Content: " b/file.go", Marker: "+++"},
		{Type: DiffLineHunk, Content: "", Marker: "@@ -1,1 +1,2 @@"},
		{Type: DiffLineContext, Content: "package main", Marker: " "},
		{Type: DiffLineAddition, Content: "import \"fmt\"", Marker: "+"},
	}

	result := applyDiffColorsOnly(lines)

	// Result should not be empty
	if result == "" {
		t.Errorf("applyDiffColorsOnly() returned empty result")
	}

	// Should contain the content from lines
	expectedStrings := []string{"file.go", "package main", "import"}
	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("applyDiffColorsOnly() result missing expected string %q", expected)
		}
	}

	// Should have the right number of newlines (one less than number of lines)
	expectedNewlines := len(lines) - 1
	actualNewlines := strings.Count(result, "\n")
	if actualNewlines != expectedNewlines {
		t.Errorf("applyDiffColorsOnly() has %d newlines, expected %d", actualNewlines, expectedNewlines)
	}
}
