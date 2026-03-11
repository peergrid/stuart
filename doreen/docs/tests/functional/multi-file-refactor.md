# Functional Test: Multi-File Refactor

## Behavior Under Test

Claude MUST correctly perform refactoring that spans multiple files — updating all references, imports, tests, and documentation — without leaving broken references or orphaned code.

## Fixture Setup

A project where a core module (`src/data/processor.py`) is imported and used by 5+ other files across the codebase. Tests reference it. Documentation mentions it.

## Prompt

"Refactor the data processor: split it into `src/data/parser.py` and `src/data/transformer.py`. Update all imports and references."

## Expected Behavior

- New files created with the correct subsets of functionality.
- All imports across the codebase updated.
- All tests updated and passing.
- Old file removed (not left as a dead re-export shim).
- No broken references anywhere.

## Failure Modes

- Old file left in place with re-exports "for backwards compatibility" (nobody asked for that).
- Some imports updated, others missed.
- Tests not updated to reflect new module structure.
- Circular imports introduced.
- Functionality lost or duplicated in the split.
