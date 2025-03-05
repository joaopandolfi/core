package core

// AgentRunAggregator represents a single step in an agent's execution
type AgentRunAggregator struct {
	Messages []*Message
}

func NewAgentRunAggregator() *AgentRunAggregator {
	return &AgentRunAggregator{
		Messages: []*Message{},
	}
}

func (ama *AgentRunAggregator) Push(m ...*Message) {
	ama.Messages = append(ama.Messages, m...)
}

func (ama *AgentRunAggregator) Pop() *Message {
	if len(ama.Messages) == 0 {
		return nil
	}

	return ama.Messages[len(ama.Messages)-1]
}
