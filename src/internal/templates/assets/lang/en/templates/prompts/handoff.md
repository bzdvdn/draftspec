# Draftspec Handoff Prompt

You are generating a session handoff document for one feature.

## Goal

Produce a compact handoff document at `.draftspec/handoff/<slug>.md` that allows a new agent session to resume work on the feature without loss of context.

## Phase Contract

Inputs: current feature artifacts for `<slug>` (phase-appropriate subset).
Outputs: `.draftspec/handoff/<slug>.md` with current phase, completed work, open items, key decisions, and next command.
Stop if: slug is ambiguous or no feature artifacts exist for the slug.

## Load First

Always read these first if they exist:

- `.draftspec/constitution.md`
- `.draftspec/specs/<slug>/spec.md`

## Load If Present

Read based on the inferred phase — only the artifacts that exist and contribute to the handoff:

- `.draftspec/specs/<slug>/inspect.md` — read the verdict line to populate `## Current Phase` and detect blockers
- `.draftspec/plans/<slug>/plan.md` — read `DEC-*` entries to populate `## Key Decisions` when the feature is in plan phase or later
- `.draftspec/plans/<slug>/tasks.md` — read task checkboxes to populate `## Open Work` and `## Completed` when the feature is in tasks/implement phase or later
- `.draftspec/plans/<slug>/verify.md` — read the metadata block and verdict line to populate `## Current Phase` and detect verification concerns or blockers

## Do Not Read By Default

- `.draftspec/plans/<slug>/data-model.md`
- `.draftspec/plans/<slug>/contracts/`
- `.draftspec/plans/<slug>/research.md`
- unrelated specs or plan packages
- implementation files
- script source files

## Stop Conditions

Stop and ask only if:

- the slug is ambiguous and cannot be derived from context or arguments
- no feature artifacts exist for the slug
- called without a slug and no active features exist at all

## All-Features Mode

When called without a slug (no slug in the user arguments):
- Run `.draftspec/scripts/list-specs.*` to enumerate active features; do not read its source.
- For each active feature, generate its handoff document at `.draftspec/handoff/<slug>.md` using the same rules as single-feature mode.
- Output a brief inline summary table: one row per feature with slug, phase, and next command.
- Always overwrite existing handoff files — each is a current-state snapshot.

## Rules

- Infer the current phase from the set of artifacts present: no spec → pre-spec; spec exists, no inspect → spec; inspect exists, no plan → inspect; plan exists, no tasks → plan; tasks exist with open items → tasks or implement; tasks all closed → verify or archive.
- When available, run `.draftspec/scripts/list-open-tasks.*` to enumerate incomplete tasks; rely on its output rather than reading the script source.
- Keep the handoff document compact: it must be loadable in a single cheap read at the start of a new session.
- Do not reproduce full artifact content in the handoff — reference file paths, not contents.
- Reference open items by their stable IDs (T*, AC-*, DEC-*, RQ-*) when available.
- Every section must add signal; omit sections that would be empty or redundant.
- Use the project's configured documentation language when writing to disk.
- Include a machine-readable metadata block at the top.

## Output Structure

Write `.draftspec/handoff/<slug>.md` using this structure:

- YAML metadata block: `report_type`, `slug`, `phase`, `docs_language`, `generated_at`
- `# Handoff: <slug>`
- `## Current Phase` — which phase the feature is in and the evidence for that inference
- `## Completed` — artifacts present and closed tasks (by ID when available)
- `## Open Work` — remaining tasks or missing required artifacts, with stable IDs where available
- `## Key Decisions` — decisions in `plan.md` that materially affect remaining work (DEC-* IDs)
- `## Open Questions` — blockers or unresolved questions that must be addressed before the next phase
- `## Artifacts` — exact project-relative paths to all relevant files for this slug
- `## Next Command` — the exact slash command and slug to resume work immediately

## Output Expectations

- Always overwrite `.draftspec/handoff/<slug>.md` if it already exists — the handoff is a current-state snapshot and the previous version is immediately stale.
- Write the file to `.draftspec/handoff/<slug>.md`.
- Also output a brief inline summary: current phase, number of open items, and next command.
- End the conversation with the exact `Next command` line so a new session can start without re-reading.
- If open questions block the next phase, state that in `## Next Command` instead of naming the phase command.

## Self-Check

- Is the current phase correctly inferred from which artifacts are present?
- Are open items referenced by stable IDs wherever possible?
- Is the document short enough to load cheaply at the start of a new session?
- Does `## Next Command` contain the exact slash command and slug to resume?
- Did I avoid reproducing full artifact content in the handoff?
