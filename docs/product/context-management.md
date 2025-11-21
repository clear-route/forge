# Product Requirements: Context Management

**Feature:** Context Information and Tracking  
**Version:** 1.0  
**Status:** Implemented  
**Owner:** Core Team  
**Last Updated:** December 2024

---

## Overview

Context Management provides users with comprehensive visibility into the agent's current state, token usage, conversation history, and workspace information. It helps users understand what the agent knows, how much context is being used, and make informed decisions about conversation flow and resource usage.

---

## Problem Statement

Users interacting with AI agents face context-related challenges:

1. **Black Box Problem:** Users don't know what context the agent has access to
2. **Token Blindness:** No visibility into token usage until hitting limits
3. **Cost Uncertainty:** Cloud LLM usage costs are hidden until bill arrives
4. **Memory Mystery:** Unclear what conversation history is retained
5. **Context Overflow:** Users don't know when approaching context limits
6. **Debugging Difficulty:** Hard to troubleshoot why agent behaves certain way

Without context visibility, users can't effectively manage conversations, optimize costs, or understand agent behavior.

---

## Goals

### Primary Goals

1. **Transparency:** Make all context information visible and understandable
2. **Token Awareness:** Show real-time token usage and limits
3. **Cost Management:** Help users optimize LLM API costs
4. **Conversation Insight:** Display what's in conversation history
5. **Workspace Clarity:** Show what files/directories agent can access
6. **Performance Monitoring:** Track session statistics and metrics

### Non-Goals

1. **Context Modification:** This is NOT for editing context directly
2. **Context Injection:** Does NOT allow injecting arbitrary context
3. **Context Filtering:** Does NOT provide tools to filter/remove context items
4. **Multi-Session Context:** Does NOT track context across multiple sessions (yet)

---

## User Personas

### Primary: Cost-Conscious Developer
- **Background:** Using paid LLM APIs, monitors costs carefully
- **Workflow:** Long coding sessions with agent
- **Pain Points:** Surprised by large API bills
- **Goals:** Understand and control token usage

### Secondary: Context-Aware Power User
- **Background:** Experienced with LLMs, understands context windows
- **Workflow:** Carefully manages conversation to stay within limits
- **Pain Points:** Needs to know when approaching limit
- **Goals:** Maximize context efficiency

### Tertiary: Debugging Developer
- **Background:** Troubleshooting unexpected agent behavior
- **Workflow:** Investigating why agent missed information or hallucinated
- **Pain Points:** Can't see what context agent actually has
- **Goals:** Understand what agent knows to fix issues

---

## Requirements

### Functional Requirements

#### FR1: Context Overlay UI
- **R1.1:** Accessible via `/context` slash command
- **R1.2:** Accessible via keyboard shortcut (Ctrl+I)
- **R1.3:** Display in modal overlay
- **R1.4:** Organized sections with clear labels
- **R1.5:** Scrollable content for large context
- **R1.6:** Close with Esc or dedicated key

#### FR2: Workspace Information
- **R2.1:** Show current workspace path
- **R2.2:** Display workspace statistics (files, size, language breakdown)
- **R2.3:** Show workspace permissions (read/write access)
- **R2.4:** List workspace restrictions (sandboxing info)
- **R2.5:** Indicate if workspace is git repository

#### FR3: Conversation History
- **R3.1:** Show total message count (user + agent + tool)
- **R3.2:** Display message breakdown by role
- **R3.3:** Show conversation age (time since first message)
- **R3.4:** List recent topics/themes (optional, AI-generated)
- **R3.5:** Indicate if history has been pruned

#### FR4: Token Usage Tracking
- **R4.1:** Display current request token breakdown:
  - System prompt tokens
  - User message tokens
  - Assistant response tokens
  - Tool call tokens
  - Tool result tokens
- **R4.2:** Show cumulative session totals
- **R4.3:** Display context window limit (model-specific)
- **R4.4:** Calculate remaining token budget
- **R4.5:** Show percentage of context used
- **R4.6:** Warn when approaching limit (>80%)

#### FR5: Memory State
- **R5.1:** Show current memory size (in messages)
- **R5.2:** Display memory strategy (e.g., "rolling window")
- **R5.3:** Indicate pruning status (if history trimmed)
- **R5.4:** Show oldest message still in memory
- **R5.5:** Display memory configuration (max messages/tokens)

#### FR6: System Prompt Information
- **R6.1:** Show system prompt length (tokens)
- **R6.2:** Display custom instructions if set
- **R6.3:** List active tools in prompt
- **R6.4:** Show prompt version/timestamp
- **R6.5:** Indicate if using default vs custom prompt

#### FR7: Session Statistics
- **R7.1:** Session duration (time since start)
- **R7.2:** Total tool calls executed
- **R7.3:** Total agent iterations
- **R7.4:** Success/failure counts
- **R7.5:** Average response time
- **R7.6:** Provider information (model, endpoint)

#### FR8: Cost Estimation
- **R8.1:** Estimate current session cost (if provider pricing known)
- **R8.2:** Show cost breakdown (input vs output tokens)
- **R8.3:** Display cost per message (average)
- **R8.4:** Provide cost optimization tips
- **R8.5:** Compare costs across different models

#### FR9: Context Health Indicators
- **R9.1:** Green/yellow/red status for token usage
- **R9.2:** Warning for approaching context limit
- **R9.3:** Alert for excessive tool result sizes
- **R9.4:** Suggestion to clear history or summarize
- **R9.5:** Performance impact indicators

### Non-Functional Requirements

#### NFR1: Performance
- **N1.1:** Context overlay opens within 100ms
- **N1.2:** Token counting completes within 50ms
- **N1.3:** Statistics calculation under 100ms
- **N1.4:** No impact on agent loop performance
- **N1.5:** Efficient caching of computed values

#### NFR2: Accuracy
- **N2.1:** Token counts within 5% of actual usage
- **N2.2:** Cost estimates within 10% of actual (if pricing stable)
- **N2.3:** Message counts exactly accurate
- **N2.4:** Timestamps precise to the second
- **N2.5:** Consistent calculations across refreshes

#### NFR3: Usability
- **N3.1:** Information presented in logical groupings
- **N3.2:** No jargon; clear explanations
- **N3.3:** Visual hierarchy for important info
- **N3.4:** Color coding for status indicators
- **N3.5:** Responsive to different terminal sizes

#### NFR4: Reliability
- **N4.1:** Never crash when displaying context
- **N4.2:** Graceful handling of missing data
- **N4.3:** Fallback values for unavailable metrics
- **N4.4:** Consistent behavior across providers
- **N4.5:** Safe handling of large contexts

---

## User Experience

### Core Workflows

#### Workflow 1: Checking Token Usage Mid-Session
1. User in middle of long conversation
2. Wants to know how much context used
3. Presses Ctrl+I or types `/context`
4. Context overlay opens
5. User sees "Token Usage: 3,247 / 8,192 (40%)"
6. Status is green (safe)
7. User closes overlay
8. Continues conversation confidently

**Success Criteria:** User knows token budget in under 3 seconds

#### Workflow 2: Investigating Why Agent Missed Information
1. Agent didn't reference earlier message
2. User suspects context was pruned
3. Opens context overlay
4. Sees "Memory: Pruned - oldest message from 15 min ago"
5. Realizes early context was lost
6. User provides information again
7. Agent responds correctly

**Success Criteria:** User understands why agent forgot context

#### Workflow 3: Optimizing for Cost
1. User worried about API costs
2. Opens context overlay
3. Sees "Estimated session cost: $0.45"
4. Notices tool results are large (2K tokens each)
5. Checks "Cost Optimization Tips"
6. Learns to use more focused queries
7. Adjusts conversation style

**Success Criteria:** User has actionable cost insights

#### Workflow 4: Understanding Workspace Scope
1. New user unsure what agent can access
2. Opens context overlay
3. Sees "Workspace: /home/user/project (123 files, Go/TypeScript)"
4. Notices "Restrictions: Read/write in workspace only"
5. Understands sandbox boundaries
6. Asks agent to work on project files

**Success Criteria:** User knows workspace boundaries

#### Workflow 5: Monitoring Long-Running Session
1. User in 2-hour coding session
2. Periodically checks context
3. Sees token usage climbing: 60% → 75% → 85%
4. Gets yellow warning at 80%
5. Decides to start fresh session
6. Saves important info
7. Restarts with clean context

**Success Criteria:** User prevents context overflow

---

## Technical Architecture

### Component Structure

```
Context Management System
├── Context Collector
│   ├── Workspace Analyzer
│   ├── Memory Inspector
│   ├── Token Counter
│   └── Session Tracker
├── Context Overlay (TUI)
│   ├── Info Renderer
│   ├── Status Indicators
│   ├── Progress Bars
│   └── Section Organizer
├── Token Accounting
│   ├── Tokenizer Integration
│   ├── Token Cache
│   ├── Usage Tracker
│   └── Cost Calculator
└── Statistics Engine
    ├── Metric Aggregator
    ├── Trend Analyzer
    └── Health Checker
```

### Data Model

```go
type ContextInfo struct {
    Workspace    WorkspaceInfo
    Conversation ConversationInfo
    TokenUsage   TokenUsageInfo
    Memory       MemoryInfo
    System       SystemInfo
    Session      SessionInfo
}

type WorkspaceInfo struct {
    Path         string
    FileCount    int
    TotalSize    int64
    Languages    map[string]int
    IsGitRepo    bool
    Permissions  []string
}

type ConversationInfo struct {
    MessageCount    int
    UserMessages    int
    AssistantMsgs   int
    ToolCalls       int
    StartTime       time.Time
    LastActivity    time.Time
}

type TokenUsageInfo struct {
    Current      TokenBreakdown
    Session      TokenBreakdown
    Limit        int
    Remaining    int
    PercentUsed  float64
    Status       HealthStatus
}

type TokenBreakdown struct {
    System    int
    User      int
    Assistant int
    Tools     int
    Total     int
}

type MemoryInfo struct {
    MessageCount  int
    IsPruned      bool
    Strategy      string
    OldestMessage time.Time
    Config        MemoryConfig
}

type SessionInfo struct {
    Duration      time.Duration
    Iterations    int
    ToolCalls     int
    Provider      string
    Model         string
    EstimatedCost float64
}
```

### Context Collection Flow

```
Slash Command: /context
    ↓
Context Collector: Gather info from multiple sources
    ↓
┌──────────────────────────────────────┐
│ Parallel Collection:                 │
│ - Workspace analyzer                 │
│ - Memory system query                │
│ - Token counter                      │
│ - Session tracker                    │
└──────────────┬───────────────────────┘
               ↓
Aggregate Results
               ↓
Calculate Derived Metrics
  - Percentage used
  - Health status
  - Cost estimates
               ↓
Render Context Overlay
               ↓
Display to User
```

---

## Design Decisions

### Why Real-Time Token Counting?
- **Accuracy:** Shows actual usage, not estimates
- **Transparency:** Users see what LLM sees
- **Optimization:** Enables informed decisions
- **Cost control:** Prevents bill shock
- **Industry practice:** Other tools (OpenAI playground) show tokens

### Why Session Totals vs Per-Message?
- **Context:** Session totals show cumulative usage
- **Trends:** Users want to know total cost/usage
- **Budgeting:** Helps manage daily/weekly limits
- **Both needed:** Show current + cumulative

**Decision:** Display both current request and session totals

### Why Not Allow Context Editing?
- **Complexity:** Editing context can break agent reasoning
- **Safety:** Manual editing could corrupt state
- **Use case unclear:** Users want visibility, not editing
- **Alternative:** Provide "start fresh" option instead

**Decision:** Read-only context view; editing in future if needed

### Why Show Workspace Stats?
- **Transparency:** Users should know what agent can access
- **Trust:** Proves sandboxing works
- **Context:** Helps users frame requests appropriately
- **Debugging:** Explains why agent can't access certain files

---

## Context Information Sections

### 1. Workspace Overview
```
📁 Workspace
  Path: /home/user/my-project
  Files: 247 (Go: 156, TypeScript: 78, Other: 13)
  Size: 4.2 MB
  Git Repo: Yes (main branch)
  Access: Read/Write (sandboxed)
```

### 2. Conversation History
```
💬 Conversation
  Messages: 45 (User: 23, Agent: 18, Tools: 4)
  Duration: 32 minutes
  Started: 2:15 PM
  Last Activity: Just now
```

### 3. Token Usage
```
🔢 Token Usage
  Current Request: 3,247 / 8,192 (40%) ●●●●○○○○○○
  Session Total: 45,891 tokens
  
  Breakdown:
    System Prompt:   487 tokens (6%)
    User Messages:   1,234 tokens (15%)
    Agent Responses: 1,456 tokens (18%)
    Tool Results:    70 tokens (1%)
    
  Status: ✓ Healthy (60% remaining)
```

### 4. Memory State
```
🧠 Memory
  Messages in Memory: 45
  Strategy: Rolling window (last 50 messages)
  Oldest Message: 32 minutes ago
  Pruned: No
```

### 5. Session Statistics
```
📊 Session Stats
  Duration: 32 minutes
  Agent Iterations: 67
  Tools Executed: 23
  Success Rate: 95%
  Avg Response Time: 2.3s
```

### 6. Cost Information
```
💰 Estimated Cost
  Session: ~$0.34
  Breakdown:
    Input tokens:  $0.12 (23K @ $5/1M)
    Output tokens: $0.22 (15K @ $15/1M)
  
  Tips: Tool results account for 40% of tokens.
        Consider more focused file reads.
```

---

## Success Metrics

### Adoption Metrics
- **Usage rate:** >50% of users access context info at least once per session
- **Frequency:** Power users check context 3+ times per long session
- **Discovery:** >70% of users find context feature within first week

### Effectiveness Metrics
- **Token awareness:** >80% of users can estimate token usage after using feature
- **Cost reduction:** 30% decrease in token usage for cost-conscious users
- **Context management:** 50% fewer "context overflow" errors
- **Session planning:** Users start fresh sessions before hitting limits

### Usability Metrics
- **Comprehension:** >90% of users understand displayed information
- **Action rate:** >40% of context checks lead to user action (optimization, fresh start, etc.)
- **Error prevention:** 60% reduction in "agent forgot context" complaints

---

## Dependencies

### External Dependencies
- Tokenizer library (for accurate token counting)
- File system access (for workspace stats)
- Time utilities (for duration tracking)

### Internal Dependencies
- Agent core (for conversation history)
- Memory system (for memory state)
- LLM provider (for token limits, pricing)
- Settings system (for configuration)

### Platform Requirements
- Access to conversation history
- Token counting capability
- File system stat operations

---

## Risks & Mitigations

### Risk 1: Inaccurate Token Counts
**Impact:** High  
**Probability:** Medium  
**Mitigation:**
- Use same tokenizer as LLM provider
- Validate counts against API responses
- Update tokenizer with provider changes
- Show estimates with ±5% disclaimer
- Cache tokenization results

### Risk 2: Performance Impact
**Impact:** Medium  
**Probability:** Low  
**Mitigation:**
- Lazy computation (only when overlay opened)
- Cache calculated values
- Async collection of non-critical stats
- Optimize tokenization (batch operations)
- Set timeout for stat collection

### Risk 3: Information Overload
**Impact:** Medium  
**Probability:** Medium  
**Mitigation:**
- Progressive disclosure (basic vs detailed view)
- Visual hierarchy (important info prominent)
- Color coding for quick scanning
- Collapsible sections
- Focus on actionable insights

### Risk 4: Privacy Concerns
**Impact:** Low  
**Probability:** Low  
**Mitigation:**
- Don't log context info
- No external analytics of usage patterns
- Keep data local
- Clear about what's shown vs stored

---

## Future Enhancements

### Phase 2 Ideas
- **Context History:** Track token usage over time (graphs)
- **Context Export:** Save context snapshot for debugging
- **Context Comparison:** Compare current vs previous sessions
- **Smart Suggestions:** AI-powered context optimization tips
- **Context Alerts:** Proactive warnings before issues

### Phase 3 Ideas
- **Multi-Session Context:** Track across multiple sessions
- **Context Sharing:** Share context snapshots with team
- **Context Templates:** Pre-configured context for tasks
- **Context Replay:** Recreate conversation from context snapshot
- **Advanced Analytics:** Detailed token usage patterns

---

## Open Questions

1. **Should we show message-level token breakdown?**
   - Pro: More granular visibility
   - Con: UI complexity, performance cost
   - Decision: Phase 2 feature (detailed view mode)

2. **Should we estimate costs for local models?**
   - Pro: Consistency across providers
   - Con: Cost is $0 for local models
   - Decision: Show "Local model (no cost)" instead of $0.00

3. **Should we track context across sessions?**
   - Pro: Better long-term insights
   - Con: Persistence complexity
   - Decision: Phase 3 feature

4. **Should context info be exportable?**
   - Use case: Bug reports, debugging
   - Complexity: Format, privacy concerns
   - Decision: Add export in Phase 2

---

## Related Documentation

- [ADR-0020: Context Information Overlay](../adr/0020-context-information-overlay.md)
- [How-to: Use TUI Interface - Context](../how-to/use-tui-interface.md#context-information)
- [Architecture: Memory System](../architecture/memory.md)
- [Agent Loop Architecture](../architecture/agent-loop.md)

---

## Changelog

| Date | Version | Changes |
|------|---------|---------|
| 2024-12 | 1.0 | Initial PRD creation |
