#!/bin/bash
# PreToolUse:Bash — Block sleep+check polling patterns
# Detects: sleep followed by check commands, repeated identical commands
set -euo pipefail

source "$(dirname "$0")/lib.sh"
install_error_trap "polling-patterns.sh"

INPUT=$(cat)
TOOL_NAME=$(echo "$INPUT" | jq -r '.tool_name // "unknown"' 2>/dev/null)

# Only check Bash tool
if [ "$TOOL_NAME" != "Bash" ]; then
    exit 0
fi

CMD=$(echo "$INPUT" | jq -r '.tool_input.command // ""' 2>/dev/null)
[ -z "$CMD" ] && exit 0

SESSION_ID=$(echo "$INPUT" | jq -r '.session_id // "unknown"' 2>/dev/null)
STATE_DIR=$(state_dir)
HISTORY_FILE="$STATE_DIR/cmd-history-${SESSION_ID}"

# --- Detect sleep+check patterns ---

# sleep followed by a check command on the same line or in a sequence
if echo "$CMD" | grep -qE '\bsleep\s+[0-9]'; then
    # Exception: external process checks that don't have completion notifications
    # e.g., gh run view, curl to check if a service is up
    if echo "$CMD" | grep -qE '(gh\s+run\s+(view|watch)|curl\s.*localhost|curl\s.*127\.0\.0\.1)' && ! echo "$CMD" | grep -qE 'while|for|until|loop'; then
        # Single check after sleep is okay, but not in a loop
        exit 0
    fi

    emit_block "BLOCKED: Do not use sleep for polling. Use 'run_in_background' and wait for the completion notification. If you need to check an external process, use a single check command without sleep."
fi

# while/until loops with sleep (polling loops)
if echo "$CMD" | grep -qE '(while|until)\s.*sleep|sleep.*&&.*(while|until)'; then
    emit_block "BLOCKED: Polling loop detected (while/until + sleep). Use 'run_in_background' for long-running tasks and wait for completion notification."
fi

# --- Detect repeated identical commands ---

# Record current command
mkdir -p "$STATE_DIR" 2>/dev/null
CMD_HASH=$(echo "$CMD" | md5sum | cut -d' ' -f1)
TIMESTAMP=$(date +%s)

# Append current command with timestamp
echo "$TIMESTAMP $CMD_HASH" >> "$HISTORY_FILE" 2>/dev/null

# Keep only last 20 entries to prevent file growth
if [ -f "$HISTORY_FILE" ]; then
    tail -20 "$HISTORY_FILE" > "$HISTORY_FILE.tmp" 2>/dev/null && mv "$HISTORY_FILE.tmp" "$HISTORY_FILE" 2>/dev/null
fi

# Check for repeated identical commands (same hash 3+ times in last 10 entries)
if [ -f "$HISTORY_FILE" ]; then
    REPEAT_COUNT=$(tail -10 "$HISTORY_FILE" | awk '{print $2}' | grep -c "^${CMD_HASH}$" 2>/dev/null || true)
    if [ "$REPEAT_COUNT" -ge 3 ]; then
        log_event "polling-patterns" "repeated-command" "hash=$CMD_HASH count=$REPEAT_COUNT"
        emit_block "BLOCKED: You have executed this same command $REPEAT_COUNT times recently. This looks like polling. If waiting for a background task, you will be notified when it finishes. If checking external status, run the check once, not in a loop."
    fi
fi

exit 0
