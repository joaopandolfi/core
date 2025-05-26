package agent

import "github.com/joaopandolfi/core"

// AgentRunAggregator represents a single step in an agent's execution
type AgentRunAggregator struct {
	Messages []*core.Message
}

func NewAgentRunAggregator() *AgentRunAggregator {
	return &AgentRunAggregator{
		Messages: []*core.Message{},
	}
}

func (ama *AgentRunAggregator) Push(m ...*core.Message) {
	ama.Messages = append(ama.Messages, m...)
}

func (ama *AgentRunAggregator) Pop() *core.Message {
	if len(ama.Messages) == 0 {
		return nil
	}

	return ama.Messages[len(ama.Messages)-1]
}
