# Doreen Docs — Conventions

This directory is the **single source of truth** for what Doreen *should* be. The code is the single source of truth for what Doreen *is*. There is no third source.

## The Contract

A fresh agent with zero context should be able to:
1. Read everything in `doreen/docs/`
2. Read the actual codebase
3. Identify the gaps — what the docs describe that the code doesn't do, and what the code does that contradicts the docs
4. Close the gaps — implement what's missing, fix what's wrong

No status fields. No checklists. No implementation tracking in the docs. The docs describe the desired state. The code is the actual state. The agent is the reconciliation loop.

## Why No Status Tracking

Status tracking in docs is a third source of truth that can drift from both the docs and the code. If a spec says `status: implemented` but the code is broken, you now have *two* things that are wrong. If a checkbox says `[x]` but the feature regressed, the checkbox is a liar.

The code never lies about what's implemented. The docs should never lie about what's desired. That's enough.

## Spec Format

Specs describe *what should exist* — behavior, structure, interfaces, constraints. They are written in plain prose and are specific enough that an agent can read one, look at the codebase, and know unambiguously whether the code matches.

A good spec reads like a contract: "The test runner MUST do X. When Y happens, it MUST produce Z." An agent can verify each statement against reality.

## Rules

1. **Specs are written when decisions are made**, not after implementation.
2. **Specs describe desired state**, not implementation status. No `status:` fields, no `[x]`/`[ ]` checklists.
3. **If a decision changes**, update the spec. The spec always reflects the *current* desired state.
4. **Never delete a spec**. If something is removed from the desired state, replace the spec content with a short note explaining why it was removed and when.
5. **Keep specs atomic**. One feature/component per file. Cross-reference by filename if needed.

## File Naming

`kebab-case.md`. Name describes the component or feature, not a ticket number or date.

Examples: `test-runner.md`, `grading-system.md`, `hook-framework.md`, `transcript-analysis.md`.
