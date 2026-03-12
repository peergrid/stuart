package transcript

import "regexp"

// Filter defines criteria for selecting records.
type Filter struct {
	Role         string // "user", "assistant", "system", "all", or ""
	RecordType   string // Filter by record type field
	ToolName     string // Filter for records containing this tool name
	Contains     string // Regex to match in message text
	AuditField   string // Dot-notation field path
	AuditPattern string // Regex to match against audit field value
	Since        string // ISO timestamp lower bound
	Until        string // ISO timestamp upper bound
	ExternalOnly bool   // Skip meta, system-reminders, tool results
	NoSidechain  bool   // Skip sidechain records

	// Compiled regexes (lazily initialized)
	containsRe *regexp.Regexp
	auditRe    *regexp.Regexp
}

// Compile prepares the filter for use. Call before Matches.
func (f *Filter) Compile() error {
	if f.Contains != "" {
		re, err := regexp.Compile("(?i)" + f.Contains)
		if err != nil {
			return err
		}
		f.containsRe = re
	}
	if f.AuditPattern != "" {
		re, err := regexp.Compile("(?i)" + f.AuditPattern)
		if err != nil {
			return err
		}
		f.auditRe = re
	}
	return nil
}

// Matches returns true if the record passes all filter criteria.
func (f *Filter) Matches(r *Record) bool {
	// Sidechain filter
	if f.NoSidechain && r.IsSidechain {
		return false
	}

	// Record type filter
	if f.RecordType != "" && r.Type != f.RecordType {
		return false
	}

	// Role filter
	if f.Role != "" && f.Role != "all" {
		msg, err := r.ParseMessage()
		if err != nil || msg == nil || msg.Role != f.Role {
			return false
		}
	}

	// Timestamp filters
	if f.Since != "" && r.Timestamp < f.Since {
		return false
	}
	if f.Until != "" && r.Timestamp > f.Until {
		return false
	}

	// TODO: implement remaining filters
	// - ExternalOnly: skip isMeta, content starting with <, tool_result-only content
	// - ToolName: check tool_use blocks for matching name
	// - Contains: regex match against text content
	// - Audit: resolve dot-path field, match against regex

	return true
}

// IsEmpty returns true if no filter criteria are set.
func (f *Filter) IsEmpty() bool {
	return f.Role == "" && f.RecordType == "" && f.ToolName == "" &&
		f.Contains == "" && f.AuditField == "" && f.Since == "" &&
		f.Until == "" && !f.ExternalOnly && !f.NoSidechain
}
