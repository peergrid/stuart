#!/bin/bash
# PreToolUse:Bash — Block sleep+check polling patterns and repeated identical commands.
# Sleep detection is stateless (pattern match on the command).
# Repeated command detection uses tq to scan recent Bash calls from the transcript.
# jq is used only for parsing hook input JSON (the hook API), not for transcripts.
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

# --- Detect sleep+check patterns (stateless) ---

# sleep followed by a check command on the same line or in a sequence
if echo "$CMD" | grep -qE '\bsleep\s+[0-9]'; then
    # Exception: external process checks that don't have completion notifications
    if echo "$CMD" | grep -qE '(gh\s+run\s+(view|watch)|curl\s.*localhost|curl\s.*127\.0\.0\.1)' && ! echo "$CMD" | grep -qE 'while|for|until|loop'; then
        exit 0
    fi

    emit_block "BLOCKED: Do not use sleep for polling. Use 'run_in_background' and wait for the completion notification. If you need to check an external process, use a single check command without sleep."
fi

# while/until loops with sleep (polling loops)
if echo "$CMD" | grep -qE '(while|until)\s.*sleep|sleep.*&&.*(while|until)'; then
    emit_block "BLOCKED: Polling loop detected (while/until + sleep). Use 'run_in_background' for long-running tasks and wait for completion notification."
fi

# --- Detect repeated identical commands via tq ---

CMD_HASH=$(echo "$CMD" | md5sum | cut -d' ' -f1)

# Get last 10 Bash commands from the transcript via tq --field (no jq needed)
REPEAT_COUNT=$("$TQ" tools --tool Bash --limit 10 --field input.command 2>/dev/null | while read -r line; do
    echo "$line" | md5sum | cut -d' ' -f1
done | grep -c "^${CMD_HASH}$" 2>/dev/null || echo 0)

if [ "$REPEAT_COUNT" -ge 3 ]; then
    log_event "polling-patterns" "repeated-command" "hash=$CMD_HASH count=$REPEAT_COUNT"
    emit_block "BLOCKED: You have executed this same command $REPEAT_COUNT times recently. This looks like polling. If waiting for a background task, you will be notified when it finishes. If checking external status, run the check once, not in a loop."
fi

exit 0
