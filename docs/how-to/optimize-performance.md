# How to Optimize Performance

Step-by-step guide to optimizing Forge agent performance for speed, cost, and efficiency.

## Overview

Performance optimization reduces latency, token usage, and API costs while maintaining agent quality. This guide covers practical optimization techniques.

**Time to complete:** 20 minutes

**What you'll learn:**
- Reduce token usage
- Minimize latency
- Optimize costs
- Monitor performance
- Benchmark improvements

---

## Prerequisites

- Running Forge agent
- Basic understanding of token costs
- Familiarity with your LLM's pricing

---

## Step 1: Baseline Performance

### Measure Current Performance

```go
type PerformanceMetrics struct {
    RequestCount    int64
    TotalLatency    time.Duration
    TotalTokens     int64
    TotalCost       float64
    mu              sync.Mutex
}

func (p *PerformanceMetrics) Record(latency time.Duration, tokens int, cost float64) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    p.RequestCount++
    p.TotalLatency += latency
    p.TotalTokens += int64(tokens)
    p.TotalCost += cost
}

func (p *PerformanceMetrics) Report() {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    avgLatency := p.TotalLatency / time.Duration(p.RequestCount)
    avgTokens := p.TotalTokens / p.RequestCount
    avgCost := p.TotalCost / float64(p.RequestCount)
    
    fmt.Printf("Performance Report:\n")
    fmt.Printf("  Requests: %d\n", p.RequestCount)
    fmt.Printf("  Avg Latency: %v\n", avgLatency)
    fmt.Printf("  Avg Tokens: %d\n", avgTokens)
    fmt.Printf("  Avg Cost: $%.4f\n", avgCost)
    fmt.Printf("  Total Cost: $%.2f\n", p.TotalCost)
}
```

### Collect Baseline

```go
metrics := &PerformanceMetrics{}

start := time.Now()
response, err := provider.Complete(ctx, messages)
latency := time.Since(start)

if err == nil {
    tokens := estimateTokens(messages, response)
    cost := calculateCost(tokens, "gpt-4")
    metrics.Record(latency, tokens, cost)
}

metrics.Report()
```

---

## Step 2: Optimize Token Usage

### Reduce Context Size

**Before:**
```go
// Using full 8K context
mem := memory.NewConversationMemory(8000)
```

**After:**
```go
// Right-sized for task
mem := memory.NewConversationMemory(4000)
```

**Impact:** ~50% reduction in tokens per request

### Shorten System Prompt

**Before:**
```go
systemPrompt := `You are a highly capable Python programming expert assistant.
You should always strive to write clean, maintainable, and well-documented code.
Follow PEP 8 style guidelines strictly. Include comprehensive type hints in all
function signatures. Provide detailed explanations for your code choices. Consider
edge cases and error handling in all implementations. Write unit tests when appropriate.`
// ~100 tokens
```

**After:**
```go
systemPrompt := `You are a Python expert. Write clean, type-hinted code following PEP 8.`
// ~15 tokens
```

**Impact:** 85% reduction in system prompt tokens

### Optimize Tool Descriptions

**Before:**
```go
func (c *Calculator) Description() string {
    return `This tool performs basic arithmetic operations including addition,
    subtraction, multiplication, and division. You can use it whenever you need
    to calculate numerical values. It supports both integers and floating-point
    numbers. The tool will return the result as a formatted string.`
}
// ~50 tokens
```

**After:**
```go
func (c *Calculator) Description() string {
    return "Performs arithmetic: add, subtract, multiply, divide"
}
// ~10 tokens
```

**Impact:** 80% reduction in tool description tokens

---

## Step 3: Reduce Latency

### Use Faster Models

**Before:**
```go
// GPT-4: ~5-10s latency
provider := openai.NewProvider("gpt-4", apiKey)
```

**After:**
```go
// GPT-3.5-Turbo: ~1-2s latency
provider := openai.NewProvider("gpt-3.5-turbo", apiKey)
```

**Impact:** 50-80% latency reduction

### Enable Streaming

**Before:**
```go
// Wait for complete response
response, err := provider.Complete(ctx, messages)
fmt.Println(response.Content)
```

**After:**
```go
// Stream for perceived performance
stream, err := provider.Stream(ctx, messages)
for chunk := range stream {
    if chunk.Error != nil {
        break
    }
    fmt.Print(chunk.Content) // Immediate feedback
}
```

**Impact:** Users see results 80% faster (perceived)

### Reduce Max Iterations

**Before:**
```go
agent, _ := core.NewAgent(
    provider,
    memory,
    tools,
    core.WithMaxIterations(20), // Allows many iterations
)
```

**After:**
```go
agent, _ := core.NewAgent(
    provider,
    memory,
    tools,
    core.WithMaxIterations(5), // Limit for simple tasks
)
```

**Impact:** Up to 75% faster for simple tasks

---

## Step 4: Optimize Costs

### Model Selection by Task

```go
func selectModel(taskComplexity string) string {
    switch taskComplexity {
    case "simple":
        return "gpt-3.5-turbo" // $0.001/1K tokens
    case "medium":
        return "gpt-4-turbo-preview" // $0.01/1K tokens
    case "complex":
        return "gpt-4" // $0.03/1K tokens
    default:
        return "gpt-3.5-turbo"
    }
}

// Usage
model := selectModel("simple")
provider := openai.NewProvider(model, apiKey)
```

**Impact:** 90% cost reduction for simple tasks

### Set Response Limits

**Before:**
```go
// No limit - responses can be very long
provider := openai.NewProvider("gpt-4", apiKey)
```

**After:**
```go
// Limit response length
provider := openai.NewProvider(
    "gpt-4",
    apiKey,
    openai.WithMaxTokens(500), // Max 500 token responses
)
```

**Impact:** 50-70% reduction in output tokens

### Aggressive Memory Pruning

**Before:**
```go
// Large memory window
mem := memory.NewConversationMemory(8000)
```

**After:**
```go
// Smaller memory window
mem := memory.NewConversationMemory(3000)

// Manual pruning when needed
if mem.EstimateTokens() > 2500 {
    mem.Prune(2000)
}
```

**Impact:** 60% reduction in context tokens

---

## Step 5: Monitor Performance

### Track Metrics

```go
type DetailedMetrics struct {
    metrics  map[string]*PerformanceMetrics
    mu       sync.Mutex
}

func (d *DetailedMetrics) Record(category string, latency time.Duration, tokens int, cost float64) {
    d.mu.Lock()
    defer d.mu.Unlock()
    
    if _, ok := d.metrics[category]; !ok {
        d.metrics[category] = &PerformanceMetrics{}
    }
    
    d.metrics[category].Record(latency, tokens, cost)
}

func (d *DetailedMetrics) ReportByCategory() {
    d.mu.Lock()
    defer d.mu.Unlock()
    
    for category, metrics := range d.metrics {
        fmt.Printf("\n%s:\n", category)
        metrics.Report()
    }
}

// Usage
metrics := &DetailedMetrics{
    metrics: make(map[string]*PerformanceMetrics),
}

metrics.Record("simple_questions", latency, tokens, cost)
metrics.Record("complex_tasks", latency, tokens, cost)
metrics.ReportByCategory()
```

### Set Performance Alerts

```go
func checkPerformance(metrics *PerformanceMetrics) []string {
    var alerts []string
    
    avgLatency := metrics.TotalLatency / time.Duration(metrics.RequestCount)
    if avgLatency > 10*time.Second {
        alerts = append(alerts, 
            fmt.Sprintf("High latency: %v", avgLatency))
    }
    
    avgTokens := metrics.TotalTokens / metrics.RequestCount
    if avgTokens > 5000 {
        alerts = append(alerts,
            fmt.Sprintf("High token usage: %d", avgTokens))
    }
    
    if metrics.TotalCost > 100.0 {
        alerts = append(alerts,
            fmt.Sprintf("High total cost: $%.2f", metrics.TotalCost))
    }
    
    return alerts
}

// Check regularly
if alerts := checkPerformance(metrics); len(alerts) > 0 {
    for _, alert := range alerts {
        log.Printf("ALERT: %s", alert)
    }
}
```

---

## Step 6: Benchmark Improvements

### Before/After Comparison

```go
func benchmarkConfiguration(name string, config func() (core.Agent, error)) {
    fmt.Printf("\nBenchmarking: %s\n", name)
    
    agent, _ := config()
    metrics := &PerformanceMetrics{}
    
    // Run 10 test requests
    for i := 0; i < 10; i++ {
        start := time.Now()
        
        // Simulate request
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        err := agent.Run(ctx, mockExecutor)
        cancel()
        
        latency := time.Since(start)
        
        if err == nil {
            tokens := 1000 // Estimate
            cost := 0.03   // Estimate
            metrics.Record(latency, tokens, cost)
        }
    }
    
    metrics.Report()
}

// Test different configurations
benchmarkConfiguration("baseline", func() (core.Agent, error) {
    return core.NewAgent(
        openai.NewProvider("gpt-4", apiKey),
        memory.NewConversationMemory(8000),
        tools,
    )
})

benchmarkConfiguration("optimized", func() (core.Agent, error) {
    return core.NewAgent(
        openai.NewProvider("gpt-3.5-turbo", apiKey,
            openai.WithMaxTokens(500)),
        memory.NewConversationMemory(3000),
        tools,
        core.WithMaxIterations(5),
    )
})
```

---

## Optimization Patterns

### Pattern 1: Tiered Models

Use cheap models for simple tasks, expensive for complex:

```go
type AdaptiveProvider struct {
    cheapProvider     *openai.Provider
    expensiveProvider *openai.Provider
}

func (a *AdaptiveProvider) Complete(ctx context.Context, messages []message.Message) (*Response, error) {
    // Analyze complexity
    complexity := analyzeComplexity(messages)
    
    if complexity == "simple" {
        log.Println("Using GPT-3.5-Turbo")
        return a.cheapProvider.Complete(ctx, messages)
    }
    
    log.Println("Using GPT-4")
    return a.expensiveProvider.Complete(ctx, messages)
}

func analyzeComplexity(messages []message.Message) string {
    // Simple heuristic: check last message length
    lastMsg := messages[len(messages)-1]
    if len(lastMsg.Content) < 100 {
        return "simple"
    }
    return "complex"
}
```

### Pattern 2: Caching

Cache repeated queries:

```go
type CachingProvider struct {
    provider *openai.Provider
    cache    map[string]*Response
    mu       sync.RWMutex
}

func (c *CachingProvider) Complete(ctx context.Context, messages []message.Message) (*Response, error) {
    // Create cache key
    key := generateCacheKey(messages)
    
    // Check cache
    c.mu.RLock()
    if cached, ok := c.cache[key]; ok {
        c.mu.RUnlock()
        log.Println("Cache hit!")
        return cached, nil
    }
    c.mu.RUnlock()
    
    // Call provider
    response, err := c.provider.Complete(ctx, messages)
    if err != nil {
        return nil, err
    }
    
    // Store in cache
    c.mu.Lock()
    c.cache[key] = response
    c.mu.Unlock()
    
    return response, nil
}

func generateCacheKey(messages []message.Message) string {
    h := sha256.New()
    for _, msg := range messages {
        h.Write([]byte(msg.Role + msg.Content))
    }
    return hex.EncodeToString(h.Sum(nil))
}
```

### Pattern 3: Batching

Batch multiple queries:

```go
type BatchProcessor struct {
    provider *openai.Provider
    requests chan BatchRequest
}

type BatchRequest struct {
    Messages []message.Message
    Response chan *Response
    Error    chan error
}

func (b *BatchProcessor) Start() {
    go func() {
        batch := make([]BatchRequest, 0, 10)
        ticker := time.NewTicker(100 * time.Millisecond)
        
        for {
            select {
            case req := <-b.requests:
                batch = append(batch, req)
                
                if len(batch) >= 10 {
                    b.processBatch(batch)
                    batch = batch[:0]
                }
                
            case <-ticker.C:
                if len(batch) > 0 {
                    b.processBatch(batch)
                    batch = batch[:0]
                }
            }
        }
    }()
}

func (b *BatchProcessor) processBatch(batch []BatchRequest) {
    for _, req := range batch {
        go func(r BatchRequest) {
            response, err := b.provider.Complete(context.Background(), r.Messages)
            if err != nil {
                r.Error <- err
            } else {
                r.Response <- response
            }
        }(req)
    }
}
```

---

## Complete Optimization Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "sync"
    "time"
    
    "github.com/yourusername/forge/pkg/core"
    "github.com/yourusername/forge/pkg/provider/openai"
    "github.com/yourusername/forge/pkg/memory"
    "github.com/yourusername/forge/pkg/tool"
)

type OptimizedAgent struct {
    agent   core.Agent
    metrics *PerformanceMetrics
}

func NewOptimizedAgent() (*OptimizedAgent, error) {
    apiKey := os.Getenv("OPENAI_API_KEY")
    
    // Optimized configuration
    provider := openai.NewProvider(
        "gpt-3.5-turbo",  // Fast, cheap model
        apiKey,
        openai.WithTemperature(0.3), // Lower for consistency
        openai.WithMaxTokens(500),   // Limit response length
        openai.WithTimeout(30*time.Second),
    )
    
    // Smaller memory window
    mem := memory.NewConversationMemory(3000)
    
    // Minimal tool set
    tools := []core.Tool{
        tool.NewTaskCompletion(),
    }
    
    // Concise system prompt
    agent, err := core.NewAgent(
        provider,
        mem,
        tools,
        core.WithMaxIterations(5),
        core.WithSystemPrompt("You are a helpful assistant. Be concise."),
    )
    
    if err != nil {
        return nil, err
    }
    
    return &OptimizedAgent{
        agent:   agent,
        metrics: &PerformanceMetrics{},
    }, nil
}

func (o *OptimizedAgent) Run(ctx context.Context, executor core.Executor) error {
    start := time.Now()
    
    err := o.agent.Run(ctx, executor)
    
    latency := time.Since(start)
    tokens := 1000 // Estimate or calculate
    cost := calculateCost(tokens, "gpt-3.5-turbo")
    
    o.metrics.Record(latency, tokens, cost)
    
    return err
}

func (o *OptimizedAgent) GetMetrics() *PerformanceMetrics {
    return o.metrics
}

type PerformanceMetrics struct {
    RequestCount int64
    TotalLatency time.Duration
    TotalTokens  int64
    TotalCost    float64
    mu           sync.Mutex
}

func (p *PerformanceMetrics) Record(latency time.Duration, tokens int, cost float64) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    p.RequestCount++
    p.TotalLatency += latency
    p.TotalTokens += int64(tokens)
    p.TotalCost += cost
}

func (p *PerformanceMetrics) Report() {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    if p.RequestCount == 0 {
        fmt.Println("No requests recorded")
        return
    }
    
    avgLatency := p.TotalLatency / time.Duration(p.RequestCount)
    avgTokens := p.TotalTokens / p.RequestCount
    avgCost := p.TotalCost / float64(p.RequestCount)
    
    fmt.Printf("\nPerformance Metrics:\n")
    fmt.Printf("  Total Requests: %d\n", p.RequestCount)
    fmt.Printf("  Avg Latency: %v\n", avgLatency)
    fmt.Printf("  Avg Tokens: %d\n", avgTokens)
    fmt.Printf("  Avg Cost: $%.4f\n", avgCost)
    fmt.Printf("  Total Cost: $%.2f\n", p.TotalCost)
    
    // Cost projections
    fmt.Printf("\nProjections:\n")
    fmt.Printf("  100 req/day: $%.2f/day = $%.2f/month\n",
        avgCost*100, avgCost*100*30)
    fmt.Printf("  1000 req/day: $%.2f/day = $%.2f/month\n",
        avgCost*1000, avgCost*1000*30)
}

func calculateCost(tokens int, model string) float64 {
    rates := map[string]float64{
        "gpt-3.5-turbo":       0.002 / 1000, // $0.002 per 1K tokens
        "gpt-4":               0.06 / 1000,  // $0.06 per 1K tokens
        "gpt-4-turbo-preview": 0.01 / 1000,  // $0.01 per 1K tokens
    }
    
    rate, ok := rates[model]
    if !ok {
        rate = 0.01 / 1000 // Default
    }
    
    return float64(tokens) * rate
}

func main() {
    agent, err := NewOptimizedAgent()
    if err != nil {
        log.Fatal(err)
    }
    
    // Use agent...
    
    // Report metrics
    agent.GetMetrics().Report()
}
```

---

## Optimization Checklist

- [ ] Baseline performance measured
- [ ] Token usage optimized (prompts, descriptions)
- [ ] Memory size appropriate for task
- [ ] Model matches task complexity
- [ ] Max iterations set appropriately
- [ ] Response limits configured
- [ ] Streaming enabled for long responses
- [ ] Performance monitoring in place
- [ ] Cost tracking implemented
- [ ] Benchmarks show improvement

---

## Next Steps

- Read [Performance Reference](../reference/performance.md) for detailed metrics
- See [Configuration Guide](../reference/configuration.md) for all options
- Learn [Error Handling](handle-errors.md) for robust systems
- Check [Testing Guide](test-tools.md) for benchmarking

You're now equipped to build highly optimized agents!