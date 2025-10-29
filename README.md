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
    "github.com/entrhq/forge/pkg/agent"
    "github.com/entrhq/forge/pkg/llm"
    "github.com/entrhq/forge/pkg/executor"
)

func main() {
    // 1. Create an LLM provider
    provider := llm.NewOpenAIProvider(config)
    
    // 2. Create an agent with the provider
    agent := agent.New(provider, options)
    
    // 3. Create an executor
    executor := executor.NewCLI()
    
    // 4. Run the agent
    executor.Run(agent)
}
```

## Architecture

Forge is built with a clean, modular architecture:

- **Agent Core** ([`pkg/agent`](pkg/agent)): Core agent interface and functionality
- **LLM Providers** ([`pkg/llm`](pkg/llm)): Pluggable LLM provider implementations
- **Executors** ([`pkg/executor`](pkg/executor)): Different execution environments for agents
- **Configuration** ([`pkg/config`](pkg/config)): Centralized configuration management
- **Types** ([`pkg/types`](pkg/types)): Shared types and interfaces

See [`docs/architecture.md`](docs/architecture.md) for detailed architecture documentation.

## Project Structure

```
forge/
├── pkg/              # Public, importable packages
│   ├── agent/        # Agent core
│   ├── llm/          # LLM provider abstractions
│   ├── executor/     # Execution plane abstractions
│   ├── config/       # Configuration
│   └── types/        # Shared types
├── internal/         # Private implementation
├── examples/         # Example applications
│   └── simple-agent/ # Basic usage example
├── docs/            # Documentation
└── .github/         # CI/CD workflows
```

## Examples

Check out the [`examples/`](examples/) directory for working examples:

- **Simple Agent** ([`examples/simple-agent`](examples/simple-agent)): Basic agent implementation

### Running Examples

```bash
# Run the simple agent example
make run-example

# Or directly with go
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

- [ ] Tool/function calling system
- [ ] Streaming response support
- [ ] State persistence and memory management
- [ ] Multi-agent coordination
- [ ] Additional LLM provider implementations
- [ ] Advanced executor implementations (HTTP API server)
- [ ] Prompt template system

## License

This project is licensed under the Apache License 2.0 - see the [`LICENSE`](LICENSE) file for details.

## Acknowledgments

Built as part of the Entr Agent Platform.

---

**Status**: 🚧 Under Active Development

This framework is currently in early development. APIs may change as we iterate on the design.
