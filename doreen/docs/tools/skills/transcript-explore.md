# Skill: Transcript Explore

## Purpose

Discover sessions, get overviews, and navigate through transcript content. This is the starting point for any transcript analysis task.

## When to Use

- "Show me recent sessions"
- "What happened in the last session?"
- "Find the session where I worked on X"
- "Walk me through what Claude did"
- "Show me what happened around turn N"
- "What agents were launched?"
- Any exploratory question about transcript content

## Workflow: Session Discovery

Start by finding what's available.

```bash
# List recent sessions for a project
dq sessions --project stuart --limit 10

# Find sessions containing specific work
dq find --project stuart --in-messages "parser refactor"

# Get overview of the latest session
dq stats --latest
```

## Workflow: Session Overview

Get the shape of a session before diving in.

```bash
# Quick stats
dq stats --project stuart --session 2ac9

# See all agents that were launched
dq agents --project stuart --session 2ac9

# Sample every 10th turn for a quick skim
dq walk --project stuart --session 2ac9 --limit 20

# See only operator messages (what the human asked for)
dq messages --project stuart --session 2ac9 --external-only
```

## Workflow: Navigating to Points of Interest

Find something specific, then look around it.

```bash
# Find and show the first error with context
dq show --project stuart --session 2ac9 --first error --context 5

# Find where a compaction happened
dq show --project stuart --session 2ac9 --first compaction --context 10

# Jump to a specific turn
dq show --project stuart --session 2ac9 --turn 42 --context 3

# Find where Claude first used a specific tool
dq show --project stuart --session 2ac9 --first "tool:Agent" --context 5
```

## Workflow: Walking Through a Session

Step through turns to understand the flow.

```bash
# Walk from the beginning, tool calls only
dq walk --project stuart --session 2ac9 --tools-only --limit 30

# Walk backward from the end to see how it finished
dq walk --project stuart --session 2ac9 --reverse --limit 10

# Walk from a specific point forward
dq walk --project stuart --session 2ac9 --from 50 --limit 20

# Walk only external messages (operator + Claude responses, no tool noise)
dq walk --project stuart --session 2ac9 --external-only
```

## Workflow: Investigating Agents

Trace what subagents did.

```bash
# List all agents with their outcomes
dq agent-trace --project stuart --latest

# Investigate a specific agent
dq agent-trace --project stuart --latest --agent "observer"

# Walk through an agent's transcript directly
dq walk /path/to/subagent.jsonl
```

## Output Modes

- Default: human-readable formatted output for terminal review
- `--json`: structured JSON when feeding results to another tool or grader
- `--jsonl`: streaming JSONL for piping through jq or other processors

## Tips

- Always start with `dq sessions` or `dq stats` to orient
- Use `--limit` liberally to avoid overwhelming output
- The `--external-only` flag is essential for seeing the "human conversation" without tool noise
- Combine `dq show --first` with `dq walk --from` to jump to interesting points then step through
- For deep dives, use `dq raw` piped to `jq` for custom field extraction
