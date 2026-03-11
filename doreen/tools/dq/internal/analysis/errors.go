package analysis

import "stuart/doreen/tools/dq/internal/transcript"

// TranscriptError describes an error found in the transcript.
type TranscriptError struct {
	Type      string `json:"type"` // "tool_error", "hook_error", "permission_denied", "runtime_error"
	Turn      int    `json:"turn"`
	Timestamp string `json:"timestamp"`
	ToolUseID string `json:"tool_use_id,omitempty"`
	Detail    string `json:"detail"`
}

// ExtractErrors finds all errors in a transcript.
func ExtractErrors(records []*transcript.Record) ([]TranscriptError, error) {
	// TODO: implement
	// 1. tool_result with is_error=true
	// 2. Hook errors (hook_errors, exit_code, "hook failed" in text)
	// 3. Permission denials ("permission" + "denied"/"rejected")
	// 4. Runtime errors (MaxFileReadTokenExceededError, EPERM, etc.)
	return nil, nil
}
