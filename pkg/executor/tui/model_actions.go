package tui

import (
	"time"

	"github.com/entrhq/forge/pkg/executor/tui/types"
)

// SetOverlay activates an overlay
func (m *model) SetOverlay(mode types.OverlayMode, overlay types.Overlay) {
	m.overlay.activate(mode, overlay)
}

// ClearOverlay closes the current overlay
func (m *model) ClearOverlay() {
	m.overlay.deactivate()
	// Refocus textarea if needed? The View logic handles focus.
	m.textarea.Focus()
}

// ShowToast displays a toast notification
func (m *model) ShowToast(message, details, icon string, isError bool) {
	m.toast = &toastNotification{
		active:    true,
		message:   message,
		details:   details,
		icon:      icon,
		isError:   isError,
		showUntil: time.Now().Add(5 * time.Second),
	}
}

// SetInput sets the textarea content
func (m *model) SetInput(value string) {
	m.textarea.SetValue(value)
	m.updateTextAreaHeight()
}

// SetCursorEnd moves the cursor to the end of input
func (m *model) SetCursorEnd() {
	m.textarea.CursorEnd()
}

// Quit triggers application exit by setting a flag that will be checked in the Update loop.
// This allows overlays and other components to request app termination without directly
// returning tea.Quit (which would break the Bubble Tea command chain).
func (m *model) Quit() {
	m.shouldQuit = true
}
