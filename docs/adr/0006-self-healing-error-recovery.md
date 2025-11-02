# 6. Self-Healing Error Recovery with Circuit Breaker

**Status:** Accepted
**Date:** 2025-10-29
**Deciders:** Forge Core Team
**Technical Story:** Enabling agents to recover from errors without infinite loops

## Context

Agent loops can encounter various errors during execution: invalid tool calls, unknown tools, execution failures, or missing required fields. Traditional approaches either terminate on first error (fragile) or retry indefinitely (potentially infinite loops). We needed a robust error recovery mechanism that allows agents to learn from mistakes and self-correct, while protecting against pathological cases.

### Problem Statement

When an LLM agent makes a mistake (e.g., malformed JSON, calls non-existent tool, provides invalid arguments), we face several challenges:

1. **Immediate Termination is Too Fragile**: A single mistake shouldn't crash the entire agent session
2. **Context Loss**: The agent loses the opportunity to understand what went wrong
3. **Infinite Loops**: Naive retry mechanisms can get stuck repeating the same error
4. **User Experience**: Agents should gracefully handle errors and provide helpful responses
5. **Resource Protection**: Need to prevent runaway loops consuming API calls

### Requirements

- Allow agents to recover from recoverable errors (tool issues, parsing failures)
- Terminate on unrecoverable errors (LLM API failures, context cancellation)
- Prevent infinite loops when the same error repeats
- Provide clear error context to help agents self-correct
- Maintain clean conversation history (no error pollution)

## Decision

We implemented a **self-healing error recovery system** with two key components:

### 1. Ephemeral Error Context

When a recoverable error occurs, we:
- Generate a descriptive error message with recovery instructions
- Pass it as `errorContext` to the next LLM call
- Include it as a temporary user message (via [`BuildMessages()`](../pkg/agent/prompts/builder.go:77))
- **Do not** add it to permanent conversation history
- Discard after the iteration completes

This gives the agent immediate context to self-correct without polluting memory.

### 2. Circuit Breaker Pattern

To prevent infinite error loops, we track errors in a **ring buffer**:

```go
type DefaultAgent struct {
    lastErrors [5]string  // Ring buffer of last 5 error messages
    errorIndex int        // Current position in ring buffer
}
```

**Trigger Condition**: Circuit breaker activates when:
1. All 5 slots in the ring buffer are filled (non-empty)
2. All 5 error messages are identical (exact string match)

**Reset Behavior**: Error tracking resets to zero when:
- Any tool executes successfully
- A different error occurs (breaking the pattern)

### Error Classification

**Terminal Errors** (immediate loop termination):
- LLM API failures (provider communication errors)
- Context cancellation (timeout, user interrupt)
- Circuit breaker activation (5 identical consecutive errors)

**Recoverable Errors** (trigger self-healing):
- No tool call found in response
- Invalid JSON in tool call
- Missing required fields (`tool_name`)
- Unknown tool (not registered)
- Tool execution failure

### Implementation

**Error Tracking** ([`trackError()`](../pkg/agent/default.go:367)):
```go
func (a *DefaultAgent) trackError(errMsg string) bool {
    // Add to ring buffer
    a.lastErrors[a.errorIndex] = errMsg
    a.errorIndex = (a.errorIndex + 1) % 5

    // Check if all 5 are identical and non-empty
    if a.lastErrors[0] == "" {
        return false // Not enough errors yet
    }

    first := a.lastErrors[0]
    for i := 1; i < 5; i++ {
        if a.lastErrors[i] != first {
            return false
        }
    }

    return true // All 5 errors are identical - trigger circuit breaker
}
```

**Error Reset** ([`resetErrorTracking()`](../pkg/agent/default.go:388)):
```go
func (a *DefaultAgent) resetErrorTracking() {
    for i := range a.lastErrors {
        a.lastErrors[i] = ""
    }
    a.errorIndex = 0
}
```

**Agent Loop Integration** ([`runAgentLoop()`](../pkg/agent/default.go:231)):
```go
func (a *DefaultAgent) runAgentLoop(ctx context.Context) error {
    var errorContext string
    
    for {
        shouldContinue, nextError, err := a.executeIteration(ctx, errorContext)
        if err != nil {
            return err // Terminal error
        }
        if !shouldContinue {
            return nil // Loop-breaking tool or circuit breaker
        }
        
        errorContext = nextError // Pass to next iteration
    }
}
```

### Error Message Format

Error messages follow a consistent structure via [`BuildErrorRecoveryMessage()`](../pkg/agent/prompts/builder.go):

```
The previous iteration failed with the following error:

[Detailed error message]

Please analyze what went wrong and try a different approach.
[Specific recovery instructions based on error type]
```

**Error-Specific Instructions:**
- **No Tool Call**: "You must provide a tool call in your response."
- **Invalid JSON**: "Ensure your tool call JSON is properly formatted."
- **Unknown Tool**: "Use one of the available tools: [tool1, tool2, ...]"
- **Execution Failure**: "The error was: [specific error message]"

## Alternatives Considered

### 1. Permanent Error History

**Approach**: Add all errors to conversation memory

**Rejected because**:
- Pollutes conversation history with implementation details
- Wastes token budget on error messages
- Confuses the agent with accumulated error context
- Makes conversation history harder to debug

### 2. Simple Retry Count

**Approach**: Track number of consecutive errors with a counter

**Rejected because**:
- Doesn't distinguish between different error types
- Could terminate on unrelated sequential errors
- Less precise than pattern matching identical errors
- Misses the key insight: repeated identical errors indicate stuck state

### 3. Exponential Backoff

**Approach**: Increase delay between retries after errors

**Rejected because**:
- Adds latency without solving the core problem
- LLM responses are already slow; backoff makes worse UX
- Doesn't help the agent understand and correct mistakes
- Delays are irrelevant when errors stem from logic, not timing

### 4. Error Hash/Checksum

**Approach**: Hash error messages to reduce memory footprint

**Rejected because**:
- String comparison is fast enough for 5 elements
- Adds complexity without meaningful benefit
- Ring buffer of 5 strings is negligible memory overhead
- Direct string comparison is more debuggable

### 5. Configurable Circuit Breaker Threshold

**Approach**: Allow users to configure the "5 consecutive errors" threshold

**Rejected because**:
- 5 is a reasonable default based on experience
- Adding configuration increases API surface
- Most users won't know the right value
- Premature optimization without proven need
- Can be added later if needed

## Consequences

### Positive

- **Self-Correction**: Agents can recover from mistakes by analyzing error messages
- **Better UX**: Agent provides helpful feedback instead of crashing
- **Resource Protection**: Circuit breaker prevents runaway API consumption
- **Clean History**: Ephemeral errors don't pollute conversation memory
- **Flexible Recovery**: Different error types get appropriate guidance
- **Event Visibility**: Error events enable observability and debugging
- **Pattern Detection**: Ring buffer identifies stuck states precisely

### Negative

- **Added Complexity**: Error tracking state in agent struct
- **Memory Overhead**: Ring buffer of 5 strings (negligible in practice)
- **Delayed Termination**: Circuit breaker requires 5 identical errors before triggering
- **State Management**: Must remember to reset error tracking on success
- **Testing Burden**: Requires comprehensive error scenario tests

### Neutral

- **Error Context Ephemeral**: Trade-off between clean history and full context
- **Fixed Threshold**: 5 consecutive errors is hardcoded (could be configurable later)
- **String Comparison**: Uses exact string match (could use fuzzy matching)

## Implementation Details

### Error Flow Example

```
User: "Calculate 100 divided by 0"

Iteration 1:
├─ Agent calls calculator tool with {a: 100, b: 0, op: "divide"}
├─ Tool returns error: "division by zero"
├─ Error tracked in ring buffer [0]: "division by zero..."
├─ Error context generated
└─ Next iteration with error context

Iteration 2:
├─ Agent receives error context message
├─ Agent calls ask_question tool to inform user
├─ Tool executes successfully
├─ Error tracking reset
└─ Agent: "I cannot divide by zero. Would you like a different calculation?"
```

### Testing

Comprehensive tests in [`pkg/agent/error_recovery_test.go`](../pkg/agent/error_recovery_test.go):

- `TestErrorTracking/TracksSingleError`: Verifies single error tracking
- `TestErrorTracking/TracksMultipleDifferentErrors`: Different errors don't trigger circuit breaker
- `TestErrorTracking/TriggersCircuitBreakerOn5IdenticalErrors`: Circuit breaker activation
- `TestErrorTracking/ResetsAfterSuccessfulIteration`: Error counter reset
- `TestErrorTracking/CircuitBreakerRequiresAllFiveSlotsFilled`: All slots must be filled

Run with:
```bash
go test ./pkg/agent -run TestErrorTracking -v
```

### Related Events

Error recovery emits several event types:
- [`ErrorEvent`](../pkg/types/events.go): Generic error notification
- [`NoToolCallEvent`](../pkg/types/events.go): No tool call found
- [`ToolResultErrorEvent`](../pkg/types/events.go): Tool execution failed

These enable external monitoring and debugging.

## Related Decisions

- [ADR-0001: Record Architecture Decisions](0001-record-architecture-decisions.md) - Establishes ADR process
- [ADR-0002: XML Format for Tool Calls](0002-xml-format-for-tool-calls.md) - Tool call parsing that can fail
- [ADR-0005: Channel-Based Agent Communication](0005-channel-based-agent-communication.md) - Event emission for errors

## References

- Circuit Breaker Pattern: Martin Fowler's [CircuitBreaker](https://martinfowler.com/bliki/CircuitBreaker.html)
- Error Recovery Document: [`docs/archive/error-recovery.md`](../archive/error-recovery.md)
- Implementation: [`pkg/agent/default.go`](../pkg/agent/default.go:366-393)
- Tests: [`pkg/agent/error_recovery_test.go`](../pkg/agent/error_recovery_test.go)