# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive documentation restructure with user-centric organization
- Getting Started guides for new users
- Complete API reference documentation
- How-to guides for common tasks
- Community resources (FAQ, troubleshooting, best practices)
- GitHub issue templates and PR template
- Code of Conduct
- Security policy

### Changed
- Reorganized documentation structure for better discoverability
- Archived historical design documents

## [0.1.0]

### Added
- Core agent framework with pluggable architecture
- Agent loop system with automatic tool execution
- Tool system with built-in loop-breaking tools:
  - `task_completion` - Signal task completion
  - `ask_question` - Request clarification from user
  - `converse` - Engage in natural conversation
- Custom tool registration API
- OpenAI-compatible LLM provider implementation
- CLI executor for running agents in terminal
- Chain-of-thought reasoning with thinking blocks
- Memory management with conversation history
- Event-driven architecture with streaming support
- Self-healing error recovery with circuit breaker pattern
- Comprehensive test coverage (196+ tests)
- Example applications demonstrating framework usage

### Features
- Interface-based design for maximum flexibility
- Tool execution with JSON schema validation
- Dynamic prompt assembly with tool schemas
- Token-based memory pruning
- Thread-safe conversation history
- Streaming response support
- Custom instruction injection
- Iteration limit controls

### Documentation
- Architecture overview
- Agent loop implementation details
- Contributing guidelines
- Example code and tutorials
- Design decision documentation

[Unreleased]: https://github.com/entrhq/forge/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/entrhq/forge/releases/tag/v0.1.0