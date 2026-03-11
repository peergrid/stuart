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

## Documentation — Single Source of Truth

`doreen/docs/` is the authoritative spec for everything in Doreen. See `doreen/docs/CONVENTIONS.md` for the full format.

**The rule:** When a decision is made, a spec is written (or updated). When code is committed, the spec's checklist is updated. If reality and the spec disagree, one of them must be fixed immediately — never leave them out of sync.

A fresh agent with no context must be able to read `doreen/docs/`, compare against the codebase, and know exactly what's done, what's broken, and what's left to build.
