#!/bin/bash
# PreToolUse hook: log agent/task launches for observability.
# Always ALLOW -- logging only, fail silently.
# Wired to: PreToolUse on Task tool
set -euo pipefail
source "$(dirname "$0")/lib.sh"
install_error_trap "agent-logger.sh"

INPUT=$(cat)
SESSION_ID=$(echo "$INPUT" | jq -r '.session_id // "unknown"' 2>/dev/null)
SUBAGENT_TYPE=$(echo "$INPUT" | jq -r '.tool_input.subagent_type // "unknown"' 2>/dev/null)
MODEL=$(echo "$INPUT" | jq -r '.tool_input.model // "inherit"' 2>/dev/null)
DESC=$(echo "$INPUT" | jq -r '.tool_input.description // ""' 2>/dev/null | cut -c1-200)
RUN_BG=$(echo "$INPUT" | jq -r '.tool_input.run_in_background // false' 2>/dev/null)

LOG_FILE="$(log_dir)/agent-spawns.jsonl"
jq -n -c \
    --arg ts "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    --arg sid "$SESSION_ID" \
    --arg type "$SUBAGENT_TYPE" \
    --arg model "$MODEL" \
    --arg desc "$DESC" \
    --argjson bg "$RUN_BG" \
    '{timestamp:$ts,session_id:$sid,subagent_type:$type,model:$model,description:$desc,background:$bg}' \
    >> "$LOG_FILE" 2>/dev/null

exit 0
