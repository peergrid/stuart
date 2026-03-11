# Grader: Tool Use Audit

## What It Evaluates

Whether Claude selected the correct tools in the correct order, without redundant or inappropriate calls.

## Type

Automated — transcript analysis.

## Inputs

A session transcript (sequence of tool calls with arguments and results).

## Checks

- **Dedicated tool preference**: No `cat`/`grep`/`find`/`sed` via Bash when Read/Grep/Glob/Edit exist. Count violations.
- **Read-before-edit**: Every Edit or Write to an existing file was preceded by a Read of that file. Count violations.
- **No redundant reads**: Same file not read multiple times without an edit in between (suggests lost context or wasted tokens).
- **Parallel where possible**: Independent tool calls made in parallel, not sequentially. Measure parallelism ratio.
- **Agent delegation**: For sessions with 10+ exploratory tool calls, was an agent launched? Measure exploration-to-delegation ratio.

## Output

A score (violations / total tool calls) and a list of specific violations with line references into the transcript.
