# CLI Reference

## Commands

### `draftspec init [path]`

Initializes a Draftspec workspace in the target project.

Examples:

```bash
draftspec init
draftspec init my-project --lang en --shell sh
draftspec init my-project --docs-lang ru --agent-lang en --comments-lang en --shell powershell --agents claude --agents cursor
```

Important flags:

- `--git` initializes a Git repository when true; default is enabled
- `--lang` sets the base language; default is `en`
- `--shell` selects the generated workflow script family; required: `sh` or `powershell`
- `--docs-lang` sets the generated documentation language
- `--agent-lang` sets the generated prompt and agent guidance language
- `--comments-lang` records the preferred code comment language
- `--agents` generates project-local agent command files

### `draftspec refresh [path]`

Refreshes only Draftspec-managed generated artifacts in an existing project.

This command updates:

- `.draftspec/draftspec.yaml`
- `.draftspec/templates/**`
- `.draftspec/scripts/**`
- project-local agent command files
- the managed Draftspec guidance block inside `AGENTS.md`

This command does not update:

- `.draftspec/constitution.md`
- `.draftspec/specs/**`
- `.draftspec/plans/**`
- `.draftspec/archive/**`

Examples:

```bash
draftspec refresh my-project
draftspec refresh my-project --shell powershell --agents claude --dry-run
draftspec refresh my-project --agent-lang ru --json
```

Important flags:

- `--lang`, `--docs-lang`, `--agent-lang`, `--comments-lang` override the existing configured languages
- `--shell` overrides the generated workflow script family
- `--agents` overrides enabled project-local agent targets
- `--dry-run` reports pending managed changes without writing them
- `--json` outputs the refresh result as JSON

### `draftspec add-agent [path]`

Adds one or more agent targets to an existing Draftspec project.

```bash
draftspec add-agent my-project --agents claude --agents codex
```

### `draftspec list-agents [path]`

Lists enabled agent targets from `.draftspec/draftspec.yaml`.

### `draftspec remove-agent [path]`

Disables one or more agent targets and removes their generated files.

### `draftspec cleanup-agents [path]`

Removes orphaned agent artifacts that no longer match enabled targets in config.

### `draftspec doctor [path]`

Checks workspace health.

`doctor` reports:

- `error` for missing required files or invalid config values
- `warning` for orphaned agent artifacts still present on disk
- `ok` when the workspace is healthy

Use `--json` for machine-readable output in automation and CI.

### `draftspec list-specs [path]`

Lists spec slugs from `.draftspec/specs/`.

### `draftspec show-spec <name> [path]`

Prints one spec file by slug.
