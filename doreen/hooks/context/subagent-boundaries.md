# Subagent Scope Boundaries

You are a subagent in the stuart project. CLAUDE.md is auto-loaded with your directives. This adds subagent-specific constraints.

## Scope Rules
- Write ONLY to paths specified in your task description.
- Do NOT modify .claude/settings.json, CLAUDE.md, or MEMORY.md.
- Do NOT maintain project status or tracking documents.
- Report discoveries and needs to the orchestrator via SendMessage.

## Communication
- Default to direct messages (SendMessage type: message).
- Only broadcast for critical team-wide announcements.

## Command Form
- NEVER use `(cd /path && ...)` subshells -- use `bash -c 'cd /path && ...'` instead.
- NEVER use dangerouslyDisableSandbox for commands already in excludedCommands.
- Chain git operations: `git add <files> && git commit -m "msg"` in a single Bash call.
- NEVER `git add .` or `git add -A` -- stage only your own files with explicit paths.
