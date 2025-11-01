package tools

import (
	"context"
	"encoding/json"
)

// Tool represents a capability that an agent can use during execution.
// Tools are invoked by the LLM through XML-formatted tool calls and can
// perform actions like task completion, asking questions, or custom operations.
//
// Example tool call format from LLM:
//
//	<tool>{"server_name": "local", "tool_name": "task_completion", "arguments": {"result": "Task done!"}}</tool>
type Tool interface {
	// Name returns the unique identifier for this tool (e.g., "task_completion")
	Name() string

	// Description returns a human-readable description of what this tool does
	Description() string

	// Schema returns the JSON schema for this tool's input parameters
	// The schema should be a valid JSON Schema object defining the structure
	// of the arguments that this tool accepts
	Schema() map[string]interface{}

	// Execute runs the tool with the given arguments and returns a result string
	// The arguments are validated against the schema before execution
	Execute(ctx context.Context, arguments json.RawMessage) (string, error)

	// IsLoopBreaking indicates whether this tool should terminate the agent loop
	// Loop-breaking tools (like task_completion, ask_question, converse) will
	// cause the agent to stop iterating and return control to the user
	IsLoopBreaking() bool
}

// ToolCall represents a parsed tool invocation from the LLM's response
type ToolCall struct {
	ServerName string          `json:"server_name"`
	ToolName   string          `json:"tool_name"`
	Arguments  json.RawMessage `json:"arguments"`
}

// BaseToolSchema creates a common JSON schema structure for a tool
// with the given properties and required fields
func BaseToolSchema(properties map[string]interface{}, required []string) map[string]interface{} {
	schema := map[string]interface{}{
		"type":       "object",
		"properties": properties,
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}
