// Package cmd implements the dq CLI command tree.
//
// Each subcommand lives in its own file. Global flags (session specifiers,
// filters, output mode) are defined here and inherited by all subcommands.
package cmd

import (
	"flag"
	"fmt"
	"os"
)

// Global flags shared across all subcommands.
var (
	flagProject    string
	flagSession    string
	flagLatest     bool
	flagAll        bool
	flagSubagents  bool
	flagJSON       bool
	flagJSONL      bool
	flagLimit      int
	flagNoTruncate bool

	// Filters
	flagRole        string
	flagType        string
	flagTool        string
	flagContains    string
	flagAudit       string
	flagSince       string
	flagUntil       string
	flagExternalOnly bool
	flagNoSidechain  bool
)

// registerGlobalFlags adds the shared flags to a FlagSet.
func registerGlobalFlags(fs *flag.FlagSet) {
	// Session specifiers
	fs.StringVar(&flagProject, "project", "", "Project name")
	fs.StringVar(&flagSession, "session", "", "Session UUID prefix")
	fs.BoolVar(&flagLatest, "latest", false, "Most recent session")
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
	fs.StringVar(&flagSince, "since", "", "Records after this ISO timestamp")
	fs.StringVar(&flagUntil, "until", "", "Records before this ISO timestamp")
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

	return fmt.Errorf("unknown command %q — run 'dq help' for usage", name)
}

func printUsage() {
	fmt.Println("dq — Doreen Query: transcript analysis tool")
	fmt.Println()
	fmt.Println("Usage: dq <command> [flags] [session-specifier]")
	fmt.Println()
	fmt.Println("Discovery:")
	fmt.Println("  sessions    List sessions with metadata")
	fmt.Println("  agents      List subagents for a session")
	fmt.Println("  find        Search sessions/agents by content")
	fmt.Println()
	fmt.Println("Navigation:")
	fmt.Println("  show        Show records around a point of interest")
	fmt.Println("  walk        Step through a transcript turn by turn")
	fmt.Println()
	fmt.Println("Query:")
	fmt.Println("  messages    Extract messages with filtering")
	fmt.Println("  tools       Extract tool calls with filtering")
	fmt.Println("  raw         Dump raw JSONL records")
	fmt.Println()
	fmt.Println("Analysis:")
	fmt.Println("  stats       Session summary statistics")
	fmt.Println("  tokens      Token consumption analysis")
	fmt.Println("  compactions Detect compaction events")
	fmt.Println("  errors      Extract errors and failures")
	fmt.Println("  audit       Tool use audit (for grading)")
	fmt.Println("  critique-data  Extract data for transcript critique grader")
	fmt.Println("  agent-trace    Trace subagent lifecycles")
	fmt.Println()
	fmt.Println("Global flags:")
	fmt.Println("  --project NAME    Project name")
	fmt.Println("  --session UUID    Session by UUID prefix")
	fmt.Println("  --latest          Most recent session")
	fmt.Println("  --all             All projects")
	fmt.Println("  --subagents       Include subagent transcripts")
	fmt.Println("  --json            Structured JSON output")
	fmt.Println("  --jsonl           Streaming JSONL output")
	fmt.Println("  --limit N         Max records to output")
	fmt.Println("  --no-truncate     Show full content")
	fmt.Println()
	fmt.Println("See doreen/docs/tools/transcript-query.md for full documentation.")
}
