# draftspec MVP

## Product statement

`draftspec` provides a minimal, file-based context system for software projects.

It helps development agents and humans work from the same project context:

- what the project is trying to do
- what a feature should do
- how the feature will be implemented
- what contracts and data models are required
- what rules are non-negotiable
- how development must be conducted
- what decisions are currently in effect
- which language generated docs and agent prompts should use

The product should feel lightweight, editable, and resilient.

Draftspec should preserve traceability through compact, stable IDs instead of broader default context or extra summary layers.

## Non-goals

- no mandatory checkpoint flow
- no quickstart wizard
- no built-in approval process
- no public CLI command for creating specs in the first version
- no public CLI command for planning in the first version
- no public CLI command for tasks in the first version
- no public CLI command for implementation in the first version
- no AI orchestration embedded into the CLI in the first version

Managed generated artifacts should remain refreshable without touching authored feature state. `refresh` should update templates, scripts, config-derived values, project-local agent files, and the managed Draftspec block in `AGENTS.md` while leaving `constitution`, `specs`, `plans`, and `archive` untouched.

## Workspace layout

```text
.draftspec/
  draftspec.yaml
  constitution.md
  specs/
    <slug>.md
  plans/
    <slug>/
      plan.md
      tasks.md
      data-model.md
      research.md
      contracts/
        api.md
        events.md
  archive/
    <slug>/
      <YYYY-MM-DD>/
      summary.md
      spec.md
      plan.md
      tasks.md
      data-model.md
      research.md
      contracts/
  templates/
    constitution.md
    spec.md
    plan.md
    tasks.md
    data-model.md
    contracts/
      api.md
      events.md
    agents-snippet.md
    prompts/
      constitution.md
      spec.md
      plan.md
      tasks.md
      implement.md
  scripts/
    inspect-spec.sh
    check-constitution.sh
    check-spec-ready.sh
    check-plan-ready.sh
    check-tasks-ready.sh
    check-implement-ready.sh
    list-open-tasks.sh
    link-agents.sh
    list-specs.sh
    show-spec.sh
AGENTS.md
```

## Phase model

The intended agent workflow is strict:

1. `constitution`
2. `spec`
3. `inspect`
4. `plan`
5. `tasks`
6. `implement`
7. `archive`

Dependency rules:

- `constitution` can be created first
- `spec` depends on the constitution
- `inspect` depends on the constitution and one spec, then conditionally loads deeper artifacts only when required
- `plan` depends on the constitution and one spec
- `tasks` depends on the constitution and one plan package, then conditionally loads deeper artifacts only when required
- `implement` depends on the constitution and one task list, then conditionally loads deeper artifacts only when required
- `archive` depends on one existing spec and archives the related plan package when present

## Language model

`draftspec init` supports a compact language configuration, an optional target path, and optional agent-target generation.

Defaults:

- `default language`: `en`
- supported values: `en`, `ru`

Controls:

- `docs language`: generated project docs and templates
- `agent language`: generated prompts and inserted `AGENTS.md` guidance
- `comments language`: preferred code comment language recorded in policy docs and config
- `shell`: generated workflow script family; supported values are `sh` and `powershell`

The language settings are stored in `.draftspec/draftspec.yaml` and reflected in:

- `.draftspec/constitution.md`
- `.draftspec/templates/agents-snippet.md`

For specification and planning work, the configured `docs language` acts as the default language for generated specs, plans, contracts, and task lists unless an existing artifact has a stronger local convention that should be preserved.

For implementation work, the configured `comments language` acts as the default language for new or edited code comments unless an existing file has a stronger local convention that should be preserved.

## Token efficiency goals

Draftspec should stay meaningfully lighter than Speckit by default.

Draftspec should also stay team-safe by default:

- each feature should be worked in a dedicated git branch
- the default feature branch naming convention should be `feature/<slug>`
- active feature state should live in the feature spec and plan package, not in a shared mutable memory file
- archive should preserve historical context without creating frequent merge conflicts across parallel work

Design constraints for low token usage:

- phase prompts should be short and explicit
- each phase should read only one feature package at a time
- optional artifacts should stay optional
- readiness scripts should enforce prerequisites instead of pushing that work into the model context
- patch existing files instead of regenerating large documents
- avoid loading unrelated feature plans, tasks, or contracts
- use `plan.md` as the tasks entrypoint
- use `tasks.md` as the implement entrypoint

Lightweight guardrails should preserve strictness without broadening default context:

- each phase should define `always load`, `load if needed`, and `never load by default` inputs
- `implement` should stay task-scoped by default and only open deeper artifacts when the active task requires them
- `verify` should stay cheap by default and only deepen into code or wider review when explicitly requested
- helper scripts and readiness checks should be preferred over repeated prompt-time reasoning for prerequisite validation
- traceability should improve through stable IDs and explicit references instead of new shared summary artifacts
- `archive` should remain a compact historical record, not a new mutable working-memory layer

Recommended stable ID scheme:

- `RQ-*` for spec requirements
- `AC-*` for acceptance criteria
- `DEC-*` for plan-level implementation decisions
- `DM-*` for data-model entities when the model needs explicit identifiers
- `API-*` for API contract entries
- `EVT-*` for event contract entries
- `T<phase>.<index>` for implementation tasks

## Plan package

Each feature plan lives under `.draftspec/plans/<slug>/`.

Required artifacts:

- `plan.md`
- `tasks.md`
- `data-model.md`
- `contracts/`

Optional artifact:

- `research.md`

`research.md` should only be created when there is genuine uncertainty, external investigation, or architectural tradeoff analysis that needs to be preserved.

## Constitution workflow

`constitution` is agent-driven and strict.

Mandatory sections:

- `Purpose`
- `Core Principles`
- `Constraints`
- `Language Policy`
- `Development Workflow`
- `Governance`
- `Last Updated`

The constitution is authoritative over specs, plans, tasks, and implementation.

## Inspect workflow

`inspect` is agent-driven.

Inputs:

- `.draftspec/constitution.md`
- `.draftspec/specs/<slug>.md`
- optional plan artifacts when they exist, with the cheapest scope preferred first:
  - `plan.md` before `tasks.md`
  - `tasks.md` before `data-model.md`, `contracts/`, or `research.md`
  - do not read implementation code by default

Outputs:

- a focused inspection report for one feature
- explicit Given/When/Then acceptance criteria, with `Given`, `When`, and `Then` kept canonical across documentation languages and inspect treating missing G/W/T as an error
- explicit cheap checks for `constitution <-> spec`, `spec <-> plan`, and `plan <-> tasks` when those downstream artifacts exist
- explicit references to stable IDs when reporting mismatches or gaps

## Archive workflow

`archive` is agent-driven.

Inputs:

- `.draftspec/specs/<slug>.md`
- optional plan artifacts when they exist

Outputs:

- `.draftspec/archive/<slug>/<YYYY-MM-DD>/summary.md`
- archived copies of the feature spec and plan artifacts

## Spec workflow

`spec` is agent-driven.

Inputs:

- `.draftspec/constitution.md`
- user request
- minimal repository context when needed

For agent-facing `/draftspec.spec`, Draftspec should support three optional command-style inputs:

- `--name <feature name>`
- `--slug <feature-slug>`
- `--branch <branch-name>`

Semantics:

- `--name` sets the canonical feature name for the current spec request
- `--slug` overrides the derived spec slug
- `--branch` overrides the execution branch name without changing the spec slug

Draftspec should support two input modes for `/draftspec.spec`:

- inline mode: the user provides the feature name and description in the same request
- staged mode: the user first provides `/draftspec.spec --name ...` and then provides the feature description in the next message

Priority rules for spec identity:

1. `--slug`
2. top-of-file `slug:`
3. slug derived from `--name`
4. slug derived from top-of-file `name:`
5. safe fallback from filename or short request text only when sufficiently specific

Priority rules for feature name:

1. `--name`
2. top-of-file `name:`
3. a concise feature name derived from the user request when unambiguous

If the request is file-based, Draftspec should still prefer explicit top-of-file metadata:

- `name: <feature name>`
- optional `slug: <feature-slug>`

If `/draftspec.spec` is invoked with `--name` but without enough feature detail to write a valid specification, Draftspec should ask for the missing description or treat the next user message as the continuation of the same spec request.

Output:

- `.draftspec/specs/<slug>.md`
- work should happen from `feature/<slug>` when the environment can create or switch branches

If the user explicitly provides `--branch <name>`, Draftspec should use that branch name instead of the default `feature/<slug>` without changing the spec slug.

The spec document itself should stay branch-agnostic. The working branch belongs to execution context, not to the persisted feature specification.

The default spec template should stay compact while making downstream traceability cheap:

- requirements should use stable `RQ-*` IDs
- acceptance criteria should use stable `AC-*` IDs
- acceptance criteria should remain in canonical Given/When/Then form
- out-of-scope boundaries should stay explicit

## Plan workflow

`plan` is responsible for translating one spec into technical design artifacts.

Inputs:

- `.draftspec/constitution.md`
- `.draftspec/specs/<slug>.md`
- repository code and docs when relevant

Outputs:

- `.draftspec/plans/<slug>/plan.md`
- `.draftspec/plans/<slug>/data-model.md`
- `.draftspec/plans/<slug>/contracts/api.md`
- `.draftspec/plans/<slug>/contracts/events.md`
- optional `.draftspec/plans/<slug>/research.md`

The default plan package should support downstream work with minimal rereading:

- `plan.md` should record significant implementation decisions with stable `DEC-*` IDs
- `plan.md` should reference relevant `AC-*` IDs for acceptance-critical behavior
- `data-model.md` should use stable `DM-*` IDs only when entity-level traceability is actually helpful
- contracts should use stable `API-*` and `EVT-*` IDs only when those boundaries exist

## Tasks workflow

`tasks` always reads:

- `.draftspec/constitution.md`
- `.draftspec/plans/<slug>/plan.md`

It then reads only as needed:

- `.draftspec/specs/<slug>.md` when intent or scope boundaries are unclear
- `.draftspec/plans/<slug>/data-model.md` when decomposition depends on entities or invariants
- `.draftspec/plans/<slug>/contracts/` when work crosses API or event boundaries
- `.draftspec/plans/<slug>/research.md` only when present and needed

It must:

- produce concrete executable tasks
- group them by implementation phase
- assign phase-scoped task IDs such as `T1.1`
- map each `AC-*` to at least one task in an explicit acceptance coverage section
- stop when the plan is underspecified or blocked by the constitution

## Implement workflow

`implement` always reads:

- `.draftspec/constitution.md`
- `.draftspec/plans/<slug>/tasks.md`

It then reads only as needed:

- `.draftspec/specs/<slug>.md` when task intent or acceptance scope is unclear
- `.draftspec/plans/<slug>/plan.md` when architectural strategy or sequencing is needed
- `.draftspec/plans/<slug>/data-model.md` when data shape or invariants matter
- `.draftspec/plans/<slug>/contracts/` when APIs or events are involved
- `.draftspec/plans/<slug>/research.md` only when present and required by the task

It must:

- execute only unfinished tasks
- respect task order and phase structure
- keep full-run execution as the default when no scope restriction is provided
- allow explicit scope restriction by one phase or specific task IDs
- reject ambiguous scope such as mixing phase and task selection in the same run
- update `tasks.md`
- emit short phase progress updates during runtime, using the configured agent language
- report implementation coverage in terms of completed task IDs and referenced `AC-*` IDs
- stop when the plan is insufficient or conflicts with the constitution

## Configuration file

```yaml
version: 1

project:
  name: my-project
  constitution_file: .draftspec/constitution.md

runtime:
  shell: sh

paths:
  specs_dir: .draftspec/specs
  plans_dir: .draftspec/plans
  templates_dir: .draftspec/templates
  scripts_dir: .draftspec/scripts

language:
  default: en
  docs: en
  agent: en
  comments: en

agents:
  update_agents_md: true
  agents_file: AGENTS.md
  targets: []

templates:
  spec: spec.md
  plan: plan.md
  tasks: tasks.md
  data_model: data-model.md
  contracts_api: contracts/api.md
  contracts_events: contracts/events.md
  constitution: constitution.md
  constitution_prompt: prompts/constitution.md
  spec_prompt: prompts/spec.md
  plan_prompt: prompts/plan.md
  tasks_prompt: prompts/tasks.md
  implement_prompt: prompts/implement.md

scripts:
  inspect_spec: inspect-spec.sh
  check_constitution: check-constitution.sh
  check_spec_ready: check-spec-ready.sh
  check_plan_ready: check-plan-ready.sh
  check_tasks_ready: check-tasks-ready.sh
  check_implement_ready: check-implement-ready.sh
  list_open_tasks: list-open-tasks.sh
  link_agents: link-agents.sh
  list_specs: list-specs.sh
  show_spec: show-spec.sh
```


`cleanup-agents` removes orphaned agent artifacts for targets that are no longer enabled in config.

`doctor` reports `error` for missing required files and `warning` for orphaned agent artifacts that remain on disk after a target has been disabled in config.

For PowerShell projects, the same generated script names should use `.ps1` extensions instead of `.sh`.
