# Product Requirements: Tool Approval System

**Feature:** Tool Approval Mechanism  
**Version:** 1.0  
**Status:** Implemented  
**Owner:** Security Team / Core Team  
**Last Updated:** December 2024

---

## Overview

The Tool Approval System provides a security-first mechanism that requires explicit human approval before the AI agent can execute potentially dangerous operations. This ensures developers maintain full control over what the agent does to their system, files, and environment while still enabling powerful automation capabilities.

---

## Problem Statement

AI coding agents that can execute arbitrary commands and modify files pose significant security and safety risks:

1. **Unintended Changes:** Agent might misunderstand requirements and make wrong changes
2. **Destructive Operations:** Accidental file deletions or overwrites
3. **Security Vulnerabilities:** Executing malicious commands or writing vulnerable code
4. **Trust Issues:** Developers hesitant to use AI tools with full system access
5. **Compliance Concerns:** Organizations need audit trails and approval workflows

Without approval mechanisms, users must choose between:
- **Full automation but no control:** Agent does whatever it wants (risky)
- **No automation:** Manual copy-paste of everything (defeats purpose)

---

## Goals

### Primary Goals

1. **Security by Default:** Require approval for all potentially dangerous operations
2. **Informed Decisions:** Provide users with all information needed to approve/deny
3. **Workflow Efficiency:** Make approval process quick and non-disruptive
4. **Flexibility:** Allow users to configure auto-approval for trusted operations
5. **Audit Trail:** Maintain clear record of what was approved/denied

### Non-Goals

1. **AI-Based Approval:** System does NOT use AI to decide what's safe
2. **Automated Learning:** Does NOT learn from user approval patterns (explicitly configured only)
3. **Multi-User Workflows:** Does NOT support approval delegation or team approvals
4. **Code Analysis:** Does NOT perform deep security analysis of proposed changes

---

## User Personas

### Primary: Security-Conscious Developer
- **Background:** Experienced developer who values security and control
- **Workflow:** Reviews all code changes carefully before committing
- **Pain Points:** Worried about AI making unintended changes
- **Goals:** Use AI assistance without giving up control over system

### Secondary: Productivity-Focused Developer
- **Background:** Developer who wants automation with minimal friction
- **Workflow:** Fast iteration, trusts tools with good defaults
- **Pain Points:** Approval prompts slow down workflow
- **Goals:** Automate trusted operations while maintaining safety net

### Tertiary: Team Lead / Engineering Manager
- **Background:** Responsible for team's code quality and security
- **Workflow:** Establishes guidelines and best practices
- **Pain Points:** Concerned about junior developers accepting dangerous operations
- **Goals:** Enable AI usage with appropriate guardrails

---

## Requirements

### Functional Requirements

#### FR1: Approval Detection
- **R1.1:** Identify which tools require approval
- **R1.2:** Support tool-level approval configuration
- **R1.3:** Allow operation-specific approval rules (e.g., write to certain paths)
- **R1.4:** Detect approval requirements before tool execution
- **R1.5:** Handle tools that sometimes need approval (conditional)

#### FR2: Approval Request UI
- **R2.1:** Display clear approval dialog when needed
- **R2.2:** Show tool name and purpose
- **R2.3:** Display all tool parameters with values
- **R2.4:** Highlight potentially dangerous values
- **R2.5:** Provide approve/deny/always options
- **R2.6:** Show preview of changes (especially for file operations)

#### FR3: User Response Handling
- **R3.1:** Accept keyboard input for approve/deny
- **R3.2:** Support "approve once" for current operation
- **R3.3:** Support "approve always" to create auto-approval rule
- **R3.4:** Support "deny" to block operation
- **R3.5:** Timeout option for approval requests (optional)
- **R3.6:** Cancel/go back option

#### FR4: Auto-Approval Rules
- **R4.1:** Allow users to configure trusted operations
- **R4.2:** Support path-based rules (e.g., "always approve reads from docs/")
- **R4.3:** Support command pattern matching (e.g., "approve 'git status'")
- **R4.4:** Support tool-specific rules (e.g., "auto-approve read_file")
- **R4.5:** Allow rule disabling without deletion
- **R4.6:** Provide rule management interface

#### FR5: Audit Logging
- **R5.1:** Log all approval requests
- **R5.2:** Log user decisions (approve/deny)
- **R5.3:** Record timestamp and context
- **R5.4:** Store tool parameters for audit
- **R5.5:** Support audit log export
- **R5.6:** Configurable log retention

#### FR6: Diff Preview (File Operations)
- **R6.1:** Show side-by-side diff for file changes
- **R6.2:** Syntax highlighting for code diffs
- **R6.3:** Line numbers and change indicators (+/-)
- **R6.4:** Support unified diff view option
- **R6.5:** Navigate large diffs easily
- **R6.6:** Show file path and operation type

#### FR7: Command Preview (Shell Execution)
- **R7.1:** Display full command to be executed
- **R7.2:** Show working directory
- **R7.3:** Highlight dangerous commands (rm, sudo, etc.)
- **R7.4:** Show environment variables if modified
- **R7.5:** Indicate timeout settings
- **R7.6:** Show expected execution context

#### FR8: Settings Integration
- **R8.1:** Approval settings accessible via settings overlay
- **R8.2:** Enable/disable approval system globally
- **R8.3:** Configure auto-approval rules
- **R8.4:** Set approval timeout (if enabled)
- **R8.5:** Manage audit log settings
- **R8.6:** Persist settings across sessions

### Non-Functional Requirements

#### NFR1: Security
- **N1.1:** Deny by default for all new tools
- **N1.2:** No bypass mechanisms (approval is mandatory)
- **N1.3:** Secure storage of approval rules
- **N1.4:** Validation of all tool parameters
- **N1.5:** Protection against injection attacks in previews

#### NFR2: Performance
- **N2.1:** Approval UI appears within 100ms of request
- **N2.2:** Diff generation under 500ms for typical files
- **N2.3:** Auto-approval rule checking under 10ms
- **N2.4:** No impact on agent loop when approval not needed
- **N2.5:** Efficient storage of audit logs

#### NFR3: Usability
- **N3.1:** Clear visual distinction of dangerous operations
- **N3.2:** Keyboard-accessible approval (no mouse required)
- **N3.3:** Intuitive approve/deny controls
- **N3.4:** Helpful error messages for denied operations
- **N3.5:** Easy rule creation from approval dialogs

#### NFR4: Reliability
- **N4.1:** Graceful handling of approval timeouts
- **N4.2:** Recovery from approval dialog crashes
- **N4.3:** Consistent approval behavior across sessions
- **N4.4:** Atomic audit log writes
- **N4.5:** Corruption-resistant settings storage

---

## User Experience

### Core Workflows

#### Workflow 1: First Tool Execution (No Auto-Approval)
1. Agent decides to use `write_file` tool
2. TUI shows "Agent wants to write file" notification
3. Approval overlay opens with:
   - Tool name: "write_file"
   - Parameters: path, content preview
   - Diff view (if file exists)
4. User reviews changes
5. User presses 'a' to approve or 'd' to deny
6. Tool executes (if approved)
7. Result displayed in chat

**Success Criteria:** User understands what will happen and makes informed decision

#### Workflow 2: Creating Auto-Approval Rule
1. Approval dialog appears for `read_file` on docs/ path
2. User sees this is safe operation
3. User presses 'A' for "Always approve"
4. Dialog prompts: "Create rule for read_file on docs/*?"
5. User confirms
6. Rule saved to settings
7. Future reads from docs/ auto-approve
8. Toast notification: "Auto-approval rule created"

**Success Criteria:** User can easily create rules for trusted operations

#### Workflow 3: Denying Dangerous Operation
1. Agent wants to execute `rm -rf /important/files`
2. Command preview shows full command
3. Red warning indicator for dangerous command
4. User presses 'd' to deny
5. Agent receives denial
6. Agent asks user for clarification or alternative approach

**Success Criteria:** User can easily block dangerous operations

#### Workflow 4: Managing Approval Rules
1. User opens settings with `/settings`
2. Navigates to "Auto-Approval" tab
3. Sees list of current rules
4. Selects rule to edit/disable/delete
5. Makes changes
6. Settings saved automatically
7. User closes settings

**Success Criteria:** User can review and modify approval rules easily

---

## Technical Architecture

### Component Structure

```
Tool Approval System
├── Approval Manager
│   ├── Request Queue
│   ├── Rule Matcher
│   ├── Decision Handler
│   └── Audit Logger
├── Approval Overlay (TUI)
│   ├── Dialog Renderer
│   ├── Diff Viewer
│   ├── Command Preview
│   └── Input Handler
├── Rule Engine
│   ├── Path Matcher
│   ├── Command Pattern Matcher
│   ├── Tool Matcher
│   └── Rule Evaluator
├── Settings Integration
│   ├── Rule Manager
│   ├── Settings Persistence
│   └── UI Components
└── Audit System
    ├── Log Writer
    ├── Log Storage
    └── Log Query
```

### Approval Flow

```
Agent Loop → Tool Call
    ↓
Approval Manager: Does tool need approval?
    ↓
┌────────────────────────────────────┐
│ Check Auto-Approval Rules          │
│ - Tool-level rules                 │
│ - Path-based rules                 │
│ - Command pattern rules            │
└────────────────┬───────────────────┘
                 ↓
         Rule Match Found?
         ├─ Yes → Auto-Approve → Execute Tool
         └─ No → Request User Approval
                 ↓
         Show Approval Overlay
                 ↓
         User Decision?
         ├─ Approve → Execute Tool
         ├─ Always → Create Rule + Execute
         └─ Deny → Return Error to Agent
```

### Data Model

```go
type ApprovalRequest struct {
    ID          string
    ToolName    string
    Parameters  map[string]interface{}
    Context     ExecutionContext
    Preview     *ChangePreview
    Timestamp   time.Time
}

type ApprovalRule struct {
    ID          string
    Type        RuleType // Tool, Path, Command
    Pattern     string
    Enabled     bool
    CreatedAt   time.Time
    LastUsed    time.Time
}

type AuditEntry struct {
    Timestamp   time.Time
    RequestID   string
    ToolName    string
    Decision    Decision // Approved, Denied, AutoApproved
    UserID      string
    Parameters  map[string]interface{}
}
```

---

## Design Decisions

### Why Approval Required for Write Operations?
- **Risk mitigation:** File writes can corrupt code or data
- **Reversibility:** Write operations are harder to undo
- **User control:** Developers want to review changes before applying
- **Industry standard:** Other tools (git, package managers) require confirmation

### Why Auto-Approval Rules Instead of AI Learning?
- **Predictability:** Explicit rules are clear and understandable
- **Control:** Users know exactly what's automated
- **Transparency:** No black box decision making
- **Debugging:** Easy to understand why something was/wasn't approved
- **Security:** No risk of AI learning wrong patterns

### Why Show Full Diff Instead of Summary?
- **Accuracy:** Summaries can miss important details
- **Trust:** Developers want to see exact changes
- **Learning:** Users learn what agent does by reviewing
- **Debugging:** Full context helps identify issues
- **Standard practice:** Developers are used to reviewing diffs

### Why Timeout is Optional?
- **Safety first:** Don't want to auto-approve by timeout
- **User pace:** Some decisions require careful consideration
- **Session integrity:** Long-running sessions shouldn't fail due to timeout
- **Flexibility:** Power users can enable if they want speed

---

## Security Considerations

### Threat Model

1. **Malicious Agent Behavior:**
   - Mitigation: All destructive operations require approval
   - Detection: Audit logs track all attempts

2. **Social Engineering:**
   - Mitigation: Clear previews show exact operations
   - Detection: Highlight dangerous patterns

3. **Accidental Approval:**
   - Mitigation: Different keys for approve vs always-approve
   - Recovery: Audit log shows what was approved

4. **Rule Exploitation:**
   - Mitigation: Strict pattern matching
   - Detection: Audit log shows rule usage

### Security Best Practices

1. **Deny by Default:** New tools require approval unless explicitly configured
2. **Least Privilege:** Auto-approval rules should be as specific as possible
3. **Audit Everything:** All approval decisions logged
4. **No Bypass:** No way to disable approval system entirely for dangerous tools
5. **Secure Defaults:** Default configuration prioritizes security over convenience

---

## Success Metrics

### Security Metrics
- **Prevention rate:** >99% of dangerous operations reviewed before execution
- **Rule accuracy:** <1% false positives (safe operations blocked)
- **Audit coverage:** 100% of tool executions logged
- **Override rate:** <5% of denials later regretted (indicates good default rules)

### Usability Metrics
- **Approval time:** p50 under 5 seconds for typical approvals
- **Rule creation rate:** >30% of users create at least one auto-approval rule
- **Denial rate:** <10% of approval requests denied (indicates agent making good suggestions)
- **Settings access:** >50% of users access approval settings within first week

### Trust Metrics
- **User confidence:** >85% of users feel in control of agent actions
- **Continued usage:** <5% of users disable agent due to approval friction
- **Feature adoption:** >70% of users understand and use approval system effectively

---

## Dependencies

### External Dependencies
- TUI framework (for overlay rendering)
- File system access (for diff generation)
- Settings system (for rule persistence)

### Internal Dependencies
- Agent core (event system for approval requests)
- Tool system (tool metadata and schemas)
- Memory system (for context in audit logs)
- Diff generation utilities

### Platform Requirements
- File system read access (for diff preview)
- Terminal with ANSI colors (for diff highlighting)
- Sufficient screen space (80x24 minimum for diff view)

---

## Risks & Mitigations

### Risk 1: Approval Fatigue
**Impact:** High  
**Probability:** Medium  
**Mitigation:**
- Smart default auto-approval rules for safe operations
- Easy rule creation from approval dialogs
- Batch approval for similar operations (future)
- Learn from usage patterns to suggest rules

### Risk 2: Users Auto-Approving Everything
**Impact:** High  
**Probability:** Low  
**Mitigation:**
- Warnings for overly broad rules
- Require confirmation for "always approve" on dangerous tools
- Audit log visibility to show rule usage
- Educational content about security

### Risk 3: Complex Rule Configuration
**Impact:** Medium  
**Probability:** Medium  
**Mitigation:**
- Simple default rules that work for most users
- Template rules for common scenarios
- Clear documentation with examples
- In-app help for rule creation

### Risk 4: Performance Impact
**Impact:** Low  
**Probability:** Low  
**Mitigation:**
- Efficient rule matching algorithms
- Cache rule evaluation results
- Async diff generation
- Optimize preview rendering

---

## Future Enhancements

### Phase 2 Ideas
- **Batch Approval:** Approve multiple similar operations at once
- **Conditional Rules:** "Approve if file size < 1000 lines"
- **Time-Based Rules:** "Auto-approve for next 1 hour"
- **Workspace Rules:** Different rules per workspace
- **Team Rules:** Share approved rules across team

### Phase 3 Ideas
- **AI-Assisted Review:** Highlight potentially problematic changes
- **Change Impact Analysis:** Show files affected by operation
- **Rollback Support:** Easy undo for approved operations
- **Approval Templates:** Pre-configured rule sets for common scenarios
- **Integration with Git:** Approve based on git diff/status

---

## Open Questions

1. **Should we support approval delegation?**
   - Use case: Junior dev gets approval from senior
   - Complexity: Requires multi-user support
   - Decision: Defer to future (Phase 3+)

2. **Should we limit auto-approval rule count?**
   - Pro: Prevents overly permissive configurations
   - Con: Power users may need many rules
   - Decision: Warn at 20+ rules but don't enforce limit

3. **Should audit logs be encrypted?**
   - Pro: Protects sensitive information
   - Con: Makes logs harder to query/debug
   - Decision: Start unencrypted, add encryption option in Phase 2

4. **Should we support approval plugins?**
   - Use case: Custom approval logic for specific tools
   - Complexity: Security implications
   - Decision: Research for Phase 3

---

## Related Documentation

- [ADR-0010: Tool Approval Mechanism](../adr/0010-tool-approval-mechanism.md)
- [ADR-0017: Auto-Approval and Settings System](../adr/0017-auto-approval-and-settings-system.md)
- [How-to: Use TUI Interface - Tool Approval](../how-to/use-tui-interface.md#tool-approval)
- [Architecture: Tool System](../architecture/tool-system.md)
- [Security Policy](../../SECURITY.md)

---

## Changelog

| Date | Version | Changes |
|------|---------|---------|
| 2024-12 | 1.0 | Initial PRD creation |
