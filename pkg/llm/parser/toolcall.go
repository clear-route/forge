package parser

import (
	"strings"
)

// ToolCallParser parses streaming content and separates <tool> tags from regular content.
// It maintains state across multiple content chunks to handle tags that span chunks.
type ToolCallParser struct {
	buffer      strings.Builder
	toolContent strings.Builder
	inToolCall  bool
	inTag       bool
}

// NewToolCallParser creates a new tool call parser.
func NewToolCallParser() *ToolCallParser {
	return &ToolCallParser{}
}

// ContentType represents the type of parsed content
type ContentType string

const (
	ContentTypeToolCall      ContentType = "tool_call"
	ContentTypeToolCallStart ContentType = "tool_call_start"
	ContentTypeRegular       ContentType = "regular"
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
	p.toolContent.Reset()
	p.inToolCall = true

	// Always return a tool call start signal (even if there's no text before)
	// This allows immediate emission of tool call start event
	return &ParsedContent{
		Type:    ContentTypeToolCallStart,
		Content: textBefore, // May be empty, which is fine
	}
}

// handleToolEnd processes the end of </tool> tag
func (p *ToolCallParser) handleToolEnd() *ParsedContent {
	// First, move any remaining buffer content (minus the "</tool>" tag) to toolContent
	bufStr := p.buffer.String()
	if len(bufStr) >= 7 {
		// Add everything except the "</tool>" tag to toolContent
		contentBeforeTag := bufStr[:len(bufStr)-7]
		p.toolContent.WriteString(contentBeforeTag)
	}
	
	p.buffer.Reset()
	p.inToolCall = false

	// Get the accumulated tool content
	content := strings.TrimSpace(p.toolContent.String())
	p.toolContent.Reset()

	// Return empty tool call even if content is empty
	return &ParsedContent{
		Type:    ContentTypeToolCall,
		Content: content,
	}
}

// flushBufferIfNotInTag flushes buffered content, keeping potential tag prefixes in buffer
// flushBufferIfNotInTag handles the complex state machine for detecting XML tag boundaries
// in streaming content. High complexity is inherent to the XML parsing state machine that must
// detect partial opening/closing tags across chunk boundaries.
//
//nolint:gocyclo
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
		switch {
		case len(text) >= 6 && text[len(text)-6:] == "</tool":
			flushText = text[:len(text)-6]
			p.buffer.Reset()
			p.buffer.WriteString("</tool")
		case len(text) >= 5 && text[len(text)-5:] == "</too":
			flushText = text[:len(text)-5]
			p.buffer.Reset()
			p.buffer.WriteString("</too")
		case len(text) >= 4 && text[len(text)-4:] == "</to":
			flushText = text[:len(text)-4]
			p.buffer.Reset()
			p.buffer.WriteString("</to")
		case len(text) >= 3 && text[len(text)-3:] == "</t":
			flushText = text[:len(text)-3]
			p.buffer.Reset()
			p.buffer.WriteString("</t")
		case len(text) >= 2 && text[len(text)-2:] == "</":
			flushText = text[:len(text)-2]
			p.buffer.Reset()
			p.buffer.WriteString("</")
		case len(text) >= 1 && text[len(text)-1:] == "<":
			flushText = text[:len(text)-1]
			p.buffer.Reset()
			p.buffer.WriteString("<")
		default:
			// No partial closing tag detected
			// Move to toolContent but don't emit yet (wait for </tool>)
			flushText = text
			p.buffer.Reset()
		}
	} else {
		// Not in tool call, check for potential <tool> prefix
		switch {
		case len(text) >= 5 && text[len(text)-5:] == "<tool":
			flushText = text[:len(text)-5]
			p.buffer.Reset()
			p.buffer.WriteString("<tool")
		case len(text) >= 4 && text[len(text)-4:] == "<too":
			flushText = text[:len(text)-4]
			p.buffer.Reset()
			p.buffer.WriteString("<too")
		case len(text) >= 3 && text[len(text)-3:] == "<to":
			flushText = text[:len(text)-3]
			p.buffer.Reset()
			p.buffer.WriteString("<to")
		case len(text) >= 2 && text[len(text)-2:] == "<t":
			flushText = text[:len(text)-2]
			p.buffer.Reset()
			p.buffer.WriteString("<t")
		case len(text) >= 1 && text[len(text)-1:] == "<":
			flushText = text[:len(text)-1]
			p.buffer.Reset()
			p.buffer.WriteString("<")
		default:
			flushText = text
			p.buffer.Reset()
		}
	}

	if flushText == "" {
		return nil
	}

	// When inside a tool call, accumulate flushed content into toolContent
	// instead of emitting it immediately. We'll emit when </tool> is detected.
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

	// Handle tool call start - return it as toolCallContent (for signaling)
	// but don't append the text content (it should go to regularContent if present)
	if newContent.Type == ContentTypeToolCallStart {
		// If there's text before the <tool> tag, add it to regularContent
		if newContent.Content != "" {
			if regularContent == nil {
				regularContent = &ParsedContent{
					Type:    ContentTypeRegular,
					Content: newContent.Content,
				}
			} else {
				regularContent.Content += newContent.Content
			}
		}
		// Return the start signal as toolCallContent
		return newContent, regularContent
	}

	if newContent.Type == ContentTypeToolCall {
		// If we already have a tool call start signal, replace it with the complete tool call
		if toolCallContent != nil && toolCallContent.Type == ContentTypeToolCallStart {
			return newContent, regularContent
		}
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

// GetAccumulatedToolContent returns the currently accumulated tool content
// This is useful for extracting partial information like tool name before the tool call is complete
func (p *ToolCallParser) GetAccumulatedToolContent() string {
	return p.toolContent.String()
}

// Flush returns any remaining buffered content and resets the parser.
// This should be called at the end of a stream to ensure all content is processed.
func (p *ToolCallParser) Flush() (toolCallContent, regularContent *ParsedContent) {
	// First flush any buffered content into toolContent if we're in a tool call
	if p.inToolCall && p.buffer.Len() > 0 {
		p.toolContent.WriteString(p.buffer.String())
		p.buffer.Reset()
	}

	// If we're in a tool call, return accumulated tool content
	if p.inToolCall {
		p.inToolCall = false
		content := strings.TrimSpace(p.toolContent.String())
		p.toolContent.Reset()
		
		if content != "" {
			toolCallContent = &ParsedContent{
				Type:    ContentTypeToolCall,
				Content: content,
			}
			return toolCallContent, nil
		}
		return nil, nil
	}

	// Otherwise, return buffered content as regular content
	text := p.buffer.String()
	p.buffer.Reset()
	
	if text == "" {
		return nil, nil
	}

	regularContent = &ParsedContent{
		Type:    ContentTypeRegular,
		Content: text,
	}
	return nil, regularContent
}

// Reset clears all parser state
func (p *ToolCallParser) Reset() {
	p.buffer.Reset()
	p.inToolCall = false
	p.inTag = false
}
