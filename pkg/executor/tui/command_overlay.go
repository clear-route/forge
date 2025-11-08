package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/entrhq/forge/pkg/types"
)

var (
	// Command overlay styles
	commandOverlayStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(salmonPink).
				Padding(1, 2).
				Width(80).
				Height(30)

	commandHeaderStyle = lipgloss.NewStyle().
				Foreground(salmonPink).
				Bold(true).
				MarginBottom(1)

	commandStatusStyle = lipgloss.NewStyle().
				Foreground(mutedGray).
				Italic(true)

	commandHelpStyle = lipgloss.NewStyle().
				Foreground(mutedGray).
				MarginTop(1)
)

// CommandExecutionOverlay displays streaming command output with cancellation support
type CommandExecutionOverlay struct {
	viewport      viewport.Model
	command       string
	workingDir    string
	executionID   string
	output        *strings.Builder
	status        string
	exitCode      int
	isRunning     bool
	width         int
	height        int
	cancelChannel chan<- *types.CancellationRequest
}

// NewCommandExecutionOverlay creates a new command execution overlay
func NewCommandExecutionOverlay(command, workingDir, executionID string, cancelChan chan<- *types.CancellationRequest) *CommandExecutionOverlay {
	vp := viewport.New(76, 20) // Slightly smaller than overlay for padding
	vp.Style = lipgloss.NewStyle()

	return &CommandExecutionOverlay{
		viewport:      vp,
		command:       command,
		workingDir:    workingDir,
		executionID:   executionID,
		output:        &strings.Builder{},
		status:        "Running...",
		isRunning:     true,
		width:         80,
		height:        30,
		cancelChannel: cancelChan,
	}
}

// Update handles messages for the command overlay
func (c *CommandExecutionOverlay) Update(msg tea.Msg) (Overlay, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if c.isRunning {
			// Handle cancellation keys
			switch msg.String() {
			case "ctrl+c", "esc":
				// Send cancellation request
				if c.cancelChannel != nil {
					c.cancelChannel <- &types.CancellationRequest{
						ExecutionID: c.executionID,
					}
				}
				c.status = "Canceling..."
				return c, nil
			}
		} else {
			// Command is finished - any key closes the overlay
			// Return nil to signal overlay should be closed
			return nil, nil
		}

		// Handle viewport scrolling (only while command is running)
		switch msg.String() {
		case "up", "k":
			c.viewport, cmd = c.viewport.Update(msg)
			return c, cmd
		case "down", "j":
			c.viewport, cmd = c.viewport.Update(msg)
			return c, cmd
		case "pgup", "b":
			c.viewport, cmd = c.viewport.Update(msg)
			return c, cmd
		case "pgdown", "f":
			c.viewport, cmd = c.viewport.Update(msg)
			return c, cmd
		case "home", "g":
			c.viewport.GotoTop()
			return c, nil
		case "end", "G":
			c.viewport.GotoBottom()
			return c, nil
		}

	case *types.AgentEvent:
		// Handle command execution events
		if msg.IsCommandExecutionEvent() {
			return c.handleCommandEvent(msg)
		}

	case tea.WindowSizeMsg:
		// Update overlay size
		c.width = min(msg.Width-4, 80)
		c.height = min(msg.Height-4, 30)

		// Update viewport size
		viewportWidth := c.width - 4   // Account for border and padding
		viewportHeight := c.height - 8 // Account for header, status, help text
		c.viewport.Width = viewportWidth
		c.viewport.Height = viewportHeight
	}

	return c, nil
}

// handleCommandEvent processes command execution events
func (c *CommandExecutionOverlay) handleCommandEvent(event *types.AgentEvent) (Overlay, tea.Cmd) {
	if event.CommandExecution == nil {
		return c, nil
	}

	data := event.CommandExecution

	// Only process events for this execution
	if data.ExecutionID != c.executionID {
		return c, nil
	}

	switch event.Type {
	case types.EventTypeCommandOutput:
		// Append new output
		c.output.WriteString(data.Output)
		c.viewport.SetContent(c.output.String())

		// Auto-scroll to bottom if we were already at the bottom
		if c.viewport.AtBottom() {
			c.viewport.GotoBottom()
		}

	case types.EventTypeCommandExecutionComplete:
		c.isRunning = false
		c.exitCode = data.ExitCode
		c.status = fmt.Sprintf("Completed in %s (exit code: %d)", data.Duration, data.ExitCode)

	case types.EventTypeCommandExecutionFailed:
		c.isRunning = false
		c.exitCode = data.ExitCode
		c.status = fmt.Sprintf("Failed in %s (exit code: %d)", data.Duration, data.ExitCode)

	case types.EventTypeCommandExecutionCanceled:
		c.isRunning = false
		c.status = "Canceled by user"
	}

	return c, nil
}

// View renders the command overlay
func (c *CommandExecutionOverlay) View() string {
	// Build header
	header := commandHeaderStyle.Render("Command Execution")

	// Build command info
	commandInfo := fmt.Sprintf("Command: %s", c.command)
	if c.workingDir != "" {
		commandInfo += fmt.Sprintf("\nWorking Dir: %s", c.workingDir)
	}

	// Build status line
	statusLine := commandStatusStyle.Render(c.status)

	// Build output viewport
	outputView := c.viewport.View()

	// Build help text
	var helpText string
	if c.isRunning {
		helpText = commandHelpStyle.Render("Ctrl+C or Esc: Cancel | ↑↓: Scroll | PgUp/PgDn: Page")
	} else {
		helpText = commandHelpStyle.Render("Press any key to close")
	}

	// Combine all parts
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		commandInfo,
		statusLine,
		"",
		outputView,
		helpText,
	)

	// Apply overall styling
	return commandOverlayStyle.
		Width(c.width).
		Height(c.height).
		Render(content)
}

// Focused returns whether this overlay should handle input
func (c *CommandExecutionOverlay) Focused() bool {
	return true
}

// Width returns the overlay width
func (c *CommandExecutionOverlay) Width() int {
	return c.width
}

// Height returns the overlay height
func (c *CommandExecutionOverlay) Height() int {
	return c.height
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
