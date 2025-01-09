package message

type MessageRole string

const (
	UserMessageRole MessageRole = "user"
	AIMessageRole   MessageRole = "ai"
)

// Message represents a single message in a conversation
type Message struct {
	Role    MessageRole
	Content string
}
