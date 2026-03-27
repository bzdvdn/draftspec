## Draftspec Self-Hosting

This repository builds the `draftspec` CLI itself.

Generated `/.draftspec/` and `/AGENTS.md` files in this repository are local test artifacts used to validate `draftspec init`. They are not product source files and should not be committed.

## Read Order

Before making meaningful changes, read:

1. `README.md`
2. `MVP.md`
3. relevant Go files under `src/cmd/` and `src/internal/`

If local self-hosted Draftspec files exist, use them only as runtime test output, not as the canonical product spec.

## Workflow Chain

The intended agent flow is:

1. `/draftspec.constitution`
2. `/draftspec.spec`
3. `/draftspec.plan`
4. `/draftspec.tasks`
5. `/draftspec.implement`

Dependency rules:

- `plan` must not run before `spec`
- `tasks` must not run before `plan`
- `implement` must not run before `tasks`
- every phase must comply with the constitution

## Language Model

`draftspec init` now supports generated language settings.

- `--lang` sets the base language and defaults to `en`
- `--docs-lang` controls generated docs and templates
- `--agent-lang` controls prompts and inserted `AGENTS.md` guidance
- `--comments-lang` records the preferred code comment language
- supported values are `en` and `ru`
- `spec`, `plan`, and `tasks` prompts should keep generated docs in the configured documentation language
- `implement` should keep new or edited code comments in the configured comments language unless a stronger local convention applies

When changing this behavior, keep config, generated docs, prompts, and scripts aligned.

## Token Discipline

Keep Draftspec light.

- Load only the current feature slug.
- Do not read unrelated plans or specs.
- Prefer scripts for readiness checks.
- Prefer patching existing docs over full rewrites.
- Treat `research.md` as optional.
- Use `plan.md` as the tasks entrypoint.
- Use `tasks.md` as the implement entrypoint.
- Load `spec`, `data-model`, and `contracts` during tasks only when the decomposition requires them.
- Load `spec`, `plan`, `data-model`, and `contracts` during implement only when the active task requires them.

## Working On The CLI

When changing CLI behavior:

- update the Go implementation
- keep embedded template assets aligned with the behavior
- update `README.md` and `MVP.md` when the product contract changes
- verify generated output from `draftspec init` when relevant

## Generated Output Policy

When testing `draftspec` in this repository:

- `/.draftspec/` is generated output
- `/AGENTS.md` is generated output
- these files may be refreshed during testing
- they should not be treated as the source of truth for product design

The source of truth lives in:

- `README.md`
- `MVP.md`
- Go source files
- embedded assets under `src/internal/templates/assets/`
