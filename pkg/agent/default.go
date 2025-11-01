package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/entrhq/forge/pkg/agent/core"
	"github.com/entrhq/forge/pkg/agent/memory"
	"github.com/entrhq/forge/pkg/agent/prompts"
	"github.com/entrhq/forge/pkg/agent/tools"
	"github.com/entrhq/forge/pkg/llm"
	"github.com/entrhq/forge/pkg/types"
)

// DefaultAgent is a basic implementation of the Agent interface.
// It processes user inputs through an LLM provider using an agent loop
// with tools, thinking, and memory management.
type DefaultAgent struct {
	provider           llm.Provider
	channels           *types.AgentChannels
	customInstructions string
	maxTurns           int
	bufferSize         int
	metadata           map[string]interface{}

	// Agent loop components
	tools   map[string]tools.Tool
	toolsMu sync.RWMutex
	memory  memory.Memory

	// Control channels
	cancelMu     sync.Mutex
	cancelStream context.CancelFunc

	// Running state
	running bool
	runMu   sync.Mutex
}

// AgentOption is a function that configures an agent
type AgentOption func(*DefaultAgent)

// WithCustomInstructions sets custom instructions for the agent
// These are user-provided instructions that will be added to the system prompt
func WithCustomInstructions(instructions string) AgentOption {
	return func(a *DefaultAgent) {
		a.customInstructions = instructions
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
		tools:      make(map[string]tools.Tool),
		memory:     memory.NewConversationMemory(),
	}

	// Register built-in tools
	a.RegisterDefaultTools()

	// Apply options
	for _, opt := range opts {
		opt(a)
	}

	// Create channels with configured buffer size
	a.channels = types.NewAgentChannels(a.bufferSize)

	return a
}

func (a *DefaultAgent) RegisterDefaultTools() {
	// Initialize built-in tools (always available)
	a.tools["task_completion"] = tools.NewTaskCompletionTool()
	a.tools["ask_question"] = tools.NewAskQuestionTool()
	a.tools["converse"] = tools.NewConverseTool()
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

// processUserInput processes a user text input using the agent loop.
func (a *DefaultAgent) processUserInput(ctx context.Context, content string) {
	// Add user message to memory
	userMsg := types.NewUserMessage(content)
	a.memory.Add(userMsg)

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

	// Run agent loop
	a.runAgentLoop(turnCtx)

	// Emit turn end
	a.emitEvent(types.NewTurnEndEvent())
}

// runAgentLoop executes the agent loop with tools and thinking
// The loop continues until a loop-breaking tool is used
func (a *DefaultAgent) runAgentLoop(ctx context.Context) {
	for {
		// Build system prompt with tool schemas
		systemPrompt := a.buildSystemPrompt()

		// Get conversation history
		history := a.memory.GetAll()

		// Build messages for this iteration
		messages := prompts.BuildMessages(systemPrompt, history, "")

		// Get response from LLM
		stream, err := a.provider.StreamCompletion(ctx, messages)
		if err != nil {
			a.emitEvent(types.NewErrorEvent(fmt.Errorf("failed to start completion: %w", err)))
			return
		}

		// Process stream and collect response
		var assistantContent string
		var toolCallContent string
		core.ProcessStream(stream, a.emitEvent, func(content, thinking, toolCall, role string) {
			assistantContent = content
			toolCallContent = toolCall
			// thinking and tool call content are already emitted via events in ProcessStream
		})

		// Check for tool call
		if toolCallContent == "" {
			// No tool call found - emit appropriate event and error
			a.emitEvent(types.NewNoToolCallEvent())
			a.emitEvent(types.NewErrorEvent(fmt.Errorf("no tool call found in response")))
			return
		}

		// Parse the tool call JSON from the extracted content
		// toolCallContent already has just the JSON (without <tool> tags)
		// Parse directly as JSON instead of wrapping with XML tags
		var toolCall tools.ToolCall
		if err := json.Unmarshal([]byte(toolCallContent), &toolCall); err != nil {
			// Log the actual content to help debug
			a.emitEvent(types.NewErrorEvent(fmt.Errorf("failed to parse tool call JSON (len=%d, content: %q): %w", len(toolCallContent), toolCallContent, err)))
			return
		}

		// Validate required fields
		if toolCall.ToolName == "" {
			a.emitEvent(types.NewErrorEvent(fmt.Errorf("tool_name is required in tool call")))
			return
		}

		// Server name defaults to "local" if not specified
		if toolCall.ServerName == "" {
			toolCall.ServerName = "local"
		}

		// Add assistant's full response to memory (content + tool call)
		fullResponse := assistantContent
		if toolCallContent != "" {
			fullResponse += "<tool>" + toolCallContent + "</tool>"
		}
		a.memory.Add(&types.Message{
			Role:    types.RoleAssistant,
			Content: fullResponse,
		})

		// Execute the tool
		tool, exists := a.getTool(toolCall.ToolName)
		if !exists {
			a.emitEvent(types.NewErrorEvent(fmt.Errorf("unknown tool: %s", toolCall.ToolName)))
			return
		}

		// Emit tool call event
		// Convert RawMessage to map for event
		var argsMap map[string]interface{}
		if err := json.Unmarshal(toolCall.Arguments, &argsMap); err != nil {
			argsMap = make(map[string]interface{})
		}
		a.emitEvent(types.NewToolCallEvent(toolCall.ToolName, argsMap))

		// Execute tool
		result, toolErr := tool.Execute(ctx, toolCall.Arguments)
		if toolErr != nil {
			a.emitEvent(types.NewToolResultErrorEvent(toolCall.ToolName, toolErr))
			a.emitEvent(types.NewErrorEvent(fmt.Errorf("tool execution failed: %w", toolErr)))
			return
		}
		a.emitEvent(types.NewToolResultEvent(toolCall.ToolName, result))

		// Check if this is a loop-breaking tool
		if tool.IsLoopBreaking() {
			return
		}

		// For non-breaking tools, add result to memory and continue loop
		a.memory.Add(types.NewUserMessage(fmt.Sprintf("Tool '%s' result:\n%s", toolCall.ToolName, result)))
	}
}

// buildSystemPrompt constructs the system prompt with tool schemas and custom instructions
func (a *DefaultAgent) buildSystemPrompt() string {
	builder := prompts.NewPromptBuilder().
		WithTools(a.GetTools())

	// Add user's custom instructions if provided
	if a.customInstructions != "" {
		builder.WithCustomInstructions(a.customInstructions)
	}

	return builder.Build()
}

// emitEvent sends an event on the event channel.
// This is a blocking send to ensure critical events like TurnEnd are not dropped.
func (a *DefaultAgent) emitEvent(event *types.AgentEvent) {
	a.channels.Event <- event
}

// RegisterTool adds a custom tool to the agent's tool registry.
// Built-in tools (task_completion, ask_question, converse) are always available
// and cannot be overridden.
func (a *DefaultAgent) RegisterTool(tool tools.Tool) error {
	if tool == nil {
		return fmt.Errorf("tool cannot be nil")
	}

	name := tool.Name()
	if name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}

	// Prevent overriding built-in tools
	builtIns := map[string]bool{
		"task_completion": true,
		"ask_question":    true,
		"converse":        true,
	}
	if builtIns[name] {
		return fmt.Errorf("cannot override built-in tool: %s", name)
	}

	a.toolsMu.Lock()
	defer a.toolsMu.Unlock()

	a.tools[name] = tool
	return nil
}

// GetTools returns a list of all available tools (built-in + custom)
func (a *DefaultAgent) GetTools() []tools.Tool {
	a.toolsMu.RLock()
	defer a.toolsMu.RUnlock()

	toolsList := make([]tools.Tool, 0, len(a.tools))
	for _, tool := range a.tools {
		toolsList = append(toolsList, tool)
	}
	return toolsList
}

// getTool retrieves a tool by name (thread-safe)
func (a *DefaultAgent) getTool(name string) (tools.Tool, bool) {
	a.toolsMu.RLock()
	defer a.toolsMu.RUnlock()

	tool, exists := a.tools[name]
	return tool, exists
}
