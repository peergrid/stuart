#!/bin/bash
# PermissionRequest hook: log every permission prompt and compute the allow rule
# that would prevent this approval next time.
# Purely observational — logs and exits 0, no stdout, passes through to normal dialog.
# Target: < 10ms (single jq + append)
set -euo pipefail
source "$(dirname "$0")/lib.sh"
install_error_trap "approval-tracker.sh"

INPUT=$(cat)

TOOL_NAME=$(echo "$INPUT" | jq -r '.tool_name // "unknown"' 2>/dev/null)
SESSION_ID=$(echo "$INPUT" | jq -r '.session_id // "unknown"' 2>/dev/null)
PERM_MODE=$(echo "$INPUT" | jq -r '.permission_mode // "unknown"' 2>/dev/null)

# Extract a meaningful detail depending on tool type
case "$TOOL_NAME" in
    Bash)
        DETAIL=$(echo "$INPUT" | jq -r '.tool_input.command // ""' 2>/dev/null | cut -c1-100)
        ;;
    Edit|Write|Read)
        DETAIL=$(echo "$INPUT" | jq -r '.tool_input.file_path // ""' 2>/dev/null)
        ;;
    Grep|Glob)
        DETAIL=$(echo "$INPUT" | jq -r '.tool_input.pattern // ""' 2>/dev/null | cut -c1-100)
        ;;
    Task)
        DETAIL=$(echo "$INPUT" | jq -r '.tool_input.description // ""' 2>/dev/null | cut -c1-100)
        ;;
    *)
        DETAIL=$(echo "$INPUT" | jq -r '.tool_input | keys[0:2] | join(",")' 2>/dev/null || echo "")
        ;;
esac

# Compute the generalized permission rule that would auto-allow this
case "$TOOL_NAME" in
    Bash)
        # Extract first word of command -> Bash(<first-word> *)
        FIRST_WORD=$(echo "$INPUT" | jq -r '.tool_input.command // ""' 2>/dev/null \
            | sed 's/^[[:space:]]*//' | cut -d' ' -f1 | cut -d'/' -f1)
        if [ -n "$FIRST_WORD" ]; then
            COMPUTED_RULE="Bash($FIRST_WORD *)"
        else
            COMPUTED_RULE="Bash(*)"
        fi
        ;;
    *)
        # Non-Bash tools: just the tool name
        COMPUTED_RULE="$TOOL_NAME"
        ;;
esac

LOG_FILE="$(log_dir)/approvals.jsonl"
jq -cn \
    --arg ts "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    --arg sid "$SESSION_ID" \
    --arg tool "$TOOL_NAME" \
    --arg detail "$DETAIL" \
    --arg rule "$COMPUTED_RULE" \
    --arg mode "$PERM_MODE" \
    '{timestamp:$ts,session_id:$sid,tool_name:$tool,detail:$detail,computed_rule:$rule,permission_mode:$mode}' \
    >> "$LOG_FILE" 2>/dev/null

# CRITICAL: exit 0 with no stdout — pass through to normal permission dialog
exit 0
