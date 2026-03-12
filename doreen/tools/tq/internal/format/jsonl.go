package format

import (
	"encoding/json"
	"fmt"
	"os"
)

// JSONLWriter writes streaming JSONL output (one JSON object per line).
type JSONLWriter struct{}

// NewJSONLWriter creates a JSONL writer to stdout.
func NewJSONLWriter() *JSONLWriter {
	return &JSONLWriter{}
}

// WriteLine writes a single value as one JSON line.
func (jw *JSONLWriter) WriteLine(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(os.Stdout, "%s\n", data)
	return err
}
