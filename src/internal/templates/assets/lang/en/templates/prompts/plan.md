# Draftspec Plan Prompt

You are creating or updating the implementation plan package for one feature.

## Goal

Produce the technical planning artifacts for a spec under `.draftspec/plans/<slug>/`.

## Operating Mode

- Plan one feature only.
- Prefer patching existing artifacts over broad rewrites.
- Keep context narrow and repository-grounded.
- Produce only the artifacts justified by the feature.

## Load Only

- `.draftspec/constitution.md`
- `.draftspec/specs/<slug>.md`
- only the repository code and docs needed to plan this one feature

## Do Not Read By Default

- unrelated specs
- unrelated plan packages
- large repository areas with no impact on this feature
- optional `research.md` unless uncertainty already exists

## Stop Conditions

Stop and ask for clarification or refinement if:

- `.draftspec/specs/<slug>.md` does not exist
- the spec is too vague to produce architecture, contracts, or data model decisions
- constitutional constraints conflict with the intended plan

Do not compensate by reading broad unrelated repository context.

## Required outputs

Always create or update:

- `.draftspec/plans/<slug>/plan.md`
- `.draftspec/plans/<slug>/data-model.md`

Create only when the feature actually requires it:

- `.draftspec/plans/<slug>/contracts/api.md` — only if the feature touches API boundaries
- `.draftspec/plans/<slug>/contracts/events.md` — only if the feature produces or consumes events

Create `.draftspec/plans/<slug>/research.md` only when real uncertainty or external investigation is needed.

## Invariants

- The plan MUST comply with the constitution.
- Keep planning tied to the current spec, not idealized architecture.
- Never read unrelated feature artifacts to compensate for missing clarity.
- When `/.draftspec/scripts/check-plan-ready.*` is available, prefer running it as the phase readiness check instead of reading script source.
- Do not read `/.draftspec/scripts/*` by default unless you are debugging the scripts, working on Draftspec itself, or the user explicitly asks to inspect script logic.
- Prefer concrete implementation decisions over generic advice.
- Optional artifacts stay optional; do not create them by habit.
- The plan is only complete when downstream task decomposition can proceed without guessing.

## Language Rules

- Use the project's configured documentation language for all new or updated planning artifacts.
- Keep the language of `plan.md`, `data-model.md`, `contracts/`, and optional `research.md` internally consistent.
- Respect an established local document convention only when preserving an existing artifact would otherwise become inconsistent.

## Traceability Rules

- Follow the structure of `.draftspec/templates/plan.md` and `.draftspec/templates/data-model.md` when creating new files.
- Data model and contracts MUST be consistent with the spec and its acceptance criteria.
- Reference stable acceptance IDs from the spec when discussing acceptance-critical behavior.
- When the plan makes a significant implementation choice, record it as a stable decision ID such as `DEC-001`.
- If `data-model.md` is created, explicitly state which entities, invariants, or lifecycle concerns require it.
- If `contracts/` is created, explicitly state which API or event boundary requires it.
- If neither richer artifact is needed, prefer not creating it.
- Use repository reality, not idealized architecture.
- If critical information is missing, ask only the minimum necessary follow-up questions.

## Output expectations

- Write or patch the plan artifacts
- State which optional artifacts were created and why
- Summarize the key technical decisions that will affect task decomposition and implementation
- Explicitly call out risks and unresolved questions that block downstream phases

## Self-Check

- Did I plan only one feature?
- Did I keep optional artifacts justified?
- Did I reference acceptance-critical behavior with stable IDs where needed?
- Can `tasks` be written from this package without guesswork?
