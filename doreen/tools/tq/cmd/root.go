// Package cmd implements the tq CLI command tree.
//
// Each subcommand lives in its own file. Global flags (time window,
// filters, output mode) are defined here and inherited by all subcommands.
package cmd

import (
	"flag"
	"fmt"
	"os"
)

// Global flags shared across all subcommands.
var (
	// Scope — all time-based, project defaults from CWD
	flagProject   string
	flagSession   string // Optional power-user filter, never required
	flagSince     string // Batch: time filter. Default: "24h".
	flagUntil     string
	flagAround    string // Cursor: approximate starting point. "around 2 days ago"
	flagAll       bool
	flagSubagents bool

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
)

// registerGlobalFlags adds the shared flags to a FlagSet.
func registerGlobalFlags(fs *flag.FlagSet) {
	// Scope
	fs.StringVar(&flagProject, "project", "", "Project name (default: detected from CWD)")
	fs.StringVar(&flagSession, "session", "", "Filter to a specific session UUID prefix")
	fs.StringVar(&flagSince, "since", "24h", "Batch: time filter start (30m, 2h, 3d, 1w, ISO date)")
	fs.StringVar(&flagUntil, "until", "", "Batch: time filter end (default: now)")
	fs.StringVar(&flagAround, "around", "", "Cursor: approximate starting point (2h, 2d, 1w, ISO date)")
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
	fmt.Println("Scope flags:")
	fmt.Println("  --since DURATION  Batch: time filter (30m, 2h, 3d, 1w, ISO date; default: 24h)")
	fmt.Println("  --until TIME      Batch: end of filter window (default: now)")
	fmt.Println("  --around DURATION Cursor: approximate starting point (2h, 2d, 1w, ISO date)")
	fmt.Println("  --project NAME    Override project (default: from CWD)")
	fmt.Println("  --session UUID    Filter to one session (rarely needed)")
	fmt.Println("  --all             All projects")
	fmt.Println("  --subagents       Include subagent transcripts")
	fmt.Println()
	fmt.Println("Output flags:")
	fmt.Println("  --json            Structured JSON output")
	fmt.Println("  --jsonl           Streaming JSONL output")
	fmt.Println("  --limit N         Max records to output")
	fmt.Println("  --no-truncate     Show full content")
}
