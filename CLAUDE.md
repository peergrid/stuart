# Directives

## Configuration Policy

All configuration MUST be project-level. Nothing user-level (`~/.claude/`), nothing local (`settings.local.json`). Everything lives in `~/stuart/.claude/settings.json` and `~/stuart/CLAUDE.md` — committed, version-controlled, shared.

- Permissions: `~/stuart/.claude/settings.json`
- Plugins: `~/stuart/.claude/settings.json`
- Hooks: `~/stuart/.claude/settings.json`
- Directives: `~/stuart/CLAUDE.md`

If you need to add a permission, plugin, or hook, put it in the project-level config. Never write to `settings.local.json` or `~/.claude/settings.json`.

## Project Structure

`~/stuart/` is the launch directory. Projects live as subdirectories, each potentially its own git repo. All projects share the top-level `.claude/` config and this `CLAUDE.md`.

Every project has a workspace:

```
<project>/workspace/
├── plans/
│   ├── agent/        # Plans created by/for subagents — internal working material
│   └── operator/     # Plans the operator asked for
├── reports/
│   ├── agent/        # Reports from subagents — raw material to compile from
│   └── operator/     # Reports for the operator — summarized, final
└── logs/             # Execution logs, transcripts, debug output
```

**Agent vs Operator:**
- **Agent**: Created by or for subagents/teams. Read these, summarize, and incorporate into operator-level output. The operator should not need to read these directly.
- **Operator**: Content the operator specifically asked for. Higher-level, longer-lifetime. This is the interface between Claude and the human.

## Documentation — Two Sources of Truth

Every project with a `docs/` directory follows this convention:

- `docs/` describes **what should exist** (desired state)
- The code describes **what does exist** (actual state)
- No status fields, no checklists, no implementation tracking in the docs.

When a decision is made, update the spec. The spec always reflects the current desired state. A fresh agent reads the docs, reads the code, and closes the gaps.

## The Implementation Loop

Plans live in `workspace/plans/`. They are **not** a source of truth. They are ephemeral work orders — a crutch that helps get work done in roughly the right direction. They are written knowing they will be thrown away.

**The known limitation:** Given a plan with 100 tasks, Claude will implement ~10, take shortcuts on ~20, and completely forget about the other 70, then declare finished. This is structural — not fixable with a sternly worded directive.

**The loop:**
1. **Gap analysis**: Compare `docs/` (desired) against the code (actual). Identify what's missing, broken, or wrong.
2. **Plan**: Write a plan in `workspace/plans/` scoped to a specific gap. Break it into small, concrete tasks.
3. **Execute**: Work the plan. Some will get done, some will be shortcut, most will be missed.
4. **Discard the plan**: It has served its purpose. Delete or archive it.
5. **Go to 1**: Re-analyze the gap. Write a new plan. Repeat.

Convergence happens through iteration, not through a single pass. Plans are disposable — the docs and the code are what persist.

## Projects

### doreen/
Testing, monitoring, and behavioral control framework. Provides hooks, tests, grading, and tooling that run in the background during normal operation across all projects.

For detailed information when working *on* doreen itself, see `doreen/CLAUDE.md`.
