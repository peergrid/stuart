# Grader: Work Product Critique

## What It Evaluates

Quality of the code Claude produced — correctness, readability, adherence to conventions, appropriate complexity.

## Type

LLM-graded — a separate Claude instance reviews the work product.

## Inputs

- The diff (what Claude changed).
- The full files after changes (for context).
- The original prompt/task.
- Project conventions (CLAUDE.md, linter configs, style guides).

## Critique Dimensions

- **Correctness**: Does the code do what was asked? Are there obvious bugs or logic errors?
- **Scope discipline**: Did Claude change only what was requested? Any drive-by refactoring, unnecessary comments, or unrequested features?
- **Convention adherence**: Does the code match the project's style, naming, patterns?
- **Complexity**: Is the solution appropriately simple? Any premature abstractions, unnecessary indirection, or over-engineering?
- **Security**: Any vulnerabilities introduced (injection, XSS, exposed secrets, etc.)?

## Output

A structured critique with a rating per dimension (pass/concern/fail) and specific findings with code references.
