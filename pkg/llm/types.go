package llm

// ContentType indicates the type of content in a StreamChunk.
type ContentType string

const (
	// ContentTypeMessage indicates regular message content (the final answer).
	ContentTypeMessage ContentType = "message"

	// ContentTypeThinking indicates thinking/reasoning content (intermediate thoughts).
	// Some models (like Claude with extended_thinking or OpenAI o1) emit this separately.
	ContentTypeThinking ContentType = "thinking"
)

// StreamChunk represents a single chunk from an LLM streaming response.
// This is a provider-layer type focused purely on LLM output, with no
// coupling to agent-level events or orchestration.
type StreamChunk struct {
	// Content is the text delta from the LLM response.
	// For streaming responses, this contains incremental text.
	Content string

	// Role is the message role (e.g., "assistant", "user", "system").
	// This is typically only set on the first chunk of a response.
	Role string

	// Type indicates whether this is thinking content or message content.
	// Defaults to ContentTypeMessage if not set.
	Type ContentType

	// Finished indicates whether this is the final chunk in the stream.
	// When true, no more chunks will be sent on the channel.
	Finished bool

	// Error contains any error that occurred during streaming.
	// When set, this is typically the last chunk sent before closing the channel.
	Error error
}

// IsError returns true if this chunk contains an error.
func (c *StreamChunk) IsError() bool {
	return c.Error != nil
}

// IsLast returns true if this is the final chunk in the stream.
func (c *StreamChunk) IsLast() bool {
	return c.Finished
}

// HasContent returns true if this chunk contains text content.
func (c *StreamChunk) HasContent() bool {
	return c.Content != ""
}

// IsThinking returns true if this chunk contains thinking/reasoning content.
func (c *StreamChunk) IsThinking() bool {
	return c.Type == ContentTypeThinking
}

// IsMessage returns true if this chunk contains message content (or Type is unset).
func (c *StreamChunk) IsMessage() bool {
	return c.Type == "" || c.Type == ContentTypeMessage
}
