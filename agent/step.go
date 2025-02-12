package agent

import (
	"github.com/agent-api/core/message"
)

// AgentStep represents a single step in an agent's execution
type AgentStep struct {
	ID string

	Message *message.Message

	// Any error that occurred during this step
	Error error
}
