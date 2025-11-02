# Thinking Content Support

The Forge agent framework supports streaming thinking/reasoning content separately from message content. This is useful for models that expose their internal reasoning process, such as:

- **OpenAI o1/o1-mini**: Streams thinking tokens before the final answer
- **Anthropic Claude with extended_thinking**: Shows reasoning steps
- **Custom models**: Any model that separates reasoning from output

## How It Works

### Automatic Parsing

The framework automatically detects and parses `<thinking>` tags from LLM responses:

```
<thinking>
Step-by-step reasoning here...
</thinking>

Final answer here.
```

The [`ThinkingParser`](../pkg/llm/parser/thinking.go:1) handles this transparently at the provider level, so thinking content is automatically separated without any configuration needed.

### StreamChunk ContentType

The [`StreamChunk`](../pkg/llm/types.go:1) type includes a `Type` field that indicates whether content is thinking or message:

```go
type StreamChunk struct {
    Content  string
    Role     string
    Type     ContentType  // "thinking" or "message"
    Finished bool
    Error    error
}
```

### Content Types

- **`ContentTypeMessage`** (default): Regular message content - the final answer
- **`ContentTypeThinking`**: Thinking/reasoning content - intermediate thoughts

### Event Flow

When processing a stream with thinking content:

1. **First thinking chunk** → `ThinkingStart` event
2. **Thinking chunks** → `ThinkingContent` events
3. **Last thinking chunk** → `ThinkingEnd` event
4. **First message chunk** → `MessageStart` event
5. **Message chunks** → `MessageContent` events  
6. **Last message chunk** → `MessageEnd` event
7. **Stream complete** → `TurnEnd` event

## Provider Implementation

The OpenAI provider automatically parses `<thinking>` tags using [`ThinkingParser`](../pkg/llm/parser/thinking.go:1). For other providers, you can use the parser:

```go
import "github.com/entrhq/forge/pkg/llm/parser"

thinkingParser := parser.NewThinkingParser()

// In your streaming loop:
for contentChunk := range apiStream {
    thinkingChunk, messageChunk := thinkingParser.Parse(contentChunk)
    
    if thinkingChunk != nil {
        // Emit thinking chunk
        chunks <- thinkingChunk
    }
    
    if messageChunk != nil {
        // Emit message chunk
        chunks <- messageChunk
    }
}
```

## Agent Processing

The [`DefaultAgent.processStream()`](../pkg/agent/default.go:152) method automatically:

- Emits thinking events for thinking content
- Emits message events for message content
- Transitions from thinking to message seamlessly
- Only saves message content to conversation history (thinking is ephemeral)

## Executor Display

The CLI executor shows thinking content by default. You can disable it if needed:

```go
executor := cli.NewExecutor(agent) // Shows thinking by default

// Or explicitly disable:
executor := cli.NewExecutor(agent,
    cli.WithShowThinking(false),  // Hide thinking process
)
```

**Example with thinking (default):**
```
You: What is 2+2?

[Thinking...]
I need to add 2 and 2. This is a simple arithmetic operation.
2 + 2 = 4
[Done thinking]