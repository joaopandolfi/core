package database

import "github.com/joaopandolfi/core"

// GormMemoryBackend implements core.MemoryBackend
// Uses GORM connection to retrieve memory from Postgres database
type GormMemoryBackend struct {
	mem []*core.Message
}

// NewGormMemoryBackend returns a new GormMemoryBackend
func NewGormMemoryBackend() *GormMemoryBackend {
	return &GormMemoryBackend{
		mem: []*core.Message{},
	}
}

// Add adds messages to the GormMemoryBackend using "append"
func (a *GormMemoryBackend) Add(m ...*core.Message) error {
	a.mem = append(a.mem, m...)
	return nil
}

// GetMaxN returns the last N number of messages
func (a *GormMemoryBackend) GetMaxN(n int) ([]*core.Message, error) {
	if n > len(a.mem) {
		n = len(a.mem)
	}

	return a.mem[:n], nil
}

// Dump returns the whole GormMemoryBackend array
func (a *GormMemoryBackend) Dump() ([]*core.Message, error) {
	return a.mem, nil
}

// Prune resets the array in the GormMemoryBackend
func (a *GormMemoryBackend) Prune() {
	a.mem = []*core.Message{}
}
