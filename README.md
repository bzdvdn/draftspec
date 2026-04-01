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

## Where Draftspec Stands Out

- Narrow context by default. Each phase is designed to load the smallest useful scope.
- Code reading should stay phase-local and targeted: enough to remove guesswork, not enough to recreate full-repository context.
- Strict workflow chain. Constitution, spec, inspect, plan, tasks, and implementation stay aligned.
- Lightweight traceability. Stable IDs and cheap readiness checks reduce prompt bloat.
- Brownfield-friendly workflow. Draftspec works well in existing repositories without forcing a heavyweight process layer.
- Branch-first collaboration. Active feature state stays local to the feature instead of spreading through shared mutable memory.
- Inspect is mandatory before planning. Each feature should persist an inspect report that confirms the spec is ready for plan work.

OpenSpec is more flexible by design and works well when teams want a looser artifact-guided workflow.

Spec Kit provides a broader and more thorough workflow surface, but usually at the cost of more artifacts, more context, and more process overhead.

Draftspec is optimized for discipline per token: strong workflow boundaries, minimal default context, and enough structure to keep agents aligned without making the workflow heavy.

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
draftspec feature <slug> [path]
draftspec feature repair <slug> [path]
draftspec features [path]
draftspec migrate [path]
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
- `data-model.md` and `contracts/` are intentionally compact but structured: entities should capture fields, invariants, and lifecycle; contracts should capture boundary IO, failures, and delivery assumptions.
- Specs use canonical `Given / When / Then` markers across documentation languages.
- Draftspec prefers stable IDs and explicit references over repeated narrative summaries: `RQ-*` for requirements, `AC-*` for acceptance criteria, `DEC-*` for plan decisions, and phase-scoped `T*` task IDs.
- Agent workflows are designed to load only the minimum context required.
- Strictness comes from phase entrypoints, templates, stable artifact structure, and readiness checks rather than large default prompts.
- Agent-facing `/draftspec.spec` is branch-first: it should work from `feature/<slug>`, support `--name` with optional `--slug` / `--branch`, and still prefer explicit `name:` / `slug:` metadata for prompt files.
- `draftspec init` requires an explicit `--shell` and generates one script family: `sh` or `powershell`.
- Generated workspaces include `.draftspec/scripts/run-draftspec.*` as the stable CLI launcher for agents; it resolves `DRAFTSPEC_BIN` first and falls back to `draftspec` from `PATH`.
- `draftspec feature repair` and `draftspec migrate` provide safe canonicalization for legacy artifacts such as old inspect report paths.
- Generated docs and prompts support English and Russian.

## Quick Example

```bash
draftspec init my-project --lang en --shell sh --agents claude --agents codex
draftspec refresh my-project --shell powershell --dry-run
draftspec doctor my-project
draftspec doctor my-project --json
```

## Install

Draftspec is distributed as a single binary via GitHub Releases.

Linux:

```bash
VERSION=v0.1.0
curl -fsSL "https://raw.githubusercontent.com/bzdvdn/draftspec/${VERSION}/scripts/install.sh" | bash -s -- --version "${VERSION}"
```

Windows (PowerShell):

```powershell
$version="v0.1.0"
$env:DRAFTSPEC_VERSION=$version
powershell -ExecutionPolicy Bypass -c "iwr -useb https://raw.githubusercontent.com/bzdvdn/draftspec/$version/scripts/install.ps1 | iex"
```

To also add the install directory to `PATH`:

- Linux: add `--add-to-path` or set `DRAFTSPEC_ADD_TO_PATH=1`
- Windows: set `$env:DRAFTSPEC_ADD_TO_PATH=1` or run the script with `-AddToPath`

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

# with version stamp
go build -ldflags "-X draftspec/src/internal/cli.Version=v0.1.0" -o bin/draftspec ./src/cmd/draftspec
```

The repository includes unit tests for config, project lifecycle, doctor checks, specs, templates, agents, and CLI-level behavior.

## License

Released under the [MIT License](LICENSE).
