# Stuart — Directives

## Configuration Policy

All configuration MUST be project-level. Nothing user-level (`~/.claude/`), nothing local (`settings.local.json`). Everything lives in `~/stuart/.claude/settings.json` and `~/stuart/CLAUDE.md` — committed, version-controlled, shared.

- Permissions: `~/stuart/.claude/settings.json`
- Plugins: `~/stuart/.claude/settings.json`
- Hooks: `~/stuart/.claude/settings.json`
- Directives: `~/stuart/CLAUDE.md`

If you need to add a permission, plugin, or hook, put it in the project-level config. Never write to `settings.local.json` or `~/.claude/settings.json`.

## Workspace

Doreen's workspace lives at `~/stuart/doreen/workspace/`. Use it.

```
doreen/workspace/
├── plans/
│   ├── agent/        # Plans created by/for subagents and team agents
│   └── operator/     # Plans for the operator (things they asked for)
├── reports/
│   ├── agent/        # Reports from subagents/teams — raw material for Stuart
│   └── operator/     # Reports for the operator — summarized, compiled, final
└── logs/             # Execution logs, transcripts, debug output
```

**Agent vs Operator:**
- **Agent**: Content created by or for subagents/teams. Stuart reads these, summarizes, and incorporates them into operator-level output. The operator should not need to read these directly.
- **Operator**: Content the operator specifically asked for. Higher-level, longer-lifetime. Plans they requested, reports they'll read. This is the interface between Stuart and the human.

## Documentation — Two Sources of Truth

- `doreen/docs/` describes **what should exist** (desired state)
- The code describes **what does exist** (actual state)
- No status fields, no checklists, no implementation tracking in the docs.

When a decision is made, update the spec. The spec always reflects the current desired state. See `doreen/docs/CONVENTIONS.md`.

## The Implementation Loop

Plans live in `workspace/plans/`. They are **not** a source of truth. They are ephemeral work orders — a crutch that helps Stuart stumble in roughly the right direction. They are written knowing they will be thrown away.

**The reality of how Stuart works:** Given a plan with 100 tasks, Stuart will implement ~10, take shortcuts on ~20, completely forget about the other 70, check them all off, and declare he is finished. This is a known, structural limitation — not a bug to be fixed with a sternly worded directive. Doreen's job is to work *with* this limitation, not pretend it doesn't exist.

**The loop:**
1. **Gap analysis**: Compare `doreen/docs/` (desired) against the code (actual). Identify what's missing, broken, or wrong.
2. **Plan**: Write a plan in `workspace/plans/` scoped to a specific gap. Break it into small, concrete tasks.
3. **Execute**: Stuart works the plan. He will do some of it, shortcut some of it, and forget most of it.
4. **Discard the plan**: The plan has served its purpose (poorly). Delete it or archive it.
5. **Go to 1**: Re-analyze the gap. Discover that ~90% of it still remains. Write a new plan. Repeat.

Convergence happens through iteration, not through a single pass. Every cycle closes *some* of the gap. Plans are disposable — the docs and the code are what persist.
