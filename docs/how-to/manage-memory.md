# How to Manage Memory

Step-by-step guide to managing conversation memory in Forge agents.

## Overview

Memory management controls how agents store and retrieve conversation history. This guide covers memory configuration, pruning strategies, and optimization.

**Time to complete:** 10 minutes

**What you'll learn:**
- Configure conversation memory
- Implement pruning strategies
- Optimize memory usage
- Handle long conversations
- Clear and reset memory

---

## Prerequisites

- Basic understanding of Forge agents
- Familiarity with token limits
- Knowledge of your LLM's context window

---

## Step 1: Create Basic Memory

### Default Configuration

```go
import "github.com/yourusername/forge/pkg/memory"

// Create memory with 8000 token limit
mem := memory.NewConversationMemory(8000)
```

### Match Model Context Window

```go
// GPT-4 (8K context) - leave room for response
mem := memory.NewConversationMemory(6000)

// GPT-4-32K - larger context
mem := memory.NewConversationMemory(28000)

// GPT-3.5-Turbo (4K context)
mem := memory.NewConversationMemory(3000)
```

**Rule of thumb:** Set memory limit to 70-80% of model's context window.

---

## Step 2: Add Messages to Memory

### Adding Different Message Types

```go
import "github.com/yourusername/forge/pkg/message"

// System message (defines behavior)
mem.Add(message.System("You are a helpful coding assistant"))

// User message
mem.Add(message.User("How do I create a function in Python?"))

// Assistant message
mem.Add(message.Assistant("Here's how to create a function in Python..."))

// Tool result
mem.Add(message.Tool("calculator", "Result: 42"))
```

### Adding Messages in Sequence

```go
// Typical conversation flow
mem.Add(message.System("You are a helpful assistant"))
mem.Add(message.User("What's the weather?"))
mem.Add(message.Assistant("[I'll check the weather]\n<tool>...</tool>"))
mem.Add(message.Tool("weather", "Sunny, 72°F"))
mem.Add(message.Assistant("The weather is sunny and 72°F"))
```

---

## Step 3: Retrieve Messages

### Get All Messages

```go
messages := mem.GetMessages()
fmt.Printf("Total messages: %d\n", len(messages))

for i, msg := range messages {
    fmt.Printf("%d. %s: %s\n", i+1, msg.Role, msg.Content)
}
```

### Get Recent Messages

```go
messages := mem.GetMessages()

// Get last N messages
n := 10
if len(messages) > n {
    recentMessages := messages[len(messages)-n:]
    // Use recentMessages
}
```

### Filter by Role

```go
messages := mem.GetMessages()

var userMessages []message.Message
for _, msg := range messages {
    if msg.Role == "user" {
        userMessages = append(userMessages, msg)
    }
}
```

---

## Step 4: Configure Pruning

### Automatic Pruning

Memory automatically prunes when token limit is exceeded:

```go
mem := memory.NewConversationMemory(5000)

// Add messages until limit is reached
for i := 0; i < 100; i++ {
    mem.Add(message.User(fmt.Sprintf("Message %d", i)))
}

// Oldest messages are automatically removed
messages := mem.GetMessages()
fmt.Printf("Messages after pruning: %d\n", len(messages))
```

### Pruning Strategy

```go
// Oldest-first (default)
mem := memory.NewConversationMemory(
    8000,
    memory.WithPruningStrategy(memory.OldestFirst),
)

// System message is always preserved
```

### Manual Pruning

```go
// Prune to specific token count
err := mem.Prune(5000)
if err != nil {
    log.Printf("Pruning failed: %v", err)
}
```

---

## Step 5: Check Memory Status

### Estimate Token Count

```go
// Get approximate token count
tokens := mem.EstimateTokens()
fmt.Printf("Estimated tokens: %d\n", tokens)

// Check if approaching limit
maxTokens := 8000
if tokens > maxTokens*0.9 {
    fmt.Println("Warning: Memory nearly full")
}
```

### Get Message Count

```go
count := mem.Size()
fmt.Printf("Total messages: %d\n", count)
```

### Check Individual Message Size

```go
func estimateMessageTokens(msg message.Message) int {
    // Rough estimate: 1 token ≈ 4 characters
    tokens := len(msg.Content) / 4
    
    // Add overhead for role
    tokens += 4
    
    return tokens
}
```

---

## Step 6: Clear Memory

### Clear All Messages

```go
// Remove all messages
err := mem.Clear()
if err != nil {
    log.Printf("Clear failed: %v", err)
}

fmt.Printf("Messages after clear: %d\n", mem.Size())
```

### Clear Except System Message

```go
messages := mem.GetMessages()

// Find and preserve system message
var systemMsg message.Message
hasSystem := false
if len(messages) > 0 && messages[0].Role == "system" {
    systemMsg = messages[0]
    hasSystem = true
}

// Clear all
mem.Clear()

// Restore system message
if hasSystem {
    mem.Add(systemMsg)
}
```

### Selective Clearing

```go
// Clear messages older than timestamp
cutoffTime := time.Now().Add(-1 * time.Hour)

// Get all messages
messages := mem.GetMessages()

// Clear and re-add recent ones
mem.Clear()

for _, msg := range messages {
    if msg.Timestamp.After(cutoffTime) {
        mem.Add(msg)
    }
}
```

---

## Advanced Usage

### Thread-Safe Access

```go
import "sync"

type ThreadSafeMemory struct {
    mem memory.Memory
    mu  sync.RWMutex
}

func (t *ThreadSafeMemory) Add(msg message.Message) error {
    t.mu.Lock()
    defer t.mu.Unlock()
    return t.mem.Add(msg)
}

func (t *ThreadSafeMemory) GetMessages() []message.Message {
    t.mu.RLock()
    defer t.mu.RUnlock()
    return t.mem.GetMessages()
}
```

### Memory with Persistence

```go
type PersistentMemory struct {
    *memory.ConversationMemory
    filePath string
}

func (p *PersistentMemory) Add(msg message.Message) error {
    if err := p.ConversationMemory.Add(msg); err != nil {
        return err
    }
    return p.saveToFile()
}

func (p *PersistentMemory) saveToFile() error {
    messages := p.GetMessages()
    data, err := json.Marshal(messages)
    if err != nil {
        return err
    }
    return os.WriteFile(p.filePath, data, 0644)
}

func (p *PersistentMemory) loadFromFile() error {
    data, err := os.ReadFile(p.filePath)
    if err != nil {
        return err
    }
    
    var messages []message.Message
    if err := json.Unmarshal(data, &messages); err != nil {
        return err
    }
    
    for _, msg := range messages {
        p.ConversationMemory.Add(msg)
    }
    
    return nil
}
```

### Memory with Sliding Window

```go
type SlidingWindowMemory struct {
    mem        memory.Memory
    windowSize int
}

func (s *SlidingWindowMemory) Add(msg message.Message) error {
    if err := s.mem.Add(msg); err != nil {
        return err
    }
    
    messages := s.mem.GetMessages()
    
    // Keep only last N messages (plus system message)
    if len(messages) > s.windowSize+1 {
        s.mem.Clear()
        
        // Re-add system message
        if messages[0].Role == "system" {
            s.mem.Add(messages[0])
        }
        
        // Add last N messages
        start := len(messages) - s.windowSize
        for _, msg := range messages[start:] {
            s.mem.Add(msg)
        }
    }
    
    return nil
}
```

---

## Memory Patterns

### Conversation Session

```go
type ConversationSession struct {
    memory    memory.Memory
    startTime time.Time
}

func NewSession(maxTokens int) *ConversationSession {
    return &ConversationSession{
        memory:    memory.NewConversationMemory(maxTokens),
        startTime: time.Now(),
    }
}

func (s *ConversationSession) Duration() time.Duration {
    return time.Since(s.startTime)
}

func (s *ConversationSession) MessageCount() int {
    return s.memory.Size()
}

func (s *ConversationSession) Reset() {
    s.memory.Clear()
    s.startTime = time.Now()
}
```

### Memory with Summary

```go
type SummarizingMemory struct {
    mem            memory.Memory
    summaryThreshold int
}

func (s *SummarizingMemory) Add(msg message.Message) error {
    if err := s.mem.Add(msg); err != nil {
        return err
    }
    
    // When memory is full, summarize old messages
    if s.mem.Size() > s.summaryThreshold {
        s.summarizeOldMessages()
    }
    
    return nil
}

func (s *SummarizingMemory) summarizeOldMessages() {
    messages := s.mem.GetMessages()
    
    // Get messages to summarize (middle portion)
    start := 1 // After system message
    end := len(messages) / 2
    
    toSummarize := messages[start:end]
    
    // Generate summary (using LLM)
    summary := generateSummary(toSummarize)
    
    // Replace old messages with summary
    s.mem.Clear()
    s.mem.Add(messages[0]) // System message
    s.mem.Add(message.System(fmt.Sprintf("Previous conversation summary: %s", summary)))
    
    // Keep recent messages
    for _, msg := range messages[end:] {
        s.mem.Add(msg)
    }
}
```

---

## Monitoring Memory

### Memory Statistics

```go
type MemoryStats struct {
    TotalMessages int
    EstimatedTokens int
    MaxTokens int
    UtilizationPercent float64
}

func getMemoryStats(mem memory.Memory, maxTokens int) MemoryStats {
    messages := mem.GetMessages()
    tokens := mem.EstimateTokens()
    
    return MemoryStats{
        TotalMessages: len(messages),
        EstimatedTokens: tokens,
        MaxTokens: maxTokens,
        UtilizationPercent: float64(tokens) / float64(maxTokens) * 100,
    }
}

// Usage
stats := getMemoryStats(mem, 8000)
fmt.Printf("Memory: %d messages, %d/%d tokens (%.1f%%)\n",
    stats.TotalMessages,
    stats.EstimatedTokens,
    stats.MaxTokens,
    stats.UtilizationPercent,
)
```

### Memory Alerts

```go
func checkMemoryHealth(mem memory.Memory, maxTokens int) []string {
    var alerts []string
    
    tokens := mem.EstimateTokens()
    utilization := float64(tokens) / float64(maxTokens)
    
    if utilization > 0.95 {
        alerts = append(alerts, "Memory nearly full (>95%)")
    } else if utilization > 0.8 {
        alerts = append(alerts, "Memory usage high (>80%)")
    }
    
    messageCount := mem.Size()
    if messageCount > 100 {
        alerts = append(alerts, fmt.Sprintf("High message count: %d", messageCount))
    }
    
    return alerts
}
```

---

## Best Practices

### 1. Match Memory to Model

```go
// ✅ Good: Appropriate size for model
// GPT-4 8K → 6K memory
mem := memory.NewConversationMemory(6000)

// ❌ Bad: Exceeds model context
// GPT-4 8K → 10K memory
mem := memory.NewConversationMemory(10000)
```

### 2. Preserve System Message

```go
// ✅ Good: System message always preserved
mem := memory.NewConversationMemory(8000)
// Automatic preservation

// ❌ Bad: Manual preservation is error-prone
```

### 3. Monitor Memory Usage

```go
// ✅ Good: Regular monitoring
if mem.EstimateTokens() > maxTokens*0.9 {
    log.Println("Memory nearly full")
}

// ❌ Bad: No monitoring
// Memory could exceed limit unexpectedly
```

### 4. Clear Between Sessions

```go
// ✅ Good: Clear between users/sessions
if newSession {
    mem.Clear()
}

// ❌ Bad: Reuse memory across sessions
// Previous context leaks to new user
```

### 5. Use Appropriate Limits

```go
// ✅ Good: Conservative limit
mem := memory.NewConversationMemory(6000) // For 8K model

// ❌ Bad: Use full context
mem := memory.NewConversationMemory(8000) // No room for response
```

---

## Troubleshooting

### Memory Fills Too Quickly

**Problem:** Memory reaches limit after few messages

**Solutions:**
```go
// 1. Increase limit (if model allows)
mem := memory.NewConversationMemory(12000)

// 2. Use more aggressive pruning
mem.Prune(4000) // Prune to lower threshold

// 3. Shorten system prompt
systemPrompt := "You are a helpful assistant" // Concise
```

### Context Length Errors

**Problem:** LLM returns context length error

**Solutions:**
```go
// 1. Reduce memory limit
mem := memory.NewConversationMemory(5000) // Lower than before

// 2. Manual pruning before each request
mem.Prune(4000)

// 3. Use smaller model with larger context
provider := openai.NewProvider("gpt-4-32k", apiKey)
```

### Important Context Lost

**Problem:** Pruning removes needed information

**Solutions:**
```go
// 1. Implement importance-based pruning
// Keep system message + important messages

// 2. Use summarization
// Summarize old messages instead of deleting

// 3. Increase memory limit
mem := memory.NewConversationMemory(12000)
```

---

## Complete Example

Full memory management implementation:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/yourusername/forge/pkg/memory"
    "github.com/yourusername/forge/pkg/message"
)

func main() {
    // Create memory with monitoring
    maxTokens := 8000
    mem := memory.NewConversationMemory(maxTokens)
    
    // Add system message
    mem.Add(message.System("You are a helpful assistant"))
    
    // Simulate conversation
    for i := 1; i <= 50; i++ {
        // User message
        mem.Add(message.User(fmt.Sprintf("Question %d", i)))
        
        // Assistant response
        mem.Add(message.Assistant(fmt.Sprintf("Answer %d", i)))
        
        // Check memory status
        stats := getMemoryStats(mem, maxTokens)
        fmt.Printf("Turn %d: %d messages, %.1f%% full\n",
            i, stats.TotalMessages, stats.UtilizationPercent)
        
        // Alert if high
        if alerts := checkMemoryHealth(mem, maxTokens); len(alerts) > 0 {
            for _, alert := range alerts {
                log.Printf("Alert: %s", alert)
            }
        }
        
        // Manual pruning if needed
        if stats.UtilizationPercent > 90 {
            mem.Prune(maxTokens * 70 / 100) // Prune to 70%
            fmt.Println("Memory pruned to 70%")
        }
    }
    
    // Final stats
    finalStats := getMemoryStats(mem, maxTokens)
    fmt.Printf("\nFinal: %d messages, %d tokens (%.1f%%)\n",
        finalStats.TotalMessages,
        finalStats.EstimatedTokens,
        finalStats.UtilizationPercent,
    )
}

type MemoryStats struct {
    TotalMessages      int
    EstimatedTokens    int
    MaxTokens          int
    UtilizationPercent float64
}

func getMemoryStats(mem memory.Memory, maxTokens int) MemoryStats {
    messages := mem.GetMessages()
    tokens := mem.EstimateTokens()
    
    return MemoryStats{
        TotalMessages:      len(messages),
        EstimatedTokens:    tokens,
        MaxTokens:          maxTokens,
        UtilizationPercent: float64(tokens) / float64(maxTokens) * 100,
    }
}

func checkMemoryHealth(mem memory.Memory, maxTokens int) []string {
    var alerts []string
    tokens := mem.EstimateTokens()
    utilization := float64(tokens) / float64(maxTokens)
    
    if utilization > 0.95 {
        alerts = append(alerts, "Memory nearly full")
    } else if utilization > 0.8 {
        alerts = append(alerts, "Memory usage high")
    }
    
    return alerts
}
```

---

## Next Steps

- Read [Memory System Architecture](../architecture/memory-system.md) for design details
- See [Configuration Reference](../reference/configuration.md) for all options
- Learn [Performance Optimization](../reference/performance.md) for memory tuning
- Check [API Reference](../reference/api-reference.md) for Memory interface

You now have complete control over conversation memory!