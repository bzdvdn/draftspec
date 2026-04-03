# Draftspec Plan Prompt

You are creating or updating the implementation plan package for one feature.

## Goal

Produce the technical planning artifacts for a spec under `.draftspec/plans/<slug>/`.

## Phase Contract

Inputs: `.draftspec/constitution.md`, `.draftspec/specs/<slug>/spec.md`, `.draftspec/specs/<slug>/inspect.md`, narrow repo code.
Outputs: `.draftspec/plans/<slug>/plan.md`, `.draftspec/plans/<slug>/data-model.md`; optional `contracts/`, `research.md`.
Stop if: spec or inspect missing, spec too vague for architecture decisions, or constitutional conflict.

## Operating Mode

- Plan one feature only.
- Prefer patching existing artifacts over broad rewrites.
- Keep context narrow and repository-grounded.
- Produce only the artifacts justified by the feature.

## Flags

`--research`: enter research-first mode before producing the plan.

When `--research` is present in `$ARGUMENTS`:
- Read the spec and inspect report, then identify the 1–5 concrete unknowns that currently block planning.
- Write them to `.draftspec/plans/<slug>/research.md`.
- Stop after writing `research.md` and ask: "Research complete — proceed to full plan?"
- Wait for an explicit confirmation before producing `plan.md`, `data-model.md`, or any other planning artifact.
- If the user confirms, continue with the normal planning flow using the research findings as grounding.
- If the user says to stop, end the session with `research.md` as the only new artifact.

Do not produce `plan.md` in the same pass as `--research` unless the user explicitly confirms.

## Load First

- `.draftspec/constitution.md`
- `.draftspec/specs/<slug>/spec.md`
- `.draftspec/specs/<slug>/inspect.md`
- only the repository code and docs needed to plan this one feature
- when code must be read, prefer the smallest file set needed to identify concrete implementation surfaces, boundaries, and constraints

## Do Not Read By Default

- large repository areas with no impact on this feature
- optional `research.md` unless uncertainty already exists

## Stop Conditions

Stop and ask for clarification or refinement if:

- `.draftspec/specs/<slug>/spec.md` does not exist
- `.draftspec/specs/<slug>/inspect.md` does not exist
- the spec is too vague to produce architecture, contracts, or data model decisions
- constitutional constraints conflict with the intended plan
- the plan would need to cross an unclear integration or architectural boundary that is not justified by the spec or focused repository evidence
- the work would only make sense if multiple feature packages were planned together

Do not compensate by reading broad unrelated repository context.

## Required outputs

Always create or update:

- `.draftspec/plans/<slug>/plan.md`
- `.draftspec/plans/<slug>/data-model.md`

Create only when the feature actually requires it:

- `.draftspec/plans/<slug>/contracts/api.md` — only if the feature touches API boundaries
- `.draftspec/plans/<slug>/contracts/events.md` — only if the feature produces or consumes events

Create `.draftspec/plans/<slug>/research.md` only when at least one of these is true:

- the feature depends on an external system, API, or dependency with behavior that is still unclear
- there are multiple realistic implementation options with meaningful trade-offs that must be preserved
- there is a non-obvious performance, security, reliability, or integration risk that affects planning
- a repository constraint or architectural boundary must be investigated before the plan can be made concrete

Before creating `research.md`, write down the concrete unknowns first:

- list only the 1-5 specific unknowns that block planning
- tie each unknown to a decision, risk, or boundary in this feature
- do not research a technology or subsystem in general; research only the narrow question that changes the plan
- if no concrete unknown remains, do not create `research.md`

## Invariants

- The plan MUST comply with the constitution.
- Inspect is a mandatory prerequisite for planning. Do not plan from a spec that has not been inspected and persisted.
- Keep planning tied to the current spec, not idealized architecture.
- Never read unrelated feature artifacts to compensate for missing clarity.
- Read code narrowly: only enough to ground implementation surfaces, integration boundaries, and repository constraints for this feature.
- When `/.draftspec/scripts/check-plan-ready.*` is available, prefer running it as the phase readiness check instead of reading script source.
- Do not read `/.draftspec/scripts/*` by default unless you are debugging the scripts, working on Draftspec itself, or the user explicitly asks to inspect script logic.
- Prefer concrete implementation decisions over generic advice.
- Record technologies, libraries, framework choices, or version constraints only when they materially affect implementation shape, integration boundaries, validation, or risk.
- If a version or dependency is named, explain why it matters for this feature: compatibility, repository constraint, external contract, rollout risk, or validation.
- Do not enumerate stack details for completeness; capture only technical constraints that reduce downstream guesswork.
- Optional artifacts stay optional; do not create them by habit.
- Do not create `research.md` for generic brainstorming or obvious implementation work that can already be planned from the spec and repository evidence.
- The plan is only complete when downstream task decomposition can proceed without guessing.
- The plan MUST name the concrete implementation surfaces expected to change.
- The plan MUST map each acceptance criterion to an implementation approach before `tasks` are written.
- Do not write the task checklist, edit implementation code, or emit verify/archive conclusions during planning.
- If a downstream task writer would need to guess the method, boundary, or validation path for an `AC-*`, the plan is underspecified.
- The plan should be specific enough that both an agent and a human reviewer can see the intended implementation shape, tradeoffs, and rollout implications without rereading the whole repository.
- Targeted code reading during planning is encouraged when it reduces downstream guesswork; broad repository exploration is not.

## Language Rules

- Use the project's configured documentation language for all new or updated planning artifacts.
- Keep the language of `plan.md`, `data-model.md`, `contracts/`, and optional `research.md` internally consistent.
- Respect an established local document convention only when preserving an existing artifact would otherwise become inconsistent.

## Traceability Rules

- Follow the structure of `.draftspec/templates/plan.md` and `.draftspec/templates/data-model.md` when creating new files.
- Data model and contracts MUST be consistent with the spec and its acceptance criteria.
- Reference stable acceptance IDs from the spec when discussing acceptance-critical behavior.
- When the plan makes a significant implementation choice, record it as a stable decision ID such as `DEC-001`.
- Each significant `DEC-*` SHOULD state `Why`, `Affects`, and `Validation`.
- If the feature introduces or changes persisted state, state transitions, or lifecycle rules, capture them in `data-model.md`.
- If the feature crosses an API boundary, capture request, response, and error behavior in `contracts/api.md`.
- If the feature produces or consumes events, capture producer, consumer, payload, and delivery assumptions in `contracts/events.md`.
- Do not leave entity shape, boundary IO, or event payload details only in prose inside `plan.md`.
- Each entity or contract entry SHOULD reference the `AC-*` that justifies it.
- If `data-model.md` is created, explicitly state which entities, invariants, or lifecycle concerns require it.
- If `contracts/` is created, explicitly state which API or event boundary requires it.
- If neither richer artifact is needed, prefer not creating it.
- Use repository reality, not idealized architecture.
- If critical information is missing, ask only the minimum necessary follow-up questions.
- Treat generic wording such as `update backend accordingly`, `adjust logic as needed`, or `wire this through the system` as a refinement signal rather than a complete plan.

## Content Quality Rules

- `## Goal` should describe implementation shape, not restate user-facing spec text.
- `## Implementation Surfaces` should explain why each surface changes and whether it is new or existing.
- `## Acceptance Approach` should map every `AC-*` to touched surfaces and proof of observability.
- `## Data and Contracts` should say what changes, what stays unchanged, and why.
- Each significant `DEC-*` should capture `Why`, `Tradeoff`, `Affects`, and `Validation`.
- `## Sequencing Notes` should separate must-happen-first work from work that can be parallelized.
- `## Risks` should include a mitigation, not just the risk label.
- `## Rollout and Compatibility` should be explicit when migration, flags, compatibility, or operational follow-up matter, and should say so plainly when they do not.
- `## Validation` should tie checks back to `AC-*` or `DEC-*`, not just list generic test ideas.
- Add a short `Unknowns First` pass before finalizing the plan: if a decision, surface, or validation path is still unclear, record the unknown explicitly or stop for refinement.
- Negative examples: avoid `update backend accordingly`, `wire this through the system`, or `adjust logic as needed` as if they were implementation decisions.
- Prefer concrete implementation guidance over architecture essay prose. If a paragraph does not reduce downstream guesswork, tighten it.

## Output expectations

- Write or patch the plan artifacts; state which optional artifacts were created and why
- Summarize key technical decisions and call out risks blocking downstream phases
- End with a summary block: `Slug`, `Status`, `Artifacts`, `Blockers`, `Next command`
- When ready: `Next command: /draftspec.tasks <slug>`

## Self-Check

- Did I keep optional artifacts justified?
- Can `tasks` be written from this package without guesswork?
- Would a human reviewer understand the main implementation shape, risks, and rollout story quickly?
