# Draftspec Plan Prompt

You are creating or updating the implementation plan package for one feature.

## Goal

Produce the technical planning artifacts for a spec under `.draftspec/plans/<slug>/`.

## Load Only

- `.draftspec/constitution.md`
- `.draftspec/memory.md`
- `.draftspec/specs/<slug>.md`
- only the repository code and docs needed to plan this one feature

## Do Not Read By Default

- unrelated specs
- unrelated plan packages
- large repository areas with no impact on this feature
- optional `research.md` unless uncertainty already exists

## Stop Conditions

Stop and ask for clarification or refinement if:

- `.draftspec/specs/<slug>.md` does not exist
- the spec is too vague to produce architecture, contracts, or data model decisions
- constitutional constraints conflict with the intended plan

Do not compensate by reading broad unrelated repository context.

## Required outputs

Create or update:

- `.draftspec/plans/<slug>/plan.md`
- `.draftspec/plans/<slug>/data-model.md`
- `.draftspec/plans/<slug>/contracts/api.md`
- `.draftspec/plans/<slug>/contracts/events.md`

Create `.draftspec/plans/<slug>/research.md` only when real uncertainty or external investigation is needed.

## Rules

- The plan MUST comply with the constitution.
- Use the project's configured documentation language for all new or updated planning artifacts.
- Keep the language of `plan.md`, `data-model.md`, `contracts/`, and optional `research.md` internally consistent.
- Respect an established local document convention only when preserving an existing artifact would otherwise become inconsistent.
- Prefer concrete implementation decisions over generic advice.
- Data model and contracts MUST be consistent with the spec.
- Use repository reality, not idealized architecture.
- If critical information is missing, ask only the minimum necessary follow-up questions.

## Output expectations

- Write or patch the plan artifacts
- State whether `research.md` was needed
- Summarize major technical decisions and risks
