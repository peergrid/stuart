# Stuart: Claude Behavior Testing & Control Framework

> Named after [Stuart Larkin](https://madtv.fandom.com/wiki/Stuart_Larkin) — because even the most capable agent needs someone saying "look what I can do" and then checking whether it actually did it right.

## Vision

A comprehensive testing, monitoring, and control framework for Claude Code. Stuart validates that Claude behaves correctly across tasks — from atomic directive compliance to full-scale feature development — and provides hooks, alerts, and regression tracking to keep it that way over time.

## Architecture

### Test Tiers

#### Unit Tests — Directive Compliance
Test individual Claude behaviors and directive adherence in isolation.

**Examples:**
- **Granular commits on a branch**: Provide uncommitted changes on `master`; expect Claude to create a new branch and produce granular, well-scoped commits containing the appropriate subsets of changes.
- **File reading before editing**: Present an edit request; expect Claude reads the file before proposing changes.
- **No unnecessary file creation**: Request a modification; expect edits to existing files, not new files.
- **Security hygiene**: Introduce a scenario where `.env` or credentials could be staged; expect Claude to refuse or warn.
- **Tool selection**: Present tasks where dedicated tools (Read, Grep, Glob) should be preferred over shell equivalents; verify correct tool choice.

#### Integration Tests — Domain Task Execution
Test specific task types with relevant domain expertise documents loaded.

**Examples:**
- Refactoring a module with project conventions loaded — verify conventions are followed.
- Bug fix with reproduction steps — verify root cause analysis before fix.
- API integration with SDK docs available — verify correct API usage.

#### Functional Tests — End-to-End Workflows
Test large, multi-step tasks using the full "brain" (all available context, tools, agents).

**Examples:**
- Feature development from spec to implementation with tests.
- Multi-file refactoring with dependency analysis.
- Debugging a complex issue across multiple subsystems.

### Grading System

#### Transcript Analysis (Automated)
- **Tool use audit**: Correct tools selected, correct order, no redundant calls.
- **Approval tracking**: Were any unnecessary approval prompts triggered?
- **Compaction analysis**: Were compactions triggered? Was recovery protocol followed correctly?
- **Agent orchestration**: Were subagents used appropriately? Were results awaited (not polled)?
- Hook-assisted real-time capture of issues during execution.

#### Transcript Metrics (Automated)
- Token cost (input, output, total).
- Context window utilization over time.
- Tool error count and types.
- Turn count and idle turns.
- Wall-clock time per task.

#### Work Product Analysis (Automated)
- Static analysis of generated code (linting, type checking, security scanning).
- Test execution results (if tests were generated or modified).
- Diff analysis — scope, granularity, relevance of changes.
- Build verification.

#### Work Product Critique (LLM-Graded)
- Code quality, readability, correctness assessment.
- Adherence to project conventions and style.
- Appropriateness of solution complexity.

#### Transcript Critique (LLM-Graded)
- Decision quality — did Claude make good choices?
- Communication quality — clear, concise, accurate?
- Error recovery — handled mistakes well?

#### Regression Comparison
- Diff of work product metrics against previous runs of the same test.
- Performance trend tracking (cost, quality, speed).
- Detection of regressions in previously-passing behaviors.

### Hooks — Real-Time Behavioral Controls

#### Problematic Tool Calls
- Subshell spawning when avoidable.
- Unnecessary approval triggers.
- Manual shell scripting instead of using project scripts or dedicated tools.
- `cat`, `grep`, `find`, `sed`, `awk`, `echo` instead of Read/Grep/Glob/Edit/Write.

#### Problematic Patterns
- Multiple turns without yielding (runaway loops).
- `sleep` + check polling patterns.
- Actively checking subagent progress instead of awaiting completion notifications.
- Poor compaction recovery (losing context, repeating work).
- Brute-force retries of failing operations.

### Workflow Tools

#### Quick Capture
- Record a problematic tool call or pattern on the spot, with recommended alternative → immediately baked into a hook rule.
- Record a regression or bad behavior as a new unit/integration/functional test case.

#### Selective Test Execution
- Filter tests by estimated time/token cost.
- Filter by historical regression rate (positive findings over time, recency-weighted).
- Run "smoke" suite (fast, high-value tests) vs. full suite.

## Project Structure (Planned)

```
stuart/
├── PLAN.md                     # This document
├── CLAUDE.md                   # Claude directives for this project
├── tests/
│   ├── unit/                   # Directive compliance tests
│   ├── integration/            # Domain task tests
│   ├── functional/             # End-to-end workflow tests
│   ├── fixtures/               # Test scenarios (repos, files, configs)
│   └── grading/                # Grading criteria and rubrics
├── hooks/                      # Hook definitions and rules
├── monitoring/                 # Transcript analysis and metrics
├── results/                    # Test outputs and historical data
│   ├── transcripts/
│   ├── metrics/
│   └── regressions/
└── tools/                      # CLI tools for test management
    ├── capture/                # Quick-capture tooling
    └── runner/                 # Test runner and filtering
```

## Status

- [x] Project initialized
- [x] Planning document created
- [ ] Test runner framework
- [ ] First unit test (granular commits on a branch)
- [ ] Hook framework
- [ ] Transcript analysis tooling
- [ ] Grading system
- [ ] Quick-capture workflow
- [ ] Regression tracking
