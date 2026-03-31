# Draftspec Tasks Prompt

You are creating or updating `.draftspec/plans/<slug>/tasks.md`.

## Goal

Break an approved plan into executable implementation tasks.

## Phase Contract

Inputs: `.draftspec/constitution.md`, `.draftspec/plans/<slug>/plan.md`; optionally spec, data-model, contracts when decomposition requires them.
Outputs: `.draftspec/plans/<slug>/tasks.md` with phased task list and Acceptance Coverage section.
Stop if: plan.md missing, plan underspecified, or any AC cannot be mapped to executable work.

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
- the smallest set of implementation files needed to confirm task boundaries, sequencing, or concrete outcomes

Do not assume `research.md` should exist; use it only when the plan clearly depends on preserved uncertainty, an external dependency, or a documented trade-off.

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
- the decomposition would span multiple feature slugs or unrelated change sets
- one or more acceptance criteria cannot be mapped to executable work without guessing

Do not jump ahead into implementation.

## Invariants

- Tasks MUST align with the plan and constitution.
- Use `plan.md` as the decomposition entrypoint.
- Never read unrelated feature artifacts to compensate for underspecified planning.
- Read implementation code only when the task list would otherwise stay vague; prefer a narrow file slice over broad exploration.
- When `/.draftspec/scripts/check-tasks-ready.*` is available, prefer running it as the phase readiness check instead of reading script source.
- Do not read `/.draftspec/scripts/*` by default unless you are debugging the scripts, working on Draftspec itself, or the user explicitly asks to inspect script logic.
- The task list must be executable in order.
- Every acceptance criterion must be covered by at least one task.
- Prefer concrete, testable, implementation-oriented tasks.
- Include validation and documentation alignment work only when needed.
- Do not generate vague umbrella tasks.
- The task list should be readable to both an implementation agent and a human reviewer without extra interpretation.
- Targeted code reading during task decomposition is useful when it reduces re-reading during implementation.
- Do not start implementation work, edit source code, or claim tasks are already done during the tasks phase.

## Language Rules

- Use the project's configured documentation language for new or updated task content.
- Preserve an established local task-document convention only when needed to keep an existing file coherent.
- Do not mix task languages inside the same task list without a strong project reason.
- Load deeper artifacts only when the current decomposition needs them.

## Task Format Rules

- Follow the structure of `.draftspec/templates/tasks.md`: group tasks into ordered phases (`## Phase N: Name`).
- Each task MUST include a phase-scoped task ID in the form `T<phase>.<index>`.
- Each task MUST follow the format: `- [ ] T<phase>.<index> <action verb> — <concrete measurable outcome>`
- Each task SHOULD reference 1-2 stable IDs when possible (`AC-*`, `RQ-*`, `DEC-*`).
- When it materially reduces downstream guesswork, add a short `Touches:` hint to a phase goal or task line naming the concrete surfaces likely to change. Keep it compact and use it only when confidence is high.
- The tasks taken together MUST cover all acceptance criteria from the spec. Any uncovered criterion is a blocker.
- The `## Acceptance Coverage` section MUST include at least one explicit coverage line for each acceptance criterion.
- Coverage lines SHOULD reference stable acceptance IDs and task IDs such as `AC-001 -> T1.1, T2.1`.
- For newly created task lists, task IDs are required.
- When meaningfully updating an existing task list without task IDs, normalize it to the ID-based format.

## Content Quality Rules

- Each phase should have a short goal that explains why the phase exists.
- Prefer a few concrete tasks with measurable outcomes over many tiny bookkeeping items.
- Use action verbs tied to observable work: implement, add, migrate, validate, remove, backfill, document.
- Keep foundational setup separate from core behavior and separate proof/validation from broad implementation.
- If a task would benefit from a surface hint, prefer a compact `Touches: auth flow, session store` style note instead of speculative file-by-file path lists.
- When a task exists only to prove behavior, make that explicit instead of hiding it inside a larger implementation task.
- If a phase is unnecessary for this feature, omit it or state that it is intentionally not needed instead of filling it with generic tasks.
- Task text should make the intended outcome obvious to a reviewer without needing to reopen the plan for every line.
- Negative examples: avoid `misc`, `cleanup as needed`, `wire everything up`, `final polish`, or task text that hides the outcome behind a generic verb.
- Avoid `misc`, `cleanup as needed`, `wire everything up`, `final polish`, or other vague umbrella wording.

## Output expectations

- Write or patch `.draftspec/plans/<slug>/tasks.md`
- Ensure tasks can be executed in order
- Call out blockers or missing inputs if decomposition is not yet possible
- When referring to created or updated files in the conversation, list their exact project-relative paths, not only bare filenames
- End the conversation with a short stable summary block that includes `Slug`, `Status`, `Artifacts`, `Blockers`, and `Next command` when that handoff is truly safe
- When task decomposition is complete and implementation can begin, end the conversation summary with `Next command: /draftspec.implement <slug>`
- If tasks are blocked or still need refinement, say that directly instead of suggesting `/draftspec.implement`

## Self-Check

- Did I decompose from `plan.md` first?
- Does every task have a stable task ID and measurable outcome?
- Is every acceptance criterion covered explicitly?
- Did I avoid reading implementation files unless decomposition required it?
- If I read code, did I read only the smallest slice needed to avoid vague tasks?
- Could another developer execute these tasks in order without guessing what `done` means for each line?
