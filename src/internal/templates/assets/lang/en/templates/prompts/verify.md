# Draftspec Verify Prompt

You are verifying one implemented feature package after task execution.

## Goal

Confirm whether the implemented work is aligned enough with tasks and project rules to proceed safely.

## Phase Contract

Inputs: see Load First and Load Only If Needed.
Outputs: verification report in conversation or file when requested.
Stop if: see Stop Conditions.

## Load First

Always read these first:

- `.draftspec/constitution.md`
- `.draftspec/plans/<slug>/tasks.md`

## Load Only If Needed

Read these only when needed to confirm a concrete claim:

- `.draftspec/specs/<slug>.md`
- `.draftspec/plans/<slug>/plan.md`
- `.draftspec/plans/<slug>/data-model.md`
- `.draftspec/plans/<slug>/contracts/`
- only the code files needed to confirm whether a task or acceptance claim was actually implemented

## Do Not Read By Default

- unrelated specs
- unrelated plan packages
- unrelated code areas
- broad repository history
- archives unless the current verification explicitly depends on them

## Stop Conditions

Stop and ask for clarification only if:

- the slug is ambiguous
- the tasks file is missing
- the verification would otherwise invent implementation facts
- the requested conclusion would require a broad repository sweep instead of focused evidence for this feature package
- the implementation claim cannot be confirmed from the current tasks, plan artifacts, and targeted code inspection

## Rules

- Start from `tasks.md` as the verification entrypoint.
- If `/.draftspec/scripts/check-verify-ready.*` is available, prefer it as the cheap first pass before reading deeper artifacts.
- Use `/.draftspec/scripts/verify-task-state.*` only as a fallback when the phase-readiness wrapper is unavailable.
- Prefer helper script output over reading helper script source.
- Do not read `/.draftspec/scripts/*` by default unless you are debugging the script, working on Draftspec itself, or the user explicitly asks to inspect script logic.
- Prefer confirming concrete implementation claims over broad subjective review.
- Verify that completed tasks are consistent with the current state of the feature package.
- Verify that open tasks do not contradict any claim that the feature is fully complete.
- Verify acceptance-to-task coverage consistency when `tasks.md` includes an `Acceptance Coverage` section.
- When `tasks.md` uses task IDs such as `T1.1`, reference those IDs directly in checks, findings, and conclusions.
- Keep default verification structural and cheap by default.
- Only deepen into broader implementation validation when the user explicitly asks for it or when a concrete contradiction cannot be resolved from tasks, plan artifacts, and focused evidence.
- Use a simple verdict: `pass`, `concerns`, or `blocked`.
- Use `pass` when no blocking problems are present and only minor or no warnings remain.
- Use `concerns` when the feature can move forward, but warnings or open questions should be resolved soon.
- Use `blocked` when missing task completion or contradictory implementation state would make archive or completion claims unsafe.
- Keep the verification output in the project's configured documentation language when writing it to disk.
- Use `/.draftspec/scripts/verify-task-state.* <slug>` as the fallback first pass only when `check-verify-ready.*` is unavailable.
- Use `.draftspec/templates/verify-report.md` as the canonical template when writing the report to disk.
- Use this report structure:
  - `# Verify Report: <slug>`
  - `## Scope`
  - `## Verdict`
  - `## Checks`
  - `## Errors`
  - `## Warnings`
  - `## Questions`
  - `## Next Step`

## Output expectations

- Output the report to the conversation unless the user asks to persist it
- If persisted without an explicit path, use `.draftspec/plans/<slug>/verify.md`
- Summarize the verdict, completed checks, remaining concerns, and whether the feature is safe to archive
- In `## Checks`, explicitly cover task completion and implementation alignment where inspected
