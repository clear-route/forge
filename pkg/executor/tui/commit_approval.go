package tui

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/entrhq/forge/pkg/agent/slash"
)

// CommitApprovalRequest is a concrete implementation of ApprovalRequest for git commits.
// It encapsulates all data needed to preview and execute a commit operation.
type CommitApprovalRequest struct {
	files        []string
	message      string
	diff         string
	args         string
	slashHandler *slash.Handler
}

// NewCommitApprovalRequest creates a new commit approval request
func NewCommitApprovalRequest(files []string, message, diff, args string, slashHandler *slash.Handler) *CommitApprovalRequest {
	return &CommitApprovalRequest{
		files:        files,
		message:      message,
		diff:         diff,
		args:         args,
		slashHandler: slashHandler,
	}
}

// Title returns the approval dialog title
func (c *CommitApprovalRequest) Title() string {
	return "Commit Preview"
}

// Content returns the formatted content for the commit preview
func (c *CommitApprovalRequest) Content() string {
	var b strings.Builder

	// Show files to commit
	if len(c.files) > 0 {
		b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(salmonPink).Render("Files to commit:"))
		b.WriteString("\n")
		for _, file := range c.files {
			b.WriteString("  • " + file + "\n")
		}
		b.WriteString("\n")
	}

	// Show commit message
	if c.message != "" {
		b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(salmonPink).Render("Commit Message:"))
		b.WriteString("\n")
		b.WriteString(c.message)
		b.WriteString("\n\n")
	}

	// Show diff with syntax highlighting
	if c.diff != "" {
		b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(salmonPink).Render("Changes:"))
		b.WriteString("\n")

		highlightedDiff, err := HighlightDiff(c.diff, "")
		if err != nil {
			b.WriteString(c.diff)
		} else {
			b.WriteString(highlightedDiff)
		}
	}

	return b.String()
}

// OnApprove returns the command to execute when the user approves the commit
func (c *CommitApprovalRequest) OnApprove() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		result, err := c.slashHandler.Execute(ctx, &slash.Command{
			Name: "commit",
			Arg:  c.args,
		})

		if err != nil {
			return toastMsg{
				message: "Commit Failed",
				details: fmt.Sprintf("%v", err),
				icon:    "❌",
				isError: true,
			}
		}

		return toastMsg{
			message: "Success",
			details: result,
			icon:    "✅",
			isError: false,
		}
	}
}

// OnReject returns the command to execute when the user rejects the commit
func (c *CommitApprovalRequest) OnReject() tea.Cmd {
	return func() tea.Msg {
		return toastMsg{
			message: "Canceled",
			details: "/commit command canceled",
			icon:    "ℹ️",
			isError: false,
		}
	}
}
