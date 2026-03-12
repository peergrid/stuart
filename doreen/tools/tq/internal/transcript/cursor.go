package transcript

// Cursor provides indexed navigation over a loaded transcript.
type Cursor struct {
	records []*Record
	pos     int // Current position (0-indexed)
}

// NewCursor creates a cursor over the given records.
func NewCursor(records []*Record) *Cursor {
	return &Cursor{records: records, pos: 0}
}

// Len returns the total number of records.
func (c *Cursor) Len() int {
	return len(c.records)
}

// Pos returns the current position (0-indexed).
func (c *Cursor) Pos() int {
	return c.pos
}

// At returns the record at position i, or nil if out of bounds.
func (c *Cursor) At(i int) *Record {
	if i < 0 || i >= len(c.records) {
		return nil
	}
	return c.records[i]
}

// Seek moves the cursor to position i.
func (c *Cursor) Seek(i int) {
	if i < 0 {
		c.pos = 0
	} else if i >= len(c.records) {
		c.pos = len(c.records) - 1
	} else {
		c.pos = i
	}
}

// SeekToTurn moves the cursor to the record with the given turn number.
func (c *Cursor) SeekToTurn(turn int) bool {
	for i, r := range c.records {
		if r.TurnNumber == turn {
			c.pos = i
			return true
		}
	}
	return false
}

// SeekToTimestamp moves the cursor to the first record at or after the timestamp.
func (c *Cursor) SeekToTimestamp(ts string) bool {
	for i, r := range c.records {
		if r.Timestamp >= ts {
			c.pos = i
			return true
		}
	}
	return false
}

// AnchorTarget specifies what kind of event to seek to.
type AnchorTarget string

const (
	AnchorCompaction  AnchorTarget = "compaction"
	AnchorError       AnchorTarget = "error"
	AnchorAgentLaunch AnchorTarget = "agent-launch"
	AnchorAgentReturn AnchorTarget = "agent-return"
)

// SeekFirst moves the cursor to the first record matching the target.
func (c *Cursor) SeekFirst(target AnchorTarget) bool {
	return c.seekNth(target, 1)
}

// SeekLast moves the cursor to the last record matching the target.
func (c *Cursor) SeekLast(target AnchorTarget) bool {
	// Find all matches, take the last
	lastIdx := -1
	for i, r := range c.records {
		if matchesTarget(r, target) {
			lastIdx = i
		}
	}
	if lastIdx >= 0 {
		c.pos = lastIdx
		return true
	}
	return false
}

// SeekNth moves the cursor to the Nth record matching the target.
func (c *Cursor) seekNth(target AnchorTarget, n int) bool {
	count := 0
	for i, r := range c.records {
		if matchesTarget(r, target) {
			count++
			if count == n {
				c.pos = i
				return true
			}
		}
	}
	return false
}

// Window returns records in [pos-before, pos+after], clamped to bounds.
func (c *Cursor) Window(before, after int) []*Record {
	start := max(c.pos-before, 0)
	end := min(c.pos+after+1, len(c.records))
	return c.records[start:end]
}

// matchesTarget checks if a record matches the given anchor target.
func matchesTarget(r *Record, target AnchorTarget) bool {
	switch target {
	case AnchorCompaction:
		return r.Type == "system" && r.Subtype == "compact_boundary"
	case AnchorError:
		// TODO: check for tool errors, hook errors, permission denials
		return false
	case AnchorAgentLaunch:
		// TODO: check for Agent tool_use in assistant content
		return false
	case AnchorAgentReturn:
		// TODO: check for Agent tool_result in user content
		return false
	default:
		return false
	}
}
