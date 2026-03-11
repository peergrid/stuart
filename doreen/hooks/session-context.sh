#!/bin/bash
# SessionStart hook: inject project context based on agent type.
# Root agent (interactive TTY) gets root-context.md.
# Subagent (no TTY) gets subagent-boundaries.md.
# Wired to: SessionStart:startup, SessionStart:compact
set -euo pipefail
source "$(dirname "$0")/lib.sh"
install_error_trap "session-context.sh"

INPUT=$(cat)
EVENT_NAME=$(echo "$INPUT" | jq -r '.hook_event_name // "SessionStart"' 2>/dev/null)
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

if is_root_agent; then
    CONTEXT_FILE="$SCRIPT_DIR/context/root-context.md"
else
    CONTEXT_FILE="$SCRIPT_DIR/context/subagent-boundaries.md"
fi

[ ! -f "$CONTEXT_FILE" ] && exit 0

CONTENT=$(cat "$CONTEXT_FILE")

jq -n --arg ctx "$CONTENT" --arg event "$EVENT_NAME" '{
  hookSpecificOutput: {
    hookEventName: $event,
    additionalContext: $ctx
  }
}'
exit 0
