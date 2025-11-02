# Quick Start

Build your first Forge agent in 5 minutes!

## What You'll Build

A simple chat agent that:
- Uses OpenAI's GPT model
- Maintains conversation history
- Runs in your terminal

## Prerequisites

- [Go 1.21+ installed](installation.md#prerequisites)
- [OpenAI API key set up](installation.md#set-up-api-keys)
- [Forge installed](installation.md#install-forge)

## Step 1: Create Your Project

```bash
mkdir my-first-agent
cd my-first-agent
go mod init my-first-agent
go get github.com/entrhq/forge
```

## Step 2: Write the Code

Create `main.go`:

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
    // 1. Create an LLM provider
    provider, err := openai.NewProvider(
        os.Getenv("OPENAI_API_KEY"),
        openai.WithModel("gpt-4o"),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // 2. Create an agent with custom instructions
    ag := agent.NewDefaultAgent(provider,
        agent.WithCustomInstructions("You are a helpful AI assistant."),
    )
    
    // 3. Create a CLI executor
    executor := cli.NewExecutor(ag)
    
    // 4. Run the agent
    if err := executor.Run(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```

## Step 3: Run Your Agent

```bash
export OPENAI_API_KEY="your-api-key"
go run main.go
```

## Step 4: Interact

You'll see a prompt where you can chat with your agent:

```
You: Hello! What can you help me with?

[The agent is greeting the user and offering assistance...]
Assistant: Hello! I'm here to help you with a variety of tasks...

You: What's 25 * 4?

[The agent needs to calculate this multiplication...]
Assistant: 25 × 4 = 100
```

## What Just Happened?

Let's break down the code:

### 1. LLM Provider

```go
provider, err := openai.NewProvider(
    os.Getenv("OPENAI_API_KEY"),
    openai.WithModel("gpt-4o"),
)
```

Creates a connection to OpenAI's API. You can customize:
- Model (gpt-4o, gpt-4o-mini, etc.)
- Temperature, max tokens, and other parameters

### 2. Agent

```go
ag := agent.NewDefaultAgent(provider,
    agent.WithCustomInstructions("You are a helpful AI assistant."),
)
```

Creates an agent that:
- Uses your LLM provider
- Follows custom instructions
- Has built-in tools (task_completion, ask_question, converse)
- Manages conversation memory automatically

### 3. Executor

```go
executor := cli.NewExecutor(ag)
```

Creates a CLI executor that:
- Handles terminal input/output
- Shows thinking in `[brackets]`
- Manages the conversation loop

### 4. Run

```go
executor.Run(context.Background())
```

Starts the agent and keeps it running until you exit (Ctrl+C).

## Customization Options

### Change the Model

```go
provider, err := openai.NewProvider(
    os.Getenv("OPENAI_API_KEY"),
    openai.WithModel("gpt-4o-mini"),        // Faster, cheaper
    openai.WithTemperature(0.7),            // Creativity (0.0-2.0)
    openai.WithMaxTokens(1000),             // Response length limit
)
```

### Customize Agent Behavior

```go
ag := agent.NewDefaultAgent(provider,
    agent.WithCustomInstructions("You are a Python expert assistant."),
    agent.WithMaxIterations(15),           // Max tool calls per turn
)
```

### Add Custom Tools

```go
// Register a custom tool (covered in detail later)
err := ag.RegisterTool(myCustomTool)
```

## Next Steps

Now that you have a working agent, learn more:

- **[Your First Agent](your-first-agent.md)** - Detailed tutorial with explanations
- **[Understanding the Agent Loop](understanding-agent-loop.md)** - How agents think and act
- **[Building Custom Tools](../guides/building-custom-tools.md)** - Add custom capabilities

## Common Issues

### "OPENAI_API_KEY not set"

Set your API key:
```bash
export OPENAI_API_KEY="your-key-here"
```

### "Package not found"

Make sure you've installed Forge:
```bash
go get github.com/entrhq/forge
go mod tidy
```

### "Rate limit exceeded"

You're making too many API calls. Consider:
- Using gpt-4o-mini (cheaper and faster)
- Adding rate limiting
- Checking your OpenAI usage limits

## Complete Example

See the full working example at [`examples/agent-chat`](../../examples/agent-chat/).

---

**🎉 Congratulations!** You've built your first Forge agent. Ready to learn more? Continue to [Your First Agent](your-first-agent.md) for a deeper dive.