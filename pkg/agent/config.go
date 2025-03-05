package agent

import (
	"github.com/go-logr/logr"

	"github.com/agent-api/core"
)

// NewAgentConfig holds configuration for agent initialization
type NewAgentConfig struct {
	// The core.Provider this agent will use
	Provider core.Provider

	// Maximum number of steps before forcing stop
	MaxSteps int

	// Initial set of tools
	Tools []core.Tool

	// Initial system prompt
	SystemPrompt string

	Logger *logr.Logger

	Memory core.MemoryBackend
}
