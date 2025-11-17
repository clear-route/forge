package tui

import "github.com/charmbracelet/lipgloss"

// Color Palette
// This is the single source of truth for all TUI colors.
// Use these constants throughout the TUI to ensure visual consistency.
var (
	// Primary Colors - Core brand colors
	salmonPink  = lipgloss.Color("#FFB3BA") // Soft pastel salmon pink - primary accent
	coralPink   = lipgloss.Color("#FFCCCB") // Lighter coral accent - secondary
	mintGreen   = lipgloss.Color("#A8E6CF") // Soft mint green - success/accept states
	mutedGray   = lipgloss.Color("#6B7280") // Muted gray - secondary text
	brightWhite = lipgloss.Color("#F9FAFB") // Bright white - primary text
	darkBg      = lipgloss.Color("#111827") // Dark background - container backgrounds

	// Semantic Colors - For specific UI states
	black    = lipgloss.Color("#000000") // Black - high contrast text on colored backgrounds
	softGray = lipgloss.Color("#9CA3AF") // Soft gray - for subtle text and descriptions

	// Diff Colors - For code diffs and syntax highlighting
	diffAddColor      = lipgloss.Color("#90EE90") // Green for additions
	diffDeleteColor   = lipgloss.Color("#FFB3BA") // Red for deletions (matches salmonPink)
	diffHunkColor     = lipgloss.Color("#87CEEB") // Cyan for hunk headers
	diffHeaderColor   = lipgloss.Color("#FFA07A") // Orange for file headers
	diffAddBgColor    = lipgloss.Color("#2d4a2b") // Dark green background for added lines
	diffDeleteBgColor = lipgloss.Color("#4a2d2d") // Dark red background for deleted lines
)

// Common Styles
// These are pre-configured styles for common UI elements.
// Use these as base styles and customize as needed.
var (
	// Text Styles
	headerStyle = lipgloss.NewStyle().
			Foreground(salmonPink).
			Bold(true)

	tipsStyle = lipgloss.NewStyle().
			Foreground(mutedGray)

	userStyle = lipgloss.NewStyle().
			Foreground(coralPink).
			Bold(true)

	thinkingStyle = lipgloss.NewStyle().
			Foreground(mutedGray).
			Italic(true)

	toolStyle = lipgloss.NewStyle().
			Foreground(mintGreen)

	toolResultStyle = lipgloss.NewStyle().
			Foreground(brightWhite)

	errorStyle = lipgloss.NewStyle().
			Foreground(salmonPink)

	// Container Styles
	statusBarStyle = lipgloss.NewStyle().
			Foreground(mutedGray).
			Background(darkBg).
			Padding(0, 1)

	inputBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(salmonPink).
			Padding(0, 1)
)

// Button Styles
// CreateButtonStyle creates a button style with the given foreground and background colors.
// Use this helper to create consistent button styling across the TUI.
func CreateButtonStyle(fg, bg lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Padding(0, 2).
		Foreground(fg).
		Background(bg)
}

// GetAcceptButtonStyle returns the style for an accept button based on selection state.
func GetAcceptButtonStyle(selected bool) lipgloss.Style {
	if selected {
		return CreateButtonStyle(black, mintGreen)
	}
	return CreateButtonStyle(mutedGray, darkBg)
}

// GetRejectButtonStyle returns the style for a reject button based on selection state.
func GetRejectButtonStyle(selected bool) lipgloss.Style {
	if selected {
		return CreateButtonStyle(black, salmonPink)
	}
	return CreateButtonStyle(mutedGray, darkBg)
}

// CreateStyledSpacer creates a spacer with the dark background color.
// Use this to create gaps between UI elements that match the container background.
func CreateStyledSpacer(width int) string {
	spacerStyle := lipgloss.NewStyle().Background(darkBg)
	return spacerStyle.Render(lipgloss.PlaceHorizontal(width, lipgloss.Left, ""))
}

// Overlay Container Styles
// CreateOverlayContainerStyle creates a standardized container style for all overlays.
// This ensures consistent appearance across all overlay types (diff viewer, command execution, etc.)
// Note: Only sets width, not height, to allow content to determine the container height naturally.
func CreateOverlayContainerStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(salmonPink).
		Background(darkBg).
		Padding(1, 2).
		Width(width)
}

// Shared text styles for overlay content
var (
	// OverlayTitleStyle is used for main overlay titles
	OverlayTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(salmonPink)

	// OverlaySubtitleStyle is used for overlay subtitles and secondary text
	OverlaySubtitleStyle = lipgloss.NewStyle().
				Foreground(mutedGray)

	// OverlayHelpStyle is used for help text and hints
	OverlayHelpStyle = lipgloss.NewStyle().
				Foreground(mutedGray).
				Italic(true)
)
