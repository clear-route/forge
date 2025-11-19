package types

import (
	tea "github.com/charmbracelet/bubbletea"
)

// StateProvider defines the read-only access to the application state
// required by subcomponents (overlays, commands).
type StateProvider interface {
	// Dimensions
	GetWidth() int
	GetHeight() int

	// Agent State
	IsThinking() bool
	IsAgentBusy() bool
	GetWorkspaceDir() string

	// UI State
	IsBashMode() bool
}

// ActionHandler defines the actions that subcomponents can trigger
// on the main application controller.
type ActionHandler interface {
	// UI Actions
	SetOverlay(mode OverlayMode, overlay Overlay)
	ClearOverlay()
	ShowToast(message, details, icon string, isError bool)

	// Input Actions
	SetInput(value string)
	SetCursorEnd()

	// System Actions
	Quit()
}

// Overlay is the interface that all overlay components must implement.
// It allows the main update loop to delegate event handling to the active overlay.
type Overlay interface {
	// Update handles messages and returns updated overlay.
	// Note: In the new architecture, this receives the StateProvider instead of *model
	Update(msg tea.Msg, state StateProvider, actions ActionHandler) (Overlay, tea.Cmd)

	// View renders the overlay
	View() string

	// Focused returns whether this overlay should handle input
	Focused() bool

	// Width returns the overlay width
	Width() int

	// Height returns the overlay height
	Height() int
}
