package agent

import "github.com/joaopandolfi/core"

// RegisterMiddleware adds a middleware to the processing chain
func (a *Agent) RegisterMiddleware(m *core.Middleware) error {
	return nil
}

// RemoveMiddleware removes a middleware by name
func (a *Agent) RemoveMiddleware(name string) error {
	return nil
}

// GetMiddleware returns a middleware by name
func (a *Agent) GetMiddleware(name string) (core.Middleware, bool) {
	return nil, false
}

// ListMiddleware returns all registered middleware in priority order
func (a *Agent) ListMiddleware() []*core.Middleware {
	return nil
}
