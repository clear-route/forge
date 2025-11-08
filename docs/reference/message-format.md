# Message Format Reference

Complete reference for message formatting in Forge.

## Table of Contents

- [Message Structure](#message-structure)
- [Message Roles](#message-roles)
- [Tool Call Format](#tool-call-format)
- [Thinking Format](#thinking-format)
- [Token Estimation](#token-estimation)
- [Message Examples](#message-examples)

---

## Message Structure

### Core Message Type

```go
type Message struct {
    Role    string  // Message role (user/assistant/system/tool)
    Content string  // Message content
}
```

### Creating Messages

```go
import "github.com/yourusername/forge/pkg/message"

// User message
msg := message.User("Hello, how are you?")

// Assistant message
msg := message.Assistant("I'm doing well, thank you!")

// System message
msg := message.System("You are a helpful assistant")

// Tool result message
msg := message.Tool("calculator", "Result: 42")
```

---

## Message Roles

### System

Defines agent behavior and constraints.

```go
msg := message.System(`You are a Python expert assistant.
- Write clean, idiomatic Python code
- Include type hints
- Follow PEP 8 style guide
- Explain your reasoning`)
```

**Characteristics:**
- **Position:** First message in conversation
- **Count:** Usually one (can be multiple)
- **Preserved:** Never pruned from memory
- **Purpose:** Set agent personality and rules

**Example Usage:**

```go
memory.Add(message.System("You are a helpful coding assistant"))
memory.Add(message.User("Write a hello world function"))
```

---

### User

Messages from the human user.

```go
msg := message.User("What is the capital of France?")
```

**Characteristics:**
- **Source:** Human input
- **Format:** Plain text question or instruction
- **Frequency:** Every turn (user speaks)

**Multi-line:**

```go
msg := message.User(`Please help me with:
1. Understanding the code
2. Fixing the bug
3. Writing tests`)
```

---

### Assistant

Messages from the AI agent.

```go
msg := message.Assistant("The capital of France is Paris.")
```

**Characteristics:**
- **Source:** LLM response
- **Format:** May include thinking + tool calls + text
- **Frequency:** Every turn (agent responds)

**With Thinking:**

```go
msg := message.Assistant(`[Let me search for this information...]

Based on my search, the capital of France is Paris.`)
```

**With Tool Call:**

```go
msg := message.Assistant(`[I'll use the calculator for this]

<tool>
{
  "server_name": "local",
  "tool_name": "calculator",
  "arguments": {"a": 5, "b": 3, "operation": "add"}
}
</tool>`)
```

---

### Tool

Results from tool execution.

```go
msg := message.Tool("calculator", "Result: 8")
```

**Characteristics:**
- **Source:** Tool execution
- **Format:** `Tool: {name}\n{result}`
- **Purpose:** Provide tool output to agent

**Full Format:**

```
Tool: calculator
Result: 8
```

**Error Format:**

```
Tool: calculator
Error: division by zero
```

---

## Tool Call Format

### XML Structure

Tools are called using XML tags with JSON payload:

```xml
<tool>
{
  "server_name": "local",
  "tool_name": "tool_name",
  "arguments": {
    "param1": "value1",
    "param2": "value2"
  }
}
</tool>
```

### Fields

- **server_name:** Always `"local"` (for local tools)
- **tool_name:** Exact tool name from `Tool.Name()`
- **arguments:** JSON object matching tool schema

### Examples

**Simple Tool:**

```xml
<tool>
{
  "server_name": "local",
  "tool_name": "task_completion",
  "arguments": {
    "result": "The calculation result is 42"
  }
}
</tool>
```

**Complex Tool:**

```xml
<tool>
{
  "server_name": "local",
  "tool_name": "api_request",
  "arguments": {
    "method": "POST",
    "url": "https://api.example.com/data",
    "headers": {
      "Content-Type": "application/json"
    },
    "body": "{\"key\": \"value\"}"
  }
}
</tool>
```

**Multiple Tools:**

Agent can call multiple tools in sequence:

```xml
[I'll need to use multiple tools for this]

<tool>
{
  "server_name": "local",
  "tool_name": "search",
  "arguments": {"query": "Python tutorials"}
}
</tool>

[Now I'll summarize the results]

<tool>
{
  "server_name": "local",
  "tool_name": "task_completion",
  "arguments": {"result": "Here are the top Python tutorials..."}
}
</tool>
```

---

## Thinking Format

### Bracket Notation

Agent thinking is enclosed in `[brackets]`:

```
[This is my internal reasoning]
```

### Purpose

1. **Transparency:** Users see agent's thought process
2. **Debugging:** Understand why agent made decisions
3. **Better Results:** Chain-of-thought improves accuracy
4. **Trust:** Users understand the reasoning

### Examples

**Simple Thinking:**

```
[Let me calculate this step by step]
```

**Multi-step Reasoning:**

```
[First, I need to understand what the user is asking]
[They want the sum of 5 and 3]
[I'll use the calculator tool for this]
```

**Decision Making:**

```
[The user asked about Python, so I should use Python-specific tools]
[I have access to a Python code analyzer]
[Let me use that instead of the general code analyzer]
```

### Disabling Thinking

```go
agent, err := core.NewAgent(
    provider,
    memory,
    tools,
    core.WithThinkingEnabled(false),
)
```

When disabled, agent outputs only final responses (no brackets).

---

## Token Estimation

### Approximate Formula

Rough estimation (OpenAI models):

```
tokens ≈ len(text) / 4
```

More accurate (including special tokens):

```go
func EstimateTokens(msg message.Message) int {
    // Base estimate
    tokens := len(msg.Content) / 4
    
    // Role overhead
    tokens += 4
    
    // Special formatting
    if strings.Contains(msg.Content, "<tool>") {
        tokens += 10 // XML overhead
    }
    
    return tokens
}
```

### Token Counts by Role

Approximate overhead per message:

| Role | Overhead | Reason |
|------|----------|--------|
| system | ~4 tokens | Role marker |
| user | ~4 tokens | Role marker |
| assistant | ~4 tokens | Role marker |
| tool | ~10 tokens | Role + tool formatting |

### Example Calculations

```go
// "Hello" (user)
// Content: 1 token
// Role: 4 tokens
// Total: ~5 tokens

msg := message.User("Hello")
// Estimated: 5 tokens

// Long message
msg := message.User("Please explain the difference between lists and tuples in Python, including their performance characteristics and common use cases.")
// Content: ~25 tokens
// Role: 4 tokens
// Total: ~29 tokens

// Tool call
msg := message.Assistant(`<tool>
{
  "server_name": "local",
  "tool_name": "calculator",
  "arguments": {"a": 5, "b": 3}
}
</tool>`)
// Content: ~20 tokens
// XML: 10 tokens
// Role: 4 tokens
// Total: ~34 tokens
```

### Memory Pruning

ConversationMemory automatically prunes when token limit exceeded:

```go
memory := memory.NewConversationMemory(8000) // 8000 token limit

// Add messages until limit
memory.Add(message.User("Message 1"))
memory.Add(message.Assistant("Response 1"))
// ... many more messages ...

// When limit exceeded, oldest messages removed (except system)
```

---

## Message Examples

### Conversation Flow

```go
// 1. System message
memory.Add(message.System("You are a helpful math tutor"))

// 2. User question
memory.Add(message.User("What is 15 + 27?"))

// 3. Assistant with thinking
memory.Add(message.Assistant(`[I'll calculate this using the calculator tool]

<tool>
{
  "server_name": "local",
  "tool_name": "calculator",
  "arguments": {"operation": "add", "a": 15, "b": 27}
}
</tool>`))

// 4. Tool result
memory.Add(message.Tool("calculator", "Result: 42"))

// 5. Assistant final response
memory.Add(message.Assistant(`[Great, the calculator confirmed the result]

15 + 27 = 42`))

// 6. User follow-up
memory.Add(message.User("Thanks! Now what's 42 divided by 6?"))
```

### Error Handling

```go
// User asks question
memory.Add(message.User("What is 10 divided by 0?"))

// Agent tries calculator
memory.Add(message.Assistant(`[I'll use the calculator]

<tool>
{
  "server_name": "local",
  "tool_name": "calculator",
  "arguments": {"operation": "divide", "a": 10, "b": 0}
}
</tool>`))

// Tool returns error
memory.Add(message.Tool("calculator", "Error: division by zero"))

// Agent adapts
memory.Add(message.Assistant(`[I see the error - division by zero is undefined]

I cannot calculate 10 divided by 0 because division by zero is undefined in mathematics. This operation would result in infinity, which cannot be represented as a normal number.`))
```

### Multi-Tool Workflow

```go
// User request
memory.Add(message.User("Find Python tutorials and create a summary"))

// Agent: First tool
memory.Add(message.Assistant(`[I'll search for Python tutorials first]

<tool>
{
  "server_name": "local",
  "tool_name": "web_search",
  "arguments": {"query": "Python tutorials for beginners"}
}
</tool>`))

// Tool 1 result
memory.Add(message.Tool("web_search", "Found 10 tutorials..."))

// Agent: Second tool
memory.Add(message.Assistant(`[Now I'll create a summary]

<tool>
{
  "server_name": "local",
  "tool_name": "task_completion",
  "arguments": {
    "result": "Here are the top Python tutorials:\n1. Tutorial A...\n2. Tutorial B..."
  }
}
</tool>`))
```

### Complex Reasoning

```go
memory.Add(message.User("Analyze this code and suggest improvements"))

memory.Add(message.Assistant(`[Let me break this down step by step]
[1. First, I'll read the code to understand it]
[2. Then identify potential issues]
[3. Finally suggest specific improvements]

<tool>
{
  "server_name": "local",
  "tool_name": "code_analyzer",
  "arguments": {"code": "def func():\n  x=1\n  return x"}
}
</tool>`))

memory.Add(message.Tool("code_analyzer", "Issues: No docstring, inconsistent spacing"))

memory.Add(message.Assistant(`[Based on the analysis, I can provide specific suggestions]

Here are the improvements:
1. Add a docstring
2. Fix spacing around the equals sign (PEP 8)
3. Consider a more descriptive function name`))
```

---

## Message Validation

### Validate Message Structure

```go
func ValidateMessage(msg message.Message) error {
    // Check role
    validRoles := []string{"system", "user", "assistant", "tool"}
    valid := false
    for _, role := range validRoles {
        if msg.Role == role {
            valid = true
            break
        }
    }
    if !valid {
        return fmt.Errorf("invalid role: %s", msg.Role)
    }
    
    // Check content
    if msg.Content == "" {
        return fmt.Errorf("content cannot be empty")
    }
    
    return nil
}
```

### Validate Tool Call Format

```go
func ValidateToolCall(content string) error {
    // Check for tool tags
    if !strings.Contains(content, "<tool>") {
        return fmt.Errorf("missing <tool> tag")
    }
    if !strings.Contains(content, "</tool>") {
        return fmt.Errorf("missing </tool> tag")
    }
    
    // Extract JSON
    start := strings.Index(content, "<tool>") + 6
    end := strings.Index(content, "</tool>")
    jsonStr := content[start:end]
    
    // Validate JSON
    var toolCall map[string]interface{}
    if err := json.Unmarshal([]byte(jsonStr), &toolCall); err != nil {
        return fmt.Errorf("invalid JSON: %w", err)
    }
    
    // Check required fields
    if _, ok := toolCall["tool_name"]; !ok {
        return fmt.Errorf("missing tool_name")
    }
    if _, ok := toolCall["arguments"]; !ok {
        return fmt.Errorf("missing arguments")
    }
    
    return nil
}
```

---

## Best Practices

### 1. Always Include System Message

```go
// ✅ Good: Clear system message
memory.Add(message.System("You are a helpful coding assistant"))

// ❌ Bad: No system message
memory.Add(message.User("Help me code"))
```

### 2. Preserve Message Order

```go
// ✅ Good: Proper order
memory.Add(message.User("Question"))
memory.Add(message.Assistant("Response with tool call"))
memory.Add(message.Tool("tool_name", "Result"))
memory.Add(message.Assistant("Final answer"))

// ❌ Bad: Out of order
memory.Add(message.Tool("tool_name", "Result"))
memory.Add(message.Assistant("Response"))
memory.Add(message.User("Question"))
```

### 3. Keep Messages Focused

```go
// ✅ Good: Single clear message
memory.Add(message.User("What is 2+2?"))

// ❌ Bad: Multiple unrelated questions
memory.Add(message.User("What is 2+2? Also tell me about Paris and write a poem"))
```

### 4. Include Context in Tool Results

```go
// ✅ Good: Clear result
message.Tool("calculator", "Result: 42 (15 + 27)")

// ❌ Bad: Ambiguous
message.Tool("calculator", "42")
```

---

## See Also

- [API Reference](api-reference.md) - Message API details
- [Architecture: Agent Loop](../architecture/agent-loop.md) - How messages flow
- [Architecture: Memory System](../architecture/memory-system.md) - Message storage
- [Tool Schema Reference](tool-schema.md) - Tool call format details