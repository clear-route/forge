package parser

import (
	"strings"
)

// ToolCallParser parses streaming content and separates <tool> tags from regular content.
// It maintains state across multiple content chunks to handle tags that span chunks.
type ToolCallParser struct {
	buffer      strings.Builder
	inToolCall  bool
	inTag       bool // true when we're buffering a potential tag (saw '<' but not yet '>')
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
// It handles <tool> tags that may span multiple chunks by buffering potential tags.
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
		bufStr := p.buffer.String()

		// Check for complete <tool> tag
		if !p.inToolCall && len(bufStr) >= 6 && bufStr[len(bufStr)-6:] == "<tool>" {
			chunk := p.handleToolStart()
			toolCallContent, regularContent = p.appendContent(toolCallContent, regularContent, chunk)
			continue
		}

		// Check for complete </tool> tag
		if p.inToolCall && len(bufStr) >= 7 && bufStr[len(bufStr)-7:] == "</tool>" {
			chunk := p.handleToolEnd()
			toolCallContent, regularContent = p.appendContent(toolCallContent, regularContent, chunk)
			continue
		}
	}

	// Flush any buffered content that's not part of a potential tag
	chunk := p.flushBufferIfNotInTag()
	toolCallContent, regularContent = p.appendContent(toolCallContent, regularContent, chunk)

	return toolCallContent, regularContent
}

// handleToolStart processes the start of <tool> tag
func (p *ToolCallParser) handleToolStart() *ParsedContent {
	// Get content before <tool>
	bufStr := p.buffer.String()
	textBefore := bufStr[:len(bufStr)-6] // Everything except "<tool>"

	p.buffer.Reset()
	p.inToolCall = true
	p.toolContent.Reset()

	if textBefore != "" {
		return &ParsedContent{
			Type:    ContentTypeRegular,
			Content: textBefore,
		}
	}
	return nil
}

// handleToolEnd processes the end of </tool> tag
func (p *ToolCallParser) handleToolEnd() *ParsedContent {
	// Get content before </tool>
	bufStr := p.buffer.String()
	textBefore := bufStr[:len(bufStr)-7] // Everything except "</tool>"

	p.buffer.Reset()
	p.inToolCall = false

	// Add any remaining buffered content to tool content
	p.toolContent.WriteString(textBefore)
	content := p.toolContent.String()
	p.toolContent.Reset()

	// Trim any trailing whitespace or stray characters
	content = strings.TrimSpace(content)

	return &ParsedContent{
		Type:    ContentTypeToolCall,
		Content: content,
	}
}

// flushBufferIfNotInTag flushes buffered content, keeping potential tag prefixes in buffer
func (p *ToolCallParser) flushBufferIfNotInTag() *ParsedContent {
	if p.buffer.Len() == 0 {
		return nil
	}

	text := p.buffer.String()

	// Keep potential tag prefixes in the buffer
	// This handles streaming where tags may be split across chunks
	var flushText string
	if p.inToolCall {
		// Inside tool call - keep partial </tool> prefixes buffered
		// Move everything else to toolContent but DON'T emit as events yet
		// This prevents premature JSON parsing before </tool> arrives
		if len(text) >= 6 && text[len(text)-6:] == "</tool" {
			flushText = text[:len(text)-6]
			p.buffer.Reset()
			p.buffer.WriteString("</tool")
		} else if len(text) >= 5 && text[len(text)-5:] == "</too" {
			flushText = text[:len(text)-5]
			p.buffer.Reset()
			p.buffer.WriteString("</too")
		} else if len(text) >= 4 && text[len(text)-4:] == "</to" {
			flushText = text[:len(text)-4]
			p.buffer.Reset()
			p.buffer.WriteString("</to")
		} else if len(text) >= 3 && text[len(text)-3:] == "</t" {
			flushText = text[:len(text)-3]
			p.buffer.Reset()
			p.buffer.WriteString("</t")
		} else if len(text) >= 2 && text[len(text)-2:] == "</" {
			flushText = text[:len(text)-2]
			p.buffer.Reset()
			p.buffer.WriteString("</")
		} else if len(text) >= 1 && text[len(text)-1:] == "<" {
			flushText = text[:len(text)-1]
			p.buffer.Reset()
			p.buffer.WriteString("<")
		} else {
			// No partial closing tag detected
			// Move to toolContent but don't emit yet (wait for </tool>)
			flushText = text
			p.buffer.Reset()
		}
	} else {
		// Not in tool call, check for potential <tool> prefix
		if len(text) >= 5 && text[len(text)-5:] == "<tool" {
			flushText = text[:len(text)-5]
			p.buffer.Reset()
			p.buffer.WriteString("<tool")
		} else if len(text) >= 4 && text[len(text)-4:] == "<too" {
			flushText = text[:len(text)-4]
			p.buffer.Reset()
			p.buffer.WriteString("<too")
		} else if len(text) >= 3 && text[len(text)-3:] == "<to" {
			flushText = text[:len(text)-3]
			p.buffer.Reset()
			p.buffer.WriteString("<to")
		} else if len(text) >= 2 && text[len(text)-2:] == "<t" {
			flushText = text[:len(text)-2]
			p.buffer.Reset()
			p.buffer.WriteString("<t")
		} else if len(text) >= 1 && text[len(text)-1:] == "<" {
			flushText = text[:len(text)-1]
			p.buffer.Reset()
			p.buffer.WriteString("<")
		} else {
			flushText = text
			p.buffer.Reset()
		}
	}

	if flushText == "" {
		return nil
	}

	// If we're in a tool call, accumulate to toolContent but DON'T emit
	// The content will only be emitted when we see the complete </tool> tag
	if p.inToolCall {
		p.toolContent.WriteString(flushText)
		return nil
	}

	// Otherwise return as regular content
	return &ParsedContent{
		Type:    ContentTypeRegular,
		Content: flushText,
	}
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
	p.inTag = false
}
