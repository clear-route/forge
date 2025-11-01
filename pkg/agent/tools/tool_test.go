package tools

import (
	"context"
	"encoding/json"
	"testing"
)

func TestTaskCompletionTool(t *testing.T) {
	tool := NewTaskCompletionTool()

	t.Run("Name", func(t *testing.T) {
		if tool.Name() != "task_completion" {
			t.Errorf("expected name 'task_completion', got '%s'", tool.Name())
		}
	})

	t.Run("IsLoopBreaking", func(t *testing.T) {
		if !tool.IsLoopBreaking() {
			t.Error("task_completion should be loop-breaking")
		}
	})

	t.Run("Execute_Success", func(t *testing.T) {
		args := json.RawMessage(`{"result": "Task completed successfully!"}`)
		result, err := tool.Execute(context.Background(), args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "Task completed successfully!" {
			t.Errorf("expected 'Task completed successfully!', got '%s'", result)
		}
	})

	t.Run("Execute_EmptyResult", func(t *testing.T) {
		args := json.RawMessage(`{"result": ""}`)
		_, err := tool.Execute(context.Background(), args)
		if err == nil {
			t.Error("expected error for empty result")
		}
	})

	t.Run("Execute_InvalidJSON", func(t *testing.T) {
		args := json.RawMessage(`invalid json`)
		_, err := tool.Execute(context.Background(), args)
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})
}

func TestAskQuestionTool(t *testing.T) {
	tool := NewAskQuestionTool()

	t.Run("Name", func(t *testing.T) {
		if tool.Name() != "ask_question" {
			t.Errorf("expected name 'ask_question', got '%s'", tool.Name())
		}
	})

	t.Run("IsLoopBreaking", func(t *testing.T) {
		if !tool.IsLoopBreaking() {
			t.Error("ask_question should be loop-breaking")
		}
	})

	t.Run("Execute_WithoutSuggestions", func(t *testing.T) {
		args := json.RawMessage(`{"question": "What is your preferred color?"}`)
		result, err := tool.Execute(context.Background(), args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "What is your preferred color?" {
			t.Errorf("expected question only, got '%s'", result)
		}
	})

	t.Run("Execute_WithSuggestions", func(t *testing.T) {
		args := json.RawMessage(`{"question": "What is your preferred color?", "suggestions": ["Red", "Blue", "Green"]}`)
		result, err := tool.Execute(context.Background(), args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "What is your preferred color?\n\nSuggested answers:\n1. Red\n2. Blue\n3. Green"
		if result != expected {
			t.Errorf("expected formatted question with suggestions, got '%s'", result)
		}
	})

	t.Run("Execute_EmptyQuestion", func(t *testing.T) {
		args := json.RawMessage(`{"question": ""}`)
		_, err := tool.Execute(context.Background(), args)
		if err == nil {
			t.Error("expected error for empty question")
		}
	})
}

func TestConverseTool(t *testing.T) {
	tool := NewConverseTool()

	t.Run("Name", func(t *testing.T) {
		if tool.Name() != "converse" {
			t.Errorf("expected name 'converse', got '%s'", tool.Name())
		}
	})

	t.Run("IsLoopBreaking", func(t *testing.T) {
		if !tool.IsLoopBreaking() {
			t.Error("converse should be loop-breaking")
		}
	})

	t.Run("Execute_Success", func(t *testing.T) {
		args := json.RawMessage(`{"message": "Hello! How can I help you today?"}`)
		result, err := tool.Execute(context.Background(), args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "Hello! How can I help you today?" {
			t.Errorf("expected message, got '%s'", result)
		}
	})

	t.Run("Execute_EmptyMessage", func(t *testing.T) {
		args := json.RawMessage(`{"message": ""}`)
		_, err := tool.Execute(context.Background(), args)
		if err == nil {
			t.Error("expected error for empty message")
		}
	})
}

func TestBaseToolSchema(t *testing.T) {
	properties := map[string]interface{}{
		"name": map[string]interface{}{
			"type":        "string",
			"description": "The name",
		},
	}
	required := []string{"name"}

	schema := BaseToolSchema(properties, required)

	if schema["type"] != "object" {
		t.Errorf("expected type 'object', got '%v'", schema["type"])
	}

	if _, ok := schema["properties"]; !ok {
		t.Error("schema should have 'properties' field")
	}

	if _, ok := schema["required"]; !ok {
		t.Error("schema should have 'required' field")
	}
}
