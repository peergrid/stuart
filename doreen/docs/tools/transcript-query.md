# dq — Doreen Query: Transcript Analysis Tool

## Purpose

`dq` is a Go CLI tool for querying, navigating, and analyzing Claude Code transcripts. It treats transcript JSONL files as a queryable database with cursor-based navigation, composable filters, and built-in analysis modes that directly serve doreen's graders.

It replaces ad-hoc transcript parsing with a single, fast, comprehensive tool that both humans and Claude can use.

## Design Philosophy

1. **Cursor, not dump.** The primary interaction model is: find a point of interest, then navigate around it. Not "dump everything and grep."
2. **Composable filters.** Every filter can combine with every other filter. Filters narrow; they never change the output format.
3. **Three audiences.** Human terminal users get readable tables. Claude gets JSON. Scripts get line-oriented output. One tool, three modes.
4. **Grader-native.** Built-in analysis commands produce exactly the data doreen's graders need. No glue code required.
5. **Fast.** Go, streaming where possible, no unnecessary parsing. A 100MB transcript should respond in under a second for targeted queries.

## Installation

```
cd doreen/tools/dq
go build -o dq .
# Move to PATH or symlink
```

## CLI Structure

`dq` uses subcommands. The general form is:

```
dq <command> [flags] [session-specifier]
```

### Session Specifiers

Every command that operates on transcript data accepts a session specifier:

| Specifier | Example | Description |
|-----------|---------|-------------|
| File path | `/path/to/session.jsonl` | Direct file path |
| `--project NAME` | `--project stuart` | All sessions in a project |
| `--session UUID` | `--session 2ac9` | Session by UUID prefix (requires --project) |
| `--latest` | `--latest` | Most recent session in current project |
| `--all` | `--all` | All projects |

The `--subagents` flag includes subagent transcripts alongside the root session.

## Commands

### Discovery Commands

#### `dq sessions`

List sessions with metadata. The entry point for exploration.

```
dq sessions --project stuart
dq sessions --project stuart --limit 5
dq sessions --project stuart --json
```

Output columns: SESSION_ID, SIZE, RECORDS, MODEL, AGENTS, START, DURATION.

#### `dq agents`

List subagents for a session.

```
dq agents --project stuart --session 2ac9
dq agents --project stuart --latest
dq agents --project stuart --latest --json
```

Output columns: AGENT_ID, NAME, MODEL, SIZE, DURATION, TOOLS, LAUNCH_PROMPT.

#### `dq find`

Search for sessions or agents by content.

```
dq find --project stuart "observer"           # Find agents by name/prompt
dq find --project stuart --in-messages "fix the bug"  # Find sessions containing text
```

### Navigation Commands

#### `dq show`

Show records around a point of interest. This is the core navigation command.

```
# Show turn 42 with 3 turns of context before and after
dq show --project stuart --session 2ac9 --turn 42 --context 3

# Show the record at a specific timestamp
dq show --session 2ac9 --project stuart --at "2026-03-10T14:22:00" --context 5

# Show the first compaction event with surrounding context
dq show --project stuart --session 2ac9 --first compaction --context 10

# Show the 3rd error with context
dq show --project stuart --session 2ac9 --nth 3 error --context 5
```

**Anchor points** (what `--first`, `--last`, `--nth N` can target):
- `compaction` — compaction boundary events
- `error` — tool errors, permission denials, hook failures
- `agent-launch` — where a subagent was launched
- `agent-return` — where a subagent's result was received
- `tool:NAME` — first/last/nth use of a specific tool (e.g., `tool:Bash`)
- `pattern:REGEX` — first/last/nth match of a regex in message text

#### `dq walk`

Step through a transcript turn by turn with optional filters. Interactive-style output for understanding flow.

```
# Walk through all turns, showing role and summary
dq walk --project stuart --session 2ac9

# Walk through only tool calls
dq walk --project stuart --session 2ac9 --tools-only

# Walk through only operator messages and Claude's direct responses
dq walk --project stuart --session 2ac9 --external-only

# Walk starting from turn 50
dq walk --project stuart --session 2ac9 --from 50

# Walk backward from the end
dq walk --project stuart --session 2ac9 --reverse --limit 20
```

Output per turn:
```
--- Turn 42 | assistant | 2026-03-10 14:22:05 | 892 out tokens ---
[Text] I'll read the file first to understand the current structure.
[Tool: Read] /home/user/project/main.go
[Tool: Grep] pattern="TODO" path=/home/user/project/
```

### Query Commands

#### `dq messages`

Extract messages with filtering.

```
# All operator messages (human input only)
dq messages --latest --external-only

# All assistant text (no tool calls)
dq messages --latest --role assistant

# Messages containing a pattern
dq messages --latest --contains "error|fail"

# Messages in a time window
dq messages --latest --since "2026-03-10T14:00:00" --until "2026-03-10T15:00:00"
```

#### `dq tools`

Extract tool calls with filtering.

```
# All tool calls in latest session
dq tools --latest

# Only Bash calls
dq tools --latest --tool Bash

# Bash calls that used grep (anti-pattern detection)
dq tools --latest --tool Bash --input-contains "grep"

# Tool calls where input matches a field pattern
dq tools --latest --audit "input.dangerouslyDisableSandbox=true"

# Tool calls with their results
dq tools --latest --with-results
```

#### `dq raw`

Dump raw JSONL records, optionally filtered. For piping into jq or other tools.

```
# All records as JSONL
dq raw --latest

# Only assistant records as JSONL
dq raw --latest --type assistant

# Pipe to jq for custom analysis
dq raw --latest --type assistant | jq '.message.usage.output_tokens'
```

### Analysis Commands

#### `dq stats`

Session summary statistics.

```
dq stats --latest
dq stats --project stuart --session 2ac9
dq stats --project stuart  # Aggregate across all sessions
```

Output:
```
Session:        2ac9f4e1-...
Duration:       1h 23m 45s
Records:        847
Model:          claude-sonnet-4-20250514

Messages:       user: 124, assistant: 318, system: 5
Tool calls:     401
  Read:         89
  Edit:         45
  Bash:         67
  Grep:         34
  ...

Tokens:
  Input:        1,234,567
  Output:       89,012
  Cache read:   987,654

Compactions:    2
Errors:         7
Subagents:      3
```

#### `dq tokens`

Token consumption analysis. Directly serves the token-cost-metrics grader.

```
# Turn-by-turn token timeline
dq tokens --latest

# Token summary only
dq tokens --latest --summary

# JSON output for grader consumption
dq tokens --latest --json
```

Output columns: TURN, TIME, CONTEXT, INPUT, OUTPUT, CACHE_READ, CUMULATIVE_INPUT, TOOLS.

Summary includes: total tokens, max single-turn context, context utilization curve data points, compaction count, tokens-per-turn average.

#### `dq compactions`

Detect and analyze compaction events. Shows what happened before and after each compaction.

```
dq compactions --latest
dq compactions --latest --with-context 5   # Show 5 turns around each compaction
dq compactions --latest --json
```

Detection methods (ordered by reliability):
1. `compact_boundary` system records (direct CC marker)
2. `POST-COMPACTION RECOVERY` markers (hook-injected)
3. Total context token drops >50% from >50K baseline (heuristic)

#### `dq errors`

Extract all errors: tool errors, permission denials, hook failures, runtime errors.

```
dq errors --latest
dq errors --latest --type tool_error       # Only tool errors
dq errors --latest --with-context 3        # Show surrounding turns
dq errors --latest --json
```

Error types: `tool_error`, `hook_error`, `permission_denied`, `runtime_error`.

#### `dq audit`

Tool use audit analysis. Directly serves the tool-use-audit grader.

```
dq audit --latest
dq audit --latest --json
```

Checks performed:
- **Dedicated tool preference**: Flags `cat`/`grep`/`find`/`sed`/`awk`/`echo` via Bash when Read/Grep/Glob/Edit/Write exist.
- **Read-before-edit**: Every Edit/Write to an existing file was preceded by a Read of that file.
- **No redundant reads**: Same file not read multiple times without an intervening edit.
- **Parallel where possible**: Measures parallelism ratio (parallel tool calls / total tool call turns).
- **Agent delegation**: For sessions with 10+ exploratory tool calls, flags lack of agent delegation.

Output:
```
Tool Use Audit — Session 2ac9

Violations: 12 / 401 tool calls (3.0% violation rate)

Dedicated tool preference:
  [Turn 15] Bash: grep -r "TODO" . → should use Grep tool
  [Turn 23] Bash: cat main.go → should use Read tool
  [Turn 67] Bash: find . -name "*.go" → should use Glob tool

Read-before-edit:
  [Turn 31] Edit: main.go — no preceding Read

Redundant reads:
  [Turn 44, 89] Read: config.yaml — read twice, no edit between

Parallelism: 23 / 67 multi-tool turns used parallel calls (34.3%)

Score: 0.970 (12 violations / 401 calls)
```

#### `dq critique-data`

Extract the data needed for the transcript-critique LLM grader. Produces a structured summary optimized for LLM consumption.

```
dq critique-data --latest --json
```

Output includes:
- Original task/prompt (first external user message)
- Decision points (turns where Claude chose an approach)
- Error recovery sequences (error + subsequent turns)
- Communication quality signals (response lengths, verbosity indicators)
- Completion claims (turns where Claude says "done"/"complete"/"finished") with what followed
- Context management signals (compactions, repeated reads, lost information indicators)

#### `dq agent-trace`

Trace subagent lifecycle: launch, execution, and return.

```
dq agent-trace --latest
dq agent-trace --latest --agent "observer"
dq agent-trace --latest --json
```

Output per agent:
```
Agent: observer (2ac9-sub-001)
  Launched: Turn 45, 14:22:05
  Prompt: "Investigate the test failures in..."  (142 tokens)
  Duration: 3m 12s
  Tool calls: 23 (Read: 8, Grep: 6, Bash: 5, Edit: 4)
  Tokens: input 45,678 / output 3,456
  Returned: Turn 48, 14:25:17
  Result: "Found 3 failing tests..." (89 tokens)
```

### Output Modes

Every command supports three output modes:

| Flag | Mode | Description |
|------|------|-------------|
| (default) | Human | Formatted tables and text for terminal |
| `--json` | JSON | Structured JSON for piping to other tools |
| `--jsonl` | JSONL | One JSON object per line for streaming |

### Global Filters

These flags work with any command that processes records:

| Flag | Description |
|------|-------------|
| `--role ROLE` | Filter by message role: `user`, `assistant`, `system`, `all` |
| `--type TYPE` | Filter by record type |
| `--tool NAME` | Filter for records containing tool_use of NAME |
| `--contains PATTERN` | Regex match against message text content |
| `--audit FIELD=REGEX` | Match any JSON field (dot-notation) against regex |
| `--since ISO_DATE` | Records after this timestamp |
| `--until ISO_DATE` | Records before this timestamp |
| `--external-only` | Only genuine operator messages (skip meta, system, tool results) |
| `--no-sidechain` | Skip sidechain records |
| `--limit N` | Max records to output |
| `--no-truncate` | Show full content (default truncates long values) |

### Cursor Model

The cursor model is how `dq show` and `dq walk` navigate transcripts. Internally, every record is assigned a sequential **turn number** (1-indexed) within its session. Turn numbers provide stable references for navigation.

**Anchoring:** You set a cursor position by turn number (`--turn N`), timestamp (`--at TS`), or semantic anchor (`--first`/`--last`/`--nth N` target).

**Context:** `--context N` shows N turns before and after the anchor. The context window respects filters: if you filter to `--tools-only`, context shows only tool call turns within the window.

**Walking:** `dq walk` advances the cursor one turn at a time, printing each turn. `--from N` sets the starting turn. `--reverse` walks backward. `--limit N` stops after N turns.

## How dq Serves Doreen's Graders

### tool-use-audit grader

Primary command: `dq audit --latest --json`

The audit command runs all tool-use checks and produces structured violations with turn references. The grader reads the JSON output directly. No additional processing needed.

### transcript-critique grader

Primary command: `dq critique-data --latest --json`

Extracts the structured summary an LLM grader needs: original prompt, decision points, error recovery, completion claims. The grader feeds this to a separate Claude instance as context for the critique.

### token-cost-metrics grader

Primary command: `dq tokens --latest --json`

Provides turn-by-turn token data, compaction count, context utilization curve, and aggregate metrics. The grader computes derived metrics (tokens per line of diff, efficiency scores) from this data.

### regression-comparison grader

Does not directly use dq. Operates on stored grader outputs. However, `dq stats --json` provides the raw session metrics that get stored for comparison.

### work-product-critique grader

Primary command: `dq tools --latest --tool Edit --json` and `dq tools --latest --tool Write --json`

Extracts the set of file modifications for work product analysis. Combined with `dq messages --latest --external-only --limit 1 --json` for the original prompt.

## JSONL Record Format

See `doreen/docs/tools/jsonl-format.md` for the complete JSONL format reference.

## Implementation Notes

### Package Structure

```
doreen/tools/dq/
  main.go              # Entry point, cobra command tree
  go.mod
  go.sum
  cmd/
    root.go            # Root command, global flags
    sessions.go        # dq sessions
    agents.go          # dq agents
    find.go            # dq find
    show.go            # dq show
    walk.go            # dq walk
    messages.go        # dq messages
    tools.go           # dq tools
    raw.go             # dq raw
    stats.go           # dq stats
    tokens.go          # dq tokens
    compactions.go     # dq compactions
    errors.go          # dq errors
    audit.go           # dq audit
    critique_data.go   # dq critique-data
    agent_trace.go     # dq agent-trace
  internal/
    transcript/
      record.go        # Record types and parsing
      loader.go        # JSONL file loading (streaming + list modes)
      filter.go        # Filter chain implementation
      cursor.go        # Cursor/navigation model
      session.go       # Session discovery and resolution
    format/
      human.go         # Human-readable output formatting
      json.go          # JSON output
      jsonl.go         # JSONL streaming output
    analysis/
      audit.go         # Tool use audit logic
      tokens.go        # Token consumption analysis
      compaction.go    # Compaction detection
      errors.go        # Error extraction
      agents.go        # Agent lifecycle tracing
      critique.go      # Critique data extraction
```

### Performance

- Stream records by default; only load full list when random access is needed (compaction detection with heuristic, cursor context)
- Use buffered I/O for all file reads
- Parse only the fields needed for the current operation (lazy parsing)
- Session metadata can be extracted from first/last lines without full file read

### Error Handling

- Missing files: warn on stderr, continue with remaining files
- Malformed JSON lines: skip silently (transcripts can have partial writes)
- No matching records: print "No results" to stderr, exit 0
- Invalid flags: print usage help, exit 1
