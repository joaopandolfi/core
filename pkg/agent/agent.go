package agent

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/agent-api/core"
	inmemory "github.com/agent-api/core/pkg/memory/inmem"
	"github.com/agent-api/core/types"
)

// Agent represents a basic AI agent with its configuration and state
type Agent struct {
	provider core.Provider
	tools    ToolMap
	memory   core.MemoryStorer

	maxSteps int

	logger *slog.Logger
}

type ToolMap map[string]types.Tool

// NewAgent creates a new agent with the given provider
func NewAgent(config *NewAgentConfig) *Agent {
	if config.MaxSteps == 0 {
		// set a sane default max steps
		config.MaxSteps = 25
	}

	if config.Memory == nil {
		config.Memory = inmemory.NewInMemoryMemStore()
	}

	return &Agent{
		provider: config.Provider,
		tools:    make(map[string]types.Tool),
		memory:   config.Memory,
		maxSteps: config.MaxSteps,
		logger:   config.Logger,
	}
}

// Run implements the main agent loop
func (a *Agent) Run(ctx context.Context, opts ...RunOptionFunc) *types.AgentRunAggregator {
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
	m := &types.Message{
		ID:         id,
		Role:       types.UserMessageRole,
		Content:    runOpts.Input,
		Images:     runOpts.Images,
		ToolCalls:  nil,
		ToolResult: nil,
		Metadata:   nil,
	}
	agg.Push(nil, m)
	a.memory.Push(m)

	for {
		a.logger.Debug("sending messages", "messages", agg.Messages)
		respMessage, respErr := a.SendMessages(ctx, agg)
		respMessage.ID = atomic.AddUint32(&id, 1)
		a.memory.Push(respMessage)

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
		//messages = []*types.Message{respMessage}

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
			a.memory.Push(toolResponses...)
		}
	}
}

type StreamRunnerResults struct {
	AggChan   <-chan types.AgentRunAggregator
	DeltaChan <-chan string
	ErrChan   <-chan error
}

// RunStream supports a streaming channel from a provider
func (a *Agent) RunStream(ctx context.Context, opts ...RunOptionFunc) *StreamRunnerResults {
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

	// buffered, non-blocking channels
	outAggChan := make(chan types.AgentRunAggregator, 10)
	outDeltaChan := make(chan string, 10)
	outErrChan := make(chan error, 10)

	result := &StreamRunnerResults{
		AggChan:   outAggChan,
		DeltaChan: outDeltaChan,
		ErrChan:   outErrChan,
	}

	// init aggregator
	agg := types.NewAgentRunAggregator()
	m := &types.Message{
		Role:       types.UserMessageRole,
		Content:    runOpts.Input,
		Images:     runOpts.Images,
		ToolCalls:  nil,
		ToolResult: nil,
		Metadata:   nil,
	}
	agg.Push(nil, m)
	a.memory.Push(m)

	a.logger.Debug("kicking run streamer")

	go func() {
		defer close(outAggChan)
		defer close(outDeltaChan)
		defer close(outErrChan)

		// Send initial aggregator state (non-blocking)
		select {
		case outAggChan <- *agg:
		default:
			// Skip if no one is listening
		}

		for {
			// Get streaming response for current messages
			msgChan, deltaChan, errChan := a.SendMessageStream(ctx, agg)

			var respMessage *types.Message
			var respErr error

			for {
				// escape inner loop if we're all done with this message stream
				allClosed := msgChan == nil && deltaChan == nil && errChan == nil
				if allClosed {
					break
				}

				select {
				case msg, ok := <-msgChan:
					if !ok {
						a.logger.Debug("send message message chan closed")
						msgChan = nil
						continue
					}
					if msg != nil {
						a.logger.Info("received message",
							"role", msg.Role,
							"content", msg.Content,
							"tool_calls", msg.ToolCalls,
						)
						respMessage = msg
						respMessage.ID = atomic.AddUint32(&id, 1)
					}

				case delta, ok := <-deltaChan:
					if !ok {
						a.logger.Debug("send message delta chan closed")
						deltaChan = nil
						continue
					}

					if delta != "" {
						select {
						case outDeltaChan <- delta:
						default:
							// Skip if no one is listening
						}
					}

					// pull errors from the downstream provider error channel.
				case err, ok := <-errChan:
					if !ok {
						a.logger.Debug("send message err chan closed")
						errChan = nil
						continue
					}
					if err != nil {
						respErr = err
						// Forward error to output channel (non-blocking)
						select {
						case outErrChan <- err:
						default:
							// Skip if no one is listening
						}
					}

				case <-ctx.Done():
					select {
					case outErrChan <- ctx.Err():
					default:
						// Skip if no one is listening
					}

					return
				}
			}

			// If we got a response message, add it to the aggregator
			if respMessage != nil {
				agg.Push(respErr, respMessage)
				a.memory.Push(respMessage)
				select {
				case outAggChan <- *agg:
				default:
					// Skip if no one is listening
				}
			}

			// If there was an error, return
			if respErr != nil {
				return
			}

			// Check stop condition
			if runOpts.StopCondition(agg) {
				a.logger.Debug("reached stop condition", "steps", len(agg.Messages))
				return
			}

			// Check max steps
			if len(agg.Messages) >= a.maxSteps {
				respErr = fmt.Errorf("exceeded maximum steps: %d - %d", len(agg.Messages), a.maxSteps)
				agg.Err = respErr
				select {
				case outErrChan <- respErr:
				default:
					// Skip if no one is listening
				}
				return
			}

			// Call tools if tool calls were present
			if respMessage != nil && len(respMessage.ToolCalls) > 0 {
				toolResponses := a.executeToolCallsParallel(ctx, respMessage.ToolCalls, id)
				agg.Push(nil, toolResponses...)
				a.memory.Push(toolResponses...)

				// Send updated aggregator after tool execution
				select {
				case outAggChan <- *agg:
				default:
					// Skip if no one is listening
				}
			}
		}
	}()

	return result
}

// SendMessage sends a message to the agent and gets a response
func (a *Agent) SendMessages(ctx context.Context, agg *types.AgentRunAggregator) (*types.Message, error) {
	toolSlice := make([]*types.Tool, 0, len(a.tools))
	for _, tool := range a.tools {
		toolSlice = append(toolSlice, &tool)
	}

	genOpts := &types.GenerateOptions{
		Messages: agg.Messages,
		Tools:    toolSlice,
	}

	a.logger.Debug("sending message with generate options", "genOpts", genOpts)
	response, err := a.provider.Generate(ctx, genOpts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SendMessage sends a message to the agent and gets a response
func (a *Agent) SendMessageStream(ctx context.Context, agg *types.AgentRunAggregator) (<-chan *types.Message, <-chan string, <-chan error) {
	toolSlice := make([]*types.Tool, 0, len(a.tools))
	for _, tool := range a.tools {
		toolSlice = append(toolSlice, &tool)
	}

	genOpts := &types.GenerateOptions{
		Messages: agg.Messages,
		Tools:    toolSlice,
	}

	a.logger.Debug("sending message with generate options", "genOpts", genOpts)
	return a.provider.GenerateStream(ctx, genOpts)
}

// CallTool sends a message to the agent and gets a response
func (a *Agent) CallTool(ctx context.Context, tc *types.ToolCall) (*types.Message, error) {
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
func (a *Agent) AddTool(tool types.Tool) error {
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
func (a *Agent) executeToolCallsParallel(ctx context.Context, toolCalls []*types.ToolCall, id uint32) []*types.Message {
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
