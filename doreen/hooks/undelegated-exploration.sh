#!/bin/bash
# PostToolUse — Warn root agent when it makes too many exploratory tool calls
# between user messages. Subagents are exempt — they exist TO explore.
# Counter resets on each user message (via exploration-reset.sh on UserPromptSubmit).
set -euo pipefail

source "$(dirname "$0")/lib.sh"
install_error_trap "undelegated-exploration.sh"

# Subagents can explore freely — this hook is only for the root agent
is_root_agent || exit 0

INPUT=$(cat)
TOOL_NAME=$(echo "$INPUT" | jq -r '.tool_name // "unknown"' 2>/dev/null)

# Only track exploratory tools: Read, Grep, Glob
case "$TOOL_NAME" in
    Read|Grep|Glob) ;;
    *) exit 0 ;;
esac

SESSION_ID=$(echo "$INPUT" | jq -r '.session_id // "unknown"' 2>/dev/null)
STATE_DIR=$(state_dir)
EXPLORE_FILE="$STATE_DIR/explore-count-${SESSION_ID}"

THRESHOLD=3

# Increment counter
COUNT=0
if [ -f "$EXPLORE_FILE" ]; then
    COUNT=$(cat "$EXPLORE_FILE" 2>/dev/null || echo 0)
    if ! echo "$COUNT" | grep -qE '^[0-9]+$'; then
        COUNT=0
    fi
fi
COUNT=$((COUNT + 1))
echo "$COUNT" > "$EXPLORE_FILE" 2>/dev/null

# Warn when threshold exceeded
if [ "$COUNT" -gt "$THRESHOLD" ]; then
    log_event "undelegated-exploration" "warn" "count=$COUNT tool=$TOOL_NAME session=$SESSION_ID"
    # Reset so next warning fires after another THRESHOLD calls
    echo "0" > "$EXPLORE_FILE" 2>/dev/null
    emit_warning "You have made $COUNT exploratory tool calls (Read/Grep/Glob) in this turn without delegating. Launch an Explore or general-purpose agent instead of consuming the main context." "PostToolUse"
fi

exit 0
