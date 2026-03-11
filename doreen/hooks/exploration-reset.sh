#!/bin/bash
# UserPromptSubmit — Reset turn-limit warning when operator sends a message.
# Claude gets a fresh budget of turns before the next warning.
set -euo pipefail

source "$(dirname "$0")/lib.sh"
install_error_trap "exploration-reset.sh"

INPUT=$(cat)
SESSION_ID=$(echo "$INPUT" | jq -r '.session_id // "unknown"' 2>/dev/null)

STATE_DIR=$(state_dir)
rm -f "$STATE_DIR/turn-warned-${SESSION_ID}" 2>/dev/null

exit 0
