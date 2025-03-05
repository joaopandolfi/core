package core

import (
	"context"
)

// AgentRunner interface defines the core capabilities required for an agent
type AgentRunner interface {
	// Run executes the agent's main loop with the given input until a stop condition is met
	Run(ctx context.Context, input string, stopCondition AgentStopCondition) ([]*AgentRunAggregator, error)

	// RunStream supports a streaming channel from a provider
	RunStream(ctx context.Context, input string, stopCondition AgentStopCondition) (<-chan AgentRunAggregator, <-chan string, <-chan error)

	// Step executes a single step of the agent's logic based on a given role
	Step(ctx context.Context, message Message) (*Message, error)

	// SendMessages provides a simpler interface for chat-style interactions
	SendMessages(ctx context.Context, content string) (*Message, error)

	// AddTool adds a new tool to the agent's capabilities
	AddTool(tool *Tool) error

	// GetTools returns the current set of available tools
	GetTools() []*Tool

	// Middleware functionality

	// ---------------------

	// RegisterMiddleware adds a middleware to the processing chain
	RegisterMiddleware(middle *Middleware) error

	// RemoveMiddleware removes a middleware by name
	RemoveMiddleware(name string) error

	// GetMiddleware returns a middleware by name
	GetMiddleware(name string) (Middleware, bool)

	// ListMiddleware returns all registered middleware in priority order
	ListMiddleware() []*Middleware
}
