# 18. Selective Tool Call Summarization with Exclusion Rules

**Status:** Accepted
**Date:** 2025-01-17
**Deciders:** Development Team
**Technical Story:** Prevent loss of critical agent context by excluding certain tool calls from summarization while maintaining efficient context management for large operational tools.

---

## Context

### Background

The current `ToolCallSummarizationStrategy` (implemented in [ADR-0015](0015-buffered-tool-call-summarization.md)) successfully reduces token usage by summarizing old tool calls and their results. However, it treats all tool calls uniformly, summarizing every tool call that meets the age and buffer criteria.

In practice, not all tool calls should be summarized. Some tool calls represent critical interaction points or contain semantic context that the agent needs to maintain throughout the conversation.

### Problem Statement

When the agent summarizes certain types of tool calls, it loses important context:

1. **Loop-Breaking Tools** (`task_completion`, `ask_question`, `converse`):
   - These represent key interaction points with the user
   - `task_completion`: Contains final results and conclusions that may be referenced later
   - `ask_question`: Contains questions asked to the user and their context
   - `converse`: Contains conversational context that maintains dialogue continuity

2. **Semantic Importance**:
   - These tools are typically small (low token count)
   - Their semantic value is high relative to their size
   - Summarizing them provides minimal token savings but loses valuable context

3. **Large Operational Tools**:
   - Tools like `read_file`, `execute_command`, `search_files`, `apply_diff` can be very large
   - These are excellent candidates for summarization (high token savings)
   - Their full output is often not needed after the immediate task

### Current Behavior Example

Consider a conversation where the agent:
1. Asks a clarifying question (turn 10)
2. User responds (turn 11)
3. Agent performs 20+ file operations (turns 12-32)
4. Completes the task (turn 33)

Currently, if the buffer trigger fires at turn 50, both the `ask_question` call from turn 10 AND all the file operations get summarized. This loses the context of what question was asked and why, which could be referenced later.

### Goals

- Preserve critical agent context by excluding certain tool calls from summarization
- Maintain efficient token management for large operational tools
- Provide sensible defaults while allowing customization
- Keep the implementation simple and maintainable
- Align with the composable strategy pattern from [ADR-0014](0014-composable-context-management.md)

### Non-Goals

- Dynamic/ML-based detection of "important" vs "unimportant" content
- Per-instance exclusion rules (exclusions apply to all instances of a tool)
- User-facing configuration UI (initially - may add later)
- Exclusion of specific tool call instances by ID

---

## Decision Drivers

* **Context Quality**: Preserve semantically important interaction points
* **Token Efficiency**: Still summarize large operational tools that provide minimal ongoing value
* **Simplicity**: Keep the exclusion mechanism straightforward and predictable
* **Defaults**: Provide sensible defaults that work for 95% of use cases
* **Flexibility**: Allow power users to customize exclusions if needed
* **Backward Compatibility**: Don't break existing behavior for users who don't need exclusions

---

## Considered Options

### Option 1: Tool Name Blacklist

**Description:** Maintain a list of tool names that should never be summarized. Check each tool call against this list during grouping.

**Pros:**
- Simple to implement and understand
- Clear semantics: "these tools are never summarized"
- Easy to configure with a simple string slice
- Fast lookup (can use map for O(1) checking)
- Aligns with how tools are already identified (by name)

**Cons:**
- All-or-nothing: can't exclude based on size or content
- Requires maintaining a list of tool names
- No automatic detection of loop-breaking tools

### Option 2: Automatic Loop-Breaking Detection

**Description:** Automatically detect loop-breaking tools by checking the `IsLoopBreaking()` method and exclude them from summarization.

**Pros:**
- Automatic: no configuration needed
- Self-maintaining: new loop-breaking tools automatically excluded
- Semantically correct: loop-breaking tools represent important context
- Type-safe: uses existing tool interface

**Cons:**
- Requires access to tool registry at summarization time
- More complex: need to parse tool names and look up tools
- Tight coupling between context management and tool system
- Can't override for specific use cases

### Option 3: Metadata-Based Exclusion

**Description:** Add metadata to messages indicating whether they should be excluded from summarization. Tools mark their messages during execution.

**Pros:**
- Very flexible: can mark specific instances
- Decentralized: each tool controls its own summarization
- Can handle dynamic exclusion logic
- No central configuration needed

**Cons:**
- Requires changes to all tools
- More complex: metadata must be checked during grouping
- Harder to reason about: exclusion rules scattered across codebase
- No central visibility of what's excluded

---

## Decision

**Chosen Option:** Option 1 - Tool Name Blacklist with Sensible Defaults

### Rationale

The tool name blacklist approach provides the best balance of simplicity, flexibility, and maintainability:

1. **Simplicity**: A string slice of tool names is easy to understand and configure
2. **Predictability**: Explicit exclusions make behavior clear and debuggable
3. **Flexibility**: Can be customized per agent configuration if needed
4. **Performance**: Fast O(1) lookup using a map internally
5. **Maintenance**: Central configuration makes it easy to see what's excluded
6. **Defaults**: We can provide sensible defaults (loop-breaking tools) out of the box

---

## Consequences

### Positive

- **Better Context Quality**: Critical interaction points (questions, task completions) are preserved in full
- **Maintained Efficiency**: Large operational tools (read_file, execute_command) are still summarized
- **Clear Behavior**: Explicit exclusion list makes it obvious what won't be summarized
- **Easy Configuration**: Simple string slice can be set during strategy creation
- **Sensible Defaults**: Works correctly out of the box for most use cases
- **Low Overhead**: Exclusion check is fast and happens during grouping, not per-message

### Negative

- **Manual Maintenance**: New tools that should be excluded must be manually added to defaults
- **No Dynamic Behavior**: Can't exclude based on tool call content or size
- **Configuration Burden**: Power users must know which tools to exclude
- **Potential for Mistakes**: Forgetting to exclude an important tool could lose context

### Neutral

- Changes which tool calls get summarized but not the summarization quality
- Excluded tools still consume tokens (no compression)
- May need to tune exclusion list based on usage patterns

---

## Implementation

See the full implementation details in the attached design document.

### Key Changes Summary

1. Add `excludedTools map[string]bool` field to `ToolCallSummarizationStrategy`
2. Update constructor to accept optional exclusion list with sensible defaults
3. Modify `groupToolCallsAndResults()` to skip excluded tools during grouping
4. Add `extractToolName()` helper to parse tool names from XML
5. Update `Summarize()` to pass exclusion set to grouping function

### Default Exclusions

By default, the following tools are excluded from summarization:
- `task_completion`
- `ask_question`
- `converse`

These can be overridden by providing a custom exclusion list to the constructor.

---

## Validation

### Success Metrics

- **Context Preservation**: Agent can reference task completions and questions from >20 turns ago
- **Token Efficiency**: Large operational tools are still summarized, maintaining 40-60% token savings
- **Default Behavior**: 95% of users don't need to customize exclusions
- **Performance**: Exclusion check adds <1ms overhead per tool call
- **Backward Compatibility**: Existing configurations continue to work

### Test Scenarios

1. **Long Session with Mixed Tools**: Verify loop-breaking tools are preserved while operational tools are summarized
2. **Custom Exclusions**: Test that custom exclusion lists work correctly
3. **No Exclusions**: Verify empty exclusion list allows all tools to be summarized
4. **Performance**: Measure grouping time with and without exclusions

---

## Related Decisions

- [ADR-0014](0014-composable-context-management.md) - Composable Context Management System
- [ADR-0015](0015-buffered-tool-call-summarization.md) - Buffered Tool Call Summarization
- [ADR-0008](0008-agent-controlled-loop-termination.md) - Agent Controlled Loop Termination

---

## References

- [Loop-Breaking Tools](../../pkg/agent/tools/tool.go) - Tool interface definition
- [Tool Call Strategy](../../pkg/agent/context/tool_call_strategy.go) - Current implementation

---

## Notes

### Design Philosophy

This ADR follows the principle of "sensible defaults with explicit overrides":
- Default behavior works correctly for most users
- Power users can customize when needed
- Explicit configuration makes behavior predictable
- Simple implementation reduces bugs

### Future Enhancements

1. **Size-Based Exclusion**: Exclude tool calls smaller than N tokens
2. **Metadata-Based Marking**: Allow tools to mark themselves as "do not summarize"
3. **Configurable Defaults**: Allow users to disable default exclusions
4. **Exclusion Metrics**: Emit events showing which tools were excluded
5. **Smart Exclusions**: Use LLM to determine if a tool call is "important"

---

**Last Updated:** 2025-01-17
**Implementation Status:** Implemented
