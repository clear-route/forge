# Forge

[![CI](https://github.com/entrhq/forge/workflows/CI/badge.svg)](https://github.com/entrhq/forge/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/entrhq/forge)](https://goreportcard.com/report/github.com/entrhq/forge)
[![GoDoc](https://pkg.go.dev/badge/github.com/entrhq/forge)](https://pkg.go.dev/github.com/entrhq/forge)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

**Forge** is an open-source, lightweight agent framework for building AI agents with pluggable components. It provides a clean, modular architecture that makes it easy to create agents with different LLM providers and execution environments.

## Features

- 🔌 **Pluggable Architecture**: Interface-based design for maximum flexibility
- 🤖 **LLM Provider Abstraction**: Support for OpenAI-compatible APIs with extensibility for custom providers
- 🚀 **Execution Plane Abstraction**: Run agents in different environments (CLI, API, custom)
- 📦 **Library-First Design**: Import as a Go module in your own applications
- 🧪 **Well-Tested**: Comprehensive test coverage with continuous integration
- 📖 **Well-Documented**: Clear documentation and examples

## Installation

```bash
go get github.com/entrhq/forge
```

## Quick Start

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
    
    // 2. Create an agent with the provider
    ag := agent.NewDefaultAgent(provider,
        agent.WithSystemPrompt("You are a helpful AI assistant."),
    )
    
    // 3. Create an executor
    executor := cli.NewExecutor(ag,
        cli.WithPrompt("You: "),
    )
    
    // 4. Run the agent
    if err := executor.Run(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```

## Architecture

Forge is built with a clean, modular architecture:

- **Agent Core** ([`pkg/agent`](pkg/agent)): Core agent interface and functionality
- **LLM Providers** ([`pkg/llm`](pkg/llm)): Pluggable LLM provider implementations
- **Executors** ([`pkg/executor`](pkg/executor)): Different execution environments for agents
- **Types** ([`pkg/types`](pkg/types)): Shared types and interfaces

See [`docs/architecture.md`](docs/architecture.md) for detailed architecture documentation.

## Project Structure

```
forge/
├── pkg/              # Public, importable packages
│   ├── agent/        # Agent core
│   ├── llm/          # LLM provider abstractions
│   ├── executor/     # Execution plane abstractions
│   └── types/        # Shared types
├── internal/         # Private implementation
├── examples/         # Example applications
│   ├── cli-chat/     # Basic CLI chat example
│   ├── cli-chat-thinking/ # CLI chat with thinking mode
│   └── simple-agent/ # Core types demonstration
├── docs/            # Documentation
└── .github/         # CI/CD workflows
```

## Examples

Check out the [`examples/`](examples/) directory for working examples:

- **CLI Chat** ([`examples/cli-chat`](examples/cli-chat)): Basic conversational agent
- **CLI Chat with Thinking** ([`examples/cli-chat-thinking`](examples/cli-chat-thinking)): Agent that shows reasoning process
- **Simple Agent** ([`examples/simple-agent`](examples/simple-agent)): Core types demonstration

### Running Examples

```bash
# Run the CLI chat example
cd examples/cli-chat
export OPENAI_API_KEY="your-api-key"
go run main.go

# Run with thinking mode
cd examples/cli-chat-thinking
export OPENAI_API_KEY="your-api-key"
go run main.go

# Explore core types
go run examples/simple-agent/main.go
```

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

We welcome contributions! Please see [`CONTRIBUTING.md`](CONTRIBUTING.md) for guidelines.

### Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment.

## Roadmap

- [x] Streaming response support
- [x] Basic CLI executor
- [x] OpenAI provider implementation
- [ ] Tool/function calling system
- [ ] State persistence and memory management
- [ ] Multi-agent coordination
- [ ] Additional LLM provider implementations (Anthropic, Google, etc.)
- [ ] Advanced executor implementations (HTTP API server, Slack bot, etc.)
- [ ] Prompt template system
- [ ] Agent collaboration and handoffs

## License

This project is licensed under the Apache License 2.0 - see the [`LICENSE`](LICENSE) file for details.

## Acknowledgments

Built as part of the Entr Agent Platform.

---

**Status**: 🚧 Under Active Development

This framework is currently in early development. APIs may change as we iterate on the design.
