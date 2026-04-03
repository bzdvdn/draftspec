# Overview

## What Draftspec Is

`draftspec` keeps project intent, specifications, plan artifacts, and tasks in plain files. It is designed to help humans and development agents share the same project context without introducing a rigid process engine.

## Core Ideas

- The constitution is the highest-priority project document.
- Every feature starts as a spec and evolves through a strict workflow.
- Each feature should be developed in its own git branch so teams can collaborate without shared-memory merge churn.
- Generated docs and prompts can use English or Russian.
- Agent workflows should load only the minimum context needed.
- Readiness checks belong in scripts whenever possible.

## Workspace Layout

```text
.draftspec/
  draftspec.yaml
  constitution.md
  specs/
    <slug>/
      spec.md
      inspect.md
      summary.md
      hotfix.md
  plans/
    <slug>/
      plan.md
      tasks.md
      data-model.md
      research.md
      contracts/
        api.md
        events.md
  archive/
    <slug>/
      <YYYY-MM-DD>/
        summary.md
        spec.md
        plan.md
        tasks.md
        data-model.md
        research.md
        contracts/
  templates/
  scripts/
AGENTS.md
```

## Public CLI Surface

The public CLI stays intentionally small:

- `draftspec init [path]`
- `draftspec add-agent [path]`
- `draftspec list-agents [path]`
- `draftspec remove-agent [path]`
- `draftspec cleanup-agents [path]`
- `draftspec doctor [path]`
- `draftspec list-specs [path]`
- `draftspec show-spec <name> [path]`

Creation and evolution of specs, plans, tasks, and implementation are agent-facing workflows, not public CLI subcommands.
