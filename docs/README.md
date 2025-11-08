# Forge Documentation

Welcome to the Forge documentation! This guide will help you find what you need.

## 📚 Documentation Structure

### 🚀 [Getting Started](getting-started/)
New to Forge? Start here!

- [Installation](getting-started/installation.md) - Set up Forge in your project
- [Quick Start](getting-started/quick-start.md) - Build your first agent in 5 minutes
- [Your First Agent](getting-started/your-first-agent.md) - Detailed tutorial
- [Understanding the Agent Loop](getting-started/understanding-agent-loop.md) - Core concepts

### 📖 [How-To Guides](guides/)
Task-oriented guides for common scenarios:

- [Building Custom Tools](guides/building-custom-tools.md) - Create your own tools
- [Implementing LLM Providers](guides/implementing-llm-providers.md) - Add new LLM support
- [Creating Executors](guides/creating-executors.md) - Custom execution environments
- [Memory Management](guides/memory-management.md) - Handle conversation history
- [Error Handling](guides/error-handling.md) - Robust error recovery
- [Streaming Responses](guides/streaming-responses.md) - Real-time response streaming

### 📚 [API Reference](reference/)
Complete technical reference:

- [Agent API](reference/api/agent.md) - Agent interface and DefaultAgent
- [Tools API](reference/api/tools.md) - Tool system reference
- [LLM Providers API](reference/api/llm-providers.md) - Provider interface
- [Executors API](reference/api/executors.md) - Executor interface
- [Types](reference/api/types.md) - Core types and events
- [Built-in Tools](reference/built-in-tools.md) - task_completion, ask_question, converse
- [Configuration](reference/configuration.md) - All configuration options

### 🏗️ [Architecture](architecture/)
Understand how Forge works:

- [Architecture Overview](architecture/overview.md) - System design and components
- [Agent Loop](architecture/agent-loop.md) - How the agent loop works
- [Tool System](architecture/tool-system.md) - Tool architecture
- [Memory System](architecture/memory-system.md) - Memory management design
- [Design Decisions](architecture/design-decisions.md) - Key architectural choices

### 💡 [Examples](examples/)
Practical code examples:

- [Basic Chat Agent](examples/basic-chat-agent.md) - Simple agent walkthrough
- [Custom Calculator Tool](examples/custom-calculator-tool.md) - Tool creation example
- [Multi-Tool Workflow](examples/multi-tool-workflow.md) - Complex agent workflows
- [Custom LLM Provider](examples/custom-llm-provider.md) - Provider implementation

### 🤝 [Community](community/)
Resources for the Forge community:

- [FAQ](community/faq.md) - Frequently asked questions
- [Troubleshooting](community/troubleshooting.md) - Common issues and solutions
- [Best Practices](community/best-practices.md) - Tips and recommendations
- [Roadmap](community/roadmap.md) - Future plans and features

## 🔍 Quick Links

- **[Contributing](../CONTRIBUTING.md)** - How to contribute to Forge
- **[Code of Conduct](../CODE_OF_CONDUCT.md)** - Community guidelines
- **[Security Policy](../SECURITY.md)** - Report security issues
- **[Changelog](../CHANGELOG.md)** - Version history

## 💬 Need Help?

- Check the [FAQ](community/faq.md) for common questions
- Review [Troubleshooting](community/troubleshooting.md) for solutions
- Open an issue on [GitHub](https://github.com/entrhq/forge/issues)
- See working examples in [`/examples`](../examples/)

## 📦 Quick Start

```bash
# Install
go get github.com/entrhq/forge

# Create a simple agent
# See getting-started/quick-start.md for full example
```

---

**Choose your path:**
- 🆕 New user? Start with [Getting Started](getting-started/)
- 🔧 Building something? Check [How-To Guides](guides/)
- 📖 Need details? See [API Reference](reference/)
- 🧠 Want to understand? Read [Architecture](architecture/)
- 👥 Contributing? Visit [Community](community/)