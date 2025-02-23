package core

import (
	"context"

	"github.com/agent-api/core/types"
)

// Provider is the interface that all agent-api LLM providers must implement.
type Provider interface {
	// GetCapabilities returns what features this provider supports through a
	// core.Capabilities struct. A provider may return an error if it cannot
	// construct or query for its capabilities.
	GetCapabilities(ctx context.Context) (*types.Capabilities, error)

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
	UseModel(ctx context.Context, model *types.Model) error

	// Generate uses the provider to generate a new message given the core.GenerateOptions
	Generate(ctx context.Context, opts *types.GenerateOptions) (*types.Message, error)

	// Generate uses the provider to stream messages. It returns:
	// * a *types.Message channel which should have complete messages to be consumed by providers.
	// * a string channel which are the streaming deltas from the provider.
	// * an error channel to surface any errors during streaming execution.j
	GenerateStream(ctx context.Context, opts *types.GenerateOptions) (<-chan *types.Message, <-chan string, <-chan error)
}
