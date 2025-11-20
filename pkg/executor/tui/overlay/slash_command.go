package overlay

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/entrhq/forge/pkg/executor/tui/syntax"
	"github.com/entrhq/forge/pkg/executor/tui/types"
)

// SlashCommandPreview displays a preview of slash command changes before execution
type SlashCommandPreview struct {
	*BaseOverlay
	commandName string
	title       string
	files       []string
	message     string
	diff        string
	prTitle     string // PR title (only for PR commands)
	prDesc      string // PR description (only for PR commands)
	selected    ApprovalChoice
	onApprove   tea.Cmd // Command to execute on approval
	onReject    tea.Cmd // Command to execute on rejection
}

// NewSlashCommandPreview creates a new slash command preview overlay
func NewSlashCommandPreview(commandName, title string, files []string, message, diff, prTitle, prDesc string, width, height int, onApprove, onReject tea.Cmd) *SlashCommandPreview {
	// Make overlay wide - 90% of screen width
	overlayWidth := max(int(float64(width)*0.9), 80)

	// Fixed viewport height for content
	const maxViewportHeight = 15
	viewportHeight := maxViewportHeight

	// Calculate total overlay height
	// Title (2) + subtitle (1) + spacing (1) + border (2) + buttons (2) + hints (1) = 9 lines
	// Plus viewport height
	overlayHeight := viewportHeight + 9

	overlay := &SlashCommandPreview{
		commandName: commandName,
		title:       title,
		files:       files,
		message:     message,
		diff:        diff,
		prTitle:     prTitle,
		prDesc:      prDesc,
		selected:    ApprovalChoiceAccept,
		onApprove:   onApprove,
		onReject:    onReject,
	}

	// Build content showing files, message, and diff preview
	content := buildPreviewContent(commandName, files, message, diff, prTitle, prDesc)

	// Configure base overlay
	baseConfig := BaseOverlayConfig{
		Width:          overlayWidth,
		Height:         overlayHeight,
		ViewportWidth:  overlayWidth - 4,
		ViewportHeight: viewportHeight,
		Content:        content,
		OnClose: func(actions types.ActionHandler) tea.Cmd {
			if actions != nil {
				actions.ClearOverlay()
			}
			return overlay.onReject
		},
		OnCustomKey: func(msg tea.KeyMsg, actions types.ActionHandler) (bool, tea.Cmd) {
			return overlay.handleCustomKeys(msg, actions)
		},
		RenderHeader: overlay.renderHeader,
		RenderFooter: overlay.renderFooter,
	}

	overlay.BaseOverlay = NewBaseOverlay(baseConfig)
	return overlay
}

// buildPreviewContent creates the content to display in the viewport
func buildPreviewContent(commandName string, files []string, message, diff, prTitle, prDesc string) string {
	var b strings.Builder

	if commandName == "pr" {
		// PR-specific layout

		// Show branch info
		if len(files) > 0 {
			b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(types.SalmonPink).Render("Branch:"))
			b.WriteString("\n")
			for _, file := range files {
				b.WriteString("  " + file + "\n")
			}
			b.WriteString("\n")
		}

		// Show PR title
		if prTitle != "" {
			b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(types.SalmonPink).Render("PR Title:"))
			b.WriteString("\n")
			b.WriteString(prTitle)
			b.WriteString("\n\n")
		}

		// Show PR description
		if prDesc != "" {
			b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(types.SalmonPink).Render("PR Description:"))
			b.WriteString("\n")
			b.WriteString(prDesc)
			b.WriteString("\n\n")
		}

		// Show commits and changes
		if diff != "" {
			b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(types.SalmonPink).Render("Commits & Changes:"))
			b.WriteString("\n")
			b.WriteString(diff)
		}
	} else {
		// Commit-specific layout

		// Show files to commit
		if len(files) > 0 {
			b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(types.SalmonPink).Render("Files to commit:"))
			b.WriteString("\n")
			for _, file := range files {
				b.WriteString("  • " + file + "\n")
			}
			b.WriteString("\n")
		}

		// Show commit message
		if message != "" {
			b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(types.SalmonPink).Render("Commit Message:"))
			b.WriteString("\n")
			b.WriteString(message)
			b.WriteString("\n\n")
		}

		// Show diff with syntax highlighting
		if diff != "" {
			b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(types.SalmonPink).Render("Changes:"))
			b.WriteString("\n")

			highlightedDiff, err := syntax.HighlightDiff(diff, "")
			if err != nil {
				b.WriteString(diff)
			} else {
				b.WriteString(highlightedDiff)
			}
		}
	}

	return b.String()
}

func (s *SlashCommandPreview) Update(msg tea.Msg, state types.StateProvider, actions types.ActionHandler) (types.Overlay, tea.Cmd) {
	// Let BaseOverlay handle standard messages
	handled, updatedBase, cmd := s.BaseOverlay.Update(msg, actions)
	s.BaseOverlay = updatedBase

	if handled {
		return s, cmd
	}

	return s, nil
}

// handleCustomKeys processes slash command-specific key presses
func (s *SlashCommandPreview) handleCustomKeys(msg tea.KeyMsg, actions types.ActionHandler) (bool, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc", "ctrl+r":
		// Close overlay before executing reject command
		if actions != nil {
			actions.ClearOverlay()
		}
		return true, s.onReject
	case "ctrl+a":
		// Close overlay before executing approve command
		if actions != nil {
			actions.ClearOverlay()
		}
		return true, s.onApprove
	case "tab":
		s.handleToggleSelection()
		return true, nil
	case "enter":
		// Close overlay before executing selected command
		if actions != nil {
			actions.ClearOverlay()
		}
		if s.selected == ApprovalChoiceAccept {
			return true, s.onApprove
		}
		return true, s.onReject
	case "left", "h":
		s.selected = ApprovalChoiceAccept
		return true, nil
	case "right", "l":
		s.selected = ApprovalChoiceReject
		return true, nil
	}
	return false, nil
}

func (s *SlashCommandPreview) handleToggleSelection() {
	if s.selected == ApprovalChoiceAccept {
		s.selected = ApprovalChoiceReject
	} else {
		s.selected = ApprovalChoiceAccept
	}
}

// renderHeader renders the slash command preview header
func (s *SlashCommandPreview) renderHeader() string {
	var b strings.Builder

	// Content width accounts for container styling
	contentWidth := s.Width() - 6

	title := s.title
	subtitle := fmt.Sprintf("/%s", s.commandName)

	// Manually center by calculating padding
	titleLen := len(title)
	subtitleLen := len(subtitle)
	titlePadding := max(0, (contentWidth-titleLen)/2)
	subtitlePadding := max(0, (contentWidth-subtitleLen)/2)

	b.WriteString(strings.Repeat(" ", titlePadding) + types.OverlayTitleStyle.Render(title))
	b.WriteString("\n")
	b.WriteString(strings.Repeat(" ", subtitlePadding) + types.OverlaySubtitleStyle.Render(subtitle))

	return b.String()
}

// renderFooter renders the viewport in a bordered box with buttons and hints
func (s *SlashCommandPreview) renderFooter() string {
	var b strings.Builder

	contentWidth := s.Width() - 6

	// Content box has its own border (2) + padding (2), so reduce width further
	contentStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(types.SalmonPink).
		Padding(0, 1).
		Width(contentWidth - 4)

	b.WriteString(contentStyle.Render(s.BaseOverlay.Viewport().View()))
	b.WriteString("\n\n")

	// Use shared button styles for consistency
	acceptStyle := types.GetAcceptButtonStyle(s.selected == ApprovalChoiceAccept)
	rejectStyle := types.GetRejectButtonStyle(s.selected == ApprovalChoiceReject)

	acceptBtn := acceptStyle.Render("✓ Execute (Enter / Ctrl+A)")
	rejectBtn := rejectStyle.Render("✗ Cancel (Esc / Ctrl+R)")

	// Use shared spacer utility for consistency
	spacer := types.CreateStyledSpacer(2)

	// Join buttons with styled spacer
	buttonsRow := acceptBtn + spacer + rejectBtn

	// Manually center buttons
	buttonsLen := lipgloss.Width(buttonsRow)
	buttonsPadding := max(0, (contentWidth-buttonsLen)/2)

	b.WriteString(strings.Repeat(" ", buttonsPadding) + buttonsRow)
	b.WriteString("\n")

	hints := "↑↓ to scroll • ← → Tab to choose • Enter to submit"
	hintsLen := len(hints)
	hintsPadding := max(0, (contentWidth-hintsLen)/2)

	b.WriteString(strings.Repeat(" ", hintsPadding) + types.OverlayHelpStyle.Render(hints))

	return b.String()
}

// View renders the overlay
func (s *SlashCommandPreview) View() string {
	return s.BaseOverlay.View(s.Width())
}
