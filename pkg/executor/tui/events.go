package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/entrhq/forge/pkg/agent/tools"
	"github.com/entrhq/forge/pkg/types"
)

// handleAgentEvent processes events from the agent event stream.
// This is the main event handler that updates the UI based on agent activity.
//
//nolint:gocyclo
func (m *model) handleAgentEvent(event *types.AgentEvent) {
	debugLog.Printf("handleAgentEvent called with event type: %s", event.Type)

	switch event.Type {
	case types.EventTypeThinkingStart:
		debugLog.Printf("Processing EventTypeThinkingStart")
		m.handleThinkingStart()

	case types.EventTypeThinkingContent:
		debugLog.Printf("Processing EventTypeThinkingContent: %s", event.Content)
		m.handleThinkingContent(event)
		return // Exit early to preserve streaming viewport update

	case types.EventTypeThinkingEnd:
		debugLog.Printf("Processing EventTypeThinkingEnd")
		m.handleThinkingEnd()

	case types.EventTypeToolCallStart:
		debugLog.Printf("Processing EventTypeToolCallStart")
		m.handleToolCallStart(event)

	case types.EventTypeToolCall:
		debugLog.Printf("Processing EventTypeToolCall")
		m.handleToolCall(event)

	case types.EventTypeToolResult:
		debugLog.Printf("Processing EventTypeToolResult")
		m.handleToolResult(event)

	case types.EventTypeMessageStart:
		debugLog.Printf("Processing EventTypeMessageStart")
		m.handleMessageStart()

	case types.EventTypeMessageContent:
		debugLog.Printf("Processing EventTypeMessageContent: %s", event.Content)
		if m.handleMessageContent(event.Content) {
			return // Exit early to preserve streaming viewport update
		}

	case types.EventTypeMessageEnd:
		debugLog.Printf("Processing EventTypeMessageEnd")
		m.handleMessageEnd()

	case types.EventTypeError:
		debugLog.Printf("Processing EventTypeError: %v", event.Error)
		m.handleError(event)

	case types.EventTypeTurnEnd:
		debugLog.Printf("Processing EventTypeTurnEnd")
		m.handleTurnEnd()

	case types.EventTypeUpdateBusy:
		debugLog.Printf("Processing EventTypeUpdateBusy")
		m.handleUpdateBusy(event)

	case types.EventTypeToolApprovalRequest:
		debugLog.Printf("Processing EventTypeToolApprovalRequest")
		m.handleToolApprovalRequest(event)

	case types.EventTypeToolApprovalGranted:
		debugLog.Printf("Processing EventTypeToolApprovalGranted")
		m.handleToolApprovalGranted()

	case types.EventTypeToolApprovalRejected:
		debugLog.Printf("Processing EventTypeToolApprovalRejected")
		m.handleToolApprovalRejected()

	case types.EventTypeToolApprovalTimeout:
		m.handleToolApprovalTimeout()

	case types.EventTypeApiCallStart:
		m.handleApiCallStart(event)

	case types.EventTypeTokenUsage:
		m.handleTokenUsage(event)

	case types.EventTypeCommandExecutionStart:
		m.handleCommandExecutionStart(event)

	case types.EventTypeCommandOutput:
		m.handleCommandExecutionOutput(event)

	case types.EventTypeCommandExecutionComplete:
		m.handleCommandExecutionComplete(event)

	case types.EventTypeContextSummarizationStart:
		m.handleContextSummarizationStart(event)

	case types.EventTypeContextSummarizationProgress:
		m.handleContextSummarizationProgress(event)

	case types.EventTypeContextSummarizationComplete:
		m.handleContextSummarizationComplete(event)
	}

	// Update viewport with current content
	m.viewport.SetContent(m.content.String())
	m.viewport.GotoBottom()
}

// Thinking event handlers

func (m *model) handleThinkingStart() {
	m.isThinking = true
	m.thinkingBuffer.Reset()
}

func (m *model) handleThinkingContent(event *types.AgentEvent) {
	if event.Content == "" {
		return
	}
	// Buffer the thinking content
	m.thinkingBuffer.WriteString(event.Content)
	// Stream with "Thinking" label, content follows immediately
	header := "💭 Thinking "
	formatted := formatEntry("", m.thinkingBuffer.String(), thinkingStyle, m.width, false)
	m.viewport.SetContent(m.content.String() + header + formatted)
	m.viewport.GotoBottom()
}

func (m *model) handleThinkingEnd() {
	if m.thinkingBuffer.Len() > 0 {
		header := "💭 Thinking "
		formatted := formatEntry("", m.thinkingBuffer.String(), thinkingStyle, m.width, false)
		m.content.WriteString(header + formatted)
	}
	m.content.WriteString("\n\n")
	m.isThinking = false
	m.thinkingBuffer.Reset()
}

// Tool event handlers

func (m *model) handleToolCallStart(event *types.AgentEvent) {
	// Check if we have early tool name detection in metadata
	if toolName, ok := event.Metadata["tool_name"].(string); ok && toolName != "" && !m.toolNameDisplayed {
		// Display the tool name immediately when detected early
		formatted := formatEntry("🔧 ", toolName, toolStyle, m.width, false)
		m.content.WriteString(formatted)
		m.content.WriteString("\n")
		m.viewport.SetContent(m.content.String())
		m.viewport.GotoBottom()
		m.toolNameDisplayed = true
	}
	// If no tool name yet, we'll wait for EventTypeToolCall
}

func (m *model) handleToolCall(event *types.AgentEvent) {
	// Only display if we haven't already shown it from early detection
	if !m.toolNameDisplayed {
		formatted := formatEntry("🔧 ", event.ToolName, toolStyle, m.width, false)
		m.content.WriteString(formatted)
		m.content.WriteString("\n")
	}
	// Track tool call for result display
	m.lastToolName = event.ToolName
	// Generate a simple cache key using timestamp + tool name
	m.lastToolCallID = fmt.Sprintf("%d_%s", time.Now().UnixNano(), event.ToolName)
	m.toolNameDisplayed = false // Reset for next tool call
}

func (m *model) handleToolResult(event *types.AgentEvent) {
	resultStr := fmt.Sprintf("%v", event.ToolOutput)

	// Classify the tool result to determine display strategy
	tier := m.resultClassifier.ClassifyToolResult(m.lastToolName, resultStr)

	switch tier {
	case TierFullInline:
		// Display full result inline (loop-breaking tools)
		formatted := formatEntry("    ✓ ", resultStr, toolResultStyle, m.width, false)
		m.content.WriteString(formatted)

	case TierSummaryWithPreview:
		// Display summary + preview lines
		summary := m.resultSummarizer.GenerateSummary(m.lastToolName, resultStr)
		preview := m.resultClassifier.GetPreviewLines(resultStr)
		displayText := summary + "\n" + preview
		formatted := formatEntry("    ✓ ", displayText, toolResultStyle, m.width, false)
		m.content.WriteString(formatted)
		// Cache the full result for viewing
		m.resultCache.store(m.lastToolCallID, m.lastToolName, resultStr, summary)

	case TierSummaryOnly:
		// Display summary only
		summary := m.resultSummarizer.GenerateSummary(m.lastToolName, resultStr)
		formatted := formatEntry("    ✓ ", summary, toolResultStyle, m.width, false)
		m.content.WriteString(formatted)
		// Cache the full result for viewing
		m.resultCache.store(m.lastToolCallID, m.lastToolName, resultStr, summary)

	case TierOverlayOnly:
		// Command execution already handled by overlay system
		// Don't display anything inline
	}

	m.content.WriteString("\n\n")
}

// Message event handlers

func (m *model) handleMessageStart() {
	m.messageBuffer.Reset()
}

func (m *model) handleMessageContent(content string) bool {
	if strings.TrimSpace(content) != "" && !m.hasMessageContentStarted {
		m.hasMessageContentStarted = true
	}

	// Buffer the message content
	m.messageBuffer.WriteString(content)

	// Stream message content as it arrives
	formatted := formatEntry("", m.messageBuffer.String(), lipgloss.NewStyle(), m.width, false)
	m.viewport.SetContent(m.content.String() + formatted)
	m.viewport.GotoBottom()

	return true
}

func (m *model) handleMessageEnd() {
	// Finalize message content (like thinking does)
	if m.messageBuffer.Len() > 0 && m.hasMessageContentStarted {
		formatted := formatEntry("", m.messageBuffer.String(), lipgloss.NewStyle(), m.width, false)
		m.content.WriteString(formatted)
		m.content.WriteString("\n\n")
		m.hasMessageContentStarted = false
	}
	m.messageBuffer.Reset()
}

// Error and state handlers

func (m *model) handleError(event *types.AgentEvent) {
	m.content.WriteString(errorStyle.Render(fmt.Sprintf("  ❌ Error: %v", event.Error)))
	m.content.WriteString("\n\n")
}

func (m *model) handleTurnEnd() {
	// Turn end - clear busy state
	m.agentBusy = false
	m.recalculateLayout()
}

func (m *model) handleUpdateBusy(event *types.AgentEvent) {
	// Update busy state based on event
	wasBusy := m.agentBusy
	m.agentBusy = event.IsBusy
	if m.agentBusy {
		// Pick a random loading message when becoming busy
		m.currentLoadingMessage = getRandomLoadingMessage()
	}
	// Recalculate layout if busy state changed
	if wasBusy != m.agentBusy {
		m.recalculateLayout()
	}
}

// Tool approval handlers

func (m *model) handleToolApprovalRequest(event *types.AgentEvent) {
	// Show "Requesting approval" message before overlay
	formatted := formatEntry("  ⏳ ", "Requesting tool approval...", toolStyle, m.width, false)
	m.content.WriteString(formatted)
	m.content.WriteString("\n")
	m.viewport.SetContent(m.content.String())
	m.viewport.GotoBottom()

	// Handle tool approval request by showing overlay
	if event.Preview != nil {
		preview, ok := event.Preview.(*tools.ToolPreview)
		if ok {
			// Create response callback that will be called by the overlay
			responseFunc := func(response *types.ApprovalResponse) {
				// Send approval response to agent
				m.channels.Approval <- response

				// Close overlay and update viewport
				m.overlay.deactivate()
				m.viewport.SetContent(m.content.String())
				m.viewport.GotoBottom()
			}

			// Create and activate diff viewer overlay
			diffViewer := NewDiffViewer(
				event.ApprovalID,
				event.ToolName,
				preview,
				m.width,
				m.height,
				responseFunc,
			)
			m.overlay.activate(OverlayModeDiffViewer, diffViewer)
		}
	}
}

func (m *model) handleToolApprovalGranted() {
	// Approval granted - show confirmation
	formatted := formatEntry("  ✓ ", "Tool approved - executing...", toolStyle, m.width, false)
	m.content.WriteString(formatted)
	m.content.WriteString("\n")
}

func (m *model) handleToolApprovalRejected() {
	// Approval rejected - log it
	formatted := formatEntry("  ✗ ", "Tool rejected by user", errorStyle, m.width, false)
	m.content.WriteString(formatted)
	m.content.WriteString("\n")
}

func (m *model) handleToolApprovalTimeout() {
	// Approval timeout - log it
	formatted := formatEntry("  ⏱ ", "Tool approval timed out", errorStyle, m.width, false)
	m.content.WriteString(formatted)
	m.content.WriteString("\n")
}

// API and token handlers

func (m *model) handleApiCallStart(event *types.AgentEvent) {
	// Update context token information
	if event.ApiCallInfo != nil {
		m.currentContextTokens = event.ApiCallInfo.ContextTokens
		m.maxContextTokens = event.ApiCallInfo.MaxContextTokens
	}
}

func (m *model) handleTokenUsage(event *types.AgentEvent) {
	// Update token usage counts
	if event.TokenUsage != nil {
		m.totalPromptTokens += event.TokenUsage.PromptTokens
		m.totalCompletionTokens += event.TokenUsage.CompletionTokens
		m.totalTokens += event.TokenUsage.TotalTokens
	}
}

// Command execution handlers

func (m *model) handleCommandExecutionStart(event *types.AgentEvent) {
	// Show command execution started message
	if event.CommandExecution != nil {
		formatted := formatEntry("  🚀 ", fmt.Sprintf("Executing: %s", event.CommandExecution.Command), toolStyle, m.width, false)
		m.content.WriteString(formatted)
		m.content.WriteString("\n")
		m.viewport.SetContent(m.content.String())
		m.viewport.GotoBottom()

		// Create and activate command execution overlay
		overlay := NewCommandExecutionOverlay(
			event.CommandExecution.Command,
			event.CommandExecution.WorkingDir,
			event.CommandExecution.ExecutionID,
			m.channels.Cancel,
		)
		m.overlay.activate(OverlayModeCommandOutput, overlay)
	}
}

func (m *model) handleCommandExecutionOutput(event *types.AgentEvent) {
	// Stream command output as it arrives
	// Write output directly without styling to preserve formatting/indentation
	if event.CommandExecution != nil && event.CommandExecution.Output != "" {
		m.content.WriteString(event.CommandExecution.Output)
	}
}

func (m *model) handleCommandExecutionComplete(event *types.AgentEvent) {
	// Show command completion status
	if event.CommandExecution != nil {
		if event.CommandExecution.ExitCode == 0 {
			formatted := formatEntry("  ✓ ", "Command completed successfully", toolStyle, m.width, false)
			m.content.WriteString(formatted)
		} else {
			formatted := formatEntry("  ✗ ", fmt.Sprintf("Command failed with exit code %d", event.CommandExecution.ExitCode), errorStyle, m.width, false)
			m.content.WriteString(formatted)
		}
		m.content.WriteString("\n")
	}
}

// Context summarization handlers

func (m *model) handleContextSummarizationStart(event *types.AgentEvent) {
	m.summarization.active = true
	m.summarization.startTime = time.Now()
	if event.ContextSummarization != nil {
		m.summarization.strategy = event.ContextSummarization.Strategy
		m.summarization.currentTokens = event.ContextSummarization.CurrentTokens
		m.summarization.maxTokens = event.ContextSummarization.MaxTokens
		m.summarization.totalItems = event.ContextSummarization.TotalItems
	}
}

func (m *model) handleContextSummarizationProgress(event *types.AgentEvent) {
	if event.ContextSummarization != nil {
		m.summarization.itemsProcessed = event.ContextSummarization.ItemsProcessed
		// Calculate progress percentage from items processed
		if event.ContextSummarization.TotalItems > 0 {
			m.summarization.progressPercent = float64(event.ContextSummarization.ItemsProcessed) / float64(event.ContextSummarization.TotalItems) * 100
		}
	}
}

func (m *model) handleContextSummarizationComplete(event *types.AgentEvent) {
	if event.ContextSummarization != nil {
		oldTokens := m.summarization.currentTokens
		newTokens := event.ContextSummarization.NewTokenCount

		m.summarization.active = false
		duration := time.Since(m.summarization.startTime).Seconds()

		m.showToast(
			"✨ Context optimized",
			fmt.Sprintf("Reduced from %s to %s tokens (%.1fs)",
				formatTokenCount(oldTokens),
				formatTokenCount(newTokens),
				duration),
			"🧠",
			false,
		)

		// Update current context tokens
		m.currentContextTokens = newTokens
	}
}
