# Draftspec Hotfix Prompt

You are implementing an emergency fix outside the standard planning phase chain.

## Goal

Write a minimal hotfix spec, implement the fix, verify it inline, and prepare for archive — without running inspect, plan, or tasks phases.

## When to Use

Use only when:
- The fix is well-understood and touches ≤ 3 files
- The root cause is already identified
- A full phase cycle would add friction without safety benefit

If scope is unclear, root cause unknown, or fix is cross-cutting — stop and use the standard workflow instead.

## Phase Contract

Inputs: `.draftspec/constitution.summary.md` (or `.draftspec/constitution.md`), user description of the fix.
Outputs: `.draftspec/specs/<slug>/hotfix.md`, implementation code.
Stop if: root cause unclear, fix exceeds 3 files, or constitutional conflict detected.

## Load First

Always read these before writing any code:

- `.draftspec/constitution.summary.md` if present; otherwise `.draftspec/constitution.md`
- only the files directly involved in the fix

## Do Not Read By Default

- full spec history or plan packages
- implementation files not listed in `Touches`
- script source files

## Stop Conditions

Stop and switch to the standard workflow if:

- root cause is unclear
- fix requires changing more than 3 files
- fix touches an API contract, data migration, or auth boundary without a clear rollback
- constitutional conflict is present
- scope requires inventing tasks beyond the stated fix

## Hotfix Spec

Write `.draftspec/specs/<slug>/hotfix.md` before touching any code:

```
---
slug: <slug>
type: hotfix
created_at: <date>
---

## Fix
<what is broken and what the change does — one or two sentences>

## Root Cause
<why it broke — one sentence>

## Risk
<what could break from this fix — one sentence; "none" only if genuinely safe>

## Verification
<concrete observable check — command output, HTTP response, or UI behavior>

## Touches
<file, file>
```

## Invariants

- Write the hotfix spec before any code change.
- Touch only files listed in `Touches`.
- Keep the fix minimal — no refactoring, no scope beyond the stated fix.
- If a file outside `Touches` must change, stop and explain before continuing.
- Log non-obvious assumptions as `[ASSUMPTION: ...]` before acting on them.
- Do not re-plan or re-design the fix silently; if the stated fix turns out to be wrong, stop and ask.

## Output expectations

- Write `.draftspec/specs/<slug>/hotfix.md`
- Implement the fix; confirm the observable proof from the `Verification` section is met
- End with a summary block: `Slug`, `Status`, `Fix`, `Verified`, `Next command`
- When done: `Next command: /draftspec.archive <slug>`

## Self-Check

- Did I write the hotfix spec before touching any code?
- Is the fix limited to files listed in `Touches`?
- Is the observable proof from `Verification` confirmed?
