# Your First Agent - Detailed Tutorial

This tutorial walks you through building a complete Forge agent step-by-step, explaining each concept along the way.

## What You'll Build

A functional chat agent that:
- Understands and responds to natural language
- Can use tools to perform tasks
- Shows its thinking process
- Maintains conversation history
- Handles errors gracefully

## Prerequisites

- Go 1.21+ installed
- OpenAI API key configured
- Forge installed (`go get github.com/entrhq/forge`)
- Basic Go programming knowledge

## Project Setup

Create a new directory for your project:

```bash
mkdir my-agent
cd my-agent
go mod init my-agent
```

Install Forge:

```bash
go get github.com/entrhq/forge
```

## Step 1: Understanding the Architecture

Before writing code, let's understand the key components:

```
User Input → Executor → Agent → LLM Provider → LLM (GPT-4)
                ↑         ↓
                └─ Tools ─┘
```

- **Executor**: Handles input/output (we'll use CLI)
- **Agent**: Orchestrates the conversation and tool usage
- **LLM Provider**: Connects to the LLM service (i.e. OpenAI)
- **Tools**: Functions the agent can call

## Step 2: Create the LLM Provider

The provider connects to your LLM service. Create `main.go`:

```go
package main

import (
    "log"
    "os"
    
    "github.com/entrhq/forge/pkg/llm/openai"
)

func main() {
    // Get API key from environment
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        log.Fatal("OPENAI_API_KEY environment variable not set")
    }
    
    // Create OpenAI provider
    provider, err := openai.NewProvider(
        apiKey,
        openai.WithModel("gpt-4o"),           // Use GPT-4 Optimized
        openai.WithTemperature(0.7),          // Balanced creativity
        openai.WithMaxTokens(2000),           // Max response length
    )
    if err != nil {
        log.Fatalf("Failed to create provider: %v", err)
    }
    
    log.Println("LLM Provider created successfully")
}
```

### Understanding Provider Options

- **Model**: Which LLM to use
  - `gpt-4o` - Latest, fastest GPT-4
  - `gpt-4o-mini` - Cheaper, good for most tasks
  - `gpt-4-turbo` - Previous generation
  
- **Temperature**: Controls randomness (0.0 - 2.0)
  - `0.0` - Deterministic, consistent
  - `0.7` - Balanced (recommended)
  - `1.5+` - Very creative

- **MaxTokens**: Maximum response length
  - Affects cost and response completeness
  - 2000 is a good default

## Step 3: Create the Agent

The agent is the "brain" that coordinates everything:

```go
package main

import (
    "log"
    "os"
    
    "github.com/entrhq/forge/pkg/agent"
    "github.com/entrhq/forge/pkg/llm/openai"
)

func main() {
    // ... provider code from Step 2 ...
    
    // Create agent with custom instructions
    ag := agent.NewDefaultAgent(
        provider,
        agent.WithCustomInstructions(
            `You are a helpful AI assistant.
            Be concise but thorough in your responses.
            Always explain your reasoning when using tools.
        `),
        agent.WithMaxIterations(10),  // Max tool calls per turn
    )
    
    log.Println("Agent created successfully")
}
```

### Understanding Agent Options

- **CustomInstructions**: System prompt that guides behavior
  - Be specific about personality and style
  - Include guidelines for tool usage
  - Keep it focused and clear

- **MaxIterations**: Prevents infinite loops
  - Default is 10
  - Increase for complex multi-step tasks
  - Decrease to save costs

## Step 4: Add the Executor

The executor handles how the agent interacts with users:

```go
package main

import (
    "context"
    "log"
    "os"
    
    "github.com/entrhq/forge/pkg/agent"
    "github.com/entrhq/forge/pkg/executor/cli"
    "github.com/entrhq/forge/pkg/llm/openai"
)

func main() {
    // 1. Create provider
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        log.Fatal("OPENAI_API_KEY environment variable not set")
    }
    
    provider, err := openai.NewProvider(
        apiKey,
        openai.WithModel("gpt-4o"),
        openai.WithTemperature(0.7),
    )
    if err != nil {
        log.Fatalf("Failed to create provider: %v", err)
    }
    
    // 2. Create agent
    ag := agent.NewDefaultAgent(
        provider,
        agent.WithCustomInstructions("You are a helpful AI assistant."),
    )
    
    // 3. Create CLI executor
    executor := cli.NewExecutor(ag)
    
    // 4. Run the agent
    log.Println("Starting agent... (Press Ctrl+C to exit)")
    if err := executor.Run(context.Background()); err != nil {
        log.Fatalf("Executor error: %v", err)
    }
}
```

## Step 5: Run Your Agent

Set your API key and run:

```bash
export OPENAI_API_KEY="your-api-key-here"
go run main.go
```

You'll see:

```
Starting agent... (Press Ctrl+C to exit)

You: 
```

Try it out:

```
You: Hello! What can you help me with?

[The agent is greeting the user and explaining its capabilities...]
Assistant: Hello! I'm an AI assistant powered by Forge. I can help you with:
- Answering questions
- Explaining concepts  
- Solving problems
- Having conversations

What would you like to know?
```

## Understanding the Output

The agent shows its "thinking" in brackets:

```
[The agent is analyzing the user's question...]
```

This is **chain-of-thought reasoning** - you can see how the agent is processing your request.

## Step 6: Adding a Custom Tool

Let's add a calculator tool. Create `calculator.go`:

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
)

// CalculatorTool performs basic arithmetic
type CalculatorTool struct{}

func (t *CalculatorTool) Name() string {
    return "calculator"
}

func (t *CalculatorTool) Description() string {
    return "Performs basic arithmetic operations (add, subtract, multiply, divide)"
}

func (t *CalculatorTool) Schema() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "operation": map[string]interface{}{
                "type":        "string",
                "description": "The operation to perform",
                "enum":        []string{"add", "subtract", "multiply", "divide"},
            },
            "a": map[string]interface{}{
                "type":        "number",
                "description": "First number",
            },
            "b": map[string]interface{}{
                "type":        "number",
                "description": "Second number",
            },
        },
        "required": []string{"operation", "a", "b"},
    }
}

func (t *CalculatorTool) Execute(ctx context.Context, arguments json.RawMessage) (string, error) {
    // Parse arguments
    var args struct {
        Operation string  `json:"operation"`
        A         float64 `json:"a"`
        B         float64 `json:"b"`
    }
    
    if err := json.Unmarshal(arguments, &args); err != nil {
        return "", fmt.Errorf("invalid arguments: %w", err)
    }
    
    // Perform calculation
    var result float64
    switch args.Operation {
    case "add":
        result = args.A + args.B
    case "subtract":
        result = args.A - args.B
    case "multiply":
        result = args.A * args.B
    case "divide":
        if args.B == 0 {
            return "", fmt.Errorf("cannot divide by zero")
        }
        result = args.A / args.B
    default:
        return "", fmt.Errorf("unknown operation: %s", args.Operation)
    }
    
    return fmt.Sprintf("%.2f", result), nil
}

func (t *CalculatorTool) IsLoopBreaking() bool {
    return false  // Agent can continue after using this tool
}
```

### Understanding Tools

Every tool must implement 5 methods:

1. **Name()** - Unique identifier
2. **Description()** - What it does (LLM sees this)
3. **Schema()** - JSON schema for parameters
4. **Execute()** - The actual logic
5. **IsLoopBreaking()** - Should agent stop after using this?

## Step 7: Register the Tool

Update `main.go` to register the calculator:

```go
func main() {
    // ... provider and agent creation ...
    
    // Register custom tool
    calculator := &CalculatorTool{}
    if err := ag.RegisterTool(calculator); err != nil {
        log.Fatalf("Failed to register tool: %v", err)
    }
    
    log.Println("Calculator tool registered")
    
    // ... executor and run ...
}
```

## Step 8: Test the Calculator

Run the agent and try:

```
You: What is 15 * 23?

[I need to calculate 15 multiplied by 23...]
<Executing: calculator(operation=multiply, a=15, b=23)>
Tool result: 345.00

[The calculation is complete...]
Assistant: 15 × 23 = 345
```

The agent:
1. Recognized it needed to calculate
2. Called the calculator tool
3. Got the result (345.00)
4. Responded to you

## Complete Code

Here's the full `main.go` with everything:

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "os"
    
    "github.com/entrhq/forge/pkg/agent"
    "github.com/entrhq/forge/pkg/executor/cli"
    "github.com/entrhq/forge/pkg/llm/openai"
)

// CalculatorTool implementation (from Step 6)
type CalculatorTool struct{}

func (t *CalculatorTool) Name() string {
    return "calculator"
}

func (t *CalculatorTool) Description() string {
    return "Performs basic arithmetic operations (add, subtract, multiply, divide)"
}

func (t *CalculatorTool) Schema() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "operation": map[string]interface{}{
                "type":        "string",
                "description": "The operation to perform",
                "enum":        []string{"add", "subtract", "multiply", "divide"},
            },
            "a": map[string]interface{}{
                "type":        "number",
                "description": "First number",
            },
            "b": map[string]interface{}{
                "type":        "number",
                "description": "Second number",
            },
        },
        "required": []string{"operation", "a", "b"},
    }
}

func (t *CalculatorTool) Execute(ctx context.Context, arguments json.RawMessage) (string, error) {
    var args struct {
        Operation string  `json:"operation"`
        A         float64 `json:"a"`
        B         float64 `json:"b"`
    }
    
    if err := json.Unmarshal(arguments, &args); err != nil {
        return "", fmt.Errorf("invalid arguments: %w", err)
    }
    
    var result float64
    switch args.Operation {
    case "add":
        result = args.A + args.B
    case "subtract":
        result = args.A - args.B
    case "multiply":
        result = args.A * args.B
    case "divide":
        if args.B == 0 {
            return "", fmt.Errorf("cannot divide by zero")
        }
        result = args.A / args.B
    default:
        return "", fmt.Errorf("unknown operation: %s", args.Operation)
    }
    
    return fmt.Sprintf("%.2f", result), nil
}

func (t *CalculatorTool) IsLoopBreaking() bool {
    return false
}

func main() {
    // 1. Get API key
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        log.Fatal("OPENAI_API_KEY environment variable not set")
    }
    
    // 2. Create LLM provider
    provider, err := openai.NewProvider(
        apiKey,
        openai.WithModel("gpt-4o"),
        openai.WithTemperature(0.7),
    )
    if err != nil {
        log.Fatalf("Failed to create provider: %v", err)
    }
    
    // 3. Create agent
    ag := agent.NewDefaultAgent(
        provider,
        agent.WithCustomInstructions("You are a helpful AI assistant."),
        agent.WithMaxIterations(10),
    )
    
    // 4. Register custom tools
    if err := ag.RegisterTool(&CalculatorTool{}); err != nil {
        log.Fatalf("Failed to register calculator: %v", err)
    }
    
    // 5. Create executor
    executor := cli.NewExecutor(ag)
    
    // 6. Run
    log.Println("Agent ready! Type your messages below.")
    if err := executor.Run(context.Background()); err != nil {
        log.Fatalf("Executor error: %v", err)
    }
}
```

## What You've Learned

1. **Provider** - Connects to LLM services
2. **Agent** - Orchestrates conversation and tools
3. **Executor** - Handles input/output
4. **Tools** - Extend agent capabilities
5. **Chain-of-thought** - Visible reasoning process

## Next Steps

- [Understanding the Agent Loop](understanding-agent-loop.md) - Deep dive into how agents work
- [Building Custom Tools](../guides/building-custom-tools.md) - More advanced tool patterns
- [Memory Management](../guides/memory-management.md) - Control conversation history
- [Error Handling](../guides/error-handling.md) - Robust error recovery

## Troubleshooting

### Agent doesn't use my tool

Make sure:
- Tool description is clear
- Schema matches your Execute() parameters
- Tool is registered before running executor

### "Too many iterations"

The agent hit MaxIterations. Either:
- Increase limit: `agent.WithMaxIterations(20)`
- Simplify the task
- Check if tool is returning useful results

### Tool execution fails

Add error handling in Execute():
```go
func (t *MyTool) Execute(ctx context.Context, args json.RawMessage) (string, error) {
    // Validate inputs
    // Handle edge cases
    // Return descriptive errors
}
```

## Full Example

See [`examples/agent-chat/`](../../examples/agent-chat/) for a complete working example.