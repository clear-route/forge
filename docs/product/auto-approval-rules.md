# Product Requirements: Auto-Approval Rules

**Feature:** Trust-Based Automation System  
**Version:** 1.0  
**Status:** Implemented  
**Owner:** Security Team / Core Team  
**Last Updated:** December 2024

---

## Overview

Auto-Approval Rules enable users to configure trusted operations that can execute without manual approval, balancing security with workflow efficiency. This system allows users to define granular rules for automatically approving specific tools, paths, commands, or patterns while maintaining comprehensive audit trails and safety guardrails.

---

## Problem Statement

Users working with AI coding agents face a tension between security and productivity:

1. **Approval Fatigue:** Repeatedly approving safe, repetitive operations slows workflow
2. **Lost Context:** Constant interruptions for approval break developer flow
3. **Inconsistent Trust:** No way to express "always trust this specific operation"
4. **Workflow Friction:** Legitimate, safe operations require same scrutiny as risky ones
5. **Time Waste:** Spending time approving obviously safe read operations
6. **Power User Frustration:** Experienced users want more control over automation

Without auto-approval, users either:
- Suffer constant interruptions (poor UX)
- Disable approval entirely (dangerous)
- Avoid using agent for repetitive tasks (underutilization)

---

## Goals

### Primary Goals

1. **Workflow Efficiency:** Enable automation of trusted operations
2. **Granular Control:** Provide precise rule definition capabilities
3. **Security First:** Maintain safety even with automation
4. **Transparency:** Clear visibility into what's automated and why
5. **Easy Management:** Simple rule creation, editing, and deletion
6. **Audit Trail:** Complete logging of all auto-approved actions

### Non-Goals

1. **AI-Based Trust:** Does NOT use AI to decide what's safe
2. **Learning System:** Does NOT automatically learn from user behavior
3. **Shared Rules:** Does NOT sync rules across users/teams (yet)
4. **Risk Analysis:** Does NOT perform deep security analysis of operations
5. **External Policy:** Does NOT integrate with organizational policy engines

---

## User Personas

### Primary: Power User Developer
- **Background:** Experienced developer who knows their workflow well
- **Workflow:** Repetitive, predictable operations during development
- **Pain Points:** Approval fatigue for known-safe operations
- **Goals:** Automate trusted operations while maintaining security

### Secondary: Team Lead
- **Background:** Sets up workflows for team members
- **Workflow:** Configures safe defaults for common tasks
- **Pain Points:** Needs consistency across team usage
- **Goals:** Define standard auto-approval rules for team

### Tertiary: New User
- **Background:** First-time Forge user, learning capabilities
- **Workflow:** Exploring features, building trust
- **Pain Points:** Overwhelmed by approval requests
- **Goals:** Start with good defaults, add rules as comfort grows

---

## Requirements

### Functional Requirements

#### FR1: Rule Types
- **R1.1:** Tool-based rules (auto-approve specific tools)
- **R1.2:** Path-based rules (auto-approve operations on specific paths)
- **R1.3:** Command pattern rules (auto-approve matching shell commands)
- **R1.4:** Composite rules (combine multiple conditions)
- **R1.5:** Exclusion rules (blacklist patterns)

#### FR2: Tool Rules
- **R2.1:** Auto-approve specific tools by name
  - Example: `read_file` always approved
- **R2.2:** Tool rules apply to all parameters
- **R2.3:** Can enable/disable per tool
- **R2.4:** Blacklist overrides whitelist
- **R2.5:** Default: No tools auto-approved except `read_file`

#### FR3: Path Rules
- **R3.1:** Whitelist path patterns (glob syntax)
  - Example: `docs/**` (all files under docs/)
- **R3.2:** Blacklist path patterns
  - Example: `!**/.env` (never auto-approve .env files)
- **R3.3:** Support wildcards (*, **, ?)
- **R3.4:** Path rules work with read/write/execute operations
- **R3.5:** Most specific rule wins (longest match)
- **R3.6:** Absolute and relative path support

#### FR4: Command Pattern Rules
- **R4.1:** Define regex patterns for shell commands
  - Example: `^git status$` (exact match)
  - Example: `^npm (test|build)` (multiple options)
- **R4.2:** Case-sensitive and case-insensitive options
- **R4.3:** Full regex support (capture groups, lookahead, etc.)
- **R4.4:** Blacklist dangerous commands (rm -rf, sudo, etc.)
- **R4.5:** Validate regex patterns before saving

#### FR5: Rule Creation
- **R5.1:** Create rules from approval dialog ("Always approve")
- **R5.2:** Create rules in settings interface
- **R5.3:** Create rules via slash command (future)
- **R5.4:** Quick-create common rules (templates)
- **R5.5:** Validate rules before activation
- **R5.6:** Provide examples for each rule type

#### FR6: Rule Management
- **R6.1:** List all active rules
- **R6.2:** Enable/disable rules without deletion
- **R6.3:** Edit existing rules
- **R6.4:** Delete rules
- **R6.5:** Reorder rules (priority)
- **R6.6:** Export/import rule sets
- **R6.7:** Search/filter rules

#### FR7: Rule Evaluation
- **R7.1:** Evaluate rules on every tool call
- **R7.2:** Check blacklist rules first
- **R7.3:** Then check whitelist rules
- **R7.4:** Default to manual approval if no match
- **R7.5:** Log evaluation decision
- **R7.6:** Performance: <10ms evaluation time

#### FR8: Safety Features
- **R8.1:** Maximum rule count (prevent overly permissive configs)
- **R8.2:** Dangerous operation warnings (even with rules)
- **R8.3:** Rule review prompt (every 30 days)
- **R8.4:** Audit log of all auto-approvals
- **R8.5:** Global kill switch (disable all auto-approval)
- **R8.6:** Rule impact preview (show what would auto-approve)

#### FR9: Default Rules
- **R9.1:** Ship with safe defaults
  - `read_file`: Always auto-approve reads
- **R9.2:** Suggest common rules based on usage (optional)
- **R9.3:** Project-specific rule templates
- **R9.4:** Easy reset to defaults

#### FR10: Audit & Logging
- **R10.1:** Log every auto-approval with:
  - Timestamp
  - Rule that matched
  - Tool/operation details
  - Full parameters
- **R10.2:** Searchable audit log
- **R10.3:** Export audit log
- **R10.4:** Retention policy (30 days default)
- **R10.5:** Visual indicator when auto-approval used

### Non-Functional Requirements

#### NFR1: Performance
- **N1.1:** Rule evaluation under 10ms per tool call
- **N1.2:** Rule creation/update under 100ms
- **N1.3:** Settings UI loads in under 200ms
- **N1.4:** No performance impact when no rules active
- **N1.5:** Efficient pattern matching (compiled regex)

#### NFR2: Security
- **N2.1:** Rules stored with restricted file permissions (600)
- **N2.2:** Rule validation prevents injection attacks
- **N2.3:** Blacklist always overrides whitelist
- **N2.4:** No way to bypass approval for truly dangerous operations
- **N2.5:** Clear warnings for risky rule configurations

#### NFR3: Usability
- **N3.1:** Intuitive rule creation interface
- **N3.2:** Clear preview of rule impact
- **N3.3:** Helpful error messages for invalid rules
- **N3.4:** Examples provided for each rule type
- **N3.5:** Easy to understand which rule matched

#### NFR4: Reliability
- **N4.1:** Invalid rules don't crash system
- **N4.2:** Graceful degradation (fall back to manual approval)
- **N4.3:** Atomic rule updates (no partial states)
- **N4.4:** Consistent evaluation across sessions
- **N4.5:** Rule file corruption recovery

---

## User Experience

### Core Workflows

#### Workflow 1: Auto-Approving Read Operations
1. User tired of approving file reads
2. Opens settings → Auto-Approval tab
3. Sees toggle: "Auto-approve read operations"
4. Enables toggle
5. Rule created: `tool:read_file → auto-approve`
6. Future file reads skip approval
7. User sees toast: "Auto-approved: read_file (rule: 'read operations')"

**Success Criteria:** No more interruptions for reads

#### Workflow 2: Creating Path-Based Rule from Dialog
1. Agent wants to read `docs/architecture.md`
2. Approval dialog appears
3. User sees "Always approve" option
4. User checks "Always approve for docs/**"
5. Dialog shows preview: "Will auto-approve all operations in docs/"
6. User confirms
7. Rule saved
8. Current operation auto-approves
9. Future docs/ operations auto-approve

**Success Criteria:** Rule created in context, immediately active

#### Workflow 3: Configuring Command Patterns
1. User frequently runs `npm test` and `npm build`
2. Opens settings → Auto-Approval tab
3. Clicks "Add Command Pattern Rule"
4. Enters regex: `^npm (test|build)$`
5. Tests pattern with examples
6. Saves rule
7. Future matching commands auto-approve

**Success Criteria:** Multiple related commands automated with one rule

#### Workflow 4: Reviewing and Disabling Rules
1. User concerned about security
2. Opens settings → Auto-Approval tab
3. Sees list of 5 active rules
4. Selects "Auto-approve npm install" rule
5. Clicks "Disable" (not delete)
6. Rule grayed out
7. Future npm install requires approval again

**Success Criteria:** Easy to temporarily disable without losing rule

#### Workflow 5: Audit Trail Review
1. User wonders what agent auto-approved
2. Opens settings → Auto-Approval tab
3. Clicks "View Audit Log"
4. Sees chronological list:
   - "10:23 AM - Auto-approved read_file docs/api.md (rule: 'docs/** whitelist')"
   - "10:24 AM - Auto-approved npm test (rule: 'npm commands')"
5. Can export log for review

**Success Criteria:** Complete transparency of auto-approvals

---

## Technical Architecture

### Component Structure

```
Auto-Approval System
├── Rule Engine
│   ├── Rule Store
│   ├── Rule Evaluator
│   ├── Pattern Matcher
│   └── Priority Resolver
├── Rule Manager
│   ├── CRUD Operations
│   ├── Validation Engine
│   ├── Import/Export
│   └── Migration Handler
├── Audit System
│   ├── Audit Logger
│   ├── Log Storage
│   ├── Query Engine
│   └── Retention Manager
├── Settings UI
│   ├── Rule List View
│   ├── Rule Editor
│   ├── Audit Viewer
│   └── Quick Actions
└── Integration
    ├── Approval System Hook
    ├── Tool System Integration
    └── Event Emitter
```

### Data Model

```go
type AutoApprovalRule struct {
    ID          string
    Type        RuleType
    Pattern     string
    Enabled     bool
    Priority    int
    CreatedAt   time.Time
    LastUsed    time.Time
    UseCount    int
    Description string
}

type RuleType int
const (
    RuleTypeTool RuleType = iota
    RuleTypePath
    RuleTypeCommand
    RuleTypeComposite
)

type RuleEvaluation struct {
    ToolName    string
    Parameters  map[string]interface{}
    MatchedRule *AutoApprovalRule
    Decision    ApprovalDecision
    Timestamp   time.Time
}

type ApprovalDecision int
const (
    DecisionManual ApprovalDecision = iota
    DecisionAutoApprove
    DecisionAutoReject
)

type AuditEntry struct {
    Timestamp   time.Time
    ToolName    string
    Decision    ApprovalDecision
    RuleID      string
    RuleName    string
    Parameters  map[string]interface{}
    UserID      string
}
```

### Rule Evaluation Flow

```
Tool Call Received
    ↓
Extract Tool & Parameters
    ↓
Auto-Approval Enabled?
    ├─ No → Manual Approval
    └─ Yes → Continue
        ↓
Load All Active Rules
        ↓
Sort by Priority
        ↓
┌────────────────────────────────┐
│ Evaluation Loop:               │
│                                │
│ For each rule:                 │
│   1. Check if rule applies     │
│   2. If blacklist → DENY       │
│   3. If whitelist → APPROVE    │
│   4. Continue to next rule     │
│                                │
│ No rule matched → Manual       │
└────────────────────────────────┘
        ↓
Log Decision to Audit
        ↓
Return Decision
        ↓
Execute (if approved) or Show Approval Dialog
```

---

## Design Decisions

### Why Whitelist + Blacklist Model?
**Rationale:**
- **Safety:** Blacklist always wins (can override broad whitelist)
- **Flexibility:** Whitelist enables automation
- **Security best practice:** Deny by default, explicit allows
- **User control:** Can carve out exceptions

**Example:**
- Whitelist: `**/*.md` (auto-approve all markdown files)
- Blacklist: `!**/secrets.md` (except secrets.md)

### Why NOT AI-Based Approval?
**Rationale:**
- **Predictability:** Users need to know exactly what's automated
- **Transparency:** AI decisions are black boxes
- **Security:** Can't audit AI reasoning
- **Trust:** Users want explicit control
- **Simplicity:** Easier to implement and debug

**Decision:** Explicit rules only, no AI learning

### Why Default Auto-Approve Read Operations?
**Rationale:**
- **Low risk:** Read operations are generally safe
- **High frequency:** Most common operation
- **User expectation:** Users expect to review writes, not reads
- **Productivity:** Biggest approval fatigue source

**Evidence:** 85% of approvals are for `read_file` in typical sessions

### Why 30-Day Rule Review Reminder?
**Rationale:**
- **Security hygiene:** Rules can become outdated
- **Awareness:** Users forget what they've automated
- **Best practice:** Regular security review
- **Balance:** Not too frequent to be annoying

---

## Rule Syntax Examples

### Tool Rules

```yaml
# Auto-approve all read operations
- type: tool
  pattern: read_file
  action: approve

# Auto-approve list operations  
- type: tool
  pattern: list_files
  action: approve

# Never auto-approve execute_command
- type: tool
  pattern: execute_command
  action: deny
```

---

### Path Rules

```yaml
# Auto-approve all operations in docs directory
- type: path
  pattern: docs/**
  action: approve

# Auto-approve reads in src, but not writes
- type: path
  pattern: src/**
  action: approve
  operations: [read]

# Never approve .env files
- type: path
  pattern: "**/.env*"
  action: deny
  priority: 1000  # High priority blacklist
```

---

### Command Rules

```yaml
# Auto-approve safe git commands
- type: command
  pattern: "^git (status|diff|log|show).*"
  action: approve

# Auto-approve npm/yarn scripts
- type: command
  pattern: "^(npm|yarn) (test|build|lint)$"
  action: approve

# Never approve dangerous commands
- type: command
  pattern: ".*(rm -rf|sudo|dd).*"
  action: deny
  priority: 1000
```

---

### Composite Rules

```yaml
# Auto-approve read operations in docs directory
- type: composite
  conditions:
    - tool: read_file
    - path: docs/**
  action: approve

# Deny writes to production configs
- type: composite
  conditions:
    - operation: write
    - path: "**/prod/**"
  action: deny
  priority: 900
```

---

## Success Metrics

### Adoption Metrics
- **Rule creation:** >60% of users create at least one auto-approval rule
- **Default acceptance:** >80% keep default "auto-approve reads" rule
- **Power user adoption:** >90% of users with 10+ sessions create rules
- **Rule count:** Average 3-5 rules per active user

### Efficiency Metrics
- **Approval reduction:** 60% fewer manual approvals with auto-approval
- **Time savings:** Average 30 seconds saved per auto-approved operation
- **Workflow interruption:** 70% reduction in approval dialog appearances
- **Session duration:** 25% longer sessions (less friction)

### Safety Metrics
- **Blacklist effectiveness:** 100% of blacklisted operations blocked
- **False positive rate:** <2% of auto-approvals were actually risky
- **Audit completeness:** 100% of auto-approvals logged
- **Rule review compliance:** >50% of users review rules when prompted

### Usability Metrics
- **Rule creation success:** >90% of rule creation attempts succeed
- **Rule understanding:** >85% of users understand what their rules do
- **Disable rate:** <10% of rules disabled due to unexpected behavior
- **Discovery:** >70% discover auto-approval within first 5 sessions

---

## Dependencies

### External Dependencies
- Glob pattern matching library
- Regex engine (standard library)
- File system access (for path validation)

### Internal Dependencies
- Tool approval system
- Settings system (for persistence)
- Audit logging system
- Event system (for notifications)

### Platform Requirements
- File system permissions
- Writable config directory
- JSON/YAML parsing

---

## Risks & Mitigations

### Risk 1: Overly Permissive Rules
**Impact:** Critical  
**Probability:** Medium  
**Mitigation:**
- Warn about broad rules (e.g., `**/*`)
- Require confirmation for risky patterns
- Rule impact preview before saving
- Regular review reminders
- Maximum rule count limit (20)
- Examples show specific, not broad rules

### Risk 2: Rule Complexity Confusion
**Impact:** Medium  
**Probability:** High  
**Mitigation:**
- Provide rule templates
- Interactive rule builder
- Clear examples for each type
- Preview what rule will match
- Help text in settings
- Validation with helpful errors

### Risk 3: Blacklist Bypass
**Impact:** Critical  
**Probability:** Low  
**Mitigation:**
- Blacklist always evaluated first
- No way to override blacklist with whitelist
- Hard-coded dangerous operation list
- Extensive testing of rule evaluation
- Security review of rule engine

### Risk 4: Performance Degradation
**Impact:** Medium  
**Probability:** Low  
**Mitigation:**
- Compile regex patterns once
- Cache rule evaluation results
- Limit number of rules (max 50)
- Optimize pattern matching
- Performance testing with many rules

### Risk 5: Audit Log Growth
**Impact:** Low  
**Probability:** Medium  
**Mitigation:**
- 30-day retention by default
- Automatic log rotation
- Size limits (10MB max)
- Export before deletion
- Configurable retention

---

## Future Enhancements

### Phase 2 Ideas
- **Smart Suggestions:** Suggest rules based on usage patterns
- **Rule Testing:** Test rules against historical operations
- **Conditional Rules:** Time-based, user-based conditions
- **Rule Groups:** Organize rules by project/task
- **Shared Rules:** Import rules from team/community

### Phase 3 Ideas
- **Machine Learning:** Learn safe patterns (with user approval)
- **Risk Scoring:** Automatic risk assessment of rules
- **Policy Integration:** Connect to organizational policy engines
- **Advanced Patterns:** More complex matching logic
- **Rule Debugging:** Detailed trace of rule evaluation

---

## Open Questions

1. **Should we suggest rules based on user behavior?**
   - Pro: Easier rule discovery
   - Con: Privacy concerns, AI complexity
   - Decision: Phase 2 feature with opt-in

2. **Should rules be workspace-specific?**
   - Pro: Different rules for different projects
   - Con: More complex management
   - Decision: Phase 2 if requested

3. **Should we support rule sharing/import from community?**
   - Pro: Faster setup for new users
   - Con: Security validation needed
   - Decision: Phase 3 with security review

4. **Should rule evaluation be customizable (plugins)?**
   - Pro: Ultimate flexibility
   - Con: Security nightmare
   - Decision: No - too risky

---

## Related Documentation

- [ADR-0017: Auto-Approval and Settings System](../adr/0017-auto-approval-and-settings-system.md)
- [Tool Approval System PRD](tool-approval-system.md)
- [Settings System PRD](settings-system.md)
- [How-to: Configure Auto-Approval Rules](../how-to/use-tui-interface.md#auto-approval-configuration)

---

## Changelog

| Date | Version | Changes |
|------|---------|---------|
| 2024-12 | 1.0 | Initial PRD creation |
