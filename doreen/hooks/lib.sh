#!/bin/bash
# Shared helper functions for doreen hooks.
# Source at the top of every hook: source "$(dirname "$0")/lib.sh"

STUART_ROOT="/home/cvk/stuart"
LOG_BASE="$STUART_ROOT/doreen/workspace/logs"
_LIB_INTENTIONAL_EXIT=0

log_dir() {
    mkdir -p "$LOG_BASE" 2>/dev/null
    echo "$LOG_BASE"
}

# Install error trap for unexpected failures.
# Usage: install_error_trap "hook-name.sh"
# Logs to doreen/workspace/logs/hook-errors.jsonl
install_error_trap() {
    local hook_name="${1:-unknown}"
    _HOOK_NAME="$hook_name"
    trap '_hook_error_handler $? ${LINENO} "${BASH_SOURCE[0]}"' ERR
}

_hook_error_handler() {
    local exit_code="$1" line="$2" source="$3"
    [ "$_LIB_INTENTIONAL_EXIT" -eq 1 ] && return 0
    mkdir -p "$LOG_BASE" 2>/dev/null
    jq -n -c \
        --arg ts "$(date -Iseconds)" \
        --arg hook "${_HOOK_NAME:-unknown}" \
        --argjson line "$line" \
        --argjson exit_code "$exit_code" \
        --arg source "$source" \
        '{ts:$ts, hook:$hook, line:$line, exit_code:$exit_code, source:$source}' \
        >> "$LOG_BASE/hook-errors.jsonl" 2>/dev/null
}

# Emit additional context back to the agent (informational, non-blocking).
# Usage: emit_context "message" "EventName"
emit_context() {
    local msg="$1"
    local event="${2:-SessionStart}"
    _LIB_INTENTIONAL_EXIT=1
    jq -n --arg msg "$msg" --arg event "$event" \
        '{"hookSpecificOutput":{"hookEventName":$event,"additionalContext":$msg}}'
    exit 0
}

# Emit a warning back to the agent (non-blocking).
emit_warning() {
    local msg="$1"
    local event="${2:-PreToolUse}"
    _LIB_INTENTIONAL_EXIT=1
    jq -n --arg msg "WARNING: $msg" --arg event "$event" \
        '{"hookSpecificOutput":{"hookEventName":$event,"additionalContext":$msg}}'
    exit 0
}

# Block a tool call with a message on stderr.
emit_block() {
    _LIB_INTENTIONAL_EXIT=1
    echo "$1" >&2
    exit 2
}

# Detect if current session is root agent (interactive TTY) or subagent.
# Returns 0 (true) if root agent, 1 (false) if subagent.
is_root_agent() {
    local tty_nr
    tty_nr=$(awk '{print $7}' /proc/$PPID/stat 2>/dev/null || echo 0)
    [ "$tty_nr" -ne 0 ] 2>/dev/null
}
