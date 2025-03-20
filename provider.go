package core

import (
	"context"
)

// Provider is the interface that all agent-api LLM providers must implement.
type Provider interface {
	// GetCapabilities returns what features this provider supports through a
	// core.Capabilities struct. A provider may return an error if it cannot
	// construct or query for its capabilities.
	GetCapabilities(ctx context.Context) (*Capabilities, error)

	// UseModel takes a context and a core.Model struct supported by the provider.
	// It returns an error if something went wrong with setting the model in
	// the provider.
	UseModel(ctx context.Context, model *Model) error

	// Generate uses the provider to generate a new message given the core.GenerateOptions
	Generate(ctx context.Context, opts *GenerateOptions) (*Message, error)

	// GenerateStream uses the provider to stream messages. It returns:
	//
	// * a *core.Message channel which should have complete messages to be consumed
	//   from the provider. I.e., these are full, complete messages.
	// * a string channel which are the streaming deltas from the provider. These
	//   are not full messages nor complete chunks: they may be only one or two words
	//   and the message deltas are provider specific.
	// * an error channel to surface any errors during streaming execution.
	GenerateStream(ctx context.Context, opts *GenerateOptions) (<-chan *Message, <-chan string, <-chan error)
}
