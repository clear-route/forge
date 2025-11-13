package tui

import (
	"context"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/entrhq/forge/pkg/agent/slash"
)

// PRApprovalRequest is a concrete implementation of ApprovalRequest for pull requests.
// It encapsulates all data needed to preview and execute a PR operation.
type PRApprovalRequest struct {
	branch       string
	prTitle      string
	prDesc       string
	changes      string
	args         string
	slashHandler *slash.Handler
}

// NewPRApprovalRequest creates a new PR approval request
func NewPRApprovalRequest(branch, prTitle, prDesc, changes, args string, slashHandler *slash.Handler) *PRApprovalRequest {
	return &PRApprovalRequest{
		branch:       branch,
		prTitle:      prTitle,
		prDesc:       prDesc,
		changes:      changes,
		args:         args,
		slashHandler: slashHandler,
	}
}

// Title returns the approval dialog title
func (p *PRApprovalRequest) Title() string {
	// Use the generated PR title if available
	if p.prTitle != "" && p.prTitle != p.args {
		return p.prTitle
	}
	// Fall back to user-provided title if any
	if p.args != "" {
		return p.args
	}
	return "Pull Request Preview"
}

// Content returns the formatted content for the PR preview
func (p *PRApprovalRequest) Content() string {
	var b strings.Builder

	// Show branch info at the top
	if p.branch != "" {
		b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(salmonPink).Render("Branch:"))
		b.WriteString("\n")
		b.WriteString("  " + p.branch + "\n")
		b.WriteString("\n")
	}

	// Show PR title if it's different from the overlay title
	// (i.e., if we're showing user-provided title in header, show generated title here)
	if p.prTitle != "" && p.args != "" && p.prTitle != p.args {
		b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(salmonPink).Render("Generated Title:"))
		b.WriteString("\n")
		b.WriteString(p.prTitle)
		b.WriteString("\n\n")
	}

	// Show PR description prominently if available
	if p.prDesc != "" {
		b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(salmonPink).Render("Description:"))
		b.WriteString("\n")
		b.WriteString(p.prDesc)
		b.WriteString("\n\n")
	}

	// Show commits and changes
	if p.changes != "" {
		b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(salmonPink).Render("Commits & Changes:"))
		b.WriteString("\n")
		b.WriteString(p.changes)
	}

	return b.String()
}

// OnApprove returns the command to execute when the user approves the PR
func (p *PRApprovalRequest) OnApprove() tea.Cmd {
	return tea.Batch(
		// First, signal that we're starting PR creation
		func() tea.Msg {
			return operationStartMsg{
				message: "Creating pull request on GitHub...",
			}
		},
		// Then execute the PR creation
		func() tea.Msg {
			ctx := context.Background()
			result, err := p.slashHandler.Execute(ctx, &slash.Command{
				Name: "pr",
				Arg:  p.args,
			})
			return operationCompleteMsg{
				result:       result,
				err:          err,
				successTitle: "Success",
				successIcon:  "🔀",
				errorTitle:   "PR Failed",
				errorIcon:    "❌",
			}
		},
	)
}

// OnReject returns the command to execute when the user rejects the PR
func (p *PRApprovalRequest) OnReject() tea.Cmd {
	return func() tea.Msg {
		return toastMsg{
			message: "Canceled",
			details: "/pr command canceled",
			icon:    "ℹ️",
			isError: false,
		}
	}
}
