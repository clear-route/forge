package llm

import (
	"errors"
	"testing"
)

func TestStreamChunk_IsError(t *testing.T) {
	tests := []struct {
		name     string
		chunk    *StreamChunk
		expected bool
	}{
		{
			name:     "chunk with error",
			chunk:    &StreamChunk{Error: errors.New("test error")},
			expected: true,
		},
		{
			name:     "chunk without error",
			chunk:    &StreamChunk{Content: "hello"},
			expected: false,
		},
		{
			name:     "empty chunk",
			chunk:    &StreamChunk{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.chunk.IsError(); got != tt.expected {
				t.Errorf("IsError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStreamChunk_IsLast(t *testing.T) {
	tests := []struct {
		name     string
		chunk    *StreamChunk
		expected bool
	}{
		{
			name:     "finished chunk",
			chunk:    &StreamChunk{Finished: true},
			expected: true,
		},
		{
			name:     "unfinished chunk",
			chunk:    &StreamChunk{Content: "hello"},
			expected: false,
		},
		{
			name:     "empty chunk",
			chunk:    &StreamChunk{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.chunk.IsLast(); got != tt.expected {
				t.Errorf("IsLast() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStreamChunk_HasContent(t *testing.T) {
	tests := []struct {
		name     string
		chunk    *StreamChunk
		expected bool
	}{
		{
			name:     "chunk with content",
			chunk:    &StreamChunk{Content: "hello"},
			expected: true,
		},
		{
			name:     "chunk with empty string",
			chunk:    &StreamChunk{Content: ""},
			expected: false,
		},
		{
			name:     "chunk with only role",
			chunk:    &StreamChunk{Role: "assistant"},
			expected: false,
		},
		{
			name:     "empty chunk",
			chunk:    &StreamChunk{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.chunk.HasContent(); got != tt.expected {
				t.Errorf("HasContent() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStreamChunk_MultipleFlags(t *testing.T) {
	t.Run("chunk can have content and be finished", func(t *testing.T) {
		chunk := &StreamChunk{
			Content:  "final message",
			Finished: true,
		}

		if !chunk.HasContent() {
			t.Error("expected chunk to have content")
		}
		if !chunk.IsLast() {
			t.Error("expected chunk to be last")
		}
		if chunk.IsError() {
			t.Error("expected chunk to not have error")
		}
	})

	t.Run("chunk can have role and content", func(t *testing.T) {
		chunk := &StreamChunk{
			Role:    "assistant",
			Content: "hello",
		}

		if chunk.Role != "assistant" {
			t.Errorf("expected role to be 'assistant', got %q", chunk.Role)
		}
		if !chunk.HasContent() {
			t.Error("expected chunk to have content")
		}
	})

	t.Run("error chunk can be finished", func(t *testing.T) {
		chunk := &StreamChunk{
			Error:    errors.New("test error"),
			Finished: true,
		}

		if !chunk.IsError() {
			t.Error("expected chunk to have error")
		}
		if !chunk.IsLast() {
			t.Error("expected chunk to be last")
		}
	})
}