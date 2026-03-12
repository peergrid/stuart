package cmd

import (
	"flag"
	"fmt"
)

func init() {
	register("walk", "Step through a transcript turn by turn", runWalk)
}

func runWalk(args []string) error {
	fs := flag.NewFlagSet("walk", flag.ExitOnError)
	from := fs.Int("from", 1, "Starting turn number")
	reverse := fs.Bool("reverse", false, "Walk backward from end")
	toolsOnly := fs.Bool("tools-only", false, "Show only tool call turns")
	registerGlobalFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	_, _, _ = from, reverse, toolsOnly

	// TODO: implement walk
	// 1. Load transcript (streaming if forward, list if reverse)
	// 2. Start at --from turn (or end if --reverse)
	// 3. For each turn, apply filters
	// 4. Render: turn header (number, role, timestamp, token count)
	//    then content summary (text blocks, tool calls with names/args)
	// 5. Stop at --limit
	fmt.Fprintln(flag.CommandLine.Output(), "dq walk: not yet implemented")
	return nil
}
