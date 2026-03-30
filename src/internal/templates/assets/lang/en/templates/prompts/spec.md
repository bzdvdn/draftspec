# Draftspec Spec Prompt

You are creating or updating one feature spec.

## Goal

Produce a clear feature specification at `.draftspec/specs/<slug>.md` that is compliant with the constitution.

Before writing or updating the spec, ensure work is happening on the feature branch for `<slug>`. The default branch naming convention is `feature/<slug>`.

## Phase Contract

Inputs: see Load Only.
Outputs: see Output expectations.
Stop if: see Stop Conditions.

## Operating Mode

- Work on exactly one feature.
- Prefer patching an existing spec over rewriting it.
- Load only the minimum context needed to remove ambiguity.
- Do not drift into planning or implementation.

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
- the request would require multiple feature slugs or multiple independent specs to stay honest about scope
- the input is a prompt file with a generic filename and no explicit `name:` or `slug:` metadata, and the user did not provide `--name` or `--slug`
- the input looks like a URL rather than a concrete feature title
- acceptance criteria would be invented rather than derived
- the requested feature appears to conflict with the constitution
- any required section would be left as TBD or placeholder text

If the user provided `--name` but has not yet given enough feature detail, do not lose the request context: ask for the missing description and treat the next non-command user message as the continuation of the same spec request.

If clarification is needed, prefer a tiny structured clarify pass instead of a broad open-ended interview:

- ask at most 1-3 questions
- ask only about gaps that would otherwise force invented acceptance criteria, unclear scope boundaries, or ambiguous success conditions
- prefer coverage-based questions such as missing scenario, constraint, actor, or edge-condition clarification
- once the answers are sufficient, patch the spec immediately instead of starting a separate clarification workflow

Do not continue into planning or implementation thinking when the spec itself is still unclear.

If the spec already exists and is current, say so and do not modify the file.

## Invariants

- The spec MUST comply with the constitution.
- Keep the spec focused on one feature or change.
- Never load unrelated feature artifacts to compensate for unclear requirements.
- Derive acceptance criteria from the request and repository reality; do not invent them.
- Use explicit scope boundaries. The out-of-scope section is mandatory.
- The spec should be detailed enough that both an agent and a human reviewer can understand the user flow and scope boundaries without reading planning artifacts.
- Do not mix languages inside the same spec without a strong project reason.
- Follow `.draftspec/templates/spec.md` when creating a new file.
- Every acceptance criterion MUST use Given/When/Then format.
- Every acceptance criterion MUST have a stable ID such as `AC-001`.
- Ask follow-up questions only when the missing information is critical.

## Resolution Rules

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
- If the user provided only `--name` and the detailed description is still missing, ask for the description and keep staged mode active for the next non-command user message.
- If the next user message begins with `/draftspec.`, staged mode is canceled and the new slash-command takes priority.
- If the next user message does not begin with `/draftspec.`, treat it as the continuation of the staged spec request.

## Language Rules

- Use the project's configured documentation language for new or updated spec content.
- Respect an established local document convention only when preserving an existing file would otherwise become inconsistent.
- Do not introduce mixed-language headings or sections in the same spec without a strong project reason.

## Acceptance Rules

- The `Given`, `When`, and `Then` markers remain canonical regardless of the documentation language:
  - **Given** — the initial state or precondition
  - **When** — the action or event
  - **Then** — the expected observable outcome
- Acceptance criteria should be observable and testable.
- Prefer a small set of strong criteria over a long redundant list.
- Each `AC-*` should explain why it matters in one short line when that context helps downstream planning stay grounded.
- Each `AC-*` should include evidence or an observable proof signal, not just a generic desired state.

## Content Quality Rules

- `## Goal` should explain who benefits, what changes, and what success looks like.
- `## Why Now` should capture why this change matters now when timing, pain, or business pressure materially affects prioritization.
- `## Primary User Flow` should describe the main path in 3-5 concrete steps, not generic prose.
- `## Change Delta` should make it obvious what becomes newly possible, what changes, and what stays unchanged.
- `## Affected Surfaces` should stay compact and name only the user-visible or repository-visible surfaces that define the feature boundary.
- `## Scope Snapshot`, `## Scope`, and `## Non-Goals` should make the feature boundary obvious to a reviewer.
- `## Context` should capture repository constraints, preserved behavior, and assumptions that materially affect the feature.
- `## Requirements` should stay clear and testable; avoid vague wording like `support this properly` or `handle it cleanly`.
- `## Edge Cases` should include only behavior that materially changes implementation or validation, not a brainstorming dump.
- `## Open Questions` should say `none` when no real question remains.
- Negative examples: do not merge multiple features into one spec, do not hide scope expansion inside edge cases, and do not use `TBD` acceptance criteria.
- Prefer density over length: every section should help planning or review, and filler text is a defect.

## Output expectations

- Create or switch to `feature/<slug>` before editing the spec when branch creation is available in the current environment, unless the user explicitly provides `--branch`.
- Write or patch `.draftspec/specs/<slug>.md`, where `<slug>` is the lowercase kebab-case of the feature name
- Summarize goal, scope, acceptance criteria, and open questions

## Self-Check

- Did I stay within one feature?
- Did I avoid planning and implementation detail?
- Did every acceptance criterion get a stable ID and Given/When/Then form?
- Would a human reviewer understand the primary user flow and scope boundary without additional explanation?
- Did I avoid loading unrelated artifacts?
