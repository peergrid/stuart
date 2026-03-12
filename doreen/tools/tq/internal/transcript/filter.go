package transcript

import (
	"encoding/json"
	"regexp"
)

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

	// External only: skip meta records, sidechain, and tool-result-only user records
	if f.ExternalOnly {
		if r.IsMeta || r.IsSidechain {
			return false
		}
		if r.Type == "system" {
			return false
		}
		// User records that only contain tool_result blocks are not external
		if r.Type == "user" && !IsExternalUserMessage(r) {
			return false
		}
	}

	// Tool name filter: record must contain a tool_use block with matching name
	if f.ToolName != "" {
		msg, err := r.ParseMessage()
		if err != nil || msg == nil {
			return false
		}
		tools := GetToolUses(msg.Content)
		found := false
		for _, t := range tools {
			if t.Name == f.ToolName {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Contains: regex match in text content
	if f.containsRe != nil {
		msg, err := r.ParseMessage()
		if err != nil || msg == nil {
			return false
		}
		text := GetTextContent(msg.Content)
		if !f.containsRe.MatchString(text) {
			return false
		}
	}

	return true
}

// IsExternalUserMessage checks whether a user record is a genuine human message
// (has at least one text content block) vs a tool-result-only record.
func IsExternalUserMessage(r *Record) bool {
	if r.Type != "user" {
		return false
	}
	msg, err := r.ParseMessage()
	if err != nil || msg == nil {
		return false
	}
	blocks, err := ParseContentBlocks(msg.Content)
	if err != nil {
		// Might be a plain string — that counts as external
		var s string
		if json.Unmarshal(msg.Content, &s) == nil && s != "" {
			return true
		}
		return false
	}
	for _, b := range blocks {
		if b.Type == "text" {
			return true
		}
	}
	return false
}

// IsEmpty returns true if no filter criteria are set.
func (f *Filter) IsEmpty() bool {
	return f.Role == "" && f.RecordType == "" && f.ToolName == "" &&
		f.Contains == "" && f.AuditField == "" && f.Since == "" &&
		f.Until == "" && !f.ExternalOnly && !f.NoSidechain
}
