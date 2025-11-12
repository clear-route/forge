package tui

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/entrhq/forge/pkg/agent/git"
	"github.com/entrhq/forge/pkg/types"
)

// CommandType indicates whether a command is handled by TUI or Agent
type CommandType int

const (
	CommandTypeTUI   CommandType = iota // Handled entirely by TUI
	CommandTypeAgent                    // Sent to agent
)

// CommandHandler processes a slash command and returns a tea.Cmd.
// The model is passed as a pointer and can be modified directly.
type CommandHandler func(m *model, args []string) tea.Cmd

// SlashCommand represents a registered command
type SlashCommand struct {
	Name        string         // Command name (without /)
	Description string         // Short description for palette
	Type        CommandType    // Where to handle the command
	Handler     CommandHandler // Handler function (for TUI commands)
	MinArgs     int            // Minimum number of arguments
	MaxArgs     int            // Maximum number of arguments (-1 for unlimited)
}

// commandRegistry holds all registered slash commands
var commandRegistry map[string]*SlashCommand

// init initializes the command registry with built-in commands
func init() {
	commandRegistry = make(map[string]*SlashCommand)

	// Register built-in commands
	registerCommand(&SlashCommand{
		Name:        "help",
		Description: "Show tips and keyboard shortcuts",
		Type:        CommandTypeTUI,
		Handler:     handleHelpCommand,
		MinArgs:     0,
		MaxArgs:     0,
	})

	registerCommand(&SlashCommand{
		Name:        "stop",
		Description: "Stop current agent operation",
		Type:        CommandTypeAgent,
		Handler:     handleStopCommand,
		MinArgs:     0,
		MaxArgs:     0,
	})

	registerCommand(&SlashCommand{
		Name:        "commit",
		Description: "Create git commit from session changes",
		Type:        CommandTypeAgent,
		Handler:     handleCommitCommand,
		MinArgs:     0,
		MaxArgs:     -1, // Unlimited for commit message
	})

	registerCommand(&SlashCommand{
		Name:        "pr",
		Description: "Create pull request from current branch",
		Type:        CommandTypeAgent,
		Handler:     handlePRCommand,
		MinArgs:     0,
		MaxArgs:     -1, // Unlimited for PR title
	})
}

// registerCommand adds a command to the registry
func registerCommand(cmd *SlashCommand) {
	commandRegistry[cmd.Name] = cmd
}

// getCommand retrieves a command from the registry
func getCommand(name string) (*SlashCommand, bool) {
	cmd, exists := commandRegistry[name]
	return cmd, exists
}

// getAllCommands returns all registered commands
func getAllCommands() []*SlashCommand {
	commands := make([]*SlashCommand, 0, len(commandRegistry))
	for _, cmd := range commandRegistry {
		commands = append(commands, cmd)
	}
	return commands
}

// parseSlashCommand parses a slash command input into command name and arguments
// Returns: commandName, args, isCommand
func parseSlashCommand(input string) (string, []string, bool) {
	trimmed := strings.TrimSpace(input)
	if !strings.HasPrefix(trimmed, "/") {
		return "", nil, false
	}

	// Remove the leading /
	trimmed = trimmed[1:]

	// Split into parts
	parts := strings.Fields(trimmed)
	if len(parts) == 0 {
		return "", nil, false
	}

	commandName := parts[0]
	args := []string{}
	if len(parts) > 1 {
		args = parts[1:]
	}

	return commandName, args, true
}

// executeSlashCommand executes a slash command
func executeSlashCommand(m model, commandName string, args []string) (model, tea.Cmd) {
	cmd, exists := getCommand(commandName)
	if !exists {
		// Unknown command - show error toast
		m.showToast("Unknown command", fmt.Sprintf("Command '/%s' not found. Type /help for available commands.", commandName), "❌", true)
		return m, nil
	}

	// Validate argument count
	if len(args) < cmd.MinArgs {
		m.showToast("Invalid arguments", fmt.Sprintf("Command '/%s' requires at least %d argument(s)", commandName, cmd.MinArgs), "❌", true)
		return m, nil
	}
	if cmd.MaxArgs != -1 && len(args) > cmd.MaxArgs {
		m.showToast("Invalid arguments", fmt.Sprintf("Command '/%s' accepts at most %d argument(s)", commandName, cmd.MaxArgs), "❌", true)
		return m, nil
	}

	// Execute the command handler
	if cmd.Handler != nil {
		cmd := cmd.Handler(&m, args)
		return m, cmd
	}

	return m, nil
}

// handleHelpCommand shows help information
func handleHelpCommand(m *model, args []string) tea.Cmd {
	// Build help content
	var helpContent strings.Builder
	helpContent.WriteString("Available Commands:\n\n")

	commands := getAllCommands()
	for _, cmd := range commands {
		helpContent.WriteString(fmt.Sprintf("  /%s\n", cmd.Name))
		helpContent.WriteString(fmt.Sprintf("    %s\n\n", cmd.Description))
	}

	helpContent.WriteString("Keyboard Shortcuts:\n\n")
	helpContent.WriteString("  Enter        Send message\n")
	helpContent.WriteString("  Alt+Enter    New line\n")
	helpContent.WriteString("  Ctrl+C       Exit\n")
	helpContent.WriteString("  Ctrl+D       Show command help\n\n")

	helpContent.WriteString("Tips:\n\n")
	helpContent.WriteString("  • Type / to see available commands\n")
	helpContent.WriteString("  • Use arrow keys to navigate command palette\n")
	helpContent.WriteString("  • Press Escape to cancel command entry\n")

	// Create and activate the help overlay
	helpOverlay := NewHelpOverlay("Help", helpContent.String())
	m.overlay.activate(OverlayModeHelp, helpOverlay)

	return nil
}

// handleStopCommand stops the current agent operation
func handleStopCommand(m *model, args []string) tea.Cmd {
	if m.channels != nil {
		// Send cancel input to agent
		m.channels.Input <- types.NewCancelInput()
		m.showToast("Stopping", "Sent stop signal to agent", "⏹️", false)
	}
	return nil
}

// handleCommitCommand creates a git commit using the slash handler
func handleCommitCommand(m *model, args []string) tea.Cmd {
	if m.slashHandler == nil {
		m.showToast("Error", "Git operations not available", "❌", true)
		return nil
	}

	commitMessage := strings.Join(args, " ")
	
	return func() tea.Msg {
		// Gather preview data in background
		ctx := context.Background()
		
		// Get modified files
		files, err := git.GetModifiedFiles(m.workspaceDir)
		if err != nil {
			return toastMsg{
				message: "Commit Failed",
				details: fmt.Sprintf("Failed to get modified files: %v", err),
				icon:    "❌",
				isError: true,
			}
		}
		
		if len(files) == 0 {
			return toastMsg{
				message: "Nothing to Commit",
				details: "No modified files found",
				icon:    "ℹ️",
				isError: false,
			}
		}
		
		// Get diff for preview
		diff, err := getDiffForFiles(m.workspaceDir, files)
		if err != nil {
			diff = "(Unable to generate diff preview)"
		}
		
		// Generate commit message if not provided
		message := commitMessage
		if message == "" {
			generatedMsg, err := m.commitGen.Generate(ctx, m.workspaceDir, files)
			if err != nil {
				return toastMsg{
					message: "Commit Failed",
					details: fmt.Sprintf("Failed to generate commit message: %v", err),
					icon:    "❌",
					isError: true,
				}
			}
			message = generatedMsg
		}
		
		// Return preview message
		return slashCommandPreviewMsg{
			commandName: "commit",
			title:       "Commit Preview",
			args:        commitMessage, // Store original args
			files:       files,
			message:     message,
			diff:        diff,
			onApprove: func() {
				// This will be called when user approves
				// The actual execution will happen in the Update handler
			},
			onReject: func() {
				// This will be called when user rejects
			},
		}
	}
}

// getDiffForFiles gets the git diff for the specified files
func getDiffForFiles(workingDir string, files []string) (string, error) {
	// Try to get diff against HEAD first (for modified tracked files)
	args := append([]string{"diff", "HEAD", "--"}, files...)
	cmd := exec.Command("git", args...)
	cmd.Dir = workingDir
	
	output, err := cmd.Output()
	if err != nil {
		// If that fails (new files not in HEAD), try without HEAD
		// This will show working directory changes
		args = append([]string{"diff", "--"}, files...)
		cmd = exec.Command("git", args...)
		cmd.Dir = workingDir
		
		output, err = cmd.Output()
		if err != nil {
			// If both fail, return a helpful message
			return "(Unable to generate diff preview - files may be untracked)", nil
		}
	}
	
	// If output is empty, the files might be new/untracked
	// Try to show them as additions
	if len(output) == 0 {
		// Get diff of what would be staged if we add these files
		args = append([]string{"diff", "--no-index", "/dev/null", "--"}, files...)
		cmd = exec.Command("git", args...)
		cmd.Dir = workingDir
		
		output, err = cmd.Output()
		if err != nil {
			return "(New/untracked files - run commit to see full content)", nil
		}
	}
	
	return string(output), nil
}

// handlePRCommand creates a pull request using the slash handler
func handlePRCommand(m *model, args []string) tea.Cmd {
	if m.slashHandler == nil {
		m.showToast("Error", "Git operations not available", "❌", true)
		return nil
	}

	prTitle := strings.Join(args, " ")
	
	return func() tea.Msg {
		// Gather preview data in background
		ctx := context.Background()
		
		// Get base branch
		base, err := git.DetectBaseBranch(m.workspaceDir)
		if err != nil {
			return toastMsg{
				message: "PR Failed",
				details: fmt.Sprintf("Failed to detect base branch: %v", err),
				icon:    "❌",
				isError: true,
			}
		}
		
		// Get current branch
		head, err := getCurrentBranch(m.workspaceDir)
		if err != nil {
			return toastMsg{
				message: "PR Failed",
				details: fmt.Sprintf("Failed to get current branch: %v", err),
				icon:    "❌",
				isError: true,
			}
		}
		
		// Get commits since base
		commits, err := git.GetCommitsSinceBase(m.workspaceDir, base, head)
		if err != nil {
			return toastMsg{
				message: "PR Failed",
				details: fmt.Sprintf("Failed to get commits: %v", err),
				icon:    "❌",
				isError: true,
			}
		}
		
		if len(commits) == 0 {
			return toastMsg{
				message: "Nothing to PR",
				details: "No commits found for pull request",
				icon:    "ℹ️",
				isError: false,
			}
		}
		
		// Get diff summary
		diffSummary, err := git.GetDiffSummary(m.workspaceDir, base, head)
		if err != nil {
			return toastMsg{
				message: "PR Failed",
				details: fmt.Sprintf("Failed to get diff summary: %v", err),
				icon:    "❌",
				isError: true,
			}
		}
		
		// Generate PR content
		prContent, err := m.prGen.Generate(ctx, commits, diffSummary, base, head, prTitle)
		if err != nil {
			return toastMsg{
				message: "PR Failed",
				details: fmt.Sprintf("Failed to generate PR content: %v", err),
				icon:    "❌",
				isError: true,
			}
		}
		
		// Build commits and changes preview
		var changesContent strings.Builder
		changesContent.WriteString(fmt.Sprintf("Commits (%d):\n", len(commits)))
		for _, commit := range commits {
			changesContent.WriteString(fmt.Sprintf("  • %s\n", commit.Message))
		}
		changesContent.WriteString("\n")
		changesContent.WriteString(diffSummary)
		
		// Return preview message with separate PR title and description
		return slashCommandPreviewMsg{
			commandName: "pr",
			title:       "Pull Request Preview",
			args:        prTitle, // Store original args
			files:       []string{fmt.Sprintf("%s → %s", head, base)},
			prTitle:     prContent.Title,
			prDesc:      prContent.Description,
			diff:        changesContent.String(),
			onApprove: func() {
				// This will be called when user approves
				// The actual execution will happen in the Update handler
			},
			onReject: func() {
				// This will be called when user rejects
			},
		}
	}
}

// getCurrentBranch gets the current git branch name
func getCurrentBranch(workingDir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = workingDir
	
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	
	return strings.TrimSpace(string(output)), nil
}

// toastMsg is a message type for showing toast notifications
type toastMsg struct {
	message string
	details string
	icon    string
	isError bool
}

// slashCommandPreviewMsg requests showing a preview overlay for a slash command
type slashCommandPreviewMsg struct {
	commandName string
	title       string
	args        string   // Original command arguments
	files       []string
	message     string
	diff        string
	prTitle     string   // PR title (only for PR commands)
	prDesc      string   // PR description (only for PR commands)
	onApprove   func()
	onReject    func()
}

// slashCommandApprovedMsg indicates the user approved the slash command
type slashCommandApprovedMsg struct {
	commandName string
	args        []string
}

// slashCommandRejectedMsg indicates the user rejected the slash command
type slashCommandRejectedMsg struct {
	commandName string
}
