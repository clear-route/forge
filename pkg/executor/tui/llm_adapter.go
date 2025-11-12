package tui

import (
	"context"
	"fmt"

	"github.com/entrhq/forge/pkg/llm"
	"github.com/entrhq/forge/pkg/types"
)

// llmAdapter adapts llm.Provider to the git.LLMClient interface
type llmAdapter struct {
	provider llm.Provider
}

// newLLMAdapter creates a new LLM adapter for git operations
func newLLMAdapter(provider llm.Provider) *llmAdapter {
	return &llmAdapter{
		provider: provider,
	}
}

// Generate implements git.LLMClient interface
func (a *llmAdapter) Generate(ctx context.Context, prompt string) (string, error) {
	if a.provider == nil {
		return "", fmt.Errorf("LLM provider not available")
	}

	// Create a simple user message with the prompt
	messages := []*types.Message{
		types.NewUserMessage(prompt),
	}

	// Use Complete to get the full response
	response, err := a.provider.Complete(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("LLM generation failed: %w", err)
	}

	return response.Content, nil
}