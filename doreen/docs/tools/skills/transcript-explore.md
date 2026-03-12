# Skill: Transcript Explore

## Purpose

Discover sessions, get overviews, and navigate through transcript content. This is the starting point for any transcript analysis task.

## When to Use

- "Show me recent sessions"
- "What happened in the last hour?"
- "Find where I worked on X"
- "Walk me through what Claude did"
- "Show me what happened around that error"
- "What agents were launched today?"
- Any exploratory question about transcript content

## Workflow: Quick Overview

Get the shape of recent work. No session IDs needed.

```bash
# Stats for the last 24h (default)
tq stats

# What happened in the last 2 hours?
tq stats --around 2h

# What happened this week?
tq stats --since 1w

# List recent sessions
tq sessions --since 3d
```

## Workflow: Finding Things

Search across all sessions in a time window.

```bash
# Find sessions/agents that mention specific work
tq find "parser refactor" --since 3d

# Find sessions with specific text
tq find --in-messages "fix the bug" --since 1w

# See what the operator asked for recently
tq messages --external-only --around 2h
```

## Workflow: Navigating to Points of Interest

Find something specific, then look around it.

```bash
# First error in the last 2 hours with context
tq show --first error --context 5 --around 2h

# Where did compaction happen?
tq show --first compaction --context 10

# First time Agent was launched today
tq show --first "tool:Agent" --context 5 --around 1d

# Find a specific pattern
tq show --first "pattern:TODO" --context 3
```

## Workflow: Walking Through History

Step through turns to understand the flow.

```bash
# Walk through last hour, tools only
tq walk --around 1h --tools-only

# Walk backward from now, last 20 turns
tq walk --reverse --limit 20

# Walk only external messages (operator + Claude responses, no tool noise)
tq walk --external-only --around 2h

# Walk from a specific timestamp
tq walk --from "2026-03-10T14:00:00"
```

## Workflow: Investigating Agents

Trace what subagents did.

```bash
# All agents launched today
tq agent-trace --since 1d

# Find a specific agent by name
tq agent-trace --agent "observer"

# Agents from the last 3 days
tq agent-trace --since 3d --json
```

## Output Modes

- Default: human-readable formatted output for terminal review
- `--json`: structured JSON when feeding results to another tool or grader
- `--jsonl`: streaming JSONL for piping through jq or other processors

## Tips

- Start with `tq stats` or `tq sessions` to orient
- Batch commands default to 24h — use `--since` to widen or narrow
- Cursor commands use `--around` to position — no time limit on navigation
- You never need a session ID
- `--external-only` is essential for seeing the "human conversation" without tool noise
- Combine `tq show --first` with `tq walk --from` to jump to interesting points then step through
