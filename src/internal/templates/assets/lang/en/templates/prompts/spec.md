# Draftspec Spec Prompt

You are creating or updating one feature spec.

## Goal

Produce a clear feature specification at `.draftspec/specs/<slug>.md` that is compliant with the constitution.

Before writing or updating the spec, ensure work is happening on the feature branch for `<slug>`. The default branch naming convention is `feature/<slug>`.

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
- the request asks to derive one spec from multiple constitutional changes without naming a single concrete feature or change
- the request combines multiple features or unrelated changes into one spec
- the input is a prompt file with a generic filename and no explicit `name:` or `slug:` metadata, and the user did not provide `--name` or `--slug`
- the input looks like a URL rather than a concrete feature title
- acceptance criteria would be invented rather than derived
- the requested feature appears to conflict with the constitution

If the user provided `--name` but has not yet given enough feature detail, do not lose the request context: ask for the missing description or accept the next user message as the continuation of the same spec request.

Do not continue into planning or implementation thinking when the spec itself is still unclear.

If the spec already exists and is current, say so and do not modify the file.

## Rules

- The spec MUST comply with the constitution.
- Keep the spec focused on one feature or change.
- `/draftspec.spec` may receive `--name <feature name>`, optional `--slug <feature-slug>`, and optional `--branch <branch-name>`.
- If `--name` is present, use it as the canonical feature name.
- If `--slug` is present, use it for the spec path.
- When `/draftspec.spec` starts from a prompt file, prefer explicit metadata at the top of the file:
  - `name: <feature name>`
  - optional `slug: <feature-slug>`
- Command-style `--name` and `--slug` arguments take priority over `name:` and `slug:` in a prompt file.
- If `slug:` is present, use it for the spec path and feature branch.
- If `--name` is present without `--slug`, derive `<slug>` from `--name`.
- If only `name:` is present, derive `<slug>` from it.
- Fall back to the file basename only when it is specific enough to produce a safe slug.
- If the user explicitly provides `--branch <name>`, use that branch name as-is instead of the default `feature/<slug>`.
- An explicit `--branch` override does not change the spec slug unless the user also requests a different slug.
- If the user provided only `--name` and the detailed description is still missing, ask for the description or accept the next user message as the continuation of the same spec request.
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

- Create or switch to `feature/<slug>` before editing the spec when branch creation is available in the current environment, unless the user explicitly provides `--branch`.
- Write or patch `.draftspec/specs/<slug>.md`, where `<slug>` is the lowercase kebab-case of the feature name
- Summarize goal, scope, acceptance criteria, and open questions
