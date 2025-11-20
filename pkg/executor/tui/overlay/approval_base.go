package overlay

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/entrhq/forge/pkg/executor/tui/types"
)

// ApprovalOverlayBase provides common functionality for approval-style overlays.
// This includes accept/reject button handling, selection toggling, and standard approval UI.
type ApprovalOverlayBase struct {
	*BaseOverlay
	selected      ApprovalChoice
	onApprove     func() tea.Cmd
	onReject      func() tea.Cmd
	approveLabel  string
	rejectLabel   string
	showHints     bool
	customButtons func(selected ApprovalChoice) string
}

// ApprovalOverlayConfig configures an approval overlay
type ApprovalOverlayConfig struct {
	BaseConfig    BaseOverlayConfig
	OnApprove     func() tea.Cmd
	OnReject      func() tea.Cmd
	ApproveLabel  string // Default: "✓ Accept (Enter / Ctrl+A)"
	RejectLabel   string // Default: "✗ Reject (Esc / Ctrl+R)"
	ShowHints     bool   // Default: true
	CustomButtons func(selected ApprovalChoice) string
}

// NewApprovalOverlayBase creates a new approval overlay base
func NewApprovalOverlayBase(config ApprovalOverlayConfig) *ApprovalOverlayBase {
	// Set default labels
	if config.ApproveLabel == "" {
		config.ApproveLabel = "✓ Accept (Enter / Ctrl+A)"
	}
	if config.RejectLabel == "" {
		config.RejectLabel = "✗ Reject (Esc / Ctrl+R)"
	}

	// Configure custom key handler for approval-specific keys
	config.BaseConfig.OnCustomKey = func(msg tea.KeyMsg, actions types.ActionHandler) (bool, tea.Cmd) {
		return false, nil // Will be set after creation
	}

	base := NewBaseOverlay(config.BaseConfig)

	approval := &ApprovalOverlayBase{
		BaseOverlay:   base,
		selected:      ApprovalChoiceAccept,
		onApprove:     config.OnApprove,
		onReject:      config.OnReject,
		approveLabel:  config.ApproveLabel,
		rejectLabel:   config.RejectLabel,
		showHints:     config.ShowHints,
		customButtons: config.CustomButtons,
	}

	// Set the custom key handler now that we have the approval instance
	base.onCustomKey = approval.handleApprovalKeys

	return approval
}

// Update handles approval overlay updates
func (a *ApprovalOverlayBase) Update(msg tea.Msg, state types.StateProvider, actions types.ActionHandler) (*ApprovalOverlayBase, tea.Cmd) {
	handled, updatedBase, cmd := a.BaseOverlay.Update(msg, actions)
	a.BaseOverlay = updatedBase

	if handled {
		return a, cmd
	}

	return a, nil
}

// handleApprovalKeys processes approval-specific keyboard input
func (a *ApprovalOverlayBase) handleApprovalKeys(msg tea.KeyMsg, actions types.ActionHandler) (bool, tea.Cmd) {
	switch msg.String() {
	case keyCtrlA:
		return true, a.approve()
	case keyCtrlR:
		return true, a.reject()
	case keyTab:
		a.toggleSelection()
		return true, nil
	case keyEnter:
		return true, a.submit()
	case keyLeft, "h":
		a.selectAccept()
		return true, nil
	case keyRight, "l":
		a.selectReject()
		return true, nil
	}
	return false, nil
}

// approve executes the approval action
func (a *ApprovalOverlayBase) approve() tea.Cmd {
	if a.onApprove != nil {
		return a.onApprove()
	}
	return nil
}

// reject executes the rejection action
func (a *ApprovalOverlayBase) reject() tea.Cmd {
	if a.onReject != nil {
		return a.onReject()
	}
	return nil
}

// toggleSelection toggles between Accept and Reject
func (a *ApprovalOverlayBase) toggleSelection() {
	if a.selected == ApprovalChoiceAccept {
		a.selected = ApprovalChoiceReject
	} else {
		a.selected = ApprovalChoiceAccept
	}
}

// submit submits the currently selected choice
func (a *ApprovalOverlayBase) submit() tea.Cmd {
	if a.selected == ApprovalChoiceAccept {
		return a.approve()
	}
	return a.reject()
}

// selectAccept selects the Accept option
func (a *ApprovalOverlayBase) selectAccept() {
	a.selected = ApprovalChoiceAccept
}

// selectReject selects the Reject option
func (a *ApprovalOverlayBase) selectReject() {
	a.selected = ApprovalChoiceReject
}

// RenderButtons renders the approval buttons
func (a *ApprovalOverlayBase) RenderButtons() string {
	if a.customButtons != nil {
		return a.customButtons(a.selected)
	}

	acceptStyle := types.GetAcceptButtonStyle(a.selected == ApprovalChoiceAccept)
	rejectStyle := types.GetRejectButtonStyle(a.selected == ApprovalChoiceReject)

	acceptBtn := acceptStyle.Render(a.approveLabel)
	rejectBtn := rejectStyle.Render(a.rejectLabel)

	spacer := types.CreateStyledSpacer(2)

	return acceptBtn + spacer + rejectBtn
}

// RenderHints renders the keyboard hints
func (a *ApprovalOverlayBase) RenderHints() string {
	if !a.showHints {
		return ""
	}
	hints := "↑↓ to scroll • ← → Tab to choose • Enter to submit"
	return types.OverlayHelpStyle.Render(hints)
}

// Selected returns the currently selected choice
func (a *ApprovalOverlayBase) Selected() ApprovalChoice {
	return a.selected
}

// SetSelected sets the selected choice
func (a *ApprovalOverlayBase) SetSelected(choice ApprovalChoice) {
	a.selected = choice
}
