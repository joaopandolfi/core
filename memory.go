package core

type MemoryBackend interface {
	// Add adds any number of messages to the memory storer backend
	Add(m ...*Message) error

	// GetMaxN gets the backend's last N messages
	GetMaxN(n int) ([]*Message, error)

	// Dump gets the backends last dump of messages
	Dump() ([]*Message, error)

	// Prune
	Prune()
}
