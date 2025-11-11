# Forge - TUI Coding Agent

A powerful terminal-based coding assistant built on the Forge agent framework. Forge provides file operations, code editing, and command execution through an intuitive chat-first interface.

## Features

- 💬 **Chat-First Interface** - Natural conversation with your AI coding assistant
- 📁 **File Operations** - Read, write, list, and search files within your workspace
- ✏️ **Smart Editing** - Diff-based code changes without full file rewrites
- 🔒 **Workspace Security** - All operations restricted to your project directory
- 🚫 **Smart Ignoring** - Automatically excludes node_modules, .git, .env, and other clutter
- 🎨 **Syntax Highlighting** - Beautiful diff previews with syntax highlighting
- ⚡ **Command Execution** - Run terminal commands with approval workflow
- 🌳 **File Navigation** - Browse your project structure

## Installation

```bash
# Build from source
cd cmd/forge
go build -o forge

# Or install to GOPATH/bin
go install github.com/entrhq/forge/cmd/forge@latest
```

## Quick Start

### Using OpenAI

1. Set your OpenAI API key:
```bash
export OPENAI_API_KEY="your-api-key-here"
```

2. Run Forge in your project directory:
```bash
forge
```

3. Start chatting with your coding assistant!

### Using Alternative Providers

Forge supports any OpenAI-compatible API through the base URL configuration:

#### OpenRouter
```bash
export OPENAI_API_KEY="sk-or-v1-..."
export OPENAI_BASE_URL="https://openrouter.ai/api/v1"
forge -model anthropic/claude-3.5-sonnet
```

#### Local LLM (e.g., LM Studio, Ollama)
```bash
export OPENAI_API_KEY="dummy-key"
export OPENAI_BASE_URL="http://localhost:1234/v1"
forge -model local-model-name
```

#### Other OpenAI-Compatible APIs
```bash
export OPENAI_BASE_URL="https://your-api.com/v1"
forge
```

## Usage

```bash
# Start in current directory
forge

# Specify workspace directory
forge -workspace /path/to/project

# Use different model
forge -model gpt-4-turbo

# Use custom API endpoint
forge -base-url https://api.example.com/v1

# Show version
forge -version

# Show help
forge -help
```

### Examples

```bash
# OpenAI (default)
export OPENAI_API_KEY="sk-..."
forge

# OpenRouter with Claude
export OPENAI_API_KEY="sk-or-..."
forge -base-url https://openrouter.ai/api/v1 -model anthropic/claude-3.5-sonnet

# Local LLM
forge -base-url http://localhost:1234/v1 -model llama3

# Custom workspace and model
forge -workspace ~/projects/myapp -model gpt-4o
```

## Configuration

### Command Line Flags

- `-api-key` - OpenAI API key (or set `OPENAI_API_KEY` env var)
- `-base-url` - OpenAI API base URL (or set `OPENAI_BASE_URL` env var) - use for OpenAI-compatible APIs
- `-model` - LLM model to use (default: `gpt-4o`)
- `-workspace` - Workspace directory (default: current directory)
- `-prompt` - Custom system prompt for the agent
- `-version` - Show version and exit

### Environment Variables

- `OPENAI_API_KEY` - Your OpenAI API key (required)
- `OPENAI_BASE_URL` - Base URL for OpenAI-compatible APIs (optional, defaults to OpenAI)

### Supported Providers

Forge works with any OpenAI-compatible API:

- **OpenAI** (default) - Standard OpenAI API
- **OpenRouter** - Access to multiple model providers
- **Azure OpenAI** - Enterprise OpenAI deployment
- **LocalAI** - Run models locally
- **LM Studio** - Local model inference
- **Ollama** - Local model management
- **Any other OpenAI-compatible API**

Simply set the `OPENAI_BASE_URL` to your provider's endpoint.

## Keyboard Shortcuts

### Conversation Mode
- `Enter` - Send message
- `Ctrl+C` / `Esc` - Quit
- `Ctrl+T` - Toggle file tree (coming soon)
- `Ctrl+O` - Toggle command output (coming soon)

### Diff Viewer (coming soon)
- `j/k` or `↓/↑` - Navigate diff lines
- `Ctrl+A` - Accept changes
- `Ctrl+R` - Reject changes
- `Esc` - Cancel

## Example Workflows

### Read and Modify Code

```
You: Read the main.go file and add error handling to the run function

Forge: [Reads file, analyzes code, proposes diff]
- Shows side-by-side diff
- You approve with Ctrl+A
- Changes are applied
```

### Search Codebase

```
You: Find all TODO comments in the codebase

Forge: [Searches with regex, shows results with context]
```

### Execute Commands

```
You: Run the tests

Forge: [Proposes: go test ./...]
- Shows command for approval
- You approve with Ctrl+A
- Streams output in real-time
```

## Security

Forge operates with workspace-level security:

- All file operations are restricted to the specified workspace directory
- Path traversal attacks are prevented
- Symlinks pointing outside the workspace are blocked
- Commands require explicit user approval before execution

## File Ignoring

Forge automatically filters out common directories and files that clutter results:

### Default Ignore Patterns

The following are always ignored:
- `node_modules/`, `vendor/` - Dependency directories
- `.git/` - Version control metadata
- `.env`, `.env.*` - Environment files
- `*.log` - Log files
- `.DS_Store` - macOS system files
- `__pycache__/`, `*.pyc` - Python cache
- `.vscode/`, `.idea/` - IDE directories
- `dist/`, `build/`, `tmp/`, `temp/` - Build artifacts
- And more (see [ignore.go](../../pkg/security/workspace/ignore.go))

### Custom Patterns

You can customize ignore behavior with two optional files in your workspace root:

1. **`.gitignore`** - Respects your existing git ignore patterns
2. **`.forgeignore`** - Forge-specific patterns (highest priority)

Both files support standard gitignore syntax:
- `*.test` - Glob patterns
- `logs/` - Directory patterns
- `!important.log` - Negation patterns
- `# comment` - Comments

**Pattern Precedence** (highest to lowest):
1. `.forgeignore` patterns
2. `.gitignore` patterns
3. Default patterns

### Examples

Create a `.forgeignore` file to customize:

```
# Ignore all test files except integration tests
*.test
!*integration.test

# Ignore specific directories
tmp/
.cache/
```

For more details, see [ADR-0016: File Ignore System](../../docs/adr/0016-file-ignore-system.md).

## Development Status

Forge is currently in active development. Current status:

- ✅ Core framework and TUI executor
- ✅ Workspace security layer
- ✅ Basic application structure
- 🚧 Coding tools (in progress)
- 🚧 Diff viewer with approval flow (in progress)
- 📅 File tree navigation (planned)
- 📅 Command output display (planned)

## Troubleshooting

### "API key is required" error

Make sure you've set the `OPENAI_API_KEY` environment variable or passed it via the `-api-key` flag.

### "Workspace directory error"

Ensure the path you're providing exists and is a directory. The workspace is where Forge can read and write files.

### TUI not rendering correctly

Ensure your terminal supports ANSI colors and has a minimum size of 80x24.

## Architecture

Forge is built on a modular architecture:

- **Agent Core** - Handles conversation loop and tool orchestration
- **Coding Tools** - Reusable tools for file operations (read, write, search, diff, execute)
- **Security Layer** - Enforces workspace boundaries on all operations
- **TUI Executor** - Provides the terminal interface with overlays for diffs and navigation

See [Architecture Documentation](../../docs/architecture/overview.md) for details.

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](../../CONTRIBUTING.md) for guidelines.

## License

Apache 2.0 - See [LICENSE](../../LICENSE) for details.

## Links

- [Documentation](../../docs/)
- [Architecture Overview](../../docs/architecture/overview.md)
- [Coding Tools Architecture](../../docs/adr/0011-coding-tools-architecture.md)
- [Tool Approval Mechanism](../../docs/adr/0010-tool-approval-mechanism.md)
- [Enhanced TUI Design](../../docs/adr/0012-enhanced-tui-executor.md)

---

**Version:** 0.1.0  
**Status:** 🚧 Under Active Development