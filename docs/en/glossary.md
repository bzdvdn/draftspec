# Glossary

## Constitution

The highest-priority project document. It defines non-negotiable rules for architecture, workflow, language policy, and governance.

## Spec

A feature-level document that describes what should be built and why.

## Inspect

A review phase that checks a spec and related artifacts for completeness, consistency, and constitutional compliance.

## Plan

A technical design phase that translates one spec into implementation-oriented artifacts.

## Plan Package

The group of files stored under `.draftspec/plans/<slug>/`, typically including `plan.md`, `tasks.md`, `data-model.md`, `contracts/`, and optional `research.md`.

## Tasks

The executable breakdown of a feature plan. In Draftspec, `tasks.md` lives inside the plan package.

## Implement

The execution phase that works from unfinished tasks and updates task state.

## Verify

A lightweight post-implementation check that confirms whether a feature is aligned enough with tasks, specs, plan artifacts, and project rules to proceed safely toward archive or further completion claims.

## Archive

The phase and storage area used to preserve a historical snapshot of a finished or inactive feature package.

## Data Model

A plan artifact that describes entities, structures, relationships, and important invariants needed for implementation.

## Contracts

Plan artifacts that define interfaces such as APIs, events, or other external interaction boundaries.

## Research

An optional plan artifact used only when uncertainty, external investigation, or architecture tradeoffs need to be preserved.

## Agent Target

A supported agent ecosystem for which Draftspec can generate project-local command or prompt files, such as `claude`, `codex`, `copilot`, `cursor`, `kilocode`, or `trae`.

## Orphaned Agent Artifact

A generated agent file that still exists on disk even though its target is no longer enabled in `.draftspec/draftspec.yaml`.

## Docs Language

The configured language used for generated project documentation such as `constitution.md`, specs, plans, and tasks.

## Agent Language

The configured language used for generated prompts and `AGENTS.md` guidance.

## Comments Language

The configured preferred language for new or edited code comments during implementation.

## Given / When / Then

The canonical BDD markers used for acceptance criteria regardless of documentation language. The surrounding text may still follow the configured docs language.

## Acceptance ID

A stable identifier for an acceptance criterion, such as `AC-001`. It helps keep traceability explicit across specs, tasks, and inspection reports.
