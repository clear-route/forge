# Forge

[![CI](https://github.com/entrhq/forge/workflows/CI/badge.svg)](https://github.com/entrhq/forge/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/entrhq/forge)](https://goreportcard.com/report/github.com/entrhq/forge)
[![GoDoc](https://pkg.go.dev/badge/github.com/entrhq/forge)](https://pkg.go.dev/github.com/entrhq/forge)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

**Forge** is an open-source, lightweight agent framework for building AI agents with pluggable components. It provides a clean, modular architecture that makes it easy to create agents with different LLM providers and execution environments.

## Features

- 🔌 **Pluggable Architecture**: Interface-based design for maximum flexibility
- 🤖 **LLM Provider Abstraction**: Support for OpenAI-compatible APIs with extensibility for custom providers
- 🛠️ **Tool System**: Agent loop with tool execution and custom tool registration
- 🧠 **Chain-of-Thought**: Built-in thinking/reasoning capabilities for transparent agent behavior
- 💾 **Memory Management**: Conversation history and context management
- 🔄 **Event-Driven**: Real-time streaming of thinking, tool calls, and messages
- 🔁 **Self-Healing Error Recovery**: Automatic error recovery with circuit breaker pattern
- 🚀 **Execution Plane Abstraction**: Run agents in different environments (CLI, API, custom)
- 📦 **Library-First Design**: Import as a Go module in your own applications
- 🧪 **Well-Tested**: Comprehensive test coverage (196+ tests passing)
- 📖 **Well-Documented**: Clear, comprehensive documentation

## Quick Start

```bash
go get github.com/entrhq/forge
```

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

## Documentation

### 📚 [Complete Documentation](docs/)

**Getting Started:**
- [Installation](docs/getting-started/installation.md) - Set up Forge in your project
- [Quick Start](docs/getting-started/quick-start.md) - Build your first agent in 5 minutes
- [Your First Agent](docs/getting-started/your-first-agent.md) - Detailed step-by-step tutorial
- [Understanding the Agent Loop](docs/getting-started/understanding-agent-loop.md) - Core concepts explained

**Guides:**
- [Building Custom Tools](docs/guides/building-custom-tools.md) - Extend agent capabilities
- [Implementing LLM Providers](docs/guides/implementing-llm-providers.md) - Add new LLM support
- [Creating Executors](docs/guides/creating-executors.md) - Custom execution environments
- [Memory Management](docs/guides/memory-management.md) - Control conversation history
- [Error Handling](docs/guides/error-handling.md) - Robust error recovery

**Reference:**
- [API Reference](docs/reference/) - Complete API documentation
- [Built-in Tools](docs/reference/built-in-tools.md) - task_completion, ask_question, converse
- [Configuration](docs/reference/configuration.md) - All configuration options

**Architecture:**
- [Architecture Overview](docs/architecture/overview.md) - System design and components
- [Agent Loop](docs/architecture/agent-loop.md) - How the agent loop works
- [Tool System](docs/architecture/tool-system.md) - Tool architecture details
- [Design Decisions](docs/architecture/design-decisions.md) - Key architectural choices

**Community:**
- [FAQ](docs/community/faq.md) - Frequently asked questions
- [Troubleshooting](docs/community/troubleshooting.md) - Common issues and solutions
- [Best Practices](docs/community/best-practices.md) - Tips and recommendations

## Architecture

Forge is built with a clean, modular architecture:

- **Agent Core** ([`pkg/agent`](pkg/agent)): Agent loop, tools, prompts, and memory
- **LLM Providers** ([`pkg/llm`](pkg/llm)): Pluggable LLM provider implementations
- **Executors** ([`pkg/executor`](pkg/executor)): Different execution environments (CLI, API, etc.)
- **Types** ([`pkg/types`](pkg/types)): Shared types, events, and interfaces

### Key Components

- **Tools** ([`pkg/agent/tools`](pkg/agent/tools)): Tool interface and built-in tools
- **Prompts** ([`pkg/agent/prompts`](pkg/agent/prompts)): Dynamic prompt assembly with tool schemas
- **Memory** ([`pkg/agent/memory`](pkg/agent/memory)): Conversation history management
- **Stream Processing** ([`pkg/agent/core`](pkg/agent/core)): Real-time parsing of thinking, tools, and messages

See [Architecture Overview](docs/architecture/overview.md) for detailed documentation.

## Examples

Check out the [`examples/`](examples/) directory for working examples:

- **Agent Chat** ([`examples/agent-chat`](examples/agent-chat)): Complete agent with tools, thinking, and custom tool registration

### Running the Example

```bash
cd examples/agent-chat
export OPENAI_API_KEY="your-api-key"
go run main.go
```

Try asking:
- "What is 15 * 23?"
- "Calculate (100 + 50) / 3"
- "What's 144 divided by 12, then add 5?"

## Development

### Prerequisites

- Go 1.21 or higher
- Make (optional, but recommended)

### Setup

```bash
# Clone the repository
git clone https://github.com/entrhq/forge.git
cd forge

# Install development tools
make install-tools

# Run tests
make test

# Run linter
make lint

# Format code
make fmt
```

### Available Make Targets

- `make test` - Run tests with coverage
- `make lint` - Run linters
- `make fmt` - Format code
- `make examples` - Build example applications
- `make run-example` - Run simple example
- `make clean` - Clean build artifacts
- `make all` - Run all checks and build examples

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

### Security

For security issues, please see our [Security Policy](SECURITY.md).

## Roadmap

- [x] Streaming response support
- [x] Basic CLI executor
- [x] OpenAI provider implementation
- [x] Tool/function calling system
- [x] Agent loop with infinite iterations
- [x] Chain-of-thought reasoning
- [x] Memory management and conversation history
- [x] Event-driven architecture
- [x] Custom tool registration
- [x] Self-healing error recovery with circuit breaker
- [ ] Auto-pruning for memory management
- [ ] Integration tests for full agent loop
- [ ] Multi-agent coordination
- [ ] Additional LLM provider implementations (Anthropic, Google, etc.)
- [ ] Advanced executor implementations (HTTP API server, Slack bot, etc.)
- [ ] Agent collaboration and handoffs

See [Roadmap](docs/community/roadmap.md) for detailed plans.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

Built as part of the Entr Agent Platform.

---

**Status**: 🚧 Under Active Development

This framework is currently in early development. APIs may change as we iterate on the design.

## Links

- [Documentation](docs/) - Complete documentation
- [Examples](examples/) - Working code examples
- [Contributing](CONTRIBUTING.md) - How to contribute
- [Changelog](CHANGELOG.md) - Version history
- [Issues](https://github.com/entrhq/forge/issues) - Report bugs or request features