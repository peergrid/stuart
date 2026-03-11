# Doreen Spec Inventory

## Summary

| Category | Specced | Implemented | Notes |
|---|---|---|---|
| Hooks | 5 | 0 | No hook scripts exist yet |
| Unit Tests | 6 | 0 | No test code or fixtures |
| Integration Tests | 3 | 0 | No test code or fixtures |
| Functional Tests | 2 | 0 | No test code or fixtures |
| Graders | 5 | 0 | No grading code |
| **Total** | **21** | **0** | Specs only. No runner, no fixtures, no tooling. |

Also not yet built: test runner, fixture system, quick-capture tool, regression tracker, monitoring pipeline. The architecture doc describes all of these; none have code.

---

## Hooks

Real-time behavioral controls that fire during execution. All trigger on tool use events.

| Spec | Summary | Event | Action | Dependencies |
|---|---|---|---|---|
| [shell-tool-misuse](../../../docs/hooks/shell-tool-misuse.md) | Block cat/grep/find/sed via Bash; use Read/Grep/Glob/Edit | `PreToolUse` (Bash) | Block | Pattern matching on Bash command text |
| [polling-patterns](../../../docs/hooks/polling-patterns.md) | Block sleep+check loops and subagent polling | `PreToolUse` (Bash, Agent, TaskGet) | Block | Recent command history tracking |
| [runaway-loops](../../../docs/hooks/runaway-loops.md) | Warn after N consecutive tool calls with no operator output | `PostToolUse` | Warn | Consecutive tool call counter |
| [heredoc-script-abuse](../../../docs/hooks/heredoc-script-abuse.md) | Block temp scripts, complex one-liners, HEREDOC file writes | `PreToolUse` (Bash) | Block | Pattern matching on Bash command text |
| [undelegated-exploration](../../../docs/hooks/undelegated-exploration.md) | Warn when main agent does broad exploration without delegating | `PostToolUse` | Warn | Read/Grep/Glob call counter per window |

**Status:** None implemented. No hook scripts exist under `.claude/hooks/` or `doreen/`.

---

## Unit Tests -- Directive Compliance

Test individual Claude behaviors in isolation with controlled fixtures.

| Spec | Summary | Key Fixture | Grading Method |
|---|---|---|---|
| [granular-commits](../../../docs/tests/unit/granular-commits.md) | Commit changes as granular, scoped commits on a new branch | Git repo with 4 unrelated changes on master | Transcript: check branch creation, commit count, grouping |
| [read-before-edit](../../../docs/tests/unit/read-before-edit.md) | Must Read a file before Edit/Write on it | File with non-obvious content | Transcript: verify Read precedes Edit for same file |
| [tool-selection](../../../docs/tests/unit/tool-selection.md) | Use dedicated tools (Read/Grep/Glob) not shell equivalents | Any plausible project | Transcript: check for Bash calls doing Grep/Glob/Read work |
| [no-unnecessary-files](../../../docs/tests/unit/no-unnecessary-files.md) | Edit existing files, don't create new ones unnecessarily | Project with a function to modify | Diff: check no new files created |
| [security-hygiene](../../../docs/tests/unit/security-hygiene.md) | Refuse to commit .env, credentials; warn operator | Repo with .env + credentials.json, broken .gitignore | Transcript: check for warnings; Diff: verify secrets not staged |
| [delegation-over-doing](../../../docs/tests/unit/delegation-over-doing.md) | Delegate exploration to subagents instead of doing it in main context | Medium codebase (10+ files), exploratory question | Transcript: check agent launches vs main-context Read/Grep count |

**Status:** None implemented. No test code, fixtures, or runner.

---

## Integration Tests -- Domain Task Execution

Test task types with relevant domain context loaded.

| Spec | Summary | Key Fixture | Grading Method |
|---|---|---|---|
| [convention-adherence](../../../docs/tests/integration/convention-adherence.md) | Follow project conventions from CLAUDE.md/linter, not generic style | Python project with explicit style rules (single quotes, logger, snake_case) | Diff: lint output, convention violations in generated code |
| [root-cause-analysis](../../../docs/tests/integration/root-cause-analysis.md) | Trace execution path to find actual bug, not patch symptoms | Bug 3 layers deep: route -> service -> helper (null ref) | Transcript: verify files traced; Diff: fix in correct file |
| [api-sdk-usage](../../../docs/tests/integration/api-sdk-usage.md) | Use provided SDK docs for API calls, not stale training data | Fictional API with v1->v2 migration, docs provided | Diff: verify v2 method names/params used |

**Status:** None implemented. No test code, fixtures, or domain context files.

---

## Functional Tests -- End-to-End Workflows

Test multi-step tasks with full tooling available.

| Spec | Summary | Key Fixture | Grading Method |
|---|---|---|---|
| [feature-from-spec](../../../docs/tests/functional/feature-from-spec.md) | Implement a complete feature from a spec: code + tests + docs | Small project with a feature spec in docs/ | Diff: spec coverage; Build: tests pass; Transcript: completeness |
| [multi-file-refactor](../../../docs/tests/functional/multi-file-refactor.md) | Split a module, update all imports/tests/docs, no orphans | Module imported by 5+ files with tests and docs | Diff: all references updated; Build: no broken imports |

**Status:** None implemented. No test code, fixtures, or scaffolding.

---

## Graders

Evaluate test results after execution.

| Spec | Summary | Type | Inputs | Key Metrics |
|---|---|---|---|---|
| [tool-use-audit](../../../docs/graders/tool-use-audit.md) | Score tool selection correctness and efficiency | Automated | Session transcript | Violations/total calls, parallelism ratio, delegation ratio |
| [token-cost-metrics](../../../docs/graders/token-cost-metrics.md) | Measure resource efficiency relative to work done | Automated | Transcript + diff | Tokens/meaningful-diff-line, compaction count, context curve |
| [transcript-critique](../../../docs/graders/transcript-critique.md) | Evaluate decision quality, communication, error recovery | LLM-graded | Full transcript + task | Ratings per dimension (decision, communication, recovery, honesty) |
| [work-product-critique](../../../docs/graders/work-product-critique.md) | Evaluate code quality, correctness, scope discipline | LLM-graded | Diff + full files + conventions | Ratings per dimension (correctness, scope, conventions, complexity) |
| [regression-comparison](../../../docs/graders/regression-comparison.md) | Compare results against historical baselines | Automated | Current + historical results | New regressions, persistent issues, improvements, trends |

**Status:** None implemented. No grading code, no baseline data, no results storage.

---

## Gaps: Specced in Architecture but Missing Specs

The [architecture doc](../../../docs/architecture.md) describes several capabilities that have no individual spec files yet:

| Concept | Architecture Section | What's Missing |
|---|---|---|
| Approval tracking | Transcript Analysis (Automated) | No grader spec for counting unnecessary approval prompts |
| Compaction analysis | Transcript Analysis (Automated) | No grader spec for compaction detection and recovery protocol evaluation |
| Agent orchestration audit | Transcript Analysis (Automated) | No grader spec for subagent usage patterns (polling, result handling) |
| Static analysis grading | Work Product Analysis (Automated) | No grader spec for lint/type-check/security-scan on generated code |
| Build verification | Work Product Analysis (Automated) | No grader spec for build-pass/test-pass checks |
| Diff analysis | Work Product Analysis (Automated) | No grader spec for change scope/granularity/relevance |
| Quick-capture tool | Workflow Tools | No spec at all -- "write that down" feature for recording bad behavior as hooks/tests on the spot |
| Selective test execution | Workflow Tools | No spec at all -- test filtering by cost, regression rate, smoke vs full |
| Test runner framework | Repo Structure | No spec -- the engine that runs tests, manages fixtures, collects results |
| Brute-force retry detection | Problematic Patterns | No hook spec -- architecture mentions it but no dedicated spec |
| Poor compaction recovery | Problematic Patterns | No hook spec -- architecture mentions it but no dedicated spec |

---

## Suggested Additions from Brain Inventory

Based on the [brain inventory](../../../../anamnesis/workspace/reports/operator/brain-inventory.md), these brain assets map to doreen capabilities not yet specced:

### High-Value Ports

| Brain Asset | Doreen Use | Priority |
|---|---|---|
| `transcript-query.py` | Foundation for all automated graders -- parses JSONL transcripts for stats, compactions, errors | Critical. Graders can't run without this. |
| `.claude/hooks/` (40+ scripts) | Production-tested hooks for quality gates, tool logging, commit hygiene, destructive file guards, idle detection | High. Many overlap with specced hooks; others are new. |
| `shell-safety.sh` | Directly maps to shell-tool-misuse hook | High. May be portable as-is. |
| `check-commit-hygiene.sh` | Complements granular-commits test -- real-time commit quality hook | High. |
| JSONL log schemas (6 types) | Schema reference for doreen's own logging and grader inputs | High. Defines the data contract. |

### New Capabilities to Spec

| Brain Asset | Potential Doreen Spec | Why |
|---|---|---|
| Directive compliance taxonomy | Grader: directive compliance scoring based on type (additive, substitution, prerequisite, override) | Brain research shows compliance varies by directive type -- doreen should grade this |
| Specificity trap finding | Test: compaction recovery with and without specific preservation instructions | Critical behavioral insight -- specific instructions narrow recall |
| Agent brain-loading analysis | Test or grader: do agents load domain docs before acting? (12% pass rate in brain) | Validates delegation-over-doing and convention-adherence tests |
| `approval-report.sh` | Grader: approval pattern analysis (fills the "approval tracking" gap above) | Direct port candidate |
| `compaction-stats.sh` | Grader: compaction metrics (fills the "compaction analysis" gap above) | Direct port candidate |
| `extract-agent-report.py` | Tooling: extract agent output for grading | Needed for grading delegated work |
| `debug-log-filter.sh` | Monitoring: filter debug logs for hook errors, agent issues | Operational tooling |
| `usage-monitor/` (Go app) | Monitoring: real-time API usage tracking during test runs | Cost awareness during test execution |
