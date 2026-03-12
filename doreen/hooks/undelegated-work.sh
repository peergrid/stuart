#!/bin/bash
# PreToolUse — Warn root agent when it takes 3+ consecutive turns without
# going idle to receive operator input. Subagents are exempt.
# State is derived from the transcript via tq — no files, no counters.
set -euo pipefail

source "$(dirname "$0")/lib.sh"
install_error_trap "undelegated-work.sh"

is_root_agent || exit 0

THRESHOLD=3

# Get recent turns in reverse. We need to find:
# 1. How many assistant turns since the last human message
# 2. Whether we already issued a warning in that span
RECENT=$("$TQ" walk --reverse --limit 50 --jsonl 2>/dev/null) || exit 0
[ -z "$RECENT" ] && exit 0

read -r TURNS WARNED < <(echo "$RECENT" | jq -rs '
  # Find first genuine human message (not tool results)
  (to_entries
   | map(select(
       .value.role == "user"
       and (.value.content // [] | any(.type == "text"))
     ))
   | .[0].key // length
  ) as $human_idx |
  .[:$human_idx] as $recent |
  ([$recent[] | select(.role == "assistant")] | length) as $turns |
  ([$recent[] | .content // [] | .[] | .text // "" |
    contains("consecutive turns without going idle")] | any) as $warned |
  "\($turns) \(if $warned then 1 else 0 end)"
' 2>/dev/null || echo "0 0")

[ "$WARNED" = "1" ] && exit 0

if [ "$TURNS" -ge "$THRESHOLD" ]; then
    log_event "undelegated-work" "warn" "turns=$TURNS"
    emit_warning "You have taken $TURNS consecutive turns without going idle. Delegate remaining work to a background agent so you remain available for operator input." "PreToolUse"
fi

exit 0
