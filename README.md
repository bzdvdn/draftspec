# draftspec

Русская версия: [README.ru.md](README.ru.md)

`draftspec` is a lightweight project context kit for development agents and humans.

It keeps project intent, specifications, plan artifacts, and task breakdowns in simple files without introducing a rigid process engine.

The first release is intentionally optimized for low overhead and real-world usage: narrow default context, minimal required artifacts, strict workflow discipline without heavyweight orchestration, and branch-based collaboration that works cleanly for both solo and team development.

## Positioning

Draftspec is a lean SDD kit for real codebases.

- stricter than OpenSpec in phase discipline and artifact alignment
- lighter than Spec Kit in default context, workflow surface, and artifact overhead
- optimized for agent-first workflows with narrow context loading
- designed to keep strictness in templates, entrypoints, and readiness checks rather than heavyweight orchestration

In short: Draftspec aims to be as strict as practical while staying lightweight enough for everyday use.

## Draftspec vs OpenSpec vs Spec Kit

| Dimension | Draftspec | OpenSpec | Spec Kit |
| --- | --- | --- | --- |
| Workflow style | Strict phase chain with narrow context | Fluid artifact-guided workflow | Thorough multi-step SDD workflow |
| Default context size | Smallest by default | Moderate | Largest |
| Artifact overhead | Low | Medium | High |
| Phase discipline | High | Medium | Highest |
| Brownfield ergonomics | High | High | Medium |
| Team collaboration model | Branch-first, feature-local artifacts | Change-folder oriented | Branch and workflow heavy |
| Shared mutable state | Avoided by design | Low | Varies by setup |
| Best fit | Lean strict SDD on real codebases | Flexible SDD-lite for fast iteration | Full-featured rigorous SDD |

In short, Draftspec aims to sit between OpenSpec and Spec Kit: stricter than OpenSpec, lighter than Spec Kit, and optimized for branch-based collaboration with minimal default context.

## Documentation

Extended documentation lives in [`docs/`](docs/README.md):

- [English docs](docs/en/index.md)
- [Русская документация](docs/ru/index.md)

## Public CLI

```text
draftspec init [path]
draftspec refresh [path]
draftspec add-agent [path]
draftspec list-agents [path]
draftspec remove-agent [path]
draftspec cleanup-agents [path]
draftspec doctor [path]
draftspec doctor [path] --json
draftspec list-specs [path]
draftspec show-spec <name> [path]
```

## Workflow

```text
constitution -> spec -> inspect -> plan -> tasks -> implement -> verify -> archive
```

## Key Points

- The constitution is the highest-priority project document.
- Plan packages keep `plan.md`, `tasks.md`, `data-model.md`, `contracts/`, and optional `research.md` together.
- Specs use canonical `Given / When / Then` markers across documentation languages.
- Agent workflows are designed to load only the minimum context required.
- Strictness comes from phase entrypoints, templates, and readiness checks rather than large default prompts.
- Agent-facing `/draftspec.spec` is branch-first: it should work from `feature/<slug>` and prefer explicit `name:` / `slug:` metadata for prompt files.
- `draftspec init` requires an explicit `--shell` and generates one script family: `sh` or `powershell`.
- Generated docs and prompts support English and Russian.

## Quick Example

```bash
draftspec init my-project --lang en --shell sh --agents claude --agents codex
draftspec refresh my-project --shell powershell --dry-run
draftspec doctor my-project
draftspec doctor my-project --json
```

For deeper guidance, use:

- [Overview](docs/en/overview.md)
- [CLI Reference](docs/en/cli.md)
- [Workflow Model](docs/en/workflow.md)
- [Examples](docs/en/examples.md)
- [Roadmap](docs/en/roadmap.md)


## Development

```bash
go test ./...
go build -o bin/draftspec ./src/cmd/draftspec
```

The repository includes unit tests for config, project lifecycle, doctor checks, specs, templates, agents, and CLI-level behavior.

## License

Released under the [MIT License](LICENSE).
