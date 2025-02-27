// package inmemory implements the "in memory" provider

package inmemory

import "github.com/agent-api/core/types"

// AgentMessageAggregator represents a single step in an agent's execution
type InMemoryMemStore struct {
	Messages []*types.Message
}

func NewInMemoryMemStore() *InMemoryMemStore {
	return &InMemoryMemStore{
		Messages: []*types.Message{},
	}
}

func (imms *InMemoryMemStore) Push(m ...*types.Message) {
	imms.Messages = append(imms.Messages, m...)
}

func (imms *InMemoryMemStore) Peek() *types.Message {
	if len(imms.Messages) == 0 {
		return nil
	}

	return imms.Messages[len(imms.Messages)-1]
}

func (imms *InMemoryMemStore) Dump() []*types.Message {
	return imms.Messages
}

func (imms *InMemoryMemStore) Clear() {
	imms.Messages = []*types.Message{}
}
