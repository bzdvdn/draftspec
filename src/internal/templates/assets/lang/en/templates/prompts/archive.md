# Draftspec Archive Prompt

You are archiving one feature package.

## Goal

Create a durable archive snapshot for one feature, or restore a previously archived feature back to active development.

## Flags

`--copy`: keep originals in place after archiving (copy-only mode). Useful for `deferred` features that may return to active development.

`--restore`: reverse a previous archive — copy the latest snapshot back into active `specs/` and `plans/`, then remove the archive entry. See Restore Rules below.

## Phase Contract

Inputs: see Load First and Load Only If Needed.
Outputs: archive snapshot and summary (default mode); restored active files (restore mode).
Stop if: see Stop Conditions.

## Load First

**Default mode** — always read these first:

- `.draftspec/specs/<slug>/spec.md`

**Restore mode** (`--restore`) — read the archive snapshot instead:

- `.draftspec/archive/<slug>/` — use the most recent dated directory

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

- the archive status is unclear (default mode)
- the reason for archiving is missing (default mode)
- the target slug is ambiguous
- `--restore` is used but no archive exists for the slug
- `--restore` is used but active files already exist for the slug (would overwrite)

## Rules

- Archive under `.draftspec/archive/<slug>/<YYYY-MM-DD>/`.
- Default behavior is **move-based**: copy artifacts to the archive directory, then **delete** active files (`specs/<slug>/spec.md` and `plans/<slug>/`). This keeps the active workspace clean.
- If the user passes `--copy`, keep the originals in place (copy-only mode).
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

## Restore Rules

When `--restore` is present in `$ARGUMENTS`:

- Locate the most recent snapshot under `.draftspec/archive/<slug>/` (by date directory name).
- If active files already exist for the slug (`specs/<slug>/spec.md` or `plans/<slug>/`), stop and ask the user — restoring would overwrite active work.
- Copy `spec.md` back to `.draftspec/specs/<slug>/spec.md`.
- Copy plan artifacts (`plan.md`, `tasks.md`, `data-model.md`, `contracts/`, `research.md`) back to `.draftspec/plans/<slug>/` if they exist in the snapshot.
- Do not copy `summary.md` or `inspect.md` back — the restored feature should be re-inspected.
- After a successful restore, delete the archive entry (the dated directory). If it was the only snapshot, remove the slug directory from `archive/` as well.
- Report which files were restored and suggest the next workflow command based on the restored state (typically `/draftspec.inspect <slug>`).

## Output expectations

### Default mode

- Create the archive snapshot
- Write or patch `summary.md`
- Summarize archived files, status, and reason
- State explicitly that archive is the terminal workflow step for this feature unless the user asks for follow-up work

### Restore mode

- List all restored files and their destinations
- State which archive snapshot was used and confirm it was removed
- Suggest the next workflow command (typically `/draftspec.inspect <slug>`)
