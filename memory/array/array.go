package array

import "github.com/joaopandolfi/core"

// ArrayMemoryBackend implements core.MemoryBackend
// with a simple, "in-memory" array of messages. This memory backend more or less
// operates like a queue where messages are stored first in, last out.
type ArrayMemoryBackend struct {
	mem []*core.Message
}

// NewArrayMemoryBackend returns a new ArrayMemoryBackend
func NewArrayMemoryBackend() *ArrayMemoryBackend {
	return &ArrayMemoryBackend{
		mem: []*core.Message{},
	}
}

// Add adds messages to the ArrayMemoryBackend using "append"
func (a *ArrayMemoryBackend) Add(m ...*core.Message) error {
	a.mem = append(a.mem, m...)
	return nil
}

// GetMaxN returns the last N number of messages
func (a *ArrayMemoryBackend) GetMaxN(n int) ([]*core.Message, error) {
	if n > len(a.mem) {
		n = len(a.mem)
	}

	return a.mem[:n], nil
}

// Dump returns the whole ArrayMemoryBackend array
func (a *ArrayMemoryBackend) Dump() ([]*core.Message, error) {
	return a.mem, nil
}

// Prune resets the array in the ArrayMemoryBackend
func (a *ArrayMemoryBackend) Prune() {
	a.mem = []*core.Message{}
}
