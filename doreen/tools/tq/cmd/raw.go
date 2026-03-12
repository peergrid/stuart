package cmd

import (
	"flag"
	"fmt"
)

func init() {
	register("raw", "Dump raw JSONL records", runRaw)
}

func runRaw(args []string) error {
	fs := flag.NewFlagSet("raw", flag.ExitOnError)
	registerGlobalFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}

	// TODO: implement raw dump
	// 1. Load and filter records (streaming)
	// 2. Write each matching record as one JSON line to stdout
	// 3. Respects all global filters
	// 4. Designed for piping to jq
	fmt.Fprintln(flag.CommandLine.Output(), "dq raw: not yet implemented")
	return nil
}
