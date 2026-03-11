#!/bin/bash
# Shared helper functions for doreen hooks
# Source this at the top of every hook: source "$(dirname "$0")/lib.sh"

LOG_BASE="/home/cvk/stuart/doreen/workspace/logs"
_LIB_INTENTIONAL_EXIT=0

log_dir() {
    mkdir -p "$LOG_BASE" 2>/dev/null
    echo "$LOG_BASE"
}

# Call after sourcing lib.sh: install_error_trap "hook-name.sh"
# Logs unexpected errors to logs/hook-errors.jsonl
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

# Emit a warning: tool proceeds, but message shown to agent as additionalContext.
# Usage: emit_warning "message" [event_name]
emit_warning() {
    local msg="$1"
    local event="${2:-PreToolUse}"
    _LIB_INTENTIONAL_EXIT=1
    jq -n --arg msg "WARNING: $msg" --arg event "$event" \
        '{"hookSpecificOutput":{"hookEventName":$event,"additionalContext":$msg}}'
    exit 0
}

# Emit a block: tool call is rejected, message shown to agent as error.
# Usage: emit_block "message"
emit_block() {
    _LIB_INTENTIONAL_EXIT=1
    echo "$1" >&2
    exit 2
}

# Log a hook event to JSONL for debugging/auditing.
# Usage: log_event "hook-name" "action" "details"
log_event() {
    local hook="$1" action="$2" details="${3:-}"
    mkdir -p "$LOG_BASE" 2>/dev/null
    jq -n -c \
        --arg ts "$(date -Iseconds)" \
        --arg hook "$hook" \
        --arg action "$action" \
        --arg details "$details" \
        '{ts:$ts, hook:$hook, action:$action, details:$details}' \
        >> "$LOG_BASE/hook-events.jsonl" 2>/dev/null
}

# State directory for hooks that need to track state across calls.
# Usage: STATE_DIR=$(state_dir)
state_dir() {
    local dir="$LOG_BASE/state"
    mkdir -p "$dir" 2>/dev/null
    echo "$dir"
}
