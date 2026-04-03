# Changelog

All notable changes to this project will be documented in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Versions follow [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- `draftspec check --all [path]`: aggregate readiness table across all features; exits with code 1 when any feature is blocked; supports `--json`
- `draftspec export <slug> [path]`: bundles all artifacts for one feature (spec, inspect report, plan, tasks, data model, research, challenge, verify) into a single markdown document; use `--output` to write to a file
- `doctor` now warns when the same stable ID (`AC-*`, `RQ-*`) appears in multiple specs — cross-spec ID collision detection
- `/draftspec.spec --amend`: targeted edit mode — update one section or add one criterion without rewriting the entire spec; does not invalidate an existing inspect report unless the change materially affects acceptance criteria or scope
- `/draftspec.handoff` without slug: generates handoff documents for all active features at once; run `list-specs` to enumerate, write one file per feature, output a summary table
- Optional workflow commands available to all agent targets:
  - `/draftspec.handoff`: generate a compact session handoff document so a new session can resume without re-reading all artifacts; always overwrites the previous snapshot
  - `/draftspec.challenge`: adversarial review of a spec or plan — finds weak assumptions, untestable acceptance criteria, and scope drift before implementation; supports `--spec` and `--plan` flags; verdict: `strong`, `concerns`, or `fragile`
  - `/draftspec.hotfix`: emergency fix outside the standard phase chain — for well-understood fixes touching ≤ 3 files with identified root cause; writes minimal spec, implements fix, verifies inline, and prepares for archive
  - `/draftspec.scope`: quick scope boundary check against the spec's in-scope / out-of-scope sections; produces no file, inline response only; verdict: `in-scope`, `drift`, or `out-of-scope`
  - `/draftspec.recap`: project-level overview of all active features with current phase and inspect verdict; no slug required; produces no file; designed as the first command in a new session
- `--research` flag for `/draftspec.plan`: enters research-first mode — agent identifies concrete unknowns, writes `research.md`, then stops and asks "Research complete — proceed to full plan?" before producing `plan.md`
- Agent targets: `windsurf` (`.windsurf/rules/`), `roocode` (`.roo/rules/`), `aider` (`.aider/CONVENTIONS.md`); total supported targets: 9
- `draftspec demo [path]`: creates a demo workspace at the given path (default: `./draftspec-demo`) with a pre-populated example feature (`export-report`) at the implement phase — spec, inspect report, plan, tasks, and data model are all populated; suggests `/draftspec.scope`, `/draftspec.challenge`, and `/draftspec.handoff` to try immediately
- `draftspec check <slug> [path]`: human-friendly feature readiness check — shows artifact presence, inspect and verify verdict, task progress, and the exact next slash command; exits with code 1 when blocked; supports `--json` for CI use
- `doctor` now warns when `constitution.md` contains unfilled placeholder content (e.g. `[PROJECT_NAME]`) and suggests running `/draftspec.constitution`

### Changed

- All agent prompts now use a unified two-level load structure: `## Load First` (always read) and `## Load If Present` (read when the file exists and is relevant); previously inconsistent across prompts (`## Load Only`, `## Load Only If Needed`, `## Load Only If Present`)
- `constitution.md` prompt: `README.md`, `AGENTS.md`, and project manifests moved from `## Load First` to `## Load If Present`
- `agents-snippet.md` (injected into `AGENTS.md`) now lists four optional commands with one-line descriptions: `challenge`, `handoff`, `scope`, `recap`

## [0.1.0] - 2026-03-31

### Added

- Initial release of the Draftspec CLI
- File-based project context system for development agents and humans
- Eight-phase workflow: `constitution → spec → inspect → plan → tasks → implement → verify → archive`
- Bilingual support for English (`en`) and Russian (`ru`) templates and prompts
- Agent integration for Claude, Codex, Copilot, Cursor, Kilocode, and Trae
- Shell support for `sh` and `powershell`
- `init` and `refresh` commands to manage the `.draftspec/` workspace
- `doctor` command with comprehensive workspace health checks
- `status` and `features` commands for feature lifecycle visibility
- `feature`, `feature repair`, `list-specs`, `show-spec`, `migrate` commands
- `add-agent`, `remove-agent`, `list-agents`, `cleanup-agents` commands for agent artifact management
- Internal CLI (`__internal`) for script delegation without exposing plumbing commands
- Phase readiness scripts for each workflow phase
- Stable IDs for traceability: `AC-*`, `RQ-*`, `DEC-*`, `T*`
- Phase Contract headers in all agent prompts with concrete inputs, outputs, and stop conditions
- Self-Check sections in inspect and verify prompts
- `--version` flag reporting the build version
- Multi-platform CI builds for linux/amd64, linux/arm64, windows/amd64, windows/arm64

[Unreleased]: https://github.com/bzdv/draftspec/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/bzdv/draftspec/releases/tag/v0.1.0
