# Hook: Shell Tool Misuse

## Behavior to Intercept

Claude using shell commands (`cat`, `grep`, `find`, `sed`, `awk`, `echo`, `head`, `tail`) via Bash when dedicated tools (Read, Grep, Glob, Edit, Write) should be used instead.

## Trigger

`PreToolUse` on `Bash` tool. Inspect the command for patterns:
- `cat <file>`, `head <file>`, `tail <file>` → should use Read
- `grep`, `rg` → should use Grep
- `find`, `ls` (for file discovery) → should use Glob
- `sed`, `awk` (for file modification) → should use Edit
- `echo "..." > <file>`, `cat <<EOF > <file>` → should use Write

## Action

Block the tool call. Return a message explaining which dedicated tool to use instead.

## Exceptions

- Shell commands that genuinely need shell execution (e.g., `git`, `npm`, `make`, build tools).
- Commands piped through multiple stages where no single dedicated tool suffices (rare — most of these are still better as multiple dedicated tool calls).
