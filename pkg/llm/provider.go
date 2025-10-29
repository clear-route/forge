// Package llm provides abstractions for LLM provider integration.
package llm

import (
	"context"

	"github.com/clear-route/forge/pkg/types"
)

// Provider interface for pluggable LLM implementations.
// Providers abstract the interaction with different LLM services,
// enabling agents to swap between providers at runtime.
type Provider interface {
	// GenerateStream sends messages to the LLM and returns a channel
	// that streams response content asynchronously as AgentEvents.
	//
	// The returned channel will receive AgentEvent instances:
	// - EventTypeMessageStart when the response begins
	// - EventTypeMessageContent for each content chunk (streaming)
	// - EventTypeMessageEnd when the response is complete
	// - EventTypeError if an error occurs
	//
	// The channel is closed when streaming completes or an error occurs.
	// Callers should continue reading until the channel is closed.
	//
	// Returns an error if streaming cannot be initiated.
	GenerateStream(ctx context.Context, messages []*types.Message) (<-chan *types.AgentEvent, error)

	// Generate sends messages to the LLM and returns the complete response.
	// This is a convenience method for non-streaming use cases.
	//
	// Returns the assistant's response message or an error.
	Generate(ctx context.Context, messages []*types.Message) (*types.Message, error)

	// GetModelInfo returns information about the LLM model being used.
	GetModelInfo() *types.ModelInfo
}
