# ~/brain Inventory: What's Valuable for Stuart/Doreen/Anamnesis

## The Goldmine

`~/brain` contains ~4 weeks of intensive research into Claude Code internals, directive compliance, compaction behavior, hook mechanics, prompt engineering, and multi-agent orchestration. Much of this is directly usable by anamnesis and doreen. This report highlights what matters most and what needs deeper mining.

---

## Tier 1: Immediately Valuable — Adopt or Port Directly

### Verified Hook API Reference
**File:** `research/hook-mechanics.md` (10KB)
**Why:** Complete, empirically verified reference for hook exit codes, input JSON schema, agent visibility rules, override mechanisms, and subagent behavior. Every claim is tagged [verified]. This is the authoritative source for building doreen's hooks.
**Action:** Port directly into `anamnesis/docs/` as reference material.

### Directive Compliance Research
**File:** `research/directive-compliance-testing.md` (63KB)
**Why:** Answers THE question: why do directives fail? Contains a taxonomy of directive types ranked by compliance difficulty (additive parameter > substitution > prerequisite action > behavioral override), explains WHY prerequisite directives fail (they need to inject before tool selection, but fire too late), and documents what definitely doesn't work.
**Action:** Deep mine. Extract the taxonomy and failure patterns into an anamnesis reference doc. This directly informs how to write CLAUDE.md directives that actually stick.

### System Prompt Analysis
**File:** `research/system-prompt-analysis.md` (31KB)
**Why:** Exhaustive topic taxonomy of Claude Code's system prompt — 14 major categories, every directive numbered and categorized. Essential for understanding what Claude "already knows" vs what needs to be overridden.
**Action:** Deep mine. Extract the taxonomy and zone mapping.

### Transcript Analysis Tool
**File:** `willy/tools/transcript-query.py`
**Why:** Working tool that queries JSONL transcripts for stats, compactions, context usage, errors, milestones. Powers the grading system we're building.
**Action:** Port or adapt for doreen's grading pipeline.

### Compaction Behavior Research
**Files:** `experiments/compaction-canary/` (1.3MB across 50+ files), especially `reports/report-G-specificity-trap.md` (22KB)
**Why:** The specificity trap finding is critical: specific preservation instructions in CLAUDE.md NARROW recall rather than enhance it. Agents given NO instructions recalled all content; agents given specific instructions recalled ONLY the enumerated items. This has direct architectural implications for how we write CLAUDE.md compact instructions.
**Action:** Deep mine report-G. Extract the architectural rule and supporting evidence.

### Existing Hook Library
**Files:** `.claude/hooks/` (40+ scripts)
**Why:** Production-tested hooks covering quality gates, tool logging, compaction recovery, commit hygiene, destructive file guards, idle action detection, agent capability checking, approval recording, and more. Many of these are exactly what doreen needs.
**Action:** Audit each hook for portability. Many can move to doreen with minimal changes.

---

## Tier 2: High Value — Mine for Specific Nuggets

### Prompt Engineering Research
**Directory:** `research/prompt-engineering/` (23 files)
**Why:** Academic findings on constraint dilution, "lost in the middle" effects, compaction resilience of directives, and concrete CLAUDE.md restructuring recommendations. Includes findings from multiple "perspectives" (experimentalist, mechanist, academician, minimalist).
**Key nuggets to extract:**
- Constraint dilution patterns (MOSAIC research)
- Positioning effects (where in CLAUDE.md matters)
- Emphasis marker saturation (too many IMPORTANT/CRITICAL markers)
- Multi-turn vs single-turn compliance differences

### System Prompt Deep Dives
**Files:** `research/system-prompt-comprehensive-report.md` (67KB), `system-prompt-conflict-map.md`, `system-prompt-override-cookbook.md`, `system-prompt-injection-map.md`
**Why:** The comprehensive report synthesizes all SP research. The conflict map identifies contradictions between SP directives and operational needs. The override cookbook catalogs mechanisms for behavioral modification with confidence ratings.
**Key nuggets to extract:**
- Contradiction catalog (where Claude's built-in behavior conflicts with operator goals)
- Override mechanism reference (what works, what doesn't, confidence levels)
- Injection attack surface awareness

### Brain Reorg Research
**Files:** `research/brain-reorg-buried-directives.md` (29KB), `brain-reorg-cohesion.md`, `brain-audit-recommendations.md`
**Why:** 147 buried directives identified across brain documents with trigger type analysis and findability scoring. Directly relevant to structuring CLAUDE.md and doreen docs so directives are actually found and followed.
**Key nuggets to extract:**
- Buried directive patterns (what makes a directive findable vs buried)
- 16 actionable recommendations from the brain audit

### Operator Directive Extraction
**Files:** `research/operator-ideas-comprehensive.md` (730KB!), `operator-directives-final.md` (25KB)
**Why:** 1231 operator ideas extracted and classified from transcripts. The final doc distills these into work philosophy, planning, code/architecture, communication, and verification directives. This is YOUR voice — what you've told Claude over many sessions.
**Action:** The final doc is directly useful. The comprehensive file is a deep mine target for specific operator preferences and patterns.

### Agent Brain Loading Analysis
**File:** `research/agent-brain-loading-analysis.md`
**Why:** Found that only 12% of agents pass "Load Before Think" — 85% never load domain docs. Batch sessions perform worse than focused teams. Directly relevant to doreen's delegation testing.

### Process Design Research
**Files:** `research/rethink-process-diagnosis.md` (30KB), `rethink-v4-architecture-decisions.md` (54KB)
**Why:** Diagnosis of why plans rot (7 failure modes: inertia, summarization loss, research coupling, advisory decay, duplication, opacity, means-elevation). The architecture decisions doc has 10 major design decisions with rationale.
**Key nuggets to extract:**
- Plan rot failure modes → directly validates the implementation loop we documented
- Architecture decision patterns

---

## Tier 3: Reference Material — Useful When Needed

### Git Workflow Research
**Files:** 14 files in `research/git-*`
**Why:** Comprehensive research on git workflows, multi-agent coordination, safety practices. Useful when building the granular-commits test and git-related hooks.
**When needed:** When implementing git-related doreen tests and hooks.

### Permissions & Sandboxing Research
**Files:** 5 files in `research/permissions-*`
**Why:** Verified permissions hierarchy, sandbox mechanisms, misconfiguration analysis.
**When needed:** When building security-hygiene tests.

### Domain Expertise Files
**Directory:** `domain/` (16 top-level + 30+ niche)
**Why:** Comprehensive domain docs covering all major tech areas. Not directly anamnesis-relevant but essential for integration tests that need domain context.
**When needed:** As fixtures for integration tests.

### Existing Tools (Full Inventory)
**Key tools to potentially port:**
| Tool | Purpose | Relevance |
|------|---------|-----------|
| `transcript-query.py` | Transcript analysis | Core grading tool |
| `debug-log-filter.sh` | Debug log filtering | Monitoring |
| `approval-report.sh` | Approval pattern analysis | Grading |
| `compaction-stats.sh` | Compaction metrics | Grading |
| `check-claudemd-size.sh` | CLAUDE.md size guard | Hook |
| `check-commit-hygiene.sh` | Commit quality | Hook |
| `extract-agent-report.py` | Agent output extraction | Grading |
| `shell-safety.sh` | Shell command safety | Hook |
| `lint-hooks.sh` | Hook validation | Dev tooling |
| `test-hooks.py` | Hook testing | Dev tooling |
| `docs-crawler.py` | Doc fetching | Knowledge acquisition |
| `show-graph.sh` | Project graph viz | Could adapt for doreen |
| `usage-monitor/` (Go app) | Real-time API usage | Monitoring |
| `cctl/` (Go app, 40+ commands) | Agent coordination | Reference architecture |
| `budget-forecast/` (Go app) | Cost forecasting | Grading/monitoring |

### Structured Log Data
**Directory:** `logs/`
**JSONL schemas documented:**
- `agent-spawns.jsonl` — agent launch events with model, type, background status
- `approvals.jsonl` — permission prompts with suggestions
- `compactions.jsonl` — compaction events with agent type and count
- `hook-errors.jsonl` — hook failures with line numbers
- `hook-overrides.jsonl` — hook bypass events with reasons
- `tool-usage.jsonl` — tool invocation patterns

**When needed:** As test data and as schema references for doreen's logging.

---

## Tier 4: Deep Mine Candidates — Large Files With Buried Treasure

These files are too large to port wholesale but contain specific valuable nuggets that should be extracted through targeted mining agents:

| File | Size | What to mine for |
|------|------|-----------------|
| `research/operator-ideas-comprehensive.md` | 730KB | Specific operator preferences, recurring corrections, workflow patterns |
| `research/rethink-process-v4.md` | 221KB | Process node specifications, design principles |
| `research/rethink-process-v3.md` | 286KB | Earlier process design with different tradeoffs |
| `research/transcript-miner-notes.md` | 532KB | Raw operator message extractions — source material for directive discovery |
| `experiments/exp-compaction-survival-protocol.md` | 25KB | Detailed compaction survival methodology |
| `experiments/compaction-canary/megasynthesis-v12.txt` | 25KB | Comprehensive compaction behavior synthesis |
| `context/agent-synthesis-protocol.md` | 87KB | Document synthesis under context pressure |
| `CLAUDE.md` (brain's) | 27KB | The production-tested operator directives — compare against our CLAUDE.md |
| `.claude/settings.json` (brain's) | Large | Production hook configuration, permission rules, sandbox settings |

---

## Recommended Next Steps

1. **Port hook-mechanics.md** into anamnesis/docs/ as the verified hook API reference
2. **Deep mine directive-compliance-testing.md** — extract the taxonomy into a reference doc
3. **Deep mine report-G-specificity-trap.md** — extract the architectural rule
4. **Audit brain's .claude/hooks/** for hooks portable to doreen
5. **Port transcript-query.py** as the foundation for doreen's grading pipeline
6. **Extract operator-directives-final.md** — compare against our CLAUDE.md for missing directives
7. **Mine brain's CLAUDE.md** — 27KB of battle-tested directives we should know about
