package context

import (
	"context"
	"fmt"
	"strings"

	"github.com/entrhq/forge/pkg/agent/memory"
	"github.com/entrhq/forge/pkg/llm"
	"github.com/entrhq/forge/pkg/types"
)

// ToolCallSummarizationStrategy summarizes old tool calls and their results
// to reduce context size while preserving semantic meaning.
type ToolCallSummarizationStrategy struct {
	// messagesOldThreshold is how many messages back to start summarizing tool calls.
	// For example, 20 means only summarize tool calls that are 20+ messages old.
	messagesOldThreshold int
}

// NewToolCallSummarizationStrategy creates a new tool call summarization strategy.
// messagesOldThreshold specifies how many messages old a tool call must be before summarization (default: 20).
func NewToolCallSummarizationStrategy(messagesOldThreshold int) *ToolCallSummarizationStrategy {
	if messagesOldThreshold <= 0 {
		messagesOldThreshold = 20 // Default to 20 messages
	}
	return &ToolCallSummarizationStrategy{
		messagesOldThreshold: messagesOldThreshold,
	}
}

// Name returns the strategy's identifier.
func (s *ToolCallSummarizationStrategy) Name() string {
	return "ToolCallSummarization"
}

// ShouldRun checks if there are old tool calls that need summarization.
func (s *ToolCallSummarizationStrategy) ShouldRun(conv *memory.ConversationMemory, currentTokens, maxTokens int) bool {
	messages := conv.GetAll()
	if len(messages) <= s.messagesOldThreshold {
		return false // Not enough message history
	}

	// Check if there are any tool calls/results in the old messages that haven't been summarized
	oldMessages := messages[:len(messages)-s.messagesOldThreshold]
	for _, msg := range oldMessages {
		// Skip if already summarized
		if isSummarized(msg) {
			continue
		}

		// Check if this is a tool result or assistant message with tool call
		if msg.Role == types.RoleTool {
			return true // Found unsummarized tool result
		}

		// Check if assistant message contains tool call indicators
		if msg.Role == types.RoleAssistant && containsToolCallIndicators(msg.Content) {
			return true // Found unsummarized tool call
		}
	}

	return false
}

// Summarize compresses old tool calls and their results using LLM-based summarization.
func (s *ToolCallSummarizationStrategy) Summarize(ctx context.Context, conv *memory.ConversationMemory, llm llm.Provider) (int, error) {
	messages := conv.GetAll()
	if len(messages) <= s.messagesOldThreshold {
		return 0, nil
	}

	// Identify old messages that can be summarized
	oldMessages := messages[:len(messages)-s.messagesOldThreshold]
	recentMessages := messages[len(messages)-s.messagesOldThreshold:]

	// Group tool calls with their results for summarization
	groups := groupToolCallsAndResults(oldMessages)
	if len(groups) == 0 {
		return 0, nil // Nothing to summarize
	}

	// Summarize each group
	summarizedMessages := make([]*types.Message, 0)
	for _, group := range groups {
		summary, err := s.summarizeGroup(ctx, group, llm)
		if err != nil {
			return 0, fmt.Errorf("failed to summarize tool call group: %w", err)
		}
		summarizedMessages = append(summarizedMessages, summary)
	}

	// Reconstruct conversation with summarized messages
	// Keep system messages, add summarized messages, then add recent messages
	newMessages := make([]*types.Message, 0)

	// Keep all system messages from old messages
	for _, msg := range oldMessages {
		if msg.Role == types.RoleSystem {
			newMessages = append(newMessages, msg)
		}
	}

	// Add summarized messages
	newMessages = append(newMessages, summarizedMessages...)

	// Add recent messages unchanged
	newMessages = append(newMessages, recentMessages...)

	// Replace conversation messages
	conv.Clear()
	for _, msg := range newMessages {
		conv.Add(msg)
	}

	// Return the number of groups processed (items summarized)
	return len(groups), nil
}

// summarizeGroup creates a concise summary of a tool call and its result using the LLM.
func (s *ToolCallSummarizationStrategy) summarizeGroup(ctx context.Context, group []*types.Message, llm llm.Provider) (*types.Message, error) {
	// Build context for summarization
	var builder strings.Builder
	for _, msg := range group {
		builder.WriteString(fmt.Sprintf("[%s]: %s\n\n", msg.Role, msg.Content))
	}

	// Create summarization prompt
	prompt := fmt.Sprintf(`You are summarizing an old tool call and its result to compress context. Provide a concise 2-3 sentence summary that captures:

1. What tool was called and why
2. Key input parameters or actions
3. Result summary (success/failure, key outcomes)

Original messages:
%s

Provide only the summary, no additional commentary:`, builder.String())

	// Call LLM for summarization
	messages := []*types.Message{
		types.NewSystemMessage("You are a helpful assistant that summarizes tool calls concisely."),
		types.NewUserMessage(prompt),
	}

	response, err := llm.Complete(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("LLM summarization failed: %w", err)
	}

	// Create summarized message
	summary := types.NewAssistantMessage(fmt.Sprintf("[SUMMARIZED] %s", response.Content))
	summary.WithMetadata("summarized", true)
	summary.WithMetadata("original_message_count", len(group))

	return summary, nil
}

// Helper functions

// isSummarized checks if a message has already been summarized.
func isSummarized(msg *types.Message) bool {
	if msg.Metadata == nil {
		return false
	}
	summarized, ok := msg.Metadata["summarized"].(bool)
	return ok && summarized
}

// containsToolCallIndicators checks if the message content contains tool call XML tags.
func containsToolCallIndicators(content string) bool {
	return strings.Contains(content, "<tool>") && strings.Contains(content, "</tool>")
}

// groupToolCallsAndResults groups related tool calls with their results.
// Returns groups of messages where each group represents a tool call sequence.
func groupToolCallsAndResults(messages []*types.Message) [][]*types.Message {
	groups := make([][]*types.Message, 0)
	currentGroup := make([]*types.Message, 0)

	for _, msg := range messages {
		// Skip already summarized messages
		if isSummarized(msg) {
			continue
		}

		// Skip system messages (they're kept separately)
		if msg.Role == types.RoleSystem {
			continue
		}

		// Add to current group if it's a tool-related message
		if msg.Role == types.RoleTool ||
			(msg.Role == types.RoleAssistant && containsToolCallIndicators(msg.Content)) {
			currentGroup = append(currentGroup, msg)

			// Complete the group when we hit a tool result
			if msg.Role == types.RoleTool && len(currentGroup) > 0 {
				groups = append(groups, currentGroup)
				currentGroup = make([]*types.Message, 0)
			}
		} else if len(currentGroup) > 0 {
			// Non-tool message ends the current group
			groups = append(groups, currentGroup)
			currentGroup = make([]*types.Message, 0)
		}
	}

	// Add any remaining group
	if len(currentGroup) > 0 {
		groups = append(groups, currentGroup)
	}

	return groups
}
