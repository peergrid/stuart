#!/bin/bash
# PreToolUse — Warn root agent when it takes 3+ consecutive turns without
# going idle to receive operator input. Subagents are exempt.
# State is derived from the transcript via tq — no files, no counters, no jq.
set -euo pipefail

source "$(dirname "$0")/lib.sh"
install_error_trap "undelegated-work.sh"

is_root_agent || exit 0

THRESHOLD=3

# Check if we already warned since the last human message
"$TQ" walk --reverse --until user --contains "consecutive turns without going idle" --exists 2>/dev/null && exit 0

# Count assistant turns since last human message
TURNS=$("$TQ" walk --reverse --until user --role assistant --count 2>/dev/null) || exit 0

if [ "$TURNS" -ge "$THRESHOLD" ]; then
    log_event "undelegated-work" "warn" "turns=$TURNS"
    emit_warning "You have taken $TURNS consecutive turns without going idle. Delegate remaining work to a background agent so you remain available for operator input." "PreToolUse"
fi

exit 0
