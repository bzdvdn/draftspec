# draftspec

`draftspec` is a lightweight project context kit for development agents and humans.

The goal is to keep project intent, feature specs, plans, task breakdowns, and working memory in one simple place without introducing a rigid workflow.

## MVP goals

- Keep the public CLI intentionally small
- Store project context in plain Markdown files
- Support development agents through templates and shell scripts
- Treat the constitution as the highest-priority project document
- Encode development workflow rules directly into the constitution
- Support plan artifacts such as `data-model.md` and `contracts/`
- Maintain project `memory.md` and keep it aligned with constitutional and implementation changes
- Keep token usage practical by making each workflow read only the minimum required artifacts
- Support English and Russian generated docs and prompts
- Avoid checkpoints, quickstarts, approvals, and heavy state machines

## Public CLI

The MVP exposes only a few commands:

```text
draftspec init
draftspec list-specs
draftspec show-spec <name>
```

## `init` language options

`draftspec init` supports a small language model for generated artifacts.

Flags:

```text
draftspec init --lang en
draftspec init --lang ru
draftspec init --lang en --docs-lang ru --agent-lang en --comments-lang en
```

Rules:

- `--lang` sets the base language for generated artifacts and defaults to `en`
- `--docs-lang` controls generated project docs such as `constitution.md`, `memory.md`, spec templates, and plan templates
- `--agent-lang` controls generated prompts and the inserted `AGENTS.md` guidance snippet
- `--comments-lang` records the preferred code comment language in config and generated policy docs

Generated language settings are written to `.draftspec/draftspec.yaml` and echoed in the constitution and memory files. During `implement`, agents should use the configured comment language for new or edited code comments unless an existing file has a different local convention that should be preserved.

## Workflow chain

The intended dependency chain is strict:

```text
constitution -> spec -> inspect -> plan -> tasks -> implement -> archive
```

Rules:

- `constitution` may be created first
- `spec` should be created against the constitution
- `plan` must not run without an existing spec
- `tasks` must not run without an existing plan package
- `implement` must not run without existing tasks
- `archive` should run only when a feature is complete, superseded, rejected, abandoned, or deferred

## Feature package model

Each feature uses two layers:

- `.draftspec/specs/<slug>.md`
- `.draftspec/plans/<slug>/...`

The plan package contains:

- `plan.md`
- `tasks.md`
- `data-model.md`
- `contracts/`
- optional `research.md`

This keeps spec intent separate from technical planning while keeping all implementation-oriented artifacts together.

## Agent-facing workflows

These are not public CLI commands in the first version. They are driven by agent rules, templates, and shell scripts.

- `constitution`
- `spec`
- `inspect`
- `plan`
- `tasks`
- `implement`
- `archive`

Each phase should read only the files it truly needs.

Recommended minimal read sets:

- `constitution`: constitution template, memory, essential repo context
- `spec`: constitution, memory, relevant user request, minimal repo context, while keeping generated spec text in the configured documentation language
- `plan`: constitution, memory, one spec, minimal repo/code context, while keeping all generated planning artifacts in the configured documentation language
- `tasks`: constitution, memory, plan first; spec, data model, contracts, and research only when required by the decomposition, while keeping the task list in the configured documentation language
- `implement`: constitution, memory, tasks first; spec, plan, data model, contracts, and research only when required by the active task
- `archive`: memory, spec, and existing plan artifacts for one slug only

This is the main way to keep Draftspec less token-hungry than Speckit.

Generated specs should use explicit `Given/When/Then` acceptance criteria. These BDD markers remain canonical regardless of `docs-lang`, while the surrounding explanatory text follows the configured documentation language. Generated tasks should include a simple acceptance-to-task coverage section so `inspect` can enforce traceability without inventing structure later.

## Token usage

What usually makes systems like `speckit` expensive is that each command pulls in many large artifacts by default.

Draftspec should stay lighter by design:

- no mandatory `quickstart.md`
- no mandatory broad context fan-out for every phase
- no large orchestration prompt that always reads everything
- `research.md` is optional
- commands should load only one feature package at a time
- prompts should prefer patching existing files over regenerating whole documents

In practice, the biggest token savers are:

- one feature per run
- one spec slug per run
- no eager loading of unrelated plans
- short prompts with strict inputs
- scripts handling readiness checks instead of the model reasoning from scratch about prerequisites
- `tasks` using `plan.md` as the decomposition entrypoint instead of loading the full feature package immediately
- `implement` using `tasks.md` as the execution entrypoint instead of loading the full feature package immediately

## Developing draftspec itself

Inside this repository, generated `/.draftspec/`, `/AGENTS.md`, and `/TESTS/` are local development artifacts.

They are useful for testing `draftspec`, but they should not be committed as product source files.
