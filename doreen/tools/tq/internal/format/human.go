// Package format provides output formatters for dq's three output modes:
// human-readable tables, structured JSON, and streaming JSONL.
package format

import (
	"fmt"
	"strings"
	"time"
)

// Truncate shortens text to maxLen characters, adding ellipsis if truncated.
func Truncate(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}

// FormatTimestamp converts an ISO timestamp to compact form (MM-DD HH:MM:SS).
func FormatTimestamp(ts string) string {
	t, err := time.Parse(time.RFC3339Nano, ts)
	if err != nil {
		// Try without timezone
		t, err = time.Parse("2006-01-02T15:04:05.000Z", ts)
		if err != nil {
			return ts
		}
	}
	return t.Format("01-02 15:04:05")
}

// FormatDuration formats seconds into a human-readable duration.
func FormatDuration(seconds int) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	if seconds < 3600 {
		return fmt.Sprintf("%dm %ds", seconds/60, seconds%60)
	}
	return fmt.Sprintf("%dh %dm %ds", seconds/3600, (seconds%3600)/60, seconds%60)
}

// FormatSize formats bytes into a human-readable size.
func FormatSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%dB", bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.1fKB", float64(bytes)/1024)
	}
	if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.1fMB", float64(bytes)/(1024*1024))
	}
	return fmt.Sprintf("%.1fGB", float64(bytes)/(1024*1024*1024))
}

// FormatTokens formats a token count with thousand separators.
func FormatTokens(n int) string {
	s := fmt.Sprintf("%d", n)
	if n < 1000 {
		return s
	}
	// Insert commas
	var result strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result.WriteByte(',')
		}
		result.WriteRune(c)
	}
	return result.String()
}

// TableRow formats columns with fixed widths for table output.
func TableRow(widths []int, values []string) string {
	var b strings.Builder
	for i, v := range values {
		if i >= len(widths) {
			b.WriteString(v)
			break
		}
		w := widths[i]
		if w > 0 {
			b.WriteString(fmt.Sprintf("%-*s", w, v))
		} else {
			b.WriteString(fmt.Sprintf("%*s", -w, v))
		}
		if i < len(values)-1 {
			b.WriteByte(' ')
		}
	}
	return b.String()
}
