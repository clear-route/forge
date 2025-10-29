package agent

import (
	"github.com/entrhq/forge/pkg/llm"
	"github.com/entrhq/forge/pkg/types"
)

// streamState tracks the current state of stream processing
type streamState struct {
	assistantContent string
	thinkingContent  string
	role             string
	messageStarted   bool
	thinkingStarted  bool
}

// ProcessStream processes a stream of chunks, emitting events and calling
// the completion handler when done. This provides reusable stream processing
// logic that any agent can use.
func ProcessStream(
	stream <-chan *llm.StreamChunk,
	emitEvent func(*types.AgentEvent),
	onComplete func(assistantContent, thinkingContent, role string),
) {
	state := &streamState{}

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
	if state.messageStarted {
		emitEvent(types.NewMessageEndEvent())
	}
	emitEvent(types.NewErrorEvent(err))
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
func handleMessageContent(content string, state *streamState, emitEvent func(*types.AgentEvent)) {
	if state.thinkingStarted {
		emitEvent(types.NewThinkingEndEvent())
		state.thinkingStarted = false
	}
	if !state.messageStarted {
		emitEvent(types.NewMessageStartEvent())
		state.messageStarted = true
	}
	emitEvent(types.NewMessageContentEvent(content))
	state.assistantContent += content
}

// finalize ends the stream processing
func finalize(state *streamState, emitEvent func(*types.AgentEvent), onComplete func(string, string, string)) {
	if state.thinkingStarted {
		emitEvent(types.NewThinkingEndEvent())
	}
	if state.messageStarted {
		emitEvent(types.NewMessageEndEvent())
	}
	role := state.role
	if role == "" {
		role = string(types.RoleAssistant)
	}
	onComplete(state.assistantContent, state.thinkingContent, role)
}
