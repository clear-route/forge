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
// It uses a buffering mechanism with dual trigger conditions to reduce LLM API calls.
type ToolCallSummarizationStrategy struct {
	// messagesOldThreshold is how many messages back to start considering tool calls for the buffer.
	// For example, 20 means only tool calls that are 20+ messages old enter the buffer.
	messagesOldThreshold int

	// minToolCallsToSummarize is the minimum number of tool calls in the buffer before triggering summarization.
	// This creates batching to reduce LLM API calls.
	minToolCallsToSummarize int

	// maxToolCallDistance is the maximum age (in messages) a tool call can be before forcing summarization.
	// If any tool call exceeds this distance, all buffered tool calls are summarized regardless of buffer size.
	maxToolCallDistance int

	// excludedTools is a set of tool names that should never be summarized.
	// These are typically loop-breaking tools or tools with high semantic value.
	excludedTools map[string]bool

	// eventChannel is used to emit progress events during parallel summarization
	eventChannel chan<- *types.AgentEvent
}

// NewToolCallSummarizationStrategy creates a new tool call summarization strategy with buffering.
// Parameters:
//   - messagesOldThreshold: Tool calls must be at least this many messages old to enter buffer (default: 20)
//   - minToolCallsToSummarize: Minimum buffer size before triggering summarization (default: 10)
//   - maxToolCallDistance: Maximum age before forcing summarization regardless of buffer size (default: 40)
//   - excludedTools: Optional list of tool names to exclude from summarization (default: loop-breaking tools)
func NewToolCallSummarizationStrategy(messagesOldThreshold, minToolCallsToSummarize, maxToolCallDistance int, excludedTools ...string) *ToolCallSummarizationStrategy {
	if messagesOldThreshold <= 0 {
		messagesOldThreshold = 20
	}
	if minToolCallsToSummarize <= 0 {
		minToolCallsToSummarize = 10
	}
	if maxToolCallDistance <= 0 {
		maxToolCallDistance = 40
	}

	// Build exclusion map
	exclusionMap := make(map[string]bool)
	if len(excludedTools) == 0 {
		// Default exclusions: loop-breaking tools that represent important interaction points
		exclusionMap["task_completion"] = true
		exclusionMap["ask_question"] = true
		exclusionMap["converse"] = true
	} else {
		// Use provided exclusions
		for _, toolName := range excludedTools {
			exclusionMap[toolName] = true
		}
	}

	return &ToolCallSummarizationStrategy{
		messagesOldThreshold:    messagesOldThreshold,
		minToolCallsToSummarize: minToolCallsToSummarize,
		maxToolCallDistance:     maxToolCallDistance,
		excludedTools:           exclusionMap,
		eventChannel:            nil, // Will be set by Manager
	}
}

// SetEventChannel sets the event channel for emitting progress events during summarization.
func (s *ToolCallSummarizationStrategy) SetEventChannel(eventChan chan<- *types.AgentEvent) {
	s.eventChannel = eventChan
}

// Name returns the strategy's identifier.
func (s *ToolCallSummarizationStrategy) Name() string {
	return "ToolCallSummarization"
}

// ShouldRun checks if buffered tool calls meet trigger conditions for summarization.
// Returns true if either:
// 1. Buffer trigger: Buffer contains >= minToolCallsToSummarize tool calls
// 2. Age trigger: Any tool call is >= maxToolCallDistance messages old
func (s *ToolCallSummarizationStrategy) ShouldRun(conv *memory.ConversationMemory, currentTokens, maxTokens int) bool {
	messages := conv.GetAll()
	totalMessages := len(messages)

	if totalMessages <= s.messagesOldThreshold {
		return false // Not enough message history
	}

	// Identify old messages that can enter the buffer
	oldMessages := messages[:totalMessages-s.messagesOldThreshold]

	// Count unsummarized tool calls in buffer and track oldest position
	bufferCount := 0
	oldestToolCallPosition := -1

	for i, msg := range oldMessages {
		// Skip if already summarized
		if isSummarized(msg) {
			continue
		}

		// Check if this is a tool-related message
		isToolMessage := msg.Role == types.RoleTool ||
			(msg.Role == types.RoleAssistant && containsToolCallIndicators(msg.Content))

		if isToolMessage {
			bufferCount++
			if oldestToolCallPosition == -1 {
				oldestToolCallPosition = i
			}
		}
	}

	// No tool calls to summarize
	if bufferCount == 0 {
		return false
	}

	// Buffer trigger: Check if buffer size meets minimum threshold
	if bufferCount >= s.minToolCallsToSummarize {
		return true
	}

	// Age trigger: Check if oldest tool call exceeds maximum distance
	if oldestToolCallPosition >= 0 {
		// Calculate distance from current position
		distance := totalMessages - oldestToolCallPosition
		if distance >= s.maxToolCallDistance {
			return true
		}
	}

	return false
}

// Summarize compresses buffered tool calls and their results using LLM-based summarization.
// All tool calls that are >= messagesOldThreshold old will be summarized when triggered,
// except for tools in the exclusion list.
func (s *ToolCallSummarizationStrategy) Summarize(ctx context.Context, conv *memory.ConversationMemory, llm llm.Provider) (int, error) {
	messages := conv.GetAll()
	if len(messages) <= s.messagesOldThreshold {
		return 0, nil
	}

	// Identify old messages that can be summarized
	oldMessages := messages[:len(messages)-s.messagesOldThreshold]
	recentMessages := messages[len(messages)-s.messagesOldThreshold:]

	// Group tool calls with their results for summarization, excluding certain tools
	groups := groupToolCallsAndResults(oldMessages, s.excludedTools)
	if len(groups) == 0 {
		return 0, nil // Nothing to summarize
	}

	// Summarize groups in parallel with progress tracking
	summarizedMessages, err := s.summarizeGroupsParallel(ctx, groups, llm)
	if err != nil {
		return 0, err
	}

	// Reconstruct conversation with summarized messages
	// Keep system messages and excluded tool messages, add summarized messages, then add recent messages
	newMessages := make([]*types.Message, 0)

	// Keep system messages and excluded tool call sequences from old messages
	inExcludedGroup := false
	excludedGroupMessages := make([]*types.Message, 0)

	for _, msg := range oldMessages {
		// Always keep system messages
		if msg.Role == types.RoleSystem {
			newMessages = append(newMessages, msg)
			continue
		}

		// Check if this starts an excluded tool group
		if msg.Role == types.RoleAssistant && containsToolCallIndicators(msg.Content) {
			toolName := extractToolName(msg.Content)
			if toolName != "" && s.excludedTools[toolName] {
				inExcludedGroup = true
				excludedGroupMessages = append(excludedGroupMessages, msg)
				continue
			}
		}

		// If we're in an excluded group, collect messages until we hit the tool result
		if inExcludedGroup {
			excludedGroupMessages = append(excludedGroupMessages, msg)
			if msg.Role == types.RoleTool {
				// End of excluded group - add all messages
				newMessages = append(newMessages, excludedGroupMessages...)
				excludedGroupMessages = make([]*types.Message, 0)
				inExcludedGroup = false
			}
		}
	}

	// Add any remaining excluded group messages
	if len(excludedGroupMessages) > 0 {
		newMessages = append(newMessages, excludedGroupMessages...)
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

// summarizeGroupsParallel processes multiple tool call groups concurrently,
// emitting progress events as each group completes.
func (s *ToolCallSummarizationStrategy) summarizeGroupsParallel(ctx context.Context, groups [][]*types.Message, llm llm.Provider) ([]*types.Message, error) {
	numGroups := len(groups)
	if numGroups == 0 {
		return nil, nil
	}

	// Create channels for results and errors
	type result struct {
		index       int
		message     *types.Message
		tokensSaved int
		err         error
	}
	resultChan := make(chan result, numGroups)

	// Process each group in a separate goroutine
	for i, group := range groups {
		go func(idx int, grp []*types.Message) {
			summary, err := s.summarizeGroup(ctx, grp, llm)

			// Calculate tokens saved for this group (approximate)
			tokensSaved := 0
			if summary != nil {
				// Rough estimate: original group size minus summary size
				for _, msg := range grp {
					tokensSaved += len(msg.Content) / 4 // Approximate tokens
				}
				tokensSaved -= len(summary.Content) / 4
			}

			resultChan <- result{index: idx, message: summary, tokensSaved: tokensSaved, err: err}

			// Emit progress event if event channel is available
			if s.eventChannel != nil {
				s.eventChannel <- types.NewContextSummarizationProgressEvent(
					s.Name(),
					idx+1,
					numGroups,
					tokensSaved,
				)
			}
		}(i, group)
	}

	// Collect results maintaining original order
	results := make([]*types.Message, numGroups)
	var firstError error

	for range numGroups {
		res := <-resultChan
		if res.err != nil && firstError == nil {
			firstError = res.err
		}
		if res.message != nil {
			results[res.index] = res.message
		}
	}

	if firstError != nil {
		return nil, fmt.Errorf("failed to summarize tool call group: %w", firstError)
	}

	return results, nil
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

// extractToolName extracts the tool name from a tool call XML content.
// The format is: <tool>{"tool_name": "name", ...}</tool>
// Returns empty string if no tool name is found.
func extractToolName(content string) string {
	// Look for "tool_name": "value" pattern in JSON
	start := strings.Index(content, `"tool_name"`)
	if start == -1 {
		return ""
	}

	// Find the colon after tool_name
	colonIdx := strings.Index(content[start:], ":")
	if colonIdx == -1 {
		return ""
	}
	start += colonIdx + 1

	// Skip whitespace
	for start < len(content) && (content[start] == ' ' || content[start] == '\t' || content[start] == '\n') {
		start++
	}

	// Expect opening quote
	if start >= len(content) || content[start] != '"' {
		return ""
	}
	start++ // Skip opening quote

	// Find closing quote
	end := strings.Index(content[start:], `"`)
	if end == -1 {
		return ""
	}

	return content[start : start+end]
}

// shouldSkipMessage returns true if the message should be skipped during grouping.
func shouldSkipMessage(msg *types.Message) bool {
	return isSummarized(msg) || msg.Role == types.RoleSystem
}

// isToolRelatedMessage checks if a message is related to a tool call or result.
func isToolRelatedMessage(msg *types.Message) bool {
	return msg.Role == types.RoleTool ||
		(msg.Role == types.RoleAssistant && containsToolCallIndicators(msg.Content))
}

// isExcludedToolCall checks if an assistant message contains an excluded tool call.
func isExcludedToolCall(msg *types.Message, excludedTools map[string]bool) bool {
	if msg.Role != types.RoleAssistant {
		return false
	}
	toolName := extractToolName(msg.Content)
	return toolName != "" && excludedTools[toolName]
}

// groupToolCallsAndResults groups related tool calls with their results,
// excluding tools specified in the excludedTools set.
// Returns groups of messages where each group represents a tool call sequence.
func groupToolCallsAndResults(messages []*types.Message, excludedTools map[string]bool) [][]*types.Message {
	groups := make([][]*types.Message, 0)
	currentGroup := make([]*types.Message, 0)
	skipCurrentGroup := false

	for _, msg := range messages {
		if shouldSkipMessage(msg) {
			continue
		}

		isToolMessage := isToolRelatedMessage(msg)

		if isToolMessage {
			if isExcludedToolCall(msg, excludedTools) {
				skipCurrentGroup = true
				currentGroup = make([]*types.Message, 0)
				continue
			}

			if skipCurrentGroup {
				if msg.Role == types.RoleTool {
					skipCurrentGroup = false
				}
				continue
			}

			currentGroup = append(currentGroup, msg)

			if msg.Role == types.RoleTool && len(currentGroup) > 0 {
				groups = append(groups, currentGroup)
				currentGroup = make([]*types.Message, 0)
			}
		} else if len(currentGroup) > 0 {
			groups = append(groups, currentGroup)
			currentGroup = make([]*types.Message, 0)
			skipCurrentGroup = false
		}
	}

	if len(currentGroup) > 0 {
		groups = append(groups, currentGroup)
	}

	return groups
}
