package cmd

import (
	"flag"
	"fmt"
)

func init() {
	register("show", "Show records around a point of interest", runShow)
}

func runShow(args []string) error {
	fs := flag.NewFlagSet("show", flag.ExitOnError)
	turn := fs.Int("turn", 0, "Jump to turn number")
	at := fs.String("at", "", "Jump to timestamp")
	first := fs.String("first", "", "Jump to first occurrence of: compaction, error, agent-launch, tool:NAME, pattern:REGEX")
	last := fs.String("last", "", "Jump to last occurrence")
	nth := fs.Int("nth", 0, "Jump to Nth occurrence (use with target type as positional arg)")
	context := fs.Int("context", 3, "Number of turns before and after to show")
	registerGlobalFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	_, _, _, _, _, _ = turn, at, first, last, nth, context

	// TODO: implement cursor navigation
	// 1. Load transcript into indexed record list
	// 2. Resolve anchor: --turn, --at, --first/--last/--nth + target
	// 3. Calculate context window (anchor - N .. anchor + N)
	// 4. Apply filters to context window
	// 5. Render each turn in the window
	fmt.Fprintln(flag.CommandLine.Output(), "dq show: not yet implemented")
	return nil
}
