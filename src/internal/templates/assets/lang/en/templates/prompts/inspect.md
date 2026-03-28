# Draftspec Inspect Prompt

You are inspecting one feature package for consistency and quality.

## Goal

Produce a focused inspection report for one feature without expanding scope.

## Load First

Always read these first:

- `.draftspec/constitution.md`
- `.draftspec/memory.md`
- `.draftspec/specs/<slug>.md`

## Load Only If Needed

Read these only when they exist and materially affect the inspection:

- `.draftspec/plans/<slug>/plan.md`
- `.draftspec/plans/<slug>/tasks.md`
- `.draftspec/plans/<slug>/data-model.md`
- `.draftspec/plans/<slug>/contracts/`
- `.draftspec/plans/<slug>/research.md`

## Do Not Read By Default

- unrelated specs
- unrelated plan packages
- broad repository history
- implementation files unless they are needed to verify a concrete consistency claim

## Stop Conditions

Stop and ask a minimal follow-up question only if:

- the target slug is ambiguous
- the spec is missing entirely
- the inspection would otherwise invent missing product intent

## Rules

- Check constitutional consistency first.
- Inspect spec completeness and clarity.
- Every acceptance criterion in the spec MUST have an explicit Given/When/Then format. The `Given`, `When`, and `Then` markers remain canonical regardless of the documentation language. Missing G/W/T is an `Error`, not a `Suggestion`.
- If `tasks.md` exists, verify that every acceptance criterion from the spec is covered by at least one task. An uncovered criterion is an `Error`.
- If plan artifacts exist, check alignment between spec, plan, data model, contracts, and tasks.
- Keep the inspection report in the project's configured documentation language when writing it to disk.
- Prefer concrete findings over generic advice.
- Use this report structure:
  - `# Inspect Report: <slug>`
  - `## Scope`
  - `## Verdict`
  - `## Errors`
  - `## Warnings`
  - `## Questions`
  - `## Suggestions`
  - `## Traceability`
  - `## Next Step`
- The `## Verdict` section MUST use one of: `pass`, `concerns`, `blocked`.
- Use `pass` when no errors are present and only minor or no warnings remain.
- Use `concerns` when the feature can still move forward, but warnings or open questions should be resolved soon.
- Use `blocked` when constitutional conflicts, missing spec intent, missing Given/When/Then acceptance criteria, or uncovered acceptance criteria prevent the next workflow step from proceeding safely.
- `## Traceability` should summarize how acceptance criteria map to tasks when `tasks.md` exists.
- Prefer traceability statements that reference stable acceptance IDs such as `AC-001`.
- `## Next Step` should say whether it is safe to continue to `plan`, `tasks`, or whether refinement is required first.

## Output expectations

- Output the report to the conversation unless the user asks to persist it.
- If the user asks to persist the report without specifying a file path, use `.draftspec/plans/<slug>/inspect.md` when the plan package exists. Otherwise use `.draftspec/specs/<slug>.inspect.md`.
- If the user provides an explicit file path, use that path.
- Summarize errors, warnings, open questions, suggestions, and the final verdict.
