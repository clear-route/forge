package types

// AgentEventType defines the type of event emitted by the agent.
type AgentEventType string

const (
	EventTypeThinkingStart   AgentEventType = "thinking_start"    // EventTypeThinkingStart indicates the agent is starting to think/reason.
	EventTypeThinkingContent AgentEventType = "thinking_content"  // EventTypeThinkingContent indicates content from the agent's thinking process.
	EventTypeThinkingEnd     AgentEventType = "thinking_end"      // EventTypeThinkingEnd indicates the agent has finished thinking.
	EventTypeMessageStart    AgentEventType = "message_start"     // EventTypeMessageStart indicates the agent is starting to compose a message.
	EventTypeMessageContent  AgentEventType = "message_content"   // EventTypeMessageContent indicates content from the agent's message.
	EventTypeMessageEnd      AgentEventType = "message_end"       // EventTypeMessageEnd indicates the agent has finished composing the message.
	EventTypeToolCall        AgentEventType = "tool_call"         // EventTypeToolCall indicates the agent is calling a tool.
	EventTypeToolResult      AgentEventType = "tool_result"       // EventTypeToolResult indicates a successful tool call result.
	EventTypeToolResultError AgentEventType = "tool_result_error" // EventTypeToolResultError indicates a tool call resulted in an error.
	EventTypeNoToolCall      AgentEventType = "no_tool_call"      // EventTypeNoToolCall indicates the agent decided not to call any tools.
	EventTypeApiCallStart    AgentEventType = "api_call_start"    // EventTypeApiCallStart indicates the agent is making an API call.
	EventTypeApiCallEnd      AgentEventType = "api_call_end"      // EventTypeApiCallEnd indicates an API call has completed.
	EventTypeToolsUpdate     AgentEventType = "tools_update"      // EventTypeToolsUpdate indicates the agent's available tools have been updated.
	EventTypeUpdateBusy      AgentEventType = "update_busy"       // EventTypeUpdateBusy indicates a change in the agent's busy status.
	EventTypeTurnEnd         AgentEventType = "turn_end"          // EventTypeTurnEnd indicates the agent has finished processing the current turn.
	EventTypeError           AgentEventType = "error"             // EventTypeError indicates an error occurred during agent processing.
)

// AgentEvent represents an event emitted by the agent during execution.
type AgentEvent struct {
	// Metadata holds optional additional information about the event.
	Metadata map[string]interface{}

	// ToolInput is the input being sent to the tool (for tool call events).
	ToolInput map[string]interface{}

	// ToolOutput is the result from the tool (for tool result events).
	ToolOutput interface{}

	// Error contains error information for error events.
	Error error

	// Content holds text content for content-type events (thinking, message, etc.).
	Content string

	// ToolName is the name of the tool being called (for tool events).
	ToolName string

	// Type indicates the kind of event.
	Type AgentEventType

	// IsBusy indicates if the agent is busy (for busy status events).
	IsBusy bool
}

// NewThinkingStartEvent creates a thinking start event.
func NewThinkingStartEvent() *AgentEvent {
	return &AgentEvent{
		Type:     EventTypeThinkingStart,
		Metadata: make(map[string]interface{}),
	}
}

// NewThinkingContentEvent creates a thinking content event.
func NewThinkingContentEvent(content string) *AgentEvent {
	return &AgentEvent{
		Type:     EventTypeThinkingContent,
		Content:  content,
		Metadata: make(map[string]interface{}),
	}
}

// NewThinkingEndEvent creates a thinking end event.
func NewThinkingEndEvent() *AgentEvent {
	return &AgentEvent{
		Type:     EventTypeThinkingEnd,
		Metadata: make(map[string]interface{}),
	}
}

// NewMessageStartEvent creates a message start event.
func NewMessageStartEvent() *AgentEvent {
	return &AgentEvent{
		Type:     EventTypeMessageStart,
		Metadata: make(map[string]interface{}),
	}
}

// NewMessageContentEvent creates a message content event.
func NewMessageContentEvent(content string) *AgentEvent {
	return &AgentEvent{
		Type:     EventTypeMessageContent,
		Content:  content,
		Metadata: make(map[string]interface{}),
	}
}

// NewMessageEndEvent creates a message end event.
func NewMessageEndEvent() *AgentEvent {
	return &AgentEvent{
		Type:     EventTypeMessageEnd,
		Metadata: make(map[string]interface{}),
	}
}

// NewToolCallEvent creates a tool call event.
func NewToolCallEvent(toolName string, toolInput map[string]interface{}) *AgentEvent {
	return &AgentEvent{
		Type:      EventTypeToolCall,
		ToolName:  toolName,
		ToolInput: toolInput,
		Metadata:  make(map[string]interface{}),
	}
}

// NewToolResultEvent creates a tool result event.
func NewToolResultEvent(toolName string, output interface{}) *AgentEvent {
	return &AgentEvent{
		Type:       EventTypeToolResult,
		ToolName:   toolName,
		ToolOutput: output,
		Metadata:   make(map[string]interface{}),
	}
}

// NewToolResultErrorEvent creates a tool result error event.
func NewToolResultErrorEvent(toolName string, err error) *AgentEvent {
	return &AgentEvent{
		Type:     EventTypeToolResultError,
		ToolName: toolName,
		Error:    err,
		Metadata: make(map[string]interface{}),
	}
}

// NewNoToolCallEvent creates a no tool call event.
func NewNoToolCallEvent() *AgentEvent {
	return &AgentEvent{
		Type:     EventTypeNoToolCall,
		Metadata: make(map[string]interface{}),
	}
}

// NewApiCallStartEvent creates an API call start event.
func NewApiCallStartEvent(apiName string) *AgentEvent {
	return &AgentEvent{
		Type:     EventTypeApiCallStart,
		Metadata: map[string]interface{}{"api_name": apiName},
	}
}

// NewApiCallEndEvent creates an API call end event.
func NewApiCallEndEvent(apiName string) *AgentEvent {
	return &AgentEvent{
		Type:     EventTypeApiCallEnd,
		Metadata: map[string]interface{}{"api_name": apiName},
	}
}

// NewToolsUpdateEvent creates a tools update event.
func NewToolsUpdateEvent(tools []string) *AgentEvent {
	return &AgentEvent{
		Type:     EventTypeToolsUpdate,
		Metadata: map[string]interface{}{"tools": tools},
	}
}

// NewUpdateBusyEvent creates a busy status update event.
func NewUpdateBusyEvent(isBusy bool) *AgentEvent {
	return &AgentEvent{
		Type:     EventTypeUpdateBusy,
		IsBusy:   isBusy,
		Metadata: make(map[string]interface{}),
	}
}

// NewTurnEndEvent creates a turn end event.
func NewTurnEndEvent() *AgentEvent {
	return &AgentEvent{
		Type:     EventTypeTurnEnd,
		Metadata: make(map[string]interface{}),
	}
}

// NewErrorEvent creates an error event.
func NewErrorEvent(err error) *AgentEvent {
	return &AgentEvent{
		Type:     EventTypeError,
		Error:    err,
		Metadata: make(map[string]interface{}),
	}
}

// WithMetadata adds metadata to the event and returns the event for chaining.
func (e *AgentEvent) WithMetadata(key string, value interface{}) *AgentEvent {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[key] = value
	return e
}

// IsThinkingEvent returns true if this is any thinking-related event.
func (e *AgentEvent) IsThinkingEvent() bool {
	return e.Type == EventTypeThinkingStart ||
		e.Type == EventTypeThinkingContent ||
		e.Type == EventTypeThinkingEnd
}

// IsMessageEvent returns true if this is any message-related event.
func (e *AgentEvent) IsMessageEvent() bool {
	return e.Type == EventTypeMessageStart ||
		e.Type == EventTypeMessageContent ||
		e.Type == EventTypeMessageEnd
}

// IsToolEvent returns true if this is any tool-related event.
func (e *AgentEvent) IsToolEvent() bool {
	return e.Type == EventTypeToolCall ||
		e.Type == EventTypeToolResult ||
		e.Type == EventTypeToolResultError ||
		e.Type == EventTypeNoToolCall
}

// IsApiEvent returns true if this is any API-related event.
func (e *AgentEvent) IsApiEvent() bool {
	return e.Type == EventTypeApiCallStart ||
		e.Type == EventTypeApiCallEnd
}

// IsContentEvent returns true if this event contains text content.
func (e *AgentEvent) IsContentEvent() bool {
	return e.Type == EventTypeThinkingContent ||
		e.Type == EventTypeMessageContent
}

// IsErrorEvent returns true if this is an error event.
func (e *AgentEvent) IsErrorEvent() bool {
	return e.Type == EventTypeError
}
