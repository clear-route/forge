package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/entrhq/forge/pkg/agent/tools"
	"github.com/entrhq/forge/pkg/types"
)

type ApprovalChoice int

const (
	ApprovalChoiceAccept ApprovalChoice = iota
	ApprovalChoiceReject
)

type DiffViewer struct {
	viewport     viewport.Model
	approvalID   string
	toolName     string
	preview      *tools.ToolPreview
	selected     ApprovalChoice
	width        int
	height       int
	focused      bool
	responseFunc func(*types.ApprovalResponse)
}

func NewDiffViewer(approvalID, toolName string, preview *tools.ToolPreview, width, height int, responseFunc func(*types.ApprovalResponse)) *DiffViewer {
	// Make overlay wide - 90% of screen width
	overlayWidth := max(int(float64(width)*0.9), 80)

	// Fixed viewport height: max 10 lines for diff content
	const maxViewportHeight = 10
	viewportHeight := maxViewportHeight

	// Calculate total overlay height
	// Title (2) + subtitle (1) + spacing (1) + border (2) + buttons (2) + hints (1) = 9 lines
	// Plus viewport height
	overlayHeight := viewportHeight + 9

	vp := viewport.New(overlayWidth-4, viewportHeight)
	vp.Style = lipgloss.NewStyle()

	// Apply syntax highlighting to the diff content
	content := ""
	if preview != nil {
		// Extract language from metadata
		language := ""
		if preview.Metadata != nil {
			if lang, ok := preview.Metadata["language"].(string); ok {
				language = lang
			}
		}

		// Apply syntax highlighting
		highlightedContent, err := HighlightDiff(preview.Content, language)
		if err != nil {
			// Fall back to original content if highlighting fails
			content = preview.Content
		} else {
			content = highlightedContent
		}
		vp.SetContent(content)
	}

	return &DiffViewer{
		viewport:     vp,
		approvalID:   approvalID,
		toolName:     toolName,
		preview:      preview,
		selected:     ApprovalChoiceAccept,
		width:        overlayWidth,
		height:       overlayHeight,
		focused:      true,
		responseFunc: responseFunc,
	}
}

func (d *DiffViewer) Update(msg tea.Msg) (Overlay, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return d.handleKeyMsg(msg)
	case tea.WindowSizeMsg:
		return d.handleWindowResize(msg)
	}
	return d, nil
}

// handleKeyMsg processes keyboard input for the diff viewer
func (d *DiffViewer) handleKeyMsg(msg tea.KeyMsg) (Overlay, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc", "ctrl+r":
		return d.handleReject()
	case "ctrl+a":
		return d.handleApprove()
	case "tab":
		return d.handleToggleSelection()
	case "enter":
		return d.handleSubmit()
	case "left", "h":
		return d.handleSelectAccept()
	case "right", "l":
		return d.handleSelectReject()
	default:
		return d.handleViewportScroll(msg)
	}
}

// handleReject sends a rejection response
func (d *DiffViewer) handleReject() (Overlay, tea.Cmd) {
	if d.responseFunc != nil {
		d.responseFunc(types.NewApprovalResponse(d.approvalID, types.ApprovalRejected))
	}
	return d, nil
}

// handleApprove sends an approval response
func (d *DiffViewer) handleApprove() (Overlay, tea.Cmd) {
	if d.responseFunc != nil {
		d.responseFunc(types.NewApprovalResponse(d.approvalID, types.ApprovalGranted))
	}
	return d, nil
}

// handleToggleSelection toggles between Accept and Reject
func (d *DiffViewer) handleToggleSelection() (Overlay, tea.Cmd) {
	if d.selected == ApprovalChoiceAccept {
		d.selected = ApprovalChoiceReject
	} else {
		d.selected = ApprovalChoiceAccept
	}
	return d, nil
}

// handleSubmit submits the currently selected choice
func (d *DiffViewer) handleSubmit() (Overlay, tea.Cmd) {
	if d.responseFunc != nil {
		decision := types.ApprovalRejected
		if d.selected == ApprovalChoiceAccept {
			decision = types.ApprovalGranted
		}
		d.responseFunc(types.NewApprovalResponse(d.approvalID, decision))
	}
	return d, nil
}

// handleSelectAccept selects the Accept option
func (d *DiffViewer) handleSelectAccept() (Overlay, tea.Cmd) {
	d.selected = ApprovalChoiceAccept
	return d, nil
}

// handleSelectReject selects the Reject option
func (d *DiffViewer) handleSelectReject() (Overlay, tea.Cmd) {
	d.selected = ApprovalChoiceReject
	return d, nil
}

// handleViewportScroll forwards scroll commands to the viewport
func (d *DiffViewer) handleViewportScroll(msg tea.KeyMsg) (Overlay, tea.Cmd) {
	var cmd tea.Cmd
	d.viewport, cmd = d.viewport.Update(msg)
	return d, cmd
}

// handleWindowResize updates dimensions when window is resized
func (d *DiffViewer) handleWindowResize(msg tea.WindowSizeMsg) (Overlay, tea.Cmd) {
	var cmd tea.Cmd
	d.width = msg.Width
	d.height = msg.Height
	d.viewport, cmd = d.viewport.Update(msg)
	return d, cmd
}

func (d *DiffViewer) View() string {
	var s strings.Builder

	// Content width accounts for outer container border (2) + padding (4) = 6 chars
	contentWidth := d.width - 6

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(salmonPink)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(mutedGray)

	title := "Tool Approval Required"
	subtitle := fmt.Sprintf("%s: %s", d.toolName, d.preview.Title)

	// Manually center by calculating padding
	titleLen := len(title)
	subtitleLen := len(subtitle)
	titlePadding := (contentWidth - titleLen) / 2
	subtitlePadding := (contentWidth - subtitleLen) / 2

	s.WriteString(strings.Repeat(" ", titlePadding) + titleStyle.Render(title))
	s.WriteString("\n")
	s.WriteString(strings.Repeat(" ", subtitlePadding) + subtitleStyle.Render(subtitle))
	s.WriteString("\n\n")

	// Diff box has its own border (2) + padding (2), so reduce width further
	diffStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(salmonPink).
		Padding(0, 1).
		Width(contentWidth - 4)

	s.WriteString(diffStyle.Render(d.viewport.View()))
	s.WriteString("\n\n")

	acceptStyle := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 2)

	rejectStyle := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 2)

	if d.selected == ApprovalChoiceAccept {
		acceptStyle = acceptStyle.
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#A8E6CF"))
		rejectStyle = rejectStyle.
			Foreground(mutedGray).
			Background(darkBg)
	} else {
		acceptStyle = acceptStyle.
			Foreground(mutedGray).
			Background(darkBg)
		rejectStyle = rejectStyle.
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FFB3BA"))
	}

	acceptBtn := acceptStyle.Render("✓ Accept (Enter / Ctrl+A)")
	rejectBtn := rejectStyle.Render("✗ Reject (Esc / Ctrl+R)")

	// Create spacer with dark background to match container
	spacerStyle := lipgloss.NewStyle().Background(darkBg)
	spacer := spacerStyle.Render("  ")

	// Join buttons with styled spacer
	buttonsRow := acceptBtn + spacer + rejectBtn

	// Manually center buttons
	buttonsLen := lipgloss.Width(buttonsRow)
	buttonsPadding := (contentWidth - buttonsLen) / 2

	s.WriteString(strings.Repeat(" ", buttonsPadding) + buttonsRow)
	s.WriteString("\n")

	hintStyle := lipgloss.NewStyle().
		Foreground(mutedGray).
		Italic(true)

	hints := "↑↓ to scroll • ← → Tab to choose • Enter to submit"
	hintsLen := len(hints)
	hintsPadding := (contentWidth - hintsLen) / 2

	s.WriteString(strings.Repeat(" ", hintsPadding) + hintStyle.Render(hints))

	containerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(salmonPink).
		Padding(1, 2).
		Width(d.width).
		Background(darkBg)

	return containerStyle.Render(s.String())
}

func (d *DiffViewer) Focused() bool {
	return d.focused
}

func (d *DiffViewer) Width() int {
	return d.width
}

func (d *DiffViewer) Height() int {
	return d.height
}
