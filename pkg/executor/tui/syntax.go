package tui

import (
	"bytes"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/lipgloss"
)

// DiffLineType represents the type of line in a diff
type DiffLineType int

const (
	DiffLineContext DiffLineType = iota
	DiffLineAddition
	DiffLineDeletion
	DiffLineHeader
	DiffLineHunk
)

// DiffLine represents a parsed line from a unified diff
type DiffLine struct {
	Type    DiffLineType
	Content string
	Marker  string // The diff marker (+, -, space, @@)
}

// Color definitions for diff markers and backgrounds
var (
	diffAddColor       = lipgloss.Color("#90EE90") // Green for additions
	diffDeleteColor    = lipgloss.Color("#FFB3BA") // Red for deletions
	diffHunkColor      = lipgloss.Color("#87CEEB") // Cyan for hunk headers
	diffHeaderColor    = lipgloss.Color("#FFA07A") // Orange for file headers
	diffAddBgColor     = lipgloss.Color("#2d4a2b") // Dark green background for added lines
	diffDeleteBgColor  = lipgloss.Color("#4a2d2d") // Dark red background for deleted lines
)

// HighlightDiff applies syntax highlighting to unified diff content
// It preserves diff markers while highlighting the code portions
func HighlightDiff(diffContent, language string) (string, error) {
	if diffContent == "" {
		return "", nil
	}

	// Parse the diff into lines
	lines := parseDiffLines(diffContent)

	// Get the lexer for the language
	lexer := getLexerForLanguage(language)
	if lexer == nil {
		// If we can't get a lexer, return the original content with just diff coloring
		return applyDiffColorsOnly(lines), nil
	}

	// Use the terminal256 formatter for better color support
	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	// Use a terminal-friendly style
	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}

	// Process each line
	var result strings.Builder
	for _, line := range lines {
		highlightedLine, err := highlightDiffLine(line, lexer, formatter, style)
		if err != nil {
			// If highlighting fails, fall back to colored marker only
			highlightedLine = applyDiffColorToLine(line)
		}
		result.WriteString(highlightedLine)
		result.WriteString("\n")
	}

	return strings.TrimSuffix(result.String(), "\n"), nil
}

// parseDiffLines parses unified diff content into structured lines
func parseDiffLines(diffContent string) []DiffLine {
	lines := strings.Split(diffContent, "\n")
	result := make([]DiffLine, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			result = append(result, DiffLine{
				Type:    DiffLineContext,
				Content: "",
				Marker:  "",
			})
			continue
		}

		var lineType DiffLineType
		var marker string
		var content string

		switch {
		case strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---"):
			lineType = DiffLineHeader
			marker = line[:3]
			content = line[3:]
		case strings.HasPrefix(line, "@@"):
			lineType = DiffLineHunk
			// Find the end of the hunk header
			endIdx := strings.Index(line[2:], "@@")
			if endIdx != -1 {
				marker = line[:endIdx+4]
				content = line[endIdx+4:]
			} else {
				marker = line
				content = ""
			}
		case strings.HasPrefix(line, "+"):
			lineType = DiffLineAddition
			marker = "+"
			content = line[1:]
		case strings.HasPrefix(line, "-"):
			lineType = DiffLineDeletion
			marker = "-"
			content = line[1:]
		case strings.HasPrefix(line, " "):
			lineType = DiffLineContext
			marker = " "
			content = line[1:]
		default:
			// Treat unknown lines as context
			lineType = DiffLineContext
			marker = ""
			content = line
		}

		result = append(result, DiffLine{
			Type:    lineType,
			Content: content,
			Marker:  marker,
		})
	}

	return result
}

// highlightDiffLine highlights a single diff line
func highlightDiffLine(line DiffLine, lexer chroma.Lexer, formatter chroma.Formatter, style *chroma.Style) (string, error) {
	// For headers and hunks, just apply the appropriate color
	switch line.Type {
	case DiffLineHeader:
		return lipgloss.NewStyle().Foreground(diffHeaderColor).Render(line.Marker + line.Content), nil
	case DiffLineHunk:
		return lipgloss.NewStyle().Foreground(diffHunkColor).Render(line.Marker + line.Content), nil
	}

	// For code lines (additions, deletions, context), highlight the code portion
	if line.Content == "" {
		return applyDiffColorToLine(line), nil
	}

	// Tokenize and highlight the code content
	iterator, err := lexer.Tokenise(nil, line.Content)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		return "", err
	}

	highlightedContent := buf.String()
	// Remove trailing newline if present
	highlightedContent = strings.TrimSuffix(highlightedContent, "\n")

	// Build the complete line with marker, syntax-highlighted content, and background
	var marker string
	var contentWithBg string

	switch line.Type {
	case DiffLineAddition:
		// Styled marker with background
		marker = lipgloss.NewStyle().
			Foreground(diffAddColor).
			Background(diffAddBgColor).
			Bold(true).
			Render("+ ")
		// Apply background to content while preserving Chroma's syntax colors
		contentWithBg = lipgloss.NewStyle().
			Background(diffAddBgColor).
			Render(highlightedContent)
	case DiffLineDeletion:
		// Styled marker with background
		marker = lipgloss.NewStyle().
			Foreground(diffDeleteColor).
			Background(diffDeleteBgColor).
			Bold(true).
			Render("- ")
		// Apply background to content while preserving Chroma's syntax colors
		contentWithBg = lipgloss.NewStyle().
			Background(diffDeleteBgColor).
			Render(highlightedContent)
	default:
		marker = "  "
		contentWithBg = highlightedContent
	}

	// Combine marker with background-styled content
	return marker + contentWithBg, nil
}

// applyDiffColorsOnly applies only diff marker colors without syntax highlighting
func applyDiffColorsOnly(lines []DiffLine) string {
	var result strings.Builder
	for _, line := range lines {
		result.WriteString(applyDiffColorToLine(line))
		result.WriteString("\n")
	}
	return strings.TrimSuffix(result.String(), "\n")
}

// applyDiffColorToLine applies color to a single line based on its type
func applyDiffColorToLine(line DiffLine) string {
	fullLine := line.Marker + line.Content

	switch line.Type {
	case DiffLineHeader:
		return lipgloss.NewStyle().Foreground(diffHeaderColor).Render(fullLine)
	case DiffLineHunk:
		return lipgloss.NewStyle().Foreground(diffHunkColor).Render(fullLine)
	case DiffLineAddition:
		return lipgloss.NewStyle().Foreground(diffAddColor).Render(fullLine)
	case DiffLineDeletion:
		return lipgloss.NewStyle().Foreground(diffDeleteColor).Render(fullLine)
	default:
		return fullLine
	}
}

// getLexerForLanguage returns a Chroma lexer for the given language
func getLexerForLanguage(language string) chroma.Lexer {
	if language == "" {
		return nil
	}

	// Try to get lexer by name
	lexer := lexers.Get(language)
	if lexer != nil {
		return lexer
	}

	// Try common aliases
	aliases := map[string]string{
		"golang": "go",
		"js":     "javascript",
		"ts":     "typescript",
		"py":     "python",
		"rb":     "ruby",
		"sh":     "bash",
		"yml":    "yaml",
	}

	if alias, ok := aliases[strings.ToLower(language)]; ok {
		lexer = lexers.Get(alias)
		if lexer != nil {
			return lexer
		}
	}

	// Try to get lexer by filename extension
	lexer = lexers.Match("file." + language)
	if lexer != nil {
		return lexer
	}

	return nil
}

// HighlightCode highlights code without diff markers (useful for other contexts)
func HighlightCode(code, language string) (string, error) {
	if code == "" {
		return "", nil
	}

	lexer := getLexerForLanguage(language)
	if lexer == nil {
		// Return original code if we can't highlight it
		return code, nil
	}

	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return code, err
	}

	var buf bytes.Buffer
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		return code, err
	}

	return buf.String(), nil
}
