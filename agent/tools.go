package agent

import "github.com/agent-api/core"

type ToolMap map[string]*core.Tool

// GetTools returns the current set of available tools
func (a *Agent) GetTools() []*core.Tool {
	return nil
}
