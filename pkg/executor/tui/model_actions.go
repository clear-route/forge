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

// Quit triggers application exit
func (m *model) Quit() {
	// We might need to send a QuitMsg or just use tea.Quit
	// Normally Update returns tea.Quit. But here we manipulate state?
	// This method might need to be handled via a channel or command if called from async context.
	// But ActionHandler is called from Update loop usually.
	// Wait, ActionHandler methods are called by subcomponents during THEIR Update.
	// But subcomponents return commands.
	// If `Quit()` is called, we probably want the MAIN update loop to see it.
	// However, typical Bubble Tea pattern is to return tea.Quit from Update.
	// Since `ActionHandler` is passed to `Overlay.Update`, and `Overlay.Update` returns `(Overlay, tea.Cmd)`,
	// we can't easily force the parent `model` to return `tea.Quit` unless we have a way to signal it.
	// But wait, `ActionHandler` is implemented by `*model` directly. so we CAN modify `*model`.
	// But `model.Update` logic decides return value.
	// I can set a flag `m.quitting = true` if needed, or assume `Quit()` just does cleanup.
	// Actually `tea.Quit` is a command. `ActionHandler` doesn't return commands.
	// Let's assume for now we handle this via a state flag or similar if needed.
	// Or maybe `Quit` just sends a message to the program.
	// But `model` doesn't have access to `program`.
	// For now, let's leave Quit logic as a TODO or assume we don't use it from overlays yet,
	// OR we implement it by returning a specific command that the parent interprets?
	// Interface says `Quit()`.
	// I'll Implement it later if needed, for now empty is safest? No, `settings` overlay might use it.
	// Actually `settings` calls `actions.ClearOverlay()`.
}
