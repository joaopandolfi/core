package core

import (
	"github.com/agent-api/core/v0/message"
	"github.com/agent-api/core/v0/tool"

	"context"
)

// ProviderCapabilities represents what features a provider supports
type ProviderCapabilities struct {
	MaxTokens         int
	SupportsStreaming bool
	SupportsTools     bool
	SupportsFunctions bool
	SupportsImages    bool
	DefaultModel      string
	AvailableModels   []string
}

// Provider is the interface that all agent-api LLM providers must implement
type Provider interface {
	// GetCapabilities returns what features this provider supports
	GetCapabilities(ctx context.Context) (*ProviderCapabilities, error)

	// GenerateResponse generates a response given a context and messages
	GenerateResponse(ctx context.Context, messages []message.Message) (*message.Message, error)

	// GenerateWithTools generates a response that can use tools
	GenerateWithTools(ctx context.Context, messages []message.Message, tools []tool.Tool) (*message.Message, error)

	// GenerateStream streams the response token by token
	GenerateStream(ctx context.Context, messages []message.Message, opts *InferenceOptions) (<-chan *message.Message, <-chan error)

	// GenerateStreamWithTools streams a response with tools token by token
	GenerateStreamWithTools(ctx context.Context, messages []message.Message, tools []tool.Tool, opts *InferenceOptions) (<-chan *message.Message, <-chan error)

	// ValidatePrompt checks if a prompt is valid for the provider
	ValidatePrompt(ctx context.Context, messages []message.Message) error

	// EstimateTokens estimates the number of tokens in a message
	EstimateTokens(ctx context.Context, message string) (int, error)

	// GetModelList returns available models for this provider
	GetModelList(ctx context.Context) ([]string, error)
}

// InferenceOptions contains parameters for LLM generation
type InferenceOptions struct {
	Temperature      float64  // Controls randomness (0.0-1.0)
	TopP             float64  // Nucleus sampling parameter
	MaxTokens        int      // Maximum tokens to generate
	StopSequences    []string // Sequences that will stop generation
	Model            string   // Specific model to use (if provider supports multiple)
	Stream           bool     // Whether to stream the response
	SystemPrompt     string   // System prompt to use
	PresencePenalty  float64  // Penalty for token presence
	FrequencyPenalty float64  // Penalty for token frequency
}
