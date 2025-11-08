// Package tui provides a terminal user interface executor for Forge agents,
// offering an interactive, Gemini-style interface for conversations.
package tui

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/entrhq/forge/pkg/agent"
	"github.com/entrhq/forge/pkg/agent/tools"
	"github.com/entrhq/forge/pkg/types"
)

var (
	// Color palette with pastel salmon pink
	salmonPink  = lipgloss.Color("#FFB3BA") // Soft pastel salmon pink
	coralPink   = lipgloss.Color("#FFCCCB") // Lighter coral accent
	mutedGray   = lipgloss.Color("#6B7280")
	brightWhite = lipgloss.Color("#F9FAFB")
	darkBg      = lipgloss.Color("#111827")

	// Styles
	headerStyle = lipgloss.NewStyle().
			Foreground(salmonPink).
			Bold(true)

	tipsStyle = lipgloss.NewStyle().
			Foreground(mutedGray)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(mutedGray).
			Background(darkBg).
			Padding(0, 1)

	inputBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(salmonPink).
			Padding(0, 1)

	userStyle     = lipgloss.NewStyle().Foreground(coralPink).Bold(true)
	thinkingStyle = lipgloss.NewStyle().Foreground(mutedGray).Italic(true)
	toolStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#A8E6CF")) // Soft mint green
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB3BA")) // Match salmon for errors
)

// Executor is a TUI-based executor that provides an interactive,
// Gemini-style interface for agent interaction.
type Executor struct {
	agent   agent.Agent
	program *tea.Program
}

// NewExecutor creates a new TUI executor for the given agent.
func NewExecutor(agent agent.Agent) *Executor {
	return &Executor{
		agent: agent,
	}
}

// Run starts the TUI executor and blocks until the user exits.
func (e *Executor) Run(ctx context.Context) error {
	// Start the agent first
	if err := e.agent.Start(ctx); err != nil {
		return fmt.Errorf("failed to start agent: %w", err)
	}

	model := initialModel()
	model.agent = e.agent
	model.channels = e.agent.GetChannels()

	e.program = tea.NewProgram(
		model,
		tea.WithAltScreen(),
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

// model represents the state of the TUI application.
type model struct {
	viewport       viewport.Model
	textarea       textarea.Model
	agent          agent.Agent
	channels       *types.AgentChannels
	content        *strings.Builder
	thinkingBuffer *strings.Builder
	messageBuffer  *strings.Builder
	overlay        *overlayState
	isThinking     bool
	width          int
	height         int
	ready          bool
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

	return model{
		viewport:       vp,
		textarea:       ta,
		content:        &strings.Builder{},
		thinkingBuffer: &strings.Builder{},
		messageBuffer:  &strings.Builder{},
		overlay:        newOverlayState(),
	}
}

// Init is the first function that will be called.
func (m model) Init() tea.Cmd {
	return textarea.Blink
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
		if currentLine == "" {
			currentLine = word
		} else if len(currentLine)+1+len(word) > width {
			// Write current line and start new one
			result.WriteString(currentLine)
			result.WriteString("\n")
			currentLine = word
		} else {
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

// handleMessageContent processes message content events with streaming support
func (m *model) handleMessageContent(content string) bool {
	if content == "" {
		return false
	}

	// Buffer and display content directly (no "Assistant:" label needed as responses come through tools)
	m.messageBuffer.WriteString(content)
	m.content.WriteString(content)

	// Update viewport immediately for streaming effect
	m.viewport.SetContent(m.content.String())
	m.viewport.GotoBottom()
	return true
}

// handleAgentEvent processes a single agent event and updates the model
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
		if m.handleMessageContent(event.Content) {
			return // Viewport already updated in handleMessageContent
		}

	case types.EventTypeMessageEnd:
		m.messageBuffer.Reset()

	case types.EventTypeError:
		m.content.WriteString("\n" + errorStyle.Render(fmt.Sprintf("  ❌ Error: %v", event.Error)))

	case types.EventTypeTurnEnd:
		// Turn end - no extra spacing needed
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
					log.Printf("DEBUG Executor: responseFunc called with decision=%v for approval %s", response.Decision, response.ApprovalID)
					
					// Send approval response to agent
					log.Printf("DEBUG Executor: Sending response to approval channel...")
					m.channels.Approval <- response
					log.Printf("DEBUG Executor: Response sent successfully")
					
					// Close overlay and update viewport
					log.Printf("DEBUG Executor: Deactivating overlay and updating viewport")
					m.overlay.deactivate()
					m.viewport.SetContent(m.content.String())
					m.viewport.GotoBottom()
					log.Printf("DEBUG Executor: Viewport updated, callback complete")
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

	// Set viewport to fill remaining space
	viewportHeight := m.height - headerHeight - inputHeight - statusBarHeight
	if viewportHeight < 5 {
		viewportHeight = 5
	}

	m.viewport.Height = viewportHeight
}

// Update is called when a message is received.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	// Store old textarea height to detect changes
	oldHeight := m.textarea.Height()
	m.textarea, tiCmd = m.textarea.Update(msg)
	newHeight := m.textarea.Height()

	// If textarea height changed, recalculate viewport height
	if oldHeight != newHeight && m.ready {
		m.recalculateLayout()
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update viewport on window resize
		m.viewport, vpCmd = m.viewport.Update(msg)
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

	case agentErrMsg:
		m.content.WriteString(errorStyle.Render(fmt.Sprintf("\n  ❌ Agent Error: %v\n", msg.err)))
		m.viewport.SetContent(m.content.String())
		m.viewport.GotoBottom()
		return m, tea.Quit

	case *types.AgentEvent:
		// Update viewport for agent events
		m.viewport, vpCmd = m.viewport.Update(msg)
		m.handleAgentEvent(msg)
		return m, tea.Batch(tiCmd, vpCmd)

	case tea.KeyMsg:
		// If overlay is active, forward all key messages to it
		if m.overlay.isActive() {
			var overlayCmd tea.Cmd
			m.overlay.overlay, overlayCmd = m.overlay.overlay.Update(msg)
			return m, overlayCmd
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
			m.handleUserInput()
			return m, tea.Batch(tiCmd, vpCmd)
		default:
			// Let viewport handle other keys (arrow keys, pgup/pgdn, etc. for scrolling)
			m.viewport, vpCmd = m.viewport.Update(msg)
		}
	}

	// Auto-adjust textarea height based on content after any key press
	m.updateTextAreaHeight()

	return m, tea.Batch(tiCmd, vpCmd)
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

// handleUserInput processes user input and sends it to the agent
func (m *model) handleUserInput() {
	input := m.textarea.Value()
	if input != "" && m.channels != nil {
		formatted := formatEntry("You: ", input, userStyle, m.width, true)
		// Strip any trailing newlines before adding our spacing
		formatted = strings.TrimRight(formatted, "\n")
		m.content.WriteString(formatted + "\n\n")
		m.viewport.SetContent(m.content.String())
		m.channels.Input <- types.NewUserInput(input)
		m.textarea.Reset()
		m.viewport.GotoBottom()
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

	// Input box
	inputBox := inputBoxStyle.Width(m.width - 4).Render(m.textarea.View())

	// Bottom status bar with three sections
	bottomLeft := "~/forge"
	bottomCenter := "Enter to send • Alt+Enter for new line"
	bottomRight := "Forge Agent"

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

	// Assemble the full UI
	baseView := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		tips,
		topStatus,
		"", // Blank line for spacing
		m.viewport.View(),
		inputBox,
		bottomBar,
	)

	// If overlay is active, render it on top
	if m.overlay.isActive() {
		return renderOverlay(baseView, m.overlay.overlay, m.width, m.height)
	}

	return baseView
}
