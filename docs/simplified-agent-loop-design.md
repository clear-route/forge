# Simplified Agent Loop Design (Revised)

## Core Principles

1. **Always Agent Loop** - Every agent runs in a loop by default
2. **Built-in Tools** - task_completion, ask_question, converse are always available
3. **Simple Tool Registration** - Users call `agent.RegisterTool()` to add custom tools
4. **No MCP** - Local tools only for now

## Simplified API

### Creating an Agent (Always has agent loop)

```go
// Create agent - loop is automatic
ag := agent.NewDefaultAgent(provider,
    agent.WithSystemPrompt("You are a helpful coding assistant."),
    agent.WithMaxIterations(10), // Optional, default 10
)

// Register custom tools
ag.RegisterTool(tools.NewWebSearchTool())
ag.RegisterTool(tools.NewFileReadTool())

// Use it (loop happens automatically)
executor := cli.NewExecutor(ag)
executor.Run(ctx)
```

### Built-in Tools (Always Available)

```go
// These are ALWAYS registered automatically:
// 1. task_completion - Exit loop and present results
// 2. ask_question - Exit loop and ask user for input
// 3. converse - Exit loop and send a message to user
```

### Custom Tool Example

```go
type WebSearchTool struct{}

func (t *WebSearchTool) Name() string {
    return "web_search"
}

func (t *WebSearchTool) Description() string {
    return "Search the web for information"
}

func (t *WebSearchTool) Schema() string {
    return `{
        "type": "object",
        "properties": {
            "query": {
                "type": "string",
                "description": "The search query"
            }
        },
        "required": ["query"]
    }`
}

func (t *WebSearchTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
    query := args["query"].(string)
    // Perform search...
    return results, nil
}
```

## DefaultAgent Structure

```go
type DefaultAgent struct {
    // Existing fields
    provider     llm.Provider
    channels     *types.AgentChannels
    systemPrompt string
    maxTurns     int
    bufferSize   int
    
    // New fields for agent loop
    tools         map[string]Tool  // Built-in + custom tools
    memory        *memory.ConversationMemory
    promptBuilder *prompts.PromptBuilder
    maxIterations int
    
    // Message history
    historyMu sync.RWMutex
    history   []*types.Message
}
```

## Agent Loop Flow (Simplified)

```
User Input
    ↓
Initialize Built-in Tools (automatic)
    ↓
┌─────────────── Agent Loop ──────────────────┐
│                                             │
│  1. Build System Prompt                     │
│     - User instructions                     │
│     - Agent loop guidance                   │
│     - Chain of thought                      │
│     - Tool schemas (built-in + custom)      │
│                                             │
│  2. Assemble Messages                       │
│     - System message                        │
│     - Conversation history                  │
│     - Current user input                    │
│                                             │
│  3. Call LLM                                │
│                                             │
│  4. Parse Response                          │
│     ├─ Has <tool> tag?                      │
│     │                                        │
│     ├─ Yes: Parse tool call                 │
│     │   ├─ Execute tool                     │
│     │   ├─ Is loop-breaking?                │
│     │   │   ├─ Yes: EXIT LOOP →             │
│     │   │   └─ No: Add result to history    │
│     │   │         Continue loop ↑           │
│     │                                        │
│     └─ No: Emit NoToolCall event            │
│         Add reminder prompt                 │
│         Continue loop ↑                     │
│                                             │
│  5. Check Max Iterations                    │
│     └─ Reached? EXIT with error             │
│                                             │
└─────────────────────────────────────────────┘
```

## Built-in Tools Design

### 1. task_completion

```go
// Exits loop and presents final results
{
  "server_name": "local",
  "tool_name": "task_completion",
  "arguments": {
    "result": "I've completed the analysis. The results show...",
    "artifacts": ["file1.txt", "report.pdf"]  // optional
  }
}
```

### 2. ask_question

```go
// Exits loop and asks user a question
{
  "server_name": "local",
  "tool_name": "ask_question",
  "arguments": {
    "question": "Which approach would you prefer: A or B?",
    "suggestions": ["Approach A", "Approach B", "Show me both"]  // optional
  }
}
```

### 3. converse

```go
// Exits loop and sends a conversational message
{
  "server_name": "local",
  "tool_name": "converse",
  "arguments": {
    "message": "I understand. Let me explain how this works..."
  }
}
```

## Package Structure (Simplified)

```
pkg/agent/
├── agent.go              # Agent interface
├── default.go            # DefaultAgent with agent loop
├── options.go            # Agent options
├── agent_test.go         # Tests
│
├── core/                 # Internal utilities
│   └── stream.go
│
├── tools/                # Tool system
│   ├── tool.go           # Tool interface
│   ├── parser.go         # Parse tool calls
│   ├── task_completion.go    # Built-in
│   ├── ask_question.go       # Built-in
│   └── converse.go           # Built-in
│
├── prompts/              # Prompt system
│   ├── builder.go        # Build dynamic prompts
│   ├── static.go         # Static templates
│   └── formatter.go      # Format tool schemas
│
└── memory/               # Memory system
    ├── memory.go         # Memory interface
    └── conversation.go   # ConversationMemory
```

## Event Types (Agent Loop)

```go
const (
    // Existing events...
    EventTypeThinkingStart
    EventTypeThinkingContent
    EventTypeThinkingEnd
    EventTypeMessageStart
    EventTypeMessageContent
    EventTypeMessageEnd
    
    // Tool events
    EventTypeToolCall         // Tool is being called
    EventTypeToolResult       // Tool succeeded
    EventTypeToolResultError  // Tool failed
    EventTypeNoToolCall       // LLM didn't call a tool
    
    // Loop events
    EventTypeLoopStart        // Loop started
    EventTypeLoopIteration    // New iteration
    EventTypeLoopEnd          // Loop completed normally
    EventTypeLoopMaxReached   // Hit max iterations
    
    // Exit events (from built-in tools)
    EventTypeTaskComplete     // task_completion called
    EventTypeQuestionAsked    // ask_question called
    EventTypeConverse         // converse called
)
```

## Prompts (Core Templates)

### System Capabilities
```
You operate in an agent loop to accomplish tasks step by step.
Available capabilities:
- Think through problems systematically
- Use tools to gather information and take actions
- Ask questions when you need clarification
- Complete tasks and present results
```

### Agent Loop Instructions
```
<agent_loop>
You iterate through these steps:
1. Analyze: Understand the current state
2. Plan: Decide what tool to use next
3. Execute: Call ONE tool per iteration
4. Review: Check if task is complete
5. Repeat or Exit: Continue loop or use a completion tool

Exit the loop by calling:
- task_completion: When task is done
- ask_question: When you need user input
- converse: When you want to discuss with user
</agent_loop>
```

### Chain of Thought
```
<chain_of_thought>
Always show your reasoning in <thinking> tags before taking action:
- What do I know?
- What do I need to find out?
- Which tool will help?
- What might go wrong?
</chain_of_thought>
```

### Tool Calling Format
```
<tool_calling>
Call tools using this XML format:

<thinking>Your reasoning here</thinking>

<tool>
{
  "server_name": "local",
  "tool_name": "tool_name_here",
  "arguments": {
    "param": "value"
  }
}
</tool>

Rules:
- ONE tool per response
- server_name must be "local"
- Never mention tool names to users
- Explain why before calling
</tool_calling>
```

## Example Usage

### Simple Agent (Just Built-ins)

```go
package main

import (
    "context"
    "github.com/entrhq/forge/pkg/agent"
    "github.com/entrhq/forge/pkg/executor/cli"
    "github.com/entrhq/forge/pkg/llm/openai"
)

func main() {
    provider, _ := openai.NewProvider(
        os.Getenv("OPENAI_API_KEY"),
        openai.WithModel("gpt-4o"),
    )
    
    // Agent loop with built-in tools only
    ag := agent.NewDefaultAgent(provider,
        agent.WithSystemPrompt("You are a helpful assistant."),
    )
    
    executor := cli.NewExecutor(ag)
    executor.Run(context.Background())
}
```

### Agent with Custom Tools

```go
func main() {
    provider, _ := openai.NewProvider(
        os.Getenv("OPENAI_API_KEY"),
        openai.WithModel("gpt-4o"),
    )
    
    ag := agent.NewDefaultAgent(provider,
        agent.WithSystemPrompt("You are a web research assistant."),
        agent.WithMaxIterations(15),
    )
    
    // Register custom tools
    ag.RegisterTool(tools.NewWebSearchTool())
    ag.RegisterTool(tools.NewFileReadTool())
    ag.RegisterTool(tools.NewCalculatorTool())
    
    executor := cli.NewExecutor(ag)
    executor.Run(context.Background())
}
```

## Implementation Order

### Phase 1: Tool Infrastructure
1. Create `pkg/agent/tools/` package
2. Define `Tool` interface
3. Implement `task_completion` built-in
4. Implement `ask_question` built-in
5. Implement `converse` built-in
6. Implement tool call parser

### Phase 2: Prompt System
1. Create `pkg/agent/prompts/` package
2. Define static prompt templates
3. Implement `PromptBuilder`
4. Implement tool schema formatter

### Phase 3: Memory System
1. Create `pkg/agent/memory/` package
2. Implement `ConversationMemory`
3. Add history management

### Phase 4: Agent Loop
1. Add fields to `DefaultAgent`
2. Add `RegisterTool()` method
3. Auto-initialize built-in tools
4. Implement loop iteration logic
5. Integrate tool execution
6. Add loop control (max iterations, break conditions)

### Phase 5: Integration
1. Update examples
2. Add tests
3. Update documentation

## Key Differences from Original Design

| Original | Simplified |
|----------|-----------|
| Opt-in agent loop | Always agent loop |
| Separate Registry class | Tools map in DefaultAgent |
| `WithToolRegistry()` option | `RegisterTool()` method |
| MCP support | Local tools only |
| Complex registry API | Simple registration |

## Benefits

1. **Simpler API** - Just `RegisterTool()`, no registry management
2. **Always Agentic** - No confusion about modes
3. **Built-in Safety** - Loop-breaking tools always available
4. **Easy Extension** - Add custom tools with simple interface
5. **Clear Flow** - One way to use the framework