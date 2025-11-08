# 9. Gemini-Inspired TUI Executor Design

**Status:** Accepted
**Date:** 2025-11-02
**Deciders:** Development Team
**Technical Story:** Implementation of a Terminal User Interface executor for Forge's flagship coding agent

---

## Context

Forge requires a Terminal User Interface (TUI) executor for its flagship coding agent to provide developers with a fluid, interactive, and visually clean experience during AI-assisted development sessions.

### Background

Modern AI coding assistants like Anthropic's Claude Code, Google's Gemini CLI, and GitHub Codex have established patterns for effective terminal-based interactions. These tools demonstrate that successful AI TUIs prioritize conversational flow and real-time feedback over complex multi-panel layouts.

### Problem Statement

The existing CLI executor provides functional interaction but lacks:
- Visual clarity for different types of agent output (thinking, tool calls, results)
- Real-time streaming feedback as the agent processes requests
- An aesthetically pleasing interface that matches modern AI tool standards
- Clear visual hierarchy to help users follow the agent's thought process

### Goals

- Create a single-pane, conversation-focused TUI that feels natural and unobtrusive
- Provide real-time streaming of agent thinking and responses
- Use color and typography to create clear visual hierarchy
- Match the aesthetic quality of leading AI CLI tools (Gemini CLI specifically)
- Maintain simplicity while providing all necessary context

### Non-Goals

- Multi-panel layouts with dedicated status, history, or file tree views
- Built-in file editing capabilities
- Complex keyboard shortcuts beyond basic navigation
- Support for simultaneous multiple conversations

---

## Decision Drivers

* **User Experience**: The interface should feel like a natural conversation with clear visual feedback
* **Visual Clarity**: Different content types (user input, thinking, tool calls) must be immediately distinguishable
* **Real-time Feedback**: Users should see agent activity as it happens, not just final results
* **Simplicity**: The interface should be intuitive without requiring documentation
* **Aesthetic Quality**: The UI should match or exceed the visual quality of competing tools

---

## Considered Options

### Option 1: Multi-Panel Layout

**Description:** A tmux-style layout with dedicated panels for conversation history, agent activity logs, file diffs, and system status.

**Pros:**
- High information density
- Simultaneous view of multiple aspects
- Familiar to power users

**Cons:**
- Complex to implement and maintain
- Overwhelming for new users
- Breaks conversational flow
- Doesn't match user's visual reference (Gemini CLI)

### Option 2: Single-Pane, Log-Style Layout

**Description:** A minimalist layout presenting all information in a single, continuous, scrollable view with inline formatting and real-time streaming.

**Pros:**
- Simple and intuitive
- Matches Gemini CLI aesthetic
- Natural conversational flow
- Easy to implement
- Focuses attention on current task

**Cons:**
- Requires scrolling to review history
- Lower information density
- No dedicated status area

### Option 3: Hybrid Minimal Layout

**Description:** Single-pane conversation with a minimal header showing context and a persistent input area.

**Pros:**
- Balances simplicity with context
- Provides session information without clutter
- Clean visual separation of sections

**Cons:**
- Slightly more complex than pure single-pane
- Header might be unnecessary for experienced users

---

## Decision

**Chosen Option:** Option 3 - Hybrid Minimal Layout

### Rationale

The hybrid approach provides the best balance between simplicity and functionality. The implementation includes:

1. **ASCII Art Header**: Welcoming branded header with Forge logo
2. **Minimal Context Bar**: Single line showing model and directory context
3. **Main Conversation Area**: Scrollable viewport for all agent interaction
4. **Persistent Input**: Salmon pink bordered input box for user queries
5. **Bottom Status**: Directory and model information

This design directly emulates the modern AI tool aesthetic while maintaining the simplicity of a single-pane conversation view. The minimal header and footer provide just enough context without cluttering the interface.

---

## Consequences

### Positive

- Clean, professional appearance matching modern AI tools
- Intuitive UX requiring no learning curve
- Real-time streaming creates engaging, responsive feel
- Simple implementation reduces maintenance burden
- Color-coded content types improve readability
- Pastel salmon pink (#FFB3BA) aesthetic is distinctive and pleasant

### Negative

- Conversation history requires scrolling to review
- No dedicated file tree or status panel for power users
- Limited screen real estate for context information
- Manual spacing adjustments needed for optimal visual flow

### Neutral

- Single conversation focus (one task at a time)
- Terminal-based (no mouse interaction beyond scrolling)
- Requires terminal with good Unicode and color support

---

## Implementation

### Technical Stack

- **TUI Framework**: `github.com/charmbracelet/bubbletea` (Elm-inspired architecture)
- **Components**: `github.com/charmbracelet/bubbles` (viewport, textarea)
- **Styling**: `github.com/charmbracelet/lipgloss` (terminal styling)

### Component Architecture

```
┌─────────────────────────────────────────┐
│  ASCII Art Header (Forge branding)     │
├─────────────────────────────────────────┤
│  Context Bar (model, directory)        │
├─────────────────────────────────────────┤
│                                         │
│  Scrollable Viewport                   │
│  - User prompts (👤 salmon pink)       │
│  - Thinking content (💭 gray italic)   │
│  - Tool calls (🔧 mint green)          │
│  - Tool results (✓ mint green)         │
│  - Messages (white)                     │
│                                         │
├─────────────────────────────────────────┤
│  Input Box (salmon pink border)        │
├─────────────────────────────────────────┤
│  Status Bar (directory, model)         │
└─────────────────────────────────────────┘
```

### Color Palette

- **Primary**: Salmon Pink (#FFB3BA) - User prompts, input border, branding
- **Accent**: Coral Pink (#FFCCCB) - Lighter accents
- **Thinking**: Muted Gray (#6B7280) - Agent thinking content
- **Tools**: Mint Green (#A8E6CF) - Tool calls and results
- **Background**: Dark (#111827)
- **Text**: Bright White (#F9FAFB)

### Key Features

1. **Real-time Streaming**: Thinking and message content streams character-by-character
2. **Custom Word Wrapping**: Manual word wrap implementation to avoid lipgloss padding issues
3. **Event-Driven Updates**: Agent events forwarded via channels to TUI
4. **Line-by-Line Styling**: Prevents block formatting issues
5. **Icon-Only Styling**: User input icon styled separately from content

### Migration Path

The TUI executor can be used alongside the existing CLI executor by changing the executor type in the main application:

```go
// Old
executor := cli.NewExecutor(agent)

// New
executor := tui.NewExecutor(agent)
```

---

## Validation

### Success Metrics

- Visual clarity: Different content types are immediately distinguishable
- Responsiveness: Agent thinking and responses stream in real-time
- Aesthetic quality: Matches or exceeds Gemini CLI visual standard
- Usability: Users can interact without referring to documentation

### Monitoring

- User feedback on visual design and usability
- Performance metrics for streaming updates
- Terminal compatibility across different environments

---

## Related Decisions

- [ADR-0005](0005-channel-based-agent-communication.md) - Channel-based agent communication enables real-time event streaming to TUI
- [ADR-0004](0004-agent-content-processing.md) - Content processing patterns used for formatting TUI output

---

## References

- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Lipgloss Styling Guide](https://github.com/charmbracelet/lipgloss)
- [Gemini CLI Design Reference](https://ai.google.dev/gemini-api/docs/cli)

---

## Notes

The implementation evolved through several iterations to achieve the desired spacing and styling:

1. Initial implementation had excessive vertical spacing
2. Lipgloss `.Width()` method was adding unwanted padding
3. Custom word wrapping function resolved formatting issues
4. Line-by-line styling prevented block formatting problems
5. Icon-only styling feature allows separate coloring of prefixes and content

The final implementation successfully achieves the Gemini-inspired aesthetic with clean spacing and proper color hierarchy.

**Last Updated:** 2025-11-05