# Changelog

All notable changes to this project will be documented in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Versions follow [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- `plan.md` template: new `## Incremental Delivery` section with `MVP (First Value)` and `Iterative Expansion` subsections — guides agents to define the smallest testable increment and plan subsequent value delivery steps while maintaining AC traceability
- `/draftspec.inspect --delta`: incremental re-check mode — after `spec --amend`, re-checks only changed sections and preserves valid findings from the previous report; falls back to full inspection when changes exceed 50%
- `/draftspec.tasks --repair <task-id-list>`: targeted repair mode — fixes specific tasks identified by verify or review without rewriting the full task list; suggests `plan --update` if the repair reveals a plan-level flaw
- `/draftspec.implement --continue`: resume mode — starts from the first unfinished task, trusts previously checked-off tasks, batch-reads only remaining surfaces; useful after session interruptions
- `/draftspec.verify --persist`: explicit flag to write the verification report to `.draftspec/plans/<slug>/verify.md`; without it the report stays in the conversation only
- Constitution post-update impact check: after updating `constitution.md`, the agent scans active specs for conflicts and flags them as `NEEDS RE-INSPECT` without modifying the specs
- Handoff now loads `verify.md` metadata — verdict and verification status are included in handoff documents for better session continuity

### Changed

- Agent adapters now include **workflow chain hint** (`constitution → spec → inspect → plan → tasks → implement → verify → archive`), **explicit script execution discipline** (execute as shell commands, trust stdout/exit code, never read source), and an **anti-pattern block** (do not skip readiness scripts, re-plan during implement, mark done without proof, or read full repo when minimal context is required)
- `inspect` helper flow now treats readiness/script output as the primary structural evidence layer: categorized findings (`structure`, `traceability`, `ambiguity`, `consistency`, `readiness`) are surfaced in `draftspec check`, `draftspec feature`, and inspect prompts so agents deepen findings instead of re-deriving them
- phase readiness checks now emit structured findings for ambiguity and acceptance/task traceability; implement readiness also warns when plan implementation surfaces drift from `tasks.md` `Surface Map` or `Touches:` references
- `/draftspec.archive` now uses **move-based** archiving by default — active files (`specs/<slug>/` and `plans/<slug>/`) are deleted after copying to `.draftspec/archive/`; pass `--copy` to keep originals in place (useful for `deferred` features)
- `/draftspec.recap` now shows **recently archived** features (last 7 days) with status, date, and reason — gives new sessions context about what was recently completed or deferred
- All `## Load If Present` sections now use **concrete trigger conditions** instead of vague "when needed" / "when relevant" — each artifact lists the specific signal that justifies reading it (e.g., "when a task references a `DEC-*`", "when verifying `AC-*` coverage")
- Compressed verbose sections: `constitution.md` Repository Evidence (14→4 lines), `implement.md` Read Discipline (16→6 lines), `plan.md` and `tasks.md` Content Quality Rules (deduplicated negative examples)
- `/draftspec.handoff` Load If Present now uses phase-aware triggers — each artifact specifies which handoff section it populates and at which phase it becomes relevant
- `/draftspec.spec` Content Quality Rules now include **positive examples** alongside negative ones for `## Requirements` and `## Edge Cases` — shows the concrete quality bar instead of only warning about bad patterns
- `/draftspec.tasks` `Touches:` field is now **MUST** (was SHOULD) — a task without `Touches:` forces the implement agent into exploratory reads; module-level names allowed only when the exact file cannot be determined yet
- `/draftspec.tasks` `## Surface Map` section is now **MUST** before the first phase — serves as the implement agent's batch-read manifest; without it the agent must scan every task line to build the read list
- `/draftspec.plan` now requires **`## Constitution Compliance`** section — explicit `no conflicts` or list of specific conflicts with constitution sections and how the plan resolves them; makes adherence reviewable instead of implicit
- `agents-snippet.md` compressed **~55%** (~1000 → ~450 tokens) — removed per-phase read hints, detailed optional command descriptions, and separate discipline sections that duplicate per-prompt rules; snippet now serves as a concise project map, not a second copy of prompt instructions

### Added

- Spec template: mandatory `## Assumptions` section — records environmental assumptions, reasonable defaults, and system dependencies explicitly so inspect can catch wrong ones early
- Spec template: optional `## Success Criteria` with `SC-*` IDs — measurable performance/reliability/UX targets separate from behavioral `AC-*` criteria; omit for purely behavioral features
- Spec template and prompt: `[NEEDS CLARIFICATION: ...]` inline markers — flag unclear requirements or AC details directly where they appear instead of only in `## Open Questions`; inspect treats remaining markers as `Error` blocking planning
- `inspect` prompt: checks for `[NEEDS CLARIFICATION]` markers (Error), missing `## Assumptions` (Warning), assumption plausibility against constitution/repo (Error if contradicted), and vague `SC-*` entries (Warning)
- `research.md` template: mandatory `## Goal`, `## Research Questions`, `## Exploration Areas`, and `## Recommendations` sections — provides a structured way to record discovery and architecture trade-offs
- `draftspec trace <slug> [path]`: scan for `@ds-task` and `@ds-test` annotations in code — provides verifiable traceability between requirements, tasks, and implementation (including test evidence)
- `.draftspec/scripts/trace.*`: new helper script for agents to run traceability checks without direct CLI dependency
- `draftspec trace --tests`: new flag to filter only test-related annotations (`@ds-test`)
- `Verify Evidence`: the `verify` prompt now instructs agents to use `draftspec trace` as primary implementation evidence; includes a `Legacy Fallback` strategy for features without annotations (manual inspection of `Touches:` files)
- `draftspec dashboard [path]`: visual dashboard of all active features, their progress, status, and Git branch alignment
- `Lazy Decomposition`: `tasks` prompt now discourages micro-tasks (1-5 lines) during initial decomposition to save tokens; `implement` prompt now allows `In-place Decomposition` (indented sub-tasks like `T1.1.1`) for the active task only
- `Decomposition Guardrails`: implementation agents are forbidden from adding new implementation surfaces (`Touches:`) or changing `AC-*` mapping during in-place decomposition
- `Smart Branching`: `doctor`, `check`, and `dashboard` now warn when the current Git branch does not match the feature slug (expected `feature/<slug>`)
- `doctor` now includes traceability and branching checks: warns about orphaned `@ds-task` annotations and invalid `AC-*` references, plus non-standard branch names
- `/draftspec.archive --restore`: reverse a previous archive — restores the latest snapshot back into active `specs/` and `plans/`, removes the archive entry, and suggests `/draftspec.inspect` as the next step; stops if active files already exist to prevent overwriting
- `/draftspec.verify --deep`: full implementation validation mode — reads all plan artifacts and traces every completed task and `AC-*` through actual code; default mode remains structural and cheap; deep mode produces comprehensive per-AC evidence report
- `/draftspec.plan --update`: targeted edit mode — update a specific section, `DEC-*`, implementation surface, or add a contract without rewriting the entire plan package; does not invalidate downstream `tasks.md` unless the change materially affects task decomposition

## [v0.2.0] 2026-04-02

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
- `agents-snippet.md` (injected into `AGENTS.md`) now lists five optional commands with one-line descriptions: `challenge`, `handoff`, `hotfix`, `scope`, `recap`

## [v0.1.0] - 2026-03-31

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
