package analysis

import "stuart/doreen/tools/dq/internal/transcript"

// TokenTimeline holds the per-turn token consumption data.
type TokenTimeline struct {
	Turns   []TokenTurn  `json:"turns"`
	Summary TokenSummary `json:"summary"`
}

// TokenTurn holds token data for a single assistant turn.
type TokenTurn struct {
	Turn            int      `json:"turn"`
	Timestamp       string   `json:"timestamp"`
	InputTokens     int      `json:"input_tokens"`
	OutputTokens    int      `json:"output_tokens"`
	CacheRead       int      `json:"cache_read"`
	CacheCreate     int      `json:"cache_create"`
	TotalContext     int      `json:"total_context"`
	CumulativeInput int      `json:"cumulative_input"`
	CumulativeOutput int     `json:"cumulative_output"`
	Tools           []string `json:"tools"`
}

// TokenSummary holds aggregate token metrics.
type TokenSummary struct {
	TotalInputTokens    int     `json:"total_input_tokens"`
	TotalOutputTokens   int     `json:"total_output_tokens"`
	TotalCacheRead      int     `json:"total_cache_read"`
	MaxSingleTurnContext int    `json:"max_single_turn_context"`
	TotalTurns          int     `json:"total_turns"`
	TokensPerTurnAvg    float64 `json:"tokens_per_turn_avg"`
	CompactionCount     int     `json:"compaction_count"`
}

// AnalyzeTokens computes the token timeline from transcript records.
func AnalyzeTokens(records []*transcript.Record) (*TokenTimeline, error) {
	// TODO: implement
	// 1. Iterate assistant records
	// 2. Extract usage fields per turn
	// 3. Compute cumulative totals
	// 4. Track tool names per turn
	// 5. Build summary with aggregates
	return &TokenTimeline{}, nil
}
