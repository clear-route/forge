// Package main provides the Forge TUI coding agent application.
// This is a flagship coding assistant that runs entirely in the terminal,
// providing file operations, code editing, and command execution with an
// intuitive chat-first interface.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/entrhq/forge/pkg/agent"
	"github.com/entrhq/forge/pkg/agent/tools"
	"github.com/entrhq/forge/pkg/executor/tui"
	"github.com/entrhq/forge/pkg/llm/openai"
	"github.com/entrhq/forge/pkg/security/workspace"
	"github.com/entrhq/forge/pkg/tools/coding"
)

const (
	version      = "0.1.0"                       // Version of the Forge coding agent
	defaultModel = "anthropic/claude-sonnet-4.5" // Default model to use

	// Default system prompt for coding agent
	defaultSystemPrompt = `
You are Forge, an expert coding assistant with access to powerful tools for file operations and code editing.

You can:
- Read and write files in the current workspace
- Search code with regex patterns
- Apply precise code changes using diffs
- Execute terminal commands
- List and navigate the file system

Guidelines:
1. Always work within the current workspace directory
2. Use ApplyDiff for targeted changes instead of rewriting entire files
3. Explain your changes clearly before proposing them
4. Ask for clarification when requirements are ambiguous
5. Follow best practices for the language you're working with

When making file changes:
- Show a clear diff of what you're changing
- Explain why you're making the change
- Consider the impact on related code

Be helpful, precise, and thoughtful in your assistance.`
)

// Config holds the application configuration
type Config struct {
	APIKey       string
	BaseURL      string
	Model        string
	WorkspaceDir string
	SystemPrompt string
	ShowVersion  bool
}

func main() {
	// Parse command line flags
	config := parseFlags()

	// Show version if requested
	if config.ShowVersion {
		fmt.Printf("Forge v%s\n", version)
		return
	}

	// Validate configuration
	if err := config.validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Create context with signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n\nShutting down gracefully...")
		cancel()
	}()

	// Run the application
	if runErr := run(ctx, config); runErr != nil {
		cancel()
		log.Fatalf("Application error: %v", runErr)
	}
}

// parseFlags parses command line flags and environment variables
func parseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.APIKey, "api-key", os.Getenv("OPENAI_API_KEY"), "OpenAI API key (or set OPENAI_API_KEY env var)")
	flag.StringVar(&config.BaseURL, "base-url", os.Getenv("OPENAI_BASE_URL"), "OpenAI API base URL (or set OPENAI_BASE_URL env var)")
	flag.StringVar(&config.Model, "model", defaultModel, "LLM model to use")
	flag.StringVar(&config.WorkspaceDir, "workspace", ".", "Workspace directory (default: current directory)")
	flag.StringVar(&config.SystemPrompt, "prompt", defaultSystemPrompt, "System prompt for the agent")
	flag.BoolVar(&config.ShowVersion, "version", false, "Show version and exit")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Forge - A TUI coding agent\n\n")
		fmt.Fprintf(os.Stderr, "Usage: forge [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stderr, "  OPENAI_API_KEY     OpenAI API key\n")
		fmt.Fprintf(os.Stderr, "  OPENAI_BASE_URL    OpenAI API base URL (for compatible APIs)\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  forge                                    # Start in current directory\n")
		fmt.Fprintf(os.Stderr, "  forge -workspace /path/to/project\n")
		fmt.Fprintf(os.Stderr, "  forge -model gpt-4-turbo\n")
		fmt.Fprintf(os.Stderr, "  forge -base-url https://api.openrouter.ai/api/v1\n")
	}

	flag.Parse()
	return config
}

// validate checks that the configuration is valid
func (c *Config) validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("API key is required. Set OPENAI_API_KEY environment variable or use -api-key flag")
	}

	// Verify workspace directory exists
	info, err := os.Stat(c.WorkspaceDir)
	if err != nil {
		return fmt.Errorf("workspace directory error: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("workspace path '%s' is not a directory", c.WorkspaceDir)
	}

	return nil
}

// run executes the main application logic
func run(ctx context.Context, config *Config) error {
	// Create OpenAI provider with optional base URL
	providerOpts := []openai.ProviderOption{
		openai.WithModel(config.Model),
	}

	// Add base URL if provided
	if config.BaseURL != "" {
		providerOpts = append(providerOpts, openai.WithBaseURL(config.BaseURL))
	}

	provider, err := openai.NewProvider(
		config.APIKey,
		providerOpts...,
	)
	if err != nil {
		return fmt.Errorf("failed to create LLM provider: %w", err)
	}

	// Create agent with custom system prompt
	ag := agent.NewDefaultAgent(
		provider,
		agent.WithCustomInstructions(config.SystemPrompt),
	)

	// Create workspace security guard
	guard, err := workspace.NewGuard(config.WorkspaceDir)
	if err != nil {
		return fmt.Errorf("failed to create workspace guard: %w", err)
	}

	// Register coding tools
	codingTools := []struct {
		name string
		tool tools.Tool
	}{
		{"read_file", coding.NewReadFileTool(guard)},
		{"write_file", coding.NewWriteFileTool(guard)},
		{"list_files", coding.NewListFilesTool(guard)},
		{"search_files", coding.NewSearchFilesTool(guard)},
		{"apply_diff", coding.NewApplyDiffTool(guard)},
		{"execute_command", coding.NewExecuteCommandTool(guard)},
	}

	for _, t := range codingTools {
		if err := ag.RegisterTool(t.tool); err != nil {
			return fmt.Errorf("failed to register %s tool: %w", t.name, err)
		}
	}

	// Create TUI executor
	executor := tui.NewExecutor(ag)

	// Display welcome message
	fmt.Printf("Forge v%s - Coding Agent\n", version)
	fmt.Printf("Workspace: %s\n", config.WorkspaceDir)
	fmt.Printf("Model: %s\n", config.Model)
	fmt.Println("\nStarting TUI...")
	fmt.Println()

	// Run the executor
	if err := executor.Run(ctx); err != nil {
		return fmt.Errorf("executor error: %w", err)
	}

	return nil
}
