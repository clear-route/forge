package core

import (
	"testing"
)

func TestExtractToolNameFromPartial(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "complete tool_name element",
			content:  "<tool>\n<server_name>local</server_name>\n<tool_name>apply_diff</tool_name>\n",
			expected: "apply_diff",
		},
		{
			name:     "tool_name with whitespace",
			content:  "<tool>\n<tool_name>  read_file  </tool_name>\n",
			expected: "read_file",
		},
		{
			name:     "incomplete tool_name - opening tag only",
			content:  "<tool>\n<tool_name>",
			expected: "",
		},
		{
			name:     "incomplete tool_name - no closing tag",
			content:  "<tool>\n<tool_name>write_file",
			expected: "",
		},
		{
			name:     "no tool_name element yet",
			content:  "<tool>\n<server_name>local</server_name>\n",
			expected: "",
		},
		{
			name:     "empty content",
			content:  "",
			expected: "",
		},
		{
			name:     "tool_name in middle of XML",
			content:  "<tool>\n<server_name>local</server_name>\n<tool_name>search_files</tool_name>\n<arguments>",
			expected: "search_files",
		},
		{
			name:     "malformed - empty tool_name",
			content:  "<tool_name></tool_name>",
			expected: "",
		},
		{
			name:     "tool_name with newlines",
			content:  "<tool_name>\napply_diff\n</tool_name>",
			expected: "apply_diff",
		},
		{
			name:     "tool_name split across chunks - incomplete",
			content:  "<tool>\n<server_name>local</server_name>\n<tool_na",
			expected: "",
		},
		{
			name:     "tool_name with underscores",
			content:  "<tool_name>execute_command</tool_name>",
			expected: "execute_command",
		},
		{
			name:     "multiple tags before tool_name",
			content:  "<tool>\n<server_name>local</server_name>\n<other>value</other>\n<tool_name>list_files</tool_name>",
			expected: "list_files",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractToolNameFromPartial(tt.content)
			if result != tt.expected {
				t.Errorf("extractToolNameFromPartial() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestExtractToolNameFromPartial_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "nested tags (should not match)",
			content:  "<tool_name><inner>test</inner></tool_name>",
			expected: "",
		},
		{
			name:     "tool_name with special characters (should not match)",
			content:  "<tool_name>tool<name</tool_name>",
			expected: "",
		},
		{
			name:     "very long tool name",
			content:  "<tool_name>this_is_a_very_long_tool_name_that_should_still_work</tool_name>",
			expected: "this_is_a_very_long_tool_name_that_should_still_work",
		},
		{
			name:     "tool_name with mixed case",
			content:  "<tool_name>ApplyDiff</tool_name>",
			expected: "ApplyDiff",
		},
		{
			name:     "tool_name with numbers",
			content:  "<tool_name>tool_v2</tool_name>",
			expected: "tool_v2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractToolNameFromPartial(tt.content)
			if result != tt.expected {
				t.Errorf("extractToolNameFromPartial() = %q, want %q", result, tt.expected)
			}
		})
	}
}
