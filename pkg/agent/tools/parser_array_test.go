package tools

import (
	"encoding/xml"
	"testing"
)

// TestParseToolCall_ArrayOfObjects tests parsing arrays of objects
// This tests the apply_diff edits parsing with native XML unmarshaling
func TestParseToolCall_ArrayOfObjects(t *testing.T) {
	// This is what the LLM generates for apply_diff
	xmlContent := `<tool>
<server_name>local</server_name>
<tool_name>apply_diff</tool_name>
<arguments>
  <path>test.go</path>
  <edits>
    <edit>
      <search>old code</search>
      <replace>new code</replace>
    </edit>
  </edits>
</arguments>
</tool>`

	toolCall, _, err := ParseToolCall(xmlContent)
	if err != nil {
		t.Fatalf("ParseToolCall failed: %v", err)
	}

	// Try to unmarshal into the actual struct that apply_diff uses
	var input struct {
		XMLName xml.Name `xml:"arguments"`
		Path    string   `xml:"path"`
		Edits   []struct {
			Search  string `xml:"search"`
			Replace string `xml:"replace"`
		} `xml:"edits>edit"`
	}

	// Now using native XML unmarshaling
	err = xml.Unmarshal(toolCall.GetArgumentsXML(), &input)
	if err != nil {
		t.Logf("Arguments XML: %s", string(toolCall.GetArgumentsXML()))
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Verify the structure
	if input.Path != "test.go" {
		t.Errorf("Expected path 'test.go', got '%s'", input.Path)
	}

	if len(input.Edits) != 1 {
		t.Errorf("Expected 1 edit, got %d", len(input.Edits))
	}

	if len(input.Edits) > 0 {
		if input.Edits[0].Search != "old code" {
			t.Errorf("Expected search 'old code', got '%s'", input.Edits[0].Search)
		}
		if input.Edits[0].Replace != "new code" {
			t.Errorf("Expected replace 'new code', got '%s'", input.Edits[0].Replace)
		}
	}
}

// TestParseToolCall_MultipleEdits tests multiple edits with native XML unmarshaling
func TestParseToolCall_MultipleEdits(t *testing.T) {
	xmlContent := `<tool>
<server_name>local</server_name>
<tool_name>apply_diff</tool_name>
<arguments>
  <path>test.go</path>
  <edits>
    <edit>
      <search>old1</search>
      <replace>new1</replace>
    </edit>
    <edit>
      <search>old2</search>
      <replace>new2</replace>
    </edit>
  </edits>
</arguments>
</tool>`

	toolCall, _, err := ParseToolCall(xmlContent)
	if err != nil {
		t.Fatalf("ParseToolCall failed: %v", err)
	}

	var input struct {
		XMLName xml.Name `xml:"arguments"`
		Path    string   `xml:"path"`
		Edits   []struct {
			Search  string `xml:"search"`
			Replace string `xml:"replace"`
		} `xml:"edits>edit"`
	}

	err = xml.Unmarshal(toolCall.GetArgumentsXML(), &input)
	if err != nil {
		t.Logf("Arguments XML: %s", string(toolCall.GetArgumentsXML()))
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(input.Edits) != 2 {
		t.Errorf("Expected 2 edits, got %d", len(input.Edits))
	}
}
