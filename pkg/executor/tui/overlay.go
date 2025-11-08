package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// OverlayMode represents the current overlay state
type OverlayMode int

const (
	// OverlayModeNone indicates no overlay is active
	OverlayModeNone OverlayMode = iota
	// OverlayModeDiffViewer shows the diff approval overlay
	OverlayModeDiffViewer
	// OverlayModeFileTree shows the file tree overlay
	OverlayModeFileTree
	// OverlayModeCommandOutput shows command output overlay
	OverlayModeCommandOutput
)

// Overlay is the base interface for all overlay components
type Overlay interface {
	// Update handles messages and returns updated overlay
	Update(msg tea.Msg) (Overlay, tea.Cmd)

	// View renders the overlay
	View() string

	// Focused returns whether this overlay should handle input
	Focused() bool

	// Width returns the overlay width
	Width() int

	// Height returns the overlay height
	Height() int
}

// overlayState tracks the active overlay and its state
type overlayState struct {
	mode    OverlayMode
	overlay Overlay
}

// newOverlayState creates a new overlay state
func newOverlayState() *overlayState {
	return &overlayState{
		mode: OverlayModeNone,
	}
}

// activate activates an overlay
func (o *overlayState) activate(mode OverlayMode, overlay Overlay) {
	o.mode = mode
	o.overlay = overlay
}

// deactivate closes the current overlay
func (o *overlayState) deactivate() {
	o.mode = OverlayModeNone
	o.overlay = nil
}

// isActive returns whether any overlay is currently active
func (o *overlayState) isActive() bool {
	return o.mode != OverlayModeNone && o.overlay != nil
}

// renderOverlay renders an overlay on top of the base content
func renderOverlay(baseView string, overlay Overlay, width, height int) string {
	if overlay == nil {
		return baseView
	}

	// Center the overlay
	overlayView := overlay.View()

	// Position the overlay
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		overlayView,
		lipgloss.WithWhitespaceChars(" "),
	)
}
