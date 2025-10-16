package agent

import "github.com/joaopandolfi/core"

type ToolMap map[string]*core.Tool

// GetTools returns the current set of available tools
func (a *Agent) GetTools() []*core.Tool {
	t := []*core.Tool{}

	for _, v := range a.tools {
		t = append(t, v)
	}
	return t
}
