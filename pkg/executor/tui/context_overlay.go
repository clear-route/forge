package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/entrhq/forge/internal/utils"
)

// ContextOverlay displays detailed context information in a modal dialog
type ContextOverlay struct {
	viewport viewport.Model
	title    string
	content  string
}

// ContextInfo contains all context statistics to display
type ContextInfo struct {
	// System prompt
	SystemPromptTokens int
	CustomInstructions bool

	// Tool system
	ToolCount          int
	ToolTokens         int
	ToolNames          []string
	CurrentToolCall    string
	HasPendingToolCall bool

	// Message history
	MessageCount       int
	ConversationTurns  int
	ConversationTokens int

	// Token usage - current context
	CurrentContextTokens int
	MaxContextTokens     int
	FreeTokens           int
	UsagePercent         float64

	// Token usage - cumulative across all API calls
	TotalPromptTokens     int
	TotalCompletionTokens int
	TotalTokens           int
}

// NewContextOverlay creates a new context information overlay
func NewContextOverlay(info *ContextInfo) *ContextOverlay {
	content := buildContextContent(info)

	vp := viewport.New(76, 20)
	vp.Style = lipgloss.NewStyle()
	vp.SetContent(content)

	return &ContextOverlay{
		viewport: vp,
		title:    "Context Information",
		content:  content,
	}
}

// buildContextContent formats the context information for display
func buildContextContent(info *ContextInfo) string {
	var b strings.Builder

	// System section
	b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(salmonPink).Render("System"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  System Prompt:      %s tokens\n", formatTokenCount(info.SystemPromptTokens)))
	if info.CustomInstructions {
		b.WriteString("  Custom Instructions: Yes\n")
	} else {
		b.WriteString("  Custom Instructions: No\n")
	}
	b.WriteString("\n")

	// Tool System section
	b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(salmonPink).Render("Tool System"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  Available Tools:    %d (%s tokens)\n", info.ToolCount, formatTokenCount(info.ToolTokens)))
	if info.HasPendingToolCall {
		b.WriteString(fmt.Sprintf("  Current Tool Call:  %s\n", info.CurrentToolCall))
	}
	b.WriteString("\n")

	// History section
	b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(salmonPink).Render("Message History"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  Messages:           %d\n", info.MessageCount))
	b.WriteString(fmt.Sprintf("  Conversation Turns: %d\n", info.ConversationTurns))
	b.WriteString(fmt.Sprintf("  Conversation:       %s tokens\n", formatTokenCount(info.ConversationTokens)))
	b.WriteString("\n")

	// Current Context section
	b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(salmonPink).Render("Current Context"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  Used:               %s / %s tokens (%.1f%%)\n",
		formatTokenCount(info.CurrentContextTokens),
		formatTokenCount(info.MaxContextTokens),
		info.UsagePercent))
	b.WriteString(fmt.Sprintf("  Free Space:         %s tokens\n", formatTokenCount(info.FreeTokens)))

	// Add a progress bar
	barWidth := 40
	filledWidth := int(float64(barWidth) * info.UsagePercent / 100.0)
	emptyWidth := barWidth - filledWidth

	var barColor lipgloss.Color
	switch {
	case info.UsagePercent < 70:
		barColor = lipgloss.Color("#98C379") // Green
	case info.UsagePercent < 90:
		barColor = lipgloss.Color("#E5C07B") // Yellow
	default:
		barColor = lipgloss.Color("#E06C75") // Red
	}

	filled := lipgloss.NewStyle().Foreground(barColor).Render(strings.Repeat("█", filledWidth))
	empty := lipgloss.NewStyle().Foreground(lipgloss.Color("#3E4451")).Render(strings.Repeat("░", emptyWidth))
	b.WriteString(fmt.Sprintf("  [%s%s]\n", filled, empty))
	b.WriteString("\n")

	// Cumulative Token Usage section
	b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(salmonPink).Render("Cumulative Usage (All API Calls)"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  Input Tokens:       %s\n", formatTokenCount(info.TotalPromptTokens)))
	b.WriteString(fmt.Sprintf("  Output Tokens:      %s\n", formatTokenCount(info.TotalCompletionTokens)))
	b.WriteString(fmt.Sprintf("  Total:              %s\n", formatTokenCount(info.TotalTokens)))

	return b.String()
}

// Update handles messages for the context overlay
func (c *ContextOverlay) Update(msg tea.Msg) (Overlay, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC, tea.KeyEnter:
			return nil, nil
		case tea.KeyUp, tea.KeyDown, tea.KeyPgUp, tea.KeyPgDown:
			c.viewport, cmd = c.viewport.Update(msg)
			return c, cmd
		}

	case tea.WindowSizeMsg:
		// Adjust viewport height if screen is too small
		c.viewport.Height = utils.Min(20, msg.Height-10)
	}

	return c, nil
}

// View renders the context overlay
func (c *ContextOverlay) View() string {
	header := OverlayTitleStyle.Render(c.title)
	viewportContent := c.viewport.View()
	footer := OverlayHelpStyle.Render("Press ESC or Enter to close • ↑/↓ to scroll")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		viewportContent,
		footer,
	)

	return CreateOverlayContainerStyle(c.viewport.Width).Render(content)
}

// Focused returns whether this overlay should handle input
func (c *ContextOverlay) Focused() bool {
	return true
}

// Width returns the overlay width
func (c *ContextOverlay) Width() int {
	return c.viewport.Width
}

// Height returns the overlay height
func (c *ContextOverlay) Height() int {
	return c.viewport.Height
}
