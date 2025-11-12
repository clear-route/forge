package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/entrhq/forge/pkg/types"
)

// CommandType indicates whether a command is handled by TUI or Agent
type CommandType int

const (
	CommandTypeTUI   CommandType = iota // Handled entirely by TUI
	CommandTypeAgent                     // Sent to agent
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

// handleCommitCommand creates a git commit
func handleCommitCommand(m *model, args []string) tea.Cmd {
	// For now, send as user input with metadata
	// Later this will be enhanced with file tracking and LLM generation
	commitMessage := strings.Join(args, " ")

	if m.channels != nil {
		if commitMessage == "" {
			// Auto-generate commit message
			input := types.NewUserInput("Create a git commit with an auto-generated conventional commit message based on the changes made in this session.")
			input.WithMetadata("command", "commit")
			input.WithMetadata("auto_generate", true)
			m.channels.Input <- input
		} else {
			// Use provided message
			input := types.NewUserInput(fmt.Sprintf("Create a git commit with this message: %s", commitMessage))
			input.WithMetadata("command", "commit")
			input.WithMetadata("message", commitMessage)
			m.channels.Input <- input
		}
		m.showToast("Commit", "Creating git commit...", "📝", false)
	}

	return nil
}

// handlePRCommand creates a pull request
func handlePRCommand(m *model, args []string) tea.Cmd {
	// For now, send as user input with metadata
	// Later this will be enhanced with base branch detection and LLM generation
	prTitle := strings.Join(args, " ")

	if m.channels != nil {
		if prTitle == "" {
			// Auto-generate PR
			input := types.NewUserInput("Create a pull request with auto-generated title and description based on the commits and changes in the current branch.")
			input.WithMetadata("command", "pr")
			input.WithMetadata("auto_generate", true)
			m.channels.Input <- input
		} else {
			// Use provided title
			input := types.NewUserInput(fmt.Sprintf("Create a pull request with this title: %s. Generate the description based on commits and changes.", prTitle))
			input.WithMetadata("command", "pr")
			input.WithMetadata("title", prTitle)
			m.channels.Input <- input
		}
		m.showToast("Pull Request", "Creating pull request...", "🔀", false)
	}

	return nil
}
