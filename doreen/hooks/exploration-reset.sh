#!/bin/bash
# UserPromptSubmit — Reset the exploration counter when the user sends a message.
# This gives the root agent a fresh budget of exploratory calls per conversation turn.
set -euo pipefail

source "$(dirname "$0")/lib.sh"
install_error_trap "exploration-reset.sh"

INPUT=$(cat)
SESSION_ID=$(echo "$INPUT" | jq -r '.session_id // "unknown"' 2>/dev/null)

STATE_DIR=$(state_dir)
echo "0" > "$STATE_DIR/explore-count-${SESSION_ID}" 2>/dev/null

exit 0
