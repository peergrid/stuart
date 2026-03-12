package cmd

import (
	"flag"
	"fmt"
)

func init() {
	register("tokens", "Token consumption analysis", runTokens)
}

func runTokens(args []string) error {
	fs := flag.NewFlagSet("tokens", flag.ExitOnError)
	summary := fs.Bool("summary", false, "Summary metrics only (no per-turn timeline)")
	registerBatchFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	_ = summary

	// TODO: implement token analysis
	// 1. Load assistant records (streaming)
	// 2. For each: extract usage fields (input, output, cache_read, cache_create)
	// 3. Compute per-turn: total_context, cumulative_input, cumulative_output
	// 4. Track tool names per turn for annotation
	// 5. If --summary: compute aggregates (total, max, average, compaction_count)
	// 6. Output: table with TURN/TIME/CONTEXT/INPUT/OUTPUT/CACHE_R/CUM_IN/TOOLS
	//    or JSON with per-turn array + summary
	fmt.Fprintln(flag.CommandLine.Output(), "dq tokens: not yet implemented")
	return nil
}
