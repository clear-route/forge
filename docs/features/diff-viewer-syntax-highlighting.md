# Diff Viewer Syntax Highlighting

The Forge TUI's diff viewer now includes syntax highlighting powered by [Chroma](https://github.com/alecthomas/chroma), providing an enhanced code review experience with color-coded diffs.

## Features

### Syntax-Highlighted Diffs

When the [`apply_diff`](../../pkg/tools/coding/apply_diff.go) tool generates a preview, the diff viewer automatically applies syntax highlighting to:

- **Code content** - Syntax tokens colored according to the programming language
- **Diff markers** - Distinct colors for additions (+), deletions (-), and headers (@@)
- **File headers** - Clearly marked file paths and metadata

### Supported Languages

Chroma provides syntax highlighting for 200+ languages, including:

- Go, Python, JavaScript, TypeScript
- Rust, C, C++, Java, Kotlin
- Ruby, PHP, Perl
- Shell scripts (Bash, Zsh)
- Markup languages (HTML, XML, Markdown, YAML, JSON)
- And many more...

### Color Scheme

The diff viewer uses the **Monokai** color scheme optimized for terminal display:

- **Additions** (`+`) - Bright green (`#A8E6CF`)
- **Deletions** (`-`) - Bright red (`#FFB3BA`)  
- **Hunk headers** (`@@`) - Cyan (`#87CEEB`)
- **File headers** (`+++`/`---`) - Orange (`#FFA07A`)
- **Code syntax** - Monokai theme colors via Chroma's `terminal256` formatter

## Implementation Details

### Architecture

```
ApplyDiffTool.GeneratePreview()
    ↓
ToolPreview (with language metadata)
    ↓
DiffViewer.NewDiffViewer()
    ↓
syntax.HighlightDiff(content, language)
    ↓
Chroma tokenization & formatting
    ↓
ANSI-colored diff in viewport
```

### Language Detection

Language is detected automatically from file extensions via the [`detectLanguage()`](../../pkg/tools/coding/apply_diff.go:246) function in the `apply_diff` tool. The language is passed as metadata to the diff viewer.

### Graceful Degradation

The syntax highlighter includes robust fallback handling:

1. **Unknown language** - Falls back to diff-marker-only coloring
2. **Highlighting errors** - Returns original content with basic diff colors
3. **No language metadata** - Attempts to detect from file extension or uses plain text

## Components

### Core Files

- [`pkg/executor/tui/syntax.go`](../../pkg/executor/tui/syntax.go) - Syntax highlighting utility
- [`pkg/executor/tui/syntax_test.go`](../../pkg/executor/tui/syntax_test.go) - Comprehensive test suite
- [`pkg/executor/tui/diff_viewer.go`](../../pkg/executor/tui/diff_viewer.go) - DiffViewer component integration

### Key Functions

#### `HighlightDiff(diffContent, language string) (string, error)`

Main entry point for diff syntax highlighting.

**Parameters:**
- `diffContent` - Unified diff text
- `language` - Programming language identifier (e.g., "go", "python")

**Returns:**
- ANSI-colored diff string
- Error (falls back gracefully on error)

**Example:**
```go
highlighted, err := HighlightDiff(diffContent, "go")
if err != nil {
    // Falls back to original content
    highlighted = diffContent
}
```

#### `parseDiffLines(diffContent string) []DiffLine`

Parses unified diff into structured lines with type information.

**Returns:**
```go
type DiffLine struct {
    Type    DiffLineType  // Addition, Deletion, Context, Header, Hunk
    Content string        // The actual code/text
    Marker  string        // Diff marker (+, -, space, @@)
}
```

#### `getLexerForLanguage(language string) chroma.Lexer`

Retrieves appropriate Chroma lexer for a language, with alias support.

**Supported aliases:**
- `golang` → `go`
- `js` → `javascript`
- `ts` → `typescript`  
- `py` → `python`
- `rb` → `ruby`
- `sh` → `bash`
- `yml` → `yaml`

## Usage Example

The syntax highlighting is automatically applied when viewing diffs:

```
User: Update the function to add error handling

Agent: [Generates apply_diff preview]
    
┌─ Tool Approval Required ─────────────────────┐
│ apply_diff: Apply 1 edit to main.go          │
│                                               │
│ --- a/main.go                                 │
│ +++ b/main.go                                 │
│ @@ -10,5 +10,8 @@ func process(data string) {│
│   func process(data string) {                 │
│ -   return result                             │
│ +   if err != nil {                           │
│ +       return nil, err                       │
│ +   }                                         │
│ +   return result, nil                        │
│   }                                           │
│                                               │
│ [✓ Accept] [✗ Reject]                        │
└───────────────────────────────────────────────┘
```

All syntax tokens (keywords, strings, etc.) are colored according to the Monokai theme, while diff markers remain distinctly colored for easy scanning.

## Performance

- **On-demand highlighting** - Applied only when diff viewer opens
- **Line-by-line processing** - Efficient memory usage for large diffs
- **Cached lexers** - Chroma reuses lexer instances

## Testing

The syntax highlighting includes comprehensive tests covering:

- ✅ Diff line parsing (all types: additions, deletions, headers, hunks)
- ✅ Language detection (Go, Python, JavaScript, TypeScript, etc.)
- ✅ Highlighting with multiple languages
- ✅ Fallback for unknown languages
- ✅ Edge cases (empty content, invalid patterns)

Run tests:
```bash
go test ./pkg/executor/tui -v -run TestSyntax
```

## Future Enhancements

Potential improvements for future versions:

- [ ] Configurable color schemes (light/dark themes)
- [ ] Side-by-side diff view with synchronized highlighting
- [ ] Custom syntax highlighting rules
- [ ] Performance optimization for very large diffs (>1000 lines)

## References

- [Chroma Documentation](https://github.com/alecthomas/chroma)
- [ADR-0012: Enhanced TUI Executor](../adr/0012-enhanced-tui-executor.md)
- [Tool Approval Mechanism](../adr/0010-tool-approval-mechanism.md)