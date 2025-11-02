# Glossary

Complete glossary of terms used in Forge.

## A

### Agent
An autonomous AI system that can reason, use tools, and accomplish tasks through iterative loops. In Forge, agents are created via the [`Agent`](api-reference.md#agent) interface.

### Agent Loop
The iterative cycle where an agent thinks, acts (uses tools), and learns from results. See [Agent Loop Architecture](../architecture/agent-loop.md).

### API Key
Authentication credential for LLM provider services (e.g., OpenAI API key). Required for provider initialization.

### Arguments
Parameters passed to a tool during execution. Defined by tool's JSON Schema. See [Tool Schema Reference](tool-schema.md).

### Assistant Message
Message from the AI agent, may contain thinking, tool calls, or text responses. See [Message Format](message-format.md#assistant).

---

## C

### Chain-of-Thought
Reasoning process where agent shows its thinking in `[brackets]`. Improves accuracy and transparency. See [Thinking Format](message-format.md#thinking-format).

### Circuit Breaker
Error handling pattern that stops retries after N consecutive failures to prevent infinite loops. See [Circuit Breaker](error-handling.md#circuit-breaker).

### Complete
Non-streaming LLM API call that returns full response at once. See [`Provider.Complete()`](api-reference.md#provider).

### Context
The conversation history sent to the LLM. Includes system, user, assistant, and tool messages.

### Context Window
Maximum number of tokens an LLM can process in a single request (e.g., GPT-4 has 8K tokens).

### ConversationMemory
Default in-memory implementation of the Memory interface. Stores messages and handles automatic pruning. See [Memory System](../architecture/memory-system.md).

---

## E

### Executor
Interface defining how agents interact with their environment (CLI, web, etc.). See [Executor](api-reference.md#executor).

### Exponential Backoff
Retry strategy where wait time increases exponentially: 1s, 2s, 4s, 8s... See [Retry Strategies](error-handling.md#retry-strategies).

---

## I

### Interface
Go language contract defining required methods. Forge uses interfaces for Agent, Provider, Tool, Memory, and Executor.

### Iteration
Single cycle of the agent loop: think → act → observe. Controlled by `WithMaxIterations()`.

---

## J

### JSON Schema
Standard for describing JSON data structures. Used to define tool parameters. See [Tool Schema Reference](tool-schema.md).

---

## L

### LLM (Large Language Model)
AI model that generates text based on input (e.g., GPT-4, Claude). Accessed through providers.

### Loop-Breaking Tool
Tool that ends the agent loop when executed (e.g., `task_completion`). Returns `IsLoopBreaking() == true`.

---

## M

### MaxIterations
Configuration limit on agent loop iterations per turn. Prevents infinite loops. Default: 10.

### Memory
Component that stores conversation history. See [Memory Package](api-reference.md#memory-package).

### Message
Unit of conversation with role (system/user/assistant/tool) and content. See [Message Structure](message-format.md#message-structure).

### Mock
Test double that simulates real component behavior. See [Mocking](testing.md#mocking).

---

## P

### Provider
Interface for LLM services (OpenAI, Anthropic, etc.). Handles API communication. See [Provider Package](api-reference.md#provider-package).

### Pruning
Automatic removal of old messages to stay within token limits. See [Memory Management](../architecture/memory-system.md#pruning).

---

## R

### Response
LLM output containing generated text and metadata. Returned by `Provider.Complete()`.

### Role
Message classification: `system`, `user`, `assistant`, or `tool`. See [Message Roles](message-format.md#message-roles).

---

## S

### Schema
See JSON Schema.

### Server Name
Identifier for tool source. Currently always `"local"` for built-in tools.

### Stream
LLM API mode that returns response in chunks as generated. See [`Provider.Stream()`](api-reference.md#provider).

### StreamChunk
Individual piece of a streaming response. Contains partial content.

### System Message
First message defining agent behavior and constraints. Never pruned. See [System](message-format.md#system).

### System Prompt
Content of the system message. Set via `WithSystemPrompt()`.

---

## T

### Temperature
LLM parameter controlling randomness (0.0 = deterministic, 2.0 = very random). See [Configuration](configuration.md#withtemperature).

### Thinking
Agent's internal reasoning displayed in `[brackets]`. Can be disabled with `WithThinkingEnabled(false)`.

### Token
Basic unit of text for LLMs. Roughly ¼ of a word. Used for pricing and limits.

### Tool
Function that agents can execute. Implements the Tool interface. See [Tool System](../architecture/tool-system.md).

### Tool Call
When agent requests to execute a tool. Uses XML format with JSON payload. See [Tool Call Format](message-format.md#tool-call-format).

### Tool Name
Unique identifier for a tool. Returned by `Tool.Name()`.

### Tool Result
Output from tool execution. Stored as tool message in memory.

---

## U

### User Message
Message from human user. Represents user input or questions. See [User](message-format.md#user).

---

## X

### XML
Format used to wrap tool calls: `<tool>{JSON}</tool>`. See [Tool Call Format](message-format.md#tool-call-format).

---

## Common Patterns

### Agent Loop Pattern

```
User Input
  ↓
[Agent Thinks]
  ↓
Agent Uses Tool
  ↓
Tool Returns Result
  ↓
[Agent Thinks About Result]
  ↓
Agent Responds or Uses Another Tool
  ↓
Loop Continues Until Task Complete
```

### Message Flow Pattern

```
System Message (behavior)
  ↓
User Message (question)
  ↓
Assistant Message (thinking + tool call)
  ↓
Tool Message (result)
  ↓
Assistant Message (final response)
```

### Error Recovery Pattern

```
Tool Execution
  ↓
Error Occurs
  ↓
Error Added to Memory
  ↓
Agent Sees Error
  ↓
Agent Adapts Approach
  ↓
Tries Again or Alternative
```

---

## Acronyms

| Acronym | Full Term | Meaning |
|---------|-----------|---------|
| AI | Artificial Intelligence | Computer systems performing tasks requiring human intelligence |
| API | Application Programming Interface | Interface for software interaction |
| CLI | Command-Line Interface | Text-based user interface |
| JSON | JavaScript Object Notation | Data interchange format |
| LLM | Large Language Model | AI model trained on text data |
| OSS | Open Source Software | Software with public source code |
| RAM | Random Access Memory | Computer memory |
| SDK | Software Development Kit | Tools for software development |
| SSE | Server-Sent Events | Server push technology |
| UUID | Universally Unique Identifier | Unique identifier |
| XML | Extensible Markup Language | Markup language for documents |

---

## Code Examples

### Creating an Agent

```go
agent, err := core.NewAgent(provider, memory, tools)
```

### Tool Call

```xml
<tool>
{
  "server_name": "local",
  "tool_name": "calculator",
  "arguments": {"operation": "add", "a": 5, "b": 3}
}
</tool>
```

### Message Types

```go
message.System("You are helpful")     // System
message.User("Hello")                 // User
message.Assistant("Hi there!")        // Assistant
message.Tool("calculator", "Result: 8") // Tool
```

---

## Related Terms

### Agent vs Tool
- **Agent:** Autonomous system that uses tools
- **Tool:** Function that agent can execute

### Memory vs Context
- **Memory:** Storage of conversation history
- **Context:** Messages sent to LLM in a request

### Complete vs Stream
- **Complete:** Get full response at once
- **Stream:** Get response in chunks

### Role vs Type
- **Role:** Message classification (user/assistant/etc.)
- **Type:** Data type in schema (string/number/etc.)

### Pruning vs Clearing
- **Pruning:** Selective removal of old messages
- **Clearing:** Remove all messages

---

## Design Patterns

### Interface-Driven Design
All major components defined as interfaces for flexibility and testability.

### Circuit Breaker Pattern
Prevent infinite retry loops by stopping after N failures.

### Strategy Pattern
Different memory pruning strategies (oldest-first, importance-based).

### Observer Pattern
Executor observes agent actions and displays them.

### Factory Pattern
Create agents with functional options pattern.

---

## Configuration Terms

### Functional Options
Pattern for configuring components:

```go
core.WithMaxIterations(10)
core.WithSystemPrompt("You are helpful")
core.WithThinkingEnabled(true)
```

### Environment Variables
Configuration stored in system environment:

```bash
OPENAI_API_KEY="sk-..."
OPENAI_BASE_URL="https://api.openai.com/v1"
```

---

## Performance Terms

### Latency
Time between request and response. Affected by model choice and context size.

### Throughput
Number of requests handled per unit time.

### Token Budget
Maximum tokens allocated for context. Controls memory pruning.

### Token Estimation
Approximating token count: `tokens ≈ len(text) / 4`

---

## Testing Terms

### Mock
Fake implementation for testing. See [Mocking](testing.md#mocking).

### Table-Driven Test
Test pattern with multiple test cases in a table structure.

### Benchmark
Performance test measuring speed and resource usage.

### Race Detector
Tool for finding concurrent access bugs: `go test -race`

---

## See Also

- [API Reference](api-reference.md) - Complete API documentation
- [Architecture](../architecture/overview.md) - System design
- [Getting Started](../getting-started/quick-start.md) - Quick start guide
- [Examples](../examples/) - Code examples