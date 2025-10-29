package types

import (
	"testing"
	"time"
)

func TestNewAgentConfig(t *testing.T) {
	config := NewAgentConfig()

	if config.MaxTurns != 0 {
		t.Errorf("NewAgentConfig MaxTurns = %v, want 0 (unlimited)", config.MaxTurns)
	}
	if config.Timeout != 0 {
		t.Errorf("NewAgentConfig Timeout = %v, want 0 (no timeout)", config.Timeout)
	}
	if !config.EnableStreaming {
		t.Error("NewAgentConfig EnableStreaming should be true by default")
	}
	if config.BufferSize != 10 {
		t.Errorf("NewAgentConfig BufferSize = %v, want 10", config.BufferSize)
	}
	if config.Metadata == nil {
		t.Error("NewAgentConfig metadata should be initialized")
	}
}

func TestAgentConfigWithSystemPrompt(t *testing.T) {
	prompt := "You are a helpful assistant"
	config := NewAgentConfig().WithSystemPrompt(prompt)

	if config.SystemPrompt != prompt {
		t.Errorf("WithSystemPrompt = %v, want %v", config.SystemPrompt, prompt)
	}
}

func TestAgentConfigWithMaxTurns(t *testing.T) {
	maxTurns := 5
	config := NewAgentConfig().WithMaxTurns(maxTurns)

	if config.MaxTurns != maxTurns {
		t.Errorf("WithMaxTurns = %v, want %v", config.MaxTurns, maxTurns)
	}
}

func TestAgentConfigWithTimeout(t *testing.T) {
	timeout := 30 * time.Second
	config := NewAgentConfig().WithTimeout(timeout)

	if config.Timeout != timeout {
		t.Errorf("WithTimeout = %v, want %v", config.Timeout, timeout)
	}
}

func TestAgentConfigWithStreaming(t *testing.T) {
	config := NewAgentConfig().WithStreaming(false)

	if config.EnableStreaming {
		t.Error("WithStreaming(false) should disable streaming")
	}

	config.WithStreaming(true)
	if !config.EnableStreaming {
		t.Error("WithStreaming(true) should enable streaming")
	}
}

func TestAgentConfigWithBufferSize(t *testing.T) {
	bufferSize := 20
	config := NewAgentConfig().WithBufferSize(bufferSize)

	if config.BufferSize != bufferSize {
		t.Errorf("WithBufferSize = %v, want %v", config.BufferSize, bufferSize)
	}
}

func TestAgentConfigWithMetadata(t *testing.T) {
	config := NewAgentConfig()
	key := "test_key"
	value := "test_value"

	result := config.WithMetadata(key, value)

	if result != config {
		t.Error("WithMetadata should return the same config for chaining")
	}
	if config.Metadata[key] != value {
		t.Errorf("WithMetadata did not set metadata correctly, got %v, want %v", config.Metadata[key], value)
	}
}

func TestAgentConfigChaining(t *testing.T) {
	config := NewAgentConfig().
		WithSystemPrompt("test prompt").
		WithMaxTurns(10).
		WithTimeout(1*time.Minute).
		WithStreaming(false).
		WithBufferSize(15).
		WithMetadata("key", "value")

	if config.SystemPrompt != "test prompt" {
		t.Error("Chaining failed for SystemPrompt")
	}
	if config.MaxTurns != 10 {
		t.Error("Chaining failed for MaxTurns")
	}
	if config.Timeout != 1*time.Minute {
		t.Error("Chaining failed for Timeout")
	}
	if config.EnableStreaming {
		t.Error("Chaining failed for EnableStreaming")
	}
	if config.BufferSize != 15 {
		t.Error("Chaining failed for BufferSize")
	}
	if config.Metadata["key"] != "value" {
		t.Error("Chaining failed for Metadata")
	}
}

func TestModelInfo(t *testing.T) {
	info := &ModelInfo{
		Name:              "gpt-4",
		MaxTokens:         8192,
		SupportsStreaming: true,
		Provider:          "openai",
	}

	if info.Name != "gpt-4" {
		t.Errorf("ModelInfo Name = %v, want gpt-4", info.Name)
	}
	if info.MaxTokens != 8192 {
		t.Errorf("ModelInfo MaxTokens = %v, want 8192", info.MaxTokens)
	}
	if !info.SupportsStreaming {
		t.Error("ModelInfo SupportsStreaming should be true")
	}
	if info.Provider != "openai" {
		t.Errorf("ModelInfo Provider = %v, want openai", info.Provider)
	}
}
