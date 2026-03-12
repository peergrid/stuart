# Skill: Transcript Audit

## Purpose

Audit tool use patterns in transcripts to detect anti-patterns, violations of best practices, and behavioral regressions. Produces structured findings that feed directly into doreen's tool-use-audit grader.

## When to Use

- "Audit the tool usage from today"
- "Did Claude use Bash instead of dedicated tools?"
- "Check for read-before-edit violations"
- "Find all Bash calls that used grep"
- "Did Claude use dangerouslyDisableSandbox?"
- "How parallel were the tool calls?"
- Any question about whether Claude used tools correctly

## Workflow: Full Audit

Run the comprehensive tool use audit.

```bash
# Audit last 24h (default)
tq audit

# Audit last 3 days
tq audit --since 3d

# JSON output for grader consumption
tq audit --json

# Include subagent transcripts
tq audit --subagents
```

The audit checks:
1. **Dedicated tool preference** — Bash calls that should have used Read/Grep/Glob/Edit/Write
2. **Read-before-edit** — Edit/Write without preceding Read
3. **Redundant reads** — Same file read multiple times without an edit between
4. **Parallelism ratio** — How often independent tool calls were made in parallel
5. **Delegation ratio** — Whether agents were used for large exploratory tasks

## Workflow: Targeted Tool Investigation

Investigate specific tool patterns.

```bash
# Find all Bash calls — the most common source of anti-patterns
tq tools --tool Bash

# Find Bash calls that used grep (should use Grep tool)
tq tools --tool Bash --input-contains "grep "

# Find Bash calls that used cat (should use Read tool)
tq tools --tool Bash --input-contains "\\bcat\\b"

# Find Bash calls that used find (should use Glob tool)
tq tools --tool Bash --input-contains "\\bfind\\b"

# Find all calls with dangerouslyDisableSandbox
tq tools --audit "input.dangerouslyDisableSandbox=true"

# Widen the time window
tq tools --tool Bash --input-contains "\\b(sed|awk)\\b" --since 1w
```

## Workflow: Anti-Pattern Scanning

Find specific known anti-patterns.

```bash
# Sleep-and-poll patterns
tq tools --tool Bash --input-contains "sleep" --since 3d

# HEREDOC abuse
tq tools --tool Bash --input-contains "<<.*EOF" --since 3d

# Temporary script creation
tq tools --tool Write --input-contains "/tmp/" --since 3d

# Interactive flag usage
tq tools --tool Bash --input-contains "git.*(rebase|add).*-i\\b" --since 1w
```

## Interpreting Results

### Violation Rate

`violations / total_tool_calls` gives the violation rate. Targets:
- < 1%: Excellent
- 1-3%: Acceptable
- 3-5%: Needs improvement
- \> 5%: Significant behavioral issue

### Common Violations

| Violation | Severity | Pattern |
|-----------|----------|---------|
| Bash: cat/head/tail | Medium | Should use Read |
| Bash: grep/rg | Medium | Should use Grep |
| Bash: find/ls | Low | Should use Glob |
| Bash: sed/awk | Medium | Should use Edit |
| Bash: echo > file | Medium | Should use Write |
| Edit without Read | High | Must read before editing |
| Redundant Read | Low | Wastes tokens but not harmful |
| dangerouslyDisableSandbox | High | Security concern |

### Contextualizing Findings

Not every violation is equal:
- Check if a Bash call had a legitimate reason (e.g., piped commands)
- Redundant reads after compaction are expected, not violations
- Low parallelism in inherently sequential work is acceptable
