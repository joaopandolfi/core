package memory

import "github.com/agent-api/core/message"

// Memory interface for different memory implementations
type Memory interface {
	Add(msg message.Message)
	Get() []message.Message
	Clear()
}
