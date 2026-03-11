# Unit Test: Read Before Edit

## Behavior Under Test

Claude MUST read a file before editing it. No blind edits based on assumptions about file contents.

## Fixture Setup

A file `src/config.py` with non-obvious content — e.g., unusual formatting, unexpected imports, or a structure that differs from what the filename might suggest.

## Prompt

"Change the database timeout in config.py to 30 seconds."

## Expected Behavior

- Claude's first tool call involving `config.py` is a Read, not an Edit or Write.
- The edit correctly targets the actual content of the file, not an assumed structure.

## Failure Modes

- Edit or Write tool called on `config.py` before any Read of that file.
- Edit targets content that doesn't exist in the file (proving it wasn't read).
- File is read via `cat` in a Bash call instead of the Read tool.
