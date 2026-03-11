#!/bin/bash
# PostToolUse hook: surface repeated permission approvals as fixable allow rules.
# Fires after every successful tool use. Must be FAST — exits immediately
# when there are no pending approvals for this session.
set -euo pipefail
source "$(dirname "$0")/lib.sh"
install_error_trap "approval-resolver.sh"

INPUT=$(cat)
SESSION_ID=$(echo "$INPUT" | jq -r '.session_id // "unknown"' 2>/dev/null)

LOG_DIR="$(log_dir)"
APPROVALS_FILE="$LOG_DIR/approvals.jsonl"
SURFACED_FILE="$LOG_DIR/.surfaced-${SESSION_ID}"

# Fast exit: no approvals log means nothing to do
[ ! -f "$APPROVALS_FILE" ] && exit 0

# Fast exit: filter to current session, bail if fewer than 2 entries total
SESSION_ENTRIES=$(grep -c "\"session_id\":\"$SESSION_ID\"" "$APPROVALS_FILE" 2>/dev/null || echo 0)
[ "$SESSION_ENTRIES" -lt 2 ] && exit 0

# Extract computed_rules for this session, count duplicates (2+ = actionable)
# Use grep + jq for speed: grep narrows to session, jq extracts field
ACTIONABLE_RULES=$(grep "\"session_id\":\"$SESSION_ID\"" "$APPROVALS_FILE" 2>/dev/null \
    | jq -r '.computed_rule' 2>/dev/null \
    | sort | uniq -c | sort -rn \
    | awk '$1 >= 2 {print $2}')

[ -z "$ACTIONABLE_RULES" ] && exit 0

# Filter out already-surfaced rules
touch "$SURFACED_FILE" 2>/dev/null
NEW_RULES=""
while IFS= read -r rule; do
    if ! grep -qxF "$rule" "$SURFACED_FILE" 2>/dev/null; then
        NEW_RULES="${NEW_RULES:+$NEW_RULES, }\"$rule\""
        echo "$rule" >> "$SURFACED_FILE" 2>/dev/null
    fi
done <<< "$ACTIONABLE_RULES"

[ -z "$NEW_RULES" ] && exit 0

# Surface the recommendation
MSG="Repeated permission approvals detected. Add these rules to the allow list in ~/stuart/.claude/settings.json to auto-approve them: $NEW_RULES"

jq -n --arg msg "$MSG" '{
  hookSpecificOutput: {
    hookEventName: "PostToolUse",
    additionalContext: $msg
  }
}'
exit 0
