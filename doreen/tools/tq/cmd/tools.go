package cmd

import (
	"flag"
	"fmt"
)

func init() {
	register("tools", "Extract tool calls with filtering", runTools)
}

func runTools(args []string) error {
	fs := flag.NewFlagSet("tools", flag.ExitOnError)
	inputContains := fs.String("input-contains", "", "Regex match in tool input")
	withResults := fs.Bool("with-results", false, "Include tool results")
	registerGlobalFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	_, _ = inputContains, withResults

	// TODO: implement tool extraction
	// 1. Load and filter records (streaming)
	// 2. For each assistant record, extract tool_use blocks
	// 3. Apply --tool name filter, --input-contains, --audit field filter
	// 4. If --with-results, find matching tool_result in next user record
	// 5. Output: [timestamp] TOOL name: input_summary
	fmt.Fprintln(flag.CommandLine.Output(), "dq tools: not yet implemented")
	return nil
}
