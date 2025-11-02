# Provider Interface Analysis

## Current Interface Assessment

### Current Design
```go
type Provider interface {
    GenerateStream(ctx context.Context, messages []*types.Message) (<-chan *types.AgentEvent, error)
    Generate(ctx context.Context, messages []*types.Message) (*types.Message, error)
    GetModelInfo() *types.ModelInfo
}
```

## Industry Standard LLM Client Patterns

### OpenAI Go SDK
```go
// Streaming
stream, err := client.CreateChatCompletionStream(ctx, ChatCompletionRequest{
    Model: "gpt-4",
    Messages: []ChatCompletionMessage{...},
    Stream: true,
})
for {
    response, err := stream.Recv()
    // Process delta
}

// Non-streaming
response, err := client.CreateChatCompletion(ctx, ChatCompletionRequest{
    Model: "gpt-4", 
    Messages: []ChatCompletionMessage{...},
})
```

### Anthropic Claude Go SDK
```go
// Streaming
stream := client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
    Model: "claude-3-opus",
    Messages: []anthropic.MessageParam{...},
})
for stream.Next() {
    event := stream.Current()
    // Process event
}

// Non-streaming
message, err := client.Messages.New(ctx, anthropic.MessageNewParams{
    Model: "claude-3-opus",
    Messages: []anthropic.MessageParam{...},
})
```

### Google Gemini (Vertex AI)
```go
// Streaming
iter := model.GenerateContentStream(ctx, genai.Text("..."))
for {
    resp, err := iter.Next()
    if err == iterator.Done { break }
    // Process chunk
}

// Non-streaming
resp, err := model.GenerateContent(ctx, genai.Text("..."))
```

## Problems with Current Design

### 1. **Event Coupling Issue** ❌
```go
GenerateStream(ctx, messages) (<-chan *types.AgentEvent, error)
```

**Problem**: Provider returns `AgentEvent` which couples:
- LLM provider layer → Agent event system
- Provider must know about ALL agent event types
- Provider cannot be used outside agent context

**Why this is wrong**:
- Providers should be **LLM-focused**, not agent-aware
- Events like `ThinkingStart`, `ToolCallStart`, `StatusUpdate` are agent concerns
- Provider should only know about LLM response streaming

### 2. **Abstraction Leakage** ❌
Provider is forced to emit agent-level events when it should only handle LLM responses.

### 3. **Reusability Problem** ❌
Can't use Provider in non-agent contexts (CLI tools, batch processing, etc.)

## Proposed Solution: Separation of Concerns

### Option A: Provider Returns LLM-Specific Stream ✅

```go
// LLM Response types (provider-specific)
type StreamChunk struct {
    Delta   string
    Role    string
    Done    bool
    Error   error
}

type Provider interface {
    // Returns LLM response chunks
    ChatCompletionStream(ctx context.Context, messages []*Message) (<-chan *StreamChunk, error)
    
    // Returns complete response
    ChatCompletion(ctx context.Context, messages []*Message) (*Message, error)
    
    GetModelInfo() *ModelInfo
}
```

**Benefits**:
- Provider only deals with LLM concerns
- Agent wraps provider and converts to AgentEvents
- Provider is reusable in any context
- Cleaner separation of concerns

### Option B: Generic Response Iterator ✅

```go
type CompletionChunk struct {
    Content  string
    Finished bool
}

type Provider interface {
    // Returns iterator-style interface
    StreamCompletion(ctx context.Context, messages []*Message) (CompletionStream, error)
    Complete(ctx context.Context, messages []*Message) (*Message, error)
    GetModelInfo() *ModelInfo
}

type CompletionStream interface {
    Next() bool
    Chunk() *CompletionChunk
    Err() error
    Close() error
}
```

**Benefits**:
- Familiar iterator pattern
- Less goroutine overhead
- Easier error handling
- Still provider-agnostic

## Recommended Approach

### Use Option A with Clear Layering

```
┌─────────────────────────────────────────┐
│           Agent Layer                   │
│  - Emits AgentEvents                    │
│  - Manages thinking/tool events         │
│  - Wraps provider calls                 │
└─────────────────┬───────────────────────┘
                  │
                  │ Uses Provider
                  ↓
┌─────────────────────────────────────────┐
│         Provider Layer                  │
│  - Returns LLM chunks                   │
│  - Handles API calls                    │
│  - Provider-specific logic              │
└─────────────────────────────────────────┘
```

### Updated Interface

```go
package llm

// StreamChunk represents a single chunk from streaming LLM response
type StreamChunk struct {
    Content  string  // Text content delta
    Role     string  // Message role (for first chunk)
    Finished bool    // True when stream is complete
    Error    error   // Any error that occurred
}

// Provider interface for LLM integrations
type Provider interface {
    // StreamCompletion streams LLM response chunks
    // Provider handles API-specific streaming and returns simple chunks
    // Agent layer converts these to AgentEvents
    StreamCompletion(ctx context.Context, messages []*types.Message) (<-chan *StreamChunk, error)
    
    // Complete returns full LLM response (convenience wrapper)
    Complete(ctx context.Context, messages []*types.Message) (*types.Message, error)
    
    // GetModelInfo returns model metadata
    GetModelInfo() *types.ModelInfo
}
```

### Agent Implementation

```go
// DefaultAgent wraps provider and emits AgentEvents
func (a *DefaultAgent) processInput(input *types.Input) {
    // Agent emits thinking event
    a.channels.Event <- &types.AgentEvent{
        Type: types.EventTypeThinkingStart,
    }
    
    // Call provider
    stream, err := a.provider.StreamCompletion(ctx, a.history)
    
    // Agent emits message start
    a.channels.Event <- &types.AgentEvent{
        Type: types.EventTypeMessageStart,
    }
    
    // Convert provider chunks to agent events
    for chunk := range stream {
        if chunk.Error != nil {
            a.channels.Event <- &types.AgentEvent{
                Type: types.EventTypeError,
                Error: chunk.Error,
            }
            break
        }
        
        a.channels.Event <- &types.AgentEvent{
            Type: types.EventTypeMessageContent,
            Content: chunk.Content,
        }
        
        if chunk.Finished {
            break
        }
    }
    
    // Agent emits message end
    a.channels.Event <- &types.AgentEvent{
        Type: types.EventTypeMessageEnd,
    }
    
    // Agent emits thinking end
    a.channels.Event <- &types.AgentEvent{
        Type: types.EventTypeThinkingEnd,
    }
}
```

## Naming Analysis

### Current Names
- `GenerateStream` - Generic but unclear
- `Generate` - Too generic

### Better Names (Chat-focused)
- `StreamCompletion` / `Complete` ✅
- `ChatStream` / `Chat` ✅
- `StreamChatCompletion` / `ChatCompletion` ✅ (matches OpenAI)

### Recommendation: Match Industry Standard
Use `ChatCompletion` / `StreamChatCompletion` to match OpenAI/Anthropic naming conventions.

## Migration Path

1. Update `Provider` interface to return `StreamChunk` instead of `AgentEvent`
2. Move event emission logic from provider → agent
3. Keep interface names generic enough for non-chat LLMs
4. Provider stays in `pkg/llm`, agent stays in `pkg/agent`

## Conclusion

**Current design violates separation of concerns by coupling provider to agent events.**

**Recommended fix**: Provider returns simple LLM chunks, Agent wraps provider and emits AgentEvents.

This makes Provider:
- ✅ Reusable outside agent context
- ✅ Testable independently
- ✅ Aligned with industry patterns
- ✅ Simpler to implement
- ✅ Easier to add new providers