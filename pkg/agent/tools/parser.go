package tools

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

const defaultServerName = "local"

// ParseToolCall extracts a tool call from an LLM response that contains
// XML-formatted tool invocations.
//
// Expected format:
//
//	<tool>{"server_name": "local", "tool_name": "task_completion", "arguments": {"result": "Done!"}}</tool>
//
// Returns the parsed ToolCall and the remaining text after removing the tool call,
// or an error if parsing fails.
func ParseToolCall(text string) (*ToolCall, string, error) {
	// Regex to match <tool>...</tool> tags
	// Uses non-greedy matching with (?s) flag to match across newlines
	toolRegex := regexp.MustCompile(`(?s)<tool>(.*?)</tool>`)

	matches := toolRegex.FindStringSubmatch(text)
	if len(matches) < 2 {
		return nil, text, fmt.Errorf("no tool call found in text")
	}

	// Extract the JSON content
	jsonContent := strings.TrimSpace(matches[1])

	// Parse the JSON into a ToolCall
	var toolCall ToolCall
	if err := json.Unmarshal([]byte(jsonContent), &toolCall); err != nil {
		return nil, text, fmt.Errorf("failed to parse tool call JSON: %w", err)
	}

	// Validate required fields
	if toolCall.ToolName == "" {
		return nil, text, fmt.Errorf("tool_name is required in tool call")
	}

	// Server name defaults to "local" if not specified
	if toolCall.ServerName == "" {
		toolCall.ServerName = defaultServerName
	}

	// Remove the tool call from the text
	remainingText := toolRegex.ReplaceAllString(text, "")
	remainingText = strings.TrimSpace(remainingText)

	return &toolCall, remainingText, nil
}

// ExtractThinkingAndToolCall separates thinking content from a tool call.
// If a tool call is found, it returns the thinking text (before the tool call),
// the tool call itself, and any remaining text after the tool call.
// If no tool call is found, it returns the entire text as thinking with nil tool call.
func ExtractThinkingAndToolCall(text string) (thinking string, toolCall *ToolCall, remaining string, err error) {
	// Find where the tool call starts in the original text
	// (?s) flag makes . match newlines
	toolRegex := regexp.MustCompile(`(?s)<tool>.*?</tool>`)
	toolLocation := toolRegex.FindStringIndex(text)

	// If no tool call found, entire text is thinking
	if toolLocation == nil {
		return text, nil, "", nil
	}

	// Text before tool call is thinking
	thinking = strings.TrimSpace(text[:toolLocation[0]])

	// Text after tool call is remaining
	remaining = strings.TrimSpace(text[toolLocation[1]:])

	// Extract just the tool call XML
	toolCallText := text[toolLocation[0]:toolLocation[1]]

	// Parse the tool call
	toolCall, _, parseErr := ParseToolCall(toolCallText)
	if parseErr != nil {
		// Shouldn't happen if regex matched, but handle it
		return text, nil, "", parseErr
	}

	return thinking, toolCall, remaining, nil
}

// HasToolCall checks if the text contains a tool call without parsing it
func HasToolCall(text string) bool {
	// (?s) flag makes . match newlines
	toolRegex := regexp.MustCompile(`(?s)<tool>.*?</tool>`)
	return toolRegex.MatchString(text)
}

// ValidateToolCall checks if a tool call has all required fields
func ValidateToolCall(tc *ToolCall) error {
	if tc == nil {
		return fmt.Errorf("tool call is nil")
	}
	if tc.ToolName == "" {
		return fmt.Errorf("tool_name is required")
	}
	if tc.ServerName == "" {
		return fmt.Errorf("server_name is required")
	}
	return nil
}
