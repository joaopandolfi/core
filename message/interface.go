package message

import "encoding/json"

type MessageRole string

const (
	UserMessageRole      MessageRole = "user"
	AssistantMessageRole MessageRole = "assistant"
	ToolMessageRole      MessageRole = "tool"
)

// For timestamps, source info, etc.
type Metadata map[string]interface{}

// Message represents a single message in a conversation with multimodal support
type Message struct {
	Role MessageRole

	// Allows for mixed content types
	Content string

	// A list of base64-encoded images (for multimodal models such as llava)
	Images []string

	// Multiple tool calls
	ToolCalls []ToolCall

	// Result from tool execution
	ToolResult *ToolResult

	// Additional context
	Metadata Metadata
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
