package agent

import (
	"context"
	"fmt"
	"sort"

	"github.com/joaopandolfi/core"
)

// RegisterMiddleware adds a middleware to the processing chain
func (a *Agent) RegisterMiddleware(m core.Middleware) error {
	if m == nil {
		return fmt.Errorf("cannot register a nil middeware")
	}
	if a.middlewares == nil {
		a.middlewares = map[string]core.Middleware{}
	}
	a.middlewares[m.Name()] = m
	return nil
}

// RemoveMiddleware removes a middleware by name
func (a *Agent) RemoveMiddleware(name string) error {
	if a.middlewares == nil {
		return fmt.Errorf("middewares not intialized")
	}

	delete(a.middlewares, name)
	return nil
}

// GetMiddleware returns a middleware by name
func (a *Agent) GetMiddleware(name string) (core.Middleware, bool) {
	if a.middlewares == nil {
		return nil, false
	}
	v, ok := a.middlewares[name]
	return v, ok
}

// ListMiddleware returns all registered middleware in priority order
func (a *Agent) ListMiddleware() []core.Middleware {
	result := []core.Middleware{}
	for _, v := range a.middlewares {
		result = append(result, v)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].GetPriority() > result[j].GetPriority()
	})
	return result
}

func (a *Agent) PreProcessMessage(ctx context.Context, m *core.Message) error {
	orderedMiddleware := a.ListMiddleware()
	for _, md := range orderedMiddleware {
		err := md.PreProcess(ctx, m)
		if err != nil {
			return fmt.Errorf("pre processing message on middleware %s: %w", md.Name(), err)
		}
	}
	return nil
}

func (a *Agent) PosProcessMessage(ctx context.Context, m *core.Message) error {
	orderedMiddleware := a.ListMiddleware()
	for _, md := range orderedMiddleware {
		err := md.PostProcess(ctx, m)
		if err != nil {
			return fmt.Errorf("pos processing message on middleware %s: %w", md.Name(), err)
		}
	}

	return nil
}
