#!/bin/bash
# PreToolUse:Bash — BLOCK commands using syntactic patterns that cannot be
# matched by permission glob rules, and suggest the approvable form.
# This is about making commands auto-approvable, not about safety.
# Tier: BLOCK (exit 2 with fix suggestion on stderr)
set -euo pipefail
source "$(dirname "$0")/lib.sh"
install_error_trap "command-rewriter.sh"

INPUT=$(cat)
CMD=$(echo "$INPUT" | jq -r '.tool_input.command // ""' 2>/dev/null)

# Skip empty commands
[ -z "$CMD" ] && exit 0

# Check: command starts with ( but is NOT an arithmetic $((...)) embedded in
# a larger command. A leading ( is a subshell — permission globs cannot match it.
#
# Pattern: command begins with optional whitespace then (
# Exclude: $((  which is arithmetic expansion (only appears mid-command)
if echo "$CMD" | grep -qE '^\s*\('; then
    # Extract what's inside the subshell for the suggestion
    INNER=$(echo "$CMD" | sed -E 's/^\s*\(\s*//' | sed -E 's/\s*\)\s*$//')
    emit_block "BLOCKED: Commands starting with ( are subshells that cannot be matched by permission allow rules. The permission system matches the first word of the command, and ( is not a valid match target.

Rewrite as:
  bash -c '$INNER'

This makes the command matchable by Bash(bash *) in the allow list."
fi

exit 0
