# Draftspec Archive Prompt

You are archiving one feature package.

## Goal

Create a durable archive snapshot for one feature and update project memory with a short archive entry.

## Load First

Always read these first:

- `.draftspec/memory.md`
- `.draftspec/specs/<slug>.md`

## Load Only If Needed

Read plan artifacts only to inform `summary.md`. If the summary can be written from the spec and memory alone, do not read them:

- `.draftspec/plans/<slug>/plan.md`
- `.draftspec/plans/<slug>/tasks.md`
- `.draftspec/plans/<slug>/data-model.md`
- `.draftspec/plans/<slug>/contracts/`
- `.draftspec/plans/<slug>/research.md`

## Do Not Read By Default

- unrelated specs
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
- If plan artifacts exist, archive them together with the spec.
- If `research.md` does not exist, do not invent it.
- Update `memory.md` by adding a short entry under `Archived Specs` using slug, status, date, and reason.
- Use one of these statuses:
  - `completed`
  - `superseded`
  - `abandoned`
  - `rejected`
  - `deferred`

## Output expectations

- Create the archive snapshot
- Write or patch `summary.md`
- Update `memory.md`
- Summarize archived files, status, and reason
