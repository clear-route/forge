# Product Requirements: Slash Commands

**Feature:** Slash Command System  
**Version:** 1.0  
**Status:** Implemented  
**Owner:** Core Team  
**Last Updated:** December 2024

---

## Overview

Slash Commands provide a quick, intuitive way for users to access TUI features and execute common operations without leaving the chat interface. Similar to Discord, Slack, and other modern chat applications, users type commands starting with "/" to trigger special functionality.

---

## Problem Statement

Terminal UI users face several challenges when trying to access features:

1. **Discoverability:** Users don't know what features are available
2. **Context Switching:** Opening menus or settings requires interrupting workflow
3. **Efficiency:** Common operations require too many steps
4. **Memory Load:** Users must remember keyboard shortcuts or menu locations
5. **Consistency:** No standard way to access special features across different tools

Traditional approaches like nested menus or complex keyboard shortcuts create friction and slow down power users.

---

## Goals

### Primary Goals

1. **Fast Access:** Provide instant access to common features with minimal typing
2. **Discoverability:** Make features easy to find through autocomplete and help
3. **Consistency:** Use familiar slash command pattern from other tools
4. **Efficiency:** Reduce clicks and keystrokes for frequent operations
5. **Extensibility:** Support adding new commands easily

### Non-Goals

1. **Full Scripting Language:** This is NOT a programming language
2. **Complex Syntax:** No parameters with flags, quoting, or escaping
3. **Command Chaining:** No piping or combining commands (initially)
4. **Remote Execution:** Commands execute locally only, not on remote servers

---

## User Personas

### Primary: Power User Developer
- **Background:** Experienced with Vim, Emacs, or other keyboard-driven tools
- **Workflow:** Minimizes mouse usage, loves keyboard shortcuts
- **Pain Points:** Slow menu navigation interrupts flow
- **Goals:** Execute commands as fast as thinking them

### Secondary: Modern Tool User
- **Background:** Uses Slack, Discord, VS Code command palette
- **Workflow:** Comfortable with slash commands and fuzzy search
- **Pain Points:** Wants familiar interaction patterns
- **Goals:** Quickly find and execute features

### Tertiary: New User
- **Background:** First time using Forge
- **Workflow:** Exploring features through trial and error
- **Pain Points:** Doesn't know what's available
- **Goals:** Discover features without reading full documentation

---

## Requirements

### Functional Requirements

#### FR1: Command Detection
- **R1.1:** Detect "/" at start of input as command trigger
- **R1.2:** Show autocomplete suggestions as user types
- **R1.3:** Support partial command matching (e.g., "/set" matches "/settings")
- **R1.4:** Distinguish commands from regular messages starting with "/"
- **R1.5:** Handle typos gracefully with suggestions

#### FR2: Command Autocomplete
- **R2.1:** Display command palette with available commands
- **R2.2:** Filter commands based on typed text
- **R2.3:** Show command descriptions in palette
- **R2.4:** Navigate suggestions with arrow keys
- **R2.5:** Select with Enter or Tab
- **R2.6:** Close palette with Esc

#### FR3: Core Commands
- **R3.1:** `/help` - Show help overlay with tips and shortcuts
- **R3.2:** `/stop` - Cancel current agent operation
- **R3.3:** `/commit` - Create git commit from session changes
- **R3.4:** `/pr` - Create pull request (with approval)
- **R3.5:** `/settings` - Open settings overlay
- **R3.6:** `/context` - Display detailed context information
- **R3.7:** `/bash` - Enter bash mode for shell commands

#### FR4: Command Execution
- **R4.1:** Execute command immediately on Enter
- **R4.2:** Provide visual feedback during execution
- **R4.3:** Show results or open relevant overlay
- **R4.4:** Handle errors gracefully with clear messages
- **R4.5:** Return to normal chat mode after completion

#### FR5: Command Help
- **R5.1:** Show full command list in help overlay
- **R5.2:** Include descriptions for each command
- **R5.3:** Provide usage examples
- **R5.4:** Document keyboard shortcuts
- **R5.5:** Make help searchable

#### FR6: Command Feedback
- **R6.1:** Visual indicator when command is recognized
- **R6.2:** Error message for unknown commands
- **R6.3:** Success confirmation for completed commands
- **R6.4:** Progress indicator for long-running commands
- **R6.5:** Toast notifications for background operations

#### FR7: Bash Mode (Special Command)
- **R7.1:** Enter dedicated bash command mode with `/bash`
- **R7.2:** Execute shell commands directly (with approval)
- **R7.3:** Show command output in real-time
- **R7.4:** Allow multiple commands in sequence
- **R7.5:** Exit bash mode with `/exit` or Ctrl+D

### Non-Functional Requirements

#### NFR1: Performance
- **N1.1:** Command detection within 10ms of "/" keystroke
- **N1.2:** Autocomplete palette opens within 50ms
- **N1.3:** Command execution starts within 100ms
- **N1.4:** No lag when typing command text
- **N1.5:** Smooth palette navigation

#### NFR2: Usability
- **N2.1:** Intuitive for users familiar with Slack/Discord
- **N2.2:** Self-documenting through autocomplete
- **N2.3:** Consistent behavior across all commands
- **N2.4:** Clear visual distinction from regular chat
- **N2.5:** Keyboard-accessible (no mouse required)

#### NFR3: Reliability
- **N3.1:** Commands never crash the TUI
- **N3.2:** Graceful handling of invalid input
- **N3.3:** Consistent state after command execution
- **N3.4:** Proper cleanup if command interrupted
- **N3.5:** Recovery from command errors

#### NFR4: Extensibility
- **N4.1:** Easy to add new commands
- **N4.2:** Commands can be enabled/disabled
- **N4.3:** Support for command aliases
- **N4.4:** Plugin architecture for custom commands (future)
- **N4.5:** Command documentation auto-generated from code

---

## User Experience

### Core Workflows

#### Workflow 1: Discovering Commands
1. User types "/" in chat input
2. Command palette appears with all commands
3. User sees descriptions for each
4. User types partial name (e.g., "set")
5. Palette filters to matching commands
6. User selects "/settings" with Enter
7. Settings overlay opens

**Success Criteria:** User finds desired command within 5 seconds

#### Workflow 2: Quick Help Access
1. User needs help with keyboard shortcuts
2. Types "/help"
3. Help overlay opens immediately
4. User sees shortcuts and tips
5. User closes with Esc
6. Returns to chat

**Success Criteria:** User gets help in 2 seconds

#### Workflow 3: Stopping Agent
1. Agent is in middle of operation
2. User wants to cancel
3. Types "/stop" or presses Ctrl+C
4. Agent stops current iteration
5. Control returns to user
6. User can start new message

**Success Criteria:** Agent stops within 1 second

#### Workflow 4: Creating Git Commit
1. User has made changes via agent
2. Types "/commit"
3. Agent analyzes changed files
4. Shows proposed commit message
5. User reviews changes in diff overlay
6. User approves
7. Commit created

**Success Criteria:** Commit created with meaningful message

#### Workflow 5: Bash Mode
1. User needs to run shell commands
2. Types "/bash"
3. UI switches to bash mode
4. User types "ls -la"
5. Command executes (with approval)
6. Output shown in overlay
7. User types more commands or "/exit"

**Success Criteria:** Multiple commands executed easily

---

## Command Specifications

### /help
**Purpose:** Show help and tips  
**Syntax:** `/help`  
**Action:** Opens help overlay  
**Category:** Navigation  
**Approval Required:** No

**Details:**
- Displays keyboard shortcuts reference
- Shows available slash commands
- Provides quick tips for common tasks
- Searchable help content
- Links to full documentation

---

### /stop
**Purpose:** Cancel current agent operation  
**Syntax:** `/stop`  
**Action:** Stops agent loop  
**Category:** Control  
**Approval Required:** No

**Details:**
- Interrupts current iteration
- Preserves conversation history
- Returns control to user immediately
- Safe to use at any time
- Alternative to Ctrl+C

---

### /commit
**Purpose:** Create git commit  
**Syntax:** `/commit`  
**Action:** Analyzes changes and creates commit  
**Category:** Git  
**Approval Required:** Yes

**Details:**
- Scans workspace for modified files
- Generates commit message based on changes
- Shows diff preview
- Requests approval before committing
- Follows conventional commit format
- Validates git repository exists

---

### /pr
**Purpose:** Create pull request  
**Syntax:** `/pr`  
**Action:** Creates PR from current branch  
**Category:** Git  
**Approval Required:** Yes

**Details:**
- Checks for uncommitted changes
- Generates PR title and description
- Shows branch diff
- Requests approval
- Pushes to remote if needed
- Opens PR on GitHub/GitLab
- Requires git remote configuration

---

### /settings
**Purpose:** Open settings  
**Syntax:** `/settings`  
**Action:** Opens settings overlay  
**Category:** Configuration  
**Approval Required:** No

**Details:**
- Multi-tab settings interface
- General, LLM, Auto-Approval, Display tabs
- Navigate with Tab/Shift+Tab
- Changes save automatically
- Validates settings before applying
- Preserves settings across sessions

---

### /context
**Purpose:** Show context information  
**Syntax:** `/context`  
**Action:** Opens context overlay  
**Category:** Information  
**Approval Required:** No

**Details:**
- Workspace statistics
- Conversation history metrics
- Token usage breakdown
- Memory state visualization
- Session totals
- Provider information

---

### /bash
**Purpose:** Enter bash mode  
**Syntax:** `/bash`  
**Action:** Switches to shell command mode  
**Category:** Development  
**Approval Required:** Commands require approval

**Details:**
- Execute shell commands directly
- Real-time output streaming
- Command history (up/down arrows)
- Working directory shown in prompt
- Exit with `/exit` or Ctrl+D
- All commands require approval
- Sandbox to workspace directory

---

## Technical Architecture

### Component Structure

```
Slash Command System
├── Command Registry
│   ├── Command Definitions
│   ├── Command Metadata
│   └── Command Handlers
├── Command Parser
│   ├── Input Detector
│   ├── Command Matcher
│   └── Parameter Extractor
├── Autocomplete Engine
│   ├── Fuzzy Matcher
│   ├── Suggestion Ranker
│   └── Palette Renderer
├── Command Executor
│   ├── Validation
│   ├── Execution
│   └── Result Handling
└── Bash Mode
    ├── Shell Interface
    ├── Command Executor
    └── Output Streamer
```

### Command Registration

```go
type Command struct {
    Name        string
    Aliases     []string
    Description string
    Category    CommandCategory
    Handler     CommandHandler
    NeedsApproval bool
    Hidden      bool
}

type CommandRegistry struct {
    commands map[string]*Command
    aliases  map[string]string
}

func (r *CommandRegistry) Register(cmd *Command) error
func (r *CommandRegistry) Execute(name string, args []string) error
func (r *CommandRegistry) Autocomplete(prefix string) []*Command
```

### Execution Flow

```
User Input: "/comm"
    ↓
Command Detector: Matches "/"
    ↓
Autocomplete Engine: Filters commands → "/commit"
    ↓
User Selects: Enter key
    ↓
Command Parser: Extracts command + args
    ↓
Command Validator: Checks requirements
    ↓
Command Handler: Execute logic
    ↓
UI Update: Show result/overlay
    ↓
Return to Chat: Normal mode
```

---

## Design Decisions

### Why Slash Commands vs Other Approaches?

**Alternatives Considered:**
1. **Ctrl+Key Shortcuts:** Hard to discover, limited keys available
2. **Menu System:** Requires multiple steps, slower
3. **Natural Language:** Ambiguous, requires AI parsing
4. **Command Palette (Ctrl+P):** Extra keystroke, less integrated

**Why Slash Commands Won:**
- Familiar pattern from Slack, Discord, VS Code
- Self-documenting through autocomplete
- Fast to type (single character trigger)
- Integrates seamlessly with chat input
- Easy to extend with new commands

### Why No Parameters in Commands?

**Current Decision:** Commands are parameter-less triggers

**Rationale:**
- Simpler mental model
- Easier autocomplete
- No quoting/escaping complexity
- Overlays provide better UI for parameters
- Bash mode available for complex operations

**Future:** May add simple parameters (e.g., `/commit "message"`) if needed

### Why Separate Bash Mode vs Shell Command?

**Alternatives:**
1. Single command: `/shell git status`
2. Always inline: Execute shell commands directly
3. Shell overlay: Dedicated shell interface

**Why Bash Mode:**
- Supports multiple sequential commands
- Clear mode distinction
- Better for extended shell work
- Streaming output works better
- Can maintain shell state/history

---

## Command Categories

### Navigation Commands
- `/help` - Access help and documentation

### Control Commands
- `/stop` - Stop current operation

### Git Commands
- `/commit` - Create commit
- `/pr` - Create pull request

### Configuration Commands
- `/settings` - Modify settings

### Information Commands
- `/context` - View context info

### Development Commands
- `/bash` - Shell command mode

---

## Success Metrics

### Adoption Metrics
- **Usage rate:** >60% of users use at least one slash command per session
- **Command frequency:** Average 5+ slash commands per hour-long session
- **Discovery rate:** >80% of users discover commands within first session
- **Help access:** >50% of users access `/help` in first session

### Efficiency Metrics
- **Time to execute:** p95 under 3 seconds from "/" to result
- **Autocomplete usage:** >70% of commands selected via autocomplete
- **Error rate:** <5% of slash commands result in error
- **Typo tolerance:** >90% of typos corrected by autocomplete

### Feature Metrics
- **Most used:** Top 3 commands account for >60% of usage
- **Bash mode:** >30% of users enter bash mode at least once
- **Settings access:** >40% use `/settings` vs keyboard shortcut
- **Commit usage:** `/commit` used in >50% of coding sessions

---

## Dependencies

### External Dependencies
- TUI input handling (Bubble Tea)
- Git binary (for /commit, /pr)
- Shell access (for /bash mode)

### Internal Dependencies
- Agent core (for /stop command)
- Settings system (for /settings)
- Context manager (for /context)
- Tool approval system (for /commit, /pr, /bash)

### Platform Requirements
- Unix-like shell environment
- Git installed and configured
- Terminal with command history support

---

## Risks & Mitigations

### Risk 1: Command Discovery
**Impact:** Medium  
**Probability:** Medium  
**Mitigation:**
- Autocomplete appears immediately on "/"
- Help overlay prominently lists commands
- Toast notification on first launch: "Try typing /"
- Command palette shows descriptions

### Risk 2: Confusion with Regular Messages
**Impact:** Low  
**Probability:** Low  
**Mitigation:**
- Clear visual distinction (different color/style)
- Error message: "Unknown command, did you mean...?"
- Option to escape "/" as literal character
- Help text explains slash command concept

### Risk 3: Too Many Commands
**Impact:** Medium  
**Probability:** High (as features grow)  
**Mitigation:**
- Categorize commands in palette
- Fuzzy search for filtering
- Hide advanced commands by default
- Command aliases for brevity

### Risk 4: Bash Mode Security
**Impact:** High  
**Probability:** Low  
**Mitigation:**
- All commands require approval
- Clear indication of bash mode (different prompt)
- Sandbox to workspace directory
- Audit log of all executed commands

---

## Future Enhancements

### Phase 2 Ideas
- **Command Aliases:** `/s` for `/settings`, `/h` for `/help`
- **Simple Parameters:** `/commit "fix: typo"` for custom messages
- **Command History:** Recent commands accessible via up/down
- **Custom Commands:** User-defined slash commands
- **Command Macros:** Combine multiple commands

### Phase 3 Ideas
- **AI-Suggested Commands:** Agent suggests relevant commands in context
- **Command Chaining:** `/commit && /pr` to run sequence
- **Conditional Commands:** `/commit if changes`
- **Remote Commands:** Execute on remote servers
- **Plugin Commands:** Third-party command extensions

---

## Open Questions

1. **Should we support command parameters?**
   - Current: No parameters, just triggers
   - Future: Simple string parameters
   - Decision: Add if user feedback demands it

2. **Should we allow custom user commands?**
   - Pro: Power users can extend functionality
   - Con: Complexity, namespace collisions
   - Decision: Phase 2 feature if requested

3. **Should we have command categories in palette?**
   - Pro: Better organization with many commands
   - Con: More complex UI
   - Decision: Implement when >10 commands

4. **Should bash mode be persistent across sessions?**
   - Pro: Better for long-running operations
   - Con: Confusing if user forgets mode
   - Decision: Reset to chat mode on session end

---

## Related Documentation

- [Slash Commands Design](../plans/slash-commands-design.md)
- [How-to: Use TUI Interface - Slash Commands](../how-to/use-tui-interface.md#slash-commands)
- [TUI Executor PRD](tui-executor.md)
- [Bash Mode Implementation](../adr/0013-streaming-command-execution.md)

---

## Changelog

| Date | Version | Changes |
|------|---------|---------|
| 2024-12 | 1.0 | Initial PRD creation |
