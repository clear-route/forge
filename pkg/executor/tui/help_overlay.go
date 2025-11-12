package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/entrhq/forge/internal/utils"
)

// HelpOverlay displays help information in a modal dialog
type HelpOverlay struct {
	viewport viewport.Model
	title    string
	content  string
}

// NewHelpOverlay creates a new help overlay
func NewHelpOverlay(title, content string) *HelpOverlay {
	vp := viewport.New(76, 20)
	vp.Style = lipgloss.NewStyle()
	vp.SetContent(content)

	return &HelpOverlay{
		viewport: vp,
		title:    title,
		content:  content,
	}
}

// Update handles messages for the help overlay
func (h *HelpOverlay) Update(msg tea.Msg) (Overlay, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC, tea.KeyEnter:
			return nil, nil
		case tea.KeyUp, tea.KeyDown, tea.KeyPgUp, tea.KeyPgDown:
			h.viewport, cmd = h.viewport.Update(msg)
			return h, cmd
		}

	case tea.WindowSizeMsg:
		// Adjust viewport height if screen is too small, but width is fixed.
		h.viewport.Height = utils.Min(20, msg.Height-10)
	}

	return h, nil
}

// View renders the help overlay
func (h *HelpOverlay) View() string {
	header := OverlayTitleStyle.Render(h.title)
	viewportContent := h.viewport.View()
	footer := OverlayHelpStyle.Render("Press ESC or Enter to close")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		viewportContent,
		footer,
	)

	// Use the viewport's width to determine the container width
	return CreateOverlayContainerStyle(h.viewport.Width).Render(content)
}

// Focused returns whether this overlay should handle input
func (h *HelpOverlay) Focused() bool {
	return true
}

// Width returns the overlay width.
// This is not used for positioning, which is handled by lipgloss.Place,
// but is part of the Overlay interface.
func (h *HelpOverlay) Width() int {
	return h.viewport.Width
}

// Height returns the overlay height.
// This is not used for positioning, but is part of the Overlay interface.
func (h *HelpOverlay) Height() int {
	return h.viewport.Height
}
