package types

// AgentChannels defines the communication channels between an agent and its executor.
// These channels enable async, event-driven communication.
type AgentChannels struct {
	// Input is the channel for sending inputs to the agent.
	// The executor writes user inputs (text, forms, or cancellations) here, and the agent reads from it.
	Input chan *Input

	// Event is the channel for receiving all events from the agent.
	// The agent emits various events (thinking, messages, tool calls, errors, etc.) on this channel.
	Event chan *AgentEvent

	// Approval is the channel for receiving approval responses from the executor.
	// When the agent requests approval, the executor sends the user's decision here.
	Approval chan *ApprovalResponse

	// Shutdown is the channel for signaling the agent to shut down.
	// The executor closes this channel to initiate graceful shutdown.
	Shutdown chan struct{}

	// Done is the channel for signaling that the agent has completed shutdown.
	// The agent closes this channel when it has fully stopped.
	Done chan struct{}
}

// NewAgentChannels creates a new AgentChannels instance with the specified buffer size.
// All channels are buffered to prevent blocking.
func NewAgentChannels(bufferSize int) *AgentChannels {
	return &AgentChannels{
		Input:    make(chan *Input, bufferSize),
		Event:    make(chan *AgentEvent, bufferSize),
		Approval: make(chan *ApprovalResponse, bufferSize),
		Shutdown: make(chan struct{}),
		Done:     make(chan struct{}),
	}
}

// Close closes all channels. This should only be called by the agent during shutdown
// to prevent send on closed channel panics.
func (c *AgentChannels) Close() {
	close(c.Event)
	close(c.Approval)
	close(c.Done)
}
