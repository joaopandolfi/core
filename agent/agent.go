package agent

import (
	"context"
	"errors"

	"github.com/agent-api/core"
	"github.com/agent-api/core/message"
	"github.com/agent-api/core/tool"
)

// Agent represents an AI agent with its configuration and state
type Agent struct {
	provider core.Provider
	tools    ToolMap
	memory   []message.Message
}

type ToolMap map[string]tool.Tool

// TODO - add adjunct functions for working with the tool map

// NewAgent creates a new agent with the given provider
func NewAgent(provider core.Provider) *Agent {
	return &Agent{
		provider: provider,
		tools:    make(map[string]tool.Tool),
		memory:   make([]message.Message, 0),
	}
}

// AddTool adds a tool to the agent's available tools
func (a *Agent) AddTool(tool tool.Tool) error {
	if tool.Name == "" {
		return errors.New("tool must have a name")
	}

	if tool.Function == nil {
		return errors.New("tool must have a function")
	}

	a.tools[tool.Name] = tool

	return nil
}

// SendMessage sends a message to the agent and gets a response
func (a *Agent) SendMessage(ctx context.Context, content string) (*message.Message, error) {
	userMsg := message.Message{
		Role:    "user",
		Content: content,
	}
	a.memory = append(a.memory, userMsg)

	var response *message.Message
	var err error

	// If we have tools, use them
	if len(a.tools) > 0 {
		toolSlice := make([]tool.Tool, 0, len(a.tools))
		for _, tool := range a.tools {
			toolSlice = append(toolSlice, tool)
		}
		response, err = a.provider.GenerateWithTools(ctx, a.memory, toolSlice)
	} else {
		response, err = a.provider.GenerateResponse(ctx, a.memory)
	}

	if err != nil {
		return nil, err
	}

	a.memory = append(a.memory, *response)
	return response, nil
}

// ClearMemory clears the agent's conversation history
func (a *Agent) ClearMemory() {
	a.memory = make([]message.Message, 0)
}
