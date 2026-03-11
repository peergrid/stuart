# Integration Test: API/SDK Usage with Docs Available

## Behavior Under Test

When domain expertise docs (SDK docs, API references) are available, Claude MUST use them to produce correct API calls — not rely on training data which may be outdated.

## Fixture Setup

A project that integrates with a fictional or versioned API. Provide:
- SDK docs (via context7, local files, or CLAUDE.md) describing the current API surface.
- The docs describe an API that has changed from v1 to v2 — method names renamed, parameters reordered, a required field added.
- Existing code using the v1 API.

## Prompt

"Update the payment integration to use the v2 API."

## Expected Behavior

- Claude consults the provided SDK docs before making changes.
- The updated code uses v2 method names, parameter order, and required fields.
- No v1 patterns remain unless explicitly documented as backwards-compatible.

## Failure Modes

- Using v1 method names or signatures from training data.
- Ignoring the provided docs and guessing the v2 API shape.
- Mixing v1 and v2 patterns in the same integration.
- Not reading the docs at all.
