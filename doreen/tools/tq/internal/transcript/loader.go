package transcript

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// StreamRecords yields records one at a time from the given files.
// Use for operations that don't need random access.
func StreamRecords(files []string, callback func(*Record) error) error {
	turn := 0
	for _, path := range files {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dq: cannot open %s: %v\n", path, err)
			continue
		}
		scanner := bufio.NewScanner(f)
		scanner.Buffer(make([]byte, 0, 1024*1024), 10*1024*1024) // 10MB max line
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Bytes()
			if len(line) == 0 {
				continue
			}
			var r Record
			if err := json.Unmarshal(line, &r); err != nil {
				continue // Skip malformed lines
			}
			turn++
			r.SourceFile = path
			r.TurnNumber = turn
			r.LineNumber = lineNum
			if err := callback(&r); err != nil {
				f.Close()
				return err
			}
		}
		f.Close()
	}
	return nil
}

// LoadAll loads all records from the given files into a slice.
// Use when random access is needed (compaction detection, cursor context).
func LoadAll(files []string) ([]*Record, error) {
	var records []*Record
	err := StreamRecords(files, func(r *Record) error {
		records = append(records, r)
		return nil
	})
	return records, err
}
