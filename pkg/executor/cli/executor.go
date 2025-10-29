// Package cli provides a command-line executor for Forge agents.
package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/entrhq/forge/pkg/agent"
	"github.com/entrhq/forge/pkg/types"
)

// Executor is a CLI-based executor that enables turn-by-turn conversation
// with an agent through terminal input/output.
type Executor struct {
	agent  agent.Agent
	reader *bufio.Reader
	writer io.Writer

	// Display options
	showThinking bool
	prompt       string
}

// ExecutorOption is a function that configures an Executor.
type ExecutorOption func(*Executor)

// WithShowThinking enables/disables displaying the agent's thinking process.
func WithShowThinking(show bool) ExecutorOption {
	return func(e *Executor) {
		e.showThinking = show
	}
}

// WithPrompt sets a custom prompt string.
func WithPrompt(prompt string) ExecutorOption {
	return func(e *Executor) {
		e.prompt = prompt
	}
}

// WithWriter sets a custom output writer (default is os.Stdout).
func WithWriter(w io.Writer) ExecutorOption {
	return func(e *Executor) {
		e.writer = w
	}
}

// NewExecutor creates a new CLI executor for the given agent.
func NewExecutor(agent agent.Agent, opts ...ExecutorOption) *Executor {
	e := &Executor{
		agent:        agent,
		reader:       bufio.NewReader(os.Stdin),
		writer:       os.Stdout,
		showThinking: true, // Show thinking by default
		prompt:       "> ",
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

// Run starts the executor and begins the conversation loop.
// Returns when the user exits or an error occurs.
func (e *Executor) Run(ctx context.Context) error {
	// Start the agent
	if err := e.agent.Start(ctx); err != nil {
		return fmt.Errorf("failed to start agent: %w", err)
	}

	channels := e.agent.GetChannels()

	// Start event handler in background
	eventsDone := make(chan struct{})
	turnEnd := make(chan struct{}, 1)
	go e.handleEvents(channels.Event, eventsDone, turnEnd)

	// Print welcome message
	fmt.Fprintln(e.writer, "Forge CLI Agent")
	fmt.Fprintln(e.writer, "Type your message and press Enter. Type 'exit' or 'quit' to end the conversation.")
	fmt.Fprintln(e.writer)

	// Main conversation loop
	for {
		// Check if context is canceled
		select {
		case <-ctx.Done():
			e.shutdown(ctx)
			<-eventsDone
			return ctx.Err()
		default:
		}

		// Read user input
		fmt.Fprint(e.writer, e.prompt)
		input, err := e.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				e.shutdown(ctx)
				<-eventsDone
				return nil
			}
			return fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)

		// Handle exit commands
		if input == "exit" || input == "quit" {
			e.shutdown(ctx)
			<-eventsDone
			return nil
		}

		// Skip empty input
		if input == "" {
			continue
		}

		// Send input to agent
		channels.Input <- types.NewUserInput(input)

		// Wait for turn to complete
		<-turnEnd
	}
}

// handleEvents processes events from the agent and renders them to the terminal.
func (e *Executor) handleEvents(events <-chan *types.AgentEvent, done chan struct{}, turnEnd chan struct{}) {
	defer close(done)

	for event := range events {
		switch event.Type {
		case types.EventTypeThinkingStart:
			if e.showThinking {
				fmt.Fprintln(e.writer, "\n[Thinking...]")
			}

		case types.EventTypeThinkingContent:
			if e.showThinking {
				fmt.Fprint(e.writer, event.Content)
			}

		case types.EventTypeThinkingEnd:
			if e.showThinking {
				fmt.Fprintln(e.writer, "\n[Done thinking]")
			}

		case types.EventTypeMessageStart:
			// Start assistant response on new line
			fmt.Fprintln(e.writer, "Assistant:")

		case types.EventTypeMessageContent:
			fmt.Fprint(e.writer, event.Content)

		case types.EventTypeMessageEnd:
			fmt.Fprintln(e.writer) // New line after message

		case types.EventTypeError:
			fmt.Fprintf(e.writer, "\n❌ Error: %v\n", event.Error)

		case types.EventTypeUpdateBusy:
			// Could show a spinner here in the future

		case types.EventTypeTurnEnd:
			// Signal that turn is complete
			select {
			case turnEnd <- struct{}{}:
			default:
			}
		}
	}
}

// shutdown gracefully shuts down the agent.
func (e *Executor) shutdown(ctx context.Context) {
	fmt.Fprintln(e.writer, "\nShutting down...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*1000000000) // 5 seconds
	defer cancel()

	if err := e.agent.Shutdown(shutdownCtx); err != nil {
		fmt.Fprintf(e.writer, "Warning: shutdown error: %v\n", err)
	}
}
