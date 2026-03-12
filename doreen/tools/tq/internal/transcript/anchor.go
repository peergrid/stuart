package transcript

import (
	"regexp"
	"strings"
)

// Anchor defines a point of interest in a transcript.
// Parsed from string specs like "tool:Bash:jq", "error", "user", "compaction".
type Anchor struct {
	Kind         AnchorKind
	ToolName     string // For tool anchors
	InputPattern *regexp.Regexp // For tool anchors with input pattern
	TextPattern  *regexp.Regexp // For pattern anchors
}

// AnchorKind identifies what category of record an anchor matches.
type AnchorKind int

const (
	AnchorKindNone        AnchorKind = iota
	AnchorKindTool                   // tool:NAME or tool:NAME:PATTERN
	AnchorKindError                  // error
	AnchorKindCompact                // compaction
	AnchorKindRole                   // user, assistant, system
	AnchorKindPattern                // pattern:REGEX
)

// ParseAnchor parses an anchor spec string.
//
// Supported formats:
//   - "tool:NAME"         — first use of tool NAME
//   - "tool:NAME:PATTERN" — first use of tool NAME where input matches PATTERN
//   - "error"             — first error record
//   - "compaction"        — first compaction boundary
//   - "user"              — first external human message
//   - "assistant"         — first assistant record
//   - "system"            — first system record
//   - "pattern:REGEX"     — first record whose text matches REGEX
func ParseAnchor(spec string) (*Anchor, error) {
	if spec == "" {
		return nil, nil
	}

	a := &Anchor{}

	switch {
	case strings.HasPrefix(spec, "tool:"):
		a.Kind = AnchorKindTool
		rest := spec[5:]
		parts := strings.SplitN(rest, ":", 2)
		a.ToolName = parts[0]
		if len(parts) == 2 && parts[1] != "" {
			re, err := regexp.Compile("(?i)" + parts[1])
			if err != nil {
				return nil, err
			}
			a.InputPattern = re
		}

	case spec == "error":
		a.Kind = AnchorKindError

	case spec == "compaction":
		a.Kind = AnchorKindCompact

	case spec == "user" || spec == "assistant" || spec == "system":
		a.Kind = AnchorKindRole
		a.ToolName = spec // reuse field for role name

	case strings.HasPrefix(spec, "pattern:"):
		a.Kind = AnchorKindPattern
		re, err := regexp.Compile("(?i)" + spec[8:])
		if err != nil {
			return nil, err
		}
		a.TextPattern = re

	default:
		// Try as a role name with any case
		lower := strings.ToLower(spec)
		if lower == "user" || lower == "assistant" || lower == "system" {
			a.Kind = AnchorKindRole
			a.ToolName = lower
			return a, nil
		}
		// Default: treat as a text pattern
		a.Kind = AnchorKindPattern
		re, err := regexp.Compile("(?i)" + spec)
		if err != nil {
			return nil, err
		}
		a.TextPattern = re
	}

	return a, nil
}

// Matches checks if a record matches this anchor.
func (a *Anchor) Matches(r *Record) bool {
	if a == nil {
		return false
	}

	switch a.Kind {
	case AnchorKindTool:
		msg, err := r.ParseMessage()
		if err != nil || msg == nil {
			return false
		}
		tools := GetToolUses(msg.Content)
		for _, t := range tools {
			if t.Name != a.ToolName {
				continue
			}
			if a.InputPattern == nil {
				return true
			}
			// Match pattern against serialized input
			if a.InputPattern.Match(t.Input) {
				return true
			}
		}
		return false

	case AnchorKindError:
		// Tool errors, permission denials, hook failures
		if r.Type == "user" {
			msg, err := r.ParseMessage()
			if err != nil || msg == nil {
				return false
			}
			blocks, err := ParseContentBlocks(msg.Content)
			if err != nil {
				return false
			}
			for _, b := range blocks {
				if b.Type == "tool_result" && b.IsError {
					return true
				}
			}
		}
		return false

	case AnchorKindCompact:
		return r.Type == "system" && r.Subtype == "compact_boundary"

	case AnchorKindRole:
		if a.ToolName == "user" {
			return IsExternalUserMessage(r)
		}
		return r.Type == a.ToolName

	case AnchorKindPattern:
		if a.TextPattern == nil {
			return false
		}
		msg, err := r.ParseMessage()
		if err != nil || msg == nil {
			return false
		}
		text := GetTextContent(msg.Content)
		return a.TextPattern.MatchString(text)

	default:
		return false
	}
}
