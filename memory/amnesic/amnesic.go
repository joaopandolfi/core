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

// Do nothing
func (a *AmnesicMemory) Add(m ...*core.Message) error {
	return nil
}

// returns a empty array of messages
func (a *AmnesicMemory) GetMaxN(n int) ([]*core.Message, error) {
	return a.Dump()
}

// returns a empty array of messages
func (a *AmnesicMemory) Dump() ([]*core.Message, error) {
	return a.mem, nil
}

// Prune resets the pointer array
func (a *AmnesicMemory) Prune() {
	a.mem = []*core.Message{}
}
