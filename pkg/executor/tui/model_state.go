package tui

// GetWidth returns the current width
func (m *model) GetWidth() int {
	return m.width
}

// GetHeight returns the current height
func (m *model) GetHeight() int {
	return m.height
}

// IsThinking returns whether the agent is currently thinking
func (m *model) IsThinking() bool {
	return m.isThinking
}

// IsAgentBusy returns whether the agent is currently busy
func (m *model) IsAgentBusy() bool {
	return m.agentBusy
}

// GetWorkspaceDir returns the current workspace directory
func (m *model) GetWorkspaceDir() string {
	return m.workspaceDir
}

// IsBashMode returns whether bash mode is active
func (m *model) IsBashMode() bool {
	return m.bashMode
}
