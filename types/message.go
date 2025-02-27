package types

import (
	"encoding/json"
	"time"
)

type MessageRole string

const (
	UserMessageRole      MessageRole = "user"
	AssistantMessageRole MessageRole = "assistant"
	SystemMessageRole    MessageRole = "system"
	ToolMessageRole      MessageRole = "tool"
)

// Message represents a single message in a conversation with multimodal support
type Message struct {
	// ID is the incrementing internal integer identifier
	ID uint32

	// The role of the message sender
	Role MessageRole

	// The primary content of the message (usually text)
	Content string

	// A list of base64-encoded images (for multimodal models such as llava
	// or llama3.2-vision)
	Images []*Image

	// Multiple tool calls
	ToolCalls []*ToolCall

	// Result from tool execution
	ToolResult []*ToolResult

	// Additional context
	Metadata *Metadata
}

// For timestamps, source info, etc.
type Metadata struct {
	Timestamp time.Time
	Source    string
	RequestID int

	ProviderProperties map[string]string
}

// ToolCall represents a specific tool invocation request
type ToolCall struct {
	// Unique identifier for tracking
	ID string

	// Name of the tool being called
	Name string

	// Structured arguments
	Arguments json.RawMessage
}

// ToolResult contains the output of a tool execution
type ToolResult struct {
	// Reference to original call
	ToolCallID string

	// Structured result
	Content interface{}

	Error string
}
