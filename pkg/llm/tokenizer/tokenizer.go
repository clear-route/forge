// Package tokenizer provides client-side token counting for LLM messages.
// This enables token tracking for any LLM provider, not just those that
// report usage in their API responses.
package tokenizer

import (
	"fmt"
	"sync"

	"github.com/entrhq/forge/pkg/types"
	tiktoken "github.com/pkoukk/tiktoken-go"
)

// Tokenizer provides token counting functionality
type Tokenizer struct {
	encoding *tiktoken.Tiktoken
	mu       sync.Mutex
}

// defaultEncoding is the encoding used for most modern models (GPT-4, Claude, etc.)
const defaultEncoding = "cl100k_base"

// New creates a new Tokenizer instance
func New() (*Tokenizer, error) {
	encoding, err := tiktoken.GetEncoding(defaultEncoding)
	if err != nil {
		return nil, fmt.Errorf("failed to get tiktoken encoding: %w", err)
	}

	return &Tokenizer{
		encoding: encoding,
	}, nil
}

// CountTokens counts the number of tokens in the given text
func (t *Tokenizer) CountTokens(text string) int {
	if text == "" {
		return 0
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	tokens := t.encoding.Encode(text, nil, nil)
	return len(tokens)
}

// CountMessageTokens counts tokens for a message with role overhead
// Different models have different formatting, but this provides a reasonable estimate
func (t *Tokenizer) CountMessageTokens(message *types.Message) int {
	if message == nil {
		return 0
	}

	// Base overhead per message (role markers, formatting, etc.)
	// This is an approximation based on OpenAI's tokenization
	tokensPerMessage := 3

	count := tokensPerMessage
	count += t.CountTokens(string(message.Role))
	count += t.CountTokens(message.Content)

	return count
}

// CountMessagesTokens counts total tokens for a slice of messages
func (t *Tokenizer) CountMessagesTokens(messages []*types.Message) int {
	total := 0
	for _, msg := range messages {
		total += t.CountMessageTokens(msg)
	}
	// Add 3 tokens for reply priming (only if there are messages)
	if len(messages) > 0 {
		total += 3
	}
	return total
}
