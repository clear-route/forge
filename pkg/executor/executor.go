// Package executor provides the execution plane abstraction
// allowing agents to run in different environments (CLI, API, etc.)
//
// Example usage:
//
//	provider, _ := openai.NewProvider(apiKey, openai.WithModel("gpt-4o"))
//	agent := agent.NewDefaultAgent(provider, agent.WithSystemPrompt("You are helpful."))
//	executor := cli.NewExecutor(agent, cli.WithPrompt("You: "))
//	executor.Run(context.Background())
package executor

import (
	"context"
)

// Executor interface defines how agents interact with their runtime environment.
// Executors provide the I/O layer for agents, handling environment-specific
// concerns like CLI interaction, HTTP requests, or other communication protocols.
//
// Executors receive the agent at construction time and manage its lifecycle.
type Executor interface {
	// Run starts the executor.
	// The executor manages the agent's lifecycle, feeding it input and
	// consuming its output based on the specific execution environment.
	//
	// This method blocks until:
	// - The context is canceled
	// - The agent completes naturally
	// - An unrecoverable error occurs
	//
	// Returns an error if execution fails, otherwise returns nil when
	// the agent completes successfully.
	Run(ctx context.Context) error

	// Stop gracefully stops the executor.
	// This signals the executor to stop sending new input to the agent
	// and begin shutdown procedures.
	//
	// Returns when the executor has fully stopped or the context is canceled.
	Stop(ctx context.Context) error
}
