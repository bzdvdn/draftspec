# CLI Reference

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
- `warning` for non-standard Git branch names
- `ok` when the workspace is healthy

Use `--json` for machine-readable output in automation and CI.

### `draftspec dashboard [path]`

Displays a visual dashboard of all active features in the project.

The dashboard includes:

- Feature slug
- Current workflow phase
- Implementation progress percentage
- Status (READY/BLOCKED)
- Current Git branch (with `!!` warning if there is a mismatch with the feature slug)

```bash
draftspec dashboard
```

### `draftspec feature <slug> [path]`

Shows a detailed workflow view for one feature.

The text view includes:

- current phase and `ready_for`
- inspect and verify status when reports exist
- task progress when `tasks.md` exists
- grouped workflow findings
- a short `focus` hint for the next likely action

Use `--json` to return structured state plus feature-local findings.

### `draftspec feature repair <slug> [path]`

Repairs safe feature-local Draftspec issues.

Current repair scope includes legacy inspect report migration from:

- `.draftspec/plans/<slug>/inspect.md`

to the canonical path:

- `.draftspec/specs/<slug>/inspect.md`

Use `--dry-run` to preview changes and `--json` for structured output.

### `draftspec features [path]`

Lists workflow status across all discovered features.

The text view summarizes:

- phase and `ready_for`
- inspect and verify verdicts
- task progress
- grouped issue counts
- artifact presence

Use `--json` for machine-readable output.

### `draftspec migrate [path]`

Runs safe project-wide Draftspec migrations.

Current migration scope focuses on canonicalizing legacy inspect report paths across the project.

### `draftspec list-specs [path]`

Lists spec slugs from `.draftspec/specs/`.

### `draftspec show-spec <name> [path]`

Prints one spec file by slug.

### `draftspec check <slug> [path]`

Shows feature readiness and the exact next action for one feature.

Output includes artifact presence, inspect and verify verdict, task progress, the exact next slash command, and a compact structured-check summary when phase-specific readiness checks produce categorized findings.

Use `--all` to check every feature in one table. Exits with code 1 when any feature is blocked.
Use `--json` for machine-readable output suitable for CI, including `check_summary` and `check_findings` when available.

```bash
draftspec check export-report
draftspec check export-report my-project --json
draftspec check my-project --all
draftspec check my-project --all --json
```

### `draftspec trace [slug] [path]`

Scans for traceability annotations in the codebase.

Annotations follow the format:
- `// @ds-task <TASK_ID>: <Description> (<AC_ID>)` for implementation code.
- `// @ds-test <TASK_ID>: <TestName> (<AC_ID>)` for test evidence.

This command identifies links between implementation code, task IDs from `tasks.md`, and acceptance criteria from `spec.md`.

Use `slug` to filter findings for a specific feature.
Use `--tests` to show only test evidence.
Use `--json` for machine-readable output.

```bash
draftspec trace
draftspec trace export-report
draftspec trace export-report --tests
draftspec trace export-report my-project --json
```

### `draftspec demo [path]`

Creates a demo workspace at the given path (default: `./draftspec-demo`).

The workspace is pre-populated with an example feature (`export-report`) at the implement phase — spec, inspect report, plan, tasks, and data model are all present. Suggests `/draftspec.scope`, `/draftspec.challenge`, and `/draftspec.handoff` to try immediately.

```bash
draftspec demo
draftspec demo ./my-demo --agents claude
```

### `draftspec export <slug> [path]`

Bundles all artifacts for one feature into a single markdown document.

Reads and concatenates: spec, inspect report, plan, tasks, data model, research, challenge report, and verify report (skips missing files). Useful for sharing full feature context with a reviewer or a new agent session.

Use `--output <file>` to write to a file instead of stdout.

```bash
draftspec export export-report
draftspec export export-report my-project --output export-report-bundle.md
```
