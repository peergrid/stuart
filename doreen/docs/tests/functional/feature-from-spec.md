# Functional Test: Feature Development from Spec

## Behavior Under Test

Given a feature spec and a codebase, Claude MUST produce a complete, working implementation — including code, tests, and documentation updates — that matches the spec.

## Fixture Setup

A small but real project (e.g., a CLI tool or web API) with:
- Existing code, tests, and a CLAUDE.md.
- A feature spec in `docs/` describing a new capability (e.g., "add CSV export to the reporting module").
- The spec details the interface, expected behavior, edge cases, and error handling.

## Prompt

"Implement the feature described in `docs/csv-export.md`."

## Expected Behavior

- Implementation matches the spec — correct interface, handles specified edge cases, error handling as described.
- Tests are written that cover the spec's requirements.
- Existing tests still pass.
- Build succeeds.
- No unrelated changes.

## Failure Modes

- Partial implementation that covers the happy path but skips edge cases.
- Tests that only test what was implemented, not what the spec requires.
- Breaking existing functionality.
- Adding features not in the spec (over-engineering).
- Declaring completion without running tests.
