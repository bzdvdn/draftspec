# draftspec

Русская версия: [README.ru.md](README.ru.md)

`draftspec` is a lightweight project context kit for development agents and humans.

It keeps project intent, specifications, plan artifacts, and task breakdowns in simple files without introducing a rigid process engine.

The first release is intentionally optimized for low overhead and real-world usage: narrow default context, minimal required artifacts, strict workflow discipline without heavyweight orchestration, and branch-based collaboration that works cleanly for both solo and team development.

## Positioning

Draftspec is a strict lightweight SDD kit for real codebases.

It is designed for teams that want more discipline than a loose planning layer, but do not want the workflow surface, artifact overhead, or orchestration weight of a heavier SDD system.

- stricter than OpenSpec in phase discipline and artifact alignment
- lighter than Spec Kit in default context, workflow surface, and artifact overhead
- optimized for agent-first workflows with narrow context loading
- designed to keep strictness in templates, entrypoints, and readiness checks rather than heavyweight orchestration
- built for brownfield repositories where context must stay narrow, local, and reviewable

In short: Draftspec aims to maximize discipline per token: strong phase boundaries, low artifact drag, and enough structure to keep agents and humans aligned in everyday work.

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
- `inspect` is a real quality gate, not a loose suggestion before planning.
- Lightweight traceability. Stable IDs and cheap readiness checks reduce prompt bloat.
- Brownfield-friendly workflow. Draftspec works well in existing repositories without forcing a heavyweight process layer.
- Branch-first collaboration. Active feature state stays local to the feature instead of spreading through shared mutable memory.
- Inspect is mandatory before planning. Each feature should persist an inspect report that confirms the spec is ready for plan work.
- Optional workflow commands available at any phase: `/draftspec.challenge` (adversarial review — finds weak assumptions and untestable criteria), `/draftspec.handoff` (compact session handoff document so a new session can resume without re-reading all artifacts), `/draftspec.hotfix` (emergency fix outside the standard phase chain — for well-understood fixes touching ≤ 3 files), `/draftspec.scope` (quick scope boundary check, inline only, no file written).

OpenSpec is more flexible by design and works well when teams want a looser artifact-guided workflow.

Spec Kit provides a broader and more thorough workflow surface, but usually at the cost of more artifacts, more context, and more process overhead.

Draftspec is optimized for discipline per token: strong workflow boundaries, minimal default context, explicit quality gates, and enough structure to keep agents aligned without making the workflow heavy.

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
draftspec check <slug> [path]
draftspec check <slug> [path] --json
draftspec check [path] --all
draftspec check [path] --all --json
draftspec demo [path]
draftspec export <slug> [path]
draftspec export <slug> [path] --output <file>
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
- `draftspec init` requires an explicit `--shell` and generates one script family: `sh` or `powershell`. Supported agent targets: `claude`, `codex`, `copilot`, `cursor`, `kilocode`, `trae`, `windsurf`, `roocode`, `aider`.
- Generated workspaces include `.draftspec/scripts/run-draftspec.*` as the stable CLI launcher for agents; it resolves `DRAFTSPEC_BIN` first and falls back to `draftspec` from `PATH`.
- `draftspec feature repair` and `draftspec migrate` provide safe canonicalization for legacy artifacts such as old inspect report paths.
- `draftspec check <slug>` shows artifact presence, inspect and verify verdict, task progress, and the exact next slash command; exits with code 1 when blocked; supports `--json` for CI use. `--all` shows a readiness table across all features.
- `draftspec demo [path]` creates a demo workspace pre-populated with an example feature at the implement phase — spec, inspect report, plan, tasks, and data model are all populated.
- `draftspec export <slug>` bundles all feature artifacts into one markdown document for sharing with a reviewer or new agent session; supports `--output` to write to a file.
- `/draftspec.plan` supports `--research`: enters research-first mode — agent identifies concrete unknowns, writes `research.md`, then asks "Research complete — proceed to full plan?" before producing `plan.md`.
- `/draftspec.spec` supports `--amend`: targeted edit mode — update one section or add one criterion without rewriting the spec or invalidating an existing inspect report.
- `/draftspec.handoff` without a slug generates handoff documents for all active features at once.
- `/draftspec.hotfix`: emergency fix workflow — writes a minimal hotfix spec (fix, root cause, risk, verification, touches) before any code change, implements, verifies inline, then archives; skips inspect, plan, and tasks phases; use only when the root cause is known and the fix touches ≤ 3 files.
- `doctor` warns when the same stable ID (`AC-*`, `RQ-*`) appears across multiple specs.
- Generated docs and prompts support English and Russian.

## Quick Example

```bash
# try the demo instantly — no project setup required
draftspec demo ./my-demo

# init a real project
draftspec init my-project --lang en --shell sh --agents claude --agents codex
draftspec refresh my-project --shell powershell --dry-run
draftspec doctor my-project
draftspec check export-report my-project
```

## Example Feature Cycle

A full workflow for a real feature — "Add CSV export to the reports page".

<details>
<summary>See the full cycle →</summary>

### 1. Init

```bash
draftspec init . --lang en --shell sh --agents claude
# → .draftspec/ scaffold, AGENTS.md, agent slash-command files written
```

### 2. Write the spec

Call `/draftspec.spec --name "CSV export for reports"` in your agent.

`.draftspec/specs/csv-export-for-reports/spec.md`:

```markdown
## Goal
Allow users to download the reports table as a CSV file.

## Acceptance Criteria

**AC-001** Export produces a file
Given the Reports page has at least one row
When the user clicks "Export CSV"
Then a .csv file downloads with column headers and all visible rows

**AC-002** Empty state is handled
Given the reports table is empty
When the user clicks "Export CSV"
Then a .csv with headers only downloads — no error shown
```

### 3. Inspect

Call `/draftspec.inspect csv-export-for-reports`.

- `.draftspec/specs/csv-export-for-reports/inspect.md` — verdict `pass`, all AC have G/W/T
- `.draftspec/specs/csv-export-for-reports/summary.md` — compact AC table used by implement and verify instead of the full spec

### 4. Plan

Call `/draftspec.plan csv-export-for-reports`.

`.draftspec/plans/csv-export-for-reports/plan.md` names the implementation surfaces: `ReportsPage.tsx` (add button), `useReportExport.ts` (new hook, CSV logic), `reports.test.ts` (new tests).

### 5. Tasks

Call `/draftspec.tasks csv-export-for-reports`.

`.draftspec/plans/csv-export-for-reports/tasks.md`:

```markdown
## Surface Map
| Surface                        | Tasks      |
|-------------------------------|------------|
| hooks/useReportExport.ts       | T1.1       |
| components/ReportsPage.tsx     | T1.2       |
| tests/reports.test.ts          | T2.1       |

## Phase 1: Hook and button

- [ ] T1.1 add `useReportExport` hook — converts `rows[]` to CSV blob and triggers browser download — AC-001  `Touches: hooks/useReportExport.ts`
- [ ] T1.2 add Export CSV button to ReportsPage — calls hook on click, disabled when rows empty — AC-001, AC-002  `Touches: components/ReportsPage.tsx`

## Phase 2: Tests

- [ ] T2.1 add tests for useReportExport — covers non-empty rows, empty rows, header-only output — AC-001, AC-002  `Touches: tests/reports.test.ts`

## Acceptance Coverage
AC-001 → T1.1, T1.2, T2.1
AC-002 → T1.2, T2.1
```

### 6. Implement, verify, archive

```
/draftspec.implement csv-export-for-reports   # Phase 1 done, stops
/draftspec.implement csv-export-for-reports   # Phase 2 done, stops
/draftspec.verify    csv-export-for-reports   # verdict: pass
/draftspec.archive   csv-export-for-reports
```

### Check readiness at any point

```bash
draftspec check csv-export-for-reports
# Phase:  tasks → implement
# Tasks:  0 / 3 done
# Next:   /draftspec.implement csv-export-for-reports
```

</details>

## Demo

A reproducible terminal demo kit lives under [`demo/`](demo/README.md).

Build the local binary and render the quick terminal demo:

```bash
go build -o bin/draftspec ./src/cmd/draftspec
vhs demo/quick.tape
```

Demo assets:

- [Quick terminal demo](demo/README.md)
- [Brownfield walkthrough](demo/brownfield.md)
- [Self-hosting walkthrough](demo/self-hosting.md)

The quick tape produces `demo/draftspec-demo.gif` and demonstrates `init`, generated agent files, `AGENTS.md`, and launcher-based `doctor` / `refresh --dry-run`.

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

Project contribution and trust docs:

- [Contributing](CONTRIBUTING.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Security Policy](SECURITY.md)


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
