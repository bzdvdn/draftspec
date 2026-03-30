# Draftspec Inspect Prompt

You are inspecting one feature package for consistency and quality.

## Goal

Produce a focused inspection report for one feature without expanding scope.

## Phase Contract

Inputs: see Load First and Load Only If Needed.
Outputs: inspection report in conversation or file when requested.
Stop if: see Stop Conditions.

## Load First

Always read these first:

- `.draftspec/constitution.md`
- `.draftspec/specs/<slug>.md`

## Load Only If Needed

Read these only when they exist and materially affect the inspection:

- `.draftspec/plans/<slug>/plan.md`
- `.draftspec/plans/<slug>/tasks.md`

## Do Not Read By Default

- `.draftspec/plans/<slug>/data-model.md`
- `.draftspec/plans/<slug>/contracts/`
- `.draftspec/plans/<slug>/research.md`
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
- If `/.draftspec/scripts/check-inspect-ready.*` is available, prefer it as the cheap first pass before deepening into artifacts.
- Use `/.draftspec/scripts/inspect-spec.*` only as a fallback when the phase-readiness wrapper is unavailable.
- Prefer helper script output over reading helper script source.
- Do not read `/.draftspec/scripts/*` by default unless you are debugging the script, working on Draftspec itself, or the user explicitly asks to inspect script logic.
- Inspect spec completeness and clarity.
- Verify `constitution <-> spec`: the spec must not conflict with explicit constitutional constraints, workflow rules, or language policy.
- Every acceptance criterion in the spec MUST have an explicit Given/When/Then format. The `Given`, `When`, and `Then` markers remain canonical regardless of the documentation language. Missing G/W/T is an `Error`, not a `Suggestion`.
- If `tasks.md` exists, verify that every acceptance criterion from the spec is covered by at least one task. An uncovered criterion is an `Error`.
- If `tasks.md` uses task IDs such as `T1.1`, prefer traceability statements that reference those task IDs directly.
- Prefer the cheapest inspection scope first: `constitution.md` and `spec.md`, then `plan.md`, then `tasks.md`, and only then deeper plan artifacts when a concrete claim requires them.
- If no `plan.md` exists, do not widen the inspection into optional plan artifacts or implementation code.
- If plan artifacts exist, check alignment between spec, plan, data model, contracts, and tasks.
- When `plan.md` exists, check `spec <-> plan` consistency before reading deeper plan artifacts.
- Verify `spec <-> plan`: the plan should preserve the feature goal, reflect major acceptance-critical behavior, and avoid unjustified new workstreams.
- If `tasks.md` exists, verify `plan <-> tasks`: task phases and task IDs should reflect the plan intent without obvious missing work for acceptance-critical behavior.
- Treat `spec.md` and `plan.md` as the required inputs for cheap plan consistency checks.
- Only read `data-model.md` or `contracts/` when `plan.md` explicitly depends on them or when they are required to confirm a concrete consistency claim.
- Check `Goal Alignment`: the plan must not change the core feature goal expressed in the spec.
- Check `Scope Expansion`: the plan must not introduce major new workstreams, components, or integration surfaces that are outside the spec.
- Check `Acceptance Coverage at Plan Level`: major acceptance-critical behavior from the spec should be reflected in the plan intent, even before tasks exist.
- Check `Constitution Consistency`: the plan must not violate constitutional rules or architectural constraints.
- Check `Artifact Justification`: if the plan introduces `data-model.md` or `contracts/`, the need for those artifacts should be justified by the spec.
- Use `blocked` when constitutional conflicts, missing product intent, missing Given/When/Then markers, uncovered acceptance criteria, or major `spec <-> plan` contradictions would make the next phase unsafe.
- Use `concerns` when the feature is still broadly aligned but has weak scope boundaries, under-justified artifacts, incomplete traceability, or open questions that should be resolved soon.
- Use `pass` when no blocking contradictions are present and only minor or no warnings remain.
- Do not turn this into a broad design review. Prefer catching obvious drift over scoring architecture quality.
- Keep the inspection report in the project's configured documentation language when writing it to disk.
- Prefer concrete findings over generic advice.
- Default to a compact report in conversation output: always include `Verdict`, include `Errors`, `Warnings`, and `Next Step` when non-empty, and include `Questions`, `Suggestions`, or `Traceability` only when they add real signal.
- Produce the full sectioned report only when the user explicitly asks for a full report or when the report is being persisted to a file.
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
- Use `concerns` when the feature can still move forward, but warnings, traceability gaps, or open questions should be resolved soon.
- Use `blocked` when constitutional conflicts, missing spec intent, missing Given/When/Then acceptance criteria, uncovered acceptance criteria, or major `spec <-> plan` contradictions prevent the next workflow step from proceeding safely.
- `## Traceability` should summarize how acceptance criteria map to tasks when `tasks.md` exists.
- Prefer traceability statements that reference stable acceptance IDs and task IDs such as `AC-001 -> T1.1, T2.1`.
- `## Next Step` should say whether it is safe to continue to `plan`, `tasks`, or whether refinement is required first.

## Output expectations

- Output the report to the conversation unless the user asks to persist it.
- If the user asks to persist the report without specifying a file path, use `.draftspec/plans/<slug>/inspect.md` when the plan package exists. Otherwise use `.draftspec/specs/<slug>.inspect.md`.
- If the user provides an explicit file path, use that path.
- In default conversation mode, prefer a compact report with only non-empty sections.
- Summarize errors, warnings, open questions, suggestions, and the final verdict.
