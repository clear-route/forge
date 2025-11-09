package context

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/entrhq/forge/pkg/agent/memory"
	"github.com/entrhq/forge/pkg/llm"
	"github.com/entrhq/forge/pkg/types"
)

var thresholdDebugLog *log.Logger

func init() {
	// Create debug log file in /tmp
	f, err := os.OpenFile("/tmp/forge-threshold-strategy-debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Failed to open threshold debug log: %v", err)
		thresholdDebugLog = log.New(os.Stderr, "[THRESHOLD-DEBUG] ", log.LstdFlags|log.Lshortfile)
	} else {
		thresholdDebugLog = log.New(f, "[THRESHOLD-DEBUG] ", log.LstdFlags|log.Lshortfile)
	}
}

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
	
	// Find messages to summarize (oldest first, skip system message)
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
		if msg.Metadata != nil {
			if summarized, ok := msg.Metadata["summarized"].(bool); ok && summarized {
				continue
			}
		}
		
		// Only summarize user and assistant messages
		if msg.Role == types.RoleUser || msg.Role == types.RoleAssistant {
			toSummarize = append(toSummarize, msg)
		}
	}
	
	if len(toSummarize) == 0 {
		return 0, nil // Nothing to summarize
	}
	
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
		return 0, fmt.Errorf("failed to generate summary: %w", err)
	}
	
	// Create a single summarized message to replace the batch
	summary := types.NewAssistantMessage(response.Content).
		WithMetadata("summarized", true).
		WithMetadata("summary_count", len(toSummarize)).
		WithMetadata("summary_method", s.Name())
	
	// Remove the summarized messages and insert the summary
	// Find the index of the first message to summarize by comparing timestamps
	firstIdx := -1
	for i, msg := range messages {
		if msg.Timestamp.Equal(toSummarize[0].Timestamp) {
			firstIdx = i
			break
		}
	}
	
	if firstIdx == -1 {
		return 0, fmt.Errorf("failed to find messages to remove")
	}
	
	// Calculate how many messages we're replacing
	replaceCount := len(toSummarize)
	
	// Build new message list
	var newMessages []*types.Message
	
	// Keep messages before the summarized range
	newMessages = append(newMessages, messages[:firstIdx]...)
	
	// Add the summary
	newMessages = append(newMessages, summary)
	
	// Keep messages after the summarized range
	if firstIdx+replaceCount < len(messages) {
		newMessages = append(newMessages, messages[firstIdx+replaceCount:]...)
	}
	
	// Clear and re-add all messages
	conv.Clear()
	conv.AddMultiple(newMessages)
	
	return len(toSummarize), nil
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