# Brownfield Walkthrough

This walkthrough shows how to introduce `draftspec` into an existing repository without trying to document the whole codebase at once.

## Goal

Use Draftspec as a strict, lightweight coordination layer for one active feature in a real repository.

## When To Use This

Use this flow when:

- the repository already exists
- the team is already shipping code
- you want better agent discipline without adding a heavy process layer
- you want feature-local artifacts instead of shared mutable planning state

## Suggested Demo Setup

Start in an existing repository root:

```bash
draftspec init . --lang en --shell sh --agents codex
draftspec doctor .
```

Expected outcome:

- `.draftspec/` exists
- `AGENTS.md` contains Draftspec guidance
- project-local agent command files exist for the chosen targets
- `doctor` reports a healthy workspace

## Feature-First Adoption

Do not try to spec the whole repository.

Pick one active feature and drive only that scope through the workflow:

1. `/draftspec.constitution`
2. `/draftspec.spec`
3. `/draftspec.inspect`
4. `/draftspec.plan`
5. `/draftspec.tasks`
6. `/draftspec.implement`

Example feature request:

```text
/draftspec.spec Add partner-specific ingestion scheduling with retry policy overrides.
```

Expected artifact growth:

- `.draftspec/specs/<slug>.md`
- `.draftspec/specs/<slug>.inspect.md`
- `.draftspec/plans/<slug>/plan.md`
- `.draftspec/plans/<slug>/tasks.md`
- optional compact plan artifacts only when needed

## What To Highlight In A Demo

For a brownfield demo, the most important points are:

- the repository does not need a full rewrite or full-repo spec effort
- the workflow stays feature-local
- `inspect` is required before planning
- code reading stays narrow and task-driven
- generated scripts help the agent validate readiness before widening context

## Good Before / After Story

Before Draftspec:

- feature intent lives in chat history or scattered notes
- agents reread too much repository context
- plans drift from specs
- branch work is hard to audit

After Draftspec:

- feature intent is persisted in small canonical files
- each phase has a clear entrypoint
- the agent is told what to read first and what not to read by default
- branch-local artifacts keep active work reviewable

## Recommended Capture Format

For public promotion, this is usually best shown as:

- one short terminal GIF for setup
- one markdown walkthrough for the workflow
- one screenshot or excerpt of the generated artifacts

That keeps the demo maintainable while still showing real-world value.
