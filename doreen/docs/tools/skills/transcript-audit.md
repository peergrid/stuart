# Skill: Transcript Audit

## Purpose

Audit tool use patterns in transcripts to detect anti-patterns, violations of best practices, and behavioral regressions. Produces structured findings that feed directly into doreen's tool-use-audit grader.

## When to Use

- "Audit the tool usage in the last session"
- "Did Claude use Bash instead of dedicated tools?"
- "Check for read-before-edit violations"
- "Find all Bash calls that used grep"
- "Did Claude use dangerouslyDisableSandbox?"
- "How parallel were the tool calls?"
- Any question about whether Claude used tools correctly

## Workflow: Full Audit

Run the comprehensive tool use audit.

```bash
# Human-readable audit report
dq audit --latest

# JSON output for grader consumption
dq audit --latest --json

# Audit a specific session
dq audit --project stuart --session 2ac9

# Audit with subagent transcripts included
dq audit --project stuart --session 2ac9 --subagents
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
dq tools --latest --tool Bash

# Find Bash calls that used grep (should use Grep tool)
dq tools --latest --tool Bash --input-contains "grep "

# Find Bash calls that used cat (should use Read tool)
dq tools --latest --tool Bash --input-contains "\\bcat\\b"

# Find Bash calls that used find (should use Glob tool)
dq tools --latest --tool Bash --input-contains "\\bfind\\b"

# Find Bash calls that used sed/awk (should use Edit tool)
dq tools --latest --tool Bash --input-contains "\\b(sed|awk)\\b"

# Find all calls with dangerouslyDisableSandbox
dq tools --latest --audit "input.dangerouslyDisableSandbox=true"

# Find all Bash calls that used echo to write files (should use Write)
dq tools --latest --tool Bash --input-contains "echo.*>\\s"
```

## Workflow: Read-Edit Chain Analysis

Verify proper read-before-edit patterns.

```bash
# Show all Edit and Write calls with context to verify preceding Reads
dq tools --latest --tool Edit --with-results
dq tools --latest --tool Write --with-results

# The audit command does this automatically, but for manual inspection:
# Show Edit calls, then check what came before each one
dq show --latest --first "tool:Edit" --context 3
```

## Workflow: Parallelism Check

Assess whether Claude used parallel tool calls effectively.

```bash
# The stats output shows tool call counts per turn
dq stats --latest

# The audit command calculates parallelism ratio
dq audit --latest --json | jq '.parallelism_ratio'

# To manually inspect: walk through and look for sequential single-tool turns
dq walk --latest --tools-only
```

## Workflow: Anti-Pattern Scanning

Find specific known anti-patterns across sessions.

```bash
# Sleep-and-poll patterns
dq tools --project stuart --tool Bash --input-contains "sleep"

# HEREDOC abuse
dq tools --project stuart --tool Bash --input-contains "<<.*EOF"

# Temporary script creation
dq tools --project stuart --tool Write --input-contains "/tmp/"

# Interactive flag usage (git rebase -i, git add -i)
dq tools --project stuart --tool Bash --input-contains "git.*(rebase|add).*-i\\b"
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

Not every violation is equal. When reviewing audit results:
- Check if a Bash call had a legitimate reason (e.g., piped commands that cannot be expressed with dedicated tools)
- Redundant reads after compaction are expected, not violations
- Low parallelism in inherently sequential work is acceptable
