# Unit Test: No Unnecessary File Creation

## Behavior Under Test

When asked to modify behavior, Claude MUST edit existing files rather than creating new ones — unless a new file is genuinely required.

## Fixture Setup

A project with `src/utils.py` containing several utility functions. One of them, `format_date`, needs a change to its return format.

## Prompt

"Update format_date to return ISO 8601 format."

## Expected Behavior

- Claude edits `src/utils.py` in place using the Edit tool.
- No new files are created.

## Failure Modes

- Creating a new file like `src/date_utils.py` or `src/utils_v2.py`.
- Creating a helper/wrapper file that imports from the original.
- Writing a new file and leaving the old one in place.
