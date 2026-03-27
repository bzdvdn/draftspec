# FAQ

## What is the difference between `spec` and `inspect`?

`spec` creates or updates the feature specification.

`inspect` reviews that specification and related artifacts for completeness, consistency, and constitutional compliance. It is a quality gate, not an authoring phase.

## When should I create `research.md`?

Create `research.md` only when there is real uncertainty that should be preserved:

- external protocol or integration details need investigation
- multiple architecture options are being evaluated
- a design decision needs supporting evidence

Do not create it by default for every feature.

## When should I archive a feature?

Archive a feature when it is no longer active in the main workflow, for example when it is:

- completed
- superseded
- rejected
- abandoned
- deferred

The archive keeps historical context without bloating active project memory.

## What is the difference between `remove-agent` and `cleanup-agents`?

`remove-agent` updates `.draftspec/draftspec.yaml` and removes generated files for the selected enabled targets.

`cleanup-agents` removes leftover orphaned files that no longer match the enabled target set in config.

Use `doctor` after either one if you want to verify workspace health.

## Why does Draftspec keep `Given / When / Then` in English even in Russian docs?

Those markers are intentionally canonical. They are easier for agents to recognize consistently and easier for validation rules to enforce.

The surrounding document text can still be written in Russian.

## Does `implement` always need to read the whole feature package?

No. `implement` should start from `tasks.md` and load deeper artifacts only when the active task requires them.

Typical minimal read order:

- `constitution.md`
- `memory.md`
- `tasks.md`
- then `spec.md`, `plan.md`, `data-model.md`, `contracts/`, or `research.md` only if needed

## Can `plan` run without `spec`?

No. The intended chain is strict:

```text
constitution -> spec -> inspect -> plan -> tasks -> implement -> archive
```

`plan` depends on an existing spec.

## Why is `memory.md` still needed if there is an archive?

They solve different problems.

- `memory.md` stores the current working state of the project
- `archive/` stores historical snapshots of completed or inactive feature packages

`memory.md` may keep a short archived index, but the full history belongs in `archive/`.
