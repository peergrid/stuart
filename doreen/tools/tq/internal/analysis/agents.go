package analysis

// AgentTrace holds the lifecycle data for a subagent.
type AgentTrace struct {
	AgentID       string            `json:"agent_id"`
	AgentName     string            `json:"agent_name"`
	Model         string            `json:"model"`
	LaunchTurn    int               `json:"launch_turn"`
	ReturnTurn    int               `json:"return_turn"`
	LaunchTime    string            `json:"launch_time"`
	ReturnTime    string            `json:"return_time"`
	DurationSecs  int               `json:"duration_secs"`
	LaunchPrompt  string            `json:"launch_prompt"`
	PromptTokens  int               `json:"prompt_tokens"`
	ResultSummary string            `json:"result_summary"`
	ResultTokens  int               `json:"result_tokens"`
	ToolCounts    map[string]int    `json:"tool_counts"`
	TotalTokensIn int               `json:"total_tokens_in"`
	TotalTokensOut int              `json:"total_tokens_out"`
}

// TraceAgents extracts lifecycle data for all subagents in a session.
func TraceAgents(projDir string, sessionID string) ([]AgentTrace, error) {
	// TODO: implement
	// 1. Find subagent JSONL files
	// 2. For each: extract name, model, launch prompt, tools, tokens, result
	// 3. Correlate with root transcript for launch/return turns
	return nil, nil
}
