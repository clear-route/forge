package tui

import (
	"strings"

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
	// OverlayModeHelp shows the help overlay
	OverlayModeHelp
	// OverlayModeSlashCommandPreview shows slash command preview
	OverlayModeSlashCommandPreview
	// OverlayModeSettings shows the settings overlay
	OverlayModeSettings
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

// renderToastOverlay renders a toast-style overlay at the bottom of the screen
// without affecting the base view's layout
func renderToastOverlay(baseView string, toastContent string) string {
	if toastContent == "" {
		return baseView
	}

	// Split base view into lines
	baseLines := strings.Split(baseView, "\n")

	// Calculate where to position the toast (bottom of screen, above input area)
	// We want to overlay it on top of the existing content
	toastLines := strings.Split(strings.TrimRight(toastContent, "\n"), "\n")
	toastHeight := len(toastLines)

	// Position toast starting from a few lines above the bottom
	// This puts it just above the input box
	startLine := len(baseLines) - 5 - toastHeight
	if startLine < 0 {
		startLine = 0
	}

	// Build result with toast overlaid
	var result strings.Builder
	for i, line := range baseLines {
		toastLineIdx := i - startLine
		if toastLineIdx >= 0 && toastLineIdx < len(toastLines) {
			// Overlay the toast line, left-aligned with small padding
			toastLine := toastLines[toastLineIdx]
			padding := 2 // Left padding for spacing from edge
			// Write toast with left padding
			result.WriteString(strings.Repeat(" ", padding))
			result.WriteString(toastLine)
		} else {
			result.WriteString(line)
		}
		if i < len(baseLines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}
