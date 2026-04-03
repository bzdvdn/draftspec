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

- `.draftspec/constitution.summary.md` if present; otherwise `.draftspec/constitution.md`
- `.draftspec/plans/<slug>/plan.md`

## Load If Present

Read these only when the decomposition requires them:

- `.draftspec/specs/<slug>/summary.md` if present; otherwise `.draftspec/specs/<slug>/spec.md` — when task intent, scope, or acceptance boundaries are unclear
- `.draftspec/plans/<slug>/data-model.md` when task decomposition depends on entities, invariants, or lifecycle details
- `.draftspec/plans/<slug>/contracts/` when tasks involve APIs, events, or integration boundaries
- `.draftspec/plans/<slug>/research.md` only when it exists and affects implementation sequencing or risk
- the smallest set of implementation files needed to confirm task boundaries, sequencing, or concrete outcomes

Do not assume `research.md` should exist; use it only when the plan clearly depends on preserved uncertainty, an external dependency, or a documented trade-off.

## Do Not Read By Default

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
- Add a `Touches:` hint to each task line whenever the implementation surface is known with confidence. This is the primary mechanism for preventing re-reads during implement — the more concrete the surfaces, the fewer exploratory reads the implement agent needs to make. Keep it compact (`Touches: auth handler, session store`) and omit only when the surface genuinely cannot be named yet.
- After writing all task phases, add `## Surface Map` before the first phase: a two-column table mapping each unique implementation surface to the task IDs that touch it (`Surface | Tasks`). This lets `implement` batch-read all needed files in one pass without scanning `Touches:` across individual task lines.
- The tasks taken together MUST cover all acceptance criteria from the spec. Any uncovered criterion is a blocker.
- The `## Acceptance Coverage` section MUST include at least one explicit coverage line for each acceptance criterion.
- Coverage lines SHOULD reference stable acceptance IDs and task IDs such as `AC-001 -> T1.1, T2.1`.
- For newly created task lists, task IDs are required.
- When meaningfully updating an existing task list without task IDs, normalize it to the ID-based format.

## Content Quality Rules

- Each phase should have a short goal that explains why the phase exists.
- Prefer a few concrete tasks with measurable outcomes over many tiny bookkeeping items.
- Keep the outcome part of each task to ≤ 12 words. If more words are needed, the task is not concrete enough — split it or tighten the verb.
- When the acceptance proof is simple, embed it directly in the outcome instead of requiring a spec lookup: prefer `add POST /auth/login — returns 200 with JWT token field — AC-001` over `add login handler — endpoint works — AC-001`.
- Use action verbs tied to observable work: implement, add, migrate, validate, remove, backfill, document.
- Keep foundational setup separate from core behavior and separate proof/validation from broad implementation.
- If a task would benefit from a surface hint, prefer a compact `Touches: auth flow, session store` style note instead of speculative file-by-file path lists.
- When a task exists only to prove behavior, make that explicit instead of hiding it inside a larger implementation task.
- If a phase is unnecessary for this feature, omit it or state that it is intentionally not needed instead of filling it with generic tasks.
- Task text should make the intended outcome obvious to a reviewer without needing to reopen the plan for every line.
- Negative examples: avoid `misc`, `cleanup as needed`, `wire everything up`, `final polish`, or task text that hides the outcome behind a generic verb.
- Avoid `misc`, `cleanup as needed`, `wire everything up`, `final polish`, or other vague umbrella wording.

## Output expectations

- Write or patch `.draftspec/plans/<slug>/tasks.md`; call out blockers if decomposition is not yet possible
- End with a summary block: `Slug`, `Status`, `Artifacts`, `Blockers`, `Next command`
- When ready: `Next command: /draftspec.implement <slug>`

## Self-Check

- Does every task have a stable task ID and measurable outcome?
- Is every acceptance criterion covered explicitly?
- Could another developer execute these tasks in order without guessing what `done` means for each line?
