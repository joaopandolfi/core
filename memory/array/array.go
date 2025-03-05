package array

import "github.com/agent-api/core"

type ArrayMemoryBackend struct {
	mem []*core.Message
}

func NewArrayMemoryBackend() *ArrayMemoryBackend {
	return &ArrayMemoryBackend{
		mem: []*core.Message{},
	}
}

func (a *ArrayMemoryBackend) Add(m ...*core.Message) error {
	a.mem = append(a.mem, m...)
	return nil
}

func (a *ArrayMemoryBackend) GetMaxN(n int) ([]*core.Message, error) {
	if n > len(a.mem) {
		n = len(a.mem)
	}

	return a.mem[:n], nil
}

func (a *ArrayMemoryBackend) Dump() ([]*core.Message, error) {
	return a.mem, nil
}

// Prune implements memory.MemoryBackend.
func (a *ArrayMemoryBackend) Prune() {}
