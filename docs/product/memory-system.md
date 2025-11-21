# Product Requirements: Memory System

**Feature:** Conversation History and Memory Management  
**Version:** 1.0  
**Status:** Implemented  
**Owner:** Core Team  
**Last Updated:** December 2024

---

## Overview

The Memory System manages conversation history, ensuring the agent has access to relevant context while staying within token limits. It implements intelligent pruning strategies, maintains message integrity, and provides tools for users to understand and control what the agent remembers.

---

## Problem Statement

AI agents operating with limited context windows face critical memory challenges:

1. **Context Window Limits:** LLMs have finite context capacity (4K-128K tokens)
2. **Cost Concerns:** Every token costs money; unnecessary context wastes resources
3. **Relevance Decay:** Older messages become less relevant but still consume tokens
4. **Information Loss:** Naive pruning loses important context
5. **Performance Degradation:** Large contexts slow down LLM processing
6. **Conversation Continuity:** Users expect agent to remember earlier discussion

Without intelligent memory management, agents either:
- Run out of context space mid-conversation
- Forget important earlier information
- Waste tokens on irrelevant old messages
- Perform poorly due to context bloat

---

## Goals

### Primary Goals

1. **Maximize Relevance:** Keep most important context within token budget
2. **Maintain Continuity:** Preserve conversation coherence across long sessions
3. **Optimize Costs:** Minimize unnecessary token usage
4. **Transparent Pruning:** Make clear what's been removed and why
5. **User Control:** Allow users to influence memory strategy
6. **Graceful Degradation:** Handle approaching limits proactively

### Non-Goals

1. **Perfect Recall:** Does NOT aim to remember everything forever
2. **External Storage:** Does NOT persist conversations to disk (yet)
3. **Semantic Search:** Does NOT use embeddings for retrieval (simple strategy)
4. **Cross-Session Memory:** Does NOT remember across multiple sessions
5. **Context Compression:** Does NOT use LLM-based summarization (yet)

---

## User Personas

### Primary: Long-Session Developer
- **Background:** Developer working on complex tasks over extended periods
- **Workflow:** 1-2 hour coding sessions with many back-and-forth exchanges
- **Pain Points:** Agent forgets earlier decisions or context
- **Goals:** Maintain conversation continuity without manual intervention

### Secondary: Cost-Conscious Developer
- **Background:** Using paid LLM APIs, monitors token usage carefully
- **Workflow:** Frequent short sessions, worried about wasted tokens
- **Pain Points:** Paying for old irrelevant messages in context
- **Goals:** Efficient memory usage without sacrificing quality

### Tertiary: Multi-Task Developer
- **Background:** Switches between different coding tasks frequently
- **Workflow:** Starts fresh conversation for each new task
- **Pain Points:** Old task context pollutes new task thinking
- **Goals:** Clean slate for each task while preserving current task context

---

## Requirements

### Functional Requirements

#### FR1: Message Storage
- **R1.1:** Store all conversation messages in memory
- **R1.2:** Track message metadata (role, timestamp, tokens, importance)
- **R1.3:** Maintain message order chronologically
- **R1.4:** Support multiple message types (user, assistant, system, tool)
- **R1.5:** Calculate and cache token counts per message
- **R1.6:** Handle large messages efficiently

#### FR2: Token Tracking
- **R2.1:** Count tokens for each message accurately
- **R2.2:** Track cumulative token usage
- **R2.3:** Monitor against context window limit
- **R2.4:** Warn when approaching limit (80% threshold)
- **R2.5:** Provide token breakdown by message type
- **R2.6:** Update counts when messages pruned

#### FR3: Pruning Strategy
- **R3.1:** Implement rolling window strategy (keep last N messages)
- **R3.2:** Always preserve system prompt (never prune)
- **R3.3:** Always preserve most recent messages (last 10)
- **R3.4:** Prune from middle of conversation (oldest first)
- **R3.5:** Prune in message pairs (user + assistant together)
- **R3.6:** Leave pruning markers ("... [15 messages removed] ...")

#### FR4: Importance Scoring
- **R4.1:** Score messages by importance (0.0-1.0)
- **R4.2:** Higher score for messages with:
  - Tool calls and results
  - User directives or constraints
  - Error messages or corrections
  - Recent messages
- **R4.3:** Lower score for:
  - Old casual conversation
  - Redundant information
  - Successfully completed one-off tasks
- **R4.4:** Use scores to guide pruning decisions

#### FR5: Smart Pruning
- **R5.1:** Detect conversation segments (task boundaries)
- **R5.2:** Prune entire completed tasks before fragmenting current task
- **R5.3:** Keep error/correction pairs even if old
- **R5.4:** Preserve context necessary for current task
- **R5.5:** Avoid breaking conversation flow unnecessarily

#### FR6: Pruning Triggers
- **R6.1:** Automatic pruning at 80% context usage
- **R6.2:** Aggressive pruning at 90% context usage
- **R6.3:** User-initiated pruning (manual clear history)
- **R6.4:** Task-boundary pruning (between major tasks)
- **R6.5:** Session start pruning (reset on new conversation)

#### FR7: Memory Visibility
- **R7.1:** Show memory stats in context overlay
- **R7.2:** Indicate when history has been pruned
- **R7.3:** Display oldest message still in memory
- **R7.4:** Show total messages vs. messages in context
- **R7.5:** Explain pruning strategy in help text

#### FR8: Memory Configuration
- **R8.1:** Configure max messages (default: 50)
- **R8.2:** Configure max tokens (default: 80% of model limit)
- **R8.3:** Configure pruning strategy (rolling window, importance-based)
- **R8.4:** Configure minimum recent messages to keep (default: 10)
- **R8.5:** Settings accessible via settings overlay

#### FR9: Context Building
- **R9.1:** Build conversation context for each LLM call
- **R9.2:** Include system prompt + message history + current message
- **R9.3:** Apply token limit to total context
- **R9.4:** Prune if over limit before sending to LLM
- **R9.5:** Validate final context size

### Non-Functional Requirements

#### NFR1: Performance
- **N1.1:** Token counting under 10ms per message
- **N1.2:** Pruning decision under 50ms
- **N1.3:** Context building under 200ms
- **N1.4:** Memory stats calculation under 100ms
- **N1.5:** No noticeable lag during conversation

#### NFR2: Accuracy
- **N2.1:** Token counts within 5% of actual LLM usage
- **N2.2:** Pruning preserves conversation coherence >90% of time
- **N2.3:** Never prune system prompt or recent messages
- **N2.4:** Importance scoring correlates with user judgment >80%
- **N2.5:** Context limits never exceeded

#### NFR3: Memory Efficiency
- **N3.1:** Memory usage scales linearly with conversation length
- **N3.2:** No memory leaks over long sessions
- **N3.3:** Efficient storage format (minimize redundancy)
- **N3.4:** Lazy computation of token counts (cache results)
- **N3.5:** Total memory under 100MB for typical session

#### NFR4: Reliability
- **N4.1:** Never lose messages unexpectedly
- **N4.2:** Graceful handling of token counting errors
- **N4.3:** Consistent behavior across sessions
- **N4.4:** Recovery from pruning errors
- **N4.5:** Atomic operations (no partial states)

---

## User Experience

### Core Workflows

#### Workflow 1: Normal Conversation (No Pruning)
1. User starts conversation
2. Exchange 20 messages (under token limit)
3. Agent has full conversation context
4. No pruning occurs
5. User continues seamlessly

**Success Criteria:** No pruning needed for typical sessions

#### Workflow 2: Approaching Context Limit
1. User in long conversation (40 messages)
2. Context usage reaches 75%
3. Memory system shows: "Context: 75% used" in overlay
4. User continues conversation
5. At 80%, oldest messages pruned automatically
6. Chat shows: "... [5 messages removed to save space] ..."
7. Agent continues with relevant recent context

**Success Criteria:** Seamless pruning without disrupting conversation

#### Workflow 3: Task Boundary Pruning
1. User completes refactoring task
2. Agent: "Task complete. Refactored 3 files."
3. User: "Now let's work on the tests"
4. Memory system detects task boundary
5. Old refactoring messages pruned
6. Test discussion starts with clean context

**Success Criteria:** Task switching clears irrelevant history

#### Workflow 4: Manual History Clear
1. User wants fresh start mid-session
2. Types: `/clear` (hypothetical command)
3. Memory system: "Clear conversation history?"
4. User confirms
5. All messages except system prompt cleared
6. User starts new topic with clean slate

**Success Criteria:** User can reset conversation anytime

#### Workflow 5: Checking Memory State
1. User notices agent missed earlier context
2. Opens context overlay (Ctrl+I)
3. Sees: "Memory: 35 messages (10 pruned 15 min ago)"
4. Realizes early context was removed
5. User re-provides important information
6. Agent continues successfully

**Success Criteria:** User understands memory state and can compensate

---

## Technical Architecture

### Component Structure

```
Memory System
├── Message Store
│   ├── Message List
│   ├── Message Metadata
│   └── Token Cache
├── Pruning Engine
│   ├── Strategy Selector
│   ├── Importance Scorer
│   ├── Pruning Executor
│   └── Marker Injector
├── Token Accountant
│   ├── Tokenizer
│   ├── Token Counter
│   ├── Budget Tracker
│   └── Limit Monitor
├── Context Builder
│   ├── Message Selector
│   ├── Context Assembler
│   └── Validator
└── Memory Stats
    ├── Usage Calculator
    ├── Metrics Collector
    └── Reporter
```

### Data Model

```go
type Memory struct {
    messages        []Message
    systemPrompt    *Message
    maxMessages     int
    maxTokens       int
    strategy        PruningStrategy
    tokenizer       Tokenizer
    prunedCount     int
    lastPruneTime   time.Time
}

type Message struct {
    ID          string
    Role        MessageRole
    Content     string
    Timestamp   time.Time
    Tokens      int
    Importance  float64
    ToolCall    *ToolCall
    ToolResult  *ToolResult
}

type PruningStrategy int
const (
    StrategyRollingWindow PruningStrategy = iota
    StrategyImportanceBased
    StrategyHybrid
)

type MemoryStats struct {
    TotalMessages   int
    MessagesInContext int
    PrunedMessages  int
    TotalTokens     int
    TokenLimit      int
    PercentUsed     float64
    OldestMessage   time.Time
    LastPrune       time.Time
}
```

### Pruning Decision Flow

```
Context Building Request
    ↓
Calculate Current Token Usage
    ↓
Usage > 80% of limit?
    ├─ No → Use full history
    └─ Yes → Trigger Pruning
        ↓
    Select Pruning Strategy
        ↓
    Score All Messages for Importance
        ↓
    Sort Messages (oldest, lowest importance first)
        ↓
    ┌──────────────────────────────┐
    │ Pruning Loop:                │
    │ 1. Skip system prompt        │
    │ 2. Skip last 10 messages     │
    │ 3. Select lowest score msg   │
    │ 4. Remove message            │
    │ 5. Recalculate tokens        │
    │ 6. Under limit? → Stop       │
    │ 7. Else continue loop        │
    └──────────────────────────────┘
        ↓
    Insert Pruning Marker
        ↓
    Return Pruned Context
```

---

## Design Decisions

### Why Rolling Window Strategy?
**Rationale:**
- **Simple:** Easy to understand and implement
- **Predictable:** Users know recent messages are kept
- **Effective:** Recent context usually most relevant
- **Fast:** O(1) pruning decisions
- **Debuggable:** Clear what's kept vs removed

**Alternatives considered:**
- Importance-based only: More complex, harder to predict
- LLM summarization: Too slow, adds cost
- No pruning: Runs out of context

**Decision:** Start with rolling window, add importance-based in Phase 2

### Why Keep Last 10 Messages Always?
**Rationale:**
- **Continuity:** Last few exchanges critical for coherence
- **User expectation:** Users assume recent context is safe
- **Task completion:** Most tasks need last few messages
- **Safety margin:** Prevents breaking current task

**Testing showed:** 10 messages covers 95% of immediate context needs

### Why Prune at 80% Not 100%?
**Rationale:**
- **Headroom:** Prevents emergency pruning mid-response
- **Efficiency:** Proactive pruning is cheaper than reactive
- **User experience:** Gradual pruning less jarring
- **Error prevention:** Avoids hitting hard limits

### Why Not Use Embeddings for Retrieval?
**Current decision:** Simple chronological pruning

**Rationale:**
- **Complexity:** Embeddings add significant overhead
- **Cost:** Requires separate embedding API calls
- **Performance:** Embedding generation is slow
- **Overkill:** Rolling window works for 90% of use cases

**Future:** Phase 3 may add semantic retrieval for long sessions

---

## Pruning Strategies Comparison

### 1. Rolling Window (Current)
**How it works:** Keep last N messages, remove oldest

**Pros:**
- Simple and fast
- Predictable behavior
- Low overhead

**Cons:**
- May remove important old messages
- No semantic understanding

**Best for:** Most conversations, bounded context

---

### 2. Importance-Based (Future)
**How it works:** Score messages by importance, remove lowest scores

**Pros:**
- Preserves critical context
- Smarter than pure chronological

**Cons:**
- More complex scoring logic
- Slower pruning decisions
- Less predictable

**Best for:** Long sessions with diverse topics

---

### 3. Hybrid (Future)
**How it works:** Combine rolling window + importance

**Pros:**
- Balance of simplicity and intelligence
- Keeps recent + important

**Cons:**
- Most complex implementation
- Tuning required

**Best for:** Power users, extended sessions

---

## Importance Scoring Algorithm

```
Base Score: 0.5 (neutral)

Modifiers:
  + 0.3 if contains tool call
  + 0.2 if contains tool result
  + 0.2 if user directive/constraint
  + 0.3 if error message
  + 0.2 if correction/clarification
  + 0.1 for each reference to earlier message
  - 0.1 for each day old
  - 0.2 if casual/social content
  - 0.3 if task marked complete

Recency Boost:
  + 0.4 if in last 5 messages
  + 0.2 if in last 10 messages
  + 0.1 if in last 20 messages

System Prompt: 1.0 (never pruned)
Final Score: Clamp to [0.0, 1.0]
```

---

## Success Metrics

### Effectiveness Metrics
- **Context preservation:** >90% of conversations maintain coherence after pruning
- **Relevance ratio:** >80% of tokens in context are relevant to current task
- **Prune accuracy:** <5% of pruned messages later needed
- **Task completion:** >95% of tasks complete without context issues

### Efficiency Metrics
- **Token waste:** <10% of context tokens are irrelevant
- **Pruning overhead:** <5% of session time spent on pruning
- **Memory usage:** <100MB for typical session
- **Cache hit rate:** >90% of token counts served from cache

### User Experience Metrics
- **Surprise rate:** <5% of users surprised by forgotten context
- **Manual intervention:** <10% of sessions require manual history management
- **Continuity score:** >85% of conversations feel coherent despite pruning
- **Transparency:** >80% of users understand memory state via overlay

---

## Dependencies

### External Dependencies
- Tokenizer library (tiktoken for OpenAI, custom for others)
- Token counting utilities

### Internal Dependencies
- LLM provider (for context limits)
- Settings system (for memory configuration)
- Context overlay (for visibility)
- Event system (for pruning notifications)

### Platform Requirements
- Sufficient RAM for conversation storage
- Fast token counting (efficient tokenizer)

---

## Risks & Mitigations

### Risk 1: Important Context Pruned
**Impact:** High  
**Probability:** Medium  
**Mitigation:**
- Keep generous recent window (10+ messages)
- Importance scoring to protect critical messages
- User can see what was pruned
- Easy to re-provide lost context
- Future: Summarization instead of deletion

### Risk 2: Excessive Token Counting Overhead
**Impact:** Medium  
**Probability:** Low  
**Mitigation:**
- Cache token counts aggressively
- Lazy computation (only when needed)
- Efficient tokenizer implementation
- Batch token counting where possible
- Monitor and optimize hot paths

### Risk 3: Inaccurate Token Counts
**Impact:** High  
**Probability:** Low  
**Mitigation:**
- Use same tokenizer as LLM provider
- Validate counts against API responses
- Conservative estimates (slight overcount)
- Regular validation tests
- Update tokenizer with provider changes

### Risk 4: Pruning Breaks Conversation Flow
**Impact:** Medium  
**Probability:** Medium  
**Mitigation:**
- Always preserve recent messages
- Prune in logical units (message pairs)
- Leave clear pruning markers
- Test pruning logic extensively
- User feedback loop for improvements

### Risk 5: Memory Leaks
**Impact:** Medium  
**Probability:** Low  
**Mitigation:**
- Hard limits on message count
- Regular cleanup of old sessions
- Memory profiling in tests
- Bounded data structures
- Clear session end handling

---

## Future Enhancements

### Phase 2 Ideas
- **Importance-Based Pruning:** Smart scoring algorithm
- **Summarization:** LLM-generated summaries of pruned content
- **User Hints:** Let users mark messages as "important"
- **Topic Detection:** Identify and preserve topic boundaries
- **Configurable Strategies:** Multiple pruning modes to choose from

### Phase 3 Ideas
- **Semantic Retrieval:** Embedding-based context selection
- **Long-Term Memory:** Persist important facts across sessions
- **Memory Visualization:** Graph of conversation structure
- **Adaptive Pruning:** Learn optimal strategy from user behavior
- **Memory Compression:** More efficient storage formats

---

## Open Questions

1. **Should we persist conversation history to disk?**
   - Pro: Resume sessions, audit trail
   - Con: Privacy, storage management, complexity
   - Decision: Phase 3 feature with encryption

2. **Should we use LLM summarization for pruned content?**
   - Pro: Preserves information in condensed form
   - Con: Cost, latency, quality variability
   - Decision: Phase 2 experiment

3. **Should users be able to edit message history?**
   - Use case: Remove incorrect information
   - Risk: Confusion about what agent "knows"
   - Decision: Phase 3, with clear UI indicators

4. **Should we support multiple memory strategies per session?**
   - Pro: Flexibility for different tasks
   - Con: UI complexity, user confusion
   - Decision: Single strategy per session for simplicity

---

## Related Documentation

- [ADR-0005: Memory and Context Management](../adr/0005-memory-and-context-management.md)
- [ADR-0014: Intelligent Context Pruning](../adr/0014-intelligent-context-pruning.md)
- [Context Management PRD](context-management.md)
- [Agent Loop Architecture PRD](agent-loop-architecture.md)
- [Architecture: Memory System](../architecture/memory.md)

---

## Changelog

| Date | Version | Changes |
|------|---------|---------|
| 2024-12 | 1.0 | Initial PRD creation |
