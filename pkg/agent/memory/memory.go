package memory

import (
	"github.com/entrhq/forge/pkg/types"
)

// Memory represents a conversation history management system
// that can store, retrieve, and prune messages.
type Memory interface {
	// Add appends a message to the conversation history
	Add(msg *types.Message)

	// GetAll returns all messages in the conversation history
	GetAll() []*types.Message

	// GetRecent returns the most recent N messages
	GetRecent(n int) []*types.Message

	// Clear removes all messages from the conversation history
	Clear()

	// Prune reduces the conversation history to fit within a token limit
	// while preserving important context (system messages, recent messages)
	Prune(maxTokens int) error

	// Count returns the number of messages in the conversation history
	Count() int
}
