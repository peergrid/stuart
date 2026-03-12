# Claude Code Transcript JSONL Format

Reference documentation for the JSONL format produced by Claude Code sessions. Each line in a `.jsonl` transcript file is one JSON record.

## File Layout

```
~/.claude/projects/<project-slug>/
  <session-uuid>.jsonl                    # Root session transcript
  <session-uuid>/subagents/
    <agent-uuid>.jsonl                    # Subagent transcripts
```

The project slug is derived from the absolute project path with `/` replaced by `-` and the leading slash dropped. For example, `/home/cvk/stuart` becomes `-home-cvk-stuart`.

## Record Structure

Every record has these top-level fields:

| Field        | Type   | Description |
|-------------|--------|-------------|
| `type`      | string | Record type: `user`, `assistant`, `system` |
| `timestamp` | string | ISO 8601 timestamp (e.g., `2026-03-10T14:22:03.456Z`) |
| `sessionId` | string | UUID of the session this record belongs to |
| `message`   | object | The message payload (structure varies by type) |
| `isSidechain` | bool | If true, this record is from a sidechain (background) operation |
| `isMeta`    | bool   | If true, this is internal metadata, not a real conversation turn |

## Record Types

### `user` Records

Represent input to Claude: operator messages, tool results, and system injections.

```json
{
  "type": "user",
  "timestamp": "2026-03-10T14:22:03.456Z",
  "sessionId": "abc123...",
  "message": {
    "role": "user",
    "content": "..."  // string or array of content blocks
  }
}
```

**Content formats:**

Plain text (operator message):
```json
{ "role": "user", "content": "Fix the bug in parser.go" }
```

Content block array (tool results, mixed content):
```json
{
  "role": "user",
  "content": [
    {
      "type": "tool_result",
      "tool_use_id": "toolu_abc123",
      "content": "file contents here...",
      "is_error": false
    },
    {
      "type": "tool_result",
      "tool_use_id": "toolu_def456",
      "content": [
        { "type": "text", "text": "Error: file not found" }
      ],
      "is_error": true
    }
  ]
}
```

Tool result content can be a string or an array of text blocks. The `is_error` field indicates whether the tool call failed.

### `assistant` Records

Represent Claude's responses: text output and tool invocations.

```json
{
  "type": "assistant",
  "timestamp": "2026-03-10T14:22:05.789Z",
  "sessionId": "abc123...",
  "message": {
    "role": "assistant",
    "model": "claude-sonnet-4-20250514",
    "content": [...],  // array of content blocks
    "usage": {
      "input_tokens": 15234,
      "output_tokens": 892,
      "cache_read_input_tokens": 12000,
      "cache_creation_input_tokens": 3000
    },
    "stop_reason": "end_turn"  // or "tool_use"
  }
}
```

**Content blocks in assistant messages:**

Text block:
```json
{ "type": "text", "text": "I'll read the file first." }
```

Tool use block:
```json
{
  "type": "tool_use",
  "id": "toolu_abc123",
  "name": "Read",
  "input": {
    "file_path": "/home/user/project/main.go"
  }
}
```

Common tool names: `Read`, `Write`, `Edit`, `Bash`, `Grep`, `Glob`, `Skill`, `Agent`, `TodoWrite`, `NotebookEdit`, `WebFetch`, `WebSearch`.

**Usage fields:**

| Field | Description |
|-------|-------------|
| `input_tokens` | Non-cached input tokens for this turn |
| `output_tokens` | Output tokens generated |
| `cache_read_input_tokens` | Tokens read from prompt cache |
| `cache_creation_input_tokens` | Tokens written to prompt cache |

Total context window consumption for a turn = `input_tokens + cache_read_input_tokens + cache_creation_input_tokens`.

### `system` Records

System-level events. The `subtype` field distinguishes them.

**Compaction boundary** (`subtype: "compact_boundary"`):
```json
{
  "type": "system",
  "subtype": "compact_boundary",
  "timestamp": "2026-03-10T15:00:00.000Z",
  "compactMetadata": {
    "trigger": "auto",      // "auto" or "manual"
    "preTokens": 180000     // context size before compaction
  }
}
```

This marks the exact point where conversation history was summarized/compacted. Everything before this point in the transcript was condensed.

## Content Patterns to Detect

### Operator vs. Internal Messages

Not all `user` records are from the human operator. To identify genuine operator input:
- Skip records where `isMeta` is true
- Skip records where content starts with `<` (system-reminder injections)
- Skip records where content starts with `/` (slash commands)
- Skip records that are pure `tool_result` arrays (these are tool return values)

### Hook Injections

Hooks inject content into `user` records. Look for patterns like `hook_errors`, `exit_code`, `hook failed` in the text content.

### Permission Denials

Tool permission denials appear in `user` records with content containing "Permission" and "denied" or "rejected".

### Post-Compaction Recovery

After a compaction, hooks may inject recovery markers. Look for `POST-COMPACTION RECOVERY` or `conversation was summarized` in user message text.

### Subagent Identification

Subagent launch prompts (first `user` record in a subagent transcript) may contain:
- `teammate_id="<name>"` — the agent's assigned name
- `"name": "<name>"` — alternative naming pattern from Task tool

## Key Relationships

- Each `tool_use` block in an assistant message has an `id` field
- The corresponding `tool_result` block in the next user message references it via `tool_use_id`
- Multiple tool calls can be made in a single assistant turn (parallel tool use)
- The `stop_reason` field: `"tool_use"` means Claude wants to call tools; `"end_turn"` means Claude is done speaking

## Session Discovery

To find transcripts for a given project:

1. Compute the project slug from its absolute path
2. Look in `~/.claude/projects/<slug>/` for `*.jsonl` files
3. Sort by modification time for chronological ordering
4. For subagents, check `<slug>/<session-uuid>/subagents/`

The "latest" session is the most recently modified `.jsonl` file in the project directory.
