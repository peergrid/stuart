#!/bin/bash
# PreToolUse — Warn root agent when it takes 3+ consecutive turns without
# going idle to receive operator input. Subagents are exempt.
# All state is derived from the transcript — no files, no counters.
set -euo pipefail

source "$(dirname "$0")/lib.sh"
install_error_trap "undelegated-exploration.sh"

is_root_agent || exit 0

INPUT=$(cat)
TRANSCRIPT=$(echo "$INPUT" | jq -r '.transcript_path // ""' 2>/dev/null)
[ -n "$TRANSCRIPT" ] && [ -f "$TRANSCRIPT" ] || exit 0

THRESHOLD=3

# From the transcript tail, extract everything since the last human message:
# - Count assistant turns
# - Check if we already warned (our warning text appears in the transcript)
read -r TURNS WARNED < <(tail -200 "$TRANSCRIPT" | jq -rs '
  reverse |
  (to_entries
   | map(select(
       .value.type == "user"
       and (.value.message.content // [] | any(.type == "text"))
     ))
   | .[0].key // length
  ) as $human_idx |
  .[:$human_idx] as $recent |
  ([$recent[] | select(.type == "assistant")] | length) as $turns |
  ([$recent[] | select(
     .type == "system" or .type == "assistant"
   ) | .message.content // [] | .[] | .text // "" |
   contains("consecutive turns without going idle")] | any) as $warned |
  "\($turns) \(if $warned then 1 else 0 end)"
' 2>/dev/null || echo "0 0")

[ "$WARNED" = "1" ] && exit 0

if [ "$TURNS" -ge "$THRESHOLD" ]; then
    log_event "undelegated-exploration" "warn" "turns=$TURNS"
    emit_warning "You have taken $TURNS consecutive turns without going idle. Delegate remaining work to a background agent so you remain available for operator input." "PreToolUse"
fi

exit 0
