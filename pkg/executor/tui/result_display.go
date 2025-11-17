// Package tui provides result display strategies for tool outputs
package tui

import (
	"fmt"
	"strings"
)

// DisplayTier represents how a tool result should be displayed in the TUI
type DisplayTier int

const (
	// TierFullInline displays the complete result inline (loop-breaking tools)
	TierFullInline DisplayTier = iota
	// TierSummaryWithPreview displays summary + first few lines
	TierSummaryWithPreview
	// TierSummaryOnly displays only a summary line
	TierSummaryOnly
	// TierOverlayOnly displays nothing inline (handled by overlay)
	TierOverlayOnly
)

// ToolResultClassifier determines how to display a tool result
type ToolResultClassifier struct {
	// Loop-breaking tools (always full inline)
	loopBreakingTools map[string]bool
	// Size threshold for summary vs preview (in lines)
	summaryThreshold int
	// Number of preview lines to show for Tier 2
	previewLines int
}

// NewToolResultClassifier creates a new classifier with default settings
func NewToolResultClassifier() *ToolResultClassifier {
	return &ToolResultClassifier{
		loopBreakingTools: map[string]bool{
			"task_completion": true,
			"ask_question":    true,
			"converse":        true,
		},
		summaryThreshold: 100, // Results >= 100 lines get summary only
		previewLines:     3,   // Show first 3 lines for preview tier
	}
}

// ClassifyToolResult determines the display tier for a tool result
func (c *ToolResultClassifier) ClassifyToolResult(toolName string, result string) DisplayTier {
	// Check if this is a loop-breaking tool (always full inline)
	if c.loopBreakingTools[toolName] {
		return TierFullInline
	}

	// Command execution uses overlay (already implemented)
	if toolName == "execute_command" {
		return TierOverlayOnly
	}

	// Count lines in result
	lineCount := strings.Count(result, "\n") + 1

	// Large results get summary only
	if lineCount >= c.summaryThreshold {
		return TierSummaryOnly
	}

	// Check for specific tools that should always be summary-only when large
	switch toolName {
	case "read_file", "search_files", "list_files":
		if lineCount >= 50 {
			return TierSummaryOnly
		}
	case "write_file", "apply_diff":
		if lineCount >= 50 {
			return TierSummaryOnly
		}
	}

	// Medium results get summary + preview
	if lineCount >= 20 {
		return TierSummaryWithPreview
	}

	// Small results get full inline
	return TierFullInline
}

// GetPreviewLines extracts the first N lines from a result
func (c *ToolResultClassifier) GetPreviewLines(result string) string {
	lines := strings.Split(result, "\n")
	if len(lines) <= c.previewLines {
		return result
	}

	preview := strings.Join(lines[:c.previewLines], "\n")
	remainingLines := len(lines) - c.previewLines
	return fmt.Sprintf("%s\n  ... [%d more lines - Ctrl+R to view]", preview, remainingLines)
}

// ToolResultSummarizer generates summaries for tool results
type ToolResultSummarizer struct{}

// NewToolResultSummarizer creates a new summarizer
func NewToolResultSummarizer() *ToolResultSummarizer {
	return &ToolResultSummarizer{}
}

// GenerateSummary creates a one-line summary for a tool result
func (s *ToolResultSummarizer) GenerateSummary(toolName string, result string) string {
	lineCount := strings.Count(result, "\n") + 1
	sizeKB := float64(len(result)) / 1024.0

	switch toolName {
	case "read_file":
		// Extract filename from result if possible (first line often has it in comments)
		filename := s.extractFilename(result)
		if filename != "" {
			return fmt.Sprintf("Read %d lines from %s (%.1f KB) [Ctrl+R to view]", lineCount, filename, sizeKB)
		}
		return fmt.Sprintf("Read %d lines (%.1f KB) [Ctrl+R to view]", lineCount, sizeKB)

	case "search_files":
		matchCount, fileCount := s.parseSearchResults(result)
		return fmt.Sprintf("Found %d matches in %d files [Ctrl+R to view]", matchCount, fileCount)

	case "list_files":
		fileCount, dirCount := s.parseListResults(result)
		return fmt.Sprintf("Listed %d files and %d directories [Ctrl+R to view]", fileCount, dirCount)

	case "write_file":
		filename := s.extractFilename(result)
		if filename != "" {
			return fmt.Sprintf("Wrote %d lines to %s (%.1f KB)", lineCount, filename, sizeKB)
		}
		return fmt.Sprintf("Wrote %d lines (%.1f KB)", lineCount, sizeKB)

	case "apply_diff":
		editCount := s.parseApplyDiffResults(result)
		filename := s.extractFilename(result)
		if filename != "" && editCount > 0 {
			return fmt.Sprintf("Applied %d edits to %s", editCount, filename)
		}
		return fmt.Sprintf("Applied changes (%.1f KB)", sizeKB)

	default:
		// Generic summary for unknown tools
		return fmt.Sprintf("%s completed (%d lines, %.1f KB) [Ctrl+R to view]", toolName, lineCount, sizeKB)
	}
}

// extractFilename attempts to extract a filename from tool result
func (s *ToolResultSummarizer) extractFilename(result string) string {
	// This is a simple heuristic - improve as needed
	lines := strings.Split(result, "\n")
	if len(lines) == 0 {
		return ""
	}

	// Look for common patterns in first few lines
	for i := 0; i < min(3, len(lines)); i++ {
		line := strings.TrimSpace(lines[i])
		// Look for file paths (simple heuristic)
		if strings.Contains(line, ".go") || strings.Contains(line, ".py") ||
			strings.Contains(line, ".js") || strings.Contains(line, ".ts") ||
			strings.Contains(line, "/") {
			// Extract just the filename if it's a path
			parts := strings.Split(line, "/")
			if len(parts) > 0 {
				return parts[len(parts)-1]
			}
			return line
		}
	}
	return ""
}

// parseSearchResults extracts match and file counts from search results
func (s *ToolResultSummarizer) parseSearchResults(result string) (matchCount, fileCount int) {
	lines := strings.Split(result, "\n")
	files := make(map[string]bool)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Count matches (non-empty lines are typically matches)
		matchCount++

		// Extract filename before colon (format: "file.go:123: match")
		if colonIdx := strings.Index(line, ":"); colonIdx > 0 {
			filename := line[:colonIdx]
			files[filename] = true
		}
	}

	fileCount = len(files)
	return matchCount, fileCount
}

// parseListResults extracts file and directory counts from list results
func (s *ToolResultSummarizer) parseListResults(result string) (fileCount, dirCount int) {
	lines := strings.Split(result, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "📁") {
			dirCount++
		} else if strings.HasPrefix(line, "📄") {
			fileCount++
		}
	}

	return fileCount, dirCount
}

// parseApplyDiffResults extracts edit count from apply_diff results
func (s *ToolResultSummarizer) parseApplyDiffResults(result string) int {
	// Look for patterns like "Applied 5 edits" in the result
	if strings.Contains(result, "Applied") && strings.Contains(result, "edit") {
		// Simple count of how many "edit" or "change" mentions
		return strings.Count(result, "edit")
	}
	// Fallback: count non-empty lines as edits
	lines := strings.Split(result, "\n")
	count := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			count++
		}
	}
	return count
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
