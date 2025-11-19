package agent

import (
	"context"
	"fmt"

	"github.com/entrhq/forge/pkg/agent/prompts"
	"github.com/entrhq/forge/pkg/agent/tools"
	"github.com/entrhq/forge/pkg/types"
)

// validateToolCallFields validates required fields in the tool call
// Returns (shouldContinue, errorContext)
func (a *DefaultAgent) validateToolCallFields(toolCall *tools.ToolCall) (bool, string) {
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

	return true, ""
}

// parseToolCallXML parses tool call XML content and handles errors
// Returns (toolCall, shouldContinue, errorContext)
func (a *DefaultAgent) parseToolCallXML(toolCallContent string) (tools.ToolCall, bool, string) {
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
			return tools.ToolCall{}, false, ""
		}

		a.emitEvent(types.NewErrorEvent(fmt.Errorf("failed to parse tool call: %w", err)))
		return tools.ToolCall{}, true, errMsg
	}

	// Use the parsed tool call
	return *parsedToolCall, true, ""
}

// validateToolCallContent checks if context was canceled and if tool call content exists
// Returns (shouldContinue, errorContext) - if errorContext is non-empty, validation failed
func (a *DefaultAgent) validateToolCallContent(ctx context.Context, toolCallContent string) (bool, string) {
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

	// Validation passed
	return true, ""
}

// processToolCall handles parsing, validation, and execution of tool calls
// Returns (shouldContinue, errorContext) following the same pattern as executeIteration
func (a *DefaultAgent) processToolCall(ctx context.Context, toolCallContent string) (bool, string) {
	// Validate content exists and context not canceled
	shouldContinue, errCtx := a.validateToolCallContent(ctx, toolCallContent)
	if errCtx != "" || !shouldContinue {
		return shouldContinue, errCtx
	}

	// Parse the tool call XML
	toolCall, shouldContinue, errCtx := a.parseToolCallXML(toolCallContent)
	if errCtx != "" || !shouldContinue {
		return shouldContinue, errCtx
	}

	// Validate required fields
	shouldContinue, errCtx = a.validateToolCallFields(&toolCall)
	if errCtx != "" || !shouldContinue {
		return shouldContinue, errCtx
	}

	// Execute the tool
	return a.executeTool(ctx, toolCall)
}
