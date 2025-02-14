package types

type GenerateOptions struct {
	// The Messages in a given generation request
	Messages []*Message

	// The Tools available to an LLM
	Tools []*Tool

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
