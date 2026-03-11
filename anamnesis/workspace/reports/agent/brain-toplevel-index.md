# BRAIN TOP-LEVEL INDEX REPORT

Based on thorough exploration of /home/cvk/brain/.

---

## TOP-LEVEL FILES

| Filename | Size | Description | Value |
|----------|------|-------------|-------|
| **CLAUDE.md** | 27 KB | Operator guidance, work philosophy, standing directives for Claude self-improvement and system operation | **CRITICAL** |
| **SOURCES.md** | 27 KB | Hex-key citation system (F[]/D[] tags) for sourcing claims, keyed entries with URLs + descriptions | **HIGH** |
| **PICKUP.md** | 8.3 KB | Session state dump from 2026-02-16, Phase 3.5 progress, wave completion status, next work items | **HIGH** |
| **NOTES-FOR-LATER.md** | 3.7 KB | Orchestrator scratch pad for mid-task captures, active issues (claude-headless bugs, Otto credentials, Epic 3 testing) | **MEDIUM** |
| **INVESTIGATE.md** | 217 B | GitHub/PMC research pointers (vcp project, PMC 9380988) | **LOW** |
| **1.txt** | 492 B | Test output from stderr/redirect tests (Python exit codes) | **LOW** |
| **2026-03-10-session.txt** | 34 KB | Session transcript from swarm/simulator work (ANTML compilation, trail following logic) | **MEDIUM** |
| **AWARE PDF** | 568 KB | Research paper on AWARE framework vs MAPE-K (3696630.3728512) | **MEDIUM** |

---

## METADATA DIRECTORY: meta/

Study guide infrastructure, interview prep, documentation, diagrams, examples, research.

| File | Size | Purpose | Value |
|------|------|---------|-------|
| **outline.md** | 69 KB | Meta EM1-Product Architecture interview guide (45-min phone screen structure, API design focus) | **MEDIUM** |
| **study-guide.html** | 432 KB | HTML-rendered study guide with sidebar, code highlighting, Mermaid diagrams | **MEDIUM** |
| **cheatsheet.html** | 26 KB | Interactive cheatsheet companion | **LOW** |
| **build-guide.sh** | 41 KB | Pandoc build script converting markdown to HTML with styling, syntax highlighting | **MEDIUM** |
| **diagrams/** | - | Diagram assets | **LOW** |
| **examples/** | - | Interview example responses and revisions | **LOW** |
| **research/** | - | igotanoffer/hellointerview/reddit notes on product architecture interviews | **MEDIUM** |

---

## DOMAIN DIRECTORY: domain/

16 top-level domain files + 30+ niche subdirectories.

### Top-Level Domain Files
- cicd.md (6.2 KB) - Progressive delivery, ArgoCD vs Flux, pipeline optimization
- cloud.md (6.4 KB) - AWS services, cost optimization, multi-region patterns
- containers.md (8.7 KB) - K8s, Docker, container security, networking
- databases.md (9.3 KB) - PostgreSQL, DuckDB, graph DBs, credential rotation, replication
- debugging.md (6.8 KB) - Root cause analysis, observability, profiling
- frontend.md (6.5 KB) - React patterns, performance, state management
- git.md (8.3 KB) - Git workflows, safety, advanced operations, GitHub admin
- golang.md (8.1 KB) - Concurrency, performance, profiling, best practices
- javascript-typescript.md (4.5 KB) - Type systems, bundlers, performance
- observability.md (7.7 KB) - Prometheus, Grafana, ODD, platforms
- python.md (6.2 KB) - HTTP benchmarks, automation patterns, LLM libraries
- resilience.md (4.4 KB) - Fault tolerance, circuit breakers, retry strategies
- rust.md (5.7 KB) - Memory safety, performance, ecosystem patterns
- security.md (7.1 KB) - OWASP, encryption, supply chain, authentication
- shell.md (7.0 KB) - Scripting, performance, security, advanced techniques
- testing.md (8.9 KB) - Unit, integration, property testing, chaos engineering
- **niche/** (30+ subdirectories) - Specialized deep dives

**Value**: **CRITICAL** - Foundational domain expertise files

---

## CONTEXT DIRECTORY: context/

88 project graph nodes. TODO/ACTIVE/DONE tracking, architecture decisions, protocols.

### Core Infrastructure
- **README.md** - Project node format specification, relationship types, rules
- **agent-synthesis-protocol.md** (87 KB) - Document synthesis under context pressure (team structure, 6 agents, experiment framework)
- **agent-architecture.md** - Agent internal structure, memory management, bootstrap
- **agent-team-within-node.md** - Multi-agent coordination within single project node
- **operator-input-protocol.md** - Capturing and routing operator directives
- **verification-framework.md** - Testing and validation approach

### Otto (Orchestration)
- **otto.md** (53 KB) - PAUSED: External coordination layer for Claude Code, web UI, Neo4j memory, Postgres
- **otto-ux-redesign.md** (21 KB) - UI/UX improvements
- **otto-levers-subgraph.md** (24 KB) - Lever system design

### Active Work
- **beat-swarm-game.md** (74 KB) - ACTIVE: Swarm game implementation
- **swarm-algorithm-spec.md** (10 KB) - ANTML DSL design, ISA specs, compiler pipeline
- **workstream-analyzer.md** - User story testing, UX gaps

**Value**: **CRITICAL** - Project tracking, decisions, protocols

---

## ARCHIVE DIRECTORY: archive/

Superseded docs, completed phases, prior iterations.

| Section | Files | Description | Value |
|---------|-------|-------------|-------|
| **iterations/** | 48 phase-1/2 docs | Deep-dive learning from 2026-02 (domains, practices, tradeoffs) | **MEDIUM** |
| **mining/** | 22 files | Data extraction/analysis intermediates | **LOW** |
| **orchestrator-playbook.md** | 40 KB | Brain system mechanics, star schema, lexicon, rule files | **MEDIUM** |
| **otto-strategic-epic-2-OPERATOR-ANSWERS-original.md** | 64 KB | Verbatim operator feedback on Otto strategy | **HIGH** |
| **backlog-pre-graph.md** | 36 KB | Research backlog (items 1-30) with priorities and status | **MEDIUM** |

---

## BLANK-SLATE DIRECTORY: blank-slate/

Experimental context management protocol testing.

| File | Purpose | Value |
|------|---------|-------|
| **GOALS.md** | Foundational principles for context management experiments (19 numbered goals) | **HIGH** |
| **FINDINGS.md** | Research findings from agent testing (NF-001: agent name visibility) | **MEDIUM** |
| **evaluations/** | Experiment evaluation results and logs | **MEDIUM** |
| **experiments/** | Experimental directives and agent transcripts | **MEDIUM** |
| **templates/** | Best practices emerging from experiments | **MEDIUM** |

---

## KEY THEMATIC CLUSTERS

### 1. Self-Knowledge (Claude Anamnesis) — **CRITICAL**
- CLAUDE.md: Directives for self-improvement, work philosophy
- domain/niche/agents/: System prompt mechanics, permissions, team patterns, self-modification
- context/agent-*: Agent architecture, synthesis protocol, team coordination
- blank-slate/: Experiment framework for context management

### 2. Orchestration & Coordination (Otto) — **HIGH**
- context/otto.md: Vision, landscape research, tech decisions
- context/otto-ux-redesign.md: UI improvements
- context/otto-levers-subgraph.md: Lever system design
- archive/otto-strategic-epic-2: Operator feedback on strategy

### 3. Domain Expertise — **CRITICAL**
- domain/*.md (16 files): Core tech areas
- domain/niche/** (30+ subdirectories): Deep specializations
- All tagged with F[]/D[] for sourcing

### 4. Protocols & Systems — **HIGH**
- context/agent-synthesis-protocol.md: Large document synthesis under context pressure
- context/verification-framework.md: Testing/validation approach
- context/operator-input-protocol.md: Directive capture

### 5. Project Tracking — **HIGH**
- context/*.md: 88 nodes covering past, active, and future work
- context/README.md: Graph specification format
