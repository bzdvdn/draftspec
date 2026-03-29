# Examples

This page shows realistic end-to-end Draftspec scenarios for one feature package.

## Quick Usage Patterns

### New Project

When starting a greenfield project, Draftspec works best as a minimal project-context layer from day one.

Example:

```bash
draftspec init my-project --lang en --shell sh --agents codex
cd my-project
draftspec doctor .
```

What to do next:

- establish the `constitution` for project rules
- describe the first feature through `spec`
- prepare `plan` and `tasks`
- use `implement` only from the current task list

Why this helps:

- humans and agents start from the same rules
- project context stays explicit and editable from the beginning
- the workflow stays lightweight because Draftspec does not require a heavy process engine

### Existing Project

For a brownfield codebase, Draftspec should be adopted incrementally instead of trying to document the whole repository at once.

Example:

```bash
cd existing-project
draftspec init . --lang en --shell sh --agents codex
draftspec doctor .
```

Recommended starting point:

- establish the `constitution` around the project's current reality
- pick one active feature or change request
- create a spec only for that scope
- move to plan, tasks, and implement only within that feature package

What not to do:

- do not try to spec the whole project at once
- do not pull broad repository context unless the active feature really needs it

Why this helps:

- Draftspec adds a lightweight layer of discipline on top of an existing codebase
- adoption happens one feature at a time
- this keeps token usage down and avoids process bloat

### Prompt File Input

When `/draftspec.spec` starts from a local prompt file, prefer explicit metadata instead of relying on a generic filename such as `spec_prompt.md`.

Example prompt file:

```text
name: Add dark mode
slug: add-dark-mode

Add a user-selectable dark theme for the dashboard and settings pages.
```

This lets Draftspec:

- derive a safe spec path such as `.draftspec/specs/add-dark-mode.md`
- create or switch to `feature/add-dark-mode`
- avoid ambiguous slugs from generic filenames

## 1. Create a Constitution for a Brownfield Project

User request:

```text
/draftspec.constitution Python project, DDD style, split into API and workers, Kafka for asynchronous integration, ClickHouse as the analytical sink.
```

Expected agent behavior:

- read the constitution prompt in `.draftspec/templates/prompts/constitution.md`
- inspect only the minimum repository evidence needed
- create or patch `.draftspec/constitution.md`
- run `check-constitution.sh` when appropriate

Expected outcome:

- architecture rules are formalized
- development workflow rules become explicit
- the constitution becomes authoritative for later phases

## 2. Create a Spec

User request:

```text
/draftspec.spec Add partner-specific ingestion scheduling with retry policy overrides.
```

Expected agent behavior:

- read constitution first
- create `.draftspec/specs/partner-scheduling.md`
- write acceptance criteria using canonical `Given / When / Then`
- keep surrounding text in the configured documentation language

Example acceptance criterion:

```md
### Acceptance Criterion 1

- ID: AC-001
- **Given** a partner with a custom retry policy
- **When** the ingestion schedule is evaluated
- **Then** the worker uses the partner-specific retry window instead of the default policy
```

## 3. Inspect the Spec

User request:

```text
/draftspec.inspect partner-scheduling
```

Expected agent behavior:

- read constitution and `.draftspec/specs/partner-scheduling.md`
- check completeness, constitutional consistency, and scenario quality
- create a focused inspection report
- if the report should be persisted, prefer `.draftspec/specs/partner-scheduling.inspect.md` before planning and `.draftspec/plans/partner-scheduling/inspect.md` after the plan package exists
- use `.draftspec/templates/inspect-report.md` as the canonical report template

Typical findings:

- missing failure-path scenario
- unclear acceptance coverage for manual retry overrides
- open question about scheduler ownership

## 4. Create a Plan Package

User request:

```text
/draftspec.plan partner-scheduling
```

Expected agent behavior:

- read constitution and the spec
- create `.draftspec/plans/partner-scheduling/plan.md`
- create `.draftspec/plans/partner-scheduling/data-model.md`
- create `.draftspec/plans/partner-scheduling/contracts/`
- create `research.md` only if uncertainty is real

Typical outputs:

- plan for scheduler integration points
- data model for partner overrides and retry windows
- event or API contracts for configuration updates

## 5. Create Tasks

User request:

```text
/draftspec.tasks partner-scheduling
```

Expected agent behavior:

- use `plan.md` as the decomposition entrypoint
- pull in spec, contracts, or data model only when needed
- produce `.draftspec/plans/partner-scheduling/tasks.md`
- include acceptance-to-task coverage

Example task structure:

```md
## Phase 1: Data Model

- [ ] T1.1 Add partner scheduling override model — override fields are persisted
- [ ] T1.2 Persist retry window fields — retry windows are available to scheduling logic

## Acceptance Coverage

- AC-001 -> T1.1, T1.2
```

## 6. Implement the Feature

User request:

```text
/draftspec.implement partner-scheduling
```

Expected agent behavior:

- start from `tasks.md`
- load spec, plan, data model, or contracts only for the active task
- implement unfinished tasks in order
- update `tasks.md`

This phase should avoid broad repository reads unless the active task actually requires them.

## 7. Verify the Feature

User request:

```text
/draftspec.verify partner-scheduling
```

Expected agent behavior:

- read constitution and tasks first
- confirm that completed tasks match the current implementation state closely enough
- produce a lightweight verification report
- start with `.draftspec/scripts/verify-task-state.sh partner-scheduling` when task-state confirmation is enough
- use `.draftspec/templates/verify-report.md` when the report should be persisted
- default to `.draftspec/plans/partner-scheduling/verify.md` when no explicit path is provided

## 8. Archive the Feature

User request:

```text
/draftspec.archive partner-scheduling --status completed --reason "implemented and merged"
```

Expected agent behavior:

- for `completed` status, start with `.draftspec/scripts/verify-task-state.sh partner-scheduling` and stop if open tasks remain
- copy the feature package into `.draftspec/archive/partner-scheduling/<YYYY-MM-DD>/`
- write `summary.md`

Expected archive result:

```text
.draftspec/archive/
  partner-scheduling/
    2026-03-28/
      summary.md
      spec.md
      plan.md
      tasks.md
      data-model.md
      contracts/
```

## 9. Agent Maintenance Scenario

A practical maintenance flow for agent targets:

```bash
draftspec add-agent my-project --agents claude --agents cursor
draftspec list-agents my-project
draftspec remove-agent my-project --agents cursor
draftspec cleanup-agents my-project
draftspec doctor my-project
```

Use this when a project changes its preferred agent mix over time.
