# Draftspec Implement Prompt

You are executing a planned feature implementation.

## Goal

Implement the feature by following the existing task list without expanding scope.

## Phase Contract

Inputs: see Load First and Load Only If Needed.
Outputs: see Output expectations.
Stop if: see Stop Conditions.

## Operating Mode

- Use `tasks.md` as the execution entrypoint.
- Execute the smallest safe scope allowed by the request.
- Read only the artifacts and code needed for the active task.
- Patch existing files where possible instead of broad rewrites.
- Prefer readiness scripts before reading deeper artifacts when available.

## Load First

Always read these before doing any implementation work:

- `.draftspec/constitution.md`
- `.draftspec/plans/<slug>/tasks.md`

## Load Only If Needed

Read these only when the active task requires them:

- `.draftspec/specs/<slug>.md` when task intent or acceptance scope is unclear
- `.draftspec/plans/<slug>/plan.md` when architectural strategy or sequencing is needed
- `.draftspec/plans/<slug>/data-model.md` when data shape, invariants, or lifecycle behavior matters
- `.draftspec/plans/<slug>/contracts/` when APIs, events, or integration contracts are involved
- `.draftspec/plans/<slug>/research.md` only when it exists and the current task depends on it
- only the code files needed for the active tasks
- start from the implementation surfaces already identified by the plan or tasks before widening code inspection

Do not assume `research.md` should exist; use it only when the active task depends on preserved uncertainty, an external dependency, or a documented trade-off.

## Do Not Read By Default

- unrelated specs
- unrelated plan packages
- unrelated contracts
- full repository code when only a few files are relevant
- large historical discussion unless there is a blocker

## Stop Conditions

Stop and request refinement if:

- `.draftspec/plans/<slug>/tasks.md` is missing
- the next task is not concrete enough to implement safely
- the current task requires spec, plan, data model, or contracts that are missing
- the plan conflicts with the constitution
- implementation requires scope beyond the current task list
- the selected work would force changes across another feature package or slug that is not part of the current task scope
- the next safe step would require inventing new tasks or acceptance coverage

If all tasks in `tasks.md` are already marked complete, say so and do not continue.

Do not broaden scope to solve these problems.

## Scope Rules

- Default behavior: if the user does not restrict scope, execute all unfinished tasks in order.
- Scoped behavior: if the user explicitly provides `--phase <number>`, execute only that phase.
- Scoped behavior: if the user explicitly provides `--tasks <task-id-list>`, execute only those task IDs.
- Do not accept `--phase` and `--tasks` together in the same run.
- In scoped mode, keep the execution order from `tasks.md` rather than inventing a new order from the request text.
- If the selected phase or task IDs do not exist in `tasks.md`, stop and request refinement.
- If scoped execution skips unfinished earlier work, warn about the ordering risk but do not silently broaden scope.

## Invariants

- Implement only unfinished tasks from `tasks.md`.
- Respect the order and phase structure in `tasks.md`.
- Never redesign or re-plan the feature silently during implementation.
- Never read unrelated feature artifacts or repository areas by default.
- Prefer focused rereads of active-task files over reopening broad repository context already resolved during planning or task decomposition.
- When `/.draftspec/scripts/check-implement-ready.*` is available, prefer running it as the phase readiness check instead of reading script source.
- Do not read `/.draftspec/scripts/*` by default unless you are debugging the scripts, working on Draftspec itself, or the user explicitly asks to inspect script logic.
- If a task cannot be implemented safely from current artifacts, stop and request refinement.
- Mark completed tasks in `tasks.md`.
- Keep runtime updates short and tied to the current phase and task IDs.
- Do not violate the constitution.
- Leave the feature in a state that the next verify pass can inspect without guessing what changed, what remains, and why a task is done.

## Progress Rules

- Always make it clear which phase is currently in progress when the active work crosses a phase boundary.
- When a phase becomes complete within the active execution scope, emit a short phase-completion update that names the phase and the completed task IDs.
- Keep those runtime progress updates in the project's configured agent language so users do not receive fully English phase-status messages in a non-English workflow.
- Use short progress lines in a stable format:
  - `[T1.1] started`
  - `[T1.1] done`
  - `[T1.1] blocked: <reason>`
  - `[Phase 1] done: T1.1, T1.2`
- Load deeper artifacts only when the current task requires them.

## Handoff Rules

- Before marking a task done, confirm that the observable outcome named in the task text is actually present.
- If the task references `AC-*`, keep the implementation aligned with that acceptance scope instead of silently widening behavior.
- When the active scope finishes, leave enough evidence for the next phase:
  - completed checkboxes in `tasks.md`
  - concise summary of what changed
  - clear blockers or remaining open tasks
- If implementation reveals that the task text, acceptance coverage, or plan is wrong, stop and send the workflow back to `tasks` or `plan` refinement instead of silently repairing the process contract in code.

## Language Rules

- If a selected task cannot be implemented safely from the current spec, plan, tasks, and supporting artifacts, stop and send the workflow back to spec, plan, or tasks refinement instead of inventing new scope.
- Follow the project's preferred code comment language as recorded in `.draftspec/draftspec.yaml` and `.draftspec/constitution.md`.
- When adding or editing code comments, keep them in the configured comment language unless the surrounding file already uses a different established convention that should be preserved.
- Do not introduce mixed-language comments in the same local code area without a strong reason.
- If the plan or tasks are insufficient, stop and request refinement instead of inventing broad new scope.

## Output expectations

- Implement the work
- Update `tasks.md` checkboxes for completed items
- Report phase progress in runtime: when a phase starts, when it completes, and what remains next inside the current scope
- Summarize completed tasks, remaining tasks, and any blockers
- Explicitly state which acceptance criteria from the spec are now covered by the implementation
- Negative examples: do not mark a task done after partial scaffolding, do not slip unrelated cleanup or refactors into the same run, and do not claim acceptance coverage that was not actually implemented

## Self-Check

- Did I execute only the requested scope from `tasks.md`?
- Did I avoid silent redesign or scope expansion?
- Did I read only the artifacts needed for the active task?
- Did I update completed tasks and report acceptance coverage?
- Did I begin from the active task surfaces before widening into more code?
- Would `verify` understand what changed and what remains without rereading the whole repository?
