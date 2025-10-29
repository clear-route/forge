// Package agent provides the core agent interface and implementation
// for the Forge agent framework.
package agent

import (
	"context"

	"github.com/clear-route/forge/pkg/types"
)

// Agent interface defines the core capabilities of a Forge agent.
// Agents are async event-driven components that process messages through
// an LLM provider and communicate via channels.
type Agent interface {
	// Start begins the agent's event loop in a goroutine.
	// The agent will listen for messages on its input channel and process them
	// asynchronously, sending responses to the output channel.
	//
	// The agent runs until:
	// - The context is canceled
	// - The shutdown channel is closed
	// - An unrecoverable error occurs
	//
	// Returns an error if the agent fails to start, otherwise returns nil
	// and continues running asynchronously.
	Start(ctx context.Context) error

	// Shutdown gracefully stops the agent.
	// This method signals the agent to stop processing new messages and
	// complete any in-flight operations before shutting down.
	//
	// Returns when the agent has fully stopped or the context is canceled.
	Shutdown(ctx context.Context) error

	// GetChannels returns the communication channels for this agent.
	// The executor uses these channels to send input and receive output.
	GetChannels() *types.AgentChannels
}
