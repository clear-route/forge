package tools

import (
	"testing"
)

// BenchmarkParseToolCall measures the performance of parsing a typical tool call
func BenchmarkParseToolCall(b *testing.B) {
	xmlInput := `<tool>
<server_name>local</server_name>
<tool_name>write_to_file</tool_name>
<arguments>
  <path>test.go</path>
  <content><![CDATA[package main

import "fmt"

func main() {
	fmt.Println("Hello, world!")
	value := &Config{name: "test"}
	count := 123
}]]></content>
</arguments>
</tool>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := ParseToolCall(xmlInput)
		if err != nil {
			b.Fatalf("ParseToolCall failed: %v", err)
		}
	}
}

// BenchmarkGetArgumentsXML measures the performance of GetArgumentsXML
func BenchmarkGetArgumentsXML(b *testing.B) {
	xmlInput := `<tool>
<server_name>local</server_name>
<tool_name>apply_diff</tool_name>
<arguments>
  <path>file.go</path>
  <edits>
    <edit>
      <search><![CDATA[old code here with lots of content]]></search>
      <replace><![CDATA[new code here with lots of content]]></replace>
    </edit>
  </edits>
</arguments>
</tool>`

	toolCall, _, err := ParseToolCall(xmlInput)
	if err != nil {
		b.Fatalf("ParseToolCall failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = toolCall.GetArgumentsXML()
	}
}

// BenchmarkParseLargeContent tests parsing performance with large content
func BenchmarkParseLargeContent(b *testing.B) {
	// Generate large content (1MB)
	largeContent := make([]byte, 1024*1024)
	for i := range largeContent {
		largeContent[i] = byte('a' + (i % 26))
	}

	xmlInput := `<tool>
<server_name>local</server_name>
<tool_name>write_to_file</tool_name>
<arguments>
  <path>large.txt</path>
  <content><![CDATA[` + string(largeContent) + `]]></content>
</arguments>
</tool>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := ParseToolCall(xmlInput)
		if err != nil {
			b.Fatalf("ParseToolCall failed: %v", err)
		}
	}
}

// BenchmarkParseComplexDiff tests parsing performance with multiple edits
func BenchmarkParseComplexDiff(b *testing.B) {
	xmlInput := `<tool>
<server_name>local</server_name>
<tool_name>apply_diff</tool_name>
<arguments>
  <path>app.ts</path>
  <edits>
    <edit>
      <search><![CDATA[const oldConfig = {
	api: "http://old.com",
	key: "old_key"
};]]></search>
      <replace><![CDATA[const newConfig = {
	api: "https://new.com",
	key: "new_key",
	verbose: true
};]]></replace>
    </edit>
    <edit>
      <search><![CDATA[function process(data) {
	return data.map(x => x * 2);
}]]></search>
      <replace><![CDATA[function process(data) {
	return data
		.filter(x => x > 0)
		.map(x => x * 2);
}]]></replace>
    </edit>
    <edit>
      <search><![CDATA[export default App;]]></search>
      <replace><![CDATA[export { App as default, Config };]]></replace>
    </edit>
  </edits>
</arguments>
</tool>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := ParseToolCall(xmlInput)
		if err != nil {
			b.Fatalf("ParseToolCall failed: %v", err)
		}
	}
}

// BenchmarkHasToolCall measures the performance of checking for tool calls
func BenchmarkHasToolCall(b *testing.B) {
	text := `This is some thinking text before the tool call.
<tool>
<server_name>local</server_name>
<tool_name>test</tool_name>
<arguments>
  <value>test</value>
</arguments>
</tool>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = HasToolCall(text)
	}
}

// BenchmarkExtractThinkingAndToolCall measures the performance of extracting both
func BenchmarkExtractThinkingAndToolCall(b *testing.B) {
	text := `Here's my thinking about this problem. I need to analyze the code
and determine the best approach for the solution.

After careful consideration, I'll use the following tool:

<tool>
<server_name>local</server_name>
<tool_name>apply_diff</tool_name>
<arguments>
  <path>main.go</path>
  <edits>
    <edit>
      <search><![CDATA[func old() {}]]></search>
      <replace><![CDATA[func new() {}]]></replace>
    </edit>
  </edits>
</arguments>
</tool>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, err := ExtractThinkingAndToolCall(text)
		if err != nil {
			b.Fatalf("ExtractThinkingAndToolCall failed: %v", err)
		}
	}
}
