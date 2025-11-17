package tools

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"math"
	"regexp"
	"strconv"
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
//	<tool_name>write_to_file</tool_name>
//	<arguments>
//	  <path>file.go</path>
//	  <content><![CDATA[package main...]]></content>
//	</arguments>
//	</tool>
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

	// Extract the XML content
	xmlContent := strings.TrimSpace(matches[1])

	// Parse using pure XML parser
	toolCall, err := parseXMLToolCall(xmlContent)
	if err != nil {
		return nil, text, fmt.Errorf("failed to parse tool call XML: %w", err)
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

	return toolCall, remainingText, nil
}

// parseXMLToolCall parses the pure XML format for tool calls
func parseXMLToolCall(xmlContent string) (*ToolCall, error) {
	// Wrap content in a root element for proper XML parsing
	wrappedXML := "<root>" + xmlContent + "</root>"
	
	// Temporary structure for parsing XML
	var parsed struct {
		ServerName string `xml:"server_name"`
		ToolName   string `xml:"tool_name"`
		Arguments  struct {
			InnerXML string `xml:",innerxml"`
		} `xml:"arguments"`
	}

	// Parse the XML structure
	if err := xml.Unmarshal([]byte(wrappedXML), &parsed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML: %w", err)
	}

	// Parse arguments recursively
	args, err := parseXMLArguments(parsed.Arguments.InnerXML)
	if err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	// Convert arguments map to JSON for compatibility
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}

	return &ToolCall{
		ServerName: parsed.ServerName,
		ToolName:   parsed.ToolName,
		Arguments:  json.RawMessage(argsJSON),
	}, nil
}

// parseXMLArguments recursively parses XML arguments into a map
func parseXMLArguments(xmlContent string) (map[string]interface{}, error) {
	xmlContent = strings.TrimSpace(xmlContent)
	if xmlContent == "" {
		return make(map[string]interface{}), nil
	}

	// Wrap in a root element for proper XML parsing
	wrapped := "<root>" + xmlContent + "</root>"

	// Parse into a generic structure
	decoder := xml.NewDecoder(strings.NewReader(wrapped))
	result := make(map[string]interface{})

	var currentKey string
	var currentValue strings.Builder
	var depth int
	var elementStack []string

	for {
		token, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("XML parsing error: %w", err)
		}

		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "root" {
				continue
			}

			elementStack = append(elementStack, t.Name.Local)
			depth++

			if depth == 1 {
				currentKey = t.Name.Local
				currentValue.Reset()
			} else {
				// Nested element - preserve as XML
				currentValue.WriteString("<")
				currentValue.WriteString(t.Name.Local)
				currentValue.WriteString(">")
			}

		case xml.EndElement:
			if t.Name.Local == "root" {
				continue
			}

			if depth == 1 {
				// Top-level element complete
				value := strings.TrimSpace(currentValue.String())
				
				// Check if this key already exists (for arrays)
				if existing, exists := result[currentKey]; exists {
					// Convert to array if not already
					switch v := existing.(type) {
					case []interface{}:
						result[currentKey] = append(v, parseValue(value))
					default:
						result[currentKey] = []interface{}{v, parseValue(value)}
					}
				} else {
					result[currentKey] = parseValue(value)
				}
				currentKey = ""
			} else {
				// Nested element end
				currentValue.WriteString("</")
				currentValue.WriteString(t.Name.Local)
				currentValue.WriteString(">")
			}

			depth--
			if len(elementStack) > 0 {
				elementStack = elementStack[:len(elementStack)-1]
			}

		case xml.CharData:
			data := string(t)
			if depth > 0 {
				currentValue.WriteString(data)
			}
		}
	}

	return result, nil
}

// parseValue attempts to parse a string value into the appropriate type
func parseValue(value string) interface{} {
	value = strings.TrimSpace(value)

	// Empty value
	if value == "" {
		return ""
	}

	// Check if it contains nested XML
	if strings.Contains(value, "<") && strings.Contains(value, ">") {
		// Try to parse as nested structure
		if nested, err := parseXMLArguments(value); err == nil && len(nested) > 0 {
			return nested
		}
	}

	// Boolean detection (case-insensitive)
	lower := strings.ToLower(value)
	if lower == "true" {
		return true
	}
	if lower == "false" {
		return false
	}

	// Null detection
	if lower == "null" {
		return nil
	}

	// Integer detection
	if i, err := strconv.ParseInt(value, 10, 64); err == nil {
		// Return int if it fits in int range
		if i >= math.MinInt && i <= math.MaxInt {
			return int(i)
		}
		return i
	}

	// Float detection
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}

	// Default: string
	return value
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
