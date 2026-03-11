# Hook: Polling Patterns

## Behavior to Intercept

Claude using `sleep` followed by status checks, or actively polling for subagent completion instead of waiting for the completion notification.

## Trigger

`PreToolUse` on `Bash` tool. Inspect the command for:
- `sleep` followed by any check command.
- Repeated identical commands within a short window (same command executed 2+ times in recent history).

Also `PreToolUse` on `Agent` tool or `TaskGet` tool when:
- A recently-launched background agent has not yet returned.
- The call is checking status rather than resuming with new work.

## Action

Block the tool call. Return a message:
- For sleep+check: "Do not poll. Use `run_in_background` and wait for the completion notification."
- For agent polling: "The agent has not completed yet. You will be notified when it finishes. Continue with other work or wait."

## Exceptions

- Polling an external process that does not provide completion notifications (e.g., `gh run view` to check CI status). Even here, prefer a single check over a sleep loop.
