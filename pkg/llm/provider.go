// Package llm provides abstractions for LLM provider integration.
package llm

import (
	"context"

	"github.com/entrhq/forge/pkg/types"
)

// Provider defines the interface for LLM integrations.
//
// Providers handle API communication with LLM services and return simple
// StreamChunk instances. This design keeps providers focused on LLM concerns
// without coupling them to agent-level events or orchestration.
//
// The Agent layer is responsible for:
// - Converting StreamChunks to AgentEvents
// - Emitting thinking, tool, and status events
// - Managing conversation state and history
//
// This separation allows providers to be:
// - Reusable in non-agent contexts (CLI tools, batch processing, etc.)
// - Testable independently of agent logic
// - Simpler to implement and maintain
type Provider interface {
	// StreamCompletion sends messages to the LLM and streams back response chunks.
	//
	// The returned channel emits StreamChunk instances:
	// - First chunk typically has Role set (e.g., "assistant")
	// - Subsequent chunks contain Content deltas
	// - Final chunk has Finished=true
	// - Error chunks have Error set
	//
	// The channel is closed when streaming completes or an error occurs.
	// Callers should continue reading until the channel is closed.
	//
	// Returns an error only if streaming cannot be initiated (e.g., invalid
	// configuration, network unavailable). Stream-time errors are sent as
	// StreamChunk instances with Error set.
	//
	// Example usage:
	//   stream, err := provider.StreamCompletion(ctx, messages)
	//   if err != nil {
	//       return err
	//   }
	//   for chunk := range stream {
	//       if chunk.IsError() {
	//           return chunk.Error
	//       }
	//       fmt.Print(chunk.Content)
	//   }
	StreamCompletion(ctx context.Context, messages []*types.Message) (<-chan *StreamChunk, error)

	// Complete sends messages to the LLM and returns the full response.
	//
	// This is a convenience wrapper around StreamCompletion for non-streaming
	// use cases. It accumulates all chunks and returns the complete message.
	//
	// Returns the assistant's response message or an error.
	Complete(ctx context.Context, messages []*types.Message) (*types.Message, error)

	// GetModelInfo returns information about the LLM model being used.
	//
	// This can be used to inspect model capabilities, pricing, token limits,
	// and other metadata.
	GetModelInfo() *types.ModelInfo
}
