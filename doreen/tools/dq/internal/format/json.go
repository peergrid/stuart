package format

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// JSONWriter writes structured JSON output.
type JSONWriter struct {
	w       io.Writer
	encoder *json.Encoder
}

// NewJSONWriter creates a JSON writer to stdout.
func NewJSONWriter() *JSONWriter {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return &JSONWriter{w: os.Stdout, encoder: enc}
}

// WriteValue writes a single JSON value with pretty printing.
func (jw *JSONWriter) WriteValue(v any) error {
	return jw.encoder.Encode(v)
}

// WriteArray writes an array of values as a JSON array.
func (jw *JSONWriter) WriteArray(values []any) error {
	data, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(jw.w, "%s\n", data)
	return err
}
