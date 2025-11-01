package parser

import (
	"strings"
)

// ToolCallParser parses streaming content and separates <tool> tags from regular content.
// It maintains state across multiple content chunks to handle tags that span chunks.
type ToolCallParser struct {
	buffer      strings.Builder
	inToolCall  bool
	toolContent strings.Builder
}

// NewToolCallParser creates a new tool call parser.
func NewToolCallParser() *ToolCallParser {
	return &ToolCallParser{}
}

// ContentType represents the type of parsed content
type ContentType string

const (
	ContentTypeToolCall ContentType = "tool_call"
	ContentTypeRegular  ContentType = "regular"
)

// ParsedContent represents parsed content with its type
type ParsedContent struct {
	Type    ContentType
	Content string
}

// Parse processes a content chunk and returns separate chunks for tool calls and regular content.
// It handles <tool> tags that may span multiple chunks by buffering from '<' to '>'.
//
// Returns:
//   - toolCallContent: Non-nil if tool call content is found
//   - regularContent: Non-nil if regular content is found
func (p *ToolCallParser) Parse(content string) (toolCallContent, regularContent *ParsedContent) {
	if content == "" {
		return nil, nil
	}

	for _, char := range content {
		p.buffer.WriteRune(char)

		if char == '<' && !p.inToolCall {
			chunk := p.handleTagStart()
			toolCallContent, regularContent = p.appendContent(toolCallContent, regularContent, chunk)
			continue
		}

		if char == '>' {
			chunk := p.handleTagEnd()
			toolCallContent, regularContent = p.appendContent(toolCallContent, regularContent, chunk)
			continue
		}
	}

	chunk := p.flushBufferIfNotInTag()
	toolCallContent, regularContent = p.appendContent(toolCallContent, regularContent, chunk)

	return toolCallContent, regularContent
}

// handleTagStart processes the start of a potential tag
func (p *ToolCallParser) handleTagStart() *ParsedContent {
	// Flush everything before '<' as regular content
	text := p.buffer.String()
	p.buffer.Reset()
	p.buffer.WriteRune('<')

	if len(text) > 1 {
		return &ParsedContent{
			Type:    ContentTypeRegular,
			Content: text[:len(text)-1],
		}
	}
	return nil
}

// handleTagEnd processes the end of a potential tag
func (p *ToolCallParser) handleTagEnd() *ParsedContent {
	tag := p.buffer.String()
	p.buffer.Reset()

	if tag == "<tool>" {
		p.inToolCall = true
		p.toolContent.Reset()
		return nil
	}

	if tag == "</tool>" {
		p.inToolCall = false
		content := p.toolContent.String()
		p.toolContent.Reset()
		return &ParsedContent{
			Type:    ContentTypeToolCall,
			Content: content,
		}
	}

	// Not a tool tag, treat as regular content
	return &ParsedContent{
		Type:    ContentTypeRegular,
		Content: tag,
	}
}

// flushBufferIfNotInTag flushes buffered content if we're not in the middle of parsing a tag
func (p *ToolCallParser) flushBufferIfNotInTag() *ParsedContent {
	if p.buffer.Len() == 0 {
		return nil
	}

	// If we're in a tool call, add buffer to tool content
	if p.inToolCall {
		p.toolContent.WriteString(p.buffer.String())
		p.buffer.Reset()
		return nil
	}

	// If buffer doesn't start with '<', it's regular content
	text := p.buffer.String()
	if !strings.HasPrefix(text, "<") {
		p.buffer.Reset()
		return &ParsedContent{
			Type:    ContentTypeRegular,
			Content: text,
		}
	}

	// Buffer starts with '<' but hasn't found '>' yet - keep buffering
	return nil
}

// appendContent appends new content to existing content based on type
func (p *ToolCallParser) appendContent(toolCallContent, regularContent, newContent *ParsedContent) (*ParsedContent, *ParsedContent) {
	if newContent == nil {
		return toolCallContent, regularContent
	}

	if newContent.Type == ContentTypeToolCall {
		if toolCallContent == nil {
			return newContent, regularContent
		}
		toolCallContent.Content += newContent.Content
		return toolCallContent, regularContent
	}

	if regularContent == nil {
		return toolCallContent, newContent
	}
	regularContent.Content += newContent.Content
	return toolCallContent, regularContent
}

// IsInToolCall returns true if currently parsing tool call content.
func (p *ToolCallParser) IsInToolCall() bool {
	return p.inToolCall
}

// Flush returns any remaining buffered content and resets the parser.
// This should be called at the end of a stream to ensure all content is processed.
func (p *ToolCallParser) Flush() (toolCallContent, regularContent *ParsedContent) {
	text := p.buffer.String()
	if text == "" && p.toolContent.Len() == 0 {
		return nil, nil
	}

	// If we have tool content buffered, return it
	if p.toolContent.Len() > 0 {
		toolCallContent = &ParsedContent{
			Type:    ContentTypeToolCall,
			Content: p.toolContent.String(),
		}
		p.toolContent.Reset()
	}

	// Return any remaining buffer as regular content
	if text != "" {
		regularContent = &ParsedContent{
			Type:    ContentTypeRegular,
			Content: text,
		}
		p.buffer.Reset()
	}

	p.inToolCall = false
	return toolCallContent, regularContent
}

// Reset clears all parser state
func (p *ToolCallParser) Reset() {
	p.buffer.Reset()
	p.toolContent.Reset()
	p.inToolCall = false
}
