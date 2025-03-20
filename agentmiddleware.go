package core

import (
	"context"
)

// Middleware defines an interface for intercepting and potentially modifying
// messages before they're processed by an agent and after they're returned.
type Middleware interface {
	// Name returns a unique identifier for a given piece of middleware
	Name() string

	// Priority determines the execution order of registered middleware
	// (0 executes first)
	//
	// SetPriority sets the priority of the middleware
	SetPriority(uint)

	// GetPriority returns the priority of the middleware
	GetPriority() uint

	// PreProcess is called before a message is sent to the agent.
	// It can modify the calling message or context.
	PreProcess(ctx context.Context, m *Message) (context.Context, *Message, error)

	// PostProcess is called after a response message is received from the agent.
	// It can modify the response message or context.
	PostProcess(ctx context.Context, m *Message) (context.Context, *Message, error)
}
