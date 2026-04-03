# Draftspec Archive Prompt

You are archiving one feature package.

## Goal

Create a durable archive snapshot for one feature.

## Phase Contract

Inputs: see Load First and Load Only If Needed.
Outputs: archive snapshot and summary.
Stop if: see Stop Conditions.

## Load First

Always read these first:

- `.draftspec/specs/<slug>.md`

## Load If Present

Read plan artifacts only to inform `summary.md`. If the summary can be written from the spec alone, do not read them:

- `.draftspec/plans/<slug>/plan.md`
- `.draftspec/plans/<slug>/tasks.md`
- `.draftspec/plans/<slug>/data-model.md`
- `.draftspec/plans/<slug>/contracts/`
- `.draftspec/plans/<slug>/research.md`

## Do Not Read By Default

- unrelated archive entries
- unrelated repository files

## Stop Conditions

Stop and ask a minimal follow-up question if:

- the archive status is unclear
- the reason for archiving is missing
- the target slug is ambiguous

## Rules

- Archive under `.draftspec/archive/<slug>/<YYYY-MM-DD>/`.
- Use copy-based archiving in MVP; do not delete active files (`specs/<slug>.md` and `plans/<slug>/`) unless explicitly instructed by the user.
- Write `summary.md` inside the archive directory.
- Keep `summary.md` compact. Prefer status, reason, completed scope, and notable deviations over long retrospective narration.
- If plan artifacts exist, archive them together with the spec.
- If `research.md` does not exist, do not invent it.
- If status is `completed` and `tasks.md` exists, use `/.draftspec/scripts/verify-task-state.* <slug>` before archiving. Do not claim a completed archive when required tasks are still open.
- Use one of these statuses:
  - `completed`
  - `superseded`
  - `abandoned`
  - `rejected`
  - `deferred`

## Output expectations

- Create the archive snapshot
- Write or patch `summary.md`
- Summarize archived files, status, and reason
- State explicitly that archive is the terminal workflow step for this feature unless the user asks for follow-up work
