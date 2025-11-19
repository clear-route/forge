package agent

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/entrhq/forge/pkg/agent/approval"
	agentcontext "github.com/entrhq/forge/pkg/agent/context"
	"github.com/entrhq/forge/pkg/agent/core"
	"github.com/entrhq/forge/pkg/agent/memory"
	"github.com/entrhq/forge/pkg/agent/prompts"
	"github.com/entrhq/forge/pkg/agent/tools"
	"github.com/entrhq/forge/pkg/llm"
	"github.com/entrhq/forge/pkg/llm/tokenizer"
	"github.com/entrhq/forge/pkg/tools/coding"
	"github.com/entrhq/forge/pkg/types"
)

var agentDebugLog *log.Logger

func init() {
	// Create debug log file in /tmp
	f, err := os.OpenFile("/tmp/forge-agent-debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Failed to open agent debug log: %v", err)
		agentDebugLog = log.New(os.Stderr, "[AGENT-DEBUG] ", log.LstdFlags|log.Lshortfile)
	} else {
		agentDebugLog = log.New(f, "[AGENT-DEBUG] ", log.LstdFlags|log.Lshortfile)
	}
}

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

	// Approval system
	approvalManager *approval.Manager
	approvalTimeout time.Duration

	// Control channels
	cancelMu     sync.Mutex
	cancelStream context.CancelFunc

	// Command execution tracking
	activeCommands sync.Map // executionID -> context.CancelFunc

	// Running state
	running bool
	runMu   sync.Mutex

	// Error recovery state
	lastErrors [5]string // Ring buffer of last 5 error messages
	errorIndex int       // Current position in ring buffer

	// Token usage tracking
	tokenizer *tokenizer.Tokenizer

	// Context management
	contextManager *agentcontext.Manager
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

// WithApprovalTimeout sets the timeout for approval requests
func WithApprovalTimeout(timeout time.Duration) AgentOption {
	return func(a *DefaultAgent) {
		// Store timeout for later use when creating approval manager
		if a.approvalManager != nil {
			// Manager already exists, we need to recreate it
			a.approvalManager = approval.NewManager(timeout, a.emitEvent)
		}
	}
}

// WithContextManager sets a context manager for the agent to handle context summarization
func WithContextManager(manager *agentcontext.Manager) AgentOption {
	return func(a *DefaultAgent) {
		a.contextManager = manager
	}
}

// NewDefaultAgent creates a new DefaultAgent with the given provider and options.
func NewDefaultAgent(provider llm.Provider, opts ...AgentOption) *DefaultAgent {
	// Create tokenizer for client-side token counting
	tok, err := tokenizer.New()
	if err != nil {
		// Fall back to nil tokenizer if initialization fails
		tok = nil
	}

	a := &DefaultAgent{
		provider:   provider,
		bufferSize: 10, // default buffer size
		tools:      make(map[string]tools.Tool),
		memory:     memory.NewConversationMemory(),
		tokenizer:  tok,
	}

	// Register built-in tools
	a.RegisterDefaultTools()

	// Apply options
	for _, opt := range opts {
		opt(a)
	}

	// Create channels with configured buffer size
	a.channels = types.NewAgentChannels(a.bufferSize)

	// Initialize approval manager with default timeout
	a.approvalTimeout = 5 * time.Minute
	a.approvalManager = approval.NewManager(a.approvalTimeout, a.emitEvent)

	// If context manager was provided, set its event channel now that channels exist
	if a.contextManager != nil {
		a.contextManager.SetEventChannel(a.channels.Event)
	}

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

	// Start a separate goroutine to handle cancellation requests
	// This ensures cancellations are processed even when the main loop is blocked
	cancelCtx, cancelStop := context.WithCancel(ctx)
	defer cancelStop()

	go func() {
		for {
			select {
			case <-cancelCtx.Done():
				return
			case cancelReq := <-a.channels.Cancel:
				if cancelReq == nil {
					return
				}
				// Handle command cancellation request immediately
				a.handleCommandCancellation(cancelReq)
			}
		}
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

			// Handle cancellation immediately (synchronously) so it can interrupt ongoing processing
			if input.IsCancel() {
				a.processInput(ctx, input)
				continue
			}

			// Process other inputs asynchronously so eventLoop can continue handling cancel requests
			go a.processInput(ctx, input)

		case approval := <-a.channels.Approval:
			if approval == nil {
				// Channel closed
				return
			}

			// Handle approval response
			a.handleApprovalResponse(approval)
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
// The loop continues until a loop-breaking tool is used or circuit breaker triggers
func (a *DefaultAgent) runAgentLoop(ctx context.Context) {
	var errorContext string

	for {
		// Check if context was canceled (e.g., via /stop command)
		select {
		case <-ctx.Done():
			// Context canceled - stop the agent loop
			// Emit a user-friendly message about the cancellation
			a.memory.Add(types.NewUserMessage("Operation stopped by user."))
			return
		default:
			// Continue with iteration
		}

		// Execute one iteration with optional error context from previous iteration
		shouldContinue, nextErrorContext := a.executeIteration(ctx, errorContext)
		if !shouldContinue {
			// Loop-breaking tool was called or terminal error occurred
			return
		}

		// Update error context for next iteration
		errorContext = nextErrorContext
	}
}

// executeIteration performs a single iteration of the agent loop
// Returns (shouldContinue, errorContext) where:
//   - shouldContinue: false if loop should terminate (loop-breaking tool or terminal error)
//   - errorContext: error message to pass to next iteration, or empty string if successful
func (a *DefaultAgent) executeIteration(ctx context.Context, errorContext string) (bool, string) {
	// Build system prompt with tool schemas
	systemPrompt := a.buildSystemPrompt()

	// Get conversation history
	history := a.memory.GetAll()

	// Build messages for this iteration with optional error context
	messages := prompts.BuildMessages(systemPrompt, history, "", errorContext)

	// Count input tokens if tokenizer is available
	var promptTokens int
	if a.tokenizer != nil {
		promptTokens = a.tokenizer.CountMessagesTokens(messages)
	}

	// Evaluate context summarization strategies before each turn
	agentDebugLog.Printf("=== Agent Loop Iteration ===")
	agentDebugLog.Printf("contextManager != nil: %v", a.contextManager != nil)
	agentDebugLog.Printf("tokenizer != nil: %v", a.tokenizer != nil)
	agentDebugLog.Printf("Prompt tokens: %d", promptTokens)

	if a.contextManager != nil && a.tokenizer != nil {
		agentDebugLog.Printf("Checking if memory is ConversationMemory...")
		// Get the conversation memory (cast from interface)
		if convMem, ok := a.memory.(*memory.ConversationMemory); ok {
			agentDebugLog.Printf("Memory is ConversationMemory - calling EvaluateAndSummarize")
			// Evaluate and run summarization strategies if needed
			_, err := a.contextManager.EvaluateAndSummarize(ctx, convMem, promptTokens)
			if err != nil {
				agentDebugLog.Printf("EvaluateAndSummarize returned error: %v", err)
				// Log error but continue - summarization failure shouldn't stop the agent
				a.emitEvent(types.NewErrorEvent(fmt.Errorf("context summarization failed: %w", err)))
			} else {
				agentDebugLog.Printf("EvaluateAndSummarize completed successfully")
			}

			// Rebuild messages after potential summarization
			history = a.memory.GetAll()
			messages = prompts.BuildMessages(systemPrompt, history, "", errorContext)

			// Recalculate tokens with updated messages
			if a.tokenizer != nil {
				promptTokens = a.tokenizer.CountMessagesTokens(messages)
				agentDebugLog.Printf("Tokens after potential summarization: %d", promptTokens)
			}
		} else {
			agentDebugLog.Printf("Memory is NOT ConversationMemory - type: %T", a.memory)
		}
	}

	// Emit API call start event with context information
	maxTokens := 0
	if a.contextManager != nil {
		maxTokens = a.contextManager.GetMaxTokens()
	}
	a.emitEvent(types.NewApiCallStartEvent("llm", promptTokens, maxTokens))

	// Get response from LLM
	stream, err := a.provider.StreamCompletion(ctx, messages)
	if err != nil {
		// Check if this is a context cancellation (user stopped the agent)
		if ctx.Err() != nil {
			return false, "" // Stop silently - user requested cancellation
		}
		// Terminal error - LLM/API failures should stop the loop
		a.emitEvent(types.NewErrorEvent(fmt.Errorf("failed to start completion: %w", err)))
		return false, ""
	}

	// Process stream and collect response
	var assistantContent string
	var toolCallContent string
	core.ProcessStream(stream, a.emitEvent, func(content, thinking, toolCall, role string) {
		assistantContent = content
		toolCallContent = toolCall
	})

	// Count completion tokens if tokenizer is available
	var completionTokens int
	if a.tokenizer != nil {
		fullResponse := assistantContent
		if toolCallContent != "" {
			fullResponse += toolCallContent
		}
		completionTokens = a.tokenizer.CountTokens(fullResponse)
	}

	// Emit token usage event if we have token counts
	if promptTokens > 0 || completionTokens > 0 {
		totalTokens := promptTokens + completionTokens
		a.emitEvent(types.NewTokenUsageEvent(promptTokens, completionTokens, totalTokens))
	}

	// Add assistant's response to memory
	fullResponse := assistantContent
	if toolCallContent != "" {
		fullResponse += "<tool>" + toolCallContent + "</tool>"
	}
	a.memory.Add(&types.Message{
		Role:    types.RoleAssistant,
		Content: fullResponse,
	})

	// Process the tool call (parse, validate, execute)
	return a.processToolCall(ctx, toolCallContent)
}

// buildSystemPrompt constructs the system prompt with tool schemas and custom instructions
func (a *DefaultAgent) buildSystemPrompt() string {
	builder := prompts.NewPromptBuilder().
		WithTools(a.getToolsList())

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

// GetTool retrieves a specific tool by name from the agent's tool registry.
// Returns nil if the tool is not found.
func (a *DefaultAgent) GetTool(name string) interface{} {
	a.toolsMu.RLock()
	defer a.toolsMu.RUnlock()

	return a.tools[name]
}

// GetTools returns a list of all available tools (built-in + custom)
// This is used internally for prompt building and memory
func (a *DefaultAgent) GetTools() []interface{} {
	a.toolsMu.RLock()
	defer a.toolsMu.RUnlock()

	toolsList := make([]interface{}, 0, len(a.tools))
	for _, tool := range a.tools {
		toolsList = append(toolsList, tool)
	}
	return toolsList
}

// GetContextInfo returns detailed context information for debugging and display
func (a *DefaultAgent) GetContextInfo() *ContextInfo {
	a.toolsMu.RLock()
	defer a.toolsMu.RUnlock()

	// Build system prompt without tools to calculate base system tokens
	baseSystemPrompt := prompts.NewPromptBuilder().
		WithCustomInstructions(a.customInstructions).
		Build()

	// Build just the tools section to calculate tool tokens
	toolsSection := ""
	if len(a.tools) > 0 {
		toolsSection = "<available_tools>\n" +
			prompts.FormatToolSchemas(a.getToolsList()) +
			"</available_tools>\n\n"
	}

	// Calculate token counts for each section
	systemPromptTokens := 0
	toolTokens := 0
	if a.tokenizer != nil {
		systemPromptTokens = a.tokenizer.CountTokens(baseSystemPrompt)
		toolTokens = a.tokenizer.CountTokens(toolsSection)
	}

	// Build full system prompt for current context calculation
	fullSystemPrompt := prompts.NewPromptBuilder().
		WithTools(a.getToolsList()).
		WithCustomInstructions(a.customInstructions).
		Build()

	// Get tool names
	toolNames := make([]string, 0, len(a.tools))
	for name := range a.tools {
		toolNames = append(toolNames, name)
	}

	// Get message history stats
	messages := a.memory.GetAll()
	messageCount := len(messages)

	// Count conversation turns (user messages)
	conversationTurns := 0
	for _, msg := range messages {
		if msg.Role == types.RoleUser {
			conversationTurns++
		}
	}

	// Calculate token counts
	conversationTokens := 0
	currentTokens := 0
	if a.tokenizer != nil {
		conversationTokens = a.tokenizer.CountMessagesTokens(messages)
		// Calculate current context tokens
		currentTokens = conversationTokens + a.tokenizer.CountTokens(fullSystemPrompt)
	} else {
		// Fallback: approximate token counting when tokenizer is unavailable
		// Use rough estimate of 1 token ≈ 4 characters
		for _, msg := range messages {
			conversationTokens += (len(msg.Content) + len(string(msg.Role)) + 12) / 4 // +12 for message overhead
		}
		currentTokens = conversationTokens + len(fullSystemPrompt)/4
	}

	// Get max tokens from context manager
	maxTokens := 0
	if a.contextManager != nil {
		maxTokens = a.contextManager.GetMaxTokens()
	}

	// Calculate free tokens and usage percentage
	freeTokens := 0
	usagePercent := 0.0
	if maxTokens > 0 {
		freeTokens = maxTokens - currentTokens
		if freeTokens < 0 {
			freeTokens = 0
		}
		usagePercent = float64(currentTokens) / float64(maxTokens) * 100.0
	}

	return &ContextInfo{
		SystemPromptTokens:    systemPromptTokens,
		CustomInstructions:    a.customInstructions != "",
		ToolCount:             len(a.tools),
		ToolTokens:            toolTokens,
		ToolNames:             toolNames,
		MessageCount:          messageCount,
		ConversationTurns:     conversationTurns,
		ConversationTokens:    conversationTokens,
		CurrentContextTokens:  currentTokens,
		MaxContextTokens:      maxTokens,
		FreeTokens:            freeTokens,
		UsagePercent:          usagePercent,
		TotalPromptTokens:     0, // These will be filled by the executor from its tracking
		TotalCompletionTokens: 0,
		TotalTokens:           0,
	}
}

// getToolsList returns tools as []tools.Tool for internal use
func (a *DefaultAgent) getToolsList() []tools.Tool {
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

// trackError adds an error to the ring buffer and checks if we've hit the circuit breaker
// Returns true if the circuit breaker should trigger (5 identical consecutive errors)
func (a *DefaultAgent) trackError(errMsg string) bool {
	// Add to ring buffer
	a.lastErrors[a.errorIndex] = errMsg
	a.errorIndex = (a.errorIndex + 1) % 5

	// Check if all 5 are identical and non-empty
	if a.lastErrors[0] == "" {
		return false // Not enough errors yet
	}

	first := a.lastErrors[0]
	for i := 1; i < 5; i++ {
		if a.lastErrors[i] != first {
			return false
		}
	}

	return true // All 5 errors are identical
}

// resetErrorTracking clears the error ring buffer after a successful iteration
func (a *DefaultAgent) resetErrorTracking() {
	for i := range a.lastErrors {
		a.lastErrors[i] = ""
	}
	a.errorIndex = 0
}

// processToolCall handles parsing, validation, and execution of tool calls
// Returns (shouldContinue, errorContext) following the same pattern as executeIteration
func (a *DefaultAgent) processToolCall(ctx context.Context, toolCallContent string) (bool, string) {
	// Check if context was canceled before processing
	if ctx.Err() != nil {
		return false, "" // Stop silently - user requested cancellation
	}

	// Check if tool call exists
	if toolCallContent == "" {
		// If context was canceled, this is expected (stream was interrupted)
		if ctx.Err() != nil {
			return false, ""
		}

		a.emitEvent(types.NewNoToolCallEvent())
		errMsg := prompts.BuildErrorRecoveryMessage(prompts.ErrorRecoveryContext{
			Type: prompts.ErrorTypeNoToolCall,
		})

		if a.trackError(errMsg) {
			a.emitEvent(types.NewErrorEvent(fmt.Errorf("circuit breaker triggered: 5 consecutive no tool call errors")))
			return false, ""
		}

		a.emitEvent(types.NewErrorEvent(fmt.Errorf("no tool call found in response")))
		return true, errMsg
	}

	// Parse the tool call (supports both XML and JSON formats)
	// Wrap content in <tool> tags since streaming parser strips them
	wrappedContent := "<tool>" + toolCallContent + "</tool>"
	parsedToolCall, _, err := tools.ParseToolCall(wrappedContent)
	if err != nil {
		// Log the actual content for debugging
		a.emitEvent(types.NewMessageContentEvent(fmt.Sprintf("\n🔍 DEBUG - Failed to parse tool call:\n%s\n", toolCallContent)))

		errMsg := prompts.BuildErrorRecoveryMessage(prompts.ErrorRecoveryContext{
			Type:    prompts.ErrorTypeInvalidXML,
			Error:   err,
			Content: toolCallContent,
		})

		if a.trackError(errMsg) {
			a.emitEvent(types.NewErrorEvent(fmt.Errorf("circuit breaker triggered: 5 consecutive parse errors")))
			return false, ""
		}

		a.emitEvent(types.NewErrorEvent(fmt.Errorf("failed to parse tool call: %w", err)))
		return true, errMsg
	}

	// Use the parsed tool call
	toolCall := *parsedToolCall

	// Validate required fields
	if toolCall.ToolName == "" {
		errMsg := prompts.BuildErrorRecoveryMessage(prompts.ErrorRecoveryContext{
			Type: prompts.ErrorTypeMissingToolName,
		})

		if a.trackError(errMsg) {
			a.emitEvent(types.NewErrorEvent(fmt.Errorf("circuit breaker triggered: 5 consecutive missing tool name errors")))
			return false, ""
		}

		a.emitEvent(types.NewErrorEvent(fmt.Errorf("tool_name is required in tool call")))
		return true, errMsg
	}

	// Server name defaults to "local" if not specified
	if toolCall.ServerName == "" {
		toolCall.ServerName = "local"
	}

	// Execute the tool
	return a.executeTool(ctx, toolCall)
}

// executeTool handles tool lookup, execution, and result processing
// Returns (shouldContinue, errorContext) following the same pattern as executeIteration
func (a *DefaultAgent) executeTool(ctx context.Context, toolCall tools.ToolCall) (bool, string) {
	// Look up the tool
	tool, exists := a.getTool(toolCall.ToolName)
	if !exists {
		errMsg := prompts.BuildErrorRecoveryMessage(prompts.ErrorRecoveryContext{
			Type:           prompts.ErrorTypeUnknownTool,
			ToolName:       toolCall.ToolName,
			AvailableTools: a.getToolsList(),
		})

		// Track error and check circuit breaker
		if a.trackError(errMsg) {
			a.emitEvent(types.NewErrorEvent(fmt.Errorf("circuit breaker triggered: 5 consecutive unknown tool errors")))
			return false, ""
		}

		a.emitEvent(types.NewErrorEvent(fmt.Errorf("unknown tool: %s", toolCall.ToolName)))
		return true, errMsg // Continue with error context
	}

	// Check if tool requires approval
	if previewable, ok := tool.(tools.Previewable); ok {
		// Generate preview
		preview, err := previewable.GeneratePreview(ctx, toolCall.GetArgumentsXML())
		if err != nil {
			// If preview generation fails, log error but continue with execution
			// (degraded mode - execute without approval)
			a.emitEvent(types.NewErrorEvent(fmt.Errorf("failed to generate preview for %s: %w", toolCall.ToolName, err)))
		} else {
			// Request approval from user
			approved, timedOut := a.requestApproval(ctx, toolCall, preview)

			if timedOut {
				// Timeout - treat as rejection and continue loop
				errMsg := fmt.Sprintf("Tool approval request timed out after %v. The tool was not executed.", a.approvalTimeout)
				a.memory.Add(types.NewUserMessage(errMsg))
				return true, ""
			}

			if !approved {
				// User rejected - continue loop without executing
				errMsg := fmt.Sprintf("Tool '%s' execution was rejected by user.", toolCall.ToolName)
				a.memory.Add(types.NewUserMessage(errMsg))
				return true, ""
			}

			// User approved - continue with execution
		}
	}

	// Emit tool call event
	var argsMap map[string]interface{}
	if err := tools.UnmarshalXMLWithFallback(toolCall.GetArgumentsXML(), &argsMap); err != nil {
		argsMap = make(map[string]interface{})
	}
	a.emitEvent(types.NewToolCallEvent(toolCall.ToolName, argsMap))

	// Inject event emitter and command registry into context for tools that support streaming events
	ctxWithEmitter := context.WithValue(ctx, coding.EventEmitterKey, coding.EventEmitter(a.emitEvent))
	ctxWithRegistry := context.WithValue(ctxWithEmitter, coding.CommandRegistryKey, &a.activeCommands)

	// Execute the tool
	result, toolErr := tool.Execute(ctxWithRegistry, toolCall.GetArgumentsXML())
	if toolErr != nil {
		a.emitEvent(types.NewToolResultErrorEvent(toolCall.ToolName, toolErr))
		errMsg := prompts.BuildErrorRecoveryMessage(prompts.ErrorRecoveryContext{
			Type:     prompts.ErrorTypeToolExecution,
			ToolName: toolCall.ToolName,
			Error:    toolErr,
		})

		// Track error and check circuit breaker
		if a.trackError(errMsg) {
			a.emitEvent(types.NewErrorEvent(fmt.Errorf("circuit breaker triggered: 5 consecutive tool execution errors")))
			return false, ""
		}

		a.emitEvent(types.NewErrorEvent(fmt.Errorf("tool execution failed: %w", toolErr)))
		return true, errMsg // Continue with error context
	}

	a.emitEvent(types.NewToolResultEvent(toolCall.ToolName, result))

	// Success! Reset error tracking
	a.resetErrorTracking()

	// Check if this is a loop-breaking tool
	if tool.IsLoopBreaking() {
		return false, "" // Stop loop
	}

	// For non-breaking tools, add result to memory and continue loop
	a.memory.Add(types.NewUserMessage(fmt.Sprintf("Tool '%s' result:\n%s", toolCall.ToolName, result)))
	return true, "" // Continue with no error
}

// handleApprovalResponse processes an approval response from the user
func (a *DefaultAgent) handleApprovalResponse(response *types.ApprovalResponse) {
	a.approvalManager.HandleResponse(response)
}

// handleCommandCancellation processes a command cancellation request
func (a *DefaultAgent) handleCommandCancellation(req *types.CancellationRequest) {
	// Look up the cancel function for this execution ID
	if cancelFunc, ok := a.activeCommands.Load(req.ExecutionID); ok {
		// Cancel the context (cancellation never returns an error)
		if cf, ok := cancelFunc.(context.CancelFunc); ok {
			cf()
		}
		// Remove from active commands
		a.activeCommands.Delete(req.ExecutionID)
	}
}

// requestApproval sends an approval request and waits for user response
// Returns (approved, timedOut) where:
//   - approved: true if user approved, false if rejected
//   - timedOut: true if the request timed out waiting for response
func (a *DefaultAgent) requestApproval(ctx context.Context, toolCall tools.ToolCall, preview *tools.ToolPreview) (bool, bool) {
	// Delegate all approval logic to the approval manager
	return a.approvalManager.RequestApproval(ctx, toolCall, preview)
}
