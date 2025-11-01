package memory

import (
	"fmt"
	"sync"

	"github.com/entrhq/forge/pkg/types"
)

// ConversationMemory is a thread-safe implementation of Memory
// that stores conversation history in memory.
type ConversationMemory struct {
	messages []*types.Message
	mu       sync.RWMutex
}

// NewConversationMemory creates a new conversation memory instance
func NewConversationMemory() *ConversationMemory {
	return &ConversationMemory{
		messages: make([]*types.Message, 0),
	}
}

// Add appends a message to the conversation history
func (cm *ConversationMemory) Add(msg *types.Message) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if msg != nil {
		cm.messages = append(cm.messages, msg)
	}
}

// GetAll returns a copy of all messages in the conversation history
func (cm *ConversationMemory) GetAll() []*types.Message {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// Return a copy to prevent external modification
	result := make([]*types.Message, len(cm.messages))
	copy(result, cm.messages)
	return result
}

// GetRecent returns the most recent N messages
func (cm *ConversationMemory) GetRecent(n int) []*types.Message {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if n <= 0 {
		return []*types.Message{}
	}

	if n >= len(cm.messages) {
		// Return all messages if n is greater than available
		result := make([]*types.Message, len(cm.messages))
		copy(result, cm.messages)
		return result
	}

	// Return the last n messages
	startIdx := len(cm.messages) - n
	result := make([]*types.Message, n)
	copy(result, cm.messages[startIdx:])
	return result
}

// Clear removes all messages from the conversation history
func (cm *ConversationMemory) Clear() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.messages = make([]*types.Message, 0)
}

// Count returns the number of messages in the conversation history
func (cm *ConversationMemory) Count() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return len(cm.messages)
}

// Prune reduces the conversation history to fit within a token limit
// while preserving important context (system messages, recent messages).
//
// Strategy:
// 1. Always keep system messages
// 2. Keep the most recent messages
// 3. Remove messages from the middle of the conversation
func (cm *ConversationMemory) Prune(maxTokens int) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if maxTokens <= 0 {
		return fmt.Errorf("maxTokens must be positive")
	}

	if len(cm.messages) == 0 {
		return nil
	}

	// Estimate tokens (rough approximation: 1 token ≈ 4 characters)
	estimateTokens := func(msg *types.Message) int {
		return len(msg.Content) / 4
	}

	// Separate system messages from conversation messages
	var systemMessages []*types.Message
	var conversationMessages []*types.Message

	for _, msg := range cm.messages {
		if msg.Role == types.RoleSystem {
			systemMessages = append(systemMessages, msg)
		} else {
			conversationMessages = append(conversationMessages, msg)
		}
	}

	// Calculate token budget
	systemTokens := 0
	for _, msg := range systemMessages {
		systemTokens += estimateTokens(msg)
	}

	remainingTokens := maxTokens - systemTokens
	if remainingTokens <= 0 {
		// System messages alone exceed the limit
		// Keep only the most recent system message
		if len(systemMessages) > 0 {
			cm.messages = []*types.Message{systemMessages[len(systemMessages)-1]}
		}
		return nil
	}

	// Keep as many recent conversation messages as possible
	keptMessages := make([]*types.Message, 0)
	currentTokens := 0

	// Add messages from newest to oldest
	for i := len(conversationMessages) - 1; i >= 0; i-- {
		msgTokens := estimateTokens(conversationMessages[i])
		if currentTokens+msgTokens <= remainingTokens {
			keptMessages = append([]*types.Message{conversationMessages[i]}, keptMessages...)
			currentTokens += msgTokens
		} else {
			// Can't fit any more messages
			break
		}
	}

	// Rebuild messages: system messages + kept conversation messages
	cm.messages = make([]*types.Message, 0, len(systemMessages)+len(keptMessages))
	cm.messages = append(cm.messages, systemMessages...)
	cm.messages = append(cm.messages, keptMessages...)

	return nil
}

// AddMultiple adds multiple messages at once (thread-safe)
func (cm *ConversationMemory) AddMultiple(messages []*types.Message) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for _, msg := range messages {
		if msg != nil {
			cm.messages = append(cm.messages, msg)
		}
	}
}

// GetByRole returns all messages with the specified role
func (cm *ConversationMemory) GetByRole(role types.MessageRole) []*types.Message {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	result := make([]*types.Message, 0)
	for _, msg := range cm.messages {
		if msg.Role == role {
			result = append(result, msg)
		}
	}
	return result
}
