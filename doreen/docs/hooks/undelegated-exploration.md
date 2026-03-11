# Hook: Undelegated Exploration

## Behavior to Intercept

The main agent performing broad codebase exploration (many Read/Grep/Glob calls) in the main context instead of delegating to an Explore or general-purpose agent.

## Trigger

`PostToolUse` — track the main agent's Read/Grep/Glob call count within a conversation window. If the count exceeds a threshold (e.g., 8+ exploratory tool calls without an Agent launch), fire.

Heuristic for "exploratory": sequential Read calls on different files, Grep calls with broad patterns, Glob calls scanning for file types. Contrast with "directed": reading a specific known file before editing it.

## Action

Warn. Insert a message:
"You are exploring the codebase extensively in the main context. Consider delegating this to an Explore agent to preserve context."

## Exceptions

- Reading files that are about to be edited (directed, not exploratory).
- The operator explicitly asked to "look at" or "show me" specific files.
- Working within a small scope (2-3 files).
