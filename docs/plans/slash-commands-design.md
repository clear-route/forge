# Slash Commands Design

## Overview

Slash commands provide a quick way for users to execute common actions and control the TUI/agent behavior through a command-line-style interface.

## User Experience

### Triggering Commands
- When the user types `/` at the beginning of the input (or after whitespace), a command palette appears above the input bar
- The palette shows available commands with descriptions
- As the user continues typing, the list filters to matching commands
- Commands can be selected with arrow keys + Enter, or by typing the full command name
- Pressing Escape cancels command entry and returns to normal input

### Command Palette UI
```
┌─────────────────────────────────────────────────────────┐
│ Available Commands:                                      │
│ > /help     Show tips and keyboard shortcuts             │
│   /stop     Stop current agent operation                 │
│   /commit   Create git commit from session changes       │
│   /pr       Create pull request from current branch      │
└─────────────────────────────────────────────────────────┘
```

## Command Categories

### 1. TUI Commands (Executor-Only)
These commands are handled entirely by the TUI executor and do not send events to the agent.

#### `/help`
- **Description**: Show help information and keyboard shortcuts
- **Behavior**: Displays an overlay with tips, shortcuts, and available commands
- **Use Case**: Quick reference for users

### 2. Agent Commands (Sent to Agent)
These commands send events/inputs to the agent, potentially affecting the agent loop.

#### `/stop`
- **Description**: Stop the current agent operation
- **Behavior**: Sends `InputTypeCancel` to agent, interrupts current processing
- **Use Case**: Stop a long-running operation or incorrect action

#### `/commit [message]`
- **Description**: Create a git commit from session changes
- **Arguments**: Optional freeform text after `/commit`
  - If provided: Use as commit message
  - If omitted: Auto-generate using LLM based on file diffs
- **Examples**:
  - `/commit` - Auto-generate message
  - `/commit chore: update dependencies` - Use provided message
  - `/commit fix navbar rendering issue` - Use provided message
- **Behavior**: 
  - Collects files modified during current session
  - Generates conventional commit message if not provided
  - Stages and commits changes
  - Shows toast with commit hash
- **Use Case**: Quickly commit agent-generated changes with meaningful messages

#### `/pr [title]`
- **Description**: Create a pull request from current branch
- **Arguments**: Optional freeform text after `/pr`
  - If provided: Use as PR title (still generate description)
  - If omitted: Auto-generate both title and description using LLM
- **Examples**:
  - `/pr` - Fully auto-generated PR
  - `/pr Add user authentication system` - Custom title, generated description
- **Behavior**:
  - Verifies all changes are committed
  - Detects base branch (where current branch diverged from)
  - Collects all commits between base and current branch
  - Analyzes actual code changes (diffs) for material changes
  - Generates comprehensive PR title and description using LLM
  - Pushes branch to remote
  - Creates PR via provider API
  - Shows toast with PR URL
- **Use Case**: Streamline workflow from changes to pull request
- **Requirements**: GITHUB_TOKEN or GITLAB_TOKEN environment variable

## Technical Architecture

### Component Structure

```
pkg/executor/tui/
├── executor.go              # Main TUI executor (existing)
├── slash_commands.go        # Command registry and handlers
├── command_palette.go       # Command palette overlay UI
└── styles.go                # Styling (existing, extend as needed)

pkg/types/
├── input.go                 # Add InputTypeSlashCommand
└── event.go                 # Potentially add EventTypeCommandExecuted
```

### Data Structures

```go
// CommandType indicates whether a command is handled by TUI or Agent
type CommandType int

const (
    CommandTypeTUI CommandType = iota   // Handled entirely by TUI
    CommandTypeAgent                     // Sent to agent
)

// CommandHandler processes a slash command
type CommandHandler func(m *model, args []string) (tea.Model, tea.Cmd)

// SlashCommand represents a registered command
type SlashCommand struct {
    Name        string          // Command name (without /)
    Description string          // Short description for palette
    Type        CommandType     // Where to handle the command
    Handler     CommandHandler  // Handler function (for TUI commands)
    MinArgs     int            // Minimum number of arguments
    MaxArgs     int            // Maximum number of arguments (-1 for unlimited)
}

// CommandPalette manages command suggestions and selection
type CommandPalette struct {
    commands       []SlashCommand
    filteredCommands []SlashCommand
    selectedIndex  int
    filter         string
    active         bool
}
```

### Execution Flow

```
┌─────────────┐
│ User types  │
│     /       │
└──────┬──────┘
       │
       v
┌─────────────────────┐
│ TUI detects slash   │
│ Activates palette   │
└──────┬──────────────┘
       │
       v
┌─────────────────────────────┐
│ User types/selects command  │
│ Presses Enter               │
└──────┬──────────────────────┘
       │
       v
┌─────────────────┐
│ Parse command   │
│ & arguments     │
└──────┬──────────┘
       │
       ├─────────────────┐
       │                 │
       v                 v
┌──────────────┐  ┌──────────────┐
│ TUI Command  │  │Agent Command │
└──────┬───────┘  └──────┬───────┘
       │                 │
       v                 v
┌──────────────┐  ┌──────────────────┐
│Execute local │  │Send to agent via │
│handler       │  │channels.Input    │
└──────────────┘  └──────────────────┘
```

### Integration Points

1. **Input Detection** (in `Update()`):
   - Detect when user types `/` at start of input
   - Activate command palette
   - Filter commands as user types

2. **Command Parsing** (in `handleUserInput()`):
   - Check if input starts with `/`
   - Parse command name and arguments
   - Route to appropriate handler

3. **TUI Command Execution** (in `slash_commands.go`):
   - Execute handler function directly on model
   - Update UI state
   - Return updated model and commands

4. **Agent Command Execution**:
   - Create appropriate Input type
   - Send to `m.channels.Input`
   - Let agent process normally

## Implementation Phases

### Phase 1: Core Infrastructure
- [ ] Define command types and structures
- [ ] Create command registry in slash_commands.go
- [ ] Add command parsing logic
- [ ] Implement /help TUI command handler

### Phase 2: Command Palette UI
- [ ] Create CommandPalette component
- [ ] Add command filtering logic
- [ ] Implement keyboard navigation (arrow keys, enter, escape)
- [ ] Style palette with Lipgloss

### Phase 3: Git Integration
- [ ] Create file modification tracker
- [ ] Implement /commit command handler
- [ ] Add commit message generation with LLM
- [ ] Implement /stop command (interrupt agent)
- [ ] Add git provider detection (GitHub/GitLab)
- [ ] Implement GitHub API client
- [ ] Implement /pr command handler
- [ ] Add PR description generation with LLM

### Phase 4: Polish & Testing
- [ ] Add comprehensive error handling for git operations
- [ ] Implement token management system
- [ ] Add integration tests for git commands
- [ ] Add command history/autocomplete
- [ ] Add fuzzy matching for command filtering

## Open Questions

1. **Input Type**: Should we add a new `InputTypeSlashCommand` or use metadata on existing `InputTypeUserInput`?
   - **Recommendation**: Use metadata on `InputTypeUserInput` for simplicity. Add `"command": true` and `"command_name": "reset"` metadata.

2. **Command Arguments**: How should we handle commands with complex arguments?
   - **Recommendation**: Start simple with space-separated args, use quotes for multi-word args later if needed.

3. **Command Discovery**: Should commands be dynamically registered or hardcoded?
   - **Recommendation**: Static registration in init() for built-in commands, with option to extend later.

4. **Error Handling**: How should invalid commands be presented?
   - **Recommendation**: Show error toast notification, keep input in textarea for correction.

5. **Git Token Storage**: Where should we store GitHub/GitLab tokens?
   - **Recommendation**: Environment variables first, then config file. Never store in conversation history.

## Testing Strategy

### Unit Tests
- Command parsing logic
- Command filtering/matching
- Individual command handlers

### Integration Tests
- Full command execution flow
- TUI command effects on model state
- Agent command routing

### Manual Testing
- UI/UX of command palette
- Keyboard navigation
- Edge cases (empty args, invalid commands)
- Command chaining/sequences

## Future Enhancements

1. **Additional TUI Commands**: 
   - `/clear` - Clear conversation history from viewport
   - `/copy` - Copy content to clipboard
   - `/export` - Export conversation to file
2. **Custom Commands**: Allow users to define custom slash commands via config file
3. **Command Aliases**: Short aliases for common commands (e.g., `/c` for `/commit`)
4. **Command History**: Arrow up/down to cycle through previous commands
5. **Interactive Commit Staging**: Allow user to select specific files to commit
6. **PR Templates**: Support .github/pull_request_template.md
7. **Draft PRs**: Option to create draft pull requests
8. **GitLab Support**: Full GitLab API integration alongside GitHub

## Security Considerations

- **Token Storage**: Never log or expose GitHub/GitLab tokens in output
- **Git Operations**: Validate all git commands to prevent injection
- **API Credentials**: Store tokens securely, use environment variables
- **Remote URLs**: Validate git remote URLs before making API calls
- **Command Injection**: Sanitize all command arguments

## Accessibility

- Command palette should be keyboard-navigable
- Provide visual feedback for command execution
- Support screen readers (announce command selection)
- Allow customization of command palette position/size
