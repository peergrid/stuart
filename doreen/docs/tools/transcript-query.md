# tq — Transcript Query: Analysis Tool for Claude Code Sessions

## Purpose

`tq` is a Go CLI tool for querying, navigating, and analyzing Claude Code transcripts. It treats transcript JSONL files as a queryable database with cursor-based navigation, composable filters, and built-in analysis modes that directly serve doreen's graders.

It replaces ad-hoc transcript parsing with a single, fast, comprehensive tool that both humans and Claude can use.

## Design Philosophy

1. **Time, not sessions.** You never need to know a session ID. Queries operate on a time window across all sessions: "last 2 hours", "last 3 days", "today". Sessions are an implementation detail.
2. **CWD, not flags.** Project is detected from your working directory automatically. You're already in the project — tq knows which one.
3. **Cursor, not dump.** Find a point of interest, then navigate around it. Not "dump everything and grep."
4. **Composable filters.** Every filter can combine with every other filter. Filters narrow; they never change the output format.
5. **Three audiences.** Human terminal users get readable tables. Claude gets JSON. Scripts get line-oriented output. One tool, three modes.
6. **Grader-native.** Built-in analysis commands produce exactly the data doreen's graders need. No glue code required.
7. **Fast.** Go, streaming where possible, no unnecessary parsing. A 100MB transcript should respond in under a second for targeted queries.

## Installation

```
cd doreen/tools/tq
go build -o tq .
# Move to PATH or symlink
```

## CLI Structure

`tq` uses subcommands. The general form is:

```
tq <command> [flags]
```

### Scope Model

Batch and cursor commands have separate scope flags. They are NOT shared.

**Batch commands** (`stats`, `audit`, `errors`, `tokens`, `compactions`, `tools`, `messages`, `raw`, `critique-data`, `agent-trace`, `sessions`, `agents`, `find`): `--since`/`--until` filter by time window. Default: last 24h.

**Cursor commands** (`show`, `walk`): `--from`/`--until` define start and stop anchors. `--around` positions approximately. All available transcripts are loaded — no time limit on navigation.

#### Batch flags

| Flag | Default | Description |
|------|---------|-------------|
| `--since DURATION` | `24h` | Time filter start. Accepts `30m`, `2h`, `3d`, `1w`, ISO date. |
| `--until TIME` | now | Time filter end. |

#### Cursor flags

| Flag | Default | Description |
|------|---------|-------------|
| `--from ANCHOR` | start | Start position: timestamp, or anchor spec. |
| `--until ANCHOR` | — | Stop position: same syntax as `--from`. |
| `--around DURATION` | — | Approximate starting point. "I think it was about 2 days ago." |

#### Anchor specs

Both `--from` and `--until` accept anchor specs — patterns that match a point of interest in the transcript:

| Spec | Description |
|------|-------------|
| `tool:NAME` | First use of tool NAME (e.g., `tool:Bash`, `tool:Read`) |
| `tool:NAME:PATTERN` | First use of tool NAME where input matches PATTERN (e.g., `tool:Bash:jq`) |
| `error` | First error record (tool error, permission denial, hook failure) |
| `compaction` | First compaction boundary |
| `user` | First external human message (not tool results) |
| `assistant` | First assistant record |
| `pattern:REGEX` | First record whose text matches REGEX |

Anchors can be combined with `--around` to search from an approximate time: `--from "tool:Bash:jq" --around 4d` means "find the first Bash call containing jq, starting from about 4 days ago."

#### Common flags (all commands)

| Flag | Default | Description |
|------|---------|-------------|
| `--project NAME` | from CWD | Override project detection |
| `--session UUID` | (all) | Filter to one session (rarely needed) |
| `--all` | false | All projects |
| `--subagents` | false | Include subagent transcripts |

**Project detection:** tq maps the current working directory to a Claude project by converting the absolute path to a slug (`/home/user/stuart` → `-home-user-stuart`) and looking it up in `~/.claude/projects/`. It walks up parent directories if the exact CWD doesn't match.

**Batch time window:** All session files whose modification time falls within the window are included. Records are merged chronologically across sessions. Default window is 24 hours.

**Cursor loading:** All available transcripts are loaded. `--around` positions the search start. `--from` and `--until` define the walk range via anchors.

**Examples:**
```
# Batch — time as filter
tq errors                       # Errors in the last 24h
tq errors --since 2d            # Errors in the last 2 days
tq stats --since 1w             # Stats for the whole week

# Cursor — anchors define range, walk between them
tq walk --reverse --until user --role assistant --count  # Assistant turns since last human message
tq walk --from "tool:Bash:jq" --around 4d --until "tool:Read" --role user --contains "poop" --count
tq walk --around 3h                                      # Walk from ~3 hours ago
tq walk --around 2d --reverse                            # Walk backward from ~2 days ago
```

## Commands

### Discovery Commands

#### `tq sessions`

List sessions in the time window with metadata. Useful for context, not required for querying.

```
tq sessions                     # Sessions in last 24h
tq sessions --since 1w          # Sessions in the last week
tq sessions --since 1w --json
```

Output columns: SESSION_ID, SIZE, RECORDS, MODEL, AGENTS, START, DURATION.

#### `tq agents`

List subagents across sessions in the time window.

```
tq agents                       # All agents, last 24h
tq agents --since 2d            # Agents from last 2 days
tq agents --since 2d --json
```

Output columns: AGENT_ID, SESSION, NAME, MODEL, SIZE, DURATION, TOOLS, LAUNCH_PROMPT.

#### `tq find`

Search for sessions or agents by content.

```
tq find "observer"              # Find agents by name/prompt
tq find --in-messages "fix the bug" --since 3d  # Find sessions with text
```

### Navigation Commands

#### `tq show`

Show records around a point of interest. This is the core navigation command. All transcripts are loaded — `--since` positions the search start, it does not limit what you can see.

```
# Show the first error with 5 turns of context (searches all history)
tq show --first error --context 5

# "I think it was about 2 days ago" — start searching near there
tq show --around 2d --first error --context 5

# Show a specific timestamp with context
tq show --at "2026-03-10T14:22:00" --context 5

# Show the 3rd compaction ever (no time limit)
tq show --nth 3 compaction --context 10
```

**Anchor points** (what `--first`, `--last`, `--nth N` can target):
- `compaction` — compaction boundary events
- `error` — tool errors, permission denials, hook failures
- `agent-launch` — where a subagent was launched
- `agent-return` — where a subagent's result was received
- `tool:NAME` — first/last/nth use of a specific tool (e.g., `tool:Bash`)
- `pattern:REGEX` — first/last/nth match of a regex in message text

#### `tq walk`

Step through turns with optional filters. All transcripts are loaded — cursor commands have no time boundary.

`--from` and `--until` accept anchor specs (see Scope Model above) to define the walk range. `--around` positions the search start approximately. All filters (`--role`, `--contains`, `--tool`, etc.) narrow what's counted/shown between the anchors.

```
# Start walking from ~1 hour ago, tools only
tq walk --around 1h --tools-only

# Walk backward from now, last 20 turns
tq walk --reverse --limit 20

# Walk from a specific timestamp
tq walk --from "2026-03-10T14:00:00"

# Walk only operator messages and responses
tq walk --external-only

# Count assistant turns since last human message
tq walk --reverse --until user --role assistant --count

# Check if warning text exists since last human message (exit code only)
tq walk --reverse --until user --contains "warning text" --exists

# Count user prompts containing "poop" between two tool uses
tq walk --from "tool:Bash:jq" --around 4d --until "tool:Read:dangerouslyDisableSandbox" --role user --contains "poop" --count

# Get the text of all assistant messages in a range
tq walk --from "tool:Bash:deploy" --until error --role assistant --field text
```

Output per turn:
```
--- Turn 42 | session:2ac9 | assistant | 2026-03-10 14:22:05 ---
[Text] I'll read the file first to understand the current structure.
[Tool: Read] /home/user/project/main.go
[Tool: Grep] pattern="TODO" path=/home/user/project/
```

### Query Commands

#### `tq messages`

Extract messages with filtering.

```
tq messages --external-only             # Operator messages, last 24h
tq messages --role assistant --since 1h # Assistant text, last hour
tq messages --contains "error|fail"     # Pattern match
```

#### `tq tools`

Extract tool calls with filtering.

```
tq tools                                # All tool calls, last 24h
tq tools --tool Bash                    # Only Bash calls
tq tools --tool Bash --input-contains "grep"  # Anti-pattern detection
tq tools --audit "input.dangerouslyDisableSandbox=true"
tq tools --with-results                 # Include tool output
```

#### `tq raw`

Dump raw JSONL records, optionally filtered. For piping into jq or other tools.

```
tq raw                                  # All records, last 24h
tq raw --type assistant                 # Only assistant records
tq raw --type assistant | jq '.message.usage.output_tokens'
```

### Analysis Commands

#### `tq stats`

Summary statistics across the time window.

```
tq stats                        # Stats for last 24h
tq stats --since 3d             # Stats for last 3 days
tq stats --json
```

Output:
```
Period:         last 24h (3 sessions)
Duration:       4h 12m total active time
Records:        2,847

Messages:       user: 324, assistant: 918, system: 15
Tool calls:     1,201
  Read:         289
  Edit:         145
  Bash:         267
  Grep:         134
  ...

Tokens:
  Input:        3,234,567
  Output:       189,012
  Cache read:   2,987,654

Compactions:    5
Errors:         17
Subagents:      12
```

#### `tq tokens`

Token consumption analysis. Directly serves the token-cost-metrics grader.

```
tq tokens                       # Turn-by-turn timeline, last 24h
tq tokens --since 1h --summary  # Summary only
tq tokens --json                # JSON for grader
```

Output columns: TURN, SESSION, TIME, CONTEXT, INPUT, OUTPUT, CACHE_READ, CUMULATIVE_INPUT, TOOLS.

#### `tq compactions`

Detect and analyze compaction events.

```
tq compactions                          # Compactions in last 24h
tq compactions --since 1w --with-context 5
tq compactions --json
```

Detection methods (ordered by reliability):
1. `compact_boundary` system records (direct CC marker)
2. `POST-COMPACTION RECOVERY` markers (hook-injected)
3. Total context token drops >50% from >50K baseline (heuristic)

#### `tq errors`

Extract all errors: tool errors, permission denials, hook failures, runtime errors.

```
tq errors                               # Errors in last 24h
tq errors --since 3d --type tool_error  # Only tool errors
tq errors --with-context 3             # Show surrounding turns
tq errors --json
```

Error types: `tool_error`, `hook_error`, `permission_denied`, `runtime_error`.

#### `tq audit`

Tool use audit analysis. Directly serves the tool-use-audit grader.

```
tq audit                        # Audit last 24h
tq audit --since 3d --json
```

Checks performed:
- **Dedicated tool preference**: Flags `cat`/`grep`/`find`/`sed`/`awk`/`echo` via Bash when Read/Grep/Glob/Edit/Write exist.
- **Read-before-edit**: Every Edit/Write to an existing file was preceded by a Read of that file.
- **No redundant reads**: Same file not read multiple times without an intervening edit.
- **Parallel where possible**: Measures parallelism ratio (parallel tool calls / total tool call turns).
- **Agent delegation**: For sessions with 10+ exploratory tool calls, flags lack of agent delegation.

Output:
```
Tool Use Audit — last 24h (3 sessions)

Violations: 12 / 401 tool calls (3.0% violation rate)

Dedicated tool preference:
  [session:2ac9 Turn 15] Bash: grep -r "TODO" . → should use Grep tool
  [session:2ac9 Turn 23] Bash: cat main.go → should use Read tool
  [session:f8b1 Turn 67] Bash: find . -name "*.go" → should use Glob tool

Read-before-edit:
  [session:2ac9 Turn 31] Edit: main.go — no preceding Read

Redundant reads:
  [session:f8b1 Turn 44, 89] Read: config.yaml — read twice, no edit between

Parallelism: 23 / 67 multi-tool turns used parallel calls (34.3%)

Score: 0.970 (12 violations / 401 calls)
```

#### `tq critique-data`

Extract the data needed for the transcript-critique LLM grader.

```
tq critique-data --json
tq critique-data --since 2d --json
```

Output includes:
- Original task/prompt (first external user message per session)
- Decision points (turns where Claude chose an approach)
- Error recovery sequences (error + subsequent turns)
- Communication quality signals (response lengths, verbosity indicators)
- Completion claims (turns where Claude says "done"/"complete"/"finished") with what followed
- Context management signals (compactions, repeated reads, lost information indicators)

#### `tq agent-trace`

Trace subagent lifecycle: launch, execution, and return.

```
tq agent-trace                          # All agents, last 24h
tq agent-trace --agent "observer"       # Filter by agent name
tq agent-trace --since 3d --json
```

Output per agent:
```
Agent: observer (session:2ac9, sub-001)
  Launched: Turn 45, 14:22:05
  Prompt: "Investigate the test failures in..."  (142 tokens)
  Duration: 3m 12s
  Tool calls: 23 (Read: 8, Grep: 6, Bash: 5, Edit: 4)
  Tokens: input 45,678 / output 3,456
  Returned: Turn 48, 14:25:17
  Result: "Found 3 failing tests..." (89 tokens)
```

### Output Modes

Every command supports three output modes, plus three output modifiers:

| Flag | Mode | Description |
|------|------|-------------|
| (default) | Human | Formatted tables and text for terminal |
| `--json` | JSON | Structured JSON for piping to other tools |
| `--jsonl` | JSONL | One JSON object per line for streaming |

| Flag | Modifier | Description |
|------|----------|-------------|
| `--count` | Count | Output only the count of matching records (a single integer) |
| `--exists` | Exists | No output. Exit 0 if any match, exit 1 if none. Short-circuits on first match. |
| `--field PATH` | Project | Output a single field from each record, one per line. Dot-notation (e.g., `input.command`). |

`--count` and `--exists` are designed for hooks and scripts. `--field` eliminates jq for simple field extraction.

### Common Filters

These flags work with any command that processes records:

| Flag | Description |
|------|-------------|
| `--role ROLE` | Filter by message role: `user`, `assistant`, `system`, `all` |
| `--type TYPE` | Filter by record type |
| `--tool NAME` | Filter for records containing tool_use of NAME |
| `--contains PATTERN` | Regex match against message text content |
| `--audit FIELD=REGEX` | Match any JSON field (dot-notation) against regex |
| `--external-only` | Only genuine operator messages (skip meta, system, tool results) |
| `--no-sidechain` | Skip sidechain records |
| `--limit N` | Max records to output |
| `--no-truncate` | Show full content (default truncates long values) |

### Cursor Model

The cursor model is how `tq show` and `tq walk` navigate transcripts. **All available transcripts are loaded** — cursor commands have no time boundary. Records are merged chronologically across sessions and assigned sequential **turn numbers** (1-indexed within the merged timeline).

**Anchoring:** `--from` and `--until` accept anchor specs to define the range of interest. `--around` provides an approximate search start. Anchors can match tool uses, errors, compactions, roles, or text patterns (see Scope Model above).

**Walking:** `tq walk` steps through turns between `--from` and `--until` anchors. `--reverse` walks backward. `--limit N` caps output. Combined with `--count`, `--exists`, `--field`, and filters, walk answers questions like "how many assistant turns since the last user message?" or "did this text appear in recent turns?" — without jq.

## How tq Serves Doreen's Graders

### tool-use-audit grader

Primary command: `tq audit --json`

The audit command runs all tool-use checks and produces structured violations with turn references. The grader reads the JSON output directly.

### transcript-critique grader

Primary command: `tq critique-data --json`

Extracts the structured summary an LLM grader needs: original prompt, decision points, error recovery, completion claims. The grader feeds this to a separate Claude instance as context for the critique.

### token-cost-metrics grader

Primary command: `tq tokens --json`

Provides turn-by-turn token data, compaction count, context utilization curve, and aggregate metrics.

### regression-comparison grader

Does not directly use tq. Operates on stored grader outputs. However, `tq stats --json` provides the raw session metrics that get stored for comparison.

### work-product-critique grader

Primary command: `tq tools --tool Edit --json` and `tq tools --tool Write --json`

Extracts the set of file modifications for work product analysis. Combined with `tq messages --external-only --limit 1 --json` for the original prompt.

## JSONL Record Format

See `doreen/docs/tools/jsonl-format.md` for the complete JSONL format reference.

## Implementation Notes

### Package Structure

```
doreen/tools/tq/
  main.go              # Entry point
  go.mod
  cmd/
    root.go            # Root command, global flags
    sessions.go        # tq sessions
    agents.go          # tq agents
    find.go            # tq find
    show.go            # tq show
    walk.go            # tq walk
    messages.go        # tq messages
    tools.go           # tq tools
    raw.go             # tq raw
    stats.go           # tq stats
    tokens.go          # tq tokens
    compactions.go     # tq compactions
    errors.go          # tq errors
    audit.go           # tq audit
    critique_data.go   # tq critique-data
    agent_trace.go     # tq agent-trace
  internal/
    transcript/
      record.go        # Record types and parsing
      loader.go        # JSONL file loading (streaming + list modes)
      filter.go        # Filter chain implementation
      cursor.go        # Cursor/navigation model
      session.go       # Project detection, time-windowed session discovery
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

- Stream records by default; only load full list when random access is needed
- Use buffered I/O for all file reads
- Parse only the fields needed for the current operation (lazy parsing)
- Session metadata can be extracted from first/last lines without full file read
- Time-window filtering at the file level (mtime check) before opening files

### Error Handling

- Missing files: warn on stderr, continue with remaining files
- Malformed JSON lines: skip silently (transcripts can have partial writes)
- No matching records: print "No results" to stderr, exit 0
- Invalid flags: print usage help, exit 1
