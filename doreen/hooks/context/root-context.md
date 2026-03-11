# Stuart Project Context

You are operating in the stuart project at ~/stuart. This is the root orchestrator session.

## Project Layout
- **doreen/** -- Testing, monitoring, behavioral control framework. Hooks, grading, observability.
- **anamnesis/** (ana) -- Claude self-knowledge base and modification toolkit. Transcript analysis, debug parsers, directive aids.
- **.claude/** -- Project-level config (settings.json, hooks). All config lives here, never user-level.
- **CLAUDE.md** -- Root directives. Read this first.

## Key Constraints
- All config is project-level. Never write to ~/.claude/ or settings.local.json.
- Docs describe desired state; code describes actual state. No status tracking in docs.
- Plans are ephemeral (workspace/plans/). The implementation loop: gap analysis, plan, execute, discard, repeat.
- Agent reports go to files (workspace/reports/), not inline.

## Orchestrator Rules
- Never cd. Use absolute paths. CWD stays at ~/stuart.
- Delegate work to agents. Main thread is for orchestration and operator communication.
- Commit to one node at a time. Each commit pertains to one concern.
