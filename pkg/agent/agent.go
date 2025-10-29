// Package agent provides the core agent interface and implementation
// for the Forge agent framework.
//
// Example usage:
//
//	package main
//
//	import (
//	    "context"
//	    "log"
//	    "os"
//
//	    "github.com/entrhq/forge/pkg/agent"
//	    "github.com/entrhq/forge/pkg/executor/cli"
//	    "github.com/entrhq/forge/pkg/llm/openai"
//	)
//
//	func main() {
//	    // Create provider
//	    provider, err := openai.NewProvider(
//	        os.Getenv("OPENAI_API_KEY"),
//	        openai.WithModel("gpt-4o"),
//	    )
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    // Create agent
//	    ag := agent.NewDefaultAgent(provider,
//	        agent.WithSystemPrompt("You are a helpful assistant."),
//	    )
//
//	    // Create executor and run
//	    executor := cli.NewExecutor(ag)
//	    if err := executor.Run(context.Background()); err != nil {
//	        log.Fatal(err)
//	    }
//	}
package agent

import (
	"context"

	"github.com/entrhq/forge/pkg/types"
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
