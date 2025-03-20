package core

// Capabilities represents what features an LLM provider supports
type Capabilities struct {
	SupportsCompletion bool
	SupportsChat       bool
	SupportsStreaming  bool
	SupportsTools      bool
	SupportsImages     bool
	DefaultModel       string

	// A provider should return the available models structs that can be iterated
	// by a downstream consumer.
	AvailableModels []*Model
}
