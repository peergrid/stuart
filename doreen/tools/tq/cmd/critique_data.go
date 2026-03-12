package cmd

import (
	"flag"
	"fmt"
)

func init() {
	register("critique-data", "Extract data for transcript critique grader", runCritiqueData)
}

func runCritiqueData(args []string) error {
	fs := flag.NewFlagSet("critique-data", flag.ExitOnError)
	registerGlobalFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}

	// TODO: implement critique data extraction
	// 1. Load all records into list
	// 2. Extract original_prompt: first external user message
	// 3. Extract decision_points: assistant turns that precede a sequence
	//    of tool calls (chose an approach)
	// 4. Extract error_recovery: error records + N subsequent turns
	// 5. Extract completion_claims: assistant messages containing
	//    done/complete/finished + what followed
	// 6. Extract communication_signals: response lengths per turn,
	//    ratio of text to tool calls
	// 7. Extract context_management: compactions, repeated reads,
	//    file reads after compaction (re-acquiring lost context)
	// 8. Output as JSON (always JSON — this is grader input)
	fmt.Fprintln(flag.CommandLine.Output(), "dq critique-data: not yet implemented")
	return nil
}
