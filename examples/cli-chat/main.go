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

	// Create agent with options
	ag := agent.NewDefaultAgent(provider,
		agent.WithCustomInstructions("You are a helpful AI assistant. Be concise and friendly."),
	)

	// Create CLI executor
	executor := cli.NewExecutor(ag,
		cli.WithPrompt("You: "),
	)

	// Run the conversation
	ctx := context.Background()
	if err := executor.Run(ctx); err != nil {
		log.Fatalf("Executor error: %v", err)
	}

	fmt.Println("Goodbye!")
}
