package amnesic

import "github.com/joaopandolfi/core"

// AmnesicMemory implements core.MemoryBackend
// amnesic memory does not store messages, only implements the interface
type AmnesicMemory struct {
	mem []*core.Message
}

// NewAmnesicMemory returns a new AmnesicMemory
func NewAmnesicMemory() *AmnesicMemory {
	return &AmnesicMemory{
		mem: []*core.Message{},
	}
}

// update whole memory - amnesic behavior
func (a *AmnesicMemory) Add(m ...*core.Message) error {
	a.mem = m
	return nil
}

// GetMaxN returns the last N number of messages
func (a *AmnesicMemory) GetMaxN(n int) ([]*core.Message, error) {
	if n > len(a.mem) {
		n = len(a.mem)
	}
	return a.mem[:n], nil
}

// returns a empty array of messages
func (a *AmnesicMemory) Dump() ([]*core.Message, error) {
	return a.mem, nil
}

// Prune resets the pointer array
func (a *AmnesicMemory) Prune() {
	a.mem = []*core.Message{}
}
