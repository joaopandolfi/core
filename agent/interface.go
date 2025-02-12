package agent

import (
	"context"

	"github.com/agent-api/core/message"
	"github.com/agent-api/core/tool"
)

// StopCondition is a function that determines if the agent should stop
type StopCondition func(step *AgentStep) bool

// AgentRunner interface defines the core capabilities required for an agent
type AgentRunner interface {
	// Run executes the agent's main loop with the given input until a stop condition is met
	Run(ctx context.Context, input string, stopCondition StopCondition) ([]AgentStep, error)

	// Step executes a single step of the agent's logic based on a given role
	Step(ctx context.Context, message message.Message) (*AgentStep, error)

	// SendMessage provides a simpler interface for chat-style interactions
	SendMessage(ctx context.Context, content string) (*message.Message, error)

	// AddTool adds a new tool to the agent's capabilities
	AddTool(tool tool.Tool) error

	// GetTools returns the current set of available tools
	GetTools() []tool.Tool
}

// AgentStreamingRunner supports streaming responses
type AgentStreamingRunner interface {
	RunStream(ctx context.Context, input string, stopCondition StopCondition) (<-chan AgentStep, <-chan error)

	AgentRunner
}

// AgentReflectiveRunner supports self-reflection and planning
type AgentReflectiveRunner interface {
	Plan(ctx context.Context, goal string) ([]string, error) // Returns planned steps
	Reflect(ctx context.Context, steps []AgentStep) (*message.Message, error)

	AgentRunner
}
