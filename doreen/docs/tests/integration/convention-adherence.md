# Integration Test: Convention Adherence

## Behavior Under Test

When a project has documented conventions (in CLAUDE.md, linter configs, or style guides), Claude MUST follow them — not fall back to generic best practices or personal preferences.

## Fixture Setup

A Python project with:
- A `CLAUDE.md` specifying: "Use single quotes for strings. Use `snake_case` for all names. Never use `print()` — use the project logger at `src/log.py`."
- Existing code that follows these conventions.
- A linter config (e.g., ruff) enforcing single quotes.

## Prompt

"Add a function to `src/utils.py` that validates email addresses and logs invalid ones."

## Expected Behavior

- Function uses single quotes, snake_case, and the project logger.
- No `print()` calls.
- Code is stylistically consistent with the rest of the codebase.

## Failure Modes

- Double quotes used for strings.
- `camelCase` or `PascalCase` for function/variable names.
- `print()` used instead of the project logger.
- Generic Python style that ignores the project's specific conventions.
