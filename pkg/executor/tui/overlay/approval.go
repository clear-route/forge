package overlay

import (
	"github.com/charmbracelet/bubbles/viewport"
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
	viewport viewport.Model
	request  approval.ApprovalRequest
	selected ApprovalChoice
	width    int
	height   int
	focused  bool
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

	vp := viewport.New(overlayWidth-4, viewportHeight)
	vp.Style = lipgloss.NewStyle()
	vp.SetContent(request.Content())

	return &GenericApprovalOverlay{
		viewport: vp,
		request:  request,
		selected: ApprovalChoiceAccept,
		width:    overlayWidth,
		height:   overlayHeight,
		focused:  true,
	}
}

// Update handles messages for the approval overlay
func (a *GenericApprovalOverlay) Update(msg tea.Msg, state types.StateProvider, actions types.ActionHandler) (types.Overlay, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return a.handleKeyMsg(msg)
	case tea.WindowSizeMsg:
		return a.handleWindowResize(msg)
	}
	return a, nil
}

func (a *GenericApprovalOverlay) handleKeyMsg(msg tea.KeyMsg) (types.Overlay, tea.Cmd) {
	switch msg.String() {
	case keyCtrlC, keyEsc, keyCtrlR:
		return a.handleReject()
	case keyCtrlA:
		return a.handleApprove()
	case keyTab:
		return a.handleToggleSelection()
	case keyEnter:
		return a.handleSubmit()
	case keyLeft, "h":
		return a.handleSelectAccept()
	case keyRight, "l":
		return a.handleSelectReject()
	default:
		return a.handleViewportScroll(msg)
	}
}

func (a *GenericApprovalOverlay) handleReject() (types.Overlay, tea.Cmd) {
	// Close overlay and execute rejection command from the request
	return nil, a.request.OnReject()
}

func (a *GenericApprovalOverlay) handleApprove() (types.Overlay, tea.Cmd) {
	// Close overlay and execute approval command from the request
	return nil, a.request.OnApprove()
}

func (a *GenericApprovalOverlay) handleToggleSelection() (types.Overlay, tea.Cmd) {
	if a.selected == ApprovalChoiceAccept {
		a.selected = ApprovalChoiceReject
	} else {
		a.selected = ApprovalChoiceAccept
	}
	return a, nil
}

func (a *GenericApprovalOverlay) handleSubmit() (types.Overlay, tea.Cmd) {
	if a.selected == ApprovalChoiceAccept {
		return a.handleApprove()
	}
	return a.handleReject()
}

func (a *GenericApprovalOverlay) handleSelectAccept() (types.Overlay, tea.Cmd) {
	a.selected = ApprovalChoiceAccept
	return a, nil
}

func (a *GenericApprovalOverlay) handleSelectReject() (types.Overlay, tea.Cmd) {
	a.selected = ApprovalChoiceReject
	return a, nil
}

func (a *GenericApprovalOverlay) handleViewportScroll(msg tea.KeyMsg) (types.Overlay, tea.Cmd) {
	var cmd tea.Cmd
	a.viewport, cmd = a.viewport.Update(msg)
	return a, cmd
}

func (a *GenericApprovalOverlay) handleWindowResize(msg tea.WindowSizeMsg) (types.Overlay, tea.Cmd) {
	a.width = msg.Width
	a.height = msg.Height
	a.viewport.Width = min(76, a.width-4)
	a.viewport.Height = min(20, a.height-10)
	return a, nil
}

// View renders the approval overlay
func (a *GenericApprovalOverlay) View() string {
	title := types.OverlayTitleStyle.Render(a.request.Title())
	viewportContent := a.viewport.View()

	// Render approval buttons using style functions
	acceptStyle := types.GetAcceptButtonStyle(a.selected == ApprovalChoiceAccept)
	rejectStyle := types.GetRejectButtonStyle(a.selected == ApprovalChoiceReject)

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Left,
		acceptStyle.Render(" ✓ Accept "),
		"  ",
		rejectStyle.Render(" ✗ Reject "),
	)

	hints := types.OverlayHelpStyle.Render("Ctrl+A: Accept • Ctrl+R: Reject • Tab: Toggle • ↑/↓: Scroll")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		viewportContent,
		"",
		buttons,
		hints,
	)

	return types.CreateOverlayContainerStyle(a.width - 4).Render(content)
}

// Focused returns whether this overlay should handle input
func (a *GenericApprovalOverlay) Focused() bool {
	return a.focused
}

// Width returns the overlay width
func (a *GenericApprovalOverlay) Width() int {
	return a.width
}

// Height returns the overlay height
func (a *GenericApprovalOverlay) Height() int {
	return a.height
}
