# Roadmap

This roadmap focuses on the next practical iterations for Draftspec rather than a long speculative backlog.

## Direction

Draftspec should continue to position itself between heavier spec-driven systems and looser change-driven systems:

- stricter than OpenSpec
- lighter than spec-kit
- optimized for agent-first workflows on real codebases

## Iteration 1

### Primary goal

Strengthen `inspect` as the central quality layer.

### Planned work

- define and refine an explicit inspection report structure and verdict semantics
- improve checks for `constitution <-> spec`
- improve checks for `spec <-> plan`
- improve checks for `plan <-> tasks`
- strengthen acceptance-to-task traceability checks

### Anti-Bloat Notes

Safe direction:

- stronger structural checks
- clearer verdict semantics
- better traceability through stable acceptance IDs

Use caution with:

- making persisted inspect reports mandatory for every feature
- reading implementation code by default during inspect
- turning inspect into a broad review engine

### Why this matters

If `inspect` is strong, every downstream phase gets better with less wasted implementation effort.

## Iteration 2

### Primary goal

Add a lightweight post-implementation verification layer.

Status: lightweight contract, prompt, readiness script, and report template are now in place. The remaining work is to deepen checks without expanding default context.

### Planned work

- introduce a small `verify` or review-oriented workflow after `implement`
- check whether completed tasks match implementation state
- check whether implementation still matches spec and plan intent
- ensure memory and task state remain synchronized

### Anti-Bloat Notes

Safe direction:

- task-state verification helpers
- memory/task synchronization checks
- optional persisted verify reports

Use caution with:

- reading code by default during verify
- turning verify into a heavy review or QA engine
- requiring verify artifacts before every downstream action

### Why this matters

This closes the gap between "tasks were executed" and "the feature is actually aligned with its intended design".

## Iteration 3

### Primary goal

Strengthen brownfield ergonomics and automation outputs.

### Planned work

- improve archive summaries and archive linkage in `memory.md`
- keep completed-archive checks cheap by reusing task-state verification
- add machine-readable outputs such as `doctor --json`
- improve config-aware helper behavior for scripts and future tooling
- continue polishing multilingual consistency in docs and prompts

### Anti-Bloat Notes

Safe direction:

- machine-readable outputs for existing checks
- better archive indexing and summaries
- config-aware helpers that reduce repeated reasoning

Use caution with:

- archive flows that require reading broad repository history
- new automation outputs that introduce mandatory new artifacts
- brownfield helpers that silently widen the default context

### Why this matters

This makes Draftspec easier to automate, easier to operate at scale, and stronger for long-lived existing codebases.

## Always-On Quality Work

Alongside feature work, Draftspec should keep improving:

- documentation consistency
- unit test coverage
- CLI ergonomics
- prompt clarity and token discipline
- brownfield workflow quality

## Not Planned Right Now

Draftspec should avoid these unless there is a very strong reason:

- a heavy orchestration engine
- mandatory checkpoint systems
- approval-gate bureaucracy
- large default prompt contexts
- required artifact sprawl for every feature
