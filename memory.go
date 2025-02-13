package core

import "github.com/agent-api/core/types"

// Memory interface for different memory implementations
type MemoryStorer interface {
	Add(msg types.Message)
	Get() []types.Message
	Clear()
}
