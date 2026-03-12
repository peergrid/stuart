package cmd

import (
	"flag"
	"fmt"
)

func init() {
	register("sessions", "List sessions with metadata", runSessions)
}

func runSessions(args []string) error {
	fs := flag.NewFlagSet("sessions", flag.ExitOnError)
	registerGlobalFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}

	// TODO: implement session listing
	// 1. Resolve project directory from --project flag
	// 2. Find all *.jsonl files (top-level only, not subagents)
	// 3. Extract metadata from each: session ID, size, record count,
	//    model, subagent count, first/last timestamp, duration
	// 4. Output as table (default) or JSON (--json)
	fmt.Fprintln(flag.CommandLine.Output(), "dq sessions: not yet implemented")
	return nil
}
