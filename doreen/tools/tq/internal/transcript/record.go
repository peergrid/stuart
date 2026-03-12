// Package transcript provides types and functions for loading, parsing,
// filtering, and navigating Claude Code transcript JSONL files.
package transcript

import "encoding/json"

// Record represents a single JSONL line from a transcript file.
type Record struct {
	Type        string          `json:"type"`                  // "user", "assistant", "system"
	Subtype     string          `json:"subtype,omitempty"`     // e.g. "compact_boundary"
	Timestamp   string          `json:"timestamp"`
	SessionID   string          `json:"sessionId"`
	Message     json.RawMessage `json:"message,omitempty"`     // Parsed lazily
	IsSidechain bool            `json:"isSidechain,omitempty"`
	IsMeta      bool            `json:"isMeta,omitempty"`

	// Compaction metadata (only for system/compact_boundary records)
	CompactMetadata *CompactMeta `json:"compactMetadata,omitempty"`

	// Internal tracking fields (not from JSONL)
	SourceFile string `json:"-"` // Which file this record came from
	TurnNumber int    `json:"-"` // Sequential position in the session
	LineNumber int    `json:"-"` // Line number in source file
}

// CompactMeta holds metadata from compaction boundary records.
type CompactMeta struct {
	Trigger   string `json:"trigger"`   // "auto" or "manual"
	PreTokens int    `json:"preTokens"` // Context size before compaction
}

// Message represents a parsed message payload.
type Message struct {
	Role    string          `json:"role"`
	Model   string          `json:"model,omitempty"`
	Content json.RawMessage `json:"content"`
	Usage   *Usage          `json:"usage,omitempty"`
}

// Usage holds token consumption data from an assistant turn.
type Usage struct {
	InputTokens             int `json:"input_tokens"`
	OutputTokens            int `json:"output_tokens"`
	CacheReadInputTokens    int `json:"cache_read_input_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
}

// TotalContext returns the total context window consumption for this turn.
func (u *Usage) TotalContext() int {
	return u.InputTokens + u.CacheReadInputTokens + u.CacheCreationInputTokens
}

// ContentBlock represents a typed block within message content.
type ContentBlock struct {
	Type      string          `json:"type"`                 // "text", "tool_use", "tool_result"
	Text      string          `json:"text,omitempty"`       // For text blocks
	ID        string          `json:"id,omitempty"`         // For tool_use blocks
	Name      string          `json:"name,omitempty"`       // For tool_use blocks
	Input     json.RawMessage `json:"input,omitempty"`      // For tool_use blocks
	ToolUseID string          `json:"tool_use_id,omitempty"` // For tool_result blocks
	IsError   bool            `json:"is_error,omitempty"`   // For tool_result blocks
	Content   json.RawMessage `json:"content,omitempty"`    // For tool_result blocks (string or block array)
}

// ParseMessage lazily parses the raw message JSON.
func (r *Record) ParseMessage() (*Message, error) {
	if r.Message == nil {
		return nil, nil
	}
	var m Message
	if err := json.Unmarshal(r.Message, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// ParseContentBlocks parses message content as an array of blocks.
// Returns nil if content is a plain string (call GetTextContent instead).
func ParseContentBlocks(raw json.RawMessage) ([]ContentBlock, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	// Check if it's a string (starts with ")
	if raw[0] == '"' {
		return nil, nil
	}
	var blocks []ContentBlock
	if err := json.Unmarshal(raw, &blocks); err != nil {
		return nil, err
	}
	return blocks, nil
}

// GetTextContent extracts plain text from message content, whether it's
// a string or an array of content blocks.
func GetTextContent(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	// Try as string first
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	// Try as block array
	blocks, err := ParseContentBlocks(raw)
	if err != nil {
		return ""
	}
	var result string
	for _, b := range blocks {
		if b.Type == "text" {
			if result != "" {
				result += "\n"
			}
			result += b.Text
		}
	}
	return result
}

// GetToolUses extracts tool_use blocks from an assistant record's content.
func GetToolUses(raw json.RawMessage) []ContentBlock {
	blocks, err := ParseContentBlocks(raw)
	if err != nil {
		return nil
	}
	var tools []ContentBlock
	for _, b := range blocks {
		if b.Type == "tool_use" {
			tools = append(tools, b)
		}
	}
	return tools
}
