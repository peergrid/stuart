package cmd

import (
	"flag"
	"fmt"
)

func init() {
	register("find", "Search sessions/agents by content", runFind)
}

func runFind(args []string) error {
	fs := flag.NewFlagSet("find", flag.ExitOnError)
	inMessages := fs.Bool("in-messages", false, "Search in message content")
	registerBatchFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	_ = inMessages

	// TODO: implement search
	// 1. Accept a search pattern as positional arg
	// 2. If --in-messages: scan message text across sessions
	// 3. Otherwise: search agent names and launch prompts
	// 4. Return matching sessions/agents with metadata
	fmt.Fprintln(flag.CommandLine.Output(), "dq find: not yet implemented")
	return nil
}
