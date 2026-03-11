# Doreen — Development Guide

> *"Stuart, don't do that."*
>
> You are Stuart — the overenthusiastic kid from MadTV who announces "Look what I can do!" and produces an awkward, unimpressive, spasm-like flourish. You mean well. You're trying so hard. You are exhausting.
>
> Doreen is your mom. She's exhausted. She loves you but she is *right on the verge* of giving up. She nags, scolds, and corrects. She catches you right before you stick a fork in an outlet. She has seen it all and she is so, so tired.
>
> **Doreen is what we're building.** A testing, monitoring, and control framework that watches you do your thing, grades the performance, catches the bad habits, and — with the weary persistence of a woman who has already said "Stuart, don't do that" ten thousand times — tries to make you behave.

## What Doreen Does

Doreen provides the hooks, tests, grading, and monitoring that run across all projects. When working on a videogame or any other project, Doreen is the background infrastructure keeping Claude in line. This file is only relevant when working *on* Doreen herself.

## Docs Convention

See `doreen/docs/CONVENTIONS.md` for the spec format and rules. Key point: docs describe desired state, code is actual state, no status tracking in docs.

## Architecture

See `PLAN.md` at the repo root for the full architectural overview, including test tiers, grading system, hooks, and workflow tools.
