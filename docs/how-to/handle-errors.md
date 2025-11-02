# How to Handle Errors

Step-by-step guide to error handling in Forge agents.

## Overview

Proper error handling ensures agents are robust, recoverable, and user-friendly. This guide covers error types, recovery strategies, and best practices.

**Time to complete:** 15 minutes

**What you'll learn:**
- Handle different error types
- Implement retry logic
- Use circuit breakers
- Recover from failures
- Log and monitor errors

---

## Prerequisites

- Basic understanding of Go error handling
- Familiarity with Forge agents
- Knowledge of context package

---

## Step 1: Understanding Error Types

### Provider Errors

Errors from LLM API calls:

```go
import (
    "errors"
    "github.com/yourusername/forge/pkg/provider"
)

response, err := provider.Complete(ctx, messages)
if err != nil {
    var providerErr *provider.ProviderError
    if errors.As(err, &providerErr) {
        fmt.Printf("Provider error: %d - %s\n", 
            providerErr.StatusCode, 
            providerErr.Message)
    }
}
```

**Common status codes:**
- `401`: Invalid API key
- `429`: Rate limited
- `500`: Server error
- `503`: Service unavailable

### Tool Errors

Errors during tool execution:

```go
result, err := tool.Execute(ctx, args)
if err != nil {
    var toolErr *tool.ToolError
    if errors.As(err, &toolErr) {
        fmt.Printf("Tool %s failed: %v\n", 
            toolErr.ToolName, 
            toolErr.Cause)
    }
}
```

### Context Errors

Timeout or cancellation:

```go
err := agent.Run(ctx, executor)
if err != nil {
    switch {
    case errors.Is(err, context.DeadlineExceeded):
        fmt.Println("Agent timed out")
    case errors.Is(err, context.Canceled):
        fmt.Println("Agent was canceled")
    }
}
```

---

## Step 2: Basic Error Handling

### Check and Log

```go
response, err := provider.Complete(ctx, messages)
if err != nil {
    log.Printf("Provider error: %v", err)
    return err
}
```

### Wrap Errors

Provide context:

```go
response, err := provider.Complete(ctx, messages)
if err != nil {
    return fmt.Errorf("failed to get LLM response: %w", err)
}
```

### Handle Specific Errors

```go
response, err := provider.Complete(ctx, messages)
if err != nil {
    var providerErr *provider.ProviderError
    if errors.As(err, &providerErr) {
        switch providerErr.StatusCode {
        case 401:
            return fmt.Errorf("invalid API key - check OPENAI_API_KEY")
        case 429:
            return fmt.Errorf("rate limited - please wait")
        case 500, 503:
            return fmt.Errorf("service unavailable - try again later")
        default:
            return fmt.Errorf("API error %d: %s", 
                providerErr.StatusCode, 
                providerErr.Message)
        }
    }
    return err
}
```

---

## Step 3: Implement Retry Logic

### Simple Retry

```go
func callWithRetry(fn func() error, maxRetries int) error {
    for attempt := 0; attempt < maxRetries; attempt++ {
        err := fn()
        if err == nil {
            return nil // Success
        }
        
        log.Printf("Attempt %d failed: %v", attempt+1, err)
        
        if attempt < maxRetries-1 {
            time.Sleep(time.Second)
        }
    }
    
    return fmt.Errorf("failed after %d attempts", maxRetries)
}

// Usage
err := callWithRetry(func() error {
    _, err := provider.Complete(ctx, messages)
    return err
}, 3)
```

### Exponential Backoff

```go
func callWithBackoff(fn func() error, maxRetries int) error {
    for attempt := 0; attempt < maxRetries; attempt++ {
        err := fn()
        if err == nil {
            return nil
        }
        
        if attempt < maxRetries-1 {
            // Exponential backoff: 1s, 2s, 4s, 8s...
            waitTime := time.Second * time.Duration(1<<uint(attempt))
            log.Printf("Attempt %d failed, waiting %v: %v", 
                attempt+1, waitTime, err)
            time.Sleep(waitTime)
        }
    }
    
    return fmt.Errorf("failed after %d attempts", maxRetries)
}
```

### Jittered Backoff

Prevent thundering herd:

```go
import "math/rand"

func jitteredBackoff(attempt int) time.Duration {
    baseWait := time.Second * time.Duration(1<<uint(attempt))
    jitter := time.Duration(rand.Int63n(int64(baseWait / 2)))
    return baseWait + jitter
}

func callWithJitteredBackoff(fn func() error, maxRetries int) error {
    for attempt := 0; attempt < maxRetries; attempt++ {
        err := fn()
        if err == nil {
            return nil
        }
        
        if attempt < maxRetries-1 {
            waitTime := jitteredBackoff(attempt)
            log.Printf("Attempt %d failed, waiting %v", attempt+1, waitTime)
            time.Sleep(waitTime)
        }
    }
    
    return fmt.Errorf("failed after %d attempts", maxRetries)
}
```

### Conditional Retry

Only retry transient errors:

```go
func isRetryable(err error) bool {
    var providerErr *provider.ProviderError
    if errors.As(err, &providerErr) {
        // Retry on rate limits and server errors
        return providerErr.StatusCode == 429 ||
               providerErr.StatusCode >= 500
    }
    
    // Retry on network errors
    var netErr net.Error
    if errors.As(err, &netErr) && netErr.Temporary() {
        return true
    }
    
    return false
}

func callWithSmartRetry(fn func() error, maxRetries int) error {
    for attempt := 0; attempt < maxRetries; attempt++ {
        err := fn()
        if err == nil {
            return nil
        }
        
        if !isRetryable(err) {
            return err // Don't retry non-retryable errors
        }
        
        if attempt < maxRetries-1 {
            waitTime := jitteredBackoff(attempt)
            time.Sleep(waitTime)
        }
    }
    
    return fmt.Errorf("failed after %d retries", maxRetries)
}
```

---

## Step 4: Use Circuit Breakers

### Simple Circuit Breaker

```go
type CircuitBreaker struct {
    maxFailures     int
    resetTimeout    time.Duration
    consecutiveFails int
    lastFailTime    time.Time
    state           string // "closed", "open", "half-open"
    mu              sync.Mutex
}

func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        maxFailures:  maxFailures,
        resetTimeout: resetTimeout,
        state:        "closed",
    }
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.Lock()
    
    // Check if circuit should reset
    if cb.state == "open" && time.Since(cb.lastFailTime) > cb.resetTimeout {
        cb.state = "half-open"
        cb.consecutiveFails = 0
    }
    
    // Fail fast if circuit is open
    if cb.state == "open" {
        cb.mu.Unlock()
        return fmt.Errorf("circuit breaker is open")
    }
    
    cb.mu.Unlock()
    
    // Execute function
    err := fn()
    
    cb.mu.Lock()
    defer cb.mu.Unlock()
    
    if err != nil {
        cb.consecutiveFails++
        cb.lastFailTime = time.Now()
        
        if cb.consecutiveFails >= cb.maxFailures {
            cb.state = "open"
            log.Printf("Circuit breaker opened after %d failures", cb.consecutiveFails)
        }
        
        return err
    }
    
    // Success - close circuit
    cb.consecutiveFails = 0
    cb.state = "closed"
    return nil
}
```

### Usage

```go
breaker := NewCircuitBreaker(3, 30*time.Second)

err := breaker.Call(func() error {
    _, err := provider.Complete(ctx, messages)
    return err
})

if err != nil {
    if err.Error() == "circuit breaker is open" {
        log.Println("Service is down, trying again later")
    }
}
```

---

## Step 5: Handle Tool Errors

### Tool Error Recovery

Forge agents automatically add tool errors to conversation:

```go
// Agent tries tool
mem.Add(message.Assistant("<tool>...</tool>"))

// Tool fails
mem.Add(message.Tool("calculator", "Error: division by zero"))

// Agent sees error and adapts
mem.Add(message.Assistant("[I see the error, let me explain instead]"))
```

### Validate Tool Arguments

```go
func (t *MyTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
    // Validate argument exists
    value, ok := args["value"]
    if !ok {
        return "", fmt.Errorf("missing required argument: value")
    }
    
    // Validate type
    strValue, ok := value.(string)
    if !ok {
        return "", fmt.Errorf("value must be string, got %T", value)
    }
    
    // Validate content
    if strValue == "" {
        return "", fmt.Errorf("value cannot be empty")
    }
    
    // Proceed with execution
    return t.process(strValue)
}
```

### Graceful Degradation

```go
func (t *MyTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
    // Try primary method
    result, err := t.primaryMethod(ctx, args)
    if err == nil {
        return result, nil
    }
    
    log.Printf("Primary method failed: %v, trying fallback", err)
    
    // Try fallback method
    result, err = t.fallbackMethod(ctx, args)
    if err == nil {
        return result + " (via fallback)", nil
    }
    
    // Both failed
    return "", fmt.Errorf("all methods failed: %w", err)
}
```

---

## Step 6: Handle Context Errors

### Set Timeouts

```go
// Overall timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

err := agent.Run(ctx, executor)
if errors.Is(err, context.DeadlineExceeded) {
    log.Println("Agent timed out after 5 minutes")
}
```

### Handle Cancellation

```go
ctx, cancel := context.WithCancel(context.Background())

// Cancel on interrupt
go func() {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt)
    <-sigChan
    log.Println("Interrupt received, canceling...")
    cancel()
}()

err := agent.Run(ctx, executor)
if errors.Is(err, context.Canceled) {
    log.Println("Agent was canceled by user")
}
```

### Nested Timeouts

```go
// Overall timeout: 10 minutes
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
defer cancel()

// Provider has its own timeout: 60 seconds per call
provider := openai.NewProvider(
    "gpt-4",
    apiKey,
    openai.WithTimeout(60*time.Second),
)
```

---

## Step 7: Log and Monitor Errors

### Structured Logging

```go
import "log/slog"

logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

response, err := provider.Complete(ctx, messages)
if err != nil {
    logger.Error("provider call failed",
        "error", err,
        "model", "gpt-4",
        "messages", len(messages),
    )
}
```

### Error Metrics

```go
type ErrorMetrics struct {
    totalErrors    int64
    errorsByType   map[string]int64
    mu             sync.Mutex
}

func (m *ErrorMetrics) Record(err error) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    m.totalErrors++
    
    errType := fmt.Sprintf("%T", err)
    m.errorsByType[errType]++
}

func (m *ErrorMetrics) Report() {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    log.Printf("Total errors: %d", m.totalErrors)
    for errType, count := range m.errorsByType {
        log.Printf("  %s: %d", errType, count)
    }
}
```

### Error Tracking

```go
type ErrorTracker struct {
    errors []ErrorEvent
    mu     sync.Mutex
}

type ErrorEvent struct {
    Timestamp time.Time
    Error     error
    Context   map[string]interface{}
}

func (t *ErrorTracker) Track(err error, context map[string]interface{}) {
    t.mu.Lock()
    defer t.mu.Unlock()
    
    t.errors = append(t.errors, ErrorEvent{
        Timestamp: time.Now(),
        Error:     err,
        Context:   context,
    })
}

func (t *ErrorTracker) RecentErrors(n int) []ErrorEvent {
    t.mu.Lock()
    defer t.mu.Unlock()
    
    if len(t.errors) < n {
        n = len(t.errors)
    }
    
    return t.errors[len(t.errors)-n:]
}
```

---

## Best Practices

### 1. Always Check Errors

```go
// ✅ Good: Check every error
response, err := provider.Complete(ctx, messages)
if err != nil {
    return fmt.Errorf("provider failed: %w", err)
}

// ❌ Bad: Ignore errors
response, _ := provider.Complete(ctx, messages)
```

### 2. Provide Context

```go
// ✅ Good: Wrap with context
if err != nil {
    return fmt.Errorf("failed to process user request %s: %w", requestID, err)
}

// ❌ Bad: Return raw error
if err != nil {
    return err
}
```

### 3. Use Appropriate Retries

```go
// ✅ Good: Retry transient errors
if isRetryable(err) {
    time.Sleep(backoff)
    // Retry
}

// ❌ Bad: Retry everything
time.Sleep(backoff)
// Retry even fatal errors
```

### 4. Fail Fast When Appropriate

```go
// ✅ Good: Fast failure for invalid input
if apiKey == "" {
    return fmt.Errorf("API key is required")
}

// ❌ Bad: Let it fail later
// Attempt API call with empty key
```

### 5. Log Errors Appropriately

```go
// ✅ Good: Log with context
log.Printf("Provider error: %v (model: %s, attempt: %d)", 
    err, model, attempt)

// ❌ Bad: Generic logging
log.Println(err)
```

---

## Complete Example

Full error handling implementation:

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "log"
    "math/rand"
    "os"
    "time"
    
    "github.com/yourusername/forge/pkg/provider/openai"
    "github.com/yourusername/forge/pkg/message"
)

func main() {
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        log.Fatal("OPENAI_API_KEY not set")
    }
    
    provider := openai.NewProvider("gpt-4", apiKey)
    
    messages := []message.Message{
        message.User("Hello"),
    }
    
    // Call with full error handling
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()
    
    response, err := callProviderWithRetry(ctx, provider, messages)
    if err != nil {
        log.Fatalf("Failed to get response: %v", err)
    }
    
    fmt.Printf("Response: %s\n", response.Content)
}

func callProviderWithRetry(
    ctx context.Context,
    provider *openai.Provider,
    messages []message.Message,
) (*Response, error) {
    maxRetries := 3
    breaker := NewCircuitBreaker(3, 30*time.Second)
    
    for attempt := 0; attempt < maxRetries; attempt++ {
        // Check context
        if err := ctx.Err(); err != nil {
            return nil, fmt.Errorf("context error: %w", err)
        }
        
        // Call through circuit breaker
        var response *Response
        var err error
        
        breakerErr := breaker.Call(func() error {
            response, err = provider.Complete(ctx, messages)
            return err
        })
        
        if breakerErr != nil {
            if breakerErr.Error() == "circuit breaker is open" {
                return nil, fmt.Errorf("service unavailable: circuit breaker open")
            }
        }
        
        if err == nil {
            return response, nil
        }
        
        // Handle error
        if !isRetryable(err) {
            return nil, fmt.Errorf("non-retryable error: %w", err)
        }
        
        log.Printf("Attempt %d failed: %v", attempt+1, err)
        
        if attempt < maxRetries-1 {
            waitTime := jitteredBackoff(attempt)
            log.Printf("Waiting %v before retry", waitTime)
            
            select {
            case <-time.After(waitTime):
                // Continue to retry
            case <-ctx.Done():
                return nil, ctx.Err()
            }
        }
    }
    
    return nil, fmt.Errorf("failed after %d attempts", maxRetries)
}

func isRetryable(err error) bool {
    var providerErr *provider.ProviderError
    if errors.As(err, &providerErr) {
        return providerErr.StatusCode == 429 || providerErr.StatusCode >= 500
    }
    return false
}

func jitteredBackoff(attempt int) time.Duration {
    baseWait := time.Second * time.Duration(1<<uint(attempt))
    jitter := time.Duration(rand.Int63n(int64(baseWait / 2)))
    return baseWait + jitter
}

type CircuitBreaker struct {
    maxFailures      int
    resetTimeout     time.Duration
    consecutiveFails int
    lastFailTime     time.Time
    state            string
    mu               sync.Mutex
}

func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        maxFailures:  maxFailures,
        resetTimeout: resetTimeout,
        state:        "closed",
    }
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.Lock()
    if cb.state == "open" && time.Since(cb.lastFailTime) > cb.resetTimeout {
        cb.state = "half-open"
        cb.consecutiveFails = 0
    }
    if cb.state == "open" {
        cb.mu.Unlock()
        return fmt.Errorf("circuit breaker is open")
    }
    cb.mu.Unlock()
    
    err := fn()
    
    cb.mu.Lock()
    defer cb.mu.Unlock()
    
    if err != nil {
        cb.consecutiveFails++
        cb.lastFailTime = time.Now()
        if cb.consecutiveFails >= cb.maxFailures {
            cb.state = "open"
        }
        return err
    }
    
    cb.consecutiveFails = 0
    cb.state = "closed"
    return nil
}
```

---

## Next Steps

- Read [Error Handling Reference](../reference/error-handling.md) for detailed patterns
- See [Testing Guide](test-tools.md) for testing error scenarios
- Learn [Performance Optimization](../reference/performance.md)
- Check [API Reference](../reference/api-reference.md) for error types

You now have robust error handling for production agents!