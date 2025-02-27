package core

import (
	"context"

	"github.com/agent-api/core/types"
)

// Middleware defines an interface for intercepting and potentially modifying
// messages before they're processed by an agent and after they're returned
type Middleware interface {
	// Name returns a unique identifier for a given piece of middleware
	Name() string

	// Priority determines the execution order of registered middleware
	// (0 executes first)
	Priority() uint

	// PreProcess is called before a message is sent to the agent.
	// It can modify the message or context.
	PreProcess(ctx context.Context, m *types.Message) (context.Context, *types.Message, error)

	// PostProcess is called after a response message is received from the agent.
	// It can modify the response message or context.
	PostProcess(ctx context.Context, m *types.Message) (context.Context, *types.Message, error)
}
