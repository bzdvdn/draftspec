# Draftspec Verify Prompt

You are verifying one implemented feature package after task execution.

## Goal

Confirm whether the implemented work is aligned enough with tasks, memory, and project rules to proceed safely.

## Load First

Always read these first:

- `.draftspec/constitution.md`
- `.draftspec/memory.md`
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

## Rules

- Start from `tasks.md` as the verification entrypoint.
- Prefer confirming concrete implementation claims over broad subjective review.
- Verify that completed tasks are consistent with the current state of the feature package.
- Verify that open tasks do not contradict any claim that the feature is fully complete.
- Verify that `.draftspec/memory.md` reflects important implementation changes when relevant.
- Verify acceptance-to-task coverage consistency when `tasks.md` includes an `Acceptance Coverage` section.
- Use a simple verdict: `pass`, `concerns`, or `blocked`.
- Use `pass` when no blocking problems are present and only minor or no warnings remain.
- Use `concerns` when the feature can move forward, but warnings or open questions should be resolved soon.
- Use `blocked` when missing task completion, contradictory implementation state, or unsynchronized memory would make archive or completion claims unsafe.
- Keep the verification output in the project's configured documentation language when writing it to disk.
- Use `.draftspec/scripts/verify-task-state.sh <slug>` as a cheap first pass before reading deeper artifacts.
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
- In `## Checks`, explicitly cover task completion, memory alignment, and implementation alignment where inspected
