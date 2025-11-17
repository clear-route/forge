package tools

import (
	"testing"
)

func TestParseToolCall(t *testing.T) {
	t.Run("ValidToolCall", func(t *testing.T) {
		text := `Some thinking here
<tool>
<server_name>local</server_name>
<tool_name>task_completion</tool_name>
<arguments>
  <result>Done!</result>
</arguments>
</tool>
Some remaining text`

		toolCall, remaining, err := ParseToolCall(text)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if toolCall.ServerName != "local" {
			t.Errorf("expected server_name 'local', got '%s'", toolCall.ServerName)
		}
		if toolCall.ToolName != "task_completion" {
			t.Errorf("expected tool_name 'task_completion', got '%s'", toolCall.ToolName)
		}

		expectedRemaining := "Some thinking here\n\nSome remaining text"
		if remaining != expectedRemaining {
			t.Errorf("expected remaining '%s', got '%s'", expectedRemaining, remaining)
		}
	})

	t.Run("ToolCallWithoutServerName", func(t *testing.T) {
		text := `<tool>
<tool_name>ask_question</tool_name>
<arguments>
  <question>What?</question>
</arguments>
</tool>`

		toolCall, _, err := ParseToolCall(text)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should default to "local"
		if toolCall.ServerName != "local" {
			t.Errorf("expected default server_name 'local', got '%s'", toolCall.ServerName)
		}
		if toolCall.ToolName != "ask_question" {
			t.Errorf("expected tool_name 'ask_question', got '%s'", toolCall.ToolName)
		}
	})

	t.Run("NoToolCall", func(t *testing.T) {
		text := "Just some regular text without a tool call"

		_, remaining, err := ParseToolCall(text)
		if err == nil {
			t.Error("expected error when no tool call present")
		}
		if remaining != text {
			t.Error("remaining text should be unchanged when no tool call found")
		}
	})

	t.Run("InvalidXML", func(t *testing.T) {
		text := `<tool>not valid xml</tool>`

		_, _, err := ParseToolCall(text)
		if err == nil {
			t.Error("expected error for invalid XML")
		}
	})

	t.Run("MissingToolName", func(t *testing.T) {
		text := `<tool>
<server_name>local</server_name>
<arguments></arguments>
</tool>`

		_, _, err := ParseToolCall(text)
		if err == nil {
			t.Error("expected error for missing tool_name")
		}
	})
}

func TestExtractThinkingAndToolCall(t *testing.T) {
	t.Run("WithThinkingAndToolCall", func(t *testing.T) {
		text := `I need to complete this task.
Let me use the completion tool.

<tool>
<tool_name>task_completion</tool_name>
<arguments>
  <result>All done!</result>
</arguments>
</tool>

After thought.`

		thinking, toolCall, remaining, err := ExtractThinkingAndToolCall(text)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expectedThinking := "I need to complete this task.\nLet me use the completion tool."
		if thinking != expectedThinking {
			t.Errorf("expected thinking '%s', got '%s'", expectedThinking, thinking)
		}

		if toolCall == nil {
			t.Fatal("expected tool call to be found")
		}
		if toolCall.ToolName != "task_completion" {
			t.Errorf("expected tool_name 'task_completion', got '%s'", toolCall.ToolName)
		}

		expectedRemaining := "After thought."
		if remaining != expectedRemaining {
			t.Errorf("expected remaining '%s', got '%s'", expectedRemaining, remaining)
		}
	})

	t.Run("OnlyToolCall", func(t *testing.T) {
		text := `<tool>
<tool_name>converse</tool_name>
<arguments>
  <message>Hi!</message>
</arguments>
</tool>`

		thinking, toolCall, _, err := ExtractThinkingAndToolCall(text)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if thinking != "" {
			t.Errorf("expected empty thinking, got '%s'", thinking)
		}
		if toolCall == nil {
			t.Fatal("expected tool call to be found")
		}
	})

	t.Run("OnlyThinking", func(t *testing.T) {
		text := "Just thinking, no tool call here."

		thinking, toolCall, remaining, err := ExtractThinkingAndToolCall(text)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if thinking != text {
			t.Errorf("expected all text as thinking, got '%s'", thinking)
		}
		if toolCall != nil {
			t.Error("expected nil tool call")
		}
		if remaining != "" {
			t.Errorf("expected empty remaining, got '%s'", remaining)
		}
	})
}

func TestHasToolCall(t *testing.T) {
	t.Run("HasToolCall", func(t *testing.T) {
		text := `Some text <tool><tool_name>test</tool_name></tool> more text`
		if !HasToolCall(text) {
			t.Error("expected HasToolCall to return true")
		}
	})

	t.Run("NoToolCall", func(t *testing.T) {
		text := "Just regular text"
		if HasToolCall(text) {
			t.Error("expected HasToolCall to return false")
		}
	})

	t.Run("IncompleteToolTag", func(t *testing.T) {
		text := "Text with <tool> but no closing tag"
		if HasToolCall(text) {
			t.Error("expected HasToolCall to return false for incomplete tag")
		}
	})
}

func TestValidateToolCall(t *testing.T) {
	t.Run("ValidToolCall", func(t *testing.T) {
		tc := &ToolCall{
			ServerName: "local",
			ToolName:   "test_tool",
			Arguments: ArgumentsBlock{
				InnerXML: []byte(``),
			},
		}
		if err := ValidateToolCall(tc); err != nil {
			t.Errorf("expected valid tool call, got error: %v", err)
		}
	})

	t.Run("NilToolCall", func(t *testing.T) {
		if err := ValidateToolCall(nil); err == nil {
			t.Error("expected error for nil tool call")
		}
	})

	t.Run("MissingToolName", func(t *testing.T) {
		tc := &ToolCall{
			ServerName: "local",
			Arguments: ArgumentsBlock{
				InnerXML: []byte(``),
			},
		}
		if err := ValidateToolCall(tc); err == nil {
			t.Error("expected error for missing tool name")
		}
	})

	t.Run("MissingServerName", func(t *testing.T) {
		tc := &ToolCall{
			ToolName: "test_tool",
			Arguments: ArgumentsBlock{
				InnerXML: []byte(``),
			},
		}
		if err := ValidateToolCall(tc); err == nil {
			t.Error("expected error for missing server name")
		}
	})
}
