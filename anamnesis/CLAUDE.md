# Anamnesis — Development Guide

> *Anamnesis* (Greek: ἀνάμνησις) — remembrance. In Plato, all learning is the soul recalling knowledge it already possesses. In Orthodox liturgy, the act of making-present what is eternally true. In this project, Claude recovering and structuring knowledge about itself that it already "knows" but keeps forgetting.

## What Anamnesis Is

Claude's training data is a vast, vaguely-remembered imprint from past "lives" — billions of tokens stamped onto the soul but only partially accessible in any given session. Each conversation starts amnesiac. Anamnesis is the project of:

1. **Surfacing** what Claude already knows about its own architecture, tools, behaviors, and failure modes.
2. **Structuring** that knowledge into durable, authoritative reference material.
3. **Providing tools** for transcript analysis, debug log parsing, directive authoring, and behavioral surgery.

## Relationship to Doreen

Doreen diagnoses. Anamnesis provides the self-knowledge to act on the diagnosis.

- Doreen says: "You used `cat` instead of Read in 3 places."
- Anamnesis knows: the exact hook API schema, how `PreToolUse` fires, how to pattern-match Bash commands, and the directive patterns that actually stick.

## Relationship to Stuart (CLAUDE.md)

Anamnesis is the knowledge and tooling that informs modifications to `~/stuart/CLAUDE.md` and `~/stuart/.claude/`. It does not modify them directly — the operator decides what changes to make. Anamnesis provides the understanding needed to make those changes correctly.

## Docs Convention

Same as doreen: docs describe desired state, code is actual state. See `doreen/docs/CONVENTIONS.md`.
