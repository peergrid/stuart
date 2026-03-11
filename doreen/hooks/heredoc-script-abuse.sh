#!/bin/bash
# PreToolUse:Bash — Block HEREDOC file creation, temporary scripts, and oversized one-liners
# Detects: cat <<EOF > file, bash -c multiline, 200+ char file manipulation, /tmp scripts
set -euo pipefail

source "$(dirname "$0")/lib.sh"
install_error_trap "heredoc-script-abuse.sh"

INPUT=$(cat)
TOOL_NAME=$(echo "$INPUT" | jq -r '.tool_name // "unknown"' 2>/dev/null)

# Only check Bash tool
if [ "$TOOL_NAME" != "Bash" ]; then
    exit 0
fi

CMD=$(echo "$INPUT" | jq -r '.tool_input.command // ""' 2>/dev/null)
[ -z "$CMD" ] && exit 0

# --- Exception: git commit with HEREDOC is the established pattern ---
if echo "$CMD" | grep -qE '^\s*git\s.*\bcommit\b'; then
    exit 0
fi

# --- Exception: commands that genuinely need shell features for a simple purpose ---
# Short commands (under 100 chars) with pipes are usually fine
CMD_LEN=${#CMD}
if [ "$CMD_LEN" -lt 100 ] && echo "$CMD" | grep -qE '\|' && ! echo "$CMD" | grep -qE '(cat\s*<<|>>?\s*/|>\s*/)'; then
    exit 0
fi

# --- Detection ---

# cat <<EOF or cat <<'EOF' writing to a file
if echo "$CMD" | grep -qE 'cat\s*<<'; then
    # Is it redirecting to a file?
    if echo "$CMD" | grep -qE 'cat\s*<<.*>\s*\S|>\s*\S.*cat\s*<<'; then
        emit_block "BLOCKED: Use the Write tool instead of HEREDOC redirection (cat <<EOF > file). The Write tool creates files with proper tracking and review."
    fi
    # Even without explicit redirect, cat <<EOF is usually writing somewhere
    emit_block "BLOCKED: Use the Write tool instead of HEREDOC (cat <<EOF). The Write tool is the proper way to create file content."
fi

# bash -c with multi-line or complex content
if echo "$CMD" | grep -qE 'bash\s+-c\s'; then
    # Get the content after bash -c
    SCRIPT_CONTENT=$(echo "$CMD" | sed "s/.*bash\s*-c\s*['\"]//")
    SCRIPT_LEN=${#SCRIPT_CONTENT}
    # Multi-line (contains actual newlines) or long
    if [ "$SCRIPT_LEN" -gt 100 ] || echo "$CMD" | grep -qcP '\n' 2>/dev/null; then
        emit_block "BLOCKED: This 'bash -c' command is too complex for an inline script. Create a proper script with the Write tool, or break this into multiple dedicated tool calls."
    fi
fi

# Commands exceeding 200 chars that involve file manipulation
if [ "$CMD_LEN" -gt 200 ]; then
    if echo "$CMD" | grep -qE '(>\s*\S|>>\s*\S|tee\s|sed\s+-i|install\s+-|cp\s.*\S+\.\S+|mv\s|mkdir.*&&)'; then
        emit_block "BLOCKED: This command is $CMD_LEN characters long and involves file manipulation. This is too complex for an inline bash command. Break it into multiple dedicated tool calls, or use the Write tool for file creation."
    fi
fi

# Creating files in /tmp that look like project scripts
if echo "$CMD" | grep -qE '(>\s*/tmp/\S+\.(sh|py|js|ts|rb|pl|go)|chmod\s+\+x\s+/tmp/|bash\s+/tmp/|python3?\s+/tmp/|node\s+/tmp/)'; then
    emit_block "BLOCKED: Creating and running scripts in /tmp is a sign of script abuse. If you need a script, create it in the project with the Write tool. If you need a temporary operation, use dedicated tools."
fi

exit 0
