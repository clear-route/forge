package overlay

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/entrhq/forge/pkg/executor/tui/syntax"
	"github.com/entrhq/forge/pkg/executor/tui/types"
)

// SlashCommandPreview displays a preview of slash command changes before execution
type SlashCommandPreview struct {
	viewport    viewport.Model
	commandName string
	title       string
	files       []string
	message     string
	diff        string
	prTitle     string // PR title (only for PR commands)
	prDesc      string // PR description (only for PR commands)
	selected    ApprovalChoice
	width       int
	height      int
	focused     bool
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

	vp := viewport.New(overlayWidth-4, viewportHeight)
	vp.Style = lipgloss.NewStyle()

	// Build content showing files, message, and diff preview
	content := buildPreviewContent(commandName, files, message, diff, prTitle, prDesc)
	vp.SetContent(content)

	return &SlashCommandPreview{
		viewport:    vp,
		commandName: commandName,
		title:       title,
		files:       files,
		message:     message,
		diff:        diff,
		prTitle:     prTitle,
		prDesc:      prDesc,
		selected:    ApprovalChoiceAccept,
		width:       overlayWidth,
		height:      overlayHeight,
		focused:     true,
		onApprove:   onApprove,
		onReject:    onReject,
	}
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return s.handleKeyMsg(msg, actions)
	case tea.WindowSizeMsg:
		return s.handleWindowResize(msg)
	}
	return s, nil
}

func (s *SlashCommandPreview) handleKeyMsg(msg tea.KeyMsg, actions types.ActionHandler) (types.Overlay, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc", "ctrl+r":
		return s.handleReject(actions)
	case "ctrl+a":
		return s.handleApprove(actions)
	case "tab":
		return s.handleToggleSelection()
	case "enter":
		return s.handleSubmit(actions)
	case "left", "h":
		return s.handleSelectAccept()
	case "right", "l":
		return s.handleSelectReject()
	default:
		return s.handleViewportScroll(msg)
	}
}

func (s *SlashCommandPreview) handleReject(actions types.ActionHandler) (types.Overlay, tea.Cmd) {
	// Close overlay and execute rejection command
	if actions != nil {
		actions.ClearOverlay()
	}
	return nil, s.onReject
}

func (s *SlashCommandPreview) handleApprove(actions types.ActionHandler) (types.Overlay, tea.Cmd) {
	// Close overlay and execute approval command
	if actions != nil {
		actions.ClearOverlay()
	}
	return nil, s.onApprove
}

func (s *SlashCommandPreview) handleToggleSelection() (types.Overlay, tea.Cmd) {
	if s.selected == ApprovalChoiceAccept {
		s.selected = ApprovalChoiceReject
	} else {
		s.selected = ApprovalChoiceAccept
	}
	return s, nil
}

func (s *SlashCommandPreview) handleSubmit(actions types.ActionHandler) (types.Overlay, tea.Cmd) {
	if s.selected == ApprovalChoiceAccept {
		return s.handleApprove(actions)
	}
	return s.handleReject(actions)
}

func (s *SlashCommandPreview) handleSelectAccept() (types.Overlay, tea.Cmd) {
	s.selected = ApprovalChoiceAccept
	return s, nil
}

func (s *SlashCommandPreview) handleSelectReject() (types.Overlay, tea.Cmd) {
	s.selected = ApprovalChoiceReject
	return s, nil
}

func (s *SlashCommandPreview) handleViewportScroll(msg tea.KeyMsg) (types.Overlay, tea.Cmd) {
	var cmd tea.Cmd
	s.viewport, cmd = s.viewport.Update(msg)
	return s, cmd
}

func (s *SlashCommandPreview) handleWindowResize(msg tea.WindowSizeMsg) (types.Overlay, tea.Cmd) {
	var cmd tea.Cmd
	s.width = msg.Width
	s.height = msg.Height
	s.viewport, cmd = s.viewport.Update(msg)
	return s, cmd
}

func (s *SlashCommandPreview) View() string {
	var b strings.Builder

	// Content width accounts for outer container border (2) + padding (4) = 6 chars
	contentWidth := s.width - 6

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
	b.WriteString("\n\n")

	// Content box has its own border (2) + padding (2), so reduce width further
	contentStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(types.SalmonPink).
		Padding(0, 1).
		Width(contentWidth - 4)

	b.WriteString(contentStyle.Render(s.viewport.View()))
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

	// Use shared overlay container style for consistency (width only, height determined by content)
	return types.CreateOverlayContainerStyle(s.width).Render(b.String())
}

func (s *SlashCommandPreview) Focused() bool {
	return s.focused
}

func (s *SlashCommandPreview) Width() int {
	return s.width
}

func (s *SlashCommandPreview) Height() int {
	return s.height
}
