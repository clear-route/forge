package overlay

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/entrhq/forge/pkg/executor/tui/types"
)

// ToolResultOverlay displays the full result of a tool call
type ToolResultOverlay struct {
	toolName string
	result   string
	viewport viewport.Model
	width    int
	height   int
}

// NewToolResultOverlay creates a new tool result overlay
func NewToolResultOverlay(toolName, result string, width, height int) *ToolResultOverlay {
	// Calculate overlay dimensions (80% of screen)
	overlayWidth := int(float64(width) * 0.8)
	overlayHeight := int(float64(height) * 0.8)

	if overlayWidth < 60 {
		overlayWidth = 60
	}
	if overlayHeight < 20 {
		overlayHeight = 20
	}

	// Create viewport for scrolling
	vp := viewport.New(overlayWidth-4, overlayHeight-6) // Account for border and header
	vp.SetContent(result)

	return &ToolResultOverlay{
		toolName: toolName,
		result:   result,
		viewport: vp,
		width:    overlayWidth,
		height:   overlayHeight,
	}
}

// Update handles messages
func (o *ToolResultOverlay) Update(msg tea.Msg, state types.StateProvider, actions types.ActionHandler) (types.Overlay, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "q", "esc", "v":
			if actions != nil {
				actions.ClearOverlay()
			}
			return nil, nil // Signal to close
		}
	}

	// Forward to viewport for scrolling
	var cmd tea.Cmd
	o.viewport, cmd = o.viewport.Update(msg)
	return o, cmd
}

// View renders the overlay
func (o *ToolResultOverlay) View() string {
	// Create header
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(types.DiffHunkColor). // Cyan/Sky blue
		Render(fmt.Sprintf("Tool Result: %s", o.toolName))

	// Create help text
	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("↑/↓: scroll • q/esc/v: close")

	// Combine header and help
	headerSection := lipgloss.JoinVertical(lipgloss.Left,
		header,
		help,
		"",
	)

	// Create the content area with viewport
	content := lipgloss.JoinVertical(lipgloss.Left,
		headerSection,
		o.viewport.View(),
	)

	// Create bordered box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(types.DiffHunkColor). // Cyan/Sky blue
		Padding(1, 2).
		Width(o.width).
		Height(o.height)

	return centerOverlay(boxStyle.Render(content), o.width, o.height)
}

// Focused returns whether this overlay should handle input
func (o *ToolResultOverlay) Focused() bool {
	return true
}

// Width returns the overlay width
func (o *ToolResultOverlay) Width() int {
	return o.width
}

// Height returns the overlay height
func (o *ToolResultOverlay) Height() int {
	return o.height
}

// centerOverlay centers an overlay on the screen
func centerOverlay(content string, width, height int) string {
	lines := strings.Split(content, "\n")
	var centered strings.Builder

	// Add vertical padding
	verticalPadding := max(0, (height-len(lines))/2)
	for i := 0; i < verticalPadding; i++ {
		centered.WriteString("\n")
	}

	// Center each line horizontally
	for _, line := range lines {
		// Calculate horizontal padding (account for ANSI codes)
		lineWidth := lipgloss.Width(line)
		horizontalPadding := max(0, (width-lineWidth)/2)
		if horizontalPadding > 0 {
			centered.WriteString(strings.Repeat(" ", horizontalPadding))
		}
		centered.WriteString(line)
		centered.WriteString("\n")
	}

	return centered.String()
}
