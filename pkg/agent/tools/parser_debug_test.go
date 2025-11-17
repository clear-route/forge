package tools

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func TestDebugXMLParsing(t *testing.T) {
	xmlContent := `
<server_name>local</server_name>
<tool_name>task_completion</tool_name>
<arguments>
  <result>Task completed successfully</result>
</arguments>`

	// Wrap in root element (same as parser does)
	wrappedXML := "<root>" + xmlContent + "</root>"

	// Test basic XML unmarshaling
	var parsed struct {
		ServerName string `xml:"server_name"`
		ToolName   string `xml:"tool_name"`
		Arguments  struct {
			InnerXML string `xml:",innerxml"`
		} `xml:"arguments"`
	}

	if err := xml.Unmarshal([]byte(wrappedXML), &parsed); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	fmt.Printf("ServerName: '%s'\n", parsed.ServerName)
	fmt.Printf("ToolName: '%s'\n", parsed.ToolName)
	fmt.Printf("Arguments InnerXML: '%s'\n", parsed.Arguments.InnerXML)

	if parsed.ServerName == "" {
		t.Error("ServerName is empty")
	}
	if parsed.ToolName == "" {
		t.Error("ToolName is empty")
	}
}
