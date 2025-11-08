# Installation

This guide covers installing Forge and setting up your development environment.

## Prerequisites

- **Go 1.21 or higher** - [Download Go](https://golang.org/dl/)
- **Git** - For cloning the repository (optional)
- **OpenAI API Key** - Or compatible LLM provider API key

## Install Forge

### Using `go get` (Recommended)

Add Forge to your Go project:

```bash
go get github.com/entrhq/forge
```

This will download Forge and all its dependencies to your project.

### Verify Installation

Create a simple file to verify the installation:

```go
// verify.go
package main

import (
    "fmt"
    "github.com/entrhq/forge/pkg/agent"
)

func main() {
    fmt.Println("Forge installed successfully!")
    fmt.Printf("Agent package available: %T\n", agent.Agent(nil))
}
```

Run it:

```bash
go run verify.go
```

You should see output confirming Forge is installed.

## Set Up API Keys

Forge requires an LLM provider API key. Currently, OpenAI-compatible providers are supported.

### OpenAI

1. Get your API key from [OpenAI Platform](https://platform.openai.com/api-keys)
2. Set it as an environment variable:

```bash
# Linux/macOS
export OPENAI_API_KEY="your-api-key-here"

# Windows (PowerShell)
$env:OPENAI_API_KEY="your-api-key-here"

# Windows (Command Prompt)
set OPENAI_API_KEY=your-api-key-here
```

### Alternative: .env File

For development, you can use a `.env` file (don't commit this!):

```bash
# .env
OPENAI_API_KEY=your-api-key-here
```

Then load it in your application using a package like [`godotenv`](https://github.com/joho/godotenv):

```go
import "github.com/joho/godotenv"

func init() {
    godotenv.Load()
}
```

## Development Tools (Optional)

For contributing to Forge or running examples:

```bash
# Clone the repository
git clone https://github.com/entrhq/forge.git
cd forge

# Install development tools
make install-tools

# Run tests
make test

# Build examples
make examples
```

## IDE Setup

### VS Code

Recommended extensions:
- **Go** by Go Team at Google
- **Go Test Explorer** for running tests

### GoLand/IntelliJ IDEA

Go support is built-in. Just open the project directory.

## Verify Complete Setup

Create a minimal agent to verify everything works:

```go
package main

import (
    "context"
    "log"
    "os"

    "github.com/entrhq/forge/pkg/agent"
    "github.com/entrhq/forge/pkg/llm/openai"
)

func main() {
    // Check API key
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        log.Fatal("OPENAI_API_KEY environment variable not set")
    }

    // Create provider
    provider, err := openai.NewProvider(apiKey)
    if err != nil {
        log.Fatal(err)
    }

    // Create agent
    ag := agent.NewDefaultAgent(provider)
    
    log.Println("Setup complete! Agent created successfully.")
}
```

## Troubleshooting

### "Package not found"

Make sure you've run `go get` and your `go.mod` includes Forge:

```bash
go mod tidy
```

### API Key Issues

Verify your API key is set:

```bash
echo $OPENAI_API_KEY  # Linux/macOS
echo %OPENAI_API_KEY%  # Windows
```

### Module Errors

Ensure you're in a Go module:

```bash
go mod init your-project-name
```

## Next Steps

- [Quick Start](quick-start.md) - Build your first agent in 5 minutes
- [Your First Agent](your-first-agent.md) - Detailed tutorial
- [Understanding the Agent Loop](understanding-agent-loop.md) - Learn core concepts

## Additional Resources

- [OpenAI API Documentation](https://platform.openai.com/docs)
- [Go Modules Documentation](https://go.dev/ref/mod)
- [Forge Examples](../../examples/)