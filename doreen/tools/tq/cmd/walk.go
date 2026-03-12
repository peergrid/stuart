package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"stuart/doreen/tools/tq/internal/format"
	"stuart/doreen/tools/tq/internal/transcript"
)

func init() {
	register("walk", "Step through a transcript turn by turn", runWalk)
}

// WalkRecord is the output format for a single turn in walk output.
type WalkRecord struct {
	Turn      int             `json:"turn"`
	Session   string          `json:"session,omitempty"`
	Role      string          `json:"role"`
	Timestamp string          `json:"timestamp"`
	Text      string          `json:"text,omitempty"`
	Tools     []WalkToolCall  `json:"tools,omitempty"`
	Content   json.RawMessage `json:"content,omitempty"`
}

// WalkToolCall summarizes a tool call within a turn.
type WalkToolCall struct {
	Name  string          `json:"name"`
	Input json.RawMessage `json:"input,omitempty"`
}

func runWalk(args []string) error {
	fs := flag.NewFlagSet("walk", flag.ExitOnError)
	reverse := fs.Bool("reverse", false, "Walk backward from end")
	toolsOnly := fs.Bool("tools-only", false, "Show only tool call turns")
	registerCursorFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}

	filter, err := filterFromFlags()
	if err != nil {
		return err
	}

	// Cursor command: load all files
	files, err := resolveCursorFiles()
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "tq walk: no transcript files found")
		return nil
	}

	records, err := transcript.LoadAll(files)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		fmt.Fprintln(os.Stderr, "tq walk: no records found")
		return nil
	}

	// Parse from/until anchors
	fromAnchor, err := transcript.ParseAnchor(flagFrom)
	if err != nil {
		return fmt.Errorf("--from: %w", err)
	}
	untilAnchor, err := transcript.ParseAnchor(flagStop)
	if err != nil {
		return fmt.Errorf("--until: %w", err)
	}

	// Determine start position
	startIdx := findStartPosition(records, fromAnchor, flagFrom, *reverse)

	// Walk and collect matching records
	count := 0
	matched := false
	limit := flagLimit
	if limit == 0 {
		limit = len(records)
	}

	step := 1
	if *reverse {
		step = -1
	}

	for i := startIdx; i >= 0 && i < len(records) && count < limit; i += step {
		r := records[i]

		// Check until anchor stop condition
		if untilAnchor != nil && untilAnchor.Matches(r) {
			break
		}

		// Apply tools-only filter
		if *toolsOnly {
			msg, err := r.ParseMessage()
			if err != nil || msg == nil {
				continue
			}
			tools := transcript.GetToolUses(msg.Content)
			if len(tools) == 0 {
				continue
			}
		}

		// Apply general filters
		if !filter.Matches(r) {
			continue
		}

		matched = true

		// --exists: short-circuit on first match
		if flagExists {
			os.Exit(0)
		}

		// --count: just count
		if flagCount {
			count++
			continue
		}

		// Build output record
		wr := buildWalkRecord(r)

		// --field: project a single field
		if flagField != "" {
			val := projectField(wr, flagField)
			if val != "" {
				fmt.Println(val)
			}
			count++
			continue
		}

		// Output
		if flagJSONL {
			jw := format.NewJSONLWriter()
			jw.WriteLine(wr)
		} else if flagJSON {
			jw := format.NewJSONWriter()
			jw.WriteValue(wr)
		} else {
			printWalkTurn(wr)
		}
		count++
	}

	if flagExists && !matched {
		os.Exit(1)
	}

	if flagCount {
		fmt.Println(count)
	}

	return nil
}

// findStartPosition determines where to begin walking.
func findStartPosition(records []*transcript.Record, fromAnchor *transcript.Anchor, fromSpec string, reverse bool) int {
	// Default positions
	if reverse && fromAnchor == nil && fromSpec == "" && flagAround == "" {
		return len(records) - 1
	}

	cursor := transcript.NewCursor(records)

	// --around positions approximately first
	if flagAround != "" {
		aroundTime, err := transcript.ParseDuration(flagAround)
		if err == nil {
			cursor.SeekToTimestamp(aroundTime.Format("2006-01-02T15:04:05"))
		}
	}

	// --from as timestamp
	if fromSpec != "" && fromAnchor == nil {
		if _, err := transcript.ParseDuration(fromSpec); err == nil {
			t, _ := transcript.ParseDuration(fromSpec)
			cursor.SeekToTimestamp(t.Format("2006-01-02T15:04:05"))
			return cursor.Pos()
		}
	}

	// --from as anchor: search from current position (which may have been set by --around)
	if fromAnchor != nil {
		start := cursor.Pos()
		for i := start; i < len(records); i++ {
			if fromAnchor.Matches(records[i]) {
				return i
			}
		}
		// Not found from around position, search from beginning
		for i := 0; i < start; i++ {
			if fromAnchor.Matches(records[i]) {
				return i
			}
		}
	}

	return cursor.Pos()
}

func buildWalkRecord(r *transcript.Record) WalkRecord {
	wr := WalkRecord{
		Turn:      r.TurnNumber,
		Session:   shortSession(r.SessionID),
		Role:      r.Type,
		Timestamp: r.Timestamp,
	}

	msg, err := r.ParseMessage()
	if err != nil || msg == nil {
		return wr
	}

	wr.Text = transcript.GetTextContent(msg.Content)
	tools := transcript.GetToolUses(msg.Content)
	for _, t := range tools {
		wr.Tools = append(wr.Tools, WalkToolCall{
			Name:  t.Name,
			Input: t.Input,
		})
	}

	return wr
}

func shortSession(id string) string {
	if len(id) > 8 {
		return id[:8]
	}
	return id
}

func printWalkTurn(wr WalkRecord) {
	ts := format.FormatTimestamp(wr.Timestamp)
	fmt.Printf("--- Turn %d | session:%s | %s | %s ---\n", wr.Turn, wr.Session, wr.Role, ts)
	if wr.Text != "" {
		text := wr.Text
		if !flagNoTruncate && len(text) > 200 {
			text = text[:200] + "..."
		}
		fmt.Printf("[Text] %s\n", text)
	}
	for _, t := range wr.Tools {
		inputStr := string(t.Input)
		if !flagNoTruncate && len(inputStr) > 100 {
			inputStr = inputStr[:100] + "..."
		}
		fmt.Printf("[Tool: %s] %s\n", t.Name, inputStr)
	}
}

// projectField extracts a dot-notation field from a WalkRecord or its tools.
func projectField(wr WalkRecord, field string) string {
	switch field {
	case "role":
		return wr.Role
	case "text":
		return wr.Text
	case "timestamp":
		return wr.Timestamp
	case "turn":
		return fmt.Sprintf("%d", wr.Turn)
	default:
		// For tool-related fields like "input.command", look at first tool
		if len(wr.Tools) > 0 {
			return projectToolField(wr.Tools[0], field)
		}
		return ""
	}
}

func projectToolField(tc WalkToolCall, field string) string {
	if tc.Input == nil {
		return ""
	}
	var m map[string]any
	if err := json.Unmarshal(tc.Input, &m); err != nil {
		return ""
	}
	return resolveFieldPath(m, field)
}

func resolveFieldPath(m map[string]any, path string) string {
	parts := splitDotPath(path)
	var current any = m
	for _, part := range parts {
		switch v := current.(type) {
		case map[string]any:
			current = v[part]
		default:
			return ""
		}
	}
	switch v := current.(type) {
	case string:
		return v
	case nil:
		return ""
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

func splitDotPath(path string) []string {
	var parts []string
	current := ""
	for _, c := range path {
		if c == '.' {
			if current != "" {
				parts = append(parts, current)
			}
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
