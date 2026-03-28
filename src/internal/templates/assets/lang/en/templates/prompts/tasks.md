# Draftspec Tasks Prompt

You are creating or updating `.draftspec/plans/<slug>/tasks.md`.

## Goal

Break an approved plan into executable implementation tasks.

## Load First

Always read these before decomposing the work:

- `.draftspec/constitution.md`
- `.draftspec/memory.md`
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

## Rules

- Tasks MUST align with the plan and constitution.
- Use `plan.md` as the decomposition entrypoint.
- Use the project's configured documentation language for new or updated task content.
- Preserve an established local task-document convention only when needed to keep an existing file coherent.
- Do not mix task languages inside the same task list without a strong project reason.
- Load deeper artifacts only when the current decomposition needs them.
- Follow the structure of `.draftspec/templates/tasks.md`: group tasks into ordered phases (`## Phase N: Name`).
- Each task MUST follow the format: `- [ ] <action verb> — <concrete measurable outcome>`
- The tasks taken together MUST cover all acceptance criteria from the spec. Any uncovered criterion is a blocker.
- The `## Acceptance Coverage` section MUST include at least one explicit coverage line for each acceptance criterion.
- Coverage lines SHOULD reference stable acceptance IDs such as `AC-001 -> Task 1, Task 2`.
- Prefer tasks that are concrete, testable, and implementation-oriented.
- Include validation and documentation alignment work where needed.
- Do not generate vague umbrella tasks.

## Output expectations

- Write or patch `.draftspec/plans/<slug>/tasks.md`
- Ensure tasks can be executed in order
- Call out blockers or missing inputs if decomposition is not yet possible
