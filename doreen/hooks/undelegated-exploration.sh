#!/bin/bash
# PostToolUse — Warn when main agent does too much codebase exploration without delegating
# Tracks Read/Grep/Glob calls. Warns at 8+ exploratory calls without Agent launch.
set -euo pipefail

source "$(dirname "$0")/lib.sh"
install_error_trap "undelegated-exploration.sh"

INPUT=$(cat)
TOOL_NAME=$(echo "$INPUT" | jq -r '.tool_name // "unknown"' 2>/dev/null)

SESSION_ID=$(echo "$INPUT" | jq -r '.session_id // "unknown"' 2>/dev/null)
STATE_DIR=$(state_dir)
EXPLORE_FILE="$STATE_DIR/explore-count-${SESSION_ID}"
RECENT_FILE="$STATE_DIR/explore-recent-${SESSION_ID}"

THRESHOLD=8

# If an Agent was launched, reset the exploration counter
# The agent is doing the exploration, which is exactly what we want
if [ "$TOOL_NAME" = "Agent" ] || [ "$TOOL_NAME" = "SubAgent" ]; then
    echo "0" > "$EXPLORE_FILE" 2>/dev/null
    # Clear recent files list
    > "$RECENT_FILE" 2>/dev/null
    exit 0
fi

# If this is an Edit or Write, the previous Read was likely directed, not exploratory
# Remove the last entry from the counter (give credit for directed reads)
if [ "$TOOL_NAME" = "Edit" ] || [ "$TOOL_NAME" = "Write" ]; then
    if [ -f "$EXPLORE_FILE" ]; then
        COUNT=$(cat "$EXPLORE_FILE" 2>/dev/null || echo 0)
        if echo "$COUNT" | grep -qE '^[0-9]+$' && [ "$COUNT" -gt 0 ]; then
            COUNT=$((COUNT - 1))
            echo "$COUNT" > "$EXPLORE_FILE" 2>/dev/null
        fi
    fi
    exit 0
fi

# Only track exploratory tools: Read, Grep, Glob
case "$TOOL_NAME" in
    Read|Grep|Glob) ;;
    *) exit 0 ;;
esac

# Track the file/pattern being accessed for heuristic analysis
ACCESSED=""
case "$TOOL_NAME" in
    Read)
        ACCESSED=$(echo "$INPUT" | jq -r '.tool_input.file_path // ""' 2>/dev/null)
        ;;
    Grep)
        ACCESSED=$(echo "$INPUT" | jq -r '.tool_input.pattern // ""' 2>/dev/null)
        ;;
    Glob)
        ACCESSED=$(echo "$INPUT" | jq -r '.tool_input.pattern // ""' 2>/dev/null)
        ;;
esac

# Record the access
if [ -n "$ACCESSED" ]; then
    echo "$ACCESSED" >> "$RECENT_FILE" 2>/dev/null
    # Keep file manageable
    if [ -f "$RECENT_FILE" ]; then
        tail -30 "$RECENT_FILE" > "$RECENT_FILE.tmp" 2>/dev/null && mv "$RECENT_FILE.tmp" "$RECENT_FILE" 2>/dev/null
    fi
fi

# Heuristic: if recent accesses are all to the same 2-3 files, this is narrow/directed work
if [ -f "$RECENT_FILE" ]; then
    UNIQUE_COUNT=$(sort -u "$RECENT_FILE" 2>/dev/null | wc -l)
    if [ "$UNIQUE_COUNT" -le 3 ]; then
        # Working within a small scope — not broad exploration
        exit 0
    fi
fi

# Increment exploration counter
COUNT=0
if [ -f "$EXPLORE_FILE" ]; then
    COUNT=$(cat "$EXPLORE_FILE" 2>/dev/null || echo 0)
    if ! echo "$COUNT" | grep -qE '^[0-9]+$'; then
        COUNT=0
    fi
fi

COUNT=$((COUNT + 1))
echo "$COUNT" > "$EXPLORE_FILE" 2>/dev/null

# Check threshold
if [ "$COUNT" -ge "$THRESHOLD" ]; then
    log_event "undelegated-exploration" "warn" "count=$COUNT tool=$TOOL_NAME session=$SESSION_ID"
    # Reset after warning; next warning fires after another THRESHOLD calls
    echo "0" > "$EXPLORE_FILE" 2>/dev/null
    > "$RECENT_FILE" 2>/dev/null
    emit_warning "You have made $COUNT exploratory tool calls (Read/Grep/Glob) in the main context without delegating to an Agent. Consider launching an Explore agent to preserve main context for decision-making and execution." "PostToolUse"
fi

exit 0
