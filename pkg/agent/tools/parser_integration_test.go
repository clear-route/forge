package tools

import (
	"encoding/json"
	"testing"
)

func TestXMLParserWithBooleanArgument(t *testing.T) {
	xmlInput := `
<tool>
<server_name>local</server_name>
<tool_name>list_files</tool_name>
<arguments>
  <path>.</path>
  <recursive>true</recursive>
</arguments>
</tool>
`

	toolCall, _, err := ParseToolCall(xmlInput)
	if err != nil {
		t.Fatalf("ParseToolCall failed: %v", err)
	}

	// Verify tool call metadata
	if toolCall.ServerName != "local" {
		t.Errorf("Expected server_name='local', got %q", toolCall.ServerName)
	}
	if toolCall.ToolName != "list_files" {
		t.Errorf("Expected tool_name='list_files', got %q", toolCall.ToolName)
	}

	// Unmarshal arguments into struct
	var args struct {
		Path      string `json:"path"`
		Recursive bool   `json:"recursive"`
	}

	if err := json.Unmarshal(toolCall.Arguments, &args); err != nil {
		t.Fatalf("Failed to unmarshal arguments: %v", err)
	}

	if args.Path != "." {
		t.Errorf("Expected path='.', got %q", args.Path)
	}

	if !args.Recursive {
		t.Errorf("Expected recursive=true, got %v", args.Recursive)
	}
}

func TestXMLParserWithNumericArguments(t *testing.T) {
	xmlInput := `
<tool>
<server_name>local</server_name>
<tool_name>read_file</tool_name>
<arguments>
  <path>test.go</path>
  <line_start>1</line_start>
  <line_end>100</line_end>
</arguments>
</tool>
`

	toolCall, _, err := ParseToolCall(xmlInput)
	if err != nil {
		t.Fatalf("ParseToolCall failed: %v", err)
	}

	var args struct {
		Path      string `json:"path"`
		LineStart int    `json:"line_start"`
		LineEnd   int    `json:"line_end"`
	}

	if err := json.Unmarshal(toolCall.Arguments, &args); err != nil {
		t.Fatalf("Failed to unmarshal arguments: %v", err)
	}

	if args.Path != "test.go" {
		t.Errorf("Expected path='test.go', got %q", args.Path)
	}
	if args.LineStart != 1 {
		t.Errorf("Expected line_start=1, got %d", args.LineStart)
	}
	if args.LineEnd != 100 {
		t.Errorf("Expected line_end=100, got %d", args.LineEnd)
	}
}

func TestXMLParserWithMixedTypes(t *testing.T) {
	xmlInput := `
<tool>
<server_name>local</server_name>
<tool_name>test_tool</tool_name>
<arguments>
  <name>test</name>
  <count>42</count>
  <ratio>3.14</ratio>
  <enabled>true</enabled>
  <disabled>false</disabled>
  <optional>null</optional>
</arguments>
</tool>
`

	toolCall, _, err := ParseToolCall(xmlInput)
	if err != nil {
		t.Fatalf("ParseToolCall failed: %v", err)
	}

	var args struct {
		Name     string   `json:"name"`
		Count    int      `json:"count"`
		Ratio    float64  `json:"ratio"`
		Enabled  bool     `json:"enabled"`
		Disabled bool     `json:"disabled"`
		Optional *string  `json:"optional"` // pointer to handle null
	}

	if err := json.Unmarshal(toolCall.Arguments, &args); err != nil {
		t.Fatalf("Failed to unmarshal arguments: %v", err)
	}

	if args.Name != "test" {
		t.Errorf("Expected name='test', got %q", args.Name)
	}
	if args.Count != 42 {
		t.Errorf("Expected count=42, got %d", args.Count)
	}
	if args.Ratio != 3.14 {
		t.Errorf("Expected ratio=3.14, got %f", args.Ratio)
	}
	if !args.Enabled {
		t.Errorf("Expected enabled=true, got %v", args.Enabled)
	}
	if args.Disabled {
		t.Errorf("Expected disabled=false, got %v", args.Disabled)
	}
	if args.Optional != nil {
		t.Errorf("Expected optional=nil, got %v", args.Optional)
	}
}

func TestXMLParserWithLargeContent(t *testing.T) {
	// Test with CDATA content (should remain as string)
	xmlInput := `
<tool>
<server_name>local</server_name>
<tool_name>write_to_file</tool_name>
<arguments>
  <path>test.go</path>
  <content><![CDATA[package main

import "fmt"

func main() {
	// This is a test
	fmt.Println("Hello, world!")
	value := true // boolean in code, not argument
	count := 123  // number in code, not argument
}]]></content>
</arguments>
</tool>
`

	toolCall, _, err := ParseToolCall(xmlInput)
	if err != nil {
		t.Fatalf("ParseToolCall failed: %v", err)
	}

	var args struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}

	if err := json.Unmarshal(toolCall.Arguments, &args); err != nil {
		t.Fatalf("Failed to unmarshal arguments: %v", err)
	}

	if args.Path != "test.go" {
		t.Errorf("Expected path='test.go', got %q", args.Path)
	}

	// Verify content is preserved as string
	expectedContent := `package main

import "fmt"

func main() {
	// This is a test
	fmt.Println("Hello, world!")
	value := true // boolean in code, not argument
	count := 123  // number in code, not argument
}`

	if args.Content != expectedContent {
		t.Errorf("Content mismatch.\nExpected:\n%s\n\nGot:\n%s", expectedContent, args.Content)
	}
}

func TestXMLParserWithArrayArguments(t *testing.T) {
	xmlInput := `
<tool>
<server_name>local</server_name>
<tool_name>search_files</tool_name>
<arguments>
  <path>src</path>
  <pattern>\.go$</pattern>
  <exclude>vendor</exclude>
  <exclude>node_modules</exclude>
  <exclude>.git</exclude>
</arguments>
</tool>
`

	toolCall, _, err := ParseToolCall(xmlInput)
	if err != nil {
		t.Fatalf("ParseToolCall failed: %v", err)
	}

	var args struct {
		Path    string   `json:"path"`
		Pattern string   `json:"pattern"`
		Exclude []string `json:"exclude"`
	}

	if err := json.Unmarshal(toolCall.Arguments, &args); err != nil {
		t.Fatalf("Failed to unmarshal arguments: %v", err)
	}

	if args.Path != "src" {
		t.Errorf("Expected path='src', got %q", args.Path)
	}
	if args.Pattern != `\.go$` {
		t.Errorf("Expected pattern='\\.go$', got %q", args.Pattern)
	}

	expectedExclude := []string{"vendor", "node_modules", ".git"}
	if len(args.Exclude) != len(expectedExclude) {
		t.Errorf("Expected %d exclude entries, got %d", len(expectedExclude), len(args.Exclude))
	}
	for i, expected := range expectedExclude {
		if i >= len(args.Exclude) || args.Exclude[i] != expected {
			t.Errorf("Expected exclude[%d]=%q, got %q", i, expected, args.Exclude[i])
		}
	}
}