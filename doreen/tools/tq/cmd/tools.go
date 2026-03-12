package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"stuart/doreen/tools/tq/internal/format"
	"stuart/doreen/tools/tq/internal/transcript"
)

func init() {
	register("tools", "Extract tool calls with filtering", runTools)
}

// ToolRecord is the output format for a single tool call.
type ToolRecord struct {
	Turn      int             `json:"turn"`
	Session   string          `json:"session,omitempty"`
	Timestamp string          `json:"timestamp"`
	Name      string          `json:"name"`
	Input     json.RawMessage `json:"input,omitempty"`
	Result    string          `json:"result,omitempty"`
	IsError   bool            `json:"is_error,omitempty"`
}

func runTools(args []string) error {
	fs := flag.NewFlagSet("tools", flag.ExitOnError)
	inputContains := fs.String("input-contains", "", "Regex match in tool input")
	withResults := fs.Bool("with-results", false, "Include tool results")
	registerBatchFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}

	filter, err := filterFromFlags()
	if err != nil {
		return err
	}

	files, err := resolveBatchFiles()
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "tq tools: no transcript files found")
		return nil
	}

	// Compile input-contains regex
	var inputContainsRe *regexp.Regexp
	if *inputContains != "" {
		inputContainsRe, err = regexp.Compile("(?i)" + *inputContains)
		if err != nil {
			return fmt.Errorf("--input-contains: %w", err)
		}
	}

	// We need all records for --with-results (to match tool_use with tool_result)
	records, err := transcript.LoadAll(files)
	if err != nil {
		return err
	}

	count := 0
	matched := false
	limit := flagLimit
	if limit == 0 {
		limit = len(records)
	}

	for i := 0; i < len(records) && count < limit; i++ {
		r := records[i]

		// Apply time/role filters from the filter (but not tool filter — we handle that)
		if filter.NoSidechain && r.IsSidechain {
			continue
		}

		msg, err := r.ParseMessage()
		if err != nil || msg == nil {
			continue
		}

		tools := transcript.GetToolUses(msg.Content)
		if len(tools) == 0 {
			continue
		}

		for _, t := range tools {
			// Apply tool name filter
			if flagTool != "" && t.Name != flagTool {
				continue
			}

			// Apply input-contains filter
			if inputContainsRe != nil && !inputContainsRe.Match(t.Input) {
				continue
			}

			if count >= limit {
				break
			}

			matched = true

			if flagExists {
				os.Exit(0)
			}

			if flagCount {
				count++
				continue
			}

			tr := ToolRecord{
				Turn:      r.TurnNumber,
				Session:   shortSession(r.SessionID),
				Timestamp: r.Timestamp,
				Name:      t.Name,
				Input:     t.Input,
			}

			// Find matching tool result if requested
			if *withResults {
				tr.Result, tr.IsError = findToolResult(records, i, t.ID)
			}

			// --field: project a single field
			if flagField != "" {
				val := projectToolRecordField(tr, flagField)
				if val != "" {
					fmt.Println(val)
				}
				count++
				continue
			}

			if flagJSONL {
				jw := format.NewJSONLWriter()
				jw.WriteLine(tr)
			} else if flagJSON {
				jw := format.NewJSONWriter()
				jw.WriteValue(tr)
			} else {
				printToolRecord(tr)
			}
			count++
		}
	}

	if flagExists && !matched {
		os.Exit(1)
	}

	if flagCount {
		fmt.Println(count)
	}

	return nil
}

// findToolResult searches for the tool_result matching a tool_use ID.
func findToolResult(records []*transcript.Record, startIdx int, toolUseID string) (string, bool) {
	// Tool results come in the next user record after the assistant record
	for i := startIdx + 1; i < len(records) && i <= startIdx+3; i++ {
		r := records[i]
		if r.Type != "user" {
			continue
		}
		msg, err := r.ParseMessage()
		if err != nil || msg == nil {
			continue
		}
		blocks, err := transcript.ParseContentBlocks(msg.Content)
		if err != nil {
			continue
		}
		for _, b := range blocks {
			if b.Type == "tool_result" && b.ToolUseID == toolUseID {
				text := transcript.GetTextContent(b.Content)
				return text, b.IsError
			}
		}
	}
	return "", false
}

// projectToolRecordField extracts a field from a ToolRecord.
func projectToolRecordField(tr ToolRecord, field string) string {
	switch field {
	case "name":
		return tr.Name
	case "timestamp":
		return tr.Timestamp
	case "turn":
		return fmt.Sprintf("%d", tr.Turn)
	case "result":
		return tr.Result
	default:
		// Resolve from input. Accept both "command" and "input.command".
		if tr.Input == nil {
			return ""
		}
		var m map[string]any
		if err := json.Unmarshal(tr.Input, &m); err != nil {
			return ""
		}
		// Strip "input." prefix if present — Input is already the input object
		path := field
		if after, ok := strings.CutPrefix(field, "input."); ok {
			path = after
		}
		return resolveFieldPath(m, path)
	}
}

func printToolRecord(tr ToolRecord) {
	ts := format.FormatTimestamp(tr.Timestamp)
	inputStr := string(tr.Input)
	if !flagNoTruncate && len(inputStr) > 120 {
		inputStr = inputStr[:120] + "..."
	}
	fmt.Printf("[%s] %s %s\n", ts, tr.Name, inputStr)
	if tr.Result != "" {
		result := tr.Result
		if !flagNoTruncate && len(result) > 200 {
			result = result[:200] + "..."
		}
		prefix := "  → "
		if tr.IsError {
			prefix = "  ✘ "
		}
		fmt.Printf("%s%s\n", prefix, result)
	}
}
