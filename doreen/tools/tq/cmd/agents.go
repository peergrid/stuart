package cmd

import (
	"flag"
	"fmt"
)

func init() {
	register("agents", "List subagents for a session", runAgents)
}

func runAgents(args []string) error {
	fs := flag.NewFlagSet("agents", flag.ExitOnError)
	registerBatchFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}

	// TODO: implement agent listing
	// 1. Resolve session from --session or --latest
	// 2. Find subagent JSONL files in <session-uuid>/subagents/
	// 3. Extract metadata: agent ID, name (from teammate_id), model,
	//    size, duration, tool counts, launch prompt preview
	// 4. Output as table or JSON
	fmt.Fprintln(flag.CommandLine.Output(), "dq agents: not yet implemented")
	return nil
}
