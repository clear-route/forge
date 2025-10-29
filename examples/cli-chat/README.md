# CLI Chat Example

This example demonstrates a complete turn-by-turn conversation with an AI agent using the Forge framework.

## Components

- **DefaultAgent**: Manages conversation state, processes user inputs, and coordinates with the LLM provider
- **OpenAI Provider**: Handles communication with OpenAI-compatible APIs
- **CLI Executor**: Provides a terminal-based interface for interacting with the agent

## Usage

### Standard OpenAI API

```bash
export OPENAI_API_KEY="your-api-key-here"
go run examples/cli-chat/main.go
```

### Custom OpenAI-Compatible Endpoint

```bash
export OPENAI_API_KEY="your-api-key-here"
export OPENAI_BASE_URL="https://your-endpoint.com/v1/ai"
go run examples/cli-chat/main.go
```

## Features

- ✅ **Turn-by-turn conversation**: Type messages and get streaming responses
- ✅ **Conversation history**: The agent remembers previous messages in the conversation
- ✅ **Streaming responses**: See the AI's response as it's generated
- ✅ **Graceful shutdown**: Type 'exit' or 'quit' to end the conversation
- ✅ **Error handling**: Clear error messages if something goes wrong

## Example Conversation

```
Forge CLI Agent
Type your message and press Enter. Type 'exit' or 'quit' to end the conversation.

You: hello
Assistant:
Hello! How can I assist you today?

You: what is 2+2?
Assistant:
2 + 2 equals 4.

You: exit
Shutting down...
Goodbye!
```

## Customization

You can customize the agent's behavior by modifying the configuration in `main.go`:

```go
config := types.NewAgentConfig().
    WithSystemPrompt("You are a helpful AI assistant. Be concise and friendly.").
    WithStreaming(true)
```

You can also customize the CLI executor:

```go
executor := cli.NewExecutor(ag,
    cli.WithShowThinking(true),  // Show the agent's thinking process
    cli.WithPrompt("Me: "),       // Custom prompt
)
```

## Architecture

The example demonstrates the key layers of the Forge framework:

1. **Provider Layer** (`pkg/llm`): Handles LLM API communication
   - Returns simple `StreamChunk` instances
   - Reusable across different contexts

2. **Agent Layer** (`pkg/agent`): Manages conversation and state
   - Converts `StreamChunk` to `AgentEvent`
   - Maintains conversation history
   - Handles cancellation and errors

3. **Executor Layer** (`pkg/executor`): Provides user interface
   - Renders events to the terminal
   - Manages input/output
   - Coordinates conversation flow

This separation of concerns makes the framework flexible and testable.