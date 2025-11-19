package agent

import (
	"github.com/entrhq/forge/pkg/agent/tools"
)

// getToolsList returns tools as []tools.Tool for internal use
func (a *DefaultAgent) getToolsList() []tools.Tool {
	a.toolsMu.RLock()
	defer a.toolsMu.RUnlock()

	toolsList := make([]tools.Tool, 0, len(a.tools))
	for _, tool := range a.tools {
		toolsList = append(toolsList, tool)
	}
	return toolsList
}

// getTool retrieves a tool by name (thread-safe)
func (a *DefaultAgent) getTool(name string) (tools.Tool, bool) {
	a.toolsMu.RLock()
	defer a.toolsMu.RUnlock()

	tool, exists := a.tools[name]
	return tool, exists
}
