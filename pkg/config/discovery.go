package config

import (
	"fmt"
)

// DiscoverToolsFromAgent initializes the auto-approval section with tools from the agent
// This should be called after the agent is created and tools are registered
func DiscoverToolsFromAgent(agent interface{}) error {
	if !IsInitialized() {
		return fmt.Errorf("config not initialized")
	}

	// Get tools from agent using the GetTools interface
	type toolGetter interface {
		GetTools() []interface{}
	}

	getter, ok := agent.(toolGetter)
	if !ok {
		return fmt.Errorf("agent does not implement GetTools() method")
	}

	tools := getter.GetTools()
	if len(tools) == 0 {
		return nil // No tools to discover
	}

	// Get the auto-approval section
	section, exists := globalManager.GetSection("auto_approval")
	if !exists {
		return fmt.Errorf("auto-approval section not found")
	}

	autoApproval, ok := section.(*AutoApprovalSection)
	if !ok {
		return fmt.Errorf("auto-approval section has wrong type")
	}

	// Extract tool names and ensure they exist in the config
	type toolNamer interface {
		Name() string
	}

	for _, tool := range tools {
		if namer, ok := tool.(toolNamer); ok {
			toolName := namer.Name()
			if toolName != "" {
				autoApproval.EnsureToolExists(toolName)
			}
		}
	}

	// Save the updated configuration
	return globalManager.SaveAll()
}