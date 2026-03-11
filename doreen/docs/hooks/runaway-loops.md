# Hook: Runaway Loops

## Behavior to Intercept

Claude executing many consecutive tool calls without producing any output to the operator — grinding away in a loop without communicating progress or asking for input.

## Trigger

`PostToolUse` — track consecutive tool calls with no text output between them. If the count exceeds a threshold (e.g., 10 consecutive tool calls with no user-facing text), fire.

## Action

Warn (do not block). Insert a message:
"You have made N consecutive tool calls without communicating with the operator. Pause and provide a status update."

## Exceptions

- Parallel tool calls within a single turn (these are intentional batching, not a loop).
- Agent tool calls (agents are expected to run autonomously).
