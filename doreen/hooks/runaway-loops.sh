#!/bin/bash
# PostToolUse — Warn when too many consecutive tool calls without user-facing output
# Tracks consecutive tool calls. Warns (does not block) at threshold of 10+.
set -euo pipefail

source "$(dirname "$0")/lib.sh"
install_error_trap "runaway-loops.sh"

INPUT=$(cat)
TOOL_NAME=$(echo "$INPUT" | jq -r '.tool_name // "unknown"' 2>/dev/null)

# Exception: Agent tool calls — agents are expected to run autonomously
if [ "$TOOL_NAME" = "Agent" ] || [ "$TOOL_NAME" = "SubAgent" ]; then
    exit 0
fi

SESSION_ID=$(echo "$INPUT" | jq -r '.session_id // "unknown"' 2>/dev/null)
STATE_DIR=$(state_dir)
COUNTER_FILE="$STATE_DIR/tool-count-${SESSION_ID}"

THRESHOLD=10

# Read current count
COUNT=0
if [ -f "$COUNTER_FILE" ]; then
    COUNT=$(cat "$COUNTER_FILE" 2>/dev/null || echo 0)
    # Validate it's a number
    if ! echo "$COUNT" | grep -qE '^[0-9]+$'; then
        COUNT=0
    fi
fi

# Increment
COUNT=$((COUNT + 1))
echo "$COUNT" > "$COUNTER_FILE" 2>/dev/null

# Check threshold
if [ "$COUNT" -ge "$THRESHOLD" ]; then
    log_event "runaway-loops" "warn" "count=$COUNT tool=$TOOL_NAME session=$SESSION_ID"
    # Reset counter after warning so we don't warn on every subsequent call
    # Next warning will fire after another THRESHOLD calls
    echo "0" > "$COUNTER_FILE" 2>/dev/null
    emit_warning "You have made $COUNT consecutive tool calls without communicating with the operator. Pause and provide a status update on what you have accomplished and what you plan to do next." "PostToolUse"
fi

exit 0
