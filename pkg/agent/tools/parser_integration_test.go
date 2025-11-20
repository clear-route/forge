package tools

import (
	"encoding/xml"
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
		XMLName   xml.Name `xml:"arguments"`
		Path      string   `xml:"path"`
		Recursive bool     `xml:"recursive"`
	}

	if err := xml.Unmarshal(toolCall.GetArgumentsXML(), &args); err != nil {
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
		XMLName   xml.Name `xml:"arguments"`
		Path      string   `xml:"path"`
		LineStart int      `xml:"line_start"`
		LineEnd   int      `xml:"line_end"`
	}

	if err := xml.Unmarshal(toolCall.GetArgumentsXML(), &args); err != nil {
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
		XMLName  xml.Name `xml:"arguments"`
		Name     string   `xml:"name"`
		Count    int      `xml:"count"`
		Ratio    float64  `xml:"ratio"`
		Enabled  bool     `xml:"enabled"`
		Disabled bool     `xml:"disabled"`
		Optional string   `xml:"optional"`
	}

	if err := xml.Unmarshal(toolCall.GetArgumentsXML(), &args); err != nil {
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
	if args.Optional != "null" {
		t.Errorf("Expected optional='null', got %q", args.Optional)
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
		XMLName xml.Name `xml:"arguments"`
		Path    string   `xml:"path"`
		Content string   `xml:"content"`
	}

	if err := xml.Unmarshal(toolCall.GetArgumentsXML(), &args); err != nil {
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
		XMLName xml.Name `xml:"arguments"`
		Path    string   `xml:"path"`
		Pattern string   `xml:"pattern"`
		Exclude []string `xml:"exclude"`
	}

	if err := xml.Unmarshal(toolCall.GetArgumentsXML(), &args); err != nil {
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

// TestXMLParserWithEntityEscaping tests the parser with XML entity-escaped content
// This verifies ADR-0024: XML entity escaping as primary method
func TestXMLParserWithEntityEscaping(t *testing.T) {
	xmlInput := `<tool>
<server_name>local</server_name>
<tool_name>write_to_file</tool_name>
<arguments>
  <path>test.go</path>
  <content>func example() *Config {
	return &amp;Config{name: &quot;test&quot;}
}
if x &lt; 10 &amp;&amp; y &gt; 5 {
	process()
}</content>
</arguments>
</tool>`

	toolCall, _, err := ParseToolCall(xmlInput)
	if err != nil {
		t.Fatalf("ParseToolCall failed: %v", err)
	}

	var args struct {
		XMLName xml.Name `xml:"arguments"`
		Path    string   `xml:"path"`
		Content string   `xml:"content"`
	}

	if err := xml.Unmarshal(toolCall.GetArgumentsXML(), &args); err != nil {
		t.Fatalf("Failed to unmarshal arguments: %v", err)
	}

	if args.Path != "test.go" {
		t.Errorf("Expected path='test.go', got %q", args.Path)
	}

	// Verify escaped entities are correctly decoded
	expectedContent := `func example() *Config {
	return &Config{name: "test"}
}
if x < 10 && y > 5 {
	process()
}`

	if args.Content != expectedContent {
		t.Errorf("Content mismatch.\nExpected:\n%s\n\nGot:\n%s", expectedContent, args.Content)
	}
}

// TestXMLParserWithMixedEscapingAndCDATA tests using both escaping and CDATA in same tool call
// This verifies both methods work together (ADR-0024)
func TestXMLParserWithMixedEscapingAndCDATA(t *testing.T) {
	xmlInput := `<tool>
<server_name>local</server_name>
<tool_name>apply_diff</tool_name>
<arguments>
  <path>test.go</path>
  <search>old &amp; code</search>
  <replace><![CDATA[new & improved code]]></replace>
</arguments>
</tool>`

	toolCall, _, err := ParseToolCall(xmlInput)
	if err != nil {
		t.Fatalf("ParseToolCall failed: %v", err)
	}

	var args struct {
		XMLName xml.Name `xml:"arguments"`
		Path    string   `xml:"path"`
		Search  string   `xml:"search"`
		Replace string   `xml:"replace"`
	}

	if err := xml.Unmarshal(toolCall.GetArgumentsXML(), &args); err != nil {
		t.Fatalf("Failed to unmarshal arguments: %v", err)
	}

	// Verify entity-escaped field
	if args.Search != "old & code" {
		t.Errorf("Expected search='old & code', got %q", args.Search)
	}

	// Verify CDATA field
	if args.Replace != "new & improved code" {
		t.Errorf("Expected replace='new & improved code', got %q", args.Replace)
	}
}

// TestXMLParserWithEntityEscaping tests the parser with XML entity-escaped content
// This verifies ADR-0024: XML entity escaping as primary method
func TestXMLParserWithAllStandardEntities(t *testing.T) {
	xmlInput := `<tool>
<server_name>local</server_name>
<tool_name>test_tool</tool_name>
<arguments>
  <content>Test &amp; &lt; &gt; &quot; &apos; entities</content>
</arguments>
</tool>`

	toolCall, _, err := ParseToolCall(xmlInput)
	if err != nil {
		t.Fatalf("ParseToolCall failed: %v", err)
	}

	var args struct {
		XMLName xml.Name `xml:"arguments"`
		Content string   `xml:"content"`
	}

	if err := xml.Unmarshal(toolCall.GetArgumentsXML(), &args); err != nil {
		t.Fatalf("Failed to unmarshal arguments: %v", err)
	}

	expected := `Test & < > " ' entities`
	if args.Content != expected {
		t.Errorf("Expected content=%q, got %q", expected, args.Content)
	}
}
