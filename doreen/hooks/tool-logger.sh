#!/bin/bash
# PreToolUse hook: log all tool calls for observability.
# Always ALLOW -- logging only, fail silently.
# Wired to: PreToolUse (all tools)
set -euo pipefail
source "$(dirname "$0")/lib.sh"
install_error_trap "tool-logger.sh"

INPUT=$(cat)
TOOL_NAME=$(echo "$INPUT" | jq -r '.tool_name // "unknown"' 2>/dev/null)
SESSION_ID=$(echo "$INPUT" | jq -r '.session_id // "unknown"' 2>/dev/null)
CMD_PREFIX=$(echo "$INPUT" | jq -r '.tool_input.command // ""' 2>/dev/null | cut -c1-100)

# Detect agent type via TTY
if is_root_agent; then
    AGENT_TYPE="root"
else
    AGENT_TYPE="subagent"
fi

LOG_FILE="$(log_dir)/tool-usage.jsonl"
jq -n -c \
    --arg ts "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    --arg tool "$TOOL_NAME" \
    --arg sid "$SESSION_ID" \
    --arg cmd "$CMD_PREFIX" \
    --arg agent "$AGENT_TYPE" \
    '{timestamp:$ts,tool_name:$tool,session_id:$sid,command_prefix:$cmd,agent_type:$agent}' \
    >> "$LOG_FILE" 2>/dev/null

exit 0
