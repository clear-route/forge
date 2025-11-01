package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/entrhq/forge/pkg/agent"
	"github.com/entrhq/forge/pkg/executor/cli"
	"github.com/entrhq/forge/pkg/llm/openai"
)

// CalculatorTool is a custom tool that performs basic arithmetic
type CalculatorTool struct{}

func (t *CalculatorTool) Name() string {
	return "calculator"
}

func (t *CalculatorTool) Description() string {
	return "Performs basic arithmetic operations (add, subtract, multiply, divide)"
}

func (t *CalculatorTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"operation": map[string]interface{}{
				"type":        "string",
				"description": "The operation to perform: add, subtract, multiply, or divide",
				"enum":        []string{"add", "subtract", "multiply", "divide"},
			},
			"a": map[string]interface{}{
				"type":        "number",
				"description": "First number",
			},
			"b": map[string]interface{}{
				"type":        "number",
				"description": "Second number",
			},
		},
		"required": []string{"operation", "a", "b"},
	}
}

func (t *CalculatorTool) Execute(ctx context.Context, arguments json.RawMessage) (string, error) {
	var args struct {
		Operation string  `json:"operation"`
		A         float64 `json:"a"`
		B         float64 `json:"b"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	var result float64
	switch args.Operation {
	case "add":
		result = args.A + args.B
	case "subtract":
		result = args.A - args.B
	case "multiply":
		result = args.A * args.B
	case "divide":
		if args.B == 0 {
			return "", fmt.Errorf("division by zero")
		}
		result = args.A / args.B
	default:
		return "", fmt.Errorf("unknown operation: %s", args.Operation)
	}

	return fmt.Sprintf("%.2f", result), nil
}

func (t *CalculatorTool) IsLoopBreaking() bool {
	return false // This tool doesn't end the conversation
}

func main() {
	// Get API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Create OpenAI provider
	provider, err := openai.NewProvider(apiKey,
		openai.WithModel("openai/gpt-4o"),
	)
	if err != nil {
		log.Fatalf("Failed to create OpenAI provider: %v", err)
	}

	// Create agent with custom instructions
	ag := agent.NewDefaultAgent(provider,
		agent.WithCustomInstructions("You are a helpful AI assistant with access to tools. You can help with calculations and conversations."),
	)

	// Register custom calculator tool
	calculator := &CalculatorTool{}
	if err := ag.RegisterTool(calculator); err != nil {
		log.Fatalf("Failed to register calculator tool: %v", err)
	}

	fmt.Println("🤖 Agent Chat Example")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Features:")
	fmt.Println("  • Agent loop with tool execution")
	fmt.Println("  • Chain-of-thought reasoning (shown in brackets)")
	fmt.Println("  • Custom calculator tool")
	fmt.Println("  • Built-in tools: task_completion, ask_question, converse")
	fmt.Println()
	fmt.Println("Try asking:")
	fmt.Println("  • What is 15 * 23?")
	fmt.Println("  • Calculate (100 + 50) / 3")
	fmt.Println("  • What's the square root of 144? (then add 5)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Create CLI executor
	executor := cli.NewExecutor(ag,
		cli.WithPrompt("You: "),
	)

	// Run the conversation
	ctx := context.Background()
	if err := executor.Run(ctx); err != nil {
		log.Fatalf("Executor error: %v", err)
	}

	fmt.Println("\nGoodbye!")
}
