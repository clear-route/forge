package agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/entrhq/forge/pkg/agent/core"
	"github.com/entrhq/forge/pkg/llm"
	"github.com/entrhq/forge/pkg/types"
)

// DefaultAgent is a basic implementation of the Agent interface.
// It processes user inputs through an LLM provider and emits events
// for thinking, messages, and errors.
type DefaultAgent struct {
	provider     llm.Provider
	channels     *types.AgentChannels
	systemPrompt string
	maxTurns     int
	bufferSize   int
	metadata     map[string]interface{}

	// Conversation history with thread-safe access
	historyMu sync.RWMutex
	history   []*types.Message

	// Control channels
	cancelMu     sync.Mutex
	cancelStream context.CancelFunc

	// Running state
	running bool
	runMu   sync.Mutex
}

// AgentOption is a function that configures an agent
type AgentOption func(*DefaultAgent)

// WithSystemPrompt sets the system prompt for the agent
func WithSystemPrompt(prompt string) AgentOption {
	return func(a *DefaultAgent) {
		a.systemPrompt = prompt
	}
}

// WithMaxTurns sets the maximum number of conversation turns
func WithMaxTurns(max int) AgentOption {
	return func(a *DefaultAgent) {
		a.maxTurns = max
	}
}

// WithBufferSize sets the channel buffer size
func WithBufferSize(size int) AgentOption {
	return func(a *DefaultAgent) {
		a.bufferSize = size
	}
}

// WithMetadata sets metadata for the agent
func WithMetadata(metadata map[string]interface{}) AgentOption {
	return func(a *DefaultAgent) {
		a.metadata = metadata
	}
}

// NewDefaultAgent creates a new DefaultAgent with the given provider and options.
func NewDefaultAgent(provider llm.Provider, opts ...AgentOption) *DefaultAgent {
	a := &DefaultAgent{
		provider:   provider,
		bufferSize: 10, // default buffer size
		history:    make([]*types.Message, 0),
	}

	// Apply options
	for _, opt := range opts {
		opt(a)
	}

	// Create channels with configured buffer size
	a.channels = types.NewAgentChannels(a.bufferSize)

	return a
}

// Start begins the agent's event loop in a goroutine.
func (a *DefaultAgent) Start(ctx context.Context) error {
	a.runMu.Lock()
	if a.running {
		a.runMu.Unlock()
		return fmt.Errorf("agent is already running")
	}
	a.running = true
	a.runMu.Unlock()

	// Add system message if configured
	if a.systemPrompt != "" {
		a.addToHistory(types.NewSystemMessage(a.systemPrompt))
	}

	// Start event loop
	go a.eventLoop(ctx)

	return nil
}

// Shutdown gracefully stops the agent.
func (a *DefaultAgent) Shutdown(ctx context.Context) error {
	// Signal shutdown
	close(a.channels.Shutdown)

	// Wait for completion or context cancellation
	select {
	case <-a.channels.Done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// GetChannels returns the communication channels for this agent.
func (a *DefaultAgent) GetChannels() *types.AgentChannels {
	return a.channels
}

// eventLoop is the main processing loop for the agent.
func (a *DefaultAgent) eventLoop(ctx context.Context) {
	defer a.channels.Close()
	defer func() {
		a.runMu.Lock()
		a.running = false
		a.runMu.Unlock()
	}()

	for {
		select {
		case <-ctx.Done():
			// Context canceled
			a.emitEvent(types.NewErrorEvent(ctx.Err()))
			return

		case <-a.channels.Shutdown:
			// Shutdown requested
			return

		case input := <-a.channels.Input:
			if input == nil {
				// Channel closed
				return
			}

			// Process the input
			a.processInput(ctx, input)
		}
	}
}

// processInput handles a single input from the user.
func (a *DefaultAgent) processInput(ctx context.Context, input *types.Input) {
	// Handle cancellation
	if input.IsCancel() {
		a.cancelMu.Lock()
		if a.cancelStream != nil {
			a.cancelStream()
			a.cancelStream = nil
		}
		a.cancelMu.Unlock()
		return
	}

	// Handle user input
	if input.IsUserInput() {
		a.processUserInput(ctx, input.Content)
		return
	}

	// Handle form input (not yet implemented)
	if input.IsFormInput() {
		a.emitEvent(types.NewErrorEvent(fmt.Errorf("form input not yet supported")))
		a.emitEvent(types.NewTurnEndEvent())
		return
	}
}

// processUserInput processes a user text input.
func (a *DefaultAgent) processUserInput(ctx context.Context, content string) {
	// Add user message to history
	userMsg := types.NewUserMessage(content)
	a.addToHistory(userMsg)

	// Create cancellable context for this turn
	turnCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	a.cancelMu.Lock()
	a.cancelStream = cancel
	a.cancelMu.Unlock()

	defer func() {
		a.cancelMu.Lock()
		a.cancelStream = nil
		a.cancelMu.Unlock()
	}()

	// Emit busy status
	a.emitEvent(types.NewUpdateBusyEvent(true))
	defer a.emitEvent(types.NewUpdateBusyEvent(false))

	// Get response from LLM
	stream, err := a.provider.StreamCompletion(turnCtx, a.getHistory())
	if err != nil {
		a.emitEvent(types.NewErrorEvent(fmt.Errorf("failed to start completion: %w", err)))
		a.emitEvent(types.NewTurnEndEvent())
		return
	}

	// Process stream chunks (handles thinking and message events) and converts to agent events
	core.ProcessStream(stream, a.emitEvent, a.handleStreamComplete)

	// Emit turn end
	a.emitEvent(types.NewTurnEndEvent())
}

// handleStreamComplete is called when stream processing completes successfully.
// It adds the assistant's message to conversation history if there was content.
func (a *DefaultAgent) handleStreamComplete(assistantContent, thinkingContent, role string) {
	if assistantContent != "" {
		a.addToHistory(&types.Message{
			Role:    types.MessageRole(role),
			Content: assistantContent,
		})
	}
}

// emitEvent sends an event on the event channel.
// This is a blocking send to ensure critical events like TurnEnd are not dropped.
func (a *DefaultAgent) emitEvent(event *types.AgentEvent) {
	a.channels.Event <- event
}

// addToHistory adds a message to the conversation history (thread-safe).
func (a *DefaultAgent) addToHistory(msg *types.Message) {
	a.historyMu.Lock()
	defer a.historyMu.Unlock()
	a.history = append(a.history, msg)
}

// getHistory returns a copy of the conversation history (thread-safe).
func (a *DefaultAgent) getHistory() []*types.Message {
	a.historyMu.RLock()
	defer a.historyMu.RUnlock()

	// Return a copy to avoid race conditions
	history := make([]*types.Message, len(a.history))
	copy(history, a.history)
	return history
}
