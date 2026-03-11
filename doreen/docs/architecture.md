# Stuart — a.k.a. Claude

> *"Look what I can do!"* — [Stuart Larkin](https://madtv.fandom.com/wiki/Stuart_Larkin), MadTV
>
> Stuart is Claude. He's the overenthusiastic, unsupervised child who — without anyone asking — proudly announces "Look what I can do!" and then produces an awkward, unimpressive, spasm-like flourish. He means well. He's trying so hard. He is exhausting.

# Doreen — the Framework

> Doreen is Stuart's mom. She's exhausted. She's been doing this for years. She loves her son but she is *right on the verge* of giving up. She nags. She scolds. She corrects. She catches him right before he sticks a fork in an outlet. She has seen it all and she is so, so tired.
>
> **Doreen is what we're building.** A testing, monitoring, and control framework that watches Claude do his thing, grades the performance, catches the bad habits, and — with the weary persistence of a woman who has already said "Stuart, don't do that" ten thousand times — tries to make him behave.

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

### Grading System — The Report Card

How Doreen evaluates whether Stuart actually did a good job or just made a mess and called it art.

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

### Hooks — Doreen's Voice

Real-time behavioral controls. The nagging. The "Stuart, NO." The sigh followed by a correction. These fire during execution and either block, warn, or log.

#### Problematic Tool Calls
- Subshell spawning when avoidable.
- Unnecessary approval triggers.
- Temporary shell scripts, complex bash one-liners, or HEREDOC abuse instead of using project scripts or dedicated tools.
- `cat`, `grep`, `find`, `sed`, `awk`, `echo` instead of Read/Grep/Glob/Edit/Write.

#### Problematic Patterns
- Multiple turns without yielding (runaway loops).
- `sleep` + check polling patterns.
- Actively checking subagent progress instead of awaiting completion notifications.
- Poor compaction recovery (losing context, repeating work).
- Brute-force retries of failing operations.
- Main agent performing work or doing exploration instead of delegating to agents or teams.

### Workflow Tools — Doreen's Toolkit

#### Quick Capture — "Write That Down"
When Stuart does something stupid mid-session, Doreen needs to be able to immediately:
- Record the problematic tool call or pattern with the recommended alternative → baked into a hook rule on the spot.
- Record the regression or bad behavior as a new test case so he never gets away with it again.

#### Selective Test Execution — "Pick Your Battles"
Doreen can't yell about everything. She has to prioritize.
- Filter tests by estimated time/token cost.
- Filter by historical regression rate (positive findings over time, recency-weighted).
- Run "smoke" suite (fast, high-value tests) vs. full suite.

## Repo Structure

`~/stuart` is the launch point. It's where Claude Code starts, so its `.claude/` and `CLAUDE.md` govern Stuart's behavior across *all* work. Projects live as subdirectories — each may be its own git repo, but they sit physically under Stuart and inherit his top-level config.

```
~/stuart/                           # Home base. Stuart lives here. This is the git repo.
├── PLAN.md                         # This document
├── CLAUDE.md                       # Stuart's directives — what he must/must not do
├── .claude/                        # Claude Code config — hooks, settings
│   ├── settings.json               # Permissions, tool config
│   └── hooks/                      # Doreen's voice — real-time behavioral controls
│
├── doreen/                         # The testing & grading framework
│   ├── tests/
│   │   ├── unit/                   # Directive compliance tests
│   │   ├── integration/            # Domain task tests
│   │   ├── functional/             # End-to-end workflow tests
│   │   └── fixtures/               # Test scenarios (repos, files, configs)
│   ├── grading/                    # Grading criteria, rubrics, critique prompts
│   ├── monitoring/                 # Transcript analysis and metrics tooling
│   ├── capture/                    # Quick-capture — "Write That Down"
│   ├── runner/                     # Test runner with filtering
│   └── results/                    # Test outputs, historical data, regressions
│       ├── transcripts/
│       ├── metrics/
│       └── regressions/
│
├── project-a/                      # A real project (its own git repo)
├── project-b/                      # Another project (its own git repo)
└── ...
```

The key insight: `.claude/` and `CLAUDE.md` at the stuart root *are* Doreen's enforcement layer — they apply everywhere. The `doreen/` directory is her brain — the tests, grading, and tooling that evaluate and improve Stuart over time. Projects are just Stuart's workspace.

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
