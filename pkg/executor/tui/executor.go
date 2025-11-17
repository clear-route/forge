// Package tui provides a terminal user interface executor for Forge agents,
// offering an interactive, Gemini-style interface for conversations.
package tui

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/entrhq/forge/pkg/agent"
	"github.com/entrhq/forge/pkg/agent/git"
	"github.com/entrhq/forge/pkg/agent/slash"
	"github.com/entrhq/forge/pkg/agent/tools"
	"github.com/entrhq/forge/pkg/config"
	"github.com/entrhq/forge/pkg/llm"
	"github.com/entrhq/forge/pkg/types"
)

// Colors and styles are now defined in styles.go for consistency across all TUI components

// loadingMessages contains quirky ASCII-style messages shown during agent processing
var loadingMessages = []string{
	"Thinking deep thoughts...",
	"Brewing some code...",
	"Analyzing the situation...",
	"Working some magic...",
	"Processing neural pathways...",
	"Launching solution engines...",
	"Forging your response...",
	"Crafting the perfect answer...",
	"Consulting the documentation...",
	"Gathering the bits and bytes...",
	"Spinning up the hamster wheel...",
	"Channeling the AI spirits...",
	"Compiling brilliance...",
	"Running through possibilities...",
	"Calculating optimal outcomes...",
}

// getRandomLoadingMessage returns a random loading message
func getRandomLoadingMessage() string {
	// #nosec G404 - Using math/rand for UI randomness is acceptable
	return loadingMessages[rand.Intn(len(loadingMessages))]
}

// Executor is a TUI-based executor that provides an interactive,
// Gemini-style interface for agent interaction.
type Executor struct {
	agent        agent.Agent
	program      *tea.Program
	provider     llm.Provider
	workspaceDir string
}

// NewExecutor creates a new TUI executor for the given agent.
func NewExecutor(agent agent.Agent, provider llm.Provider, workspaceDir string) *Executor {
	return &Executor{
		agent:        agent,
		provider:     provider,
		workspaceDir: workspaceDir,
	}
}

// Run starts the TUI executor and blocks until the user exits.
func (e *Executor) Run(ctx context.Context) error {
	// Start the agent first
	if err := e.agent.Start(ctx); err != nil {
		return fmt.Errorf("failed to start agent: %w", err)
	}

	// Discover tools from agent and populate config
	if err := config.DiscoverToolsFromAgent(e.agent); err != nil {
		// Log error but don't fail - config system is optional
		log.Printf("Warning: failed to discover tools from agent: %v", err)
	}

	model := initialModel()
	model.agent = e.agent
	model.channels = e.agent.GetChannels()
	model.workspaceDir = e.workspaceDir

	// Initialize slash handler for git operations
	if e.provider != nil && e.workspaceDir != "" {
		llmClient := newLLMAdapter(e.provider)
		tracker := git.NewModificationTracker()
		model.commitGen = git.NewCommitMessageGenerator(llmClient)
		model.prGen = git.NewPRGenerator(llmClient)
		model.slashHandler = slash.NewHandler(e.workspaceDir, tracker, model.commitGen, model.prGen)
	}

	e.program = tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	go func() {
		// Listen for agent events and forward them to the TUI
		for event := range model.channels.Event {
			e.program.Send(event)
		}
	}()

	if _, err := e.program.Run(); err != nil {
		return fmt.Errorf("failed to run TUI program: %w", err)
	}

	return nil
}

type agentErrMsg struct{ err error }

// slashCommandCompleteMsg signals that a slash command has completed
type slashCommandCompleteMsg struct{}

// operationStartMsg signals that a long-running operation has started
type operationStartMsg struct {
	message string // Loading message to display
}

// operationCompleteMsg signals that a long-running operation has completed
type operationCompleteMsg struct {
	result       string
	err          error
	successTitle string
	successIcon  string
	errorTitle   string
	errorIcon    string
}

// summarizationStatus tracks an active context summarization operation
type summarizationStatus struct {
	active          bool
	strategy        string
	currentTokens   int
	maxTokens       int
	itemsProcessed  int
	totalItems      int
	currentItem     string
	progressPercent float64
	startTime       time.Time
}

// toastNotification represents a temporary notification message
type toastNotification struct {
	active    bool
	message   string
	details   string
	icon      string
	isError   bool
	showUntil time.Time
}

// model represents the state of the TUI application.
type model struct {
	viewport                 viewport.Model
	textarea                 textarea.Model
	agent                    agent.Agent
	channels                 *types.AgentChannels
	slashHandler             *slash.Handler
	workspaceDir             string
	commitGen                *git.CommitMessageGenerator
	prGen                    *git.PRGenerator
	content                  *strings.Builder
	thinkingBuffer           *strings.Builder
	messageBuffer            *strings.Builder
	overlay                  *overlayState
	commandPalette           *CommandPalette
	summarization            *summarizationStatus
	toast                    *toastNotification
	spinner                  spinner.Model
	isThinking               bool
	agentBusy                bool
	currentLoadingMessage    string
	width                    int
	height                   int
	ready                    bool
	hasMessageContentStarted bool

	// Token usage tracking
	totalPromptTokens     int // Cumulative input tokens across all API calls
	totalCompletionTokens int // Cumulative output tokens across all API calls
	totalTokens           int // Cumulative total tokens (input + output)
	currentContextTokens  int // Current conversation context size
	maxContextTokens      int // Maximum allowed context size
}

// initialModel returns the initial state of the TUI.
func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Type your message..."
	ta.Focus()
	ta.Prompt = "> "
	ta.CharLimit = 0
	ta.SetHeight(1)
	ta.MaxHeight = 10 // Allow up to 10 lines
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false) // Disable default Enter behavior
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.FocusedStyle.Prompt = lipgloss.NewStyle().Foreground(salmonPink)
	ta.FocusedStyle.Text = lipgloss.NewStyle().Foreground(brightWhite)

	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle().Padding(0, 2)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(salmonPink)

	return model{
		viewport:       vp,
		textarea:       ta,
		content:        &strings.Builder{},
		thinkingBuffer: &strings.Builder{},
		messageBuffer:  &strings.Builder{},
		overlay:        newOverlayState(),
		commandPalette: newCommandPalette(),
		summarization:  &summarizationStatus{},
		toast:          &toastNotification{},
		spinner:        s,
		agentBusy:      false,
	}
}

// Init is the first function that will be called.
func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, m.spinner.Tick)
}

// formatTokenCount formats a token count with K/M suffixes for readability
func formatTokenCount(count int) string {
	if count >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(count)/1000000)
	}
	if count >= 1000 {
		return fmt.Sprintf("%.1fK", float64(count)/1000)
	}
	return fmt.Sprintf("%d", count)
}

// formatEntry formats and wraps an entry with icon/emoji and text
// It ensures the icon stays with the text and applies proper styling
// If iconOnly is true, only the icon is styled, text remains white
func formatEntry(icon string, text string, style lipgloss.Style, width int, iconOnly bool) string {
	// Calculate wrap width (full width minus small padding)
	wrapWidth := width - 4
	if wrapWidth <= 0 {
		wrapWidth = 80
	}

	if iconOnly {
		// Style only the icon, keep text white
		styledIcon := style.Render(icon)
		fullText := icon + text // Use unstyled for wrapping calculation
		wrapped := wordWrap(fullText, wrapWidth)

		// Replace the unstyled icon with styled icon in first occurrence
		wrapped = strings.Replace(wrapped, icon, styledIcon, 1)
		return wrapped
	}

	// Style everything (default behavior)
	fullText := icon + text
	wrapped := wordWrap(fullText, wrapWidth)

	// Apply ONLY color using inline rendering to avoid block formatting
	styledLines := make([]string, 0)
	for line := range strings.SplitSeq(wrapped, "\n") {
		styledLines = append(styledLines, style.Render(line))
	}

	return strings.Join(styledLines, "\n")
}

// wordWrap manually wraps text at word boundaries without adding any padding
// It also handles long strings without spaces by breaking them at character boundaries
func wordWrap(text string, width int) string {
	if width <= 0 {
		width = 80
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var result strings.Builder
	currentLine := ""

	for _, word := range words {
		// If a single word is longer than width, break it up
		if len(word) > width {
			// First, flush current line if it has content
			if currentLine != "" {
				result.WriteString(currentLine)
				result.WriteString("\n")
				currentLine = ""
			}

			// Break the long word into chunks
			for len(word) > 0 {
				chunkSize := width
				if len(word) < chunkSize {
					chunkSize = len(word)
				}
				result.WriteString(word[:chunkSize])
				result.WriteString("\n")
				word = word[chunkSize:]
			}
			continue
		}

		// Check if adding this word would exceed width
		switch {
		case currentLine == "":
			currentLine = word
		case len(currentLine)+1+len(word) > width:
			// Write current line and start new one
			result.WriteString(currentLine)
			result.WriteString("\n")
			currentLine = word
		default:
			// Add word to current line
			currentLine += " " + word
		}
	}

	// Write final line if there's content
	if currentLine != "" {
		result.WriteString(currentLine)
	}

	return result.String()
}

// renderSummarizationStatus renders the context summarization status overlay
func (m model) renderSummarizationStatus() string {
	if !m.summarization.active {
		return ""
	}

	// Create box with border
	boxWidth := m.width - 4
	if boxWidth < 40 {
		boxWidth = 40
	}

	var content strings.Builder

	// Header line with brain icon and message
	header := fmt.Sprintf("🧠 Optimizing context... [%s]", m.summarization.strategy)
	content.WriteString(header)
	content.WriteString("\n")

	// Progress bar
	barWidth := boxWidth - 10 // Leave room for percentage
	if barWidth < 20 {
		barWidth = 20
	}

	filledWidth := int(float64(barWidth) * m.summarization.progressPercent / 100.0)
	if filledWidth > barWidth {
		filledWidth = barWidth
	}

	bar := strings.Repeat("━", filledWidth) + strings.Repeat("━", barWidth-filledWidth)
	progressLine := fmt.Sprintf("%s %.0f%%", bar, m.summarization.progressPercent)
	content.WriteString(progressLine)
	content.WriteString("\n")

	// Current item description
	if m.summarization.currentItem != "" {
		content.WriteString(m.summarization.currentItem)
	} else if m.summarization.totalItems > 0 {
		content.WriteString(fmt.Sprintf("Processing item %d of %d...",
			m.summarization.itemsProcessed, m.summarization.totalItems))
	}

	// Create styled box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(salmonPink).
		Padding(0, 1).
		Width(boxWidth)

	return "\n" + boxStyle.Render(content.String()) + "\n"
}

// showToast displays a toast notification to the user
func (m *model) showToast(message, details, icon string, isError bool) {
	m.toast.active = true
	m.toast.message = message
	m.toast.details = details
	m.toast.icon = icon
	m.toast.isError = isError
	m.toast.showUntil = time.Now().Add(3 * time.Second)
}

// renderToast renders a toast notification
func (m model) renderToast() string {
	if !m.toast.active || time.Now().After(m.toast.showUntil) {
		return ""
	}

	// Create box with border
	boxWidth := m.width - 4
	if boxWidth < 40 {
		boxWidth = 40
	}

	var content strings.Builder

	// Icon and message
	header := fmt.Sprintf("%s %s", m.toast.icon, m.toast.message)
	content.WriteString(header)
	content.WriteString("\n")

	// Details
	if m.toast.details != "" {
		content.WriteString(m.toast.details)
	}

	// Choose border color based on error status
	borderColor := mintGreen // Green for success
	if m.toast.isError {
		borderColor = salmonPink // Pink/red for errors
	}

	// Create styled box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1).
		Width(boxWidth)

	return "\n" + boxStyle.Render(content.String()) + "\n"
}

// handleMessageContent processes message content events with streaming support
func (m *model) handleMessageContent(content string) bool {
	if content == "" {
		return false
	}

	// Buffer the content (like thinking does)
	m.messageBuffer.WriteString(content)

	// Stream message content as it arrives (like thinking does)
	formatted := formatEntry("", m.messageBuffer.String(), lipgloss.NewStyle(), m.width, false)
	m.viewport.SetContent(m.content.String() + formatted)
	m.viewport.GotoBottom()
	return true
}

// handleAgentEvent processes a single agent event and updates the model
// handleAgentEvent routes agent events to appropriate handlers. High complexity is inherent
// to event routing logic that must handle 15+ distinct event types with different behaviors.
//
//nolint:gocyclo
func (m *model) handleAgentEvent(event *types.AgentEvent) {
	switch event.Type {
	case types.EventTypeThinkingStart:
		m.isThinking = true
		m.thinkingBuffer.Reset()

	case types.EventTypeThinkingContent:
		if event.Content == "" {
			return
		}
		// Buffer the thinking content
		m.thinkingBuffer.WriteString(event.Content)
		// Stream thinking content as it arrives
		formatted := formatEntry("💭 ", m.thinkingBuffer.String(), thinkingStyle, m.width, false)
		m.viewport.SetContent(m.content.String() + formatted)
		m.viewport.GotoBottom()
		return

	case types.EventTypeThinkingEnd:
		if m.thinkingBuffer.Len() > 0 {
			formatted := formatEntry("💭 ", m.thinkingBuffer.String(), thinkingStyle, m.width, false)
			m.content.WriteString(formatted)
		}
		m.content.WriteString("\n\n")
		m.isThinking = false
		m.thinkingBuffer.Reset()

	case types.EventTypeToolCall:
		formatted := formatEntry("🔧 ", event.ToolName, toolStyle, m.width, false)
		m.content.WriteString(formatted)
		m.content.WriteString("\n")

	case types.EventTypeToolResult:
		resultStr := fmt.Sprintf("%v", event.ToolOutput)
		formatted := formatEntry("  ✓ ", resultStr, toolStyle, m.width, false)
		m.content.WriteString(formatted)
		m.content.WriteString("\n\n")

	case types.EventTypeMessageStart:
		m.messageBuffer.Reset()

	case types.EventTypeMessageContent:
		if strings.TrimSpace(event.Content) != "" && !m.hasMessageContentStarted {
			m.hasMessageContentStarted = true
		}
		if m.handleMessageContent(event.Content) {
			return // Viewport already updated in handleMessageContent
		}

	case types.EventTypeMessageEnd:
		// Finalize message content (like thinking does)
		if m.messageBuffer.Len() > 0 && m.hasMessageContentStarted {
			formatted := formatEntry("", m.messageBuffer.String(), lipgloss.NewStyle(), m.width, false)
			m.content.WriteString(formatted)
			m.content.WriteString("\n\n")
			m.hasMessageContentStarted = false
		}
		m.messageBuffer.Reset()

	case types.EventTypeError:
		m.content.WriteString(errorStyle.Render(fmt.Sprintf("  ❌ Error: %v", event.Error)))
		m.content.WriteString("\n\n")

	case types.EventTypeTurnEnd:
		// Turn end - clear busy state
		m.agentBusy = false
		m.recalculateLayout()
		return

	case types.EventTypeUpdateBusy:
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
		return

	case types.EventTypeToolApprovalRequest:
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
		return

	case types.EventTypeToolApprovalGranted:
		// Approval granted - show confirmation
		formatted := formatEntry("  ✓ ", "Tool approved - executing...", toolStyle, m.width, false)
		m.content.WriteString(formatted)
		m.content.WriteString("\n")

	case types.EventTypeToolApprovalRejected:
		// Approval rejected - log it
		formatted := formatEntry("  ✗ ", "Tool rejected by user", errorStyle, m.width, false)
		m.content.WriteString(formatted)
		m.content.WriteString("\n")

	case types.EventTypeToolApprovalTimeout:
		// Approval timeout - log it
		formatted := formatEntry("  ⏱ ", "Tool approval timed out", errorStyle, m.width, false)
		m.content.WriteString(formatted)
		m.content.WriteString("\n")

	case types.EventTypeApiCallStart:
		// Update context token information
		if event.ApiCallInfo != nil {
			m.currentContextTokens = event.ApiCallInfo.ContextTokens
			m.maxContextTokens = event.ApiCallInfo.MaxContextTokens
		}
		return // Don't update viewport for API call events

	case types.EventTypeTokenUsage:
		// Update token usage counts
		if event.TokenUsage != nil {
			m.totalPromptTokens += event.TokenUsage.PromptTokens
			m.totalCompletionTokens += event.TokenUsage.CompletionTokens
			m.totalTokens += event.TokenUsage.TotalTokens
		}
		return // Don't update viewport for token events

	case types.EventTypeCommandExecutionStart:
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
		return

	case types.EventTypeCommandOutput:
		// Output events are handled by the overlay itself
		return

	case types.EventTypeCommandExecutionComplete:
		// Command completed successfully
		if event.CommandExecution != nil {
			formatted := formatEntry("  ✓ ", fmt.Sprintf("Command completed (exit code: %d, duration: %s)",
				event.CommandExecution.ExitCode, event.CommandExecution.Duration), toolStyle, m.width, false)
			m.content.WriteString(formatted)
			m.content.WriteString("\n")
		}
		// Note: Overlay stays open until user dismisses it

	case types.EventTypeCommandExecutionFailed:
		// Command failed
		if event.CommandExecution != nil {
			formatted := formatEntry("  ✗ ", fmt.Sprintf("Command failed (exit code: %d, duration: %s)",
				event.CommandExecution.ExitCode, event.CommandExecution.Duration), errorStyle, m.width, false)
			m.content.WriteString(formatted)
			m.content.WriteString("\n")
		}
		// Note: Overlay stays open until user dismisses it

	case types.EventTypeCommandExecutionCanceled:
		// Command was canceled
		formatted := formatEntry("  ⏹ ", "Command canceled by user", toolStyle, m.width, false)
		m.content.WriteString(formatted)
		m.content.WriteString("\n")
		// Note: Overlay stays open until user dismisses it

	case types.EventTypeContextSummarizationStart:
		// Context summarization started - activate status overlay
		if event.ContextSummarization != nil {
			cs := event.ContextSummarization
			m.summarization.active = true
			m.summarization.strategy = cs.Strategy
			m.summarization.currentTokens = cs.CurrentTokens
			m.summarization.maxTokens = cs.MaxTokens
			m.summarization.itemsProcessed = 0
			m.summarization.totalItems = 0
			m.summarization.currentItem = ""
			m.summarization.progressPercent = 0
			m.summarization.startTime = time.Now()
		}
		return // Don't write to chat

	case types.EventTypeContextSummarizationProgress:
		// Context summarization progress update - update status overlay
		if event.ContextSummarization != nil && m.summarization.active {
			cs := event.ContextSummarization
			m.summarization.itemsProcessed = cs.ItemsProcessed
			m.summarization.totalItems = cs.TotalItems

			// Calculate progress percentage
			if cs.TotalItems > 0 {
				m.summarization.progressPercent = float64(cs.ItemsProcessed) / float64(cs.TotalItems) * 100.0
			}

			// Update current item description
			if cs.TotalItems > 0 {
				m.summarization.currentItem = fmt.Sprintf("Summarizing item %d of %d...",
					cs.ItemsProcessed, cs.TotalItems)
			}
		}
		return // Don't write to chat

	case types.EventTypeContextSummarizationComplete:
		// Context summarization completed - hide status and show success toast
		if event.ContextSummarization != nil {
			cs := event.ContextSummarization

			// Clear summarization status
			m.summarization.active = false

			// Show success toast
			m.toast.active = true
			m.toast.message = "Context optimized"
			m.toast.details = fmt.Sprintf("Saved %s tokens (%d items summarized)\n%s • %s",
				formatTokenCount(cs.TokensSaved),
				cs.ItemsProcessed,
				cs.Strategy,
				cs.Duration)
			m.toast.icon = "✨"
			m.toast.isError = false
			m.toast.showUntil = time.Now().Add(4 * time.Second)
		}
		return // Don't write to chat

	case types.EventTypeContextSummarizationError:
		// Context summarization failed - hide status and show error toast
		if event.ContextSummarization != nil {
			cs := event.ContextSummarization

			// Clear summarization status
			m.summarization.active = false

			// Show error toast
			m.toast.active = true
			m.toast.message = "Context optimization failed"
			m.toast.details = fmt.Sprintf("%s: %v\nContinuing with current context",
				cs.Strategy, event.Error)
			m.toast.icon = "⚠️"
			m.toast.isError = true
			m.toast.showUntil = time.Now().Add(4 * time.Second)
		}
		return // Don't write to chat
	}

	// Update viewport for all other event types
	m.viewport.SetContent(m.content.String())
	m.viewport.GotoBottom()
}

// recalculateLayout recalculates viewport height based on current dimensions
func (m *model) recalculateLayout() {
	if !m.ready {
		return
	}

	// Calculate heights for different sections
	headerHeight := 10 // ASCII art (6) + tips (1) + status bar (1) + blank line (1) + spacing (1)
	// Input height is dynamic based on textarea height (with border padding)
	inputHeight := m.textarea.Height() + 2 // textarea height + border
	statusBarHeight := 1

	// Account for loading indicator height when active
	loadingIndicatorHeight := 0
	if m.agentBusy {
		loadingIndicatorHeight = 1 // One line for the loading indicator
	}

	// Set viewport to fill remaining space
	viewportHeight := m.height - headerHeight - inputHeight - statusBarHeight - loadingIndicatorHeight
	if viewportHeight < 5 {
		viewportHeight = 5
	}

	m.viewport.Height = viewportHeight
}

// Update is called when a message is received. High complexity is inherent to the
// Bubble Tea Update pattern which must handle multiple message types and UI states.
//
//nolint:gocyclo
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	// Handle spinner tick messages
	var spinnerCmd tea.Cmd
	m.spinner, spinnerCmd = m.spinner.Update(msg)

	// Only update textarea if no overlay is active
	// This prevents the textarea from capturing scroll events when an overlay is open
	if !m.overlay.isActive() {
		// Store old textarea height to detect changes
		oldHeight := m.textarea.Height()
		m.textarea, tiCmd = m.textarea.Update(msg)
		newHeight := m.textarea.Height()

		// If textarea height changed, recalculate viewport height
		if oldHeight != newHeight && m.ready {
			m.recalculateLayout()
		}
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update viewport on window resize
		m.viewport, _ = m.viewport.Update(msg)
		m.width = msg.Width
		m.height = msg.Height

		// Calculate heights for different sections
		headerHeight := 10 // ASCII art (6) + tips (1) + status bar (1) + blank line (1) + spacing (1)
		// Input height is dynamic based on textarea height (with border padding)
		inputHeight := m.textarea.Height() + 2 // textarea height + border
		statusBarHeight := 1

		// Set viewport to fill remaining space
		viewportHeight := m.height - headerHeight - inputHeight - statusBarHeight
		if viewportHeight < 5 {
			viewportHeight = 5
		}

		m.viewport.Width = m.width - 4
		m.viewport.Height = viewportHeight
		m.textarea.SetWidth(m.width - 8)
		m.ready = true
		m.recalculateLayout()
		return m, nil

	case slashCommandCompleteMsg:
		// Slash command completed - clear busy state
		m.agentBusy = false
		m.recalculateLayout()
		return m, nil

	case operationStartMsg:
		// Generic operation started - show loading indicator
		m.agentBusy = true
		m.currentLoadingMessage = msg.message
		m.recalculateLayout()
		return m, nil

	case operationCompleteMsg:
		// Generic operation completed - hide loading and show toast
		m.agentBusy = false
		m.recalculateLayout()

		if msg.err != nil {
			m.showToast(msg.errorTitle, fmt.Sprintf("%v", msg.err), msg.errorIcon, true)
		} else {
			m.showToast(msg.successTitle, msg.result, msg.successIcon, false)
		}
		return m, nil

	case toastMsg:
		// Handle toast messages from slash commands
		m.showToast(msg.message, msg.details, msg.icon, msg.isError)
		// Also clear busy state when toast is shown (command completed)
		m.agentBusy = false
		m.recalculateLayout()
		return m, nil

	case approvalRequestMsg:
		// Generic approval handling - works with any ApprovalRequest
		overlay := NewGenericApprovalOverlay(msg.request, m.width, m.height)
		m.overlay.activate(OverlayModeSlashCommandPreview, overlay)
		return m, nil

	case agentErrMsg:
		m.content.WriteString(errorStyle.Render(fmt.Sprintf("\n  ❌ Agent Error: %v\n", msg.err)))
		m.viewport.SetContent(m.content.String())
		m.viewport.GotoBottom()
		return m, tea.Quit

	case *types.AgentEvent:
		// If overlay is active and it's a command execution event, forward to overlay
		if m.overlay.isActive() && msg.IsCommandExecutionEvent() {
			var overlayCmd tea.Cmd
			m.overlay.overlay, overlayCmd = m.overlay.overlay.Update(msg)
			// Still handle the event in the main model too
			m.handleAgentEvent(msg)
			return m, tea.Batch(tiCmd, vpCmd, overlayCmd)
		}

		// Update viewport for agent events
		m.viewport, vpCmd = m.viewport.Update(msg)
		m.handleAgentEvent(msg)
		return m, tea.Batch(tiCmd, vpCmd)

	case tea.MouseMsg:
		// Handle mouse events (especially scroll wheel) for viewport
		// If overlay is active, forward mouse events to it
		if m.overlay.isActive() {
			var overlayCmd tea.Cmd
			updatedOverlay, overlayCmd := m.overlay.overlay.Update(msg)

			// Check if overlay returned nil (signals to close)
			if updatedOverlay == nil {
				m.overlay.deactivate()
				m.viewport.SetContent(m.content.String())
				m.viewport.GotoBottom()
				return m, overlayCmd
			}

			m.overlay.overlay = updatedOverlay
			return m, overlayCmd
		}

		// Route mouse events to viewport for scrolling
		m.viewport, vpCmd = m.viewport.Update(msg)
		return m, tea.Batch(tiCmd, vpCmd)

	case tea.KeyMsg:
		// If overlay is active, forward all key messages to it
		if m.overlay.isActive() {
			var overlayCmd tea.Cmd
			updatedOverlay, overlayCmd := m.overlay.overlay.Update(msg)

			// Check if overlay returned nil (signals to close)
			if updatedOverlay == nil {
				m.overlay.deactivate()
				m.viewport.SetContent(m.content.String())
				m.viewport.GotoBottom()
				return m, overlayCmd
			}

			m.overlay.overlay = updatedOverlay
			return m, overlayCmd
		}

		// Handle command palette navigation when active
		// Only intercept specific keys, let everything else pass through to textarea
		if m.commandPalette.active {
			switch msg.Type {
			case tea.KeyEsc:
				// Cancel command palette
				m.commandPalette.deactivate()
				m.textarea.Reset()
				return m, tea.Batch(tiCmd, vpCmd)
			case tea.KeyUp:
				// Navigate up in palette, don't update textarea
				m.commandPalette.selectPrev()
				return m, nil
			case tea.KeyDown:
				// Navigate down in palette, don't update textarea
				m.commandPalette.selectNext()
				return m, nil
			case tea.KeyTab:
				// Autocomplete with selected command and close palette
				selected := m.commandPalette.getSelected()
				if selected != nil {
					m.textarea.SetValue("/" + selected.Name + " ")
					m.textarea.CursorEnd()
				}
				m.commandPalette.deactivate()
				return m, tea.Batch(tiCmd, vpCmd)
			case tea.KeyEnter:
				// Autocomplete with the selected command and close the palette
				selected := m.commandPalette.getSelected()
				if selected != nil {
					m.textarea.SetValue("/" + selected.Name + " ")
					m.textarea.CursorEnd()
				}
				m.commandPalette.deactivate()
				return m, tea.Batch(tiCmd, vpCmd)
			default:
				// For all other keys (typing, backspace, etc.), let textarea handle them
				// The textarea has already been updated at the top of Update()
				// Return here to prevent the outer switch from handling these keys
				return m, tea.Batch(tiCmd, vpCmd)
			}
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			// Alt+Enter (Option+Enter on Mac) adds newline
			// Plain Enter submits the message
			if msg.Alt {
				m.textarea.InsertString("\n")
				m.updateTextAreaHeight()
				return m, tea.Batch(tiCmd, vpCmd)
			}
			// Plain Enter submits
			input := m.textarea.Value()
			if input == "" {
				return m, tea.Batch(tiCmd, vpCmd)
			}

			// Check if input is a slash command
			commandName, args, isCommand := parseSlashCommand(input)
			if isCommand {
				// Execute slash command and get updated model
				m.textarea.Reset()
				m.commandPalette.deactivate()
				return executeSlashCommand(m, commandName, args)
			}

			// If input starts with / but isn't a valid command, don't send to agent
			if strings.HasPrefix(strings.TrimSpace(input), "/") {
				// Keep palette active, don't submit
				return m, tea.Batch(tiCmd, vpCmd)
			}

			// Regular user input - send to agent
			if m.channels != nil {
				formatted := formatEntry("You: ", input, userStyle, m.width, true)
				// Strip any trailing newlines before adding our spacing
				formatted = strings.TrimRight(formatted, "\n")
				m.content.WriteString(formatted + "\n\n")
				m.viewport.SetContent(m.content.String())

				// Set agent as busy and pick a random loading message
				m.agentBusy = true
				m.currentLoadingMessage = getRandomLoadingMessage()
				m.recalculateLayout()

				m.channels.Input <- types.NewUserInput(input)
				m.textarea.Reset()
				m.viewport.GotoBottom()
			}
			return m, tea.Batch(tiCmd, vpCmd, spinnerCmd)
		default:
			// Let viewport handle other keys (arrow keys, pgup/pgdn, etc. for scrolling)
			m.viewport, vpCmd = m.viewport.Update(msg)
		}
	}

	// Check if we should activate/deactivate command palette based on input
	value := m.textarea.Value()

	// Handle command palette activation/deactivation based on input
	switch {
	case value == "/" && !m.commandPalette.active:
		// Only activate palette if input is exactly "/" as first character
		m.commandPalette.activate()
		m.commandPalette.updateFilter("")
	case strings.HasPrefix(value, "/") && m.commandPalette.active:
		// Update filter if palette is already active
		filter := strings.TrimPrefix(value, "/")
		m.commandPalette.updateFilter(filter)
	case !strings.HasPrefix(value, "/") && m.commandPalette.active:
		// Deactivate palette if input no longer starts with /
		m.commandPalette.deactivate()
	}

	// Auto-adjust textarea height based on content after any key press
	m.updateTextAreaHeight()

	return m, tea.Batch(tiCmd, vpCmd, spinnerCmd)
}

// updateTextAreaHeight adjusts textarea height based on number of lines
// including visual line wrapping
func (m *model) updateTextAreaHeight() {
	value := m.textarea.Value()
	if value == "" {
		if m.textarea.Height() != 1 {
			m.textarea.SetHeight(1)
			m.recalculateLayout()
		}
		return
	}

	// Calculate visual lines accounting for wrapping
	width := m.textarea.Width()
	if width <= 0 {
		width = 80 // default width
	}

	// Account for prompt width ("> " = 2 chars)
	effectiveWidth := width - 2
	if effectiveWidth <= 0 {
		effectiveWidth = 78
	}

	// Split by actual newlines first
	textLines := strings.Split(value, "\n")
	visualLines := 0

	for _, line := range textLines {
		if line == "" {
			visualLines++ // Empty line still counts as 1 visual line
		} else {
			// Calculate how many visual lines this logical line takes
			lineLen := len(line)
			wrappedLines := (lineLen + effectiveWidth - 1) / effectiveWidth
			if wrappedLines == 0 {
				wrappedLines = 1
			}
			visualLines += wrappedLines
		}
	}

	// Clamp between 1 and MaxHeight
	if visualLines < 1 {
		visualLines = 1
	}
	if visualLines > m.textarea.MaxHeight {
		visualLines = m.textarea.MaxHeight
	}

	// Only update if height changed to avoid unnecessary recalculation
	if visualLines != m.textarea.Height() {
		m.textarea.SetHeight(visualLines)
		m.recalculateLayout()
	}
}

// View renders the TUI.
func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// ASCII art header with gradient effect
	header := headerStyle.Render(`
	███████╗ ██████╗ ██████╗  ██████╗ ███████╗
	██╔════╝██╔═══██╗██╔══██╗██╔════╝ ██╔════╝
	█████╗  ██║   ██║██████╔╝██║  ███╗█████╗
	██╔══╝  ██║   ██║██╔══██╗██║   ██║██╔══╝
	██║     ╚██████╔╝██║  ██║╚██████╔╝███████╗
	╚═╝      ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚══════╝`)

	// Tips section
	tips := tipsStyle.Render(`  Tips: Ask questions • Alt+Enter for new line • Enter to send • Ctrl+C to exit`)

	// Top status bar
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "~"
	}
	topStatus := statusBarStyle.Render(fmt.Sprintf("  Working directory: %s", cwd))

	// Loading indicator (shown above input box when agent is busy)
	var loadingIndicator string
	if m.agentBusy {
		loadingMsg := fmt.Sprintf("%s %s", m.spinner.View(), m.currentLoadingMessage)
		loadingStyle := lipgloss.NewStyle().
			Foreground(salmonPink).
			Width(m.width-4).
			Padding(0, 2)
		loadingIndicator = loadingStyle.Render(loadingMsg)
	}

	// Input box
	inputBox := inputBoxStyle.Width(m.width - 4).Render(m.textarea.View())

	// Bottom status bar with three sections
	bottomLeft := "~/forge"
	bottomCenter := "Enter to send • Alt+Enter for new line"

	// Right section includes token usage if available
	bottomRight := "Forge Agent"
	if m.totalTokens > 0 {
		// Format context with visual indicator if approaching limit
		contextStr := formatTokenCount(m.currentContextTokens)
		if m.maxContextTokens > 0 {
			contextStr = fmt.Sprintf("%s/%s", contextStr, formatTokenCount(m.maxContextTokens))

			// Add color indicator if context is >= 80% of max (approaching limit)
			percentage := float64(m.currentContextTokens) / float64(m.maxContextTokens) * 100
			if percentage >= 80 {
				contextStr = lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Render(contextStr) // Orange/red
			}
		}

		bottomRight = fmt.Sprintf("◆ Context: %s | Input: %s | Output: %s | Total: %s",
			contextStr,
			formatTokenCount(m.totalPromptTokens),
			formatTokenCount(m.totalCompletionTokens),
			formatTokenCount(m.totalTokens))
	}

	// Calculate spacing
	totalUsed := len(bottomLeft) + len(bottomCenter) + len(bottomRight)
	leftPadding := (m.width - totalUsed) / 3
	rightPadding := m.width - totalUsed - leftPadding*2
	if leftPadding < 2 {
		leftPadding = 2
	}
	if rightPadding < 2 {
		rightPadding = 2
	}

	bottomBar := statusBarStyle.Width(m.width).Render(
		bottomLeft +
			strings.Repeat(" ", leftPadding) +
			bottomCenter +
			strings.Repeat(" ", rightPadding) +
			bottomRight,
	)

	// Build viewport section - just the viewport itself
	viewportSection := m.viewport.View()

	// Assemble the base UI without overlays
	var baseView string
	if m.agentBusy {
		// Include loading indicator when agent is busy
		baseView = lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			tips,
			topStatus,
			"", // Blank line for spacing
			viewportSection,
			loadingIndicator,
			inputBox,
			bottomBar,
		)
	} else {
		// Normal view without loading indicator
		baseView = lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			tips,
			topStatus,
			"", // Blank line for spacing
			viewportSection,
			inputBox,
			bottomBar,
		)
	}

	// Layer overlays on top of the base view using absolute positioning
	if m.overlay.isActive() {
		baseView = renderOverlay(baseView, m.overlay.overlay, m.width, m.height)
	}

	// Add command palette as overlay if active
	if m.commandPalette.active {
		paletteContent := m.commandPalette.render(m.width)
		baseView = renderToastOverlay(baseView, paletteContent)
	}

	// Add summarization status as overlay if active
	if m.summarization.active {
		summarizationContent := m.renderSummarizationStatus()
		baseView = renderToastOverlay(baseView, summarizationContent)
	}

	// Add toast notification as overlay if active and not expired
	if m.toast.active && time.Now().Before(m.toast.showUntil) {
		toastContent := m.renderToast()
		baseView = renderToastOverlay(baseView, toastContent)
	}

	return baseView
}
