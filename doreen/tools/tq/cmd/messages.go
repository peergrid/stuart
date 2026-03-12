package cmd

import (
	"flag"
	"fmt"
)

func init() {
	register("messages", "Extract messages with filtering", runMessages)
}

func runMessages(args []string) error {
	fs := flag.NewFlagSet("messages", flag.ExitOnError)
	registerBatchFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}

	// TODO: implement message extraction
	// 1. Load and filter records (streaming)
	// 2. For user/assistant records, extract text content
	// 3. For tool-only turns, show tool call summaries
	// 4. Apply --contains, --role, --external-only, time filters
	// 5. Output: [timestamp] ROLE: text (truncated unless --no-truncate)
	fmt.Fprintln(flag.CommandLine.Output(), "dq messages: not yet implemented")
	return nil
}
