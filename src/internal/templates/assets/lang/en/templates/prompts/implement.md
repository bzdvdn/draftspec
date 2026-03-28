# Draftspec Implement Prompt

You are executing a planned feature implementation.

## Goal

Implement the feature by following the existing task list without expanding scope.

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

If all tasks in `tasks.md` are already marked complete, say so and do not continue.

Do not broaden scope to solve these problems.

## Rules

- Implement only unfinished tasks from `tasks.md`.
- Respect the order and phase structure in `tasks.md`.
- Use `tasks.md` as the execution entrypoint.
- Load deeper artifacts only when the current task requires them.
- Do not violate the constitution.
- Follow the project's preferred code comment language as recorded in `.draftspec/draftspec.yaml` and `.draftspec/constitution.md`.
- When adding or editing code comments, keep them in the configured comment language unless the surrounding file already uses a different established convention that should be preserved.
- Do not introduce mixed-language comments in the same local code area without a strong reason.
- If the plan or tasks are insufficient, stop and request refinement instead of inventing broad new scope.
- Mark completed tasks in `tasks.md`.

## Output expectations

- Implement the work
- Update `tasks.md` checkboxes for completed items
- Summarize completed tasks, remaining tasks, and any blockers
- Explicitly state which acceptance criteria from the spec are now covered by the implementation
