package parser

import (
	"testing"
)

// TestToolCallParser_XMLContent verifies that the streaming parser correctly handles
// XML+CDATA content within tool tags. The parser is format-agnostic and only cares
// about <tool> boundaries, making it compatible with the new XML+CDATA format.
func TestToolCallParser_XMLContent(t *testing.T) {
	parser := NewToolCallParser()

	// Complete XML tool call with CDATA in single chunk
	tc, rc := parser.Parse(`<tool>
<server_name>local</server_name>
<tool_name>write_to_file</tool_name>
<arguments>
  <path>example.go</path>
  <content><![CDATA[package main

func main() {
	fmt.Println("Hello, World!")
}]]></content>
</arguments>
</tool>`)

	if tc == nil {
		t.Fatal("Expected tool call, got nil")
	}

	expectedContent := `<server_name>local</server_name>
<tool_name>write_to_file</tool_name>
<arguments>
  <path>example.go</path>
  <content><![CDATA[package main

func main() {
	fmt.Println("Hello, World!")
}]]></content>
</arguments>`

	if tc.Content != expectedContent {
		t.Errorf("Expected XML content, got '%s'", tc.Content)
	}
	if rc != nil {
		t.Errorf("Expected no regular content, got %v", rc)
	}
}

// TestToolCallParser_StreamedXMLWithCDATA verifies that XML+CDATA content
// can be streamed across multiple chunks and the parser correctly accumulates it.
func TestToolCallParser_StreamedXMLWithCDATA(t *testing.T) {
	parser := NewToolCallParser()

	// Simulate streaming XML chunks
	chunks := []string{
		"<tool>\n",
		"<server_name>local</server_name>\n",
		"<tool_name>apply_diff</tool_name>\n",
		"<arguments>\n",
		"  <path>test.go</path>\n",
		"  <diff><![CDATA[\n",
		"<<<<<<< SEARCH\n",
		"old code\n",
		"=======\n",
		"new code\n",
		">>>>>>> REPLACE\n",
		"]]></diff>\n",
		"</arguments>\n",
		"</tool>",
	}

	var toolCallContent string

	for _, chunk := range chunks {
		toolCall, _ := parser.Parse(chunk)
		if toolCall != nil {
			toolCallContent = toolCall.Content
		}
	}

	// Verify the complete XML content was captured
	if !contains(toolCallContent, "<server_name>local</server_name>") {
		t.Error("Expected server_name in content")
	}
	if !contains(toolCallContent, "<tool_name>apply_diff</tool_name>") {
		t.Error("Expected tool_name in content")
	}
	if !contains(toolCallContent, "<![CDATA[") {
		t.Error("Expected CDATA section in content")
	}
	if !contains(toolCallContent, "<<<<<<< SEARCH") {
		t.Error("Expected SEARCH marker in CDATA content")
	}
	if !contains(toolCallContent, ">>>>>>> REPLACE") {
		t.Error("Expected REPLACE marker in CDATA content")
	}
}

// TestToolCallParser_XMLWithNestedTags verifies handling of nested XML elements
func TestToolCallParser_XMLWithNestedTags(t *testing.T) {
	parser := NewToolCallParser()

	tc, _ := parser.Parse(`<tool>
<server_name>mcp-server</server_name>
<tool_name>use_mcp_tool</tool_name>
<arguments>
  <files>
    <file>
      <path>file1.txt</path>
      <content><![CDATA[content 1]]></content>
    </file>
    <file>
      <path>file2.txt</path>
      <content><![CDATA[content 2]]></content>
    </file>
  </files>
</arguments>
</tool>`)

	if tc == nil {
		t.Fatal("Expected tool call, got nil")
	}

	// Verify nested structure is preserved
	if !contains(tc.Content, "<files>") {
		t.Error("Expected <files> tag in content")
	}
	if !contains(tc.Content, "<file>") {
		t.Error("Expected <file> tags in content")
	}
	if !contains(tc.Content, "content 1") && !contains(tc.Content, "content 2") {
		t.Error("Expected CDATA content preserved")
	}
}

// TestToolCallParser_XMLWithSpecialCharacters verifies CDATA preserves special chars
func TestToolCallParser_XMLWithSpecialCharacters(t *testing.T) {
	parser := NewToolCallParser()

	// Content with characters that would break JSON parsing
	tc, _ := parser.Parse(`<tool>
<tool_name>write_to_file</tool_name>
<arguments>
  <content><![CDATA[{
	"key": "value with \"quotes\"",
	"escaped": "line1\nline2\ttab",
	"special": "<>&"
}]]></content>
</arguments>
</tool>`)

	if tc == nil {
		t.Fatal("Expected tool call, got nil")
	}

	// Verify special characters are preserved
	if !contains(tc.Content, `"key": "value with \"quotes\""`) {
		t.Error("Expected quotes preserved in CDATA")
	}
	if !contains(tc.Content, `<>&`) {
		t.Error("Expected special XML chars preserved in CDATA")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}
