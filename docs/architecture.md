# Forge Architecture

## Overview

Forge is a lightweight, modular agent framework designed for building AI agents with pluggable components.

## Core Components

### Agent Core (`pkg/agent`)
The heart of the framework, defining the agent interface and core functionality.

### LLM Abstraction (`pkg/llm`)
Provides a provider-agnostic interface for integrating with various LLM services.
- OpenAI-compatible providers
- Extensible for custom providers

### Executor Interface (`pkg/executor`)
Abstracts the execution environment, allowing agents to run in different contexts:
- CLI applications
- API servers
- Custom implementations

### Configuration (`pkg/config`)
Centralized configuration management for agent setup and runtime options.

### Types (`pkg/types`)
Shared types and interfaces used across the framework.

## Internal Packages

### Core (`internal/core`)
Private business logic implementation details.

### Utils (`internal/utils`)
Internal utility functions not exposed to external consumers.

## Design Principles

1. **Interface-Driven**: All major components are defined as interfaces for maximum flexibility
2. **Pluggable Architecture**: Easy to swap implementations without changing application code
3. **Library-First**: Designed as an importable library, not a standalone application
4. **Clean Boundaries**: Clear separation between public (`pkg/`) and private (`internal/`) APIs

## Usage Pattern

```go
import (
    "github.com/clear-route/forge/pkg/agent"
    "github.com/clear-route/forge/pkg/llm"
    "github.com/clear-route/forge/pkg/executor"
)

// 1. Create an LLM provider
provider := llm.NewOpenAIProvider(config)

// 2. Create an agent
agent := agent.New(provider, options)

// 3. Create an executor
executor := executor.NewCLI()

// 4. Run the agent
executor.Run(agent)
```

## Extension Points

- **Custom LLM Providers**: Implement the `llm.Provider` interface
- **Custom Executors**: Implement the `executor.Executor` interface
- **Custom Agents**: Implement the `agent.Agent` interface

## Future Enhancements

- Tool/function calling system
- State persistence
- Multi-agent coordination
- Streaming support
- Advanced memory management