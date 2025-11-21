# Product Requirements: TUI Executor

**Feature:** Terminal User Interface Executor  
**Version:** 1.0  
**Status:** Implemented  
**Owner:** Core Team  
**Last Updated:** December 2024

---

## Overview

The TUI Executor provides a rich, interactive terminal-based interface for the Forge coding agent. Unlike traditional CLI tools that simply print output, the TUI offers a modern chat-like experience with real-time updates, visual feedback, and keyboard-driven navigation - all without leaving the terminal.

---

## Problem Statement

Developers working with AI coding assistants face several challenges:

1. **Context Switching:** Moving between terminal, IDE, and web browser disrupts flow
2. **Poor Visibility:** CLI output is hard to scan, especially for code diffs and long results
3. **Limited Interactivity:** Traditional CLI can't show real-time updates or allow mid-stream interaction
4. **Discoverability:** Command-line tools have poor feature discoverability
5. **Visual Clutter:** Raw tool output mixed with conversation makes it hard to follow

Traditional CLI tools are too limiting, while switching to web UIs breaks the developer's terminal-centric workflow.

---

## Goals

### Primary Goals

1. **Seamless Terminal Integration:** Provide a rich UI that runs entirely in the terminal
2. **Real-Time Feedback:** Show agent thinking, tool execution, and results as they happen
3. **Visual Clarity:** Present code, diffs, and structured data in readable formats
4. **Keyboard Efficiency:** Enable power users to navigate and control without mouse
5. **Contextual Help:** Make features discoverable through UI elements and help systems

### Non-Goals

1. **Web UI:** This is NOT a web-based interface (see separate web UI PRD if planned)
2. **GUI Application:** This is NOT a standalone desktop application
3. **IDE Plugin:** This is NOT an IDE extension (though may inspire one)
4. **Multi-User:** This is NOT designed for concurrent multi-user sessions

---

## User Personas

### Primary: Terminal-Centric Developer
- **Background:** Senior/mid-level developer comfortable with command-line tools
- **Workflow:** Lives in terminal, uses vim/neovim or Emacs, minimal mouse usage
- **Pain Points:** Needs AI assistance but doesn't want to leave terminal
- **Goals:** Fast, keyboard-driven AI interaction with minimal context switching

### Secondary: Full-Stack Developer
- **Background:** Uses mix of IDE and terminal, familiar with modern dev tools
- **Workflow:** Terminal for git/deployment, IDE for coding
- **Pain Points:** Switching between tools for AI help is disruptive
- **Goals:** Quick access to AI assistance within existing terminal workflow

### Tertiary: DevOps Engineer
- **Background:** Automation-focused, heavy shell script and infrastructure-as-code user
- **Workflow:** SSH sessions, remote servers, terminal multiplexers (tmux/screen)
- **Pain Points:** No GUI available on remote servers, needs terminal-only tools
- **Goals:** Reliable AI assistant that works in any terminal environment

---

## Requirements

### Functional Requirements

#### FR1: Chat Interface
- **R1.1:** Display conversation history in scrollable viewport
- **R1.2:** Show user messages, agent responses, and tool executions
- **R1.3:** Support markdown-style formatting for code blocks
- **R1.4:** Auto-scroll to bottom on new messages
- **R1.5:** Allow manual scrolling to review history
- **R1.6:** Preserve scroll position when new content arrives (if user scrolled up)

#### FR2: Input System
- **R2.1:** Multi-line text input area for user messages
- **R2.2:** Support for long messages with automatic wrapping
- **R2.3:** Command history navigation (up/down arrows)
- **R2.4:** Slash command detection and execution
- **R2.5:** Input validation before submission
- **R2.6:** Clear indication of input focus state

#### FR3: Real-Time Updates
- **R3.1:** Stream agent responses as they arrive (token-by-token)
- **R3.2:** Show "thinking" indicator when agent is processing
- **R3.3:** Display tool execution status in real-time
- **R3.4:** Update UI when tool results arrive
- **R3.5:** Handle long-running operations with progress indication

#### FR4: Visual Elements
- **R4.1:** Syntax highlighting for code snippets
- **R4.2:** Distinct styling for user vs agent messages
- **R4.3:** Visual indicators for different message types (thinking, tool call, result)
- **R4.4:** Color-coded status indicators (success, error, warning)
- **R4.5:** Borders and spacing for readability

#### FR5: Overlay System
- **R5.1:** Support multiple overlay types (help, settings, approval, etc.)
- **R5.2:** Modal overlays that dim background content
- **R5.3:** Keyboard navigation within overlays
- **R5.4:** Close overlays with ESC or dedicated key
- **R5.5:** Maintain chat state when overlay is active
- **R5.6:** Stack overlays when needed (e.g., help on top of settings)

#### FR6: Status Bar
- **R6.1:** Show current mode (chat, bash, overlay active)
- **R6.2:** Display token usage information
- **R6.3:** Show keyboard hints for current context
- **R6.4:** Indicate connection status (if applicable)
- **R6.5:** Show workspace path

#### FR7: Keyboard Navigation
- **R7.1:** Comprehensive keyboard shortcuts for all operations
- **R7.2:** Vi-style navigation options for power users
- **R7.3:** Standard navigation (arrows, page up/down, home/end)
- **R7.4:** Tab completion for slash commands
- **R7.5:** Keyboard shortcut help accessible via overlay

#### FR8: Toast Notifications
- **R8.1:** Non-intrusive notifications for background events
- **R8.2:** Auto-dismiss after timeout
- **R8.3:** Different styles for info, success, warning, error
- **R8.4:** Stack multiple toasts if needed
- **R8.5:** Manual dismiss option

### Non-Functional Requirements

#### NFR1: Performance
- **N1.1:** Render updates within 16ms (60 FPS target)
- **N1.2:** Handle viewport with 10,000+ lines of history without lag
- **N1.3:** Start up in under 500ms on typical hardware
- **N1.4:** Memory usage under 100MB for typical session
- **N1.5:** Smooth scrolling even with rapid updates

#### NFR2: Compatibility
- **N2.1:** Work on all major terminal emulators (iTerm2, Terminal.app, Windows Terminal, etc.)
- **N2.2:** Support standard terminal sizes (80x24 minimum, up to 200x100+)
- **N2.3:** Graceful degradation on terminals with limited color support
- **N2.4:** Work over SSH with minimal latency impact
- **N2.5:** Compatible with terminal multiplexers (tmux, screen)

#### NFR3: Accessibility
- **N3.1:** All features accessible via keyboard (no mouse required)
- **N3.2:** Clear focus indicators
- **N3.3:** Readable with different color schemes
- **N3.4:** Support for high-contrast modes
- **N3.5:** Screen reader compatible (where terminal emulator supports it)

#### NFR4: Reliability
- **N4.1:** Gracefully handle terminal resize events
- **N4.2:** Recover from render errors without crashing
- **N4.3:** Preserve session state on unexpected termination
- **N4.4:** Handle rapid input without dropping characters
- **N4.5:** Proper cleanup on exit (restore terminal state)

#### NFR5: Usability
- **N5.1:** Intuitive for users familiar with modern chat applications
- **N5.2:** Discoverable features through help system
- **N5.3:** Consistent UI patterns across all overlays
- **N5.4:** Clear visual hierarchy
- **N5.5:** Minimal learning curve for basic usage

---

## User Experience

### Core Workflows

#### Workflow 1: Starting a Chat Session
1. User launches `forge` command
2. TUI renders with welcome message
3. Input cursor is focused and ready
4. User types message and presses Enter
5. Message appears in viewport
6. Agent response streams in real-time
7. User continues conversation

**Success Criteria:** User can start chatting within 5 seconds of launch

#### Workflow 2: Reviewing Tool Execution
1. Agent decides to use a tool
2. TUI shows "Agent is using [tool_name]" indicator
3. Tool approval overlay appears (if required)
4. User reviews and approves/denies
5. Tool executes, result appears in chat
6. Agent continues reasoning
7. User can scroll back to review tool details

**Success Criteria:** User can understand what tool was used and why

#### Workflow 3: Accessing Settings
1. User types `/settings` or presses Ctrl+,
2. Settings overlay opens
3. User navigates tabs with Tab/Shift+Tab
4. User modifies settings with arrow keys and Enter
5. Changes are saved automatically
6. User closes overlay with ESC
7. Chat resumes

**Success Criteria:** User can change settings in under 30 seconds

#### Workflow 4: Viewing Context Info
1. User types `/context` or presses Ctrl+I
2. Context overlay displays current session info
3. User reviews token usage, workspace, history
4. User closes overlay
5. User adjusts behavior based on context info

**Success Criteria:** User understands current session state clearly

---

## Technical Architecture

### Component Structure

```
TUI Executor
├── Main Event Loop (Bubble Tea)
├── Chat Viewport
│   ├── Message Renderer
│   ├── Scroll Controller
│   └── Syntax Highlighter
├── Input Area
│   ├── Text Buffer
│   ├── Command Parser
│   └── History Manager
├── Overlay Manager
│   ├── Base Overlay Component
│   ├── Help Overlay
│   ├── Settings Overlay
│   ├── Context Overlay
│   ├── Approval Overlay
│   ├── Diff Viewer Overlay
│   ├── Command Execution Overlay
│   └── Result List Overlay
├── Status Bar
├── Toast Manager
└── Event System
    ├── Agent Events
    ├── User Input Events
    └── System Events
```

### Key Technologies

- **Bubble Tea:** TUI framework (event-driven architecture)
- **Lip Gloss:** Styling and layout
- **Bubbles:** Reusable components (viewport, textarea, etc.)
- **Chroma:** Syntax highlighting
- **Go standard library:** Core functionality

### Event Flow

```
User Input → Input Handler → Command/Message Router
                                  ↓
                    ┌─────────────┴─────────────┐
                    ↓                           ↓
              Slash Command              Chat Message
                    ↓                           ↓
              Execute Locally           Send to Agent Loop
                    ↓                           ↓
              Update UI               Stream Events Back
                                                ↓
                                      Update UI in Real-Time
```

---

## Design Decisions

### Why Bubble Tea?
- **Mature ecosystem:** Well-tested framework with active community
- **Event-driven:** Natural fit for streaming LLM responses
- **Composable:** Easy to build complex UIs from simple components
- **Performance:** Efficient rendering with minimal overhead
- **Go-native:** Integrates seamlessly with rest of Forge codebase

### Why Terminal-Only?
- **Developer preference:** Target users live in terminal
- **Zero dependencies:** No browser or GUI framework needed
- **SSH-friendly:** Works over remote connections
- **Fast:** No browser overhead or network latency
- **Focus:** Keeps users in their workflow

### Why Modal Overlays vs Split Panes?
- **Screen space:** Overlays don't permanently consume space
- **Context:** Overlays clearly indicate separate mode
- **Simplicity:** Easier mental model than managing multiple panes
- **Flexibility:** Can show different overlays as needed
- **Mobile-like:** Familiar pattern from mobile UIs

---

## Success Metrics

### Adoption Metrics
- **Usage rate:** >80% of Forge users choose TUI over CLI
- **Session duration:** Average 15+ minutes per session (indicates engagement)
- **Retention:** >60% of users return within 7 days

### Performance Metrics
- **Render latency:** p95 frame time <16ms (60 FPS)
- **Input lag:** <50ms from keypress to visual feedback
- **Memory usage:** <100MB for typical session
- **CPU usage:** <5% idle, <30% during active streaming

### Usability Metrics
- **Feature discovery:** >70% find help overlay within first session
- **Error rate:** <5% of user inputs result in confusion or error
- **Task completion:** >90% of intended actions complete successfully
- **Learning curve:** Users can perform basic tasks within 5 minutes

### Quality Metrics
- **Crash rate:** <0.1% of sessions
- **Render errors:** <1% of updates
- **Terminal compatibility:** Works on >95% of tested emulators
- **User satisfaction:** >4.5/5 rating for UI/UX

---

## Dependencies

### External Dependencies
- Bubble Tea framework (v0.25+)
- Lip Gloss styling library
- Bubbles component library
- Chroma syntax highlighter

### Internal Dependencies
- Agent core (event system)
- Tool system (for approval overlays)
- Settings manager
- Memory system (for context display)

### Platform Requirements
- Terminal emulator with ANSI color support
- Minimum 80x24 terminal size
- Unicode support (recommended)
- Go 1.21+ runtime

---

## Risks & Mitigations

### Risk 1: Terminal Compatibility Issues
**Impact:** High  
**Probability:** Medium  
**Mitigation:**
- Test on all major terminal emulators
- Provide fallback rendering for limited terminals
- Document known compatibility issues
- Graceful degradation strategy

### Risk 2: Performance Degradation with Large History
**Impact:** Medium  
**Probability:** Medium  
**Mitigation:**
- Implement virtual scrolling for viewport
- Limit rendered history (window-based rendering)
- Provide history export/clear functionality
- Optimize render pipeline

### Risk 3: Complex UI Overwhelming Users
**Impact:** Medium  
**Probability:** Low  
**Mitigation:**
- Progressive disclosure of features
- Clear help documentation
- Onboarding flow for first-time users
- Keyboard shortcut cheat sheet

### Risk 4: SSH Latency Impact
**Impact:** Low  
**Probability:** Low  
**Mitigation:**
- Minimize render updates
- Batch UI changes where possible
- Provide "low-latency" mode
- Test over high-latency connections

---

## Future Enhancements

### Phase 2 Ideas
- **Themes:** User-selectable color themes
- **Split view:** Side-by-side code editing
- **Mouse support:** Optional mouse interaction
- **Custom keybindings:** User-configurable shortcuts
- **Session tabs:** Multiple concurrent sessions

### Phase 3 Ideas
- **Collaboration:** Shared TUI sessions
- **Plugins:** Custom overlay types
- **Macros:** Record and replay command sequences
- **Advanced search:** Full-text search across history
- **Export:** Save conversations in various formats

---

## Open Questions

1. **Should we support mouse interaction?**
   - Pro: Easier for less experienced users
   - Con: Adds complexity, not needed for target users
   - Decision: Defer to Phase 2, keyboard-first for now

2. **How much history should we keep in memory?**
   - Current: Unlimited (constrained by memory)
   - Alternative: Last 1000 messages, archive rest
   - Decision: TBD based on user feedback

3. **Should we support multiple themes out of the box?**
   - Current: Single default theme
   - Alternative: Light/dark/high-contrast
   - Decision: Start with one, add more based on demand

---

## Related Documentation

- [ADR-0009: TUI Executor Design](../adr/0009-tui-executor-design.md)
- [ADR-0012: Enhanced TUI Executor](../adr/0012-enhanced-tui-executor.md)
- [ADR-0021: Early Tool Call Detection](../adr/0021-early-tool-call-detection.md)
- [How-to: Use TUI Interface](../how-to/use-tui-interface.md)
- [Architecture: Event System](../architecture/events.md)

---

## Changelog

| Date | Version | Changes |
|------|---------|---------|
| 2024-12 | 1.0 | Initial PRD creation |
