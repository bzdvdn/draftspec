# CLI Reference

## Commands

### `draftspec init [path]`

Initializes a Draftspec workspace in the target project.

Examples:

```bash
draftspec init
draftspec init my-project --lang en
draftspec init my-project --docs-lang ru --agent-lang en --comments-lang en --agents claude --agents cursor
```

Important flags:

- `--git` initializes a Git repository when true; default is enabled
- `--lang` sets the base language; default is `en`
- `--docs-lang` sets the generated documentation language
- `--agent-lang` sets the generated prompt and agent guidance language
- `--comments-lang` records the preferred code comment language
- `--agents` generates project-local agent command files

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

### `draftspec list-specs [path]`

Lists spec slugs from `.draftspec/specs/`.

### `draftspec show-spec <name> [path]`

Prints one spec file by slug.
