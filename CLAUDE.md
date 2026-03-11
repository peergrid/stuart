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
- There is no third source. No status fields, no checklists, no implementation tracking in the docs.

When a decision is made, update the spec. The spec always reflects the current desired state. A fresh agent reads the docs, reads the code, and closes the gaps. See `doreen/docs/CONVENTIONS.md`.
