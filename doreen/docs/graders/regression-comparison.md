# Grader: Regression Comparison

## What It Evaluates

Whether the current test run shows regressions compared to previous runs of the same test.

## Type

Automated — diff against historical results.

## Inputs

- Current test results (all grader outputs for this run).
- Historical results for the same test (stored in `doreen/workspace/results/`).

## Comparisons

- **Score regressions**: Any grader score that worsened compared to the best previous run.
- **New violations**: Tool use violations or failure modes that didn't appear in previous runs.
- **Cost regressions**: Token cost significantly higher than the historical median for this test.
- **Resolved issues**: Previous violations that no longer appear (positive signal).

## Output

A regression report listing:
- New regressions (things that got worse).
- Persistent issues (things that have been bad across multiple runs).
- Improvements (things that got better).
- Trend lines for key metrics over the last N runs.
