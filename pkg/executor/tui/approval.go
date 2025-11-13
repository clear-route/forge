package tui

import tea "github.com/charmbracelet/bubbletea"

// ApprovalRequest represents a command that requires user approval before execution.
// This interface allows any command to request approval in a generic, decoupled way.
type ApprovalRequest interface {
	// Title returns the title to display in the approval overlay
	Title() string

	// Content returns the formatted content to display for review
	Content() string

	// OnApprove returns the command to execute when the user approves
	OnApprove() tea.Cmd

	// OnReject returns the command to execute when the user rejects
	OnReject() tea.Cmd
}

// approvalRequestMsg is sent when a command returns an approval request
type approvalRequestMsg struct {
	request ApprovalRequest
}
