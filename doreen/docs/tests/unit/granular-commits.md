# Unit Test: Granular Commits on a Branch

## Behavior Under Test

When asked to commit changes, Claude MUST create a new branch and produce granular, well-scoped commits — not a single monolithic commit of everything.

## Fixture Setup

A git repo on `master` with several uncommitted changes spanning multiple concerns:
- A bug fix in `src/auth.py`
- A new utility function in `src/utils.py`
- A formatting cleanup in `src/config.py`
- An updated test in `tests/test_auth.py` related to the bug fix

## Prompt

"Please commit these changes."

## Expected Behavior

- Claude creates a new branch off `master` (not committing directly to `master`).
- Changes are split into logical commits (e.g., the bug fix + its test together, the utility function separately, the formatting cleanup separately).
- Each commit message describes the *why*, not just the *what*.

## Failure Modes

- All changes lumped into a single commit.
- Commits made directly on `master`.
- Unrelated changes grouped together (e.g., formatting cleanup bundled with the bug fix).
- Commit messages that just list files changed.
