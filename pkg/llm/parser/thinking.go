// Package parser provides utilities for parsing structured content from LLM streams.
package parser

import (
	"strings"

	"github.com/clear-route/forge/pkg/llm"
)

// ThinkingParser parses streaming content and separates <thinking> tags from regular content.
// It maintains state across multiple content chunks to handle tags that span chunks.
type ThinkingParser struct {
	buffer     strings.Builder
	inThinking bool
	inTag      bool // true when we're buffering a potential tag (saw '<' but not yet '>')
}

// NewThinkingParser creates a new thinking parser.
func NewThinkingParser() *ThinkingParser {
	return &ThinkingParser{}
}

// Parse processes a content chunk and returns separate chunks for thinking and message content.
// It handles <thinking> tags that may span multiple chunks by buffering from '<' to '>'.
//
// Returns:
//   - thinkingChunk: Non-nil if thinking content is found (with Type = ContentTypeThinking)
//   - messageChunk: Non-nil if message content is found (with Type = ContentTypeMessage)
func (p *ThinkingParser) Parse(content string) (thinkingChunk, messageChunk *llm.StreamChunk) {
	if content == "" {
		return nil, nil
	}

	for _, ch := range content {
		p.buffer.WriteRune(ch)

		if ch == '<' {
			chunk := p.handleTagStart()
			thinkingChunk, messageChunk = p.appendChunk(thinkingChunk, messageChunk, chunk)
			continue
		}

		if ch == '>' && p.inTag {
			chunk := p.handleTagEnd()
			thinkingChunk, messageChunk = p.appendChunk(thinkingChunk, messageChunk, chunk)
			continue
		}
	}

	// Emit accumulated non-tag content at end of chunk
	chunk := p.flushBufferIfNotInTag()
	thinkingChunk, messageChunk = p.appendChunk(thinkingChunk, messageChunk, chunk)

	return
}

// handleTagStart processes the start of a potential tag
func (p *ThinkingParser) handleTagStart() *llm.StreamChunk {
	// If we were buffering non-tag content, emit it first
	if p.buffer.Len() > 1 && !p.inTag {
		text := p.buffer.String()[:p.buffer.Len()-1] // Exclude the '<'

		chunk := p.createChunk(text)

		p.buffer.Reset()
		p.buffer.WriteRune('<')
		p.inTag = true

		return chunk
	}

	p.inTag = true
	return nil
}

// handleTagEnd processes the end of a tag
func (p *ThinkingParser) handleTagEnd() *llm.StreamChunk {
	p.inTag = false
	tag := p.buffer.String()
	p.buffer.Reset()

	if tag == "<thinking>" {
		p.inThinking = true
		return nil
	}

	if tag == "</thinking>" {
		p.inThinking = false
		return nil
	}

	// Not a thinking tag, treat as regular content
	return p.createChunk(tag)
}

// flushBufferIfNotInTag flushes buffer content if not currently in a tag
func (p *ThinkingParser) flushBufferIfNotInTag() *llm.StreamChunk {
	if !p.inTag && p.buffer.Len() > 0 {
		text := p.buffer.String()
		p.buffer.Reset()
		return p.createChunk(text)
	}
	return nil
}

// createChunk creates a chunk with appropriate type based on current mode
func (p *ThinkingParser) createChunk(text string) *llm.StreamChunk {
	if text == "" {
		return nil
	}

	if p.inThinking {
		return &llm.StreamChunk{
			Content: text,
			Type:    llm.ContentTypeThinking,
		}
	}

	return &llm.StreamChunk{
		Content: text,
		Type:    llm.ContentTypeMessage,
	}
}

// appendChunk appends a new chunk to existing chunks based on type
func (p *ThinkingParser) appendChunk(thinkingChunk, messageChunk, newChunk *llm.StreamChunk) (*llm.StreamChunk, *llm.StreamChunk) {
	if newChunk == nil {
		return thinkingChunk, messageChunk
	}

	if newChunk.Type == llm.ContentTypeThinking {
		if thinkingChunk == nil {
			return newChunk, messageChunk
		}
		thinkingChunk.Content += newChunk.Content
		return thinkingChunk, messageChunk
	}

	if messageChunk == nil {
		return thinkingChunk, newChunk
	}
	messageChunk.Content += newChunk.Content
	return thinkingChunk, messageChunk
}

// IsInThinking returns true if currently parsing thinking content.
func (p *ThinkingParser) IsInThinking() bool {
	return p.inThinking
}

// Flush returns any buffered content that hasn't been emitted yet.
// This should be called at the end of a stream to ensure all content is processed.
func (p *ThinkingParser) Flush() (thinkingChunk, messageChunk *llm.StreamChunk) {
	text := p.buffer.String()
	if len(text) == 0 {
		return nil, nil
	}

	// Emit whatever is in the buffer
	p.buffer.Reset()
	p.inTag = false

	if p.inThinking {
		return &llm.StreamChunk{
			Content: text,
			Type:    llm.ContentTypeThinking,
		}, nil
	}

	return nil, &llm.StreamChunk{
		Content: text,
		Type:    llm.ContentTypeMessage,
	}
}

// Reset resets the parser state for a new stream.
func (p *ThinkingParser) Reset() {
	p.buffer.Reset()
	p.inThinking = false
	p.inTag = false
}
