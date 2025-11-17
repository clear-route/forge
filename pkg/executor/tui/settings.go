package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/entrhq/forge/pkg/config"
)

const (
	sectionIDAutoApproval     = "auto_approval"
	sectionIDCommandWhitelist = "command_whitelist"
)

// SettingsOverlay displays and manages application settings
type SettingsOverlay struct {
	width   int
	height  int
	focused bool
}

// NewSettingsOverlay creates a new settings overlay
func NewSettingsOverlay(width, height int) *SettingsOverlay {
	return &SettingsOverlay{
		width:   width,
		height:  height,
		focused: true,
	}
}

// Update handles messages for the settings overlay
func (s *SettingsOverlay) Update(msg tea.Msg) (Overlay, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case keyEsc, "q":
			// Close overlay
			return nil, nil
		}
	}

	return s, nil
}

// View renders the settings overlay
func (s *SettingsOverlay) View() string {
	// Get configuration manager
	if !config.IsInitialized() {
		return s.renderError("Configuration not initialized")
	}

	manager := config.Global()

	var content strings.Builder

	// Title using Forge colors
	title := OverlayTitleStyle.Render("⚙️  Settings")
	content.WriteString(title)
	content.WriteString("\n\n")

	// Instructions using Forge subtitle style
	instructions := OverlaySubtitleStyle.Render("Press [Esc] or [q] to close • Settings are saved to ~/.forge/config.json")
	content.WriteString(instructions)
	content.WriteString("\n\n")

	// Render each section
	sections := manager.GetSections()
	for i, section := range sections {
		if i > 0 {
			content.WriteString("\n")
		}
		content.WriteString(s.renderSection(section))
	}

	content.WriteString("\n\n")
	content.WriteString(OverlayHelpStyle.Render("Note: Full settings editor with interactive controls coming soon!"))

	// Use Forge overlay container style
	boxStyle := CreateOverlayContainerStyle(s.width - 4).Height(s.height - 4)

	return lipgloss.Place(
		s.width,
		s.height,
		lipgloss.Center,
		lipgloss.Center,
		boxStyle.Render(content.String()),
	)
}

// renderSection renders a single configuration section
func (s *SettingsOverlay) renderSection(section config.Section) string {
	var out strings.Builder

	// Section title using mint green (success color)
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(mintGreen)

	out.WriteString(titleStyle.Render("▸ " + section.Title()))
	out.WriteString("\n")

	// Section description using muted gray
	if desc := section.Description(); desc != "" {
		descStyle := lipgloss.NewStyle().
			Foreground(mutedGray).
			Italic(true)
		out.WriteString("  ")
		out.WriteString(descStyle.Render(desc))
		out.WriteString("\n")
	}

	// Section data
	data := section.Data()
	if len(data) == 0 {
		out.WriteString("  (no settings configured)\n")
		return out.String()
	}

	// Render based on section type
	switch section.ID() {
	case sectionIDAutoApproval:
		s.renderAutoApprovalData(&out, data)
	case sectionIDCommandWhitelist:
		s.renderCommandWhitelistData(&out, data)
	default:
		// Generic rendering
		for key, value := range data {
			out.WriteString(fmt.Sprintf("  %s: %v\n", key, value))
		}
	}

	return out.String()
}

// renderAutoApprovalData renders auto-approval settings
func (s *SettingsOverlay) renderAutoApprovalData(out *strings.Builder, data map[string]interface{}) {
	if len(data) == 0 {
		out.WriteString("  All available tools discovered from agent.\n")
		out.WriteString("  Edit ~/.forge/config.json to enable auto-approval for specific tools.\n")
		return
	}

	// Count enabled tools
	enabledCount := 0
	for _, v := range data {
		if enabled, ok := v.(bool); ok && enabled {
			enabledCount++
		}
	}

	out.WriteString(fmt.Sprintf("  %d/%d tools auto-approved\n", enabledCount, len(data)))

	// Show first few enabled tools
	shown := 0
	maxShow := 3
	for tool, v := range data {
		if enabled, ok := v.(bool); ok && enabled {
			if shown < maxShow {
				out.WriteString(fmt.Sprintf("    ✓ %s\n", tool))
				shown++
			}
		}
	}
	if enabledCount > maxShow {
		out.WriteString(fmt.Sprintf("    ... and %d more\n", enabledCount-maxShow))
	}
}

// renderCommandWhitelistData renders command whitelist settings
func (s *SettingsOverlay) renderCommandWhitelistData(out *strings.Builder, data map[string]interface{}) {
	patterns, ok := data["patterns"]
	if !ok {
		out.WriteString("  No commands whitelisted\n")
		return
	}

	patternsList, ok := patterns.([]interface{})
	if !ok || len(patternsList) == 0 {
		out.WriteString("  No commands whitelisted\n")
		return
	}

	out.WriteString(fmt.Sprintf("  %d command pattern(s) whitelisted\n", len(patternsList)))

	// Show first few patterns
	maxShow := 3
	for i, p := range patternsList {
		if i >= maxShow {
			out.WriteString(fmt.Sprintf("    ... and %d more\n", len(patternsList)-maxShow))
			break
		}

		if patternMap, ok := p.(map[string]interface{}); ok {
			pattern := patternMap["pattern"]
			desc := patternMap["description"]
			out.WriteString(fmt.Sprintf("    ✓ %s", pattern))
			if desc != nil && desc != "" {
				out.WriteString(fmt.Sprintf(" - %s", desc))
			}
			out.WriteString("\n")
		}
	}
}

// renderError renders an error message
func (s *SettingsOverlay) renderError(message string) string {
	errStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("1")).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("1")).
		Padding(1, 2).
		Width(s.width - 4).
		Height(s.height - 4)

	content := errStyle.Render("Error: ") + message

	return lipgloss.Place(
		s.width,
		s.height,
		lipgloss.Center,
		lipgloss.Center,
		boxStyle.Render(content),
	)
}

// Focused returns whether this overlay should handle input
func (s *SettingsOverlay) Focused() bool {
	return s.focused
}

// Width returns the overlay width
func (s *SettingsOverlay) Width() int {
	return s.width
}

// Height returns the overlay height
func (s *SettingsOverlay) Height() int {
	return s.height
}

// handleSettingsCommand opens the settings overlay
func handleSettingsCommand(m *model, args []string) interface{} {
	// Show interactive settings overlay
	m.overlay.activate(OverlayModeSettings, NewInteractiveSettingsOverlay(m.width, m.height))
	return nil
}
