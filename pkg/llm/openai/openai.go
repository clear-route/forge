// Package openai provides an OpenAI-compatible LLM provider implementation.
package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/clear-route/forge/pkg/llm"
	"github.com/clear-route/forge/pkg/types"
	"github.com/openai/openai-go"
)

// Provider implements the LLM provider interface for OpenAI-compatible APIs.
type Provider struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
	model      string
	modelInfo  *types.ModelInfo
}

// ProviderOption is a function that configures a Provider.
type ProviderOption func(*Provider)

// WithModel sets the model to use for completions.
func WithModel(model string) ProviderOption {
	return func(p *Provider) {
		p.model = model
	}
}

// WithBaseURL sets a custom base URL for OpenAI-compatible APIs.
// This enables using Azure OpenAI, local models, or other compatible services.
func WithBaseURL(baseURL string) ProviderOption {
	return func(p *Provider) {
		// Store in modelInfo metadata for reference
		if p.modelInfo.Metadata == nil {
			p.modelInfo.Metadata = make(map[string]interface{})
		}
		p.modelInfo.Metadata["base_url"] = baseURL
	}
}

// NewProvider creates a new OpenAI provider with the given API key.
//
// If apiKey is empty, it will attempt to read from the OPENAI_API_KEY environment variable.
// If baseURL is not provided via WithBaseURL option, it will check OPENAI_BASE_URL environment variable.
//
// The default model is "gpt-4".
//
// Example:
//
//	// Standard OpenAI
//	provider, _ := openai.NewProvider("sk-...", openai.WithModel("gpt-4"))
//
//	// Azure OpenAI
//	provider, _ := openai.NewProvider("your-key",
//	    openai.WithBaseURL("https://your-resource.openai.azure.com"),
//	    openai.WithModel("gpt-4o"))
//
//	// Local OpenAI-compatible API
//	provider, _ := openai.NewProvider("local",
//	    openai.WithBaseURL("http://localhost:8080/v1"))
func NewProvider(apiKey string, opts ...ProviderOption) (*Provider, error) {
	// Use environment variable if no API key provided
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required (provide via parameter or OPENAI_API_KEY environment variable)")
	}

	// Check for custom base URL from environment
	baseURL := os.Getenv("OPENAI_BASE_URL")

	// Create provider with defaults
	p := &Provider{
		model:      "gpt-4o", // Default model
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}

	// Apply options (may override baseURL via WithBaseURL)
	for _, opt := range opts {
		opt(p)
	}

	// Check if baseURL was set via options (stored in metadata)
	if p.modelInfo != nil && p.modelInfo.Metadata != nil {
		if customURL, ok := p.modelInfo.Metadata["base_url"].(string); ok {
			baseURL = customURL
		}
	}

	// Store base URL
	p.baseURL = baseURL
	if p.baseURL == "" {
		p.baseURL = "https://api.openai.com/v1"
	}

	// Initialize model info (if not already set by options)
	if p.modelInfo == nil {
		p.modelInfo = &types.ModelInfo{
			Metadata: make(map[string]interface{}),
		}
	}

	p.modelInfo.Provider = "openai"
	p.modelInfo.Name = p.model
	p.modelInfo.SupportsStreaming = true
	p.modelInfo.MaxTokens = 8192 // Default, varies by model

	// Store base URL in metadata if set
	if baseURL != "" {
		p.modelInfo.Metadata["base_url"] = baseURL
	}

	return p, nil
}

// StreamCompletion sends messages to the OpenAI API and streams back response chunks.
//
// The returned channel emits StreamChunk instances as the response is generated.
// The channel is closed when streaming completes or an error occurs.
//
// This implementation uses raw HTTP streaming to handle SSE events directly,
// which provides better compatibility with OpenAI-compatible APIs that may
// include SSE comments or have slight format variations.
func (p *Provider) StreamCompletion(ctx context.Context, messages []*types.Message) (<-chan *llm.StreamChunk, error) {
	// Convert our message format to OpenAI format
	openaiMessages := convertToOpenAIMessages(messages)

	// Create request body
	reqBody := map[string]interface{}{
		"model":    p.model,
		"messages": openaiMessages,
		"stream":   true,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := p.baseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Accept", "text/event-stream")

	// Send request
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Create output channel
	chunks := make(chan *llm.StreamChunk, 10)

	// Stream in goroutine
	go func() {
		defer close(chunks)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		firstChunk := true

		for scanner.Scan() {
			line := scanner.Text()

			// Skip empty lines and SSE comments
			if line == "" || strings.HasPrefix(line, ":") {
				continue
			}

			// SSE data lines start with "data: "
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			// Extract JSON data
			data := strings.TrimPrefix(line, "data: ")

			// Check for stream end
			if data == "[DONE]" {
				chunks <- &llm.StreamChunk{Finished: true}
				return
			}

			// Parse JSON chunk
			var chunk struct {
				Choices []struct {
					Delta struct {
						Role    string `json:"role"`
						Content string `json:"content"`
					} `json:"delta"`
					FinishReason *string `json:"finish_reason"`
				} `json:"choices"`
			}

			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				// Skip malformed chunks silently
				continue
			}

			if len(chunk.Choices) == 0 {
				continue
			}

			delta := chunk.Choices[0].Delta

			// Create stream chunk
			streamChunk := &llm.StreamChunk{}

			// Set role on first chunk
			if firstChunk && delta.Role != "" {
				streamChunk.Role = delta.Role
				firstChunk = false
			}

			// Set content if present
			if delta.Content != "" {
				streamChunk.Content = delta.Content
			}

			// Check for finish
			if chunk.Choices[0].FinishReason != nil && *chunk.Choices[0].FinishReason == "stop" {
				streamChunk.Finished = true
			}

			// Only send if there's content, role, or finished flag
			if streamChunk.Content != "" || streamChunk.Role != "" || streamChunk.Finished {
				select {
				case chunks <- streamChunk:
				case <-ctx.Done():
					chunks <- &llm.StreamChunk{Error: ctx.Err()}
					return
				}
			}
		}

		if err := scanner.Err(); err != nil {
			chunks <- &llm.StreamChunk{Error: fmt.Errorf("stream read error: %w", err)}
		}
	}()

	return chunks, nil
}

// Complete sends messages to the OpenAI API and returns the full response.
//
// This is a convenience wrapper around StreamCompletion that accumulates
// all chunks into a single message.
func (p *Provider) Complete(ctx context.Context, messages []*types.Message) (*types.Message, error) {
	stream, err := p.StreamCompletion(ctx, messages)
	if err != nil {
		return nil, err
	}

	var content string
	var role string

	for chunk := range stream {
		if chunk.IsError() {
			return nil, chunk.Error
		}

		if chunk.Role != "" {
			role = chunk.Role
		}

		content += chunk.Content
	}

	// Default to assistant role if not set
	if role == "" {
		role = string(types.RoleAssistant)
	}

	return &types.Message{
		Role:    types.MessageRole(role),
		Content: content,
	}, nil
}

// GetModelInfo returns information about the OpenAI model being used.
func (p *Provider) GetModelInfo() *types.ModelInfo {
	return p.modelInfo
}

// convertToOpenAIMessages converts our Message format to OpenAI's ChatCompletionMessageParamUnion format.
func convertToOpenAIMessages(messages []*types.Message) []openai.ChatCompletionMessageParamUnion {
	openaiMessages := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages))

	for _, msg := range messages {
		switch msg.Role {
		case types.RoleSystem:
			openaiMessages = append(openaiMessages, openai.SystemMessage(msg.Content))
		case types.RoleUser:
			openaiMessages = append(openaiMessages, openai.UserMessage(msg.Content))
		case types.RoleAssistant:
			openaiMessages = append(openaiMessages, openai.AssistantMessage(msg.Content))
		default:
			// Default to user message for unknown roles
			openaiMessages = append(openaiMessages, openai.UserMessage(msg.Content))
		}
	}

	return openaiMessages
}
