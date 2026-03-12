package analysis

import "stuart/doreen/tools/tq/internal/transcript"

// CompactionEvent describes a detected compaction.
type CompactionEvent struct {
	Type          string `json:"type"`      // "compact_boundary", "post_compaction_marker", "context_drop"
	Turn          int    `json:"turn"`
	Timestamp     string `json:"timestamp"`
	Trigger       string `json:"trigger,omitempty"`         // For compact_boundary
	PreTokens     int    `json:"pre_tokens,omitempty"`      // For compact_boundary
	PrevContext   int    `json:"prev_context,omitempty"`     // For context_drop
	NewContext    int    `json:"new_context,omitempty"`      // For context_drop
	DropPercent   float64 `json:"drop_percent,omitempty"`   // For context_drop
	MarkerText    string `json:"marker_text,omitempty"`     // For post_compaction_marker
}

// DetectCompactions finds compaction events in a transcript.
func DetectCompactions(records []*transcript.Record) ([]CompactionEvent, error) {
	// TODO: implement
	// 1. Method 1: compact_boundary system records
	// 2. Method 2: POST-COMPACTION RECOVERY markers
	// 3. Method 3: Total context token drops >50% from >50K baseline
	return nil, nil
}
