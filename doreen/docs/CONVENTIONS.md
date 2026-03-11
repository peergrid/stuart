# Doreen Docs — Conventions

This directory is the **single source of truth** for Doreen. It must be kept current as decisions are made, changes are executed, and features are implemented.

## The Contract

A fresh agent with zero context should be able to:
1. Read everything in `doreen/docs/`
2. Compare the specs against the actual codebase
3. Fix any bugs (reality doesn't match an implemented spec)
4. Implement any planned specs (spec exists but isn't built yet)

If the docs can't support that workflow, they're broken.

## Spec Format

Every spec file uses YAML frontmatter:

```markdown
---
status: planned | in-progress | implemented
owner: who decided/wrote this
created: YYYY-MM-DD
updated: YYYY-MM-DD
---
```

**Status meanings:**
- `planned` — Decided, specced, not yet built. A fresh agent should implement this.
- `in-progress` — Partially built. The checklist below shows what's done and what remains.
- `implemented` — Fully built and matching the spec. A fresh agent should verify this still holds and fix any drift.

## Requirements Within Specs

Each spec contains requirements as a checklist. This is the granular status tracker — not a separate system, not a task board, just the spec itself.

```markdown
## Requirements

- [x] Thing that is built and working
- [ ] Thing that is specced but not yet built
- [x] Another built thing
- [ ] Another planned thing
```

A fresh agent scans for `- [ ]` to find unfinished work and `- [x]` to find things to verify.

## Rules

1. **Specs are written when decisions are made**, not after implementation.
2. **Checkboxes are checked when code is committed**, not when code is planned.
3. **If reality contradicts a spec**, either fix the code or update the spec — never leave them out of sync.
4. **Never delete a spec**. If something is removed, mark it `status: deprecated` with a note explaining why.
5. **Keep specs atomic**. One feature/component per file. Cross-reference by filename if needed.

## File Naming

`kebab-case.md`. Name describes the component or feature, not a ticket number or date.

Examples: `test-runner.md`, `grading-system.md`, `hook-framework.md`, `transcript-analysis.md`.
