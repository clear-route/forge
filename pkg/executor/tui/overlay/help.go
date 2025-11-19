package overlay

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/entrhq/forge/pkg/executor/tui/types"
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
func (h *HelpOverlay) Update(msg tea.Msg, state types.StateProvider, actions types.ActionHandler) (types.Overlay, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC, tea.KeyEnter:
			// Instead of returning nil, nil which meant "close" in the old system,
			// we explicitly ask to clear the overlay.
			// However, the interface contract says we return the updated *Overlay*.
			// If we want to close, we should return nil or have actions.ClearOverlay().
			// In the old system: "return nil, nil" meant "I am done, remove me".
			return nil, nil
		case tea.KeyUp, tea.KeyDown, tea.KeyPgUp, tea.KeyPgDown:
			h.viewport, cmd = h.viewport.Update(msg)
			return h, cmd
		}

	case tea.WindowSizeMsg:
		// Adjust viewport height if screen is too small, but width is fixed.
		h.viewport.Height = min(20, msg.Height-10)
	}

	return h, nil
}

// View renders the help overlay
func (h *HelpOverlay) View() string {
	header := types.OverlayTitleStyle.Render(h.title)
	viewportContent := h.viewport.View()
	footer := types.OverlayHelpStyle.Render("Press ESC or Enter to close")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		viewportContent,
		footer,
	)

	// Use the viewport's width to determine the container width
	return types.CreateOverlayContainerStyle(h.viewport.Width).Render(content)
}

// Focused returns whether this overlay should handle input
func (h *HelpOverlay) Focused() bool {
	return true
}

// Width returns the overlay width.
func (h *HelpOverlay) Width() int {
	return h.viewport.Width
}

// Height returns the overlay height.
func (h *HelpOverlay) Height() int {
	return h.viewport.Height
}
