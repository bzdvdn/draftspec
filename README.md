# draftspec

`draftspec` is a lightweight project context kit for development agents and humans.

It keeps project intent, specifications, plan artifacts, task breakdowns, and working memory in simple files without introducing a rigid process engine.

## Documentation

Extended documentation lives in [`docs/`](docs/README.md):

- [English docs](docs/en/index.md)
- [Русская документация](docs/ru/index.md)

## Public CLI

```text
draftspec init [path]
draftspec add-agent [path]
draftspec list-agents [path]
draftspec remove-agent [path]
draftspec cleanup-agents [path]
draftspec doctor [path]
draftspec list-specs [path]
draftspec show-spec <name> [path]
```

## Workflow

```text
constitution -> spec -> inspect -> plan -> tasks -> implement -> archive
```

## Key Points

- The constitution is the highest-priority project document.
- Plan packages keep `plan.md`, `tasks.md`, `data-model.md`, `contracts/`, and optional `research.md` together.
- Specs use canonical `Given / When / Then` markers across documentation languages.
- Agent workflows are designed to load only the minimum context required.
- Generated docs and prompts support English and Russian.

## Quick Example

```bash
draftspec init my-project --lang en --agents claude --agents codex
draftspec doctor my-project
```

For deeper guidance, use:

- [Overview](docs/en/overview.md)
- [CLI Reference](docs/en/cli.md)
- [Workflow Model](docs/en/workflow.md)
- [Examples](docs/en/examples.md)


## Development

```bash
go test ./...
go build -o bin/draftspec ./src/cmd/draftspec
```

The repository includes unit tests for config, project lifecycle, doctor checks, specs, templates, agents, and CLI-level behavior.

## License

Released under the [MIT License](LICENSE).

