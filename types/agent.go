package types

// AgentRunAggregator represents a single step in an agent's execution
type AgentRunAggregator struct {
	Messages []*Message

	// Any error that occurred during this step
	Err error
}

func NewAgentRunAggregator() *AgentRunAggregator {
	return &AgentRunAggregator{
		Messages: []*Message{},
		Err:      nil,
	}
}

func (ama *AgentRunAggregator) Push(e error, m ...*Message) {
	ama.Messages = append(ama.Messages, m...)
	ama.Err = e
}

func (ama *AgentRunAggregator) Pop() (*Message, error) {
	if len(ama.Messages) == 0 {
		return nil, nil
	}

	return ama.Messages[len(ama.Messages)-1], ama.Err
}

// StopCondition is a function that determines if the agent should stop
// after its completed a step (i.e., a full "start" -> "doing work" -> "done" cycle)
type AgentStopCondition func(step *AgentRunAggregator) bool
