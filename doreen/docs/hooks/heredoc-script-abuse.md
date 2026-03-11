# Hook: HEREDOC and Temporary Script Abuse

## Behavior to Intercept

Claude creating temporary shell scripts, writing complex multi-line bash via HEREDOCs, or constructing elaborate one-liners instead of using project scripts, dedicated tools, or the Write tool to create a proper script.

## Trigger

`PreToolUse` on `Bash` tool. Inspect the command for:
- `cat <<EOF` or `cat <<'EOF'` writing to a file (should use Write tool).
- `bash -c '...'` with multi-line or complex content.
- Commands exceeding a character threshold (e.g., 200+ chars) that involve file manipulation.
- Creating files in `/tmp` that are really project scripts in disguise.

## Action

Block the tool call. Return a message explaining:
- If writing a file: "Use the Write tool instead of HEREDOC/echo redirection."
- If running a complex script: "This is too complex for an inline bash command. Create a proper script with the Write tool, or break this into multiple dedicated tool calls."

## Exceptions

- Git commit messages using HEREDOC (this is the established pattern for multi-line commit messages).
- Commands that genuinely need shell features like pipes or subshells for a simple, clear purpose.
