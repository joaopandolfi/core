package agent

// StopCondition is a function that determines if the agent should stop
// after its completed a step (i.e., a full "start" -> "doing work" -> "done" cycle)
type AgentStopCondition func(step *AgentRunAggregator) bool
