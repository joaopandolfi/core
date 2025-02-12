package core

import (
	"context"

	"github.com/agent-api/core/message"
	"github.com/agent-api/core/model"
	"github.com/agent-api/core/tool"
)

type GenerateOptions struct {
	// The Messages in a given generation request
	Messages []message.Message

	// The Tools available to an LLM
	Tools []tool.Tool

	// Controls generation randomness (0.0-1.0)
	Temperature float64

	// Nucleus sampling parameter
	TopP float64

	// Maximum tokens to generate
	MaxTokens int

	// Sequences that will stop generation
	StopSequences []string

	// Penalty for token presence
	PresencePenalty float64

	// Penalty for token frequency
	FrequencyPenalty float64
}

// Provider is the interface that all agent-api LLM providers must implement.
type Provider interface {
	// GetCapabilities returns what features this provider supports through a
	// core.Capabilities struct. A provider may return an error if it cannot
	// construct or query for its capabilities.
	GetCapabilities(ctx context.Context) (*Capabilities, error)

	// UseModel takes a context and a model string ID (i.e., "qwen2.5") and configuration
	// options through a core.ModelKnobs struct. It returns:
	//
	// 1. An "ok" boolean defining if the provider supports the given model by
	//    the given options.
	// 2. The constructed core.Model itself.
	// 3. An error.
	//
	// A provider implementation may choose to return (true, nil, error) where
	// some pre-check, pre-authentication, or query to the API failed causing an
	// error despite the Model itself being supported.
	UseModel(ctx context.Context, model *model.Model) error

	// Generate uses the provider to generate a new message given the core.GenerateOptions
	Generate(ctx context.Context, opts *GenerateOptions) (*message.Message, error)

	// Generate uses the provider to stream a new message channel given the core.GenerateOptions
	GenerateStream(ctx context.Context, opts *GenerateOptions) (<-chan *message.Message, <-chan error)
}
