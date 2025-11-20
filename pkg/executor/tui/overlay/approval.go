package overlay

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/entrhq/forge/pkg/executor/tui/approval"
	"github.com/entrhq/forge/pkg/executor/tui/types"
)

// Key binding constants for approval overlay
const (
	keyCtrlA = "ctrl+a"
	keyCtrlC = "ctrl+c"
	keyCtrlR = "ctrl+r"
	keyTab   = "tab"
	keyEnter = "enter"
	keyLeft  = "left"
	keyRight = "right"
	keyEsc   = "esc"
)

// GenericApprovalOverlay displays an approval request for any command.
// It is completely agnostic of the specific command being approved.
type GenericApprovalOverlay struct {
	*ApprovalOverlayBase
	request approval.ApprovalRequest
}

// NewGenericApprovalOverlay creates a new generic approval overlay
func NewGenericApprovalOverlay(request approval.ApprovalRequest, width, height int) *GenericApprovalOverlay {
	// Make overlay wide - 90% of screen width
	overlayWidth := max(int(float64(width)*0.9), 80)

	// Fixed viewport height for content
	const maxViewportHeight = 15
	viewportHeight := maxViewportHeight

	// Calculate total overlay height
	// Title (2) + subtitle (1) + spacing (1) + border (2) + buttons (2) + hints (1) = 9 lines
	// Plus viewport height
	overlayHeight := viewportHeight + 9

	overlay := &GenericApprovalOverlay{
		request: request,
	}

	// Configure approval overlay
	approvalConfig := ApprovalOverlayConfig{
		BaseConfig: BaseOverlayConfig{
			Width:          overlayWidth,
			Height:         overlayHeight,
			ViewportWidth:  overlayWidth - 4,
			ViewportHeight: viewportHeight,
			Content:        request.Content(),
			RenderHeader:   overlay.renderHeader,
		},
		OnApprove:    request.OnApprove,
		OnReject:     request.OnReject,
		ApproveLabel: " ✓ Accept ",
		RejectLabel:  " ✗ Reject ",
		ShowHints:    true,
	}

	overlay.ApprovalOverlayBase = NewApprovalOverlayBase(approvalConfig)
	return overlay
}

// Update handles messages for the approval overlay
func (a *GenericApprovalOverlay) Update(msg tea.Msg, state types.StateProvider, actions types.ActionHandler) (types.Overlay, tea.Cmd) {
	updatedApproval, cmd := a.ApprovalOverlayBase.Update(msg, state, actions)
	a.ApprovalOverlayBase = updatedApproval
	return a, cmd
}

// renderHeader renders the approval overlay header
func (a *GenericApprovalOverlay) renderHeader() string {
	return types.OverlayTitleStyle.Render(a.request.Title())
}

// View renders the approval overlay
func (a *GenericApprovalOverlay) View() string {
	var sections []string

	// Header
	sections = append(sections, a.renderHeader())
	sections = append(sections, "")

	// Viewport content
	sections = append(sections, a.Viewport().View())
	sections = append(sections, "")

	// Buttons
	sections = append(sections, a.RenderButtons())

	// Hints - custom hints for generic approval
	hints := types.OverlayHelpStyle.Render("Ctrl+A: Accept • Ctrl+R: Reject • Tab: Toggle • ↑/↓: Scroll")
	sections = append(sections, hints)

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return types.CreateOverlayContainerStyle(a.Width() - 4).Render(content)
}
