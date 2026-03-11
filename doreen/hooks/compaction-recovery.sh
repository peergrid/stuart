#!/bin/bash
# SessionStart(compact) hook: re-inject critical context after compaction.
# Compaction means we're context-constrained -- keep this minimal.
# Wired to: SessionStart (source=compact)
set -euo pipefail
source "$(dirname "$0")/lib.sh"
install_error_trap "compaction-recovery.sh"

INPUT=$(cat)
SOURCE=$(echo "$INPUT" | jq -r '.source // empty')

# Only fire on compaction events
if [ "$SOURCE" != "compact" ]; then
    exit 0
fi

# Build minimal recovery context
RECOVERY="## Post-Compaction Recovery

You just experienced context compaction. Key reminders:

**Project**: ~/stuart -- doreen/ (testing/hooks), anamnesis/ (self-knowledge).
**Config**: Project-level only (.claude/settings.json, CLAUDE.md). Never user-level.
**Docs = desired state, code = actual state.** No status tracking in docs.
**Plans are ephemeral** -- workspace/plans/. Re-read yours if you had one.
**Never cd.** Use absolute paths from ~/stuart.
**Agent reports go to files**, not inline output."

jq -n --arg ctx "$RECOVERY" '{
  hookSpecificOutput: {
    hookEventName: "SessionStart",
    additionalContext: $ctx
  }
}'
exit 0
