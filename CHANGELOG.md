# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

#### Documentation
- **TUI Interface Guide** - Comprehensive guide for using the Terminal User Interface
  - Complete coverage of all keyboard shortcuts and navigation
  - Detailed slash command reference (/help, /stop, /commit, /pr, /settings, /context, /bash)
  - Overlay system documentation (help, settings, context, approval, diff viewer, command execution, result list)
  - Tool approval workflow explanation with security best practices
  - Settings configuration guide with all tabs (General, LLM, Auto-Approval, Display)
  - Tips & troubleshooting for common issues
- **Built-in Tools Reference** - Complete documentation for all loop-breaking and operational tools
  - Detailed parameter schemas with examples
  - Use cases and best practices
  - XML encoding guidelines for tool calls
- Comprehensive documentation restructure with user-centric organization
- Getting Started guides for new users
- Complete API reference documentation
- How-to guides for common tasks
- Community resources (FAQ, troubleshooting, best practices)
- GitHub issue templates and PR template
- Code of Conduct
- Security policy

#### Features Completed
- ✅ **TUI Executor** - Full-featured terminal interface with Bubble Tea
  - Gemini-inspired chat interface with viewport scrolling
  - Command palette with slash command system
  - Multiple overlay types (help, settings, context, approval, diff viewer)
  - Syntax-highlighted diff viewer with side-by-side and unified modes
  - Real-time command execution with streaming output
  - Toast notifications for status updates
  - Token usage tracking and display
- ✅ **Tool Approval System** - Security-first approval workflow
  - Interactive approval dialogs with detailed tool information
  - Auto-approval rules with path/command patterns
  - Configurable approval requirements per tool
  - Settings-based approval configuration
- ✅ **Slash Commands** - Rich command system for TUI
  - `/help` - Show help overlay with tips and shortcuts
  - `/stop` - Cancel current agent operation
  - `/commit` - Create git commits from session changes
  - `/pr` - Create pull requests (requires approval)
  - `/settings` - Interactive settings configuration
  - `/context` - Display detailed context information
  - `/bash` - Enter bash mode for shell commands
- ✅ **Settings System** - Comprehensive configuration management
  - Multi-tab settings interface (General, LLM, Auto-Approval, Display)
  - Per-provider API key management
  - Auto-approval rule configuration
  - Display preferences (theme, syntax highlighting, diff style)
  - Settings persistence to `~/.config/forge/settings.json`
- ✅ **Context Overlay** - Detailed session information display
  - Workspace statistics and path information
  - Conversation history metrics
  - Token usage breakdown (system, user, assistant, tools)
  - Cumulative session tracking
  - Memory state visualization
- ✅ **Intelligent Result Display** - Smart tool result rendering
  - Result caching system (configurable size, default 20)
  - Collapsible tool result messages
  - Preview with configurable line count
  - Result list overlay for browsing all cached results
  - Automatic truncation with expand/collapse controls

#### Architecture & Infrastructure
- 25 Architecture Decision Records documenting design decisions
- Event-driven architecture with streaming support
- Self-healing error recovery with circuit breaker pattern
- Workspace security with path validation and sandboxing
- Memory management with conversation history and token-based pruning
- Comprehensive test coverage (196+ tests across all packages)

### Changed
- Reorganized documentation structure for better discoverability
- Archived historical design documents
- Enhanced CHANGELOG with detailed feature tracking and proper categorization

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