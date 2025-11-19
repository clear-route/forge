# Phase 2 Refactoring: Detailed Planning Document

**Created:** 2024-12-19
**Status:** Planning
**Target:** Simplify agent loop and improve code maintainability

---

## Executive Summary

Phase 2 focuses on **simplifying the agent loop** in `pkg/agent/default.go` (898 lines, 32 functions). The primary goal is to reduce complexity by extracting helper methods, standardizing error handling, and improving code organization without changing behavior.

**Key Metrics:**
- Current file size: 898 lines
- Current function count: 32 functions
- Largest functions: `executeIteration()`, `executeTool()`, `processToolCall()`
- Target: All functions <10 cyclomatic complexity

---

## Current State Analysis

### File Structure: `pkg/agent/default.go`

**Complex Functions Requiring Refactoring:**

1. **`executeIteration()` (lines 364-470, ~107 lines)**
   - Responsibilities: Message building, context management, LLM calling, token tracking, stream processing
   - Complexity: High - orchestrates 5+ different concerns
   - Refactoring potential: Extract 5-6 helper methods

2. **`executeTool()` (lines 775-871, ~95 lines)**
   - Responsibilities: Tool lookup, approval flow, execution, result processing, error handling
   - Complexity: High - handles approval, execution, and circuit breaker logic
   - Refactoring potential: Extract 4-5 helper methods

3. **`processToolCall()` (lines 696-773, ~76 lines)**
   - Responsibilities: Parsing, validation, circuit breaker logic
   - Complexity: Medium-high - error handling and validation
   - Refactoring potential: Extract 3-4 helper methods

4. **`GetContextInfo()` (lines 565-643, ~79 lines)**
   - Responsibilities: Token counting, metrics calculation
   - Complexity: Medium - lots of calculations
   - Refactoring potential: Extract calculation helpers

**Supporting Infrastructure:**

- Error tracking: `trackError()`, `resetErrorTracking()` (lines 666-694)
- Tool management: `getTool()`, `getToolsList()` (lines 646-664)
- Approval system: `requestApproval()`, `handleApprovalResponse()` (lines 873-898)
- Event loop: `eventLoop()`, `processInput()` (lines 207-332)

---

## Task 2.1: Simplify Agent Loop Methods

**Estimated:** 8 hours
**Branch:** `refactor/simplify-agent-loop`

### Goals

1. **Reduce `executeIteration()` from 107 lines to <30 lines**
2. **Reduce `executeTool()` from 95 lines to <40 lines**
3. **Reduce `processToolCall()` from 76 lines to <30 lines**
4. **All functions achieve <10 cyclomatic complexity**
5. **Zero behavior changes - maintain exact same logic**

### Proposed Extraction Strategy

#### A. `executeIteration()` Refactoring

**Current structure:**
```
executeIteration() (107 lines)
├── Build messages
├── Check context limits
├── Handle summarization
├── Emit API call event
├── Stream completion
├── Process stream
├── Count tokens
├── Emit token usage
├── Add to memory
└── Process tool call
```

**Proposed helpers:**

1. **`buildIterationMessages(systemPrompt, errorContext)`** (15-20 lines)
   - Build messages from memory
   - Add error context if present
   - Return messages array

2. **`evaluateContextManagement(messages, systemPrompt)`** (20-25 lines)
   - Count tokens
   - Check if over limit
   - Trigger summarization if needed
   - Return updated messages

3. **`callLLMWithContext(ctx, messages, promptTokens)`** (15-20 lines)
   - Emit API call start event
   - Call provider.StreamCompletion
   - Handle cancellation
   - Return stream or error

4. **`processStreamResponse(stream)`** (20-25 lines)
   - Process stream events
   - Collect assistant content
   - Collect tool call content
   - Return both strings

5. **`trackTokenUsage(assistantContent, toolCallContent, promptTokens)`** (15-20 lines)
   - Count completion tokens
   - Calculate totals
   - Emit token usage event
   - Return completion tokens

6. **`addResponseToMemory(assistantContent, toolCallContent)`** (10-15 lines)
   - Build full response
   - Create message
   - Add to memory

**New `executeIteration()` structure:**
```go
func (a *DefaultAgent) executeIteration(ctx context.Context, errorContext string) (bool, string) {
    systemPrompt := a.buildSystemPrompt()
    
    // Build and manage messages
    messages := a.buildIterationMessages(systemPrompt, errorContext)
    messages, promptTokens := a.evaluateContextManagement(messages, systemPrompt)
    
    // Call LLM
    stream, err := a.callLLMWithContext(ctx, messages, promptTokens)
    if err != nil {
        return false, "" // Handle errors internally
    }
    
    // Process response
    assistantContent, toolCallContent := a.processStreamResponse(stream)
    a.trackTokenUsage(assistantContent, toolCallContent, promptTokens)
    a.addResponseToMemory(assistantContent, toolCallContent)
    
    // Process tool call
    return a.processToolCall(ctx, toolCallContent)
}
```

#### B. `executeTool()` Refactoring

**Current structure:**
```
executeTool() (95 lines)
├── Look up tool
├── Check if previewable
├── Generate preview
├── Request approval
├── Handle timeout/rejection
├── Emit tool call event
├── Inject context values
├── Execute tool
├── Handle error with circuit breaker
├── Emit result
├── Reset error tracking
└── Check if loop-breaking
```

**Proposed helpers:**

1. **`lookupTool(toolName)`** (15-20 lines)
   - Look up tool in registry
   - Build error message if not found
   - Track error and check circuit breaker
   - Return tool or error context

2. **`checkToolApproval(ctx, tool, toolCall)`** (25-30 lines)
   - Check if tool is previewable
   - Generate preview
   - Request approval
   - Handle timeout/rejection
   - Return (approved, shouldContinue, errorContext)

3. **`executeToolWithContext(ctx, tool, toolCall)`** (20-25 lines)
   - Emit tool call event
   - Inject event emitter and registry
   - Execute tool
   - Return result or error

4. **`handleToolSuccess(toolCall, result, tool)`** (15-20 lines)
   - Emit result event
   - Reset error tracking
   - Check if loop-breaking
   - Add to memory if continuing
   - Return shouldContinue

5. **`handleToolExecutionError(toolName, err)`** (20-25 lines)
   - Emit error events
   - Build error recovery message
   - Track error and check circuit breaker
   - Return errorContext

**New `executeTool()` structure:**
```go
func (a *DefaultAgent) executeTool(ctx context.Context, toolCall tools.ToolCall) (bool, string) {
    // Look up tool
    tool, errCtx := a.lookupTool(toolCall.ToolName)
    if errCtx != "" {
        return true, errCtx
    }
    
    // Check approval
    approved, shouldContinue, errCtx := a.checkToolApproval(ctx, tool, toolCall)
    if !approved {
        return shouldContinue, errCtx
    }
    
    // Execute tool
    result, err := a.executeToolWithContext(ctx, tool, toolCall)
    if err != nil {
        return a.handleToolExecutionError(toolCall.ToolName, err)
    }
    
    // Handle success
    return a.handleToolSuccess(toolCall, result, tool)
}
```

#### C. `processToolCall()` Refactoring

**Current structure:**
```
processToolCall() (76 lines)
├── Check context cancellation
├── Check if tool call exists
├── Parse XML
├── Handle parse errors
├── Validate tool name
├── Set default server name
└── Execute tool
```

**Proposed helpers:**

1. **`validateToolCallContent(ctx, toolCallContent)`** (20-25 lines)
   - Check context cancellation
   - Check if content exists
   - Emit no tool call event if missing
   - Track error and check circuit breaker
   - Return error context

2. **`parseToolCallXML(toolCallContent)`** (20-25 lines)
   - Wrap content in <tool> tags
   - Parse XML
   - Handle parse errors
   - Track error and check circuit breaker
   - Return parsed tool call or error context

3. **`validateToolCall(toolCall)`** (15-20 lines)
   - Check tool name exists
   - Set default server name if needed
   - Track error if validation fails
   - Return validated tool call or error context

**New `processToolCall()` structure:**
```go
func (a *DefaultAgent) processToolCall(ctx context.Context, toolCallContent string) (bool, string) {
    // Validate content
    if errCtx := a.validateToolCallContent(ctx, toolCallContent); errCtx != "" {
        return a.shouldContinueAfterError(errCtx)
    }
    
    // Parse XML
    toolCall, errCtx := a.parseToolCallXML(toolCallContent)
    if errCtx != "" {
        return true, errCtx
    }
    
    // Validate and normalize
    toolCall, errCtx = a.validateToolCall(toolCall)
    if errCtx != "" {
        return true, errCtx
    }
    
    // Execute
    return a.executeTool(ctx, toolCall)
}
```

---

## Task 2.2: Standardize Error Handling

**Estimated:** 6 hours
**Branch:** `refactor/standardize-errors`

### Goals

1. **Create dedicated error package** (`pkg/agent/errors`)
2. **Define structured error types**
3. **Standardize error creation and wrapping**
4. **Improve error context and debugging**

### Proposed Error Package Structure

```
pkg/agent/errors/
├── errors.go       # Core error types and constructors
├── types.go        # ErrorType enum and constants
└── recovery.go     # Error recovery message building
```

#### Error Type Definitions

```go
// ErrorType represents categories of agent errors
type ErrorType int

const (
    ErrorTypeNoToolCall ErrorType = iota
    ErrorTypeInvalidXML
    ErrorTypeMissingToolName
    ErrorTypeUnknownTool
    ErrorTypeToolExecution
    ErrorTypeCircuitBreaker
    ErrorTypeLLMFailure
    ErrorTypeContextCanceled
    ErrorTypeApprovalTimeout
    ErrorTypeApprovalRejected
)

// AgentError is a structured error with context
type AgentError struct {
    Type      ErrorType
    Message   string
    Cause     error
    Context   map[string]interface{}
    Timestamp time.Time
}
```

#### Migration Strategy

1. **Phase 1:** Create error package (2 hours)
   - Define error types
   - Implement constructors
   - Add tests

2. **Phase 2:** Update agent code (3 hours)
   - Replace inline error handling
   - Use structured errors
   - Update tests

3. **Phase 3:** Integration (1 hour)
   - Verify all error paths
   - Update documentation
   - Run full test suite

---

## Task 2.3: Consolidate Overlay Components

**Estimated:** 4 hours
**Branch:** `refactor/consolidate-overlays`

### Goals

1. **Reduce duplication in overlay code**
2. **Standardize overlay styling**
3. **Improve overlay component reusability**

### Current Overlay Structure

```
pkg/executor/tui/overlay/
├── settings.go (1,151 lines) ⚠️ LARGEST FILE
├── palette.go
├── approval.go
├── command.go
└── types.go
```

### Refactoring Strategy

**Focus on `settings.go` (1,151 lines):**

1. **Extract rendering helpers** (2 hours)
   - Common list rendering
   - Common input rendering
   - Common button rendering

2. **Split into modules** (2 hours)
   - `settings/model.go` - State and types
   - `settings/init.go` - Initialization
   - `settings/update.go` - Update logic
   - `settings/view.go` - View rendering
   - `settings/helpers.go` - Helper functions

---

## Task 2.4: Implement Structured Logging

**Estimated:** 6 hours
**Branch:** `refactor/structured-logging`

### Goals

1. **Replace `agentDebugLog` with structured logger**
2. **Add log levels (DEBUG, INFO, WARN, ERROR)**
3. **Enable optional file logging**
4. **Add context-aware logging**

### Proposed Logging Structure

```go
// pkg/agent/logging/logger.go
type Logger struct {
    level   Level
    output  io.Writer
    context map[string]interface{}
}

func (l *Logger) Debug(msg string, fields ...Field)
func (l *Logger) Info(msg string, fields ...Field)
func (l *Logger) Warn(msg string, fields ...Field)
func (l *Logger) Error(msg string, fields ...Field)
func (l *Logger) WithContext(fields ...Field) *Logger
```

### Migration Strategy

1. **Create logging package** (2 hours)
2. **Replace agentDebugLog** (3 hours)
3. **Add context to key operations** (1 hour)

---

## Testing Strategy

### Test Coverage Goals

- **Maintain current coverage:** 20.8%
- **Add tests for new helper methods**
- **Ensure no behavioral changes**

### Verification Steps (Per Task)

1. ✅ **Unit tests pass:** `make test`
2. ✅ **No linter errors:** `make lint`
3. ✅ **Code formatted:** `make fmt`
4. ✅ **Complexity check:** `gocyclo -over 10 pkg/agent/`
5. ✅ **Manual TUI testing:** Basic functionality verification

### Regression Prevention

- **Golden tests:** Capture expected behavior before refactoring
- **Integration tests:** Verify end-to-end flows
- **Manual testing:** Test TUI interactions

---

## Risk Assessment

### High Risk Areas

1. **`executeIteration()` refactoring**
   - Risk: Breaking agent loop behavior
   - Mitigation: Extract one helper at a time, test after each extraction

2. **Error handling changes**
   - Risk: Missing error cases or changing error recovery
   - Mitigation: Comprehensive error path testing

3. **Context management**
   - Risk: Breaking summarization or token tracking
   - Mitigation: Verify token counts before/after

### Low Risk Areas

1. **Overlay consolidation** (isolated changes)
2. **Logging improvements** (additive changes)
3. **Helper method extraction** (when done carefully)

---

## Success Criteria

### Phase 2 Complete When:

- ✅ All functions in `pkg/agent/default.go` have <10 cyclomatic complexity
- ✅ File is <700 lines (down from 898)
- ✅ All tests passing
- ✅ No linter errors
- ✅ Manual TUI testing successful
- ✅ Error handling standardized
- ✅ Logging structured and configurable

### Metrics to Track

- Lines of code per file
- Function count
- Cyclomatic complexity
- Test coverage
- Build time
- Memory usage (if impacted)

---

## Timeline

**Week 2-3 Breakdown:**

- **Days 1-2:** Task 2.1 (Simplify Agent Loop) - 8 hours
- **Days 3-4:** Task 2.2 (Standardize Errors) - 6 hours
- **Day 5:** Task 2.3 (Consolidate Overlays) - 4 hours
- **Days 6-7:** Task 2.4 (Structured Logging) - 6 hours

**Total:** 24 hours estimated

---

## Open Questions

1. **Should we introduce a service layer pattern?**
   - Pro: Better separation of concerns
   - Con: Adds abstraction complexity
   - Decision: Defer to Phase 3

2. **Should we split `default.go` into multiple files?**
   - Pro: Smaller files, better organization
   - Con: May fragment related code
   - Decision: Evaluate after Task 2.1 completion

3. **Should we add metrics/telemetry?**
   - Pro: Better observability
   - Con: Scope creep for Phase 2
   - Decision: Defer to Phase 4

---

## Next Steps

1. **Review this plan** with team
2. **Get approval** for approach
3. **Start Task 2.1** branch creation
4. **Begin extraction** of `executeIteration()` helpers

**Ready to proceed?** ✅
