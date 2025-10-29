package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/clear-route/forge/pkg/llm/openai"
	"github.com/clear-route/forge/pkg/types"
)

func main() {
	// Check if API key is available
	apiKey := os.Getenv("OPENAI_API_KEY")
	baseURL := os.Getenv("OPENAI_BASE_URL")

	if apiKey == "" {
		fmt.Println("⚠️  No OPENAI_API_KEY environment variable set")
		fmt.Println("Set it with: export OPENAI_API_KEY=your-key-here")
		fmt.Println("\nOptionally set custom base URL:")
		fmt.Println("  export OPENAI_BASE_URL=https://your-custom-endpoint.com")
		fmt.Println("\nRunning in demo mode (will show structure without making API calls)...")
		demoMode()
		return
	}

	if baseURL != "" {
		fmt.Printf("✅ Using custom base URL: %s\n", baseURL)
	}
	fmt.Println("✅ OpenAI API key found, running live test...")
	liveTest(apiKey)
}

func demoMode() {
	fmt.Println("\n=== Demo Mode: Provider Structure ===\n")

	// Create provider (will fail without key, but shows structure)
	provider, err := openai.NewProvider("", openai.WithModel("gpt-4"))
	if err != nil {
		fmt.Printf("Expected error (no API key): %v\n\n", err)
	}

	if provider != nil {
		modelInfo := provider.GetModelInfo()
		fmt.Printf("Model Info:\n")
		fmt.Printf("  Provider: %s\n", modelInfo.Provider)
		fmt.Printf("  Name: %s\n", modelInfo.Name)
		fmt.Printf("  Supports Streaming: %v\n", modelInfo.SupportsStreaming)
		fmt.Printf("  Max Tokens: %d\n", modelInfo.MaxTokens)
	}

	fmt.Println("\n=== Sample Usage ===\n")
	fmt.Println("```go")
	fmt.Println(`// Standard OpenAI`)
	fmt.Println(`provider, _ := openai.NewProvider("sk-...", openai.WithModel("gpt-4"))`)
	fmt.Println()
	fmt.Println(`// Custom base URL (Azure, local, etc.)`)
	fmt.Println(`provider, _ := openai.NewProvider("key",`)
	fmt.Println(`    openai.WithBaseURL("https://custom.openai.azure.com"),`)
	fmt.Println(`    openai.WithModel("gpt-4"))`)
	fmt.Println()
	fmt.Println(`// Use the provider`)
	fmt.Println(`messages := []*types.Message{types.NewUserMessage("Hello!")}`)
	fmt.Println(`stream, _ := provider.StreamCompletion(ctx, messages)`)
	fmt.Println(`for chunk := range stream {`)
	fmt.Println(`    fmt.Print(chunk.Content)`)
	fmt.Println(`}`)
	fmt.Println("```")
}

func liveTest(apiKey string) {
	ctx := context.Background()

	// Create OpenAI provider
	fmt.Println("\n1. Creating OpenAI provider...")
	provider, err := openai.NewProvider(apiKey, openai.WithModel("openai/gpt-4o"))
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}

	// Show model info
	modelInfo := provider.GetModelInfo()
	fmt.Printf("   ✓ Provider: %s\n", modelInfo.Provider)
	fmt.Printf("   ✓ Model: %s\n", modelInfo.Name)
	fmt.Printf("   ✓ Supports Streaming: %v\n", modelInfo.SupportsStreaming)
	if baseURL, ok := modelInfo.Metadata["base_url"].(string); ok {
		fmt.Printf("   ✓ Base URL: %s\n", baseURL)
	}
	fmt.Println()

	// Create test messages
	fmt.Println("2. Creating test conversation...")
	messages := []*types.Message{
		types.NewSystemMessage("You are a helpful assistant. Keep responses very brief (1-2 sentences)."),
		types.NewUserMessage("What is 2+2? Just give the number."),
	}
	fmt.Printf("   ✓ System message set\n")
	fmt.Printf("   ✓ User message: %s\n\n", messages[1].Content)

	// Test streaming
	fmt.Println("3. Testing StreamCompletion()...")
	stream, err := provider.StreamCompletion(ctx, messages)
	if err != nil {
		log.Fatalf("Failed to start stream: %v", err)
	}

	fmt.Print("   Response: ")
	var fullResponse string
	chunkCount := 0

	for chunk := range stream {
		if chunk.IsError() {
			log.Fatalf("\n   ✗ Stream error: %v", chunk.Error)
		}

		if chunk.HasContent() {
			fmt.Print(chunk.Content)
			fullResponse += chunk.Content
			chunkCount++
		}

		if chunk.IsLast() {
			fmt.Println()
			break
		}
	}

	fmt.Printf("   ✓ Received %d chunks\n", chunkCount)
	fmt.Printf("   ✓ Full response: %q\n\n", fullResponse)

	// Test non-streaming
	fmt.Println("4. Testing Complete() (non-streaming)...")
	messages = append(messages, types.NewUserMessage("What is 3+3? Just the number."))

	response, err := provider.Complete(ctx, messages)
	if err != nil {
		log.Fatalf("Failed to complete: %v", err)
	}

	fmt.Printf("   ✓ Response role: %s\n", response.Role)
	fmt.Printf("   ✓ Response content: %q\n\n", response.Content)

	fmt.Println("✅ All tests passed!")
	fmt.Println("\nThe OpenAI provider is working correctly.")
	fmt.Println("Next steps: Implement DefaultAgent to wrap this provider.")
}
