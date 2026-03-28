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

`verify` is intentionally lightweight:

- it starts from `tasks.md`
- it can use `.draftspec/scripts/verify-task-state.sh <slug>` as a cheap first-pass helper
- it can use `.draftspec/scripts/verify-memory-sync.sh <slug>` for coarse memory/task sync signals before touching code
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
