package tui

import (
	"fmt"
	"log"
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
	overlayWidth := int(float64(width) * 0.9)
	if overlayWidth < 80 {
		overlayWidth = 80
	}

	// Fixed viewport height: max 10 lines for diff content
	const maxViewportHeight = 10
	viewportHeight := maxViewportHeight

	// Calculate total overlay height
	// Title (2) + subtitle (1) + spacing (1) + border (2) + buttons (2) + hints (1) = 9 lines
	// Plus viewport height
	overlayHeight := viewportHeight + 9

	vp := viewport.New(overlayWidth-4, viewportHeight)
	vp.Style = lipgloss.NewStyle()

	if preview != nil {
		vp.SetContent(preview.Content)
	}

	// DEBUG: Log overlay dimensions
	log.Printf("DEBUG DiffViewer: screen=%dx%d overlay=%dx%d viewport=%dx%d",
		width, height, overlayWidth, overlayHeight, overlayWidth-4, viewportHeight)

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
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if d.responseFunc != nil {
				d.responseFunc(types.NewApprovalResponse(d.approvalID, types.ApprovalRejected))
			}
			return d, nil

		case "ctrl+a":
			if d.responseFunc != nil {
				d.responseFunc(types.NewApprovalResponse(d.approvalID, types.ApprovalGranted))
			}
			return d, nil

		case "ctrl+r":
			if d.responseFunc != nil {
				d.responseFunc(types.NewApprovalResponse(d.approvalID, types.ApprovalRejected))
			}
			return d, nil

		case "tab":
			// Tab switches between Accept/Reject
			if d.selected == ApprovalChoiceAccept {
				d.selected = ApprovalChoiceReject
			} else {
				d.selected = ApprovalChoiceAccept
			}
			return d, nil

		case "enter":
			// Enter submits the selected choice
			log.Printf("DEBUG DiffViewer: Enter key pressed, selected=%d", d.selected)
			if d.responseFunc != nil {
				var decision types.ApprovalDecision
				if d.selected == ApprovalChoiceAccept {
					decision = types.ApprovalGranted
					log.Printf("DEBUG DiffViewer: Calling responseFunc with GRANTED for approval %s", d.approvalID)
				} else {
					decision = types.ApprovalRejected
					log.Printf("DEBUG DiffViewer: Calling responseFunc with REJECTED for approval %s", d.approvalID)
				}
				d.responseFunc(types.NewApprovalResponse(d.approvalID, decision))
				log.Printf("DEBUG DiffViewer: responseFunc completed, returning")
			} else {
				log.Printf("DEBUG DiffViewer: ERROR - responseFunc is nil!")
			}
			return d, nil

		case "left", "h":
			// Left arrow or 'h' selects Accept
			d.selected = ApprovalChoiceAccept
			return d, nil

		case "right", "l":
			// Right arrow or 'l' selects Reject
			d.selected = ApprovalChoiceReject
			return d, nil

		default:
			// All other keys (including up/down, pgup/pgdn, j/k for vim) go to viewport for scrolling
			d.viewport, cmd = d.viewport.Update(msg)
			return d, cmd
		}

	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height
		d.viewport, cmd = d.viewport.Update(msg)
		return d, cmd
	}

	return d, nil
}

func (d *DiffViewer) View() string {
	var s strings.Builder

	// Content width accounts for outer container border (2) + padding (4) = 6 chars
	contentWidth := d.width - 6

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(salmonPink).
		Align(lipgloss.Center).
		Width(contentWidth)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(mutedGray).
		Align(lipgloss.Center).
		Width(contentWidth)

	title := "Tool Approval Required"
	subtitle := fmt.Sprintf("%s: %s", d.toolName, d.preview.Title)

	s.WriteString(titleStyle.Render(title))
	s.WriteString("\n")
	s.WriteString(subtitleStyle.Render(subtitle))
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
		Padding(0, 2).
		MarginRight(2)

	rejectStyle := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 2)

	if d.selected == ApprovalChoiceAccept {
		acceptStyle = acceptStyle.
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#A8E6CF"))
		rejectStyle = rejectStyle.
			Foreground(mutedGray)
	} else {
		acceptStyle = acceptStyle.
			Foreground(mutedGray)
		rejectStyle = rejectStyle.
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FFB3BA"))
	}

	acceptBtn := acceptStyle.Render("✓ Accept (Enter / Ctrl+A)")
	rejectBtn := rejectStyle.Render("✗ Reject (Esc / Ctrl+R)")

	buttonsRow := lipgloss.JoinHorizontal(lipgloss.Left, acceptBtn, rejectBtn)
	buttonsStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(contentWidth)

	s.WriteString(buttonsStyle.Render(buttonsRow))
	s.WriteString("\n")

	hintStyle := lipgloss.NewStyle().
		Foreground(mutedGray).
		Italic(true).
		Align(lipgloss.Center).
		Width(contentWidth)

	s.WriteString(hintStyle.Render("↑↓ to scroll • ← → Tab to choose • Enter to submit"))

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
