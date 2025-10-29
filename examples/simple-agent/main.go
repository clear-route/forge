package main

import (
	"fmt"
	"log"

	"github.com/clear-route/forge/pkg/types"
)

func main() {
	fmt.Println("Forge Simple Agent Example")
	log.Println("This example demonstrates the core types and interfaces")

	// Demonstrate Message types
	fmt.Println("\n--- Message Types ---")
	systemMsg := types.NewSystemMessage("You are a helpful AI assistant")
	userMsg := types.NewUserMessage("Hello, how are you?")
	assistantMsg := types.NewAssistantMessage("I'm doing well, thank you!")

	fmt.Printf("System: %s\n", systemMsg.Content)
	fmt.Printf("User: %s\n", userMsg.Content)
	fmt.Printf("Assistant: %s\n", assistantMsg.Content)

	// Demonstrate Message metadata
	userMsg.WithMetadata("user_id", "user123").WithMetadata("session_id", "session456")
	fmt.Printf("User message metadata: %+v\n", userMsg.Metadata)

	// Demonstrate Input types
	fmt.Println("\n--- Input Types ---")
	userInput := types.NewUserInput("What is the weather like?")
	cancelInput := types.NewCancelInput()
	formInput := types.NewFormInput(map[string]string{
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   "30",
	})

	fmt.Printf("User Input: %s (type: %s, isUserInput: %v)\n",
		userInput.Content, userInput.Type, userInput.IsUserInput())
	fmt.Printf("Cancel Input: type=%s, isCancel=%v\n",
		cancelInput.Type, cancelInput.IsCancel())
	fmt.Printf("Form Input: type=%s, isFormInput=%v, fields=%d\n",
		formInput.Type, formInput.IsFormInput(), len(formInput.FormData))
	for key, value := range formInput.FormData {
		fmt.Printf("  - %s: %s\n", key, value)
	}

	// Demonstrate AgentConfig
	fmt.Println("\n--- Agent Configuration ---")
	config := types.NewAgentConfig().
		WithSystemPrompt("You are a helpful assistant").
		WithMaxTurns(10).
		WithStreaming(true).
		WithBufferSize(20)

	fmt.Printf("Config - Streaming: %v, MaxTurns: %d, BufferSize: %d\n",
		config.EnableStreaming, config.MaxTurns, config.BufferSize)

	// Demonstrate AgentChannels
	fmt.Println("\n--- Agent Channels ---")
	channels := types.NewAgentChannels(config.BufferSize)
	fmt.Printf("Channels created with buffer size: %d\n", config.BufferSize)
	fmt.Printf("- Input channel capacity: %d\n", cap(channels.Input))
	fmt.Printf("- Event channel capacity: %d\n", cap(channels.Event))

	// Demonstrate AgentEvent types
	fmt.Println("\n--- Agent Events (Streaming Support) ---")

	// Thinking events
	thinkingStart := types.NewThinkingStartEvent()
	thinkingContent := types.NewThinkingContentEvent("Analyzing the request...")
	thinkingEnd := types.NewThinkingEndEvent()
	fmt.Printf("Thinking: start=%s, content=%s, end=%s\n",
		thinkingStart.Type, thinkingContent.Type, thinkingEnd.Type)

	// Message events (for streaming LLM responses)
	messageStart := types.NewMessageStartEvent()
	messageContent1 := types.NewMessageContentEvent("Hello ")
	messageContent2 := types.NewMessageContentEvent("world!")
	messageEnd := types.NewMessageEndEvent()
	fmt.Printf("Message streaming: %s -> %s ('%s') -> %s ('%s') -> %s\n",
		messageStart.Type, messageContent1.Type, messageContent1.Content,
		messageContent2.Type, messageContent2.Content, messageEnd.Type)

	// Tool events
	toolCall := types.NewToolCallEvent("weather_api", map[string]interface{}{
		"location": "San Francisco",
	})
	toolResult := types.NewToolResultEvent("weather_api", "Sunny, 72°F")
	fmt.Printf("Tool: %s (tool=%s) -> %s\n", toolCall.Type, toolCall.ToolName, toolResult.Type)

	// Status events
	turnEnd := types.NewTurnEndEvent()
	fmt.Printf("Turn complete: %s\n", turnEnd.Type)

	// Demonstrate AgentError
	fmt.Println("\n--- Error Handling ---")
	err := types.NewAgentError(types.ErrorCodeLLMFailure, "Failed to connect to LLM")
	err.WithMetadata("provider", "openai").WithMetadata("model", "gpt-4")

	fmt.Printf("Error: %v\n", err)
	fmt.Printf("Error code: %s\n", err.Code)
	fmt.Printf("Error metadata: %+v\n", err.Metadata)

	// Demonstrate ModelInfo
	fmt.Println("\n--- Model Info ---")
	modelInfo := &types.ModelInfo{
		Name:              "gpt-4",
		MaxTokens:         8192,
		SupportsStreaming: true,
		Provider:          "openai",
		Metadata:          make(map[string]interface{}),
	}

	fmt.Printf("Model: %s\n", modelInfo.Name)
	fmt.Printf("Provider: %s\n", modelInfo.Provider)
	fmt.Printf("Max Tokens: %d\n", modelInfo.MaxTokens)
	fmt.Printf("Supports Streaming: %v\n", modelInfo.SupportsStreaming)

	fmt.Println("\n--- Framework Ready ---")
	log.Println("All core types and interfaces are defined and ready for implementation!")
	log.Println("Next steps:")
	log.Println("  1. Implement a concrete Agent")
	log.Println("  2. Implement an Executor (CLI, API, etc.)")
	log.Println("  3. Implement an LLM Provider")
	log.Println("  4. Connect them together for a working agent!")
}
