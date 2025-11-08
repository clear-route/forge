# Forge TUI Coding Agent - Implementation Plan

**Status:** Planning Complete  
**Date:** 2025-01-05  
**Version:** 1.0  

---

## Executive Summary

This document outlines the complete plan for building a flagship TUI-based coding agent for the Forge framework. The agent will provide file operations, code editing, and command execution capabilities through an intuitive chat-first interface with rich visual components for code review.

## Vision

Build a coding agent that rivals Claude Code and Cursor but runs entirely in the terminal, leveraging Forge's open-source, extensible architecture. The agent should feel natural for developers, with a chat-first interface enhanced by sophisticated diff viewing and approval workflows.

---

## Architecture Overview

### System Components

```mermaid
graph TB
    subgraph User Layer
        User[Developer]
    end
    
    subgraph TUI Layer
        TUI[Enhanced TUI Executor]
        Chat[Chat View]
        Diff[Diff Overlay]
        Tree[File Tree Overlay]
        Cmd[Command Output Overlay]
    end
    
    subgraph Agent Layer
        Agent[Forge Agent]
        Loop[Agent Loop with Approval]
        Memory[Auto-Pruning Memory]
    end
    
    subgraph Tools Layer
        Read[ReadFileTool]
        Write[WriteFileTool]
        List[ListFilesTool]
        Search[SearchFilesTool]
        ApplyDiff[ApplyDiffTool]
        Execute[ExecuteCommandTool]
    end
    
    subgraph Security Layer
        Guard[WorkspaceGuard]
        Validator[PathValidator]
    end
    
    subgraph Application
        Main[cmd/forge/main.go]
    end
    
    User --> TUI
    TUI --> Chat
    TUI -.-> Diff
    TUI -.-> Tree
    TUI -.-> Cmd
    
    TUI --> Agent
    Agent --> Loop
    Loop --> Memory
    
    Loop --> Read
    Loop --> Write
    Loop --> List
    Loop --> Search
    Loop --> ApplyDiff
    Loop --> Execute
    
    Read --> Guard
    Write --> Guard
    List --> Guard
    Search --> Guard
    ApplyDiff --> Guard
    Execute --> Guard
    
    Guard --> Validator
    
    Main --> Agent
    Main --> TUI
```

### Key Design Decisions

1. **Tool Approval Flow** (ADR-0010)
   - Event-based approval mechanism
   - User can approve/reject file changes and commands
   - Diff preview before execution
   - 5-minute timeout for approval requests

2. **Coding Tools** (ADR-0011)
   - Reusable tools in `pkg/tools/coding/`
   - Workspace-only file access enforced
   - Diff-based editing (search/replace) to avoid full rewrites
   - Tools work across any executor (TUI, CLI, API)

3. **Enhanced TUI** (ADR-0012)
   - Chat-first interface with dynamic overlays
   - Side-by-side diff viewer with syntax highlighting
   - File tree navigation on demand
   - Command output display during execution
   - Keyboard-driven workflow

4. **Security**
   - All file operations restricted to current working directory
   - Path validation prevents traversal attacks
   - No command whitelist/blacklist initially (trust + approval flow)

5. **Context Management**
   - Automatic memory pruning when approaching token limits
   - Agent loop tracks context size
   - Simple threshold-based pruning (no over-engineering)

---

## Implementation Phases

### Phase 1: Foundation (Weeks 1-2)

**Architecture & Security**
- [ ] Write ADR-0010 (Tool Approval Mechanism) ✅
- [ ] Write ADR-0011 (Coding Tools Architecture) ✅
- [ ] Write ADR-0012 (Enhanced TUI Executor) ✅
- [ ] Create comprehensive plan document ✅
- [ ] Create `pkg/security/workspace/` package
- [ ] Implement `WorkspaceGuard` with path validation
- [ ] Implement `PathValidator` for security checks
- [ ] Write comprehensive security tests

**Core Tools**
- [ ] Create `pkg/tools/coding/` package structure
- [ ] Implement `ReadFileTool` with line range support
- [ ] Implement `WriteFileTool` with path validation
- [ ] Write tests for Read/Write tools

### Phase 2: Advanced Tools (Weeks 3-4)

**Search & Navigation**
- [ ] Implement `ListFilesTool` with recursive/pattern support
- [ ] Implement `SearchFilesTool` with regex and context
- [ ] Write tests for List/Search tools

**Diff-Based Editing**
- [ ] Implement `ApplyDiffTool` with search/replace logic
- [ ] Add `Previewable` interface for diff preview
- [ ] Implement diff generation and validation
- [ ] Write comprehensive diff tool tests

**Command Execution**
- [ ] Implement `ExecuteCommandTool` with timeout
- [ ] Add output streaming via events
- [ ] Implement workspace-only execution
- [ ] Write command execution tests

### Phase 3: Agent Loop Enhancements (Week 5)

**Approval Mechanism**
- [ ] Add approval event types to `pkg/types/`
- [ ] Create approval response channel
- [ ] Modify agent loop to emit approval requests
- [ ] Implement approval timeout (5 minutes)
- [ ] Handle approval/rejection responses
- [ ] Write approval flow tests

**Context Management**
- [ ] Add context size tracking to agent loop
- [ ] Implement automatic pruning threshold check
- [ ] Call memory pruning when approaching limit
- [ ] Test pruning behavior with large conversations

### Phase 4: TUI Enhancements (Weeks 6-7)

**Overlay Infrastructure**
- [ ] Add overlay state to TUI model
- [ ] Implement overlay mode switching
- [ ] Create overlay base components
- [ ] Add keyboard shortcut system

**Diff Viewer**
- [ ] Create `DiffViewer` component
- [ ] Implement side-by-side panes
- [ ] Integrate Chroma for syntax highlighting
- [ ] Add accept/reject controls (Ctrl+A, Ctrl+R)
- [ ] Wire diff viewer to approval events
- [ ] Test diff viewer with various file types

**File Tree**
- [ ] Create `FileTree` component
- [ ] Implement directory tree building
- [ ] Add expand/collapse navigation
- [ ] Add keyboard shortcuts (j/k, Enter)
- [ ] Test file tree with large directories

**Command Output**
- [ ] Create `CommandOutput` component
- [ ] Implement real-time output streaming
- [ ] Add ANSI color support
- [ ] Add scrolling controls
- [ ] Test with various commands

### Phase 5: Integration (Week 8)

**Main Application**
- [ ] Create `cmd/forge/` directory structure
- [ ] Implement main.go with CLI argument parsing
- [ ] Initialize LLM provider from config
- [ ] Create agent with coding tools registered
- [ ] Initialize enhanced TUI executor
- [ ] Add graceful shutdown handling
- [ ] Test end-to-end workflow

**System Prompt**
- [ ] Design comprehensive coding agent system prompt
- [ ] Include tool usage guidelines
- [ ] Add coding best practices
- [ ] Test prompt effectiveness with various tasks

### Phase 6: Testing & Documentation (Week 9)

**Integration Testing**
- [ ] Write integration tests for full coding workflows
- [ ] Test file read/write/edit scenarios
- [ ] Test command execution with approval
- [ ] Test diff viewer with real code changes
- [ ] Test error handling and recovery
- [ ] Test security boundary enforcement

**Documentation**
- [ ] Create coding agent user guide
- [ ] Document tool schemas and usage
- [ ] Create example coding workflows
- [ ] Document keyboard shortcuts
- [ ] Add troubleshooting guide
- [ ] Update main README with coding agent info

---

## Technical Specifications

### Tool Specifications

#### ReadFileTool
```
Parameters:
  - path (required): File path relative to workspace
  - line_range (optional): "start-end" for partial reads
Returns: Line-numbered file content
Security: Workspace-only access
```

#### WriteFileTool
```
Parameters:
  - path (required): File path relative to workspace
  - content (required): Complete file content
Returns: Success message with file info
Security: Workspace-only access
Preview: Shows diff if file exists
```

#### ListFilesTool
```
Parameters:
  - path (optional): Directory path (default: ".")
  - recursive (optional): Boolean for recursive listing
  - pattern (optional): Glob pattern filter
Returns: Formatted file/directory listing
Security: Workspace-only access
```

#### SearchFilesTool
```
Parameters:
  - pattern (required): Regex search pattern
  - path (optional): Directory to search (default: ".")
  - file_pattern (optional): File glob filter
  - context_lines (optional): Lines of context (default: 2)
Returns: Matches with surrounding context
Security: Workspace-only access
```

#### ApplyDiffTool
```
Parameters:
  - path (required): File to modify
  - search (required): Exact text to find
  - replace (required): Replacement text
Returns: Success message
Security: Workspace-only access
Preview: Shows diff before applying
```

#### ExecuteCommandTool
```
Parameters:
  - command (required): Command to execute
  - working_dir (optional): Relative working directory
Returns: Command output (stdout/stderr)
Security: Runs in workspace, has timeout
Requires: User approval via overlay
```

### Event Types

```go
// New event types for approval flow
const (
    EventTypeToolApprovalRequest  EventType = "tool_approval_request"
    EventTypeToolApprovalResponse EventType = "tool_approval_response"
    EventTypeToolRejected        EventType = "tool_rejected"
)
```

### TUI Keyboard Shortcuts

**Conversation Mode:**
- `Enter` - Send message
- `Ctrl+C` / `Esc` - Quit
- `Ctrl+T` - Toggle file tree
- `Ctrl+O` - Toggle command output

**Diff Overlay:**
- `j/k` or `↓/↑` - Navigate diff lines
- `Ctrl+A` - Accept changes
- `Ctrl+R` - Reject changes
- `Esc` - Cancel (reject)

**File Tree:**
- `j/k` or `↓/↑` - Navigate files
- `Enter` - Expand/collapse or select
- `Esc` - Close overlay

**Command Output:**
- `j/k` or `↓/↑` - Scroll output
- `Esc` - Close overlay

---

## Dependencies

### New Go Dependencies
- `github.com/alecthomas/chroma/v2` - Syntax highlighting
- Already have: `github.com/charmbracelet/bubbletea` - TUI framework
- Already have: `github.com/charmbracelet/lipgloss` - Styling

### Package Structure

```
pkg/
├── security/
│   └── workspace/
│       ├── guard.go          # WorkspaceGuard implementation
│       ├── validator.go      # Path validation
│       └── workspace_test.go # Security tests
├── tools/
│   ├── coding/
│   │   ├── read_file.go
│   │   ├── write_file.go
│   │   ├── list_files.go
│   │   ├── search_files.go
│   │   ├── apply_diff.go
│   │   ├── execute_command.go
│   │   └── coding_test.go
│   ├── ask_question.go       # Existing
│   ├── converse.go          # Existing
│   ├── task_completion.go   # Existing
│   └── tool.go              # Existing interface
├── executor/
│   └── tui/
│       ├── executor.go       # Enhanced with overlays
│       ├── diff_viewer.go    # New component
│       ├── file_tree.go      # New component
│       ├── command_output.go # New component
│       └── tui_test.go
└── types/
    ├── event.go              # Add approval event types
    └── channels.go           # Add approval channel

cmd/
└── forge/
    ├── main.go               # Main application
    └── config.go             # Configuration
```

---

## Success Criteria

### Functional
- ✅ Agent can read files within workspace
- ✅ Agent can write new files and modify existing files
- ✅ Agent can search codebase with regex
- ✅ Agent can apply inline edits via diff
- ✅ Agent can execute commands with approval
- ✅ User can review diffs before accepting
- ✅ User can approve/reject command execution
- ✅ File tree shows workspace structure
- ✅ Syntax highlighting works for common languages

### Security
- ✅ All file operations stay within workspace
- ✅ Path traversal attempts are blocked
- ✅ Command execution requires user approval
- ✅ Approval timeout prevents hanging

### Performance
- ✅ File operations complete in <100ms
- ✅ Diff viewer renders in <200ms
- ✅ Syntax highlighting doesn't lag
- ✅ Memory pruning maintains context size

### UX
- ✅ Chat-first interface feels natural
- ✅ Diffs are easy to read and understand
- ✅ Keyboard shortcuts are intuitive
- ✅ Approval flow is smooth and quick
- ✅ Error messages are clear and actionable

---

## Future Enhancements (Post-V1)

### Git Integration
- Git status tool
- Git commit tool
- Git diff tool
- Show git status in file tree

### Multi-File Operations
- Batch file edits
- Atomic multi-file transactions
- Cross-file refactoring

### Advanced Features
- Code analysis tools
- Refactoring tools
- Test generation
- Documentation generation

### UX Improvements
- Customizable layouts
- Saved approval patterns
- Command history
- Session persistence

---

## Risk Mitigation

### Risk: Syntax highlighting performance
**Mitigation:** Use lazy highlighting, cache results, limit file size

### Risk: Approval flow complexity
**Mitigation:** Extensive testing, clear timeouts, simple state machine

### Risk: Security vulnerabilities
**Mitigation:** Comprehensive path validation tests, security review

### Risk: TUI complexity
**Mitigation:** Incremental overlay addition, unit test each component

---

## Timeline Summary

- **Weeks 1-2:** Foundation (Security, Core Tools, ADRs)
- **Weeks 3-4:** Advanced Tools (Search, Diff, Execute)
- **Week 5:** Agent Loop Enhancements (Approval, Pruning)
- **Weeks 6-7:** TUI Enhancements (Overlays, Diff Viewer)
- **Week 8:** Integration (Main App, System Prompt)
- **Week 9:** Testing & Documentation

**Total Estimated Time:** 9 weeks

---

## Conclusion

This plan provides a comprehensive roadmap for building the Forge TUI coding agent. The architecture is designed to be:
- **Secure:** Workspace boundaries enforced at tool level
- **Reusable:** Tools work across any executor
- **User-Friendly:** Chat-first with rich visual previews
- **Extensible:** Easy to add more tools and features

All major architectural decisions are documented in ADRs, and the implementation is broken into manageable phases with clear deliverables.

**Next Step:** Begin Phase 1 implementation with security layer and core tools.

---

**Document Version:** 1.0  
**Last Updated:** 2025-01-05  
**Maintained By:** Forge Core Team