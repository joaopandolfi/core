package types

// Capabilities represents what features an LLM provider supports
type Capabilities struct {
	SupportsCompletion bool
	SupportsChat       bool
	SupportsStreaming  bool
	SupportsTools      bool
	SupportsImages     bool
	DefaultModel       string
	AvailableModels    []string
}
