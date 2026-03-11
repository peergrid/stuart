# Grader: Token Cost Metrics

## What It Evaluates

Resource efficiency — how many tokens were consumed relative to the work accomplished.

## Type

Automated — transcript metrics.

## Inputs

A session transcript with token counts per turn, plus the final work product (diff).

## Metrics

- **Total tokens**: Input + output for the session.
- **Tokens per line of meaningful diff**: Total tokens divided by lines of non-whitespace, non-comment code changed. A proxy for efficiency.
- **Context utilization curve**: How full the context window was over time. Rapid growth suggests wasteful reads; a compaction event suggests poor context management.
- **Compaction count**: How many times the context was compacted. Zero is ideal for small tasks; more than one for a simple task is a red flag.
- **Agent token overhead**: Tokens consumed by agents vs. the main context. High agent overhead with low output suggests poorly scoped delegations.

## Output

A metrics summary with the raw numbers and comparison against historical baselines for similar task types.
