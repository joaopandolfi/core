package agent

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/agent-api/core"
	"github.com/agent-api/core/types"
)

// DefaultAgent represents a basic AI agent with its configuration and state
type DefaultAgent struct {
	provider core.Provider
	tools    ToolMap
	memory   []*types.Message

	maxSteps int

	logger *slog.Logger
}

type ToolMap map[string]types.Tool

// NewAgent creates a new agent with the given provider
func NewAgent(config *NewAgentConfig) *DefaultAgent {
	if config.MaxSteps == 0 {
		// set a sane default max steps
		config.MaxSteps = 5
	}

	return &DefaultAgent{
		provider: config.Provider,
		tools:    make(map[string]types.Tool),
		memory:   make([]*types.Message, 0),
		maxSteps: config.MaxSteps,
		logger:   config.Logger,
	}
}

// Run implements the main agent loop
func (a *DefaultAgent) Run(ctx context.Context, input string, stopCondition types.AgentStopCondition) ([]*types.AgentStep, error) {
	var steps []*types.AgentStep

	currentStep := &types.AgentStep{
		ID: "1",
		Message: &types.Message{
			Role:       types.UserMessageRole,
			Content:    input,
			ToolCalls:  nil,
			ToolResult: nil,
			Metadata:   nil,
		},
		Error: nil,
	}

	steps = append(steps, currentStep)

	for {
		a.logger.Debug("sending message", "message", currentStep.Message)
		respMessage, err := a.SendMessage(ctx, currentStep.Message)

		a.logger.Debug("response message", "message", respMessage)
		respStep := &types.AgentStep{
			ID:      "2",
			Message: respMessage,
			Error:   err,
		}
		steps = append(steps, respStep)
		if err != nil {
			return steps, err
		}

		// Check stop condition
		if stopCondition(respStep) {
			a.logger.Debug("reached stop condition", "steps", len(steps))
			return steps, nil
		}

		// Check max steps
		if len(steps) >= a.maxSteps {
			return steps, fmt.Errorf("exceeded maximum steps: %d - %d", len(steps), a.maxSteps)
		}

		// 2 "send" scenarios:
		//    * "user" message
		//    * "tool" results message
		//
		// 1 "receive" scenario:
		//    * LLM responds with "content" and "tool_calls". Either or may be empty

		// Prepare next input based on tool results
		if len(respStep.Message.ToolCalls) > 0 {
			a.logger.Debug("calling tool", "tool", respStep.Message.ToolCalls[0].Name)
			toolMessage, err := a.CallTool(ctx, respStep.Message.ToolCalls[0])

			a.logger.Debug("tool response message", "message", toolMessage)
			currentStep = &types.AgentStep{
				ID:      "3",
				Message: toolMessage,
				Error:   err,
			}
		} else {
			currentStep = &types.AgentStep{
				ID:      "4",
				Message: respMessage,
				Error:   err,
			}
		}
	}
}

// RunStream supports a streaming channel from a provider
func (a *DefaultAgent) RunStream(ctx context.Context, input string, stopCondition types.AgentStopCondition) (<-chan types.AgentStep, <-chan string, <-chan error) {
	// TODO - need to implement step stream
	//stepsChan := make(chan *types.AgentStep)

	currentStep := &types.AgentStep{
		ID: "1",
		Message: &types.Message{
			Role:       types.UserMessageRole,
			Content:    input,
			ToolCalls:  nil,
			ToolResult: nil,
			Metadata:   nil,
		},
		Error: nil,
	}

	a.logger.Debug("adding to chan", "current step", currentStep)
	//stepsChan <- currentStep

	a.logger.Debug("sending streaming message", "message", currentStep.Message)
	msgChan, deltaChan, errChan := a.SendMessageStream(ctx, currentStep.Message)

	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				a.logger.Info("stream message channel closed")
				return nil, nil, nil
			}
			if msg != nil {
				a.logger.Info("received message",
					"role", msg.Role,
					"content", msg.Content,
					"tool_calls", msg.ToolCalls,
				)
			}

		case delta, ok := <-deltaChan:
			if !ok {
				a.logger.Info("stream delta chan closed")
				return nil, nil, nil
			}
			if delta != "" {
				print(delta)
			}

		case err, ok := <-errChan:
			if !ok {
				a.logger.Info("stream error chan closed")
				return nil, nil, nil
			}
			if err != nil {
				panic(err)
			}

		case <-ctx.Done():
			return nil, nil, nil

		case <-time.After(30 * time.Second):
			a.logger.Error("stream timeout")
			panic("stream timeout")
		}
	}
}

// SendMessage sends a message to the agent and gets a response
func (a *DefaultAgent) SendMessage(ctx context.Context, m *types.Message) (*types.Message, error) {
	a.memory = append(a.memory, m)

	var response *types.Message
	var err error

	toolSlice := make([]*types.Tool, 0, len(a.tools))
	for _, tool := range a.tools {
		toolSlice = append(toolSlice, &tool)
	}

	genOpts := &types.GenerateOptions{
		Messages: a.memory,
		Tools:    toolSlice,
	}

	a.logger.Debug("sending message with generate options", "genOpts", genOpts)
	response, err = a.provider.Generate(ctx, genOpts)
	if err != nil {
		return nil, err
	}

	a.memory = append(a.memory, response)
	return response, nil
}

// SendMessage sends a message to the agent and gets a response
func (a *DefaultAgent) SendMessageStream(ctx context.Context, m *types.Message) (<-chan *types.Message, <-chan string, <-chan error) {
	a.memory = append(a.memory, m)

	toolSlice := make([]*types.Tool, 0, len(a.tools))
	for _, tool := range a.tools {
		toolSlice = append(toolSlice, &tool)
	}

	genOpts := &types.GenerateOptions{
		Messages: a.memory,
		Tools:    toolSlice,
	}

	a.logger.Debug("sending message with generate options", "genOpts", genOpts)
	return a.provider.GenerateStream(ctx, genOpts)
}

// CallTool sends a message to the agent and gets a response
func (a *DefaultAgent) CallTool(ctx context.Context, tc *types.ToolCall) (*types.Message, error) {
	// Find the corresponding tool
	var toolToCall *types.Tool

	for _, t := range a.tools {
		if t.Name == tc.Name {
			toolToCall = &t
			break
		}
	}

	if toolToCall == nil {
		return nil, fmt.Errorf("tool %s not found", tc.Name)
	}

	// Call the tool
	result, err := toolToCall.WrappedToolFunction(ctx, []byte(tc.Arguments))
	if err != nil {
		return nil, fmt.Errorf("tool execution failed: %w", err)
	}

	// Add the tool response to messages
	return &types.Message{
		Role:    types.ToolMessageRole,
		Content: fmt.Sprintf("%v", result),
		ToolResult: &types.ToolResult{
			ToolCallID: tc.ID,
			Content:    result,
			Error:      "",
		},
	}, nil
}

// AddTool adds a tool to the agent's available tools
func (a *DefaultAgent) AddTool(tool types.Tool) error {
	if tool.Name == "" {
		return errors.New("tool must have a name")
	}

	if tool.WrappedToolFunction == nil {
		return errors.New("tool must have a function")
	}

	a.tools[tool.Name] = tool

	return nil
}

// Example stop condition
func DefaultStopCondition(step *types.AgentStep) bool {
	// Stop if there's an error
	if step.Error != nil {
		return true
	}

	// Stop if no tool calls were made and we got a response
	if len(step.Message.ToolCalls) == 0 && len(step.Message.Content) != 0 {
		return true
	}

	return false
}
