# Workflow Model

## Strict Phase Chain

```text
constitution -> spec -> inspect -> plan -> tasks -> implement -> verify -> archive
```

## Phase Responsibilities

### `constitution`

Defines the non-negotiable rules of the project.

Mandatory sections:

- `Purpose`
- `Core Principles`
- `Constraints`
- `Language Policy`
- `Development Workflow`
- `Governance`
- `Last Updated`

### `spec`

Captures one feature request as a concrete spec. Acceptance criteria should use canonical `Given / When / Then` markers even when the surrounding document language is Russian.

### `inspect`

Checks consistency and quality for a single feature. It can flag missing scenarios, weak acceptance criteria, constitutional conflicts, plan drift, or missing task coverage.

A full inspection report should use a stable structure:

- `# Inspect Report: <slug>`
- `## Scope`
- `## Verdict`
- `## Errors`
- `## Warnings`
- `## Questions`
- `## Suggestions`
- `## Traceability`
- `## Next Step`

`Verdict` should be one of:

- `pass`
- `concerns`
- `blocked`

Suggested semantics:

- `pass`: no blocking problems; only minor or no warnings remain
- `concerns`: the workflow can continue, but warnings or open questions should be resolved soon
- `blocked`: the next workflow step would otherwise proceed on missing or contradictory information

When an inspection report should be persisted to disk, Draftspec should prefer these default paths:

- `.draftspec/plans/<slug>/inspect.md` when the plan package already exists
- `.draftspec/specs/<slug>.inspect.md` when no plan package exists yet

Use `.draftspec/templates/inspect-report.md` as the canonical template when the report is written to disk.

Stable acceptance IDs such as `AC-001` make traceability lighter and easier to validate.

### `plan`

Produces technical design artifacts for one feature package:

- `plan.md`
- `data-model.md`
- `contracts/`
- optional `research.md`

### `tasks`

Turns the plan package into executable tasks. `tasks.md` lives next to other plan artifacts inside `.draftspec/plans/<slug>/`.

### `implement`

Executes unfinished tasks, updates `tasks.md`, and refreshes `memory.md`.

### `verify`

Runs a lightweight post-implementation check to confirm that completed work is aligned enough with tasks, memory, and project rules to move forward safely.

A full verification report should use a stable structure:

- `# Verify Report: <slug>`
- `## Scope`
- `## Verdict`
- `## Checks`
- `## Errors`
- `## Warnings`
- `## Questions`
- `## Next Step`

`Verdict` should be one of:

- `pass`
- `concerns`
- `blocked`

Suggested semantics:

- `pass`: no blocking problems are present; only minor or no warnings remain
- `concerns`: the feature can move forward, but warnings or open questions should be resolved soon
- `blocked`: archive or completion claims would otherwise proceed on contradictory implementation state, unfinished required work, or unsynchronized memory

When a verification report should be persisted to disk, Draftspec should prefer `.draftspec/plans/<slug>/verify.md`.

Use `.draftspec/templates/verify-report.md` as the canonical template when the report is written to disk.

Use `.draftspec/scripts/verify-task-state.sh <slug>` as the cheapest first-pass helper when you only need task-state confirmation.

### `archive`

Copies a completed, superseded, rejected, abandoned, or deferred feature package into `.draftspec/archive/<slug>/<YYYY-MM-DD>/`.

When archiving with status `completed`, Draftspec should first use `.draftspec/scripts/verify-task-state.sh <slug>` and treat remaining open tasks as a blocker.

## Why This Chain Exists

The chain keeps the product strict without becoming bureaucratic:

- architecture and workflow rules come first
- user intent becomes a spec before technical planning starts
- technical planning happens before task breakdown
- implementation follows tasks instead of improvisation
- lightweight verification closes the gap between implementation and archive
- completed feature packages can be archived without bloating active memory
