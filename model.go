package core

// Model is a metadata type that providers or end users can define for use during
// message generation.
type Model struct {
	// The raw string ID of the model (i.e., "qwen2.5:latest")
	ID string

	// The configured maxiumum tokens to use during execution
	MaxTokens int
}
