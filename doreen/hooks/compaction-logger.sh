#!/bin/bash
# SessionStart(compact) hook: log compaction events for observability.
# Wired to: SessionStart (source=compact)
set -euo pipefail
source "$(dirname "$0")/lib.sh"
install_error_trap "compaction-logger.sh"

INPUT=$(cat)
SOURCE=$(echo "$INPUT" | jq -r '.source // empty')

# Only fire on compaction events
if [ "$SOURCE" != "compact" ]; then
    exit 0
fi

SESSION_ID=$(echo "$INPUT" | jq -r '.session_id // "unknown"' 2>/dev/null)

# Count transcript records if path available
TRANSCRIPT_PATH=$(echo "$INPUT" | jq -r '.transcript_path // ""' 2>/dev/null)
RECORD_COUNT=0
if [ -n "$TRANSCRIPT_PATH" ] && [ -f "$TRANSCRIPT_PATH" ]; then
    RECORD_COUNT=$(wc -l < "$TRANSCRIPT_PATH" 2>/dev/null || echo 0)
fi

LOG_FILE="$(log_dir)/compactions.jsonl"
jq -n -c \
    --arg ts "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    --arg sid "$SESSION_ID" \
    --argjson records "$RECORD_COUNT" \
    '{timestamp:$ts,session_id:$sid,record_count:$records}' \
    >> "$LOG_FILE" 2>/dev/null

exit 0
