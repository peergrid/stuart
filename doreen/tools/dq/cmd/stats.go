package cmd

import (
	"flag"
	"fmt"
)

func init() {
	register("stats", "Session summary statistics", runStats)
}

func runStats(args []string) error {
	fs := flag.NewFlagSet("stats", flag.ExitOnError)
	registerGlobalFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}

	// TODO: implement stats
	// 1. Load all records (streaming accumulation)
	// 2. Accumulate: record count by type/role, tool call counts,
	//    model distribution, token totals (input/output/cache),
	//    first/last timestamps, session IDs, file count
	// 3. Calculate duration from timestamps
	// 4. Output formatted summary or JSON
	fmt.Fprintln(flag.CommandLine.Output(), "dq stats: not yet implemented")
	return nil
}
