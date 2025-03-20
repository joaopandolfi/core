package agent

import (
	"context"

	"github.com/agent-api/core"
)

// Step executes a single step of the agent's logic based on a given role
func (a *Agent) Step(ctx context.Context, message core.Message) (*core.Message, error) {
	return nil, nil
}
