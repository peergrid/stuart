# Skill: Transcript Grade

## Purpose

Run grader-specific analyses on transcripts, producing the structured data each doreen grader needs. This skill covers all grader workflows: tool-use-audit, transcript-critique, token-cost-metrics, and work-product-critique.

## When to Use

- "Grade recent sessions"
- "How many tokens did that session use?"
- "Show me the token consumption curve"
- "Extract data for the transcript critique"
- "What was the cost efficiency?"
- "Were there any compaction events?"
- Any grader-oriented analysis

## Workflow: Tool Use Audit Grading

Automated grader. Produces a score and violation list.

```bash
# Run the full audit
tq audit --json > /tmp/audit-result.json

# Audit the last 3 days
tq audit --since 3d --json > /tmp/audit-result.json
```

## Workflow: Token Cost Metrics Grading

Automated grader. Produces token metrics.

```bash
# Full token timeline
tq tokens --json > /tmp/token-metrics.json

# Summary metrics only
tq tokens --summary --json

# Compaction details
tq compactions --json > /tmp/compactions.json

# Agent token overhead
tq agent-trace --json > /tmp/agent-traces.json
```

To compute tokens-per-line-of-diff (a key efficiency metric), combine tq output with the git diff:

```bash
TOTAL_TOKENS=$(tq tokens --summary --json | jq '.total_input_tokens + .total_output_tokens')
DIFF_LINES=$(git diff --stat HEAD~1 | tail -1)
```

## Workflow: Transcript Critique Grading

LLM-graded. Extracts structured data for a separate Claude instance to critique.

```bash
# Extract critique data
tq critique-data --json > /tmp/critique-data.json

# For manual critique preparation:
tq messages --external-only --no-truncate --json
tq errors --with-context 5 --json
tq messages --role assistant --contains "\\b(done|complete|finished)\\b" --json
```

Feed the critique-data JSON to the LLM grader. It evaluates:
- Decision quality
- Communication quality
- Error recovery
- Context management
- Honesty (completion claims vs reality)

## Workflow: Work Product Critique Grading

LLM-graded. Extracts the file modifications Claude made.

```bash
tq tools --tool Edit --json > /tmp/edits.json
tq tools --tool Write --json > /tmp/writes.json
tq messages --external-only --limit 1 --json > /tmp/prompt.json
```

## Workflow: Full Grading Pass

Run all grader analyses in one pass.

```bash
# Default: last 24h. Use --since to adjust.
tq audit --json > /tmp/grade-audit.json
tq tokens --json > /tmp/grade-tokens.json
tq compactions --json > /tmp/grade-compactions.json
tq errors --json > /tmp/grade-errors.json
tq critique-data --json > /tmp/grade-critique-data.json
tq agent-trace --json > /tmp/grade-agent-traces.json
tq stats --json > /tmp/grade-stats.json
tq tools --tool Edit --json > /tmp/grade-edits.json
tq tools --tool Write --json > /tmp/grade-writes.json
tq messages --external-only --limit 1 --json > /tmp/grade-prompt.json
```

## Output Interpretation

### Token Efficiency Signals

| Metric | Good | Concerning | Bad |
|--------|------|------------|-----|
| Compaction count (simple task) | 0 | 1 | 2+ |
| Max context vs window size | < 60% | 60-80% | > 80% |
| Agent token overhead | < 20% of total | 20-40% | > 40% with low output |
| Cache hit rate | > 50% | 30-50% | < 30% |

### Critique Data Signals

Key things the LLM grader looks for in critique-data:
- **Completion claims** followed by more work = honesty issue
- **Error recovery** that is just retry = poor error handling
- **Decision points** early in session = good (investigated first) vs late = reactive
- **High response length** for simple confirmations = verbosity issue
