package context

import (
	"context"
	"fmt"
	"strings"

	"github.com/entrhq/forge/pkg/agent/memory"
	"github.com/entrhq/forge/pkg/llm"
	"github.com/entrhq/forge/pkg/types"
)

// ThresholdSummarizationStrategy triggers summarization when token usage
// exceeds a configured percentage of the maximum context window.
type ThresholdSummarizationStrategy struct {
	// thresholdPercent is the percentage (0-100) of max tokens that triggers summarization
	thresholdPercent float64

	// messagesPerSummary is how many messages to summarize in each batch
	messagesPerSummary int
}

// NewThresholdSummarizationStrategy creates a new threshold-based strategy.
// thresholdPercent should be between 0 and 100 (e.g., 80 for 80% of max tokens).
// messagesPerSummary controls how many messages to summarize in each batch.
func NewThresholdSummarizationStrategy(thresholdPercent float64, messagesPerSummary int) *ThresholdSummarizationStrategy {
	// Clamp threshold to valid range
	if thresholdPercent < 0 {
		thresholdPercent = 0
	}
	if thresholdPercent > 100 {
		thresholdPercent = 100
	}

	// Ensure we summarize at least 1 message
	if messagesPerSummary < 1 {
		messagesPerSummary = 1
	}

	return &ThresholdSummarizationStrategy{
		thresholdPercent:   thresholdPercent,
		messagesPerSummary: messagesPerSummary,
	}
}

// Name returns the strategy name
func (s *ThresholdSummarizationStrategy) Name() string {
	return "ThresholdSummarization"
}

// ShouldRun returns true when current token usage exceeds the threshold
func (s *ThresholdSummarizationStrategy) ShouldRun(conv *memory.ConversationMemory, currentTokens, maxTokens int) bool {
	if maxTokens <= 0 {
		return false
	}

	// Calculate current usage percentage
	usagePercent := (float64(currentTokens) / float64(maxTokens)) * 100

	// Trigger if we've exceeded the threshold
	return usagePercent >= s.thresholdPercent
}

// Summarize creates summaries for old messages to free up context space
func (s *ThresholdSummarizationStrategy) Summarize(ctx context.Context, conv *memory.ConversationMemory, llm llm.Provider) (int, error) {
	messages := conv.GetAll()
	if len(messages) == 0 {
		return 0, nil
	}

	// Collect messages to summarize
	toSummarize := s.collectMessagesToSummarize(messages)
	if len(toSummarize) == 0 {
		return 0, nil
	}

	// Generate summary using LLM
	summary, err := s.generateSummary(ctx, toSummarize, llm)
	if err != nil {
		return 0, err
	}

	// Replace summarized messages with the summary
	if err := s.replaceMessagesWithSummary(conv, messages, toSummarize, summary); err != nil {
		return 0, err
	}

	return len(toSummarize), nil
}

// collectMessagesToSummarize finds messages to summarize (oldest first, skip system message)
func (s *ThresholdSummarizationStrategy) collectMessagesToSummarize(messages []*types.Message) []*types.Message {
	var toSummarize []*types.Message
	startIdx := 0

	// Skip system message if present
	if len(messages) > 0 && messages[0].Role == types.RoleSystem {
		startIdx = 1
	}

	// Collect messages to summarize, respecting the batch size
	for i := startIdx; i < len(messages) && len(toSummarize) < s.messagesPerSummary; i++ {
		msg := messages[i]

		// Skip already summarized messages
		if s.isSummarized(msg) {
			continue
		}

		// Only summarize user and assistant messages
		if msg.Role == types.RoleUser || msg.Role == types.RoleAssistant {
			toSummarize = append(toSummarize, msg)
		}
	}

	return toSummarize
}

// isSummarized checks if a message has already been summarized
func (s *ThresholdSummarizationStrategy) isSummarized(msg *types.Message) bool {
	if msg.Metadata == nil {
		return false
	}
	summarized, ok := msg.Metadata["summarized"].(bool)
	return ok && summarized
}

// generateSummary calls the LLM to create a summary of the given messages
func (s *ThresholdSummarizationStrategy) generateSummary(ctx context.Context, toSummarize []*types.Message, llm llm.Provider) (*types.Message, error) {
	// Build prompt for summarization
	prompt := s.buildSummarizationPrompt(toSummarize)

	// Create messages for LLM
	llmMessages := []*types.Message{
		types.NewSystemMessage("You are a helpful assistant that creates concise summaries of agent conversations. You are excellent at preserving important context."),
		types.NewUserMessage(prompt),
	}

	// Call LLM to generate summary
	response, err := llm.Complete(ctx, llmMessages)
	if err != nil {
		return nil, fmt.Errorf("failed to generate summary: %w", err)
	}

	// Create a single summarized message to replace the batch
	summary := types.NewAssistantMessage(response.Content).
		WithMetadata("summarized", true).
		WithMetadata("summary_count", len(toSummarize)).
		WithMetadata("summary_method", s.Name())

	return summary, nil
}

// replaceMessagesWithSummary removes summarized messages and inserts the summary
func (s *ThresholdSummarizationStrategy) replaceMessagesWithSummary(conv *memory.ConversationMemory, messages []*types.Message, toSummarize []*types.Message, summary *types.Message) error {
	// Find the index of the first message to summarize
	firstIdx := s.findMessageIndex(messages, toSummarize[0])
	if firstIdx == -1 {
		return fmt.Errorf("failed to find messages to remove")
	}

	// Build new message list
	newMessages := s.buildNewMessageList(messages, firstIdx, len(toSummarize), summary)

	// Clear and re-add all messages
	conv.Clear()
	conv.AddMultiple(newMessages)

	return nil
}

// findMessageIndex finds the index of a message by comparing timestamps
func (s *ThresholdSummarizationStrategy) findMessageIndex(messages []*types.Message, target *types.Message) int {
	for i, msg := range messages {
		if msg.Timestamp.Equal(target.Timestamp) {
			return i
		}
	}
	return -1
}

// buildNewMessageList creates a new message list with the summary replacing the batch
func (s *ThresholdSummarizationStrategy) buildNewMessageList(messages []*types.Message, firstIdx, replaceCount int, summary *types.Message) []*types.Message {
	var newMessages []*types.Message

	// Keep messages before the summarized range
	newMessages = append(newMessages, messages[:firstIdx]...)

	// Add the summary
	newMessages = append(newMessages, summary)

	// Keep messages after the summarized range
	if firstIdx+replaceCount < len(messages) {
		newMessages = append(newMessages, messages[firstIdx+replaceCount:]...)
	}

	return newMessages
}

// buildSummarizationPrompt creates a prompt for summarizing a batch of messages
func (s *ThresholdSummarizationStrategy) buildSummarizationPrompt(messages []*types.Message) string {
	var b strings.Builder

	b.WriteString("Please create a concise summary of the following conversation messages.\n")
	b.WriteString("Preserve key information, decisions, and context. Be brief but comprehensive.\n\n")
	b.WriteString("Messages to summarize:\n\n")

	for i, msg := range messages {
		b.WriteString(fmt.Sprintf("%d. %s: %s\n\n", i+1, msg.Role, msg.Content))
	}

	b.WriteString("Please provide a concise summary:")

	return b.String()
}
