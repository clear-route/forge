# Agent Loop System Design

## Overview

This document outlines the design for evolving Forge from a single-turn agent to a full agent loop system with tool calling, dynamic prompt assembly, planning context, and loop control tools.

## Current Architecture

```
┌─────────────┐
│   Executor  │
│   (CLI)     │
└──────┬──────┘
       │ Input
       ▼
┌─────────────┐      ┌──────────────┐
│    Agent    │─────▶│   Provider   │
│  (Default)  │◀─────│   (OpenAI)   │
└──────┬──────┘      └──────────────┘
       │ Events
       ▼
┌─────────────┐
│   Executor  │
│  (Render)   │
└─────────────┘
```

**Current Flow:**
1. User sends message
2. Agent processes once through LLM
3. Agent emits response events
4. Executor renders response
5. Loop back to step 1

## Target Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                        Executor (CLI)                        │
│  ┌────────────────────────────────────────────────────────┐  │
│  │              Agent Loop Controller                     │  │
│  │  • Manages iteration loop                              │  │
│  │  • Detects loop-breaking tools                        │  │
│  │  • Handles tool execution                             │  │
│  └────────────────────────────────────────────────────────┘  │
└───────────┬──────────────────────────────────────┬───────────┘
            │                                      │
            ▼                                      ▼
    ┌───────────────┐                    ┌─────────────────┐
    │  Prompt       │                    │  Tool Registry  │
    │  Builder      │                    │  • Local Tools  │
    │               │                    │  • MCP Tools    │
    └───────┬───────┘                    └─────────┬───────┘
            │                                      │
            │ Dynamic System Prompt                │ Tool Schemas
            ▼                                      │
    ┌───────────────────────────────────────────┐  │
    │           DefaultAgent                     │  │
    │  ┌─────────────────────────────────────┐  │  │
    │  │  Message History                    │  │  │
    │  │  • System messages                  │  │  │
    │  │  │• User messages                   │  │  │
    │  │  • Assistant messages               │  │  │
    │  │  • Tool messages                    │  │  │
    │  └─────────────────────────────────────┘  │  │
    │                                            │  │
    │  ┌─────────────────────────────────────┐  │  │
    │  │  Planning Context                   │  │  │
    │  │  • Current plan                     │  │  │
    │  │  • Completed steps                  │  │  │
    │  │  • Next actions                     │  │  │
    │  └─────────────────────────────────────┘  │  │
    └────────────────┬──────────────────────────┘  │
                     │                              │
                     ▼                              │
            ┌────────────────┐                      │
            │   Provider     │                      │
            │   (OpenAI)     │                      │
            └────────┬───────┘                      │
                     │                              │
                     ▼                              │
            ┌────────────────┐                      │
            │  LLM Response  │                      │
            │  with Tool Call│──────────────────────┘
            └────────────────┘
```

## Key Components

### 1. Tool System (`pkg/agent/tools/`)

**Tool Interface:**
```go
type Tool interface {
    Name() string
    Description() string
    Schema() string  // JSON schema for parameters
    Execute(ctx context.Context, args map[string]interface{}) (interface{}, error)
}
```

**Tool Registry:**
```go
type Registry struct {
    localTools map[string]Tool
    mcpServers map[string]*MCPServer
}
```

**Built-in Loop Control Tools:**
- `task_completion` - Exits loop, presents final results
- `ask_question` - Exits loop, asks user for input
- `converse` - Exits loop, engages in conversation

**MCP Integration:**
```go
type MCPServer struct {
    name  string
    tools map[string]Tool
}
```

### 2. Prompt System (`pkg/agent/prompts/`)

**Static Prompts:**
```go
const (
    SystemCapabilitiesPrompt
    AgentLoopPrompt
    ChainOfThoughtPrompt
    ToolCallingPrompt
    McpToolCallingPrompt
    ToolUseRulesPrompt
)
```

**Dynamic Prompt Builder:**
```go
type PromptBuilder struct {
    instructions         string
    additionalCapabilities []string
    localTools          map[string]Tool
    mcpTools            map[string][]Tool
}

func (pb *PromptBuilder) BuildMessages(
    history []Message,
    planContext string,
    ragContext map[string]string,
    finalMessage Message,
) []Message
```

**Message Assembly Order:**
1. System message with:
   - User instructions (persona/task)
   - Agent loop instructions
   - System capabilities
   - Chain of thought guidance
2. Conversation history
3. Current turn message with:
   - Tool schemas
   - Tool calling instructions
   - Planning context (if any)
   - RAG context (if any)
   - User's current input

### 3. Agent Loop (`pkg/agent/`)

**Enhanced DefaultAgent:**
```go
type DefaultAgent struct {
    // Existing fields
    provider     llm.Provider
    channels     *types.AgentChannels
    systemPrompt string
    
    // New fields
    promptBuilder *prompts.PromptBuilder
    toolRegistry  *tools.Registry
    memory        *memory.ConversationMemory
    planContext   string
}
```

**Agent Loop Flow:**
```
1. Receive user input
2. Build dynamic system prompt
3. Assemble messages (history + tools + planning + input)
4. Call LLM
5. Parse response for tool calls
6. If tool call detected:
   a. Emit tool_call event
   b. Execute tool
   c. Emit tool_result event
   d. If loop-breaking tool → exit loop
   e. Else → add to history, goto step 2
7. If no tool call:
   a. Emit no_tool_call event
   b. Re-prompt with reminder
   c. Goto step 3
8. If conversation/task_completion → exit loop
```

### 4. Memory System (`pkg/agent/memory/`)

**Conversation Memory:**
```go
type ConversationMemory struct {
    messages    []types.Message
    maxMessages int
    maxTokens   int
}

func (m *ConversationMemory) Add(msg *types.Message)
func (m *ConversationMemory) GetHistory() []types.Message
func (m *ConversationMemory) Prune() // Remove old messages
```

### 5. Planning System (`pkg/agent/planning/`)

**Plan Context:**
```go
type PlanContext struct {
    goal          string
    steps         []PlanStep
    currentStep   int
}

type PlanStep struct {
    description string
    completed   bool
    result      string
}
```

## Event System Updates

**New Event Types:**
```go
const (
    // Existing events...
    EventTypeToolCall
    EventTypeToolResult  
    EventTypeToolResultError
    EventTypeNoToolCall
    
    // New events
    EventTypeLoopStart        // Agent loop started
    EventTypeLoopIteration    // New loop iteration
    EventTypeLoopEnd          // Agent loop ended
    EventTypePlanUpdate       // Planning context updated
)
```

## Tool Call Format

The agent will use XML-based tool calling format in responses:

```xml
<thinking>
I need to search for information about the weather.
</thinking>

<tool>
{
  "server_name": "local",
  "tool_name": "web_search",
  "arguments": {
    "query": "current weather in San Francisco"
  }
}
</tool>
```

## Prompt Templates

### System Capabilities
```
<system_capabilities>
- Analyze user messages and determine the best course of action
- Maintain conversational context and remember previous interactions
- Use task_completion tool to mark tasks as complete
- Use ask_question tool to gather clarifying information
- Utilize various tools to complete tasks step by step
- Provide clear explanations of reasoning process
</system_capabilities>
```

### Agent Loop
```
<agent_loop>
You operate in an iterative agent loop:
1. Analyze Events: Understand user needs and current state
2. Select Tools: Choose next tool based on analysis
3. Iterate: One tool per iteration, repeat until task complete
4. Submit Results: Use task_completion to present final deliverables
5. Ask Questions: Use ask_question when clarification needed
</agent_loop>
```

### Chain of Thought
```
<chain_of_thought>
Before each action, outline your thought process in <thinking> tags:
- What information do you have
- What steps are needed
- What tool will help accomplish this
- What challenges might arise

This ensures systematic thinking and clear communication.
</chain_of_thought>
```

### Tool Calling
```
<tool_calling>
Format tool calls as XML with embedded JSON:
<tool>
{
  "server_name": "local",
  "tool_name": "tool_name_here",
  "arguments": { ... }
}
</tool>

Rules:
- One tool per turn
- Must include server_name field
- Never mention tool names to users
- Explain why you're using the tool before calling it
</tool_calling>
```

## Implementation Phases

### Phase 1: Tool Infrastructure
1. Create `pkg/agent/tools/` package
2. Implement Tool interface
3. Implement Registry
4. Create built-in loop control tools:
   - task_completion
   - ask_question
   - converse
5. Add tool call parsing to agent

### Phase 2: Prompt System
1. Create `pkg/agent/prompts/` package
2. Define static prompt constants
3. Implement PromptBuilder
4. Add prompt assembly logic
5. Integrate with DefaultAgent

### Phase 3: Memory System  
1. Create `pkg/agent/memory/` package
2. Implement ConversationMemory
3. Add history management to agent
4. Add context window pruning

### Phase 4: Agent Loop
1. Modify DefaultAgent to support iterations
2. Add tool call detection and parsing
3. Add tool execution flow
4. Implement loop-breaking logic
5. Update event emissions

### Phase 5: Planning (Optional Enhancement)
1. Create `pkg/agent/planning/` package
2. Implement PlanContext
3. Add planning prompt integration
4. Add plan tracking to agent loop

### Phase 6: MCP Integration
1. Define MCP server interface
2. Implement MCP tool wrapping
3. Add MCP-specific prompts
4. Integrate with tool registry

## API Changes

### Agent Construction
```go
// Before
ag := agent.NewDefaultAgent(provider,
    agent.WithSystemPrompt("You are helpful."),
)

// After
registry := tools.NewRegistry()
registry.RegisterLocal(tools.BuiltinTaskCompletion())
registry.RegisterLocal(tools.BuiltinAskQuestion())
registry.RegisterMCP("filesystem", fsServer)

mem := memory.NewConversationMemory(
    memory.WithMaxMessages(100),
)

ag := agent.NewDefaultAgent(provider,
    agent.WithSystemPrompt("You are helpful."),
    agent.WithToolRegistry(registry),
    agent.WithMemory(mem),
)
```

### Executor Usage
```go
// Executor stays the same - loop logic is internal to agent
executor := cli.NewExecutor(ag)
err := executor.Run(ctx)
```

## Benefits

1. **True Agentic Behavior**: Multi-turn reasoning and tool use
2. **Flexible Tool System**: Easy to add local or MCP tools
3. **Loop Control**: Agent can decide when to exit loop
4. **Planning Support**: Track complex multi-step tasks
5. **Memory Management**: Automatic context window handling
6. **Backward Compatible**: Existing simple usage still works

## Migration Path

Existing code continues to work:
```go
// This still works - simple single-turn mode
ag := agent.NewDefaultAgent(provider,
    agent.WithSystemPrompt("You are helpful."),
)
```

New capabilities are opt-in:
```go
// Enable agent loop with tools
ag := agent.NewDefaultAgent(provider,
    agent.WithSystemPrompt("You are helpful."),
    agent.WithToolRegistry(registry), // Enables loop mode
)
```

## File Structure

```
pkg/agent/
├── agent.go              # Agent interface
├── default.go            # DefaultAgent implementation  
├── options.go            # Agent options
├── agent_test.go         # Tests
│
├── core/                 # Internal utilities
│   └── stream.go         # Stream processing
│
├── tools/                # Tool system
│   ├── tool.go           # Tool interface
│   ├── registry.go       # Tool registry
│   ├── parser.go         # Parse tool calls from LLM
│   ├── builtin/          # Built-in tools
│   │   ├── task_completion.go
│   │   ├── ask_question.go
│   │   └── converse.go
│   └── mcp/              # MCP integration
│       ├── server.go
│       └── client.go
│
├── prompts/              # Prompt system
│   ├── builder.go        # PromptBuilder
│   ├── static.go         # Static prompt constants
│   └── formatter.go      # Format tools/plans/etc
│
├── memory/               # Memory system
│   ├── memory.go         # Memory interface
│   ├── conversation.go   # ConversationMemory
│   └── window.go         # Context window mgmt
│
└── planning/             # Planning system (future)
    ├── plan.go
    └── context.go
```

## Next Steps

1. Review this design with stakeholders
2. Create detailed implementation tasks
3. Begin with Phase 1 (Tool Infrastructure)
4. Iterate through phases incrementally
5. Add comprehensive tests at each phase
6. Update documentation and examples