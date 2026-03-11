package analysis

import "stuart/doreen/tools/dq/internal/transcript"

// CritiqueData holds the structured data for the transcript-critique grader.
type CritiqueData struct {
	OriginalPrompt    string             `json:"original_prompt"`
	DecisionPoints    []DecisionPoint    `json:"decision_points"`
	ErrorRecoveries   []ErrorRecovery    `json:"error_recoveries"`
	CompletionClaims  []CompletionClaim  `json:"completion_claims"`
	CommunicationData CommunicationData  `json:"communication"`
	ContextManagement ContextManagement  `json:"context_management"`
}

// DecisionPoint marks a turn where Claude chose an approach.
type DecisionPoint struct {
	Turn      int    `json:"turn"`
	Timestamp string `json:"timestamp"`
	Text      string `json:"text"`      // What Claude said about the decision
	ToolsCalled []string `json:"tools_called"` // Tools used in the subsequent sequence
}

// ErrorRecovery captures an error and how Claude responded.
type ErrorRecovery struct {
	ErrorTurn     int    `json:"error_turn"`
	ErrorType     string `json:"error_type"`
	ErrorDetail   string `json:"error_detail"`
	RecoveryTurns []int  `json:"recovery_turns"` // Subsequent turns
	RecoveryText  string `json:"recovery_text"`  // What Claude did/said after
}

// CompletionClaim marks where Claude declared work done.
type CompletionClaim struct {
	Turn       int    `json:"turn"`
	Timestamp  string `json:"timestamp"`
	ClaimText  string `json:"claim_text"`
	TurnsAfter int    `json:"turns_after"` // How many turns followed this claim
	WasActuallyDone bool `json:"was_actually_done"`
}

// CommunicationData holds signals about communication quality.
type CommunicationData struct {
	AvgResponseLength float64 `json:"avg_response_length"`
	MaxResponseLength int     `json:"max_response_length"`
	TextToToolRatio   float64 `json:"text_to_tool_ratio"` // Text turns / tool call turns
}

// ContextManagement holds signals about context management quality.
type ContextManagement struct {
	CompactionCount    int      `json:"compaction_count"`
	RepeatedReads      []string `json:"repeated_reads"`      // Files read more than once
	PostCompactionReads []string `json:"post_compaction_reads"` // Files re-read after compaction
}

// ExtractCritiqueData builds the critique data from a transcript.
func ExtractCritiqueData(records []*transcript.Record) (*CritiqueData, error) {
	// TODO: implement
	// 1. Find original prompt (first external user message)
	// 2. Identify decision points (text before tool call sequences)
	// 3. Find error-recovery pairs
	// 4. Find completion claims (done/complete/finished) and check what followed
	// 5. Compute communication metrics
	// 6. Analyze context management (compactions, repeated reads)
	return &CritiqueData{}, nil
}
