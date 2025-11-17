package core

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/entrhq/forge/pkg/llm"
	"github.com/entrhq/forge/pkg/llm/parser"
	"github.com/entrhq/forge/pkg/types"
)

// streamState tracks the current state of stream processing
type streamState struct {
	assistantContent string
	thinkingContent  string
	toolCallContent  string
	toolCallBuffer   string // buffer content until tool name is detected
	role             string
	messageStarted   bool
	thinkingStarted  bool
	toolCallStarted  bool
	toolNameDetected bool // tracks if we've detected and emitted the tool name
	toolNameEmitted  bool // tracks if we've emitted buffered content after tool name
	toolCallParser   *parser.ToolCallParser
}

// ProcessStream processes a stream of chunks, emitting events and calling
// the completion handler when done. This provides reusable stream processing
// logic that any agent can use.
func ProcessStream(
	stream <-chan *llm.StreamChunk,
	emitEvent func(*types.AgentEvent),
	onComplete func(assistantContent, thinkingContent, toolCallContent, role string),
) {
	state := &streamState{
		toolCallParser: parser.NewToolCallParser(),
	}

	for chunk := range stream {
		if chunk.IsError() {
			handleError(chunk.Error, state, emitEvent)
			return
		}

		if chunk.Role != "" {
			state.role = chunk.Role
		}

		if chunk.HasContent() {
			handleContent(chunk, state, emitEvent)
		}

		if chunk.IsLast() {
			finalize(state, emitEvent, onComplete)
			return
		}
	}

	// Stream ended without explicit finish
	finalize(state, emitEvent, onComplete)
}

// handleError handles error chunks and cleans up state
func handleError(err error, state *streamState, emitEvent func(*types.AgentEvent)) {
	if state.thinkingStarted {
		emitEvent(types.NewThinkingEndEvent())
	}
	if state.toolCallStarted {
		emitEvent(types.NewToolCallEndEvent())
	}
	if state.messageStarted {
		emitEvent(types.NewMessageEndEvent())
	}

	// Don't emit error events for context cancellation - this is expected when user stops the agent
	if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
		// Check if error message contains "context canceled" (some wrapped errors don't preserve type)
		if err.Error() != "stream read error: context canceled" {
			emitEvent(types.NewErrorEvent(err))
		}
	}
}

// handleContent processes content chunks based on type
func handleContent(chunk *llm.StreamChunk, state *streamState, emitEvent func(*types.AgentEvent)) {
	if chunk.IsThinking() {
		handleThinkingContent(chunk.Content, state, emitEvent)
	} else {
		handleMessageContent(chunk.Content, state, emitEvent)
	}
}

// handleThinkingContent processes thinking content
func handleThinkingContent(content string, state *streamState, emitEvent func(*types.AgentEvent)) {
	if !state.thinkingStarted {
		emitEvent(types.NewThinkingStartEvent())
		state.thinkingStarted = true
	}
	emitEvent(types.NewThinkingContentEvent(content))
	state.thinkingContent += content
}

// handleMessageContent processes message content
// It also parses out tool call XML tags and emits them separately
func handleMessageContent(content string, state *streamState, emitEvent func(*types.AgentEvent)) {
	// End thinking if it was active
	if state.thinkingStarted {
		emitEvent(types.NewThinkingEndEvent())
		state.thinkingStarted = false
	}

	// Parse content for tool calls
	toolCallContent, regularContent := state.toolCallParser.Parse(content)

	// Handle tool call start signal - emit immediately when <tool> is detected
	if toolCallContent != nil && toolCallContent.Type == "tool_call_start" {
		// Close any active message before starting tool call
		if state.messageStarted {
			emitEvent(types.NewMessageEndEvent())
			state.messageStarted = false
		}

		// Emit the tool call start event for immediate UI feedback
		if !state.toolCallStarted {
			emitEvent(types.NewToolCallStartEvent())
			state.toolCallStarted = true
		}
		return
	}

	// Check for tool name in accumulated content after tool call start
	checkAndEmitToolName(state, emitEvent)

	// Handle complete tool call content (when </tool> is detected)
	if toolCallContent != nil && toolCallContent.Type == "tool_call" && toolCallContent.Content != "" {
		handleToolCallContent(toolCallContent.Content, state, emitEvent)
	}

	// Handle regular message content
	if regularContent != nil && regularContent.Content != "" {
		handleRegularContent(regularContent.Content, state, emitEvent)
	}
}

// checkAndEmitToolName checks for tool name in accumulated content and emits early detection event
func checkAndEmitToolName(state *streamState, emitEvent func(*types.AgentEvent)) {
	if !state.toolCallStarted || state.toolNameDetected {
		return
	}

	accumulatedContent := state.toolCallParser.GetAccumulatedToolContent()
	toolName := extractToolNameFromPartial(accumulatedContent)

	if toolName != "" {
		state.toolNameDetected = true
		event := types.NewToolCallStartEvent()
		event.Metadata["tool_name"] = toolName
		emitEvent(event)
	}
}

// handleRegularContent processes regular message content
func handleRegularContent(content string, state *streamState, emitEvent func(*types.AgentEvent)) {
	// End tool call if it was active
	if state.toolCallStarted {
		emitEvent(types.NewToolCallEndEvent())
		state.toolCallStarted = false
		state.toolNameDetected = false
	}

	if !state.messageStarted {
		emitEvent(types.NewMessageStartEvent())
		state.messageStarted = true
	}

	emitEvent(types.NewMessageContentEvent(content))
	state.assistantContent += content
}

// handleToolCallContent processes tool call XML content
func handleToolCallContent(content string, state *streamState, emitEvent func(*types.AgentEvent)) {
	// End message if it was active
	if state.messageStarted {
		emitEvent(types.NewMessageEndEvent())
		state.messageStarted = false
	}

	// Emit initial tool call start event only once
	if !state.toolCallStarted {
		emitEvent(types.NewToolCallStartEvent())
		state.toolCallStarted = true
	}

	// Accumulate tool call content
	state.toolCallContent += content

	// If we haven't detected the tool name yet, buffer the content
	if !state.toolNameDetected {
		state.toolCallBuffer += content

		// Try to detect tool name from accumulated buffer
		if toolName := extractToolNameFromPartial(state.toolCallContent); toolName != "" {
			state.toolNameDetected = true

			// Emit EventTypeToolCallStart with the tool name in metadata
			// This provides early feedback to the UI
			event := types.NewToolCallStartEvent()
			event.Metadata["tool_name"] = toolName
			emitEvent(event)

			// Now emit all buffered content at once
			emitEvent(types.NewToolCallContentEvent(state.toolCallBuffer))
			state.toolNameEmitted = true
		}
		// Don't emit content events until we have the tool name
		return
	}

	// After tool name is detected, emit content normally
	emitEvent(types.NewToolCallContentEvent(content))
}

// finalize ends the stream processing
func finalize(state *streamState, emitEvent func(*types.AgentEvent), onComplete func(string, string, string, string)) {
	// Flush any remaining content from tool call parser
	toolCallContent, regularContent := state.toolCallParser.Flush()
	if toolCallContent != nil && toolCallContent.Content != "" {
		handleToolCallContent(toolCallContent.Content, state, emitEvent)
	}
	if regularContent != nil && regularContent.Content != "" {
		if state.toolCallStarted {
			emitEvent(types.NewToolCallEndEvent())
			state.toolCallStarted = false
		}
		if !state.messageStarted {
			emitEvent(types.NewMessageStartEvent())
			state.messageStarted = true
		}
		emitEvent(types.NewMessageContentEvent(regularContent.Content))
		state.assistantContent += regularContent.Content
	}

	// End any active sections
	if state.thinkingStarted {
		emitEvent(types.NewThinkingEndEvent())
	}
	if state.toolCallStarted {
		emitEvent(types.NewToolCallEndEvent())
	}
	if state.messageStarted {
		emitEvent(types.NewMessageEndEvent())
	}

	role := state.role
	if role == "" {
		role = string(types.RoleAssistant)
	}
	onComplete(state.assistantContent, state.thinkingContent, state.toolCallContent, role)
}

// extractToolNameFromPartial attempts to extract the tool name from partial XML content.
// It looks for the <tool_name>value</tool_name> pattern and returns the tool name if found.
// Returns empty string if the pattern is not yet complete or malformed.
func extractToolNameFromPartial(content string) string {
	// Pattern: <tool_name>value</tool_name>
	// Must be strict to avoid false positives
	// Matches: opening tag, whitespace (optional), non-empty value (no < or >), whitespace (optional), closing tag
	re := regexp.MustCompile(`<tool_name>\s*([^<>\s][^<>]*?)\s*</tool_name>`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}
