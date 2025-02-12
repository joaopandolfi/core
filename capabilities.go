package core

import "context"

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

// Extended capabilities for provider-specific features
type ExtendedCapabilities interface {
	GetExtendedCapability(key string) (interface{}, bool)
	ListExtendedCapabilities() []string
}

// Core feature interfaces
type SafetyAwareProvider interface {
	SetSafeMode(ctx context.Context, enabled bool) error
	GetSafeMode(ctx context.Context) bool

	// SafetyAwareProvider should implement the core agent-api.Provider
	Provider
}
