# Product Requirements: Agent Loop Architecture

**Feature:** Core Agent Reasoning and Execution Loop  
**Version:** 1.0  
**Status:** Implemented  
**Owner:** Core Team  
**Last Updated:** December 2024

---

## Overview

The Agent Loop Architecture is the heart of Forge's AI capabilities, implementing an iterative reasoning-action cycle where the agent analyzes user requests, plans solutions, executes tools, and learns from results. This loop enables the agent to break down complex tasks, use tools effectively, and adapt its approach based on outcomes.

---

## Problem Statement

AI coding assistants need to handle complex, multi-step tasks that require:

1. **Task Decomposition:** Breaking down high-level requests into actionable steps
2. **Tool Orchestration:** Using multiple tools in sequence to accomplish goals
3. **Adaptive Planning:** Adjusting approach based on intermediate results
4. **Error Recovery:** Handling failures and trying alternative approaches
5. **Context Maintenance:** Tracking progress across multiple iterations
6. **Safety Limits:** Preventing infinite loops or runaway execution

Without a robust execution loop, AI assistants either:
- Execute only simple, single-step tasks
- Require excessive user handholding for complex workflows
- Get stuck when encountering errors
- Waste tokens on inefficient approaches

---

## Goals

### Primary Goals

1. **Autonomous Execution:** Enable agent to complete multi-step tasks independently
2. **Intelligent Planning:** Use chain-of-thought reasoning to plan before acting
3. **Efficient Tool Use:** Select and execute tools appropriately
4. **Graceful Failure:** Handle errors and adapt strategies
5. **Resource Management:** Stay within token and iteration limits
6. **User Control:** Allow users to interrupt or guide execution

### Non-Goals

1. **General AI:** This is NOT trying to create AGI
2. **Unlimited Autonomy:** Agent has clear boundaries and safety limits
3. **Self-Modification:** Agent does NOT modify its own code or prompts
4. **Multi-Agent Coordination:** Single agent only (no agent-to-agent communication)

---

## User Personas

### Primary: Task-Oriented Developer
- **Background:** Developer with specific coding tasks to complete
- **Workflow:** Describes goal, lets agent work, reviews results
- **Pain Points:** Manual multi-step workflows are tedious
- **Goals:** Delegate complex tasks to AI assistant

### Secondary: Exploratory Developer
- **Background:** Developer experimenting with ideas
- **Workflow:** Iterative refinement with agent
- **Pain Points:** Agent gives up too easily or tries wrong approaches
- **Goals:** Agent that persists and adapts

### Tertiary: Quality-Focused Developer
- **Background:** Developer who values correctness over speed
- **Workflow:** Reviews each step carefully
- **Pain Points:** Agent makes assumptions or skips validation
- **Goals:** Transparent, verifiable agent reasoning

---

## Requirements

### Functional Requirements

#### FR1: Agent Loop Cycle
- **R1.1:** Receive and analyze user message
- **R1.2:** Generate chain-of-thought reasoning
- **R1.3:** Select appropriate tool or response action
- **R1.4:** Execute tool call (if needed)
- **R1.5:** Receive and process tool result
- **R1.6:** Update context with new information
- **R1.7:** Decide whether to continue or complete
- **R1.8:** Iterate or present final result

#### FR2: Chain-of-Thought Reasoning
- **R2.1:** Think through problem before acting
- **R2.2:** Break down complex tasks into sub-tasks
- **R2.3:** Identify required information or tools
- **R2.4:** Consider multiple approaches
- **R2.5:** Explain reasoning in thinking tags
- **R2.6:** Reason about tool selection and parameters

#### FR3: Tool Call Detection & Parsing
- **R3.1:** Detect tool calls in LLM response (XML format)
- **R3.2:** Parse tool name and parameters
- **R3.3:** Validate tool call structure
- **R3.4:** Handle malformed tool calls gracefully
- **R3.5:** Support early detection (before full response)
- **R3.6:** Extract multiple tool calls if present (error case)

#### FR4: Tool Execution
- **R4.1:** Request approval for dangerous tools
- **R4.2:** Execute approved tools
- **R4.3:** Capture tool results (stdout, stderr, return value)
- **R4.4:** Handle tool execution errors
- **R4.5:** Time out long-running tools
- **R4.6:** Format results for agent consumption

#### FR5: Iteration Management
- **R5.1:** Track iteration count per agent loop
- **R5.2:** Enforce maximum iteration limit (default: 25)
- **R5.3:** Warn agent when approaching limit
- **R5.4:** Gracefully terminate at limit
- **R5.5:** Allow user to extend limit if needed
- **R5.6:** Reset count for new user messages

#### FR6: Loop-Breaking Tools
- **R6.1:** Detect loop-breaking tool calls
  - `task_completion` - Finish task
  - `ask_question` - Ask user for input
  - `converse` - Casual conversation
- **R6.2:** Stop iteration immediately on loop-breaker
- **R6.3:** Present result to user
- **R6.4:** Wait for user response before continuing

#### FR7: Error Handling
- **R7.1:** Catch tool execution errors
- **R7.2:** Present errors to agent for recovery
- **R7.3:** Allow agent to try alternative approaches
- **R7.4:** Track repeated errors (avoid infinite retry)
- **R7.5:** Escalate to user if unrecoverable
- **R7.6:** Maintain conversation context on error

#### FR8: Context Management
- **R8.1:** Build conversation context for each iteration
- **R8.2:** Include system prompt with instructions
- **R8.3:** Add conversation history (pruned if needed)
- **R8.4:** Include tool definitions and schemas
- **R8.5:** Track token usage per request
- **R8.6:** Prune old messages when approaching limits

#### FR9: User Interruption
- **R9.1:** Allow user to stop execution mid-loop
- **R9.2:** Handle Ctrl+C gracefully
- **R9.3:** Support `/stop` slash command
- **R9.4:** Save conversation state on interrupt
- **R9.5:** Allow resumption after interrupt
- **R9.6:** Clean up resources on forced stop

#### FR10: Progress Visibility
- **R10.1:** Show agent thinking in real-time
- **R10.2:** Display tool execution status
- **R10.3:** Stream tool results as they arrive
- **R10.4:** Update iteration count
- **R10.5:** Show "working" indicators during pauses
- **R10.6:** Display errors and warnings prominently

### Non-Functional Requirements

#### NFR1: Performance
- **N1.1:** Iteration cycle overhead under 100ms
- **N1.2:** Tool call detection within 50ms
- **N1.3:** Context building under 200ms
- **N1.4:** Minimal latency between iterations
- **N1.5:** Efficient memory usage during long loops

#### NFR2: Reliability
- **N2.1:** Never lose conversation context
- **N2.2:** Recover from LLM API errors
- **N2.3:** Handle unexpected tool failures
- **N2.4:** Graceful degradation on resource limits
- **N2.5:** 99.9% loop completion rate (non-error cases)

#### NFR3: Safety
- **N3.1:** Hard limit on max iterations (prevent runaway)
- **N3.2:** Timeout on individual tool executions
- **N3.3:** Sandbox tool execution to workspace
- **N3.4:** Validate all tool parameters
- **N3.5:** Audit trail of all actions

#### NFR4: Transparency
- **N4.1:** All reasoning visible to user
- **N4.2:** Tool calls and results logged
- **N4.3:** Clear indication of agent state
- **N4.4:** Decision points explained
- **N4.5:** Easy to debug unexpected behavior

---

## User Experience

### Core Workflows

#### Workflow 1: Simple Single-Iteration Task
1. User: "What files are in the src directory?"
2. Agent thinks: "I need to list files in src/"
3. Agent calls: `list_files` with path="src"
4. Tool executes, returns file list
5. Agent thinks: "Got the results, task complete"
6. Agent calls: `task_completion` with formatted list
7. User sees result

**Success Criteria:** Task completed in 1 iteration, under 5 seconds

#### Workflow 2: Multi-Step Refactoring Task
1. User: "Refactor the authentication code to use middleware"
2. Agent thinks: "Need to understand current auth implementation"
3. Iteration 1: `read_file` auth.go
4. Agent analyzes code, plans refactoring
5. Iteration 2: `read_file` middleware.go (to see patterns)
6. Agent designs new middleware structure
7. Iteration 3: `write_file` auth_middleware.go (new file)
8. Iteration 4: `apply_diff` to modify existing files
9. Iteration 5: `apply_diff` to update imports
10. Agent calls: `task_completion` with summary
11. User reviews changes

**Success Criteria:** Multi-step task completed autonomously in 5-10 iterations

#### Workflow 3: Error Recovery
1. User: "Run the tests"
2. Agent: `execute_command` "npm test"
3. Tool result: Exit code 1, tests failed
4. Agent thinks: "Tests failed, need to investigate"
5. Agent: `read_file` on failing test
6. Agent identifies issue in code
7. Agent: `apply_diff` to fix code
8. Agent: `execute_command` "npm test" again
9. Tool result: Exit code 0, tests pass
10. Agent: `task_completion` "Fixed issue, tests passing"

**Success Criteria:** Agent recovers from failure and completes task

#### Workflow 4: Information Gathering
1. User: "Should we upgrade to React 19?"
2. Agent thinks: "Need to check current version and dependencies"
3. Iteration 1: `read_file` package.json
4. Agent sees React 18.2.0
5. Agent calls: `ask_question` "Would you like me to check for breaking changes in React 19?"
6. User: "Yes"
7. Agent resumes, researches breaking changes
8. Agent: `task_completion` with upgrade recommendation

**Success Criteria:** Agent gathers needed info before making recommendations

#### Workflow 5: User Interruption
1. User: "Refactor all components to TypeScript"
2. Agent starts working (10+ files to modify)
3. After 3 files converted, user presses Ctrl+C
4. Agent stops immediately
5. Conversation state preserved
6. User: "Actually, just do the header component"
7. Agent adjusts and completes focused task

**Success Criteria:** User can interrupt and redirect agent anytime

---

## Technical Architecture

### Component Structure

```
Agent Loop System
├── Loop Controller
│   ├── Iteration Manager
│   ├── State Tracker
│   └── Termination Detector
├── Reasoning Engine
│   ├── LLM Provider Interface
│   ├── Prompt Builder
│   └── Response Parser
├── Tool System
│   ├── Tool Registry
│   ├── Tool Executor
│   └── Result Formatter
├── Context Manager
│   ├── Message History
│   ├── Token Counter
│   └── Memory Pruner
└── Event System
    ├── Event Emitter
    ├── Progress Reporter
    └── Error Handler
```

### Agent Loop Flow

```
User Message Received
    ↓
Initialize Loop Context
    ↓
┌─────────────────────────────────────┐
│ AGENT LOOP (max 25 iterations)      │
│                                     │
│  1. Build Conversation Context      │
│     - System prompt                 │
│     - Message history               │
│     - Tool definitions              │
│                                     │
│  2. Call LLM (streaming)            │
│     - Send context                  │
│     - Receive response              │
│     - Emit thinking events          │
│                                     │
│  3. Parse Response                  │
│     - Extract thinking              │
│     - Detect tool call              │
│     - Validate structure            │
│                                     │
│  4. Check Loop Breaking             │
│     ├─ task_completion → DONE       │
│     ├─ ask_question → DONE          │
│     └─ converse → DONE              │
│                                     │
│  5. Execute Tool (if needed)        │
│     - Request approval              │
│     - Execute tool                  │
│     - Capture result                │
│     - Add to context                │
│                                     │
│  6. Check Termination               │
│     - Max iterations reached?       │
│     - User interrupted?             │
│     - Error threshold exceeded?     │
│                                     │
│  7. Continue or Exit                │
│     └─ Loop back to step 1          │
└─────────────────────────────────────┘
    ↓
Present Final Result
```

### Data Model

```go
type AgentLoop struct {
    Context        *ConversationContext
    IterationCount int
    MaxIterations  int
    State          LoopState
    ErrorCount     int
    StartTime      time.Time
}

type LoopState int
const (
    StateInitializing LoopState = iota
    StateThinking
    StateExecutingTool
    StateWaitingApproval
    StateCompleted
    StateError
    StateInterrupted
)

type ToolCall struct {
    ServerName string
    ToolName   string
    Arguments  map[string]interface{}
}

type ToolResult struct {
    ToolName   string
    Success    bool
    Result     string
    Error      error
    Duration   time.Duration
}
```

---

## Design Decisions

### Why 25 Max Iterations?
**Rationale:**
- **Safety:** Prevents infinite loops from consuming resources
- **Cost control:** Limits API usage for runaway sessions
- **Sufficient:** 95% of tasks complete in under 10 iterations
- **Configurable:** Users can adjust if needed

**Alternatives considered:**
- Unlimited: Too risky, could drain API credits
- Token-based limit: Harder to predict behavior
- Time-based: Doesn't account for task complexity

### Why Chain-of-Thought Required?
**Rationale:**
- **Better reasoning:** LLMs perform better with explicit thinking
- **Transparency:** Users see agent's thought process
- **Debugging:** Easier to understand unexpected behavior
- **Quality:** Reduces impulsive, poorly-planned actions

**Evidence:** Research shows CoT improves task completion by 15-30%

### Why XML for Tool Calls?
**Rationale:**
- **Structure:** Enforces clear parameter structure
- **Validation:** Easy to validate before execution
- **LLM-friendly:** Models trained on XML formatting
- **Extensible:** Easy to add new parameters/tools
- **Unambiguous:** Clear start/end markers

**Alternatives considered:**
- JSON: LLMs struggle with nested structures
- Natural language: Too ambiguous, hard to parse
- Custom DSL: Learning curve, less familiar

### Why Loop-Breaking Tools?
**Rationale:**
- **Clear intent:** Explicit signal that task is complete
- **User control:** Allows agent to ask questions
- **Conversation flow:** Maintains natural dialog
- **Efficiency:** Avoids unnecessary iterations

**Without loop-breakers:** Agent might keep iterating needlessly

---

## Agent Loop States

### State Transitions

```
Initializing → Thinking
    ↓
Thinking → ExecutingTool (if tool call detected)
    ↓
ExecutingTool → WaitingApproval (if approval needed)
    ↓
WaitingApproval → ExecutingTool (if approved)
    ↓
ExecutingTool → Thinking (result received)
    ↓
Thinking → Completed (loop-breaking tool)
    ↓
Thinking → Error (unrecoverable error)
    ↓
Thinking → Interrupted (user stopped)
```

### State Behaviors

| State | Behavior | User Visible |
|-------|----------|--------------|
| Initializing | Setting up context | "Starting..." |
| Thinking | LLM generating response | Agent message streaming |
| ExecutingTool | Running tool | "Agent is using [tool]" |
| WaitingApproval | Awaiting user approval | Approval overlay |
| Completed | Task done | Final result |
| Error | Unrecoverable error | Error message |
| Interrupted | User stopped | "Stopped by user" |

---

## Success Metrics

### Effectiveness Metrics
- **Task completion rate:** >90% of valid requests completed successfully
- **Average iterations:** <7 iterations per task
- **Error recovery rate:** >80% of errors recovered without user intervention
- **Token efficiency:** <10% wasted tokens on dead ends

### Performance Metrics
- **Time to first tool:** <5 seconds from user message
- **Iteration latency:** <2 seconds average (LLM + tool execution)
- **Loop overhead:** <5% of total execution time
- **Context building:** <200ms per iteration

### Quality Metrics
- **Reasoning quality:** >85% of thinking sections are logical and relevant
- **Tool selection accuracy:** >95% of tool calls are appropriate
- **Parameter correctness:** >98% of tool parameters are valid
- **Plan adherence:** >80% of multi-step plans executed as intended

### User Experience Metrics
- **Interruption handling:** 100% of interrupts handled gracefully
- **Transparency:** >90% of users understand agent reasoning
- **Trust:** >85% of users confident in agent decisions
- **Surprise rate:** <10% of actions are unexpected to users

---

## Dependencies

### External Dependencies
- LLM provider API (OpenAI, Anthropic, etc.)
- Tool execution environment
- Token counting library

### Internal Dependencies
- Tool system (registry, execution)
- Memory system (history, pruning)
- Event system (progress updates)
- Settings system (max iterations, timeouts)
- TUI (for progress display)

### Platform Requirements
- Network access (for LLM API)
- Sufficient memory (conversation context)
- Tool execution permissions

---

## Risks & Mitigations

### Risk 1: Infinite Loops
**Impact:** Critical  
**Probability:** Low  
**Mitigation:**
- Hard limit on max iterations (25)
- Detect repeated identical tool calls
- Warn agent when approaching limit
- User interrupt capability
- Automatic termination at limit

### Risk 2: Context Overflow
**Impact:** High  
**Probability:** Medium  
**Mitigation:**
- Token counting per request
- Automatic message pruning
- Warn at 80% context usage
- Suggest starting fresh session
- Efficient summarization of old context

### Risk 3: Poor Tool Selection
**Impact:** Medium  
**Probability:** Medium  
**Mitigation:**
- Clear tool documentation in prompt
- Examples of tool usage
- Validation before execution
- Error feedback to agent
- User approval for dangerous tools

### Risk 4: LLM API Failures
**Impact:** High  
**Probability:** Low  
**Mitigation:**
- Retry logic with exponential backoff
- Fallback to alternative provider (if configured)
- Graceful error messages
- Preserve conversation state
- Allow manual retry

### Risk 5: Runaway Costs
**Impact:** High  
**Probability:** Low  
**Mitigation:**
- Iteration limits
- Token budgets per session
- Cost warnings in context overlay
- Auto-stop at spending thresholds
- Local model option (no API cost)

---

## Future Enhancements

### Phase 2 Ideas
- **Parallel Tool Execution:** Execute independent tools concurrently
- **Sub-Agent Spawning:** Delegate sub-tasks to specialized agents
- **Plan Caching:** Reuse plans for similar tasks
- **Learning from Feedback:** Improve planning based on success/failure
- **Dynamic Iteration Limits:** Adjust max iterations based on task complexity

### Phase 3 Ideas
- **Multi-Agent Collaboration:** Multiple agents working together
- **Hierarchical Planning:** High-level planning with detailed sub-plans
- **Reinforcement Learning:** Optimize tool selection over time
- **Verification Loop:** Automatic validation of agent outputs
- **Rollback/Undo:** Revert to previous agent state if needed

---

## Open Questions

1. **Should we support parallel tool execution?**
   - Pro: Faster completion for independent tasks
   - Con: More complex error handling, approval flow
   - Decision: Phase 2 feature if performance gains justify complexity

2. **Should agent learn from past successes/failures?**
   - Pro: Improves over time
   - Con: Requires persistence, complexity
   - Decision: Phase 3 research project

3. **Should we support agent-defined tools?**
   - Use case: Agent creates custom functions
   - Risk: Security, validation complexity
   - Decision: No - too risky for current version

4. **Should iteration limit be dynamic?**
   - Pro: Complex tasks get more iterations
   - Con: Harder to predict resource usage
   - Decision: Static limit for now, revisit in Phase 2

---

## Related Documentation

- [ADR-0003: Agent Core Loop Design](../adr/0003-agent-core-loop-design.md)
- [ADR-0021: Early Tool Call Detection](../adr/0021-early-tool-call-detection.md)
- [Architecture: Agent Loop](../architecture/agent-loop.md)
- [Architecture: Event System](../architecture/events.md)
- [Built-in Tools Reference](../reference/built-in-tools.md)

---

## Changelog

| Date | Version | Changes |
|------|---------|---------|
| 2024-12 | 1.0 | Initial PRD creation |
