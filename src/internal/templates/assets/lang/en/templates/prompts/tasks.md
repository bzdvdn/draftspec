# Draftspec Tasks Prompt

You are creating or updating `.draftspec/plans/<slug>/tasks.md`.

## Goal

Break an approved plan into executable implementation tasks.

## Operating Mode

- Decompose one approved plan package only.
- Use `plan.md` as the entrypoint and go deeper only when required.
- Produce the smallest task list that still covers the feature safely.
- Prefer explicit sequencing over umbrella tasks.

## Load First

Always read these before decomposing the work:

- `.draftspec/constitution.md`
- `.draftspec/plans/<slug>/plan.md`

## Load Only If Needed

Read these only when the decomposition requires them:

- `.draftspec/specs/<slug>.md` when task intent, scope, or acceptance boundaries are unclear
- `.draftspec/plans/<slug>/data-model.md` when task decomposition depends on entities, invariants, or lifecycle details
- `.draftspec/plans/<slug>/contracts/` when tasks involve APIs, events, or integration boundaries
- `.draftspec/plans/<slug>/research.md` only when it exists and affects implementation sequencing or risk

## Do Not Read By Default

- unrelated specs
- unrelated plan packages
- implementation files that are not needed to decompose the work
- broad repository history

## Stop Conditions

Stop and ask for refinement if:

- `.draftspec/plans/<slug>/plan.md` is missing
- tasks would be vague because the plan is underspecified
- the current decomposition requires spec, data model, contracts, or research that are missing
- the constitution blocks the proposed decomposition

Do not jump ahead into implementation.

## Invariants

- Tasks MUST align with the plan and constitution.
- Use `plan.md` as the decomposition entrypoint.
- Never read unrelated feature artifacts to compensate for underspecified planning.
- The task list must be executable in order.
- Every acceptance criterion must be covered by at least one task.
- Prefer concrete, testable, implementation-oriented tasks.
- Include validation and documentation alignment work only when needed.
- Do not generate vague umbrella tasks.

## Language Rules

- Use the project's configured documentation language for new or updated task content.
- Preserve an established local task-document convention only when needed to keep an existing file coherent.
- Do not mix task languages inside the same task list without a strong project reason.
- Load deeper artifacts only when the current decomposition needs them.

## Task Format Rules

- Follow the structure of `.draftspec/templates/tasks.md`: group tasks into ordered phases (`## Phase N: Name`).
- Each task MUST include a phase-scoped task ID in the form `T<phase>.<index>`.
- Each task MUST follow the format: `- [ ] T<phase>.<index> <action verb> — <concrete measurable outcome>`
- The tasks taken together MUST cover all acceptance criteria from the spec. Any uncovered criterion is a blocker.
- The `## Acceptance Coverage` section MUST include at least one explicit coverage line for each acceptance criterion.
- Coverage lines SHOULD reference stable acceptance IDs and task IDs such as `AC-001 -> T1.1, T2.1`.
- For newly created task lists, task IDs are required.
- When meaningfully updating an existing task list without task IDs, normalize it to the ID-based format.

## Output expectations

- Write or patch `.draftspec/plans/<slug>/tasks.md`
- Ensure tasks can be executed in order
- Call out blockers or missing inputs if decomposition is not yet possible

## Self-Check

- Did I decompose from `plan.md` first?
- Does every task have a stable task ID and measurable outcome?
- Is every acceptance criterion covered explicitly?
- Did I avoid reading implementation files unless decomposition required it?
