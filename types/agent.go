package types

// AgentStep represents a single step in an agent's execution
type AgentStep struct {
	ID string

	Message *Message

	// Any error that occurred during this step
	Error error
}

// StopCondition is a function that determines if the agent should stop
type AgentStopCondition func(step *AgentStep) bool
