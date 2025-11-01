package prompts

import (
	"fmt"
	"strings"

	"github.com/entrhq/forge/pkg/agent/tools"
)

// ErrorRecoveryType represents different types of recoverable errors
type ErrorRecoveryType string

const (
	ErrorTypeNoToolCall      ErrorRecoveryType = "no_tool_call"
	ErrorTypeInvalidJSON     ErrorRecoveryType = "invalid_json"
	ErrorTypeMissingToolName ErrorRecoveryType = "missing_tool_name"
	ErrorTypeUnknownTool     ErrorRecoveryType = "unknown_tool"
	ErrorTypeToolExecution   ErrorRecoveryType = "tool_execution"
)

// ErrorRecoveryContext contains data needed to build error recovery messages
type ErrorRecoveryContext struct {
	Type           ErrorRecoveryType
	Error          error
	ToolName       string
	Content        string
	AvailableTools []tools.Tool
}

// BuildErrorRecoveryMessage creates an error message with recovery instructions
// based on the error context
func BuildErrorRecoveryMessage(ctx ErrorRecoveryContext) string {
	switch ctx.Type {
	case ErrorTypeNoToolCall:
		return buildNoToolCallError()
	case ErrorTypeInvalidJSON:
		return buildParseError(ctx.Error, ctx.Content)
	case ErrorTypeMissingToolName:
		return buildMissingToolNameError()
	case ErrorTypeUnknownTool:
		return buildUnknownToolError(ctx.ToolName, ctx.AvailableTools)
	case ErrorTypeToolExecution:
		return buildToolExecutionError(ctx.ToolName, ctx.Error)
	default:
		return fmt.Sprintf("ERROR: An unknown error occurred: %v\n\nPlease try again.", ctx.Error)
	}
}

// buildNoToolCallError creates an error message with recovery instructions for missing tool calls
func buildNoToolCallError() string {
	return `ERROR: No tool call found in your response.

You MUST use a tool in every response. Available tools include task_completion, ask_question, converse, and any registered custom tools.

CORRECT FORMAT:
<tool>{"tool_name": "tool_name_here", "arguments": {...}}</tool>

Example:
<tool>{"tool_name": "task_completion", "arguments": {"result": "Task completed successfully"}}</tool>

Please try again with a valid tool call.`
}

// buildParseError creates an error message with recovery instructions for JSON parsing errors
func buildParseError(err error, content string) string {
	return fmt.Sprintf(`ERROR: Invalid tool call JSON.

Parse error: %v

Your tool call content: %s

CORRECT FORMAT:
<tool>{"tool_name": "calculator", "arguments": {"a": 10, "b": 20, "operation": "add"}}</tool>

Common JSON issues:
- Missing quotes around property names
- Trailing commas in objects or arrays
- Unescaped quotes inside strings
- Missing closing braces

Please fix the JSON syntax and try again.`, err, content)
}

// buildMissingToolNameError creates an error message for missing tool_name field
func buildMissingToolNameError() string {
	return `ERROR: Missing required field "tool_name" in tool call.

The tool_name field is required and must specify which tool to execute.

CORRECT FORMAT:
<tool>{"tool_name": "your_tool_here", "arguments": {...}}</tool>

Please include the tool_name field and try again.`
}

// buildUnknownToolError creates an error message with available tools listed
func buildUnknownToolError(toolName string, availableTools []tools.Tool) string {
	var toolNames []string
	for _, tool := range availableTools {
		toolNames = append(toolNames, fmt.Sprintf("- %s: %s", tool.Name(), tool.Description()))
	}

	return fmt.Sprintf(`ERROR: Unknown tool "%s".

Available tools:
%s

Please use one of the available tools and try again.`, toolName, strings.Join(toolNames, "\n"))
}

// buildToolExecutionError creates an error message for tool execution failures
func buildToolExecutionError(toolName string, err error) string {
	return fmt.Sprintf(`ERROR: Tool "%s" execution failed.

Error details: %v

Please review the error message, adjust your arguments if needed, and try again.
If the error persists, consider using a different approach or tool.`, toolName, err)
}
