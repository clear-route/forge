package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/entrhq/forge/pkg/agent"
	"github.com/entrhq/forge/pkg/executor/cli"
	"github.com/entrhq/forge/pkg/llm/openai"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Get optional base URL
	baseURL := os.Getenv("OPENAI_BASE_URL")

	// Create OpenAI provider
	var provider *openai.Provider
	var err error

	if baseURL != "" {
		fmt.Printf("Using custom OpenAI endpoint: %s\n", baseURL)
		provider, err = openai.NewProvider(apiKey, openai.WithModel("openai/gpt-4o"))
	} else {
		provider, err = openai.NewProvider(apiKey, openai.WithModel("gpt-4o"))
	}

	if err != nil {
		log.Fatalf("Failed to create OpenAI provider: %v", err)
	}

	// Create agent with custom instructions
	// Note: The agent loop system now automatically handles thinking and tool calls
	// You can provide additional instructions here that will be added to the base system prompt
	customInstructions := `You are a helpful AI assistant.
When solving problems, think step by step and be thorough in your reasoning.`

	// Create agent with options
	ag := agent.NewDefaultAgent(provider,
		agent.WithCustomInstructions(customInstructions),
	)

	// Create CLI executor
	// Thinking is shown by default when the model includes <thinking> tags
	executor := cli.NewExecutor(ag,
		cli.WithPrompt("You: "),
	)

	fmt.Println("=== Agent Loop Chat with Thinking ===")
	fmt.Println("The agent automatically shows its reasoning in [Thinking...] blocks")
	fmt.Println("The agent uses tools to complete tasks or ask questions")
	fmt.Println("Try: 'What is 15 * 23?' or 'Help me write a haiku'")
	fmt.Println()

	// Run the conversation
	ctx := context.Background()
	if err := executor.Run(ctx); err != nil {
		log.Fatalf("Executor error: %v", err)
	}

	fmt.Println("Goodbye!")
}
