package overlay

import (
	"testing"
	"time"

	"github.com/entrhq/forge/pkg/executor/tui/types"
)

func TestResultListItem_Description(t *testing.T) {
	tests := []struct {
		name        string
		summary     string
		expected    string
		shouldTrunc bool
	}{
		{
			name:        "short summary - no truncation",
			summary:     "Short summary",
			expected:    "Short summary",
			shouldTrunc: false,
		},
		{
			name:        "exactly 77 characters - no truncation",
			summary:     "1234567890123456789012345678901234567890123456789012345678901234567890123456",
			expected:    "1234567890123456789012345678901234567890123456789012345678901234567890123456",
			shouldTrunc: false,
		},
		{
			name:        "exactly 78 characters - should truncate",
			summary:     "123456789012345678901234567890123456789012345678901234567890123456789012345678",
			expected:    "12345678901234567890123456789012345678901234567890123456789012345678901234567...",
			shouldTrunc: true,
		},
		{
			name:        "long summary - should truncate",
			summary:     "This is a very long summary that exceeds the maximum display width and should be truncated with ellipsis",
			expected:    "This is a very long summary that exceeds the maximum display width and should...",
			shouldTrunc: true,
		},
		{
			name:        "empty summary",
			summary:     "",
			expected:    "",
			shouldTrunc: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := resultListItem{
				result: &types.CachedResult{
					Summary: tt.summary,
				},
			}

			result := item.Description()

			if result != tt.expected {
				t.Errorf("Description() = %q (len=%d), want %q (len=%d)",
					result, len(result), tt.expected, len(tt.expected))
			}
		})
	}
}

func TestResultListItem_Title(t *testing.T) {
	tests := []struct {
		name     string
		toolName string
		expected string
	}{
		{
			name:     "read_file tool",
			toolName: "read_file",
			expected: "read_file",
		},
		{
			name:     "empty tool name",
			toolName: "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := resultListItem{
				result: &types.CachedResult{
					ToolName:  tt.toolName,
					Timestamp: time.Now(),
				},
			}

			result := item.Title()
			// Title includes timestamp, so just verify it contains the tool name
			if tt.toolName != "" && result[:len(tt.toolName)] != tt.toolName {
				t.Errorf("Title() should start with %q, got %q", tt.toolName, result)
			}
		})
	}
}

func TestResultListItem_FilterValue(t *testing.T) {
	tests := []struct {
		name     string
		toolName string
		expected string
	}{
		{
			name:     "returns tool name for filtering",
			toolName: "read_file",
			expected: "read_file",
		},
		{
			name:     "empty tool name",
			toolName: "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := resultListItem{
				result: &types.CachedResult{
					ToolName: tt.toolName,
				},
			}

			result := item.FilterValue()
			if result != tt.expected {
				t.Errorf("FilterValue() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestResultListModel_Activity(t *testing.T) {
	m := NewResultListModel()
	if m.IsActive() {
		t.Error("New model should not be active")
	}

	m.Activate([]*types.CachedResult{}, 100, 100)
	if !m.IsActive() {
		t.Error("Model should be active after Activate")
	}

	m.Deactivate()
	if m.IsActive() {
		t.Error("Model should not be active after Deactivate")
	}
}
