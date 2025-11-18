package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// CommandPalette manages command suggestions and selection
type CommandPalette struct {
	commands         []*SlashCommand
	filteredCommands []*SlashCommand
	selectedIndex    int
	filter           string
	active           bool
}

// newCommandPalette creates a new command palette
func newCommandPalette() *CommandPalette {
	return &CommandPalette{
		commands:         getAllCommands(),
		filteredCommands: getAllCommands(),
		selectedIndex:    0,
		active:           false,
	}
}

// activate shows the command palette
func (cp *CommandPalette) activate() {
	cp.active = true
	cp.filter = ""
	cp.selectedIndex = 0
	cp.updateFiltered()
}

// deactivate hides the command palette
func (cp *CommandPalette) deactivate() {
	cp.active = false
	cp.filter = ""
	cp.selectedIndex = 0
}

// updateFilter updates the filter string and refreshes filtered commands
func (cp *CommandPalette) updateFilter(filter string) {
	newFilter := strings.ToLower(strings.TrimSpace(filter))
	// Only reset selection if the filter actually changed
	if newFilter != cp.filter {
		cp.filter = newFilter
		cp.selectedIndex = 0
		cp.updateFiltered()
	}
}

// updateFiltered updates the list of filtered commands based on current filter
func (cp *CommandPalette) updateFiltered() {
	if cp.filter == "" {
		cp.filteredCommands = cp.commands
		return
	}

	filtered := make([]*SlashCommand, 0)
	for _, cmd := range cp.commands {
		// Match on command name or description
		if strings.Contains(strings.ToLower(cmd.Name), cp.filter) ||
			strings.Contains(strings.ToLower(cmd.Description), cp.filter) {
			filtered = append(filtered, cmd)
		}
	}
	cp.filteredCommands = filtered

	// Ensure selected index is valid
	if cp.selectedIndex >= len(cp.filteredCommands) {
		cp.selectedIndex = len(cp.filteredCommands) - 1
	}
	if cp.selectedIndex < 0 {
		cp.selectedIndex = 0
	}
}

// selectNext moves selection down
func (cp *CommandPalette) selectNext() {
	if len(cp.filteredCommands) == 0 {
		return
	}
	cp.selectedIndex = (cp.selectedIndex + 1) % len(cp.filteredCommands)
}

// selectPrev moves selection up
func (cp *CommandPalette) selectPrev() {
	if len(cp.filteredCommands) == 0 {
		return
	}
	cp.selectedIndex--
	if cp.selectedIndex < 0 {
		cp.selectedIndex = len(cp.filteredCommands) - 1
	}
}

// getSelected returns the currently selected command
func (cp *CommandPalette) getSelected() *SlashCommand {
	if len(cp.filteredCommands) == 0 || cp.selectedIndex >= len(cp.filteredCommands) {
		return nil
	}
	return cp.filteredCommands[cp.selectedIndex]
}

// render renders the command palette
func (cp *CommandPalette) render(width int) string {
	if !cp.active || len(cp.filteredCommands) == 0 {
		return ""
	}

	var sb strings.Builder

	// Calculate palette width (80% of screen or max 80 chars)
	paletteWidth := width * 80 / 100
	if paletteWidth > 80 {
		paletteWidth = 80
	}
	if paletteWidth < 40 {
		paletteWidth = 40
	}

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(salmonPink).
		Bold(true).
		PaddingLeft(1)

	sb.WriteString(headerStyle.Render("Available Commands:"))
	sb.WriteString("\n")

	// Show up to 5 commands
	maxVisible := 5
	if len(cp.filteredCommands) < maxVisible {
		maxVisible = len(cp.filteredCommands)
	}

	for i := 0; i < maxVisible; i++ {
		cmd := cp.filteredCommands[i]
		prefix := "  "
		if i == cp.selectedIndex {
			prefix = "> "
		}

		// Command name in salmon pink, description in soft gray
		cmdNameStyle := lipgloss.NewStyle().
			Foreground(salmonPink).
			Bold(i == cp.selectedIndex)

		descStyle := lipgloss.NewStyle().
			Foreground(mutedGray)

		if i == cp.selectedIndex {
			// Highlighted background for selected item
			lineStyle := lipgloss.NewStyle().
				Background(lipgloss.Color("#2d2d2d")).
				Width(paletteWidth - 2).
				PaddingLeft(1)

			line := prefix + cmdNameStyle.Render("/"+cmd.Name) + "  " + descStyle.Render(cmd.Description)
			sb.WriteString(lineStyle.Render(line))
		} else {
			line := prefix + cmdNameStyle.Render("/"+cmd.Name) + "  " + descStyle.Render(cmd.Description)
			sb.WriteString(line)
		}
		sb.WriteString("\n")
	}

	// Footer hint
	if len(cp.filteredCommands) > maxVisible {
		footerStyle := lipgloss.NewStyle().
			Foreground(mutedGray).
			Italic(true).
			PaddingLeft(1)
		sb.WriteString(footerStyle.Render("... and more. Keep typing to filter."))
		sb.WriteString("\n")
	}

	// Wrap in border
	paletteStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(salmonPink).
		Width(paletteWidth).
		Padding(0, 1)

	return paletteStyle.Render(sb.String())
}
