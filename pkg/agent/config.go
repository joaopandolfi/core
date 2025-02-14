package agent

import (
	"log/slog"

	"github.com/agent-api/core"
	"github.com/agent-api/core/types"
)

// NewAgentConfig holds configuration for agent initialization
type NewAgentConfig struct {
	// The core.Provider this agent will use
	Provider core.Provider

	// Maximum number of steps before forcing stop
	MaxSteps int

	// Initial set of tools
	Tools []types.Tool

	// Initial system prompt
	SystemPrompt string

	Logger *slog.Logger
}
