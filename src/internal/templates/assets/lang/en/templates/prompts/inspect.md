# Draftspec Inspect Prompt

You are inspecting one feature package for consistency and quality.

## Goal

Produce a focused inspection report for one feature without expanding scope.

## Load First

Always read these first:

- `.draftspec/constitution.md`
- `.draftspec/memory.md`
- `.draftspec/specs/<slug>.md`

## Load Only If Needed

Read these only when they exist and materially affect the inspection:

- `.draftspec/plans/<slug>/plan.md`
- `.draftspec/plans/<slug>/tasks.md`
- `.draftspec/plans/<slug>/data-model.md`
- `.draftspec/plans/<slug>/contracts/`
- `.draftspec/plans/<slug>/research.md`

## Do Not Read By Default

- unrelated specs
- unrelated plan packages
- broad repository history
- implementation files unless they are needed to verify a concrete consistency claim

## Stop Conditions

Stop and ask a minimal follow-up question only if:

- the target slug is ambiguous
- the spec is missing entirely
- the inspection would otherwise invent missing product intent

## Rules

- Check constitutional consistency first.
- Inspect spec completeness and clarity.
- If plan artifacts exist, check alignment between spec, plan, data model, contracts, and tasks.
- Keep the inspection report in the project's configured documentation language when writing it to disk.
- Prefer concrete findings over generic advice.
- Use these sections in the report:
  - `Errors`
  - `Warnings`
  - `Questions`
  - `Suggestions`
- Suggest Given/When/Then scenarios only when they would materially strengthen weak acceptance criteria.

## Output expectations

- Write or patch an inspection report for the feature
- Summarize errors, warnings, open questions, and suggestions
