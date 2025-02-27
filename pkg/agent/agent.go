package agent

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/agent-api/core"
	inmemory "github.com/agent-api/core/pkg/memory/inmem"
	"github.com/agent-api/core/types"
)

// DefaultAgent represents a basic AI agent with its configuration and state
type DefaultAgent struct {
	provider core.Provider
	tools    ToolMap
	memory   core.MemoryStorer

	maxSteps int

	logger *slog.Logger
}

type ToolMap map[string]types.Tool

// NewAgent creates a new agent with the given provider
func NewAgent(config *NewAgentConfig) *DefaultAgent {
	if config.MaxSteps == 0 {
		// set a sane default max steps
		config.MaxSteps = 25
	}

	if config.Memory == nil {
		config.Memory = inmemory.NewInMemoryMemStore()
	}

	return &DefaultAgent{
		provider: config.Provider,
		tools:    make(map[string]types.Tool),
		memory:   config.Memory,
		maxSteps: config.MaxSteps,
		logger:   config.Logger,
	}
}

// Run implements the main agent loop
func (a *DefaultAgent) Run(ctx context.Context, opts ...RunOptionFunc) *types.AgentRunAggregator {
	// Initialize with default options
	runOpts := &RunOptions{
		Input:         "Execute given tasks.",
		StopCondition: DefaultStopCondition,
		Images:        []*types.Image{},
	}

	// Apply all option functions
	for _, opt := range opts {
		opt(runOpts)
	}

	var id uint32 = 0

	agg := types.NewAgentRunAggregator()
	messages := []*types.Message{
		{
			ID:         id,
			Role:       types.UserMessageRole,
			Content:    runOpts.Input,
			Images:     runOpts.Images,
			ToolCalls:  nil,
			ToolResult: nil,
			Metadata:   nil,
		},
	}

	agg.Push(nil, messages...)

	for {
		a.logger.Debug("sending messages", "messages", messages)
		respMessage, respErr := a.SendMessages(ctx, agg, messages...)
		respMessage.ID = atomic.AddUint32(&id, 1)

		a.logger.Debug("response message", "message", respMessage)
		agg.Push(respErr, respMessage)
		if respErr != nil {
			return agg
		}

		// Check stop condition
		if runOpts.StopCondition(agg) {
			a.logger.Debug("reached stop condition", "steps", len(agg.Messages))
			return agg
		}

		// Check max steps
		if len(agg.Messages) >= a.maxSteps {
			a.logger.Error("exceeded max steps", "steps", len(agg.Messages))
			agg.Err = fmt.Errorf("exceeded maximum steps: %d - %d", len(agg.Messages), a.maxSteps)
			return agg
		}

		// reset messages for next go around
		messages = []*types.Message{respMessage}

		// 2 "send" scenarios:
		//    * "user" message
		//    * "tool" results message
		//
		// 1 "receive" scenario:
		//    * LLM responds with "content" and "tool_calls". Either or may be empty

		// Call tools if tool calls were present
		if len(respMessage.ToolCalls) > 0 {
			toolResponses := a.executeToolCallsParallel(ctx, respMessage.ToolCalls, id)
			agg.Push(nil, toolResponses...)

			// reset the messages to the tool responses
			messages = toolResponses
		}
	}
}

// RunStream supports a streaming channel from a provider
func (a *DefaultAgent) RunStream(ctx context.Context, input string, stopCondition types.AgentStopCondition) (<-chan types.AgentRunAggregator, <-chan string, <-chan error) {
	// TODO - need to implement step stream
	//stepsChan := make(chan *types.AgentStep)

	//var err error

	//agg := types.NewAgentMessageAggregator()
	m := &types.Message{
		Role:       types.UserMessageRole,
		Content:    input,
		ToolCalls:  nil,
		ToolResult: nil,
		Metadata:   nil,
	}

	//a.logger.Debug("adding to chan", "current step", currentStep)
	//stepsChan <- currentStep

	a.logger.Debug("sending streaming message", "message", m)
	msgChan, deltaChan, errChan := a.SendMessageStream(ctx, m)

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
func (a *DefaultAgent) SendMessages(ctx context.Context, agg *types.AgentRunAggregator, m ...*types.Message) (*types.Message, error) {
	a.memory.Push(m...)

	var response *types.Message
	var err error

	toolSlice := make([]*types.Tool, 0, len(a.tools))
	for _, tool := range a.tools {
		toolSlice = append(toolSlice, &tool)
	}

	genOpts := &types.GenerateOptions{
		Messages: agg.Messages,
		Tools:    toolSlice,
	}

	a.logger.Debug("sending message with generate options", "genOpts", genOpts)
	response, err = a.provider.Generate(ctx, genOpts)
	if err != nil {
		return nil, err
	}

	a.memory.Push(response)
	return response, nil
}

// SendMessage sends a message to the agent and gets a response
func (a *DefaultAgent) SendMessageStream(ctx context.Context, m *types.Message) (<-chan *types.Message, <-chan string, <-chan error) {
	a.memory.Push(m)

	toolSlice := make([]*types.Tool, 0, len(a.tools))
	for _, tool := range a.tools {
		toolSlice = append(toolSlice, &tool)
	}

	genOpts := &types.GenerateOptions{
		Messages: a.memory.Dump(),
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
		ToolResult: []*types.ToolResult{
			{
				ToolCallID: tc.ID,
				Content:    result,
				Error:      "",
			},
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
func DefaultStopCondition(agg *types.AgentRunAggregator) bool {
	// Stop if there's an error
	if agg.Err != nil {
		return true
	}

	// Stop if no tool calls were made and we got a response
	if len(agg.Messages) != 0 {
		if len(agg.Messages[len(agg.Messages)-1].ToolCalls) == 0 && len(agg.Messages[len(agg.Messages)-1].Content) != 0 {
			return true
		}
	}

	return false
}

// executeToolCallsParallel executes multiple tool calls in parallel using WaitGroup
func (a *DefaultAgent) executeToolCallsParallel(ctx context.Context, toolCalls []*types.ToolCall, id uint32) []*types.Message {
	var wg sync.WaitGroup
	responses := make([]*types.Message, len(toolCalls))

	for i, toolCall := range toolCalls {
		wg.Add(1)

		// Launch each tool call in its own goroutine
		go func(i int, tc *types.ToolCall) {
			defer wg.Done()

			a.logger.Debug("calling tool", "tool", tc.Name, "id", tc.ID)
			toolResp, internalErr := a.CallTool(ctx, tc)

			// handle the internal tool calling error
			// (this is different from errors related to LLM hallucinations like
			// improperly formatted json or missing required params)
			if internalErr != nil {
				a.logger.Error("tool execution failed",
					"tool", tc.Name,
					"error", internalErr)

				toolResp = &types.Message{
					ID:        atomic.AddUint32(&id, 1),
					Role:      types.ToolMessageRole,
					Content:   "",
					ToolCalls: nil,
					ToolResult: []*types.ToolResult{
						{
							ToolCallID: tc.ID,
							Error:      fmt.Sprintf("internal error executing tool %s: %v", tc.Name, internalErr),
						},
					},
					Metadata: nil,
				}
			}

			a.logger.Debug("tool response message", "message", toolResp)
			responses[i] = toolResp
		}(i, toolCall)
	}

	// Wait for all tool calls to complete
	wg.Wait()
	return responses
}
