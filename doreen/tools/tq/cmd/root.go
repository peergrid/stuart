// Package cmd implements the tq CLI command tree.
//
// Each subcommand lives in its own file. Common flags (filters, output mode)
// are defined here. Time/scope flags are split: batch commands get --since/--until,
// cursor commands get --from/--until/--around with anchor semantics.
package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"stuart/doreen/tools/tq/internal/transcript"
)

// Global flags shared across all subcommands.
var (
	// Scope — project and session
	flagProject   string
	flagSession   string // Optional power-user filter, never required
	flagAll       bool
	flagSubagents bool

	// Batch time flags (registered by batch commands only)
	flagSince string
	flagUntil string

	// Cursor flags (registered by cursor commands only)
	flagAround string
	flagFrom   string     // Anchor spec or timestamp
	flagStop   string     // Anchor spec or timestamp (walk --until)

	// Output
	flagJSON       bool
	flagJSONL      bool
	flagLimit      int
	flagNoTruncate bool

	// Filters
	flagRole         string
	flagType         string
	flagTool         string
	flagContains     string
	flagAudit        string
	flagExternalOnly bool
	flagNoSidechain  bool

	// Output modifiers
	flagCount  bool
	flagExists bool
	flagField  string
)

// registerCommonFlags adds flags shared by ALL subcommands.
func registerCommonFlags(fs *flag.FlagSet) {
	// Scope
	fs.StringVar(&flagProject, "project", "", "Project name (default: detected from CWD)")
	fs.StringVar(&flagSession, "session", "", "Filter to a specific session UUID prefix")
	fs.BoolVar(&flagAll, "all", false, "All projects")
	fs.BoolVar(&flagSubagents, "subagents", false, "Include subagent transcripts")

	// Output
	fs.BoolVar(&flagJSON, "json", false, "Output structured JSON")
	fs.BoolVar(&flagJSONL, "jsonl", false, "Output streaming JSONL")
	fs.IntVar(&flagLimit, "limit", 0, "Max records to output")
	fs.BoolVar(&flagNoTruncate, "no-truncate", false, "Show full content")

	// Filters
	fs.StringVar(&flagRole, "role", "", "Filter by role: user|assistant|system|all")
	fs.StringVar(&flagType, "type", "", "Filter by record type")
	fs.StringVar(&flagTool, "tool", "", "Filter for tool_use by name")
	fs.StringVar(&flagContains, "contains", "", "Regex match in message text")
	fs.StringVar(&flagAudit, "audit", "", "Match JSON field by regex (field=regex)")
	fs.BoolVar(&flagExternalOnly, "external-only", false, "Only genuine operator messages")
	fs.BoolVar(&flagNoSidechain, "no-sidechain", false, "Skip sidechain records")

	// Output modifiers
	fs.BoolVar(&flagCount, "count", false, "Output only the count of matching records")
	fs.BoolVar(&flagExists, "exists", false, "Exit 0 if any match, exit 1 if none (no output)")
	fs.StringVar(&flagField, "field", "", "Project a single field from each record (dot-notation, e.g. input.command)")
}

// registerBatchFlags adds --since/--until for batch commands.
func registerBatchFlags(fs *flag.FlagSet) {
	registerCommonFlags(fs)
	fs.StringVar(&flagSince, "since", "24h", "Time filter start (30m, 2h, 3d, 1w, ISO date)")
	fs.StringVar(&flagUntil, "until", "", "Time filter end (default: now)")
}

// registerCursorFlags adds --from/--until/--around for cursor commands.
func registerCursorFlags(fs *flag.FlagSet) {
	registerCommonFlags(fs)
	fs.StringVar(&flagFrom, "from", "", "Start anchor: timestamp or spec (tool:Bash:jq, error, user, pattern:REGEX)")
	fs.StringVar(&flagStop, "until", "", "Stop anchor: same syntax as --from")
	fs.StringVar(&flagAround, "around", "", "Approximate starting point (2h, 2d, 1w, ISO date)")
}

// subcommand maps name -> handler.
type subcommand struct {
	name  string
	brief string
	run   func(args []string) error
}

var subcommands []subcommand

func register(name, brief string, run func(args []string) error) {
	subcommands = append(subcommands, subcommand{name, brief, run})
}

// Execute parses args and dispatches to the appropriate subcommand.
func Execute() error {
	if len(os.Args) < 2 {
		printUsage()
		return nil
	}

	name := os.Args[1]

	// Help
	if name == "help" || name == "--help" || name == "-h" {
		printUsage()
		return nil
	}

	for _, sc := range subcommands {
		if sc.name == name {
			return sc.run(os.Args[2:])
		}
	}

	return fmt.Errorf("unknown command %q — run 'tq help' for usage", name)
}

// filterFromFlags builds a Filter from global flag values and compiles it.
func filterFromFlags() (*transcript.Filter, error) {
	var auditField, auditPattern string
	if flagAudit != "" {
		parts := strings.SplitN(flagAudit, "=", 2)
		if len(parts) == 2 {
			auditField = parts[0]
			auditPattern = parts[1]
		}
	}
	f := &transcript.Filter{
		Role:         flagRole,
		RecordType:   flagType,
		ToolName:     flagTool,
		Contains:     flagContains,
		AuditField:   auditField,
		AuditPattern: auditPattern,
		ExternalOnly: flagExternalOnly,
		NoSidechain:  flagNoSidechain,
	}
	if err := f.Compile(); err != nil {
		return nil, err
	}
	return f, nil
}

// resolveBatchFiles finds transcript files for batch commands using --since/--until.
func resolveBatchFiles() ([]string, error) {
	projDir, err := resolveProjectDir()
	if err != nil {
		return nil, err
	}

	sinceStr := flagSince
	if sinceStr == "" {
		sinceStr = "24h"
	}
	since, err := transcript.ParseDuration(sinceStr)
	if err != nil {
		return nil, err
	}
	until := time.Now()
	if flagUntil != "" {
		until, err = transcript.ParseDuration(flagUntil)
		if err != nil {
			return nil, err
		}
	}
	return transcript.FindSessionsInWindow(projDir, since, until)
}

// resolveCursorFiles finds all transcript files for cursor commands.
func resolveCursorFiles() ([]string, error) {
	projDir, err := resolveProjectDir()
	if err != nil {
		return nil, err
	}
	since := time.Now().AddDate(-1, 0, 0) // 1 year back
	return transcript.FindSessionsInWindow(projDir, since, time.Now())
}

func resolveProjectDir() (string, error) {
	if flagProject != "" {
		return transcript.ProjectDirFromName(flagProject)
	}
	return transcript.ProjectDirFromCWD()
}

func printUsage() {
	fmt.Println("tq — Transcript Query: analysis tool for Claude Code sessions")
	fmt.Println()
	fmt.Println("Usage: tq <command> [flags]")
	fmt.Println()
	fmt.Println("Project is detected from CWD. Time window defaults to last 24 hours.")
	fmt.Println("All queries span every session in the time window automatically.")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  tq errors                     Errors in the last 24h")
	fmt.Println("  tq errors --since 2d          Errors in the last 2 days")
	fmt.Println("  tq tools --since 1h           Tool calls in the last hour")
	fmt.Println("  tq stats                      Stats for last 24h, all sessions")
	fmt.Println("  tq audit --since 3d           Tool use audit, last 3 days")
	fmt.Println()
	fmt.Println("  tq walk --reverse --until user --role assistant --count")
	fmt.Println("  tq walk --from \"tool:Bash:jq\" --around 4d --until \"tool:Read\" --count")
	fmt.Println("  tq tools --tool Bash --limit 10 --field input.command")
	fmt.Println()
	fmt.Println("Discovery:")
	fmt.Println("  sessions    List sessions with metadata")
	fmt.Println("  agents      List subagents across sessions")
	fmt.Println("  find        Search sessions/agents by content")
	fmt.Println()
	fmt.Println("Navigation:")
	fmt.Println("  show        Show records around a point of interest")
	fmt.Println("  walk        Step through turns across sessions")
	fmt.Println()
	fmt.Println("Query:")
	fmt.Println("  messages    Extract messages with filtering")
	fmt.Println("  tools       Extract tool calls with filtering")
	fmt.Println("  raw         Dump raw JSONL records")
	fmt.Println()
	fmt.Println("Analysis:")
	fmt.Println("  stats       Summary statistics")
	fmt.Println("  tokens      Token consumption analysis")
	fmt.Println("  compactions Detect compaction events")
	fmt.Println("  errors      Extract errors and failures")
	fmt.Println("  audit       Tool use audit (for grading)")
	fmt.Println("  critique-data  Extract data for transcript critique grader")
	fmt.Println("  agent-trace    Trace subagent lifecycles")
	fmt.Println()
	fmt.Println("Batch scope flags (stats, errors, tools, audit, etc.):")
	fmt.Println("  --since DURATION  Time filter start (30m, 2h, 3d, 1w, ISO date; default: 24h)")
	fmt.Println("  --until TIME      Time filter end (default: now)")
	fmt.Println()
	fmt.Println("Cursor scope flags (walk, show):")
	fmt.Println("  --from ANCHOR     Start: timestamp or anchor spec (tool:Bash:jq, error, user)")
	fmt.Println("  --until ANCHOR    Stop: same syntax as --from")
	fmt.Println("  --around DURATION Approximate starting point (2h, 2d, 1w, ISO date)")
	fmt.Println()
	fmt.Println("Common flags:")
	fmt.Println("  --project NAME    Override project (default: from CWD)")
	fmt.Println("  --session UUID    Filter to one session (rarely needed)")
	fmt.Println("  --all             All projects")
	fmt.Println("  --subagents       Include subagent transcripts")
	fmt.Println()
	fmt.Println("Filter flags:")
	fmt.Println("  --role ROLE       Filter by role: user|assistant|system|all")
	fmt.Println("  --tool NAME       Filter for tool_use by name")
	fmt.Println("  --contains REGEX  Regex match in message text")
	fmt.Println("  --external-only   Only genuine operator messages")
	fmt.Println()
	fmt.Println("Output flags:")
	fmt.Println("  --json            Structured JSON output")
	fmt.Println("  --jsonl           Streaming JSONL output")
	fmt.Println("  --limit N         Max records to output")
	fmt.Println("  --no-truncate     Show full content")
	fmt.Println("  --count           Output only the count of matching records")
	fmt.Println("  --exists          Exit 0 if any match, exit 1 if none (no output)")
	fmt.Println("  --field PATH      Project a single field (dot-notation, e.g. input.command)")
}
