# Product Requirements: Tool Approval System

**Feature:** Tool Approval Mechanism  
**Version:** 1.0  
**Status:** Implemented  
**Owner:** Security Team / Core Team  
**Last Updated:** December 2024

---

## Product Vision

Empower developers to use AI automation confidently by maintaining full control over what the agent can do to their system, files, and environment. Every potentially dangerous operation requires explicit human approval, ensuring safety without sacrificing the power of autonomous AI assistance.

**Strategic Alignment:** Trust is foundational to AI adoption—developers must feel in control. The approval system makes Forge safe enough for production use while maintaining powerful automation capabilities.

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

## Key Value Propositions

### For Security-Conscious Developers
- **Full Control:** Review every file change and command before execution
- **Transparent Operations:** See exactly what agent wants to do before approving
- **Safe Defaults:** Dangerous operations blocked by default, safe ones streamlined

### For Productivity-Focused Developers
- **Flexible Automation:** Auto-approve trusted operations to reduce friction
- **Smart Workflows:** Create approval rules for repetitive safe tasks
- **Fast Decisions:** Keyboard-driven approval process, no context switching

### For Team Leads & Managers
- **Compliance Ready:** Complete audit trail of all agent actions
- **Risk Management:** Organizational control over automation boundaries
- **Educational Value:** Team learns safe AI usage through approval patterns

---

## Target Users & Use Cases

### Primary: Security-Conscious Developer
**Profile:**
- Experienced developer who values security and control
- Reviews all code changes carefully before committing
- Works on production codebases or sensitive systems

**Key Use Cases:**
- Reviewing agent-proposed file changes before applying
- Blocking dangerous shell commands
- Maintaining audit trail for compliance

**Pain Points Addressed:**
- Worry about AI making unintended changes
- Fear of accidental destructive operations
- Lack of visibility into AI actions

---

### Secondary: Productivity-Focused Developer
**Profile:**
- Developer who wants automation with minimal friction
- Fast iteration workflow, trusts tools with good defaults
- Willing to configure for optimal productivity

**Key Use Cases:**
- Auto-approving safe read operations
- Creating rules for trusted write patterns
- Streamlining repetitive approval workflows

**Pain Points Addressed:**
- Approval prompts slowing down workflow
- Repetitive approval of safe operations
- Manual configuration overhead

---

### Tertiary: Team Lead / Engineering Manager
**Profile:**
- Responsible for team's code quality and security
- Establishes guidelines and best practices
- Manages compliance and audit requirements

**Key Use Cases:**
- Reviewing team's AI usage patterns through audit logs
- Setting organization-wide approval policies
- Ensuring junior developers don't approve dangerous operations

**Pain Points Addressed:**
- Lack of visibility into team's AI tool usage
- Compliance and security requirements
- Training team on safe AI assistance

---

## Product Requirements

### Priority 0 (Must Have)

#### P0-1: Visual Approval Interface
**Description:** Clear, informative interface showing what agent wants to do before execution

**User Stories:**
- As a developer, I want to see exactly what changes the agent will make before approving
- As a user, I want preview of file changes with syntax highlighting and diff view

**Acceptance Criteria:**
- Approval overlay appears automatically when agent requests dangerous operation
- File changes shown as side-by-side diff with syntax highlighting
- Shell commands displayed with full context (working directory, parameters)
- Tool parameters shown in readable format
- Approve/Deny buttons with keyboard shortcuts

**Example:**
```
┌─ Agent Approval Required ─────────────────┐
│ write_file                                │
│                                           │
│ File: src/auth.go                         │
│                                           │
│ Changes:                                  │
│ - func authenticate(user string)         │
│ + func authenticate(user string) error   │
│                                           │
│ [A] Approve  [D] Deny                     │
└───────────────────────────────────────────┘
```

---

#### P0-2: Keyboard-Driven Workflow
**Description:** Fast, keyboard-accessible approval without requiring mouse

**User Stories:**
- As a developer, I want to approve/deny without leaving the keyboard
- As a user, I want common actions mapped to intuitive shortcuts

**Acceptance Criteria:**
- Ctrl+A or Enter to approve
- Ctrl+R, Ctrl+C, or ESC to deny
- Tab to navigate between Approve/Deny buttons
- No mouse required for any approval action

---

#### P0-3: Automatic Safety Blocks
**Description:** System automatically requires approval for dangerous operations

**User Stories:**
- As a developer, I want dangerous operations blocked by default
- As a user, I want protection without manual configuration

**Acceptance Criteria:**
- File write operations (write_file, apply_diff) require approval
- Shell command execution (execute_command) requires approval
- Read-only operations (read_file, list_files, search_files) do not require approval by default
- New tools require approval unless explicitly configured as safe

---

#### P0-4: Clear Denial Feedback
**Description:** Agent receives clear feedback when operations are denied

**User Stories:**
- As a developer, I want the agent to understand when I reject an operation
- As a user, I want the agent to explain alternatives after denial

**Acceptance Criteria:**
- Denied operations return error to agent
- Agent acknowledges denial in conversation
- Agent offers alternative approaches or asks for clarification
- Conversation continues smoothly after denial

**Example:**
```
User denies file write
Agent: "I understand you don't want to modify that file. 
Would you like me to show you the proposed changes instead?"
```

---

### Priority 1 (Should Have)

#### P1-1: Auto-Approval Rules
**Description:** Users can configure trusted operations to auto-approve without prompts

**User Stories:**
- As a developer, I want to auto-approve safe operations I do frequently
- As a user, I want to reduce approval friction for trusted patterns

**Acceptance Criteria:**
- Settings interface to create auto-approval rules
- Support tool-level rules (e.g., "auto-approve all read_file")
- Support path-based rules (e.g., "auto-approve writes to test/ directory")
- Support command whitelist for execute_command
- Rules saved persistently across sessions
- Clear indication when auto-approval rule used

**Example Rules:**
- Auto-approve: read_file, list_files, search_files (safe read operations)
- Auto-approve: write_file to docs/ (documentation changes)
- Auto-approve: execute_command "git status" (safe git commands)

---

#### P1-2: Approval Timeout
**Description:** Configurable timeout for approval requests

**User Stories:**
- As a developer, I want control over how long approval waits
- As a user, I want sessions to fail gracefully if I step away

**Acceptance Criteria:**
- Configurable timeout setting (default: disabled)
- Timeout triggers automatic denial (safe default)
- Clear notification when timeout occurs
- Agent receives timeout feedback

---

#### P1-3: Audit Trail
**Description:** Complete log of all approval requests and decisions

**User Stories:**
- As a team lead, I want to review what operations were approved/denied
- As a developer, I want compliance-ready audit logs

**Acceptance Criteria:**
- Log all approval requests with timestamp
- Record user decisions (approved/denied)
- Include tool name and parameters
- Support log export for compliance
- Configurable log retention period

---

#### P1-4: Dangerous Command Warnings
**Description:** Visual warnings for particularly risky operations

**User Stories:**
- As a developer, I want clear warnings for destructive commands
- As a user, I want extra attention drawn to risky operations

**Acceptance Criteria:**
- Red/warning indicators for dangerous patterns (rm -rf, sudo, etc.)
- Warning text explaining specific risk
- Explicit confirmation required for most dangerous operations
- Educational messaging about risks

**Example:**
```
⚠️  WARNING: Destructive Command
   This will permanently delete files!
   
   Command: rm -rf /important/files
   
   [A] I understand, approve  [D] Deny
```

---

### Priority 2 (Nice to Have)

#### P2-1: Rule Management Interface
**Description:** Comprehensive interface for managing auto-approval rules

**User Stories:**
- As a power user, I want to review and edit all my approval rules
- As a developer, I want to enable/disable rules temporarily

**Acceptance Criteria:**
- Settings tab showing all active rules
- Enable/disable individual rules
- Edit rule patterns
- Delete unused rules
- See last used timestamp for each rule

---

#### P2-2: Batch Approval
**Description:** Approve multiple similar operations at once

**User Stories:**
- As a developer, I want to approve a batch of similar file changes
- As a user, I want to reduce repetitive approval prompts

**Acceptance Criteria:**
- Detect multiple similar pending operations
- Option to approve all at once
- Clear preview of what will be approved in batch
- Safety limit on batch size

---

#### P2-3: Approval Templates
**Description:** Pre-configured rule sets for common scenarios

**User Stories:**
- As a new user, I want recommended approval rules for my workflow
- As a developer, I want quick setup without manual configuration

**Acceptance Criteria:**
- Templates for common scenarios ("Documentation Work", "Testing", "Refactoring")
- One-click activation of template rules
- Clear explanation of what each template enables
- Ability to customize after applying template

---

## User Experience Flow

### First Approval Experience

```
Agent wants to write file
    ↓
Approval overlay appears
    ↓
User sees diff preview with changes
    ↓
User presses 'A' to approve
    ↓
File written, result shown in chat
    ↓
Agent continues task
```

**Experience:** Clear, informative—user feels in control

---

### Denial and Alternative Flow

```
Agent wants to execute dangerous command
    ↓
Approval overlay with red warning
    ↓
User reads command: "rm -rf node_modules"
    ↓
User presses 'D' to deny
    ↓
Agent: "Command denied. Would you like me to
       move files to trash instead?"
    ↓
User: "Yes, move to trash"
    ↓
Agent uses safer alternative
```

**Experience:** Safety-first with intelligent alternatives

---

### Auto-Approval Rule Creation

```
Approval request for read_file on docs/
    ↓
User realizes this is always safe
    ↓
User opens settings (Ctrl+,)
    ↓
Navigates to Auto-Approval tab
    ↓
Adds rule: "read_file docs/**"
    ↓
Saves settings
    ↓
Future docs/ reads auto-approve
```

**Experience:** Easy customization for trusted patterns

---

### Audit Review Flow

```
Team lead wants to review AI usage
    ↓
Opens audit log from settings
    ↓
Sees chronological list of approvals
    ↓
Filters by team member or date range
    ↓
Reviews denied operations for training
    ↓
Exports log for compliance
```

**Experience:** Transparency for oversight and compliance

---

## User Interface & Interaction Design

### Approval Overlay - File Write

**Visual Layout:**
```
┌─ Tool Approval Required ──────────────────────────┐
│ write_file                                        │
│                                                   │
│ File: src/components/Header.tsx                   │
│                                                   │
│ ┌─ Changes ───────────────────────────────────┐  │
│ │  1 | import React from 'react';             │  │
│ │  2 | import { Logo } from './Logo';         │  │
│ │  3 |                                        │  │
│ │  4 | export function Header() {             │  │
│ │ -5 |   return <header><Logo /></header>;   │  │
│ │ +5 |   return (                             │  │
│ │ +6 |     <header className="app-header">   │  │
│ │ +7 |       <Logo />                         │  │
│ │ +8 |     </header>                          │  │
│ │ +9 |   );                                   │  │
│ │ 10 | }                                      │  │
│ └─────────────────────────────────────────────┘  │
│                                                   │
│ ⌨  [Ctrl+A] Approve    [Ctrl+R] Reject           │
└───────────────────────────────────────────────────┘
```

**Design Principles:**
- Syntax-highlighted diff for code changes
- Line numbers for easy reference
- Clear +/- indicators for additions/deletions
- Keyboard shortcuts visible
- Scrollable for large changes

---

### Approval Overlay - Command Execution

**Visual Layout:**
```
┌─ Tool Approval Required ──────────────────────────┐
│ execute_command                                   │
│                                                   │
│ Command: npm run build                            │
│ Working Directory: /workspace/my-project          │
│ Timeout: 30 seconds                               │
│                                                   │
│ This command will:                                │
│ • Compile TypeScript files                        │
│ • Bundle assets                                   │
│ • Generate build/ directory                       │
│                                                   │
│ ⌨  [Ctrl+A] Approve    [Ctrl+R] Reject           │
└───────────────────────────────────────────────────┘
```

**Design Principles:**
- Clear command display
- Context information (directory, timeout)
- Human-readable explanation of effects
- Safe, professional appearance

---

### Approval Overlay - Dangerous Command

**Visual Layout:**
```
┌─ ⚠️  DANGEROUS OPERATION ─────────────────────────┐
│ execute_command                                   │
│                                                   │
│ ⚠️  WARNING: This command will permanently        │
│    delete files without recovery!                 │
│                                                   │
│ Command: rm -rf node_modules                      │
│ Working Directory: /workspace/my-project          │
│                                                   │
│ Affected: ~15,000 files in node_modules/          │
│                                                   │
│ ⚠️  Consider: Use trash instead of permanent      │
│    deletion, or run 'npm prune' for cleanup       │
│                                                   │
│ ⌨  [Ctrl+A] I understand, approve                │
│    [Ctrl+R] Reject (recommended)                  │
└───────────────────────────────────────────────────┘
```

**Design Principles:**
- Red/warning color scheme
- Explicit risk explanation
- Impact assessment (file count)
- Safer alternative suggestions
- Bias toward rejection in button labels

---

### Auto-Approval Settings Panel

**Visual Layout:**
```
┌─ Settings: Auto-Approval Rules ───────────────────┐
│                                                   │
│ ✅ Read Operations (recommended)                  │
│    Auto-approve: read_file, list_files,           │
│    search_files                                   │
│    [Disable]                                      │
│                                                   │
│ ✅ Documentation Writes                           │
│    Auto-approve: write_file docs/**               │
│    Last used: 2 hours ago                         │
│    [Disable] [Edit]                               │
│                                                   │
│ ✅ Safe Git Commands                              │
│    Auto-approve: git status, git log, git diff    │
│    Last used: 15 minutes ago                      │
│    [Disable] [Edit]                               │
│                                                   │
│ [+ Add New Rule]                                  │
│                                                   │
└───────────────────────────────────────────────────┘
```

**Design Principles:**
- Visual checkboxes show enabled/disabled
- Last used timestamp shows rule activity
- Quick enable/disable actions
- Edit option for customization
- Clear rule descriptions

---

## Success Metrics

### Security Metrics

**Dangerous Operation Prevention:**
- Target: >99% of dangerous operations reviewed before execution
- Measure: Track approval requests vs. auto-approved operations

**Denial Rate:**
- Target: <10% of approvals denied (indicates agent making good suggestions)
- Measure: Approved vs. denied requests

**Audit Coverage:**
- Target: 100% of tool executions logged
- Measure: Compare tool executions to audit entries

---

### Usability Metrics

**Approval Speed:**
- Target: p50 approval time under 5 seconds
- Measure: Time from request to decision

**Auto-Approval Adoption:**
- Target: >30% of users create at least one auto-approval rule
- Measure: User accounts with rules configured

**Rule Effectiveness:**
- Target: >50% reduction in approval prompts after rule creation
- Measure: Approval requests before vs. after rule configuration

---

### Trust Metrics

**User Confidence:**
- Target: >85% of users feel in control of agent actions
- Measure: Post-session surveys

**Feature Satisfaction:**
- Target: >80% satisfaction with approval system
- Measure: User ratings of approval experience

**Continued Usage:**
- Target: <5% of users disable agent due to approval friction
- Measure: Churn analysis with exit surveys

---

## User Enablement

### Discoverability

**First-Time Experience:**
- Tutorial on first approval: "This is how you control what the agent does"
- Tooltip hints on keyboard shortcuts
- Link to approval documentation in overlay

**Progressive Disclosure:**
- Beginner: See approval overlays, learn the flow
- Intermediate: Create first auto-approval rule
- Advanced: Configure comprehensive rule sets

---

### Learning Path

**Beginner:**
1. Experience first approval overlay
2. Understand diff preview and command display
3. Learn keyboard shortcuts (Ctrl+A, Ctrl+R)

**Intermediate:**
1. Identify repetitive safe approvals
2. Create first auto-approval rule
3. Monitor rule usage in settings

**Advanced:**
1. Configure comprehensive rule sets for workflow
2. Use approval templates for quick setup
3. Review audit logs for optimization

---

### Support Materials

**Documentation:**
- "Understanding Tool Approval" - System overview
- "Creating Approval Rules" - Rule configuration guide
- "Approval Best Practices" - Security recommendations

**In-App Help:**
- Tooltips in approval overlay explaining controls
- Help text in settings for rule configuration
- Examples of common approval rules

**Video Tutorials:**
- "Your First Approval" - Walkthrough
- "Streamlining with Auto-Approval" - Rule creation
- "Staying Safe with AI" - Security best practices

---

## Risk & Mitigation

### Risk 1: Approval Fatigue
**Impact:** High - Users annoyed by excessive prompts  
**Probability:** Medium  
**User Impact:** Reduced productivity, frustration

**Mitigation:**
- Smart default auto-approval rules for safe operations
- Easy rule creation from approval dialogs
- Clear documentation on optimal rule configuration
- Template rule sets for common workflows

---

### Risk 2: Users Auto-Approving Everything
**Impact:** High - Defeats purpose of security system  
**Probability:** Low  
**User Impact:** Security vulnerabilities, accidental damage

**Mitigation:**
- Warnings for overly broad rules
- Educational content about security risks
- Audit log visibility showing rule usage
- Require explicit confirmation for dangerous tool auto-approval
- No "disable all approvals" option

---

### Risk 3: Complex Rule Configuration
**Impact:** Medium - Users struggle to create effective rules  
**Probability:** Medium  
**User Impact:** Either too restrictive or too permissive rules

**Mitigation:**
- Simple default rules that work for most users
- Template rule sets for common scenarios
- Clear documentation with examples
- In-app help for rule creation
- Preview showing what rule will match

---

### Risk 4: Poor Denial Experience
**Impact:** Medium - Agent confused by denials  
**Probability:** Low  
**User Impact:** Conversation breakdown, task failure

**Mitigation:**
- Clear denial feedback to agent
- Agent trained to offer alternatives after denial
- User can explain denial reason in chat
- Denial doesn't end conversation, just blocks operation

---

### Risk 5: Compliance Gaps
**Impact:** Low - Audit logs insufficient for compliance  
**Probability:** Low  
**User Impact:** Cannot meet regulatory requirements

**Mitigation:**
- Complete logging of all operations
- Export functionality for external audit systems
- Configurable retention periods
- Timestamp and parameter logging
- Future: Encryption option for sensitive logs

---

## Dependencies & Integration Points

### Feature Dependencies

**Agent Loop:**
- Approval system integrates with tool execution pipeline
- Approval requests pause agent loop until decision
- Denial feedback returns to agent as error

**Tool System:**
- Tools declare if they require approval
- Tool metadata includes preview generation
- Tool schemas used for parameter display

**Settings System:**
- Auto-approval rules stored in user configuration
- Settings UI for rule management
- Persistent storage across sessions

**Event System:**
- Approval events communicated to UI
- Real-time updates during approval flow
- Audit events logged for all decisions

---

### User-Facing Integrations

**TUI Display:**
- Approval overlay rendering
- Diff viewer for file changes
- Command preview formatting

**Keyboard Controls:**
- Ctrl+A, Ctrl+R for approve/deny
- Tab navigation in overlay
- ESC to close and deny

**Settings Interface:**
- Auto-approval rule management
- Timeout configuration
- Audit log access

---

## Constraints & Trade-offs

### Product Constraints

**Security vs. Convenience:**
- **Trade-off:** Full security vs. streamlined workflow
- **Decision:** Secure by default, convenience through configuration
- **Rationale:** Trust and safety are non-negotiable; users can optimize

**Automatic vs. Explicit:**
- **Trade-off:** AI-learned rules vs. explicit configuration
- **Decision:** Explicit configuration only
- **Rationale:** Predictability and transparency over convenience

**Flexibility vs. Simplicity:**
- **Trade-off:** Powerful rule patterns vs. simple configuration
- **Decision:** Simple patterns with expansion path
- **Rationale:** Most users need simple rules; power users can learn advanced

---

### Design Constraints

**Screen Space:**
- **Constraint:** Approval overlay must fit in terminal
- **Trade-off:** Information density vs. readability
- **Decision:** Scrollable viewport with 15-line limit
- **Rationale:** Balance detail with terminal constraints

**Keyboard-Only:**
- **Constraint:** Must work without mouse
- **Trade-off:** UI complexity vs. accessibility
- **Decision:** All actions keyboard-accessible
- **Rationale:** Terminal users expect keyboard workflow

**Preview Generation:**
- **Constraint:** Previews must generate quickly
- **Trade-off:** Accuracy vs. speed
- **Decision:** Fast approximation over perfect rendering
- **Rationale:** Approval speed more important than perfect preview

---

## Competitive Analysis

### GitHub Copilot
**Approach:** No approval system - all suggestions applied manually  
**Strengths:** Simple, no interruptions  
**Weaknesses:** No autonomous execution, no file operations  
**Differentiation:** We enable powerful automation with safety

### Cursor
**Approach:** Command palette for dangerous operations  
**Strengths:** Clear command preview  
**Weaknesses:** Binary yes/no, no rule system  
**Differentiation:** Configurable auto-approval for trusted operations

### Aider
**Approach:** All file changes require git commit review  
**Strengths:** Uses familiar git workflow  
**Weaknesses:** No real-time approval, command execution unrestricted  
**Differentiation:** Real-time approval with comprehensive coverage

### Windsurf
**Approach:** Approval for all tool calls  
**Strengths:** Maximum control  
**Weaknesses:** High friction, approval fatigue  
**Differentiation:** Smart defaults and auto-approval reduce friction

---

## Go-to-Market Considerations

### Positioning

**Primary Message:**  
"Use AI automation confidently—Forge keeps you in control with transparent approval for every dangerous operation. Start safe, streamline with auto-approval rules as you build trust."

**Key Differentiators:**
- Security-first with safe defaults
- Transparent previews before execution
- Flexible auto-approval for trusted operations
- Complete audit trail for compliance

---

### Target Segments

**Early Adopters:**
- Security-conscious developers
- Organizations with compliance requirements
- Teams managing production systems

**Value Propositions by Segment:**
- **Enterprise:** Compliance-ready audit logs, organizational control
- **Indie Developers:** Safety without sacrificing automation
- **Teams:** Educational tool for safe AI adoption

---

### Documentation Needs

**Essential Documentation:**
1. "Understanding Tool Approval" - How the system works
2. "Creating Approval Rules" - Configuration guide
3. "Approval Best Practices" - Security recommendations
4. "Audit Trail Guide" - Compliance usage

**FAQ Topics:**
- "Why do I need to approve file writes?"
- "How do I auto-approve safe operations?"
- "What happens when I deny an operation?"
- "Can I disable approval for certain tools?"

---

### Support Considerations

**Common Support Requests:**
1. Creating auto-approval rules
2. Understanding why operation requires approval
3. Configuring approval timeout
4. Accessing audit logs

**Support Resources:**
- Approval overlay help text
- Settings panel documentation
- Template rule sets
- Security best practices guide

---

## Evolution & Roadmap

### Version History

**v1.0 (Current):**
- Visual approval overlay with diff preview
- Keyboard-driven approval workflow
- Basic auto-approval rule system
- Audit logging

---

### Future Enhancements

#### Phase 2: Enhanced Automation
- Batch approval for similar operations
- Conditional rules ("approve if file < 1000 lines")
- Time-based rules ("auto-approve for next hour")
- Workspace-specific rule sets

**User Value:** More flexible automation without sacrificing control

---

#### Phase 3: Team & Compliance
- Approval delegation (junior dev → senior approval)
- Team-shared rule templates
- Enhanced audit log encryption
- Approval analytics dashboard

**User Value:** Enterprise-ready compliance and team workflows

---

#### Phase 4: AI-Assisted Review
- AI highlights potentially problematic changes
- Change impact analysis (files affected)
- Rollback support for approved operations
- Smart rule suggestions based on usage

**User Value:** Intelligent assistance in making approval decisions

---

### Open Questions

**Question 1: Should we support approval delegation?**
- **Pro:** Teams need senior review workflow
- **Con:** Requires multi-user architecture
- **Current Direction:** Phase 3 feature

**Question 2: Should audit logs be encrypted by default?**
- **Pro:** Protects sensitive information
- **Con:** Makes debugging harder
- **Current Direction:** Optional encryption in Phase 2

**Question 3: Should we AI-learn approval patterns?**
- **Pro:** Reduce configuration burden
- **Con:** Unpredictable, security risk
- **Current Direction:** Explicit rules only

**Question 4: Maximum rule count limit?**
- **Pro:** Prevents overly permissive configurations
- **Con:** Power users may need many rules
- **Current Direction:** Warn at 20+, no hard limit

---

## Technical References

- **Architecture Documentation:** `docs/architecture/tool-approval.md`
- **Implementation Details:** See ADR-0010 (Tool Approval Mechanism)
- **Auto-Approval System:** See ADR-0017 (Auto-Approval and Settings)
- **Related Features:** Tool System PRD, Settings System PRD

---

## Changelog

### 2024-12-XX
- Transformed to product-focused PRD format
- Removed technical implementation details
- Enhanced user experience sections
- Added competitive analysis
- Expanded go-to-market considerations

### 2024-12 (Original)
- Initial PRD with technical architecture
- Component structure and data models
- Code examples and flow diagrams
