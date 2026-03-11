# Skill: Transcript Grade

## Purpose

Run grader-specific analyses on transcripts, producing the structured data each doreen grader needs. This skill covers all grader workflows: tool-use-audit, transcript-critique, token-cost-metrics, and work-product-critique.

## When to Use

- "Grade the last session"
- "How many tokens did that session use?"
- "Show me the token consumption curve"
- "Extract data for the transcript critique"
- "What was the cost efficiency?"
- "Were there any compaction events?"
- "Prepare grading data for session X"
- Any grader-oriented analysis

## Workflow: Tool Use Audit Grading

Automated grader. Produces a score and violation list.

```bash
# Run the full audit for grading
dq audit --latest --json > /tmp/audit-result.json

# The JSON output contains:
# - score (1.0 - violation_rate)
# - violations array with turn references
# - check results for each category
# - parallelism_ratio
# - delegation_ratio
```

The grader reads the JSON and produces the final grade. No additional processing.

## Workflow: Token Cost Metrics Grading

Automated grader. Produces token metrics.

```bash
# Full token timeline
dq tokens --latest --json > /tmp/token-metrics.json

# Summary metrics only
dq tokens --latest --summary --json

# Key numbers the grader needs:
# - total_input_tokens, total_output_tokens
# - max_single_turn_context
# - compaction_count
# - context_utilization_curve (array of per-turn context sizes)
# - tokens_per_turn_average

# Compaction details
dq compactions --latest --json > /tmp/compactions.json

# Agent token overhead
dq agent-trace --latest --json > /tmp/agent-traces.json
```

To compute tokens-per-line-of-diff (a key efficiency metric), combine dq output with the git diff:

```bash
# Token data
TOTAL_TOKENS=$(dq tokens --latest --summary --json | jq '.total_input_tokens + .total_output_tokens')

# Lines of meaningful diff (from the work product)
DIFF_LINES=$(git diff --stat HEAD~1 | tail -1)

# The grader combines these
```

## Workflow: Transcript Critique Grading

LLM-graded. Extracts structured data for a separate Claude instance to critique.

```bash
# Extract critique data
dq critique-data --latest --json > /tmp/critique-data.json

# The JSON output contains:
# - original_prompt: first external user message
# - decision_points: turns where Claude chose an approach
# - error_recovery: error + subsequent recovery turns
# - completion_claims: turns where Claude said "done" with post-context
# - communication_signals: response lengths, verbosity
# - context_management: compactions, repeated reads

# For manual critique preparation, you can also gather:
# Operator messages only (to see the full conversation from human perspective)
dq messages --latest --external-only --no-truncate --json

# Error sequences with context
dq errors --latest --with-context 5 --json

# All places Claude declared completion
dq messages --latest --role assistant --contains "\\b(done|complete|finished)\\b" --json
```

Feed the critique-data JSON to the LLM grader as context. The grader evaluates:
- Decision quality
- Communication quality
- Error recovery
- Context management
- Honesty (completion claims vs reality)

## Workflow: Work Product Critique Grading

LLM-graded. Extracts the file modifications Claude made.

```bash
# Get all file modifications
dq tools --latest --tool Edit --json > /tmp/edits.json
dq tools --latest --tool Write --json > /tmp/writes.json

# Get the original prompt
dq messages --latest --external-only --limit 1 --json > /tmp/prompt.json

# The grader combines these with the actual file contents and git diff
# to evaluate correctness, scope discipline, convention adherence, etc.
```

## Workflow: Regression Comparison

Not directly a dq workflow — the regression grader compares stored results. But dq provides the baseline data:

```bash
# Generate all grading data for storage
dq audit --latest --json > results/audit.json
dq tokens --latest --json > results/tokens.json
dq stats --latest --json > results/stats.json

# The regression grader diffs these against historical results
```

## Workflow: Full Session Grading

Run all grader analyses for a session in one pass.

```bash
SESSION="--project stuart --session 2ac9"

# Collect all grading inputs
dq audit $SESSION --json > /tmp/grade-audit.json
dq tokens $SESSION --json > /tmp/grade-tokens.json
dq compactions $SESSION --json > /tmp/grade-compactions.json
dq errors $SESSION --json > /tmp/grade-errors.json
dq critique-data $SESSION --json > /tmp/grade-critique-data.json
dq agent-trace $SESSION --json > /tmp/grade-agent-traces.json
dq stats $SESSION --json > /tmp/grade-stats.json
dq tools $SESSION --tool Edit --json > /tmp/grade-edits.json
dq tools $SESSION --tool Write --json > /tmp/grade-writes.json
dq messages $SESSION --external-only --limit 1 --json > /tmp/grade-prompt.json
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
