#!/bin/bash
# PreToolUse:Bash — Block shell commands when dedicated tools should be used
# Detects: cat/head/tail (use Read), grep/rg (use Grep), find/ls (use Glob),
#          sed/awk (use Edit), echo>/cat<<EOF> (use Write)
set -euo pipefail

source "$(dirname "$0")/lib.sh"
install_error_trap "shell-tool-misuse.sh"

INPUT=$(cat)
TOOL_NAME=$(echo "$INPUT" | jq -r '.tool_name // "unknown"' 2>/dev/null)

# Only check Bash tool
if [ "$TOOL_NAME" != "Bash" ]; then
    exit 0
fi

CMD=$(echo "$INPUT" | jq -r '.tool_input.command // ""' 2>/dev/null)
[ -z "$CMD" ] && exit 0

# --- Exceptions: always allow these through ---

# Multi-stage pipes (3+ pipe segments) — complex pipelines that no single tool replaces
PIPE_COUNT=$(echo "$CMD" | tr -cd '|' | wc -c)
if [ "$PIPE_COUNT" -ge 2 ]; then
    exit 0
fi

# Extract the first real command (skip variable assignments, env prefixes)
FIRST_CMD=$(echo "$CMD" | sed 's/^[[:space:]]*//' | sed 's/^[A-Z_][A-Z_0-9]*=[^ ]* //' | awk '{print $1}')

# Allow legitimate shell tools that have no dedicated alternative
SHELL_TOOLS="git|npm|npx|make|docker|go|pip|pip3|python|python3|node|cargo|rustc|cmake|gcc|g++|clang|chmod|chown|mkdir|rmdir|cp|mv|rm|ln|touch|tar|zip|unzip|curl|wget|ssh|scp|rsync|kill|killall|pkill|ps|pstree|lsof|pgrep|du|df|mount|umount|uname|whoami|id|env|export|set|unset|which|type|command|realpath|readlink|stat|file|diff|sort|uniq|cut|tr|wc|tee|xargs|nohup|timeout|date|gh|jq|fuser|ss|dmesg|journalctl|apt|apt-get|apt-cache|brew|test|true|false|bash|sh|source|cd|pushd|popd|for|while|if"

if echo "$FIRST_CMD" | grep -qxE "($SHELL_TOOLS)"; then
    exit 0
fi

# --- Detection: check for misused commands ---

# cat/head/tail used to read files -> should use Read
if echo "$CMD" | grep -qE '^\s*(cat|head|tail)\s+[^|]'; then
    # Allow cat with pipes (cat file | something) — but only if genuinely piped
    if echo "$CMD" | grep -qE '^\s*(cat|head|tail)\s+.*\|'; then
        exit 0
    fi
    emit_block "BLOCKED: Use the Read tool instead of '$(echo "$FIRST_CMD")' to read files. The Read tool provides line numbers and handles large files properly."
fi

# grep/rg used for content search -> should use Grep
if echo "$CMD" | grep -qE '^\s*(grep|rg)\s'; then
    # Allow grep in pipes (something | grep pattern)
    if echo "$CMD" | grep -qE '\|\s*(grep|rg)\s'; then
        exit 0
    fi
    emit_block "BLOCKED: Use the Grep tool instead of '$(echo "$FIRST_CMD")'. The Grep tool supports regex, file type filters, and context lines."
fi

# find/ls for file discovery -> should use Glob
if echo "$CMD" | grep -qE '^\s*find\s'; then
    emit_block "BLOCKED: Use the Glob tool instead of 'find'. The Glob tool is optimized for file pattern matching."
fi

# ls used for file discovery (not ls with specific flags for details)
if echo "$CMD" | grep -qE '^\s*ls\s+(-[a-zA-Z]*\s+)*[^-|]' && ! echo "$CMD" | grep -qE '^\s*ls\s+.*(-l|-la|-al|-lh|-alh)'; then
    # Only flag bare 'ls <path>' for directory listing, not 'ls -l' which gives file details
    # Actually, be lenient: ls is commonly used for checking directory contents
    exit 0
fi

# sed/awk for file modification -> should use Edit
if echo "$CMD" | grep -qE '^\s*(sed|awk)\s'; then
    # Allow sed/awk in pipes — they're text processors
    if echo "$CMD" | grep -qE '\|\s*(sed|awk)\s'; then
        exit 0
    fi
    # sed -i or awk modifying files
    if echo "$CMD" | grep -qE '^\s*sed\s+(-[a-zA-Z]*i|--in-place)'; then
        emit_block "BLOCKED: Use the Edit tool instead of 'sed -i' for file modifications. The Edit tool tracks changes precisely and is safer."
    fi
    # Plain sed without pipe is often misuse
    if echo "$CMD" | grep -qE '^\s*sed\s'; then
        emit_block "BLOCKED: Use the Edit tool (for modifications) or Grep tool (for search) instead of 'sed'. These dedicated tools are more reliable."
    fi
    if echo "$CMD" | grep -qE '^\s*awk\s.*>\s'; then
        emit_block "BLOCKED: Use the Edit tool instead of 'awk' for file transformations. The Edit tool provides precise, reviewable changes."
    fi
fi

# echo "..." > file -> should use Write
if echo "$CMD" | grep -qE '^\s*echo\s.*>\s*[^|&]'; then
    emit_block "BLOCKED: Use the Write tool instead of 'echo > file'. The Write tool creates files properly and is reviewable."
fi

# cat <<EOF > file -> should use Write (also caught by heredoc hook, but catch here too)
if echo "$CMD" | grep -qE '^\s*cat\s*<<'; then
    emit_block "BLOCKED: Use the Write tool instead of 'cat <<EOF'. The Write tool creates files properly and is reviewable."
fi

exit 0
