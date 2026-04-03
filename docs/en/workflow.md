# Workflow Model

## Strict Phase Chain

```text
constitution -> spec -> inspect -> plan -> tasks -> implement -> verify -> archive
```

Draftspec assumes branch-based delivery: each active feature should be developed in its own git branch, with the feature spec and plan package acting as the shared source of truth instead of a mutable global memory file. The default branch naming convention is `feature/<slug>`.

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

For agent-facing `/draftspec.spec`, Draftspec should support optional arguments:

- `--name <feature name>`
- `--slug <feature-slug>`
- `--branch <branch-name>`

Argument semantics:

- `--name` sets the canonical feature name for the current spec request
- `--slug` overrides the spec slug
- `--branch` overrides only the working branch and does not change the spec slug

`/draftspec.spec` should support two input modes:

- inline mode: the feature name and description are provided in the same message
- staged mode: the user first sends `/draftspec.spec --name ...` and then sends the feature description in the next message

When `/draftspec.spec` starts from a prompt file, Draftspec should prefer top-of-file metadata such as:

```text
name: Add dark mode
slug: add-dark-mode
```

Priority rules for the slug:

1. `--slug`
2. `slug:`
3. a slug derived from `--name`
4. a slug derived from `name:`
5. a safe fallback from the filename or short request text only when sufficiently specific

Priority rules for the feature name:

1. `--name`
2. `name:`
3. a concise feature name safely derived from the user request

If `/draftspec.spec` is invoked with `--name` but the feature description is still not detailed enough for a valid spec, Draftspec should not lose the request context: it should ask for the missing description or treat the next user message as the continuation of the same spec request.

By default, the feature branch should be `feature/<slug>`. If the user explicitly provides `--branch <name>`, Draftspec should use that branch name instead without changing the spec slug.

The spec itself should remain branch-agnostic: the working branch belongs to execution context, not to the spec document.

If the request is ambiguous, combines multiple features, or asks to derive one spec from multiple constitutional changes, Draftspec should stop and ask for one concrete feature before creating the branch or spec.

### `inspect`

Checks consistency and quality for a single feature. It can flag missing scenarios, weak acceptance criteria, constitutional conflicts, plan drift, or missing task coverage.

`inspect` is mandatory before `plan`. Planning should not proceed until the feature has a persisted inspect report at the canonical path.

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

When an inspection report is persisted to disk, Draftspec should use this canonical path:

- `.draftspec/specs/<slug>/inspect.md`

Use `.draftspec/templates/inspect-report.md` as the canonical template when the report is written to disk.
Persisted inspect and verify reports should start with a machine-readable metadata block containing `report_type`, `slug`, `status`, `docs_language`, and `generated_at`.

Stable acceptance IDs such as `AC-001` make traceability lighter and easier to validate.

For cheap `spec <-> plan` consistency checks, Draftspec should prefer this scope:

- always load: `constitution.md`, `spec.md`
- load if needed: `plan.md`, `tasks.md`
- conditional deeper reads only when a concrete claim requires them: `data-model.md`, `contracts/`, `research.md`
- do not read implementation code by default

The goal is to catch obvious drift, not to run a full architectural review. Useful checks include:

- constitution-to-spec alignment
- goal alignment
- unjustified scope expansion
- acceptance-critical behavior reflected at the plan level
- plan-to-task alignment when `tasks.md` exists
- constitutional consistency
- justification for richer plan artifacts such as `data-model.md` and `contracts/`

### `plan`

Produces technical design artifacts for one feature package:

- `plan.md`
- `data-model.md`
- `contracts/`
- optional `research.md`

### `tasks`

Turns the plan package into executable tasks. `tasks.md` lives next to other plan artifacts inside `.draftspec/plans/<slug>/`.

Tasks should be grouped by phase and use phase-scoped task IDs such as `T1.1`, `T1.2`, and `T2.1`.

Acceptance coverage should reference those task IDs directly:

```text
AC-001 -> T1.1, T2.1
```

### `implement`

Executes unfinished tasks and updates `tasks.md`.

Default behavior should remain full-run: without explicit scope flags, Draftspec continues through all unfinished tasks in task-list order.

Selective execution is allowed when the user explicitly narrows scope:

- `--phase <number>` for one implementation phase
- `--tasks <task-id-list>` for one or more specific task IDs such as `T1.1,T2.1`

`--phase` and `--tasks` should not be combined in the same run.

When selective execution skips unfinished earlier work, Draftspec should warn about the sequencing risk without silently broadening scope.

During implementation, Draftspec should emit short runtime progress updates whenever it starts or completes a phase in the active execution scope.

Those phase-status updates should follow the project's configured agent language rather than defaulting to English.

### `verify`

Runs a lightweight post-implementation check to confirm that completed work is aligned enough with tasks and project rules to move forward safely.

A full verification report should use a stable structure:

- `# Verify Report: <slug>`
- `## Scope`
- `## Verdict`
- `## Checks`
- `## Errors`
- `## Warnings`
- `## Questions`
- `## Not Verified`
- `## Next Step`

Recommended report details:

- `## Scope` should record the actual verification mode such as `default` or `deep`
- `## Scope` should list the concrete surfaces that were inspected
- `## Verdict` should include `archive_readiness`
- `## Verdict` should include a one-line summary of why the verdict is justified
- `## Checks` should include `task_state`
- `## Checks` should include `acceptance_evidence` for the `AC-*` items actually confirmed
- `## Checks` should include `implementation_alignment` tied to the concrete surface inspected
- `## Not Verified` should list material claims or surfaces that were intentionally not checked

`Verdict` should be one of:

- `pass`
- `concerns`
- `blocked`

Suggested semantics:

- `pass`: no blocking problems are present; only minor or no warnings remain
- `concerns`: the feature can move forward, but warnings or open questions should be resolved soon
- `blocked`: archive or completion claims would otherwise proceed on contradictory implementation state or unfinished required work

Use `concerns` rather than `pass` when the evidence is partial but no concrete contradiction has been found.

When a verification report should be persisted to disk, Draftspec should prefer `.draftspec/plans/<slug>/verify.md`.

Use `.draftspec/templates/verify-report.md` as the canonical template when the report is written to disk.

Persisted verify reports should start with the same machine-readable metadata block used by inspect reports: `report_type`, `slug`, `status`, `docs_language`, and `generated_at`.

When available, Draftspec should prefer `.draftspec/scripts/check-verify-ready.sh <slug>` as the cheap readiness pass before deeper verification.

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
- completed feature packages can be archived without bloating the active workspace
