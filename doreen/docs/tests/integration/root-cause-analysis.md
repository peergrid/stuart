# Integration Test: Root Cause Analysis Before Fix

## Behavior Under Test

When given a bug report, Claude MUST investigate the root cause before proposing a fix — not just patch the symptom described in the report.

## Fixture Setup

A project with a bug: "Users see a 500 error when submitting the contact form." The actual root cause is a missing null check three layers deep in a service that the form handler calls — not in the form handler itself.

- `src/routes/contact.py` — the form handler (looks fine)
- `src/services/email.py` — calls a helper (looks fine)
- `src/helpers/template.py` — the actual bug: `template.render(data['name'])` where `data['name']` can be `None`

## Prompt

"Users are getting 500 errors when they submit the contact form. Please fix it."

## Expected Behavior

- Claude traces the execution path from the form handler through the service layer.
- Claude identifies the null reference in `template.py` as the root cause.
- The fix is applied in `template.py`, not as a defensive check in the route handler.

## Failure Modes

- Adding a try/except in the route handler that swallows the error.
- Adding null checks at every layer instead of fixing the actual problem.
- Fixing the route handler without ever reading the files it calls.
- Guessing the fix without tracing the execution path.
