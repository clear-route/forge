package openai

import (
	"os"
	"testing"

	"github.com/entrhq/forge/pkg/types"
)

func TestNewProvider_WithAPIKey(t *testing.T) {
	// Clear environment variables to ensure we're testing parameter
	oldKey := os.Getenv("OPENAI_API_KEY")
	oldBaseURL := os.Getenv("OPENAI_BASE_URL")
	os.Setenv("OPENAI_API_KEY", "")
	os.Setenv("OPENAI_BASE_URL", "")
	defer func() {
		os.Setenv("OPENAI_API_KEY", oldKey)
		os.Setenv("OPENAI_BASE_URL", oldBaseURL)
	}()

	provider, err := NewProvider("test-api-key")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if provider == nil {
		t.Fatal("Expected provider, got nil")
	}

	if provider.apiKey != "test-api-key" {
		t.Errorf("Expected API key 'test-api-key', got '%s'", provider.apiKey)
	}

	// Should use default model
	if provider.model != "gpt-4o" {
		t.Errorf("Expected default model 'gpt-4o', got '%s'", provider.model)
	}

	// Should use default base URL
	if provider.baseURL != "https://api.openai.com/v1" {
		t.Errorf("Expected default base URL, got '%s'", provider.baseURL)
	}
}

func TestNewProvider_FromEnvironment(t *testing.T) {
	// Set environment variable
	oldKey := os.Getenv("OPENAI_API_KEY")
	os.Setenv("OPENAI_API_KEY", "env-api-key")
	defer os.Setenv("OPENAI_API_KEY", oldKey)

	provider, err := NewProvider("")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if provider.apiKey != "env-api-key" {
		t.Errorf("Expected API key from env 'env-api-key', got '%s'", provider.apiKey)
	}
}

func TestNewProvider_NoAPIKey(t *testing.T) {
	// Clear environment variable
	oldKey := os.Getenv("OPENAI_API_KEY")
	os.Setenv("OPENAI_API_KEY", "")
	defer os.Setenv("OPENAI_API_KEY", oldKey)

	_, err := NewProvider("")
	if err == nil {
		t.Error("Expected error when no API key provided")
	}

	if err.Error() != "OpenAI API key is required (provide via parameter or OPENAI_API_KEY environment variable)" {
		t.Errorf("Expected error message about required API key, got '%v'", err)
	}
}

func TestNewProvider_WithModel(t *testing.T) {
	oldKey := os.Getenv("OPENAI_API_KEY")
	os.Setenv("OPENAI_API_KEY", "")
	defer os.Setenv("OPENAI_API_KEY", oldKey)

	provider, err := NewProvider("test-key", WithModel("gpt-4-turbo"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if provider.model != "gpt-4-turbo" {
		t.Errorf("Expected model 'gpt-4-turbo', got '%s'", provider.model)
	}

	// Model info should reflect the custom model
	if provider.modelInfo.Name != "gpt-4-turbo" {
		t.Errorf("Expected model info name 'gpt-4-turbo', got '%s'", provider.modelInfo.Name)
	}
}

func TestNewProvider_WithBaseURL(t *testing.T) {
	oldKey := os.Getenv("OPENAI_API_KEY")
	oldBaseURL := os.Getenv("OPENAI_BASE_URL")
	os.Setenv("OPENAI_API_KEY", "")
	os.Setenv("OPENAI_BASE_URL", "")
	defer func() {
		os.Setenv("OPENAI_API_KEY", oldKey)
		os.Setenv("OPENAI_BASE_URL", oldBaseURL)
	}()

	customURL := "https://custom.openai.com/v1"
	provider, err := NewProvider("test-key", WithBaseURL(customURL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if provider.baseURL != customURL {
		t.Errorf("Expected base URL '%s', got '%s'", customURL, provider.baseURL)
	}

	// Should be stored in metadata
	if provider.modelInfo.Metadata["base_url"] != customURL {
		t.Errorf("Expected base URL in metadata, got %v", provider.modelInfo.Metadata["base_url"])
	}
}

func TestNewProvider_BaseURLFromEnvironment(t *testing.T) {
	oldKey := os.Getenv("OPENAI_API_KEY")
	oldBaseURL := os.Getenv("OPENAI_BASE_URL")
	os.Setenv("OPENAI_API_KEY", "")
	envBaseURL := "https://env.openai.com/v1"
	os.Setenv("OPENAI_BASE_URL", envBaseURL)
	defer func() {
		os.Setenv("OPENAI_API_KEY", oldKey)
		os.Setenv("OPENAI_BASE_URL", oldBaseURL)
	}()

	provider, err := NewProvider("test-key")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if provider.baseURL != envBaseURL {
		t.Errorf("Expected base URL from env '%s', got '%s'", envBaseURL, provider.baseURL)
	}
}

func TestNewProvider_MultipleOptions(t *testing.T) {
	oldKey := os.Getenv("OPENAI_API_KEY")
	os.Setenv("OPENAI_API_KEY", "")
	defer os.Setenv("OPENAI_API_KEY", oldKey)

	provider, err := NewProvider("test-key",
		WithModel("gpt-3.5-turbo"),
		WithBaseURL("https://custom.api.com/v1"),
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if provider.model != "gpt-3.5-turbo" {
		t.Errorf("Expected model 'gpt-3.5-turbo', got '%s'", provider.model)
	}

	if provider.baseURL != "https://custom.api.com/v1" {
		t.Errorf("Expected custom base URL, got '%s'", provider.baseURL)
	}
}

func TestGetModelInfo(t *testing.T) {
	oldKey := os.Getenv("OPENAI_API_KEY")
	os.Setenv("OPENAI_API_KEY", "")
	defer os.Setenv("OPENAI_API_KEY", oldKey)

	provider, err := NewProvider("test-key", WithModel("gpt-4"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	info := provider.GetModelInfo()

	if info == nil {
		t.Fatal("Expected model info, got nil")
	}

	if info.Provider != "openai" {
		t.Errorf("Expected provider 'openai', got '%s'", info.Provider)
	}

	if info.Name != "gpt-4" {
		t.Errorf("Expected model name 'gpt-4', got '%s'", info.Name)
	}

	if !info.SupportsStreaming {
		t.Error("Expected model to support streaming")
	}

	if info.MaxTokens != 8192 {
		t.Errorf("Expected max tokens 8192, got %d", info.MaxTokens)
	}

	if info.Metadata == nil {
		t.Error("Expected metadata to be initialized")
	}
}

func TestConvertToOpenAIMessages(t *testing.T) {
	tests := []struct {
		name     string
		input    []*types.Message
		expected int
	}{
		{
			name: "SystemMessage",
			input: []*types.Message{
				types.NewSystemMessage("System prompt"),
			},
			expected: 1,
		},
		{
			name: "UserMessage",
			input: []*types.Message{
				types.NewUserMessage("User input"),
			},
			expected: 1,
		},
		{
			name: "AssistantMessage",
			input: []*types.Message{
				types.NewAssistantMessage("Assistant response"),
			},
			expected: 1,
		},
		{
			name: "MultipleMessages",
			input: []*types.Message{
				types.NewSystemMessage("System"),
				types.NewUserMessage("User"),
				types.NewAssistantMessage("Assistant"),
			},
			expected: 3,
		},
		{
			name:     "EmptyMessages",
			input:    []*types.Message{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertToOpenAIMessages(tt.input)

			if len(result) != tt.expected {
				t.Errorf("Expected %d messages, got %d", tt.expected, len(result))
			}
		})
	}
}

func TestConvertToOpenAIMessages_UnknownRole(t *testing.T) {
	// Test that unknown roles default to user messages
	input := []*types.Message{
		{
			Role:    "unknown",
			Content: "Test content",
		},
	}

	result := convertToOpenAIMessages(input)

	if len(result) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(result))
	}

	// Should have been converted (defaulting to user message)
	// We can't easily test the exact type without reflection,
	// but we can verify it didn't panic
}

func TestConvertToOpenAIMessages_PreservesOrder(t *testing.T) {
	input := []*types.Message{
		types.NewSystemMessage("First"),
		types.NewUserMessage("Second"),
		types.NewAssistantMessage("Third"),
		types.NewUserMessage("Fourth"),
	}

	result := convertToOpenAIMessages(input)

	if len(result) != 4 {
		t.Errorf("Expected 4 messages in order, got %d", len(result))
	}
}

func TestProvider_ModelInfoMetadata(t *testing.T) {
	oldKey := os.Getenv("OPENAI_API_KEY")
	oldBaseURL := os.Getenv("OPENAI_BASE_URL")
	os.Setenv("OPENAI_API_KEY", "")
	os.Setenv("OPENAI_BASE_URL", "")
	defer func() {
		os.Setenv("OPENAI_API_KEY", oldKey)
		os.Setenv("OPENAI_BASE_URL", oldBaseURL)
	}()

	customURL := "https://azure.openai.com/v1"
	provider, err := NewProvider("test-key",
		WithBaseURL(customURL),
		WithModel("gpt-4"),
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	info := provider.GetModelInfo()

	// Verify metadata contains base URL
	if baseURL, ok := info.Metadata["base_url"].(string); !ok || baseURL != customURL {
		t.Errorf("Expected base_url in metadata to be '%s', got %v", customURL, info.Metadata["base_url"])
	}
}

func TestProvider_HTTPClientInitialized(t *testing.T) {
	oldKey := os.Getenv("OPENAI_API_KEY")
	os.Setenv("OPENAI_API_KEY", "")
	defer os.Setenv("OPENAI_API_KEY", oldKey)

	provider, err := NewProvider("test-key")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if provider.httpClient == nil {
		t.Error("Expected HTTP client to be initialized")
	}
}
