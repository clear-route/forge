// Package tui provides a terminal user interface executor for Forge agents,
// offering an interactive, Gemini-style interface for conversations.
package tui

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/entrhq/forge/pkg/agent"
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
	ta.ShowLineNumbers = false
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
func wordWrap(text string, width int) string {
	if width <= 0 {
		width = 80
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var result strings.Builder
	currentLine := words[0]

	for _, word := range words[1:] {
		// Check if adding this word would exceed width
		if len(currentLine)+1+len(word) > width {
			// Write current line and start new one
			result.WriteString(currentLine)
			result.WriteString("\n")
			currentLine = word
		} else {
			// Add word to current line
			currentLine += " " + word
		}
	}

	// Write final line
	result.WriteString(currentLine)
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
	}

	// Update viewport for all other event types
	m.viewport.SetContent(m.content.String())
	m.viewport.GotoBottom()
}

// Update is called when a message is received.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Calculate heights for different sections
		headerHeight := 9 // ASCII art (6) + tips (1) + status bar (1) + spacing (1)
		inputHeight := 3  // Input box with border
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
		return m, nil

	case agentErrMsg:
		m.content.WriteString(errorStyle.Render(fmt.Sprintf("\n  ❌ Agent Error: %v\n", msg.err)))
		m.viewport.SetContent(m.content.String())
		m.viewport.GotoBottom()
		return m, tea.Quit

	case *types.AgentEvent:
		m.handleAgentEvent(msg)
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.handleUserInput()
		}
	}

	return m, tea.Batch(tiCmd, vpCmd)
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
	tips := tipsStyle.Render(`  Tips: Ask questions • Tools run automatically • Ctrl+C to exit`)

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
	bottomCenter := "Press Ctrl+C to exit"
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
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		tips,
		topStatus,
		m.viewport.View(),
		inputBox,
		bottomBar,
	)
}
