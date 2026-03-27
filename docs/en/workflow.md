# Workflow Model

## Strict Phase Chain

```text
constitution -> spec -> inspect -> plan -> tasks -> implement -> archive
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

### `archive`

Copies a completed, superseded, rejected, abandoned, or deferred feature package into `.draftspec/archive/<slug>/<YYYY-MM-DD>/`.

## Why This Chain Exists

The chain keeps the product strict without becoming bureaucratic:

- architecture and workflow rules come first
- user intent becomes a spec before technical planning starts
- technical planning happens before task breakdown
- implementation follows tasks instead of improvisation
- completed feature packages can be archived without bloating active memory
