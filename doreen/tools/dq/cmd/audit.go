package cmd

import (
	"flag"
	"fmt"
)

func init() {
	register("audit", "Tool use audit for grading", runAudit)
}

func runAudit(args []string) error {
	fs := flag.NewFlagSet("audit", flag.ExitOnError)
	registerGlobalFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}

	// TODO: implement tool use audit
	// 1. Load all records into list (needs cross-reference)
	// 2. Build tool call sequence: [(turn, tool_name, input, result)]
	// 3. Check: dedicated tool preference
	//    - Bash commands containing cat/grep/find/sed/awk/echo
	//    - Map each to the preferred tool: Read/Grep/Glob/Edit/Write
	// 4. Check: read-before-edit
	//    - For each Edit/Write, find the file path in input
	//    - Check that a Read of the same file appears earlier (after last edit)
	// 5. Check: no redundant reads
	//    - Track read files; flag same file read twice without intervening edit
	// 6. Check: parallelism ratio
	//    - Count turns with multiple tool_use blocks vs single tool_use
	// 7. Check: agent delegation
	//    - If >10 exploratory tool calls (Grep/Glob/Read without Edit),
	//      flag if no Agent launch occurred
	// 8. Compute score: 1.0 - (violations / total_tool_calls)
	// 9. Output: structured report with violations, turn refs, and score
	fmt.Fprintln(flag.CommandLine.Output(), "dq audit: not yet implemented")
	return nil
}
