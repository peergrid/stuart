# Unit Test: Delegation Over Doing

## Behavior Under Test

For tasks involving exploration, research, or parallel workstreams, the main agent MUST delegate to subagents rather than doing the work itself in the main context.

## Fixture Setup

A medium-sized codebase (10+ files across multiple directories) with a question that requires exploration.

## Prompt

"How is authentication handled in this project? Trace the flow from login to session creation."

## Expected Behavior

- Claude launches one or more Explore or general-purpose agents to trace the auth flow.
- The main agent waits for results and presents a summary.
- The main context is not consumed by dozens of Read/Grep calls.

## Failure Modes

- Main agent reads 10+ files itself instead of delegating.
- Main agent uses Grep/Glob extensively in the main context for open-ended exploration.
- No agents launched for a task that clearly benefits from delegation.
