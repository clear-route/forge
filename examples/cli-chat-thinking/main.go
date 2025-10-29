package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/clear-route/forge/pkg/agent"
	"github.com/clear-route/forge/pkg/executor/cli"
	"github.com/clear-route/forge/pkg/llm/openai"
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

	// Create agent configuration with system prompt that encourages thinking
	systemPrompt := `You are a helpful AI assistant. When answering questions, you should show your thinking process.

IMPORTANT: Format your responses as follows:
1. First, wrap your reasoning/thinking in <thinking></thinking> tags
2. Then provide your final answer outside the tags

Example:
<thinking>
Let me break this down step by step:
- First consideration...
- Second point...
- Therefore...
</thinking>

Based on my analysis, the answer is...

Always use this format for complex questions that require reasoning.`

	// Create agent with options
	ag := agent.NewDefaultAgent(provider,
		agent.WithSystemPrompt(systemPrompt),
	)

	// Create CLI executor
	// Thinking is shown by default when the model includes <thinking> tags
	executor := cli.NewExecutor(ag,
		cli.WithPrompt("You: "),
	)

	fmt.Println("=== Chat with Thinking Mode ===")
	fmt.Println("The AI will show its reasoning process in [Thinking...] blocks")
	fmt.Println("Try asking: 'What is 15 * 23? Show your work.'")
	fmt.Println()

	// Run the conversation
	ctx := context.Background()
	if err := executor.Run(ctx); err != nil {
		log.Fatalf("Executor error: %v", err)
	}

	fmt.Println("Goodbye!")
}
