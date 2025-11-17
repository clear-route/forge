package core

import (
	"context"
	"errors"

	"github.com/entrhq/forge/pkg/llm"
	"github.com/entrhq/forge/pkg/llm/parser"
	"github.com/entrhq/forge/pkg/types"
)

// streamState tracks the current state of stream processing
type streamState struct {
	assistantContent string
	thinkingContent  string
	toolCallContent  string
	role             string
	messageStarted   bool
	thinkingStarted  bool
	toolCallStarted  bool
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

	// Handle tool call content
	if toolCallContent != nil && toolCallContent.Content != "" {
		handleToolCallContent(toolCallContent.Content, state, emitEvent)
	}

	// Handle regular message content
	if regularContent != nil && regularContent.Content != "" {
		// End tool call if it was active
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
}

// handleToolCallContent processes tool call XML content
func handleToolCallContent(content string, state *streamState, emitEvent func(*types.AgentEvent)) {
	// End message if it was active
	if state.messageStarted {
		emitEvent(types.NewMessageEndEvent())
		state.messageStarted = false
	}

	if !state.toolCallStarted {
		emitEvent(types.NewToolCallStartEvent())
		state.toolCallStarted = true
	}
	emitEvent(types.NewToolCallContentEvent(content))
	state.toolCallContent += content
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
