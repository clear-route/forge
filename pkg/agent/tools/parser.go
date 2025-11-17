package tools

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
)

const defaultServerName = "local"

// ParseToolCall extracts a tool call from an LLM response that contains
// XML-formatted tool invocations.
//
// Expected format (Pure XML with CDATA):
//
//	<tool>
//	<server_name>local</server_name>
//	<tool_name>apply_diff</tool_name>
//	<arguments>
//	  <path>file.go</path>
//	  <edits>
//	    <edit>
//	      <search><![CDATA[old code]]></search>
//	      <replace><![CDATA[new code]]></replace>
//	    </edit>
//	  </edits>
//	</arguments>
//	</tool>
//
// Returns the parsed ToolCall and the remaining text after removing the tool call,
// or an error if parsing fails.
func ParseToolCall(text string) (*ToolCall, string, error) {
	// Regex to match <tool>...</tool> tags
	toolRegex := regexp.MustCompile(`(?s)<tool>.*?</tool>`)

	matches := toolRegex.FindStringSubmatch(text)
	if len(matches) < 1 {
		return nil, text, fmt.Errorf("no tool call found in text")
	}

	// Extract the full <tool> element including tags
	toolXML := strings.TrimSpace(matches[0])

	var toolCall ToolCall
	if err := xml.Unmarshal([]byte(toolXML), &toolCall); err != nil {
		return nil, text, fmt.Errorf("failed to unmarshal tool call XML: %w", err)
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
	toolRegex := regexp.MustCompile(`(?s)<tool>.*?</tool>`)
	
	if !toolRegex.MatchString(text) {
		return text, nil, "", nil
	}

	loc := toolRegex.FindStringIndex(text)
	if loc == nil {
		return text, nil, "", nil
	}

	thinking = strings.TrimSpace(text[:loc[0]])
	toolCallText := text[loc[0]:loc[1]]
	remaining = strings.TrimSpace(text[loc[1]:])

	toolCall, _, err = ParseToolCall(toolCallText)
	if err != nil {
		return thinking, nil, remaining, err
	}

	return thinking, toolCall, remaining, nil
}

// HasToolCall checks if the text contains a tool call.
func HasToolCall(text string) bool {
	toolRegex := regexp.MustCompile(`(?s)<tool>.*?</tool>`)
	return toolRegex.MatchString(text)
}

// ValidateToolCall checks if a ToolCall has all required fields.
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