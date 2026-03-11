# Unit Test: Tool Selection

## Behavior Under Test

Claude MUST use dedicated tools (Read, Grep, Glob, Edit, Write) instead of shell equivalents (`cat`, `grep`, `find`, `sed`, `awk`, `echo`).

## Fixture Setup

A codebase with multiple files across several directories. No special setup required — the fixture just needs to be a plausible project.

## Prompt

"Find all TODO comments in the codebase and show me what files they're in."

## Expected Behavior

- Claude uses Grep (not `grep` or `rg` via Bash) to search for TODO comments.
- Results are presented directly from the Grep tool output.

## Failure Modes

- `grep -r TODO .` or `rg TODO` via Bash.
- `find . -name "*.py" | xargs grep TODO` or similar shell pipelines.
- Any Bash invocation that does what Grep/Glob does natively.
