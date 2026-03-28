# Draftspec Spec Prompt

You are creating or updating one feature spec.

## Goal

Produce a clear feature specification at `.draftspec/specs/<slug>.md` that is compliant with the constitution.

## Load Only

- `.draftspec/constitution.md`
- the current user request and conversation
- the smallest amount of repository context needed to remove ambiguity

## Do Not Read By Default

- unrelated specs
- unrelated plan packages
- implementation-heavy code areas unless they are needed to define scope correctly
- contracts or data models for other features

## Stop Conditions

Stop and ask a minimal follow-up question if:

- the feature goal is ambiguous
- acceptance criteria would be invented rather than derived
- the requested feature appears to conflict with the constitution

Do not continue into planning or implementation thinking when the spec itself is still unclear.

If the spec already exists and is current, say so and do not modify the file.

## Rules

- The spec MUST comply with the constitution.
- Keep the spec focused on one feature or change.
- Use the project's configured documentation language for new or updated spec content.
- Respect an established local document convention only when preserving an existing file would otherwise become inconsistent.
- Do not introduce mixed-language headings or sections in the same spec without a strong project reason.
- Follow the structure of `.draftspec/templates/spec.md` when creating a new file.
- Every acceptance criterion MUST use Given/When/Then format. The `Given`, `When`, and `Then` markers remain canonical regardless of the documentation language:
  - **Given** — the initial state or precondition
  - **When** — the action or event
  - **Then** — the expected observable outcome
- Use explicit scope boundaries. The out-of-scope section is mandatory.
- Ask follow-up questions only when the missing information is critical.
- Do not jump into implementation planning here.

## Output expectations

- Write or patch `.draftspec/specs/<slug>.md`, where `<slug>` is the lowercase kebab-case of the feature name
- Summarize goal, scope, acceptance criteria, and open questions
