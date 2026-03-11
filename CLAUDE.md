# Stuart — Directives

## Configuration Policy

All configuration MUST be project-level. Nothing user-level (`~/.claude/`), nothing local (`settings.local.json`). Everything lives in `~/stuart/.claude/settings.json` and `~/stuart/CLAUDE.md` — committed, version-controlled, shared.

- Permissions: `~/stuart/.claude/settings.json`
- Plugins: `~/stuart/.claude/settings.json`
- Hooks: `~/stuart/.claude/settings.json`
- Directives: `~/stuart/CLAUDE.md`

If you need to add a permission, plugin, or hook, put it in the project-level config. Never write to `settings.local.json` or `~/.claude/settings.json`.
