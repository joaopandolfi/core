package tool

import "context"

// Tool represents a function that the agent can use
type Tool struct {
	Name        string
	Description string
	Function    func(ctx context.Context, args map[string]interface{}) (interface{}, error)
	Schema      map[string]interface{} // JSON Schema for the tool's parameters
}
