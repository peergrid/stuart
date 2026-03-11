#!/bin/bash
# PreToolUse — Warn root agent when it takes 3+ consecutive turns without
# going idle to receive operator input. Each "turn" is one assistant response
# cycle. Multiple tool calls in one turn (batching) is fine. The problem is
# multiple turns in a row, which blocks queued operator messages.
# Subagents are exempt. Warned flag resets on UserPromptSubmit.
set -euo pipefail

source "$(dirname "$0")/lib.sh"
install_error_trap "undelegated-exploration.sh"

# Subagents run autonomously — no turn limit
is_root_agent || exit 0

INPUT=$(cat)
SESSION_ID=$(echo "$INPUT" | jq -r '.session_id // "unknown"' 2>/dev/null)

STATE_DIR=$(state_dir)
WARNED_FILE="$STATE_DIR/turn-warned-${SESSION_ID}"

# Already warned this run — don't nag
[ -f "$WARNED_FILE" ] && exit 0

TRANSCRIPT=$(echo "$INPUT" | jq -r '.transcript_path // ""' 2>/dev/null)
[ -n "$TRANSCRIPT" ] && [ -f "$TRANSCRIPT" ] || exit 0

THRESHOLD=3

# Count assistant turns since last human user message in transcript.
# Human messages have content blocks with type "text".
# Tool-result messages have content blocks with type "tool_result" only.
TURNS=$(tail -200 "$TRANSCRIPT" | jq -s '
  reverse |
  (to_entries
   | map(select(
       .value.type == "user"
       and (.value.message.content // [] | any(.type == "text"))
     ))
   | .[0].key // length
  ) as $human_idx |
  [.[:$human_idx] | .[] | select(.type == "assistant")] | length
' 2>/dev/null || echo 0)

if [ "$TURNS" -ge "$THRESHOLD" ]; then
    touch "$WARNED_FILE" 2>/dev/null
    log_event "undelegated-exploration" "warn" "turns=$TURNS session=$SESSION_ID"
    emit_warning "You have taken $TURNS consecutive turns without going idle. Delegate remaining work to a background agent so you remain available for operator input." "PreToolUse"
fi

exit 0
