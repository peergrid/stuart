package cmd

import (
	"flag"
	"fmt"
)

func init() {
	register("errors", "Extract errors and failures", runErrors)
}

func runErrors(args []string) error {
	fs := flag.NewFlagSet("errors", flag.ExitOnError)
	errorType := fs.String("error-type", "", "Filter by error type: tool_error|hook_error|permission_denied|runtime_error")
	withContext := fs.Int("with-context", 0, "Show N turns around each error")
	registerBatchFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	_, _ = errorType, withContext

	// TODO: implement error extraction
	// 1. Load records (streaming, or list if --with-context)
	// 2. Detect tool_result with is_error=true
	// 3. Detect hook errors (hook_errors/exit_code/hook failed in text)
	// 4. Detect permission denials
	// 5. Detect runtime errors (MaxFileReadTokenExceededError, EPERM, etc.)
	// 6. Group by error type
	// 7. If --with-context: show surrounding turns
	// 8. Output: grouped error list with timestamps and details
	fmt.Fprintln(flag.CommandLine.Output(), "dq errors: not yet implemented")
	return nil
}
