# Agents

## Supported Agent Targets

Draftspec can generate project-local command or prompt files for:

- `claude`
- `codex`
- `copilot`
- `cursor`
- `kilocode`
- `trae`
- `all`

## Generated Locations

- Claude: `.claude/commands/`
- Codex: `.codex/prompts/`
- Copilot: `.github/prompts/`
- Cursor: `.cursor/rules/`
- Kilo Code: `.kilocode/rules/`
- Trae: `.trae/project_rules.md`

These generated files are thin wrappers around the canonical Draftspec prompts in `.draftspec/templates/prompts/`.

## Agent Discipline

The agent-facing workflows are:

- `constitution`
- `spec`
- `inspect`
- `plan`
- `tasks`
- `implement`
- `verify`
- `archive`

Each prompt is designed to:

- read only the minimum required context
- stop when prerequisites are missing
- respect the configured documentation and agent languages
- preserve constitutional authority over specs, plans, tasks, and implementation

`spec` should stay branch-first:

- it should create or switch to `feature/<slug>` before writing `.draftspec/specs/<slug>.md` when the environment allows it
- it should support `--name`, optional `--slug`, and optional `--branch` for chat-oriented input
- if `/draftspec.spec` is invoked with `--name` but without enough description, it should preserve context and ask for or accept the next message as the continuation of the spec request
- when the input comes from a local prompt file, it should prefer top-of-file `name:` and optional `slug:` metadata over a generic filename
- if the request is ambiguous, multi-feature, URL-like, or tries to derive one spec from multiple constitutional changes, it should stop and ask for one concrete feature

`verify` is intentionally lightweight:

- it starts from `tasks.md`
- it can use `.draftspec/scripts/verify-task-state.sh <slug>` as a cheap first-pass helper
- it reads deeper artifacts only when needed to confirm a concrete claim
- it is meant to confirm readiness for archive or follow-up refinement, not to become a heavy review engine

For `archive`, a `completed` status should reuse `verify-task-state.sh` before creating the snapshot.

## Maintenance Commands

Use the public CLI to manage agent targets safely:

```bash
draftspec add-agent my-project --agents claude --agents cursor
draftspec list-agents my-project
draftspec remove-agent my-project --agents cursor
draftspec cleanup-agents my-project
draftspec doctor my-project
```
