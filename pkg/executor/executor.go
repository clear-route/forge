// Package executor provides the execution plane abstraction
// allowing agents to run in different environments (CLI, API, etc.)
package executor

import (
	"context"

	"github.com/clear-route/forge/pkg/agent"
)

// Executor interface defines how agents interact with their runtime environment.
// Executors provide the I/O layer for agents, handling environment-specific
// concerns like CLI interaction, HTTP requests, or other communication protocols.
type Executor interface {
	// Run starts the executor with the given agent.
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
	Run(ctx context.Context, agent agent.Agent) error

	// Stop gracefully stops the executor.
	// This signals the executor to stop sending new input to the agent
	// and begin shutdown procedures.
	//
	// Returns when the executor has fully stopped or the context is canceled.
	Stop(ctx context.Context) error
}
