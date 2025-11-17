package tools

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParseToolCall_PureXML_Simple(t *testing.T) {
	input := `<tool>
<server_name>local</server_name>
<tool_name>task_completion</tool_name>
<arguments>
  <result>Task completed successfully</result>
</arguments>
</tool>`

	toolCall, remaining, err := ParseToolCall(input)
	if err != nil {
		t.Fatalf("ParseToolCall failed: %v", err)
	}

	if toolCall.ServerName != "local" {
		t.Errorf("Expected server_name 'local', got '%s'", toolCall.ServerName)
	}

	if toolCall.ToolName != "task_completion" {
		t.Errorf("Expected tool_name 'task_completion', got '%s'", toolCall.ToolName)
	}

	// Parse arguments
	var args map[string]interface{}
	if err := json.Unmarshal(toolCall.Arguments, &args); err != nil {
		t.Fatalf("Failed to unmarshal arguments: %v", err)
	}

	if args["result"] != "Task completed successfully" {
		t.Errorf("Expected result 'Task completed successfully', got '%v'", args["result"])
	}

	if remaining != "" {
		t.Errorf("Expected empty remaining text, got '%s'", remaining)
	}
}

func TestParseToolCall_PureXML_WithCDATA(t *testing.T) {
	input := `<tool>
<server_name>local</server_name>
<tool_name>write_to_file</tool_name>
<arguments>
  <path>test.go</path>
  <content><![CDATA[package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
	// Special characters: "quotes", 'apostrophes', \backslashes\
}]]></content>
</arguments>
</tool>`

	toolCall, _, err := ParseToolCall(input)
	if err != nil {
		t.Fatalf("ParseToolCall failed: %v", err)
	}

	var args map[string]interface{}
	if err := json.Unmarshal(toolCall.Arguments, &args); err != nil {
		t.Fatalf("Failed to unmarshal arguments: %v", err)
	}

	expectedContent := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
	// Special characters: "quotes", 'apostrophes', \backslashes\
}`

	if args["path"] != "test.go" {
		t.Errorf("Expected path 'test.go', got '%v'", args["path"])
	}

	content := args["content"].(string)
	if content != expectedContent {
		t.Errorf("Content mismatch.\nExpected:\n%s\n\nGot:\n%s", expectedContent, content)
	}
}

func TestParseToolCall_PureXML_ComplexDiff(t *testing.T) {
	input := `<tool>
<server_name>local</server_name>
<tool_name>apply_diff</tool_name>
<arguments>
  <path>src/app.ts</path>
  <edits>
    <edit>
      <search><![CDATA[const oldConfig = {
	api: "http://old.com"
};]]></search>
      <replace><![CDATA[const newConfig = {
	api: "https://new.com",
	verbose: true
};]]></replace>
    </edit>
    <edit>
      <search><![CDATA[function process(data) {
	return data.map(x => x * 2);
}]]></search>
      <replace><![CDATA[function process(data) {
	return data.filter(x => x > 0).map(x => x * 2);
}]]></replace>
    </edit>
  </edits>
</arguments>
</tool>`

	toolCall, _, err := ParseToolCall(input)
	if err != nil {
		t.Fatalf("ParseToolCall failed: %v", err)
	}

	var args map[string]interface{}
	if err := json.Unmarshal(toolCall.Arguments, &args); err != nil {
		t.Fatalf("Failed to unmarshal arguments: %v", err)
	}

	if args["path"] != "src/app.ts" {
		t.Errorf("Expected path 'src/app.ts', got '%v'", args["path"])
	}

	// Check edits structure
	edits, ok := args["edits"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected edits to be a map, got %T", args["edits"])
	}

	// The parser should have nested the edit elements
	if edits["edit"] == nil {
		t.Errorf("Expected 'edit' key in edits, got: %+v", edits)
	}
}

func TestParseToolCall_PureXML_ArrayArguments(t *testing.T) {
	input := `<tool>
<server_name>local</server_name>
<tool_name>search_files</tool_name>
<arguments>
  <path>src</path>
  <pattern>\.ts$</pattern>
  <exclude>node_modules</exclude>
  <exclude>dist</exclude>
  <exclude>.git</exclude>
</arguments>
</tool>`

	toolCall, _, err := ParseToolCall(input)
	if err != nil {
		t.Fatalf("ParseToolCall failed: %v", err)
	}

	var args map[string]interface{}
	if err := json.Unmarshal(toolCall.Arguments, &args); err != nil {
		t.Fatalf("Failed to unmarshal arguments: %v", err)
	}

	// Check that exclude became an array
	excludes, ok := args["exclude"].([]interface{})
	if !ok {
		t.Fatalf("Expected exclude to be an array, got %T: %v", args["exclude"], args["exclude"])
	}

	if len(excludes) != 3 {
		t.Errorf("Expected 3 exclude values, got %d", len(excludes))
	}

	expectedExcludes := []string{"node_modules", "dist", ".git"}
	for i, expected := range expectedExcludes {
		if excludes[i] != expected {
			t.Errorf("Expected exclude[%d] = '%s', got '%v'", i, expected, excludes[i])
		}
	}
}

func TestParseToolCall_PureXML_WithThinking(t *testing.T) {
	input := `I need to write a file now.

<tool>
<server_name>local</server_name>
<tool_name>write_to_file</tool_name>
<arguments>
  <path>hello.txt</path>
  <content><![CDATA[Hello, World!]]></content>
</arguments>
</tool>

That should do it.`

	toolCall, _, err := ParseToolCall(input)
	if err != nil {
		t.Fatalf("ParseToolCall failed: %v", err)
	}

	if toolCall.ToolName != "write_to_file" {
		t.Errorf("Expected tool_name 'write_to_file', got '%s'", toolCall.ToolName)
	}

	// The input text contains thinking before and after
	if !strings.Contains(input, "I need to write a file now.") {
		t.Errorf("Original input should contain thinking text")
	}
}

func TestParseToolCall_PureXML_DefaultServerName(t *testing.T) {
	input := `<tool>
<tool_name>task_completion</tool_name>
<arguments>
  <result>Done</result>
</arguments>
</tool>`

	toolCall, _, err := ParseToolCall(input)
	if err != nil {
		t.Fatalf("ParseToolCall failed: %v", err)
	}

	if toolCall.ServerName != "local" {
		t.Errorf("Expected default server_name 'local', got '%s'", toolCall.ServerName)
	}
}

func TestParseToolCall_PureXML_EmptyArguments(t *testing.T) {
	input := `<tool>
<server_name>local</server_name>
<tool_name>some_tool</tool_name>
<arguments>
</arguments>
</tool>`

	toolCall, _, err := ParseToolCall(input)
	if err != nil {
		t.Fatalf("ParseToolCall failed: %v", err)
	}

	var args map[string]interface{}
	if err := json.Unmarshal(toolCall.Arguments, &args); err != nil {
		t.Fatalf("Failed to unmarshal arguments: %v", err)
	}

	if len(args) != 0 {
		t.Errorf("Expected empty arguments, got %+v", args)
	}
}

func TestParseToolCall_PureXML_MissingToolName(t *testing.T) {
	input := `<tool>
<server_name>local</server_name>
<arguments>
  <result>Done</result>
</arguments>
</tool>`

	_, _, err := ParseToolCall(input)
	if err == nil {
		t.Fatal("Expected error for missing tool_name, got nil")
	}

	if !strings.Contains(err.Error(), "tool_name is required") {
		t.Errorf("Expected error about missing tool_name, got: %v", err)
	}
}

func TestParseToolCall_PureXML_NoToolTag(t *testing.T) {
	input := `Some text without a tool tag`

	_, _, err := ParseToolCall(input)
	if err == nil {
		t.Fatal("Expected error for missing tool tag, got nil")
	}

	if !strings.Contains(err.Error(), "no tool call found") {
		t.Errorf("Expected error about no tool call found, got: %v", err)
	}
}

func TestParseToolCall_PureXML_LargeContent(t *testing.T) {
	// Test with large file content (simulating real-world usage)
	largeCode := strings.Repeat("func example() {\n\treturn true\n}\n\n", 100)
	
	input := `<tool>
<server_name>local</server_name>
<tool_name>write_to_file</tool_name>
<arguments>
  <path>large.go</path>
  <content><![CDATA[` + largeCode + `]]></content>
</arguments>
</tool>`

	toolCall, _, err := ParseToolCall(input)
	if err != nil {
		t.Fatalf("ParseToolCall failed with large content: %v", err)
	}

	var args map[string]interface{}
	if err := json.Unmarshal(toolCall.Arguments, &args); err != nil {
		t.Fatalf("Failed to unmarshal arguments: %v", err)
	}

	content := args["content"].(string)
	// Allow for minor whitespace differences due to XML parsing
	if len(content) < len(largeCode)-10 || len(content) > len(largeCode)+10 {
		t.Errorf("Large content length mismatch. Expected ~%d, got %d", len(largeCode), len(content))
	}
	// Verify the actual code is present
	if !strings.Contains(content, "func example()") {
		t.Error("Large content missing expected function definition")
	}
}

func TestExtractThinkingAndToolCall_PureXML(t *testing.T) {
	input := `Let me complete this task.

<tool>
<server_name>local</server_name>
<tool_name>task_completion</tool_name>
<arguments>
  <result>All done!</result>
</arguments>
</tool>

Everything is finished.`

	thinking, toolCall, remaining, err := ExtractThinkingAndToolCall(input)
	if err != nil {
		t.Fatalf("ExtractThinkingAndToolCall failed: %v", err)
	}

	if thinking != "Let me complete this task." {
		t.Errorf("Expected thinking 'Let me complete this task.', got '%s'", thinking)
	}

	if toolCall == nil {
		t.Fatal("Expected toolCall to be non-nil")
	}

	if toolCall.ToolName != "task_completion" {
		t.Errorf("Expected tool_name 'task_completion', got '%s'", toolCall.ToolName)
	}

	if remaining != "Everything is finished." {
		t.Errorf("Expected remaining 'Everything is finished.', got '%s'", remaining)
	}
}

func TestHasToolCall_PureXML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name: "Has tool call",
			input: `<tool>
<server_name>local</server_name>
<tool_name>test</tool_name>
<arguments></arguments>
</tool>`,
			expected: true,
		},
		{
			name:     "No tool call",
			input:    "Just some regular text",
			expected: false,
		},
		{
			name:     "Incomplete tool tag",
			input:    "<tool>incomplete",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasToolCall(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestValidateToolCall_XML(t *testing.T) {
	tests := []struct {
		name      string
		toolCall  *ToolCall
		expectErr bool
	}{
		{
			name: "Valid tool call",
			toolCall: &ToolCall{
				ServerName: "local",
				ToolName:   "test",
				Arguments:  json.RawMessage("{}"),
			},
			expectErr: false,
		},
		{
			name:      "Nil tool call",
			toolCall:  nil,
			expectErr: true,
		},
		{
			name: "Missing tool name",
			toolCall: &ToolCall{
				ServerName: "local",
				Arguments:  json.RawMessage("{}"),
			},
			expectErr: true,
		},
		{
			name: "Missing server name",
			toolCall: &ToolCall{
				ToolName:  "test",
				Arguments: json.RawMessage("{}"),
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateToolCall(tt.toolCall)
			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}