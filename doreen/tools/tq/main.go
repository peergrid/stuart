// tq — Doreen Query: transcript analysis tool for Claude Code sessions.
//
// Treats JSONL transcript files as a queryable database with cursor-based
// navigation, composable filters, and built-in analysis modes that serve
// doreen's grading system.
//
// See doreen/docs/tools/transcript-query.md for the full specification.
package main

import (
	"fmt"
	"os"

	"stuart/doreen/tools/tq/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "tq: %v\n", err)
		os.Exit(1)
	}
}
