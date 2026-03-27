# Self-Hosting and Development

## Developing Draftspec Itself

Inside the Draftspec repository, generated `/.draftspec/`, `/AGENTS.md`, and `/TESTS/` are local development artifacts. They are useful for smoke tests and self-hosting experiments, but they are not product source files.

## Recommended Local Workflow

```bash
go test ./...
go build -o bin/draftspec ./src/cmd/draftspec
./bin/draftspec init TESTS/demo --git=false --lang en --agents claude --agents cursor
./bin/draftspec doctor TESTS/demo
```

## Current Test Coverage

The repository already includes unit tests for:

- config loading, defaults, save, and path resolution
- project initialization and agent lifecycle operations
- workspace health checks through `doctor`
- spec listing, reading, and template-based creation
- localized template asset consistency
- agent file generation
- CLI-level command behavior

## Why `doctor` and `cleanup-agents` Matter

When you test multiple agent targets, it is easy to leave stale generated files behind. Draftspec separates these concerns:

- `remove-agent` updates config and removes files for selected enabled targets
- `cleanup-agents` removes leftover artifacts for targets that are no longer enabled
- `doctor` reports missing files as `error` and leftover artifacts as `warning`

## Source of Truth

The main sources of truth in this repository are:

- `src/` for Go code
- `src/internal/templates/assets/lang/` for localized generated assets
- `README.md` for product summary
- `MVP.md` for the current product model
- `docs/` for broader documentation
