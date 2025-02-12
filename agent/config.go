package agent

import (
	"github.com/agent-api/core"
	"github.com/agent-api/core/tool"
)

// AgentConfig holds configuration for agent initialization
type AgentConfig struct {
	// The core.Provider this agent will use
	Provider core.Provider

	// Maximum number of steps before forcing stop
	MaxSteps int

	// Initial set of tools
	Tools []tool.Tool

	// Initial system prompt
	SystemPrompt string
}
