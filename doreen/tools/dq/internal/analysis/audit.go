// Package analysis implements the analytical modes for dq:
// tool use auditing, token analysis, compaction detection,
// error extraction, agent tracing, and critique data extraction.
package analysis

import "stuart/doreen/tools/dq/internal/transcript"

// AuditResult holds the output of a tool use audit.
type AuditResult struct {
	TotalToolCalls int              `json:"total_tool_calls"`
	Violations     []AuditViolation `json:"violations"`
	ViolationRate  float64          `json:"violation_rate"`
	Score          float64          `json:"score"`

	// Per-check results
	DedicatedToolViolations []AuditViolation `json:"dedicated_tool_violations"`
	ReadBeforeEditViolations []AuditViolation `json:"read_before_edit_violations"`
	RedundantReadViolations  []AuditViolation `json:"redundant_read_violations"`
	ParallelismRatio         float64          `json:"parallelism_ratio"`
	DelegationRatio          float64          `json:"delegation_ratio"`
}

// AuditViolation describes a single tool use violation.
type AuditViolation struct {
	Turn       int    `json:"turn"`
	Category   string `json:"category"` // "dedicated_tool", "read_before_edit", "redundant_read"
	Tool       string `json:"tool"`
	Detail     string `json:"detail"`
	Suggestion string `json:"suggestion"`
}

// RunAudit performs the full tool use audit on a set of records.
func RunAudit(records []*transcript.Record) (*AuditResult, error) {
	// TODO: implement
	// 1. Build tool call sequence with turn numbers
	// 2. Check dedicated tool preference (Bash calls with cat/grep/find/sed/awk/echo)
	// 3. Check read-before-edit (Edit/Write preceded by Read of same file)
	// 4. Check redundant reads (same file read twice without edit between)
	// 5. Calculate parallelism ratio
	// 6. Calculate delegation ratio
	// 7. Aggregate violations and compute score
	return &AuditResult{}, nil
}
