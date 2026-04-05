## Draftspec

Primary project context lives in `.draftspec/`. Languages: docs=[DOCS_LANGUAGE], agent=[AGENT_LANGUAGE], comments=[COMMENTS_LANGUAGE]

Workflow chain: `constitution → spec → inspect → plan → tasks → implement → verify → archive`
- `/draftspec.constitution`: create or patch `.draftspec/constitution.md`
- `/draftspec.spec`: create or refine `.draftspec/specs/<slug>/spec.md`; `--amend` for targeted edits
- `/draftspec.inspect`: check one feature for consistency and quality
- `/draftspec.plan`: create or patch `.draftspec/plans/<slug>/`; `--update` for targeted edits, `--research` for research-first
- `/draftspec.tasks`: create or patch `.draftspec/plans/<slug>/tasks.md`
- `/draftspec.implement`: execute unfinished tasks from `tasks.md`
- `/draftspec.verify`: verify one feature package; `--deep` for full per-AC code tracing
- `/draftspec.archive`: archive to `.draftspec/archive/` (move-based); `--copy` keeps originals, `--restore` unarchives

Optional (any point): `/draftspec.challenge` (adversarial review; `--spec`/`--plan`), `/draftspec.handoff` (session handoff), `/draftspec.hotfix` (emergency fix ≤ 3 files), `/draftspec.scope` (boundary check; `--plan`/`--tasks`), `/draftspec.recap` (project overview)

Read discipline:
- Do not skip phases; load only the current feature slug by default
- Prefer readiness scripts over reading deeper artifacts; use `./.draftspec/scripts/run-draftspec.sh` for CLI access
- Never load: unrelated specs/plans, broad repo scans, script source, files already read this session (unless you edited them)
- Use the configured comment language for new/edited code comments; preserve existing file conventions

Before meaningful changes: review `constitution.md`, the relevant `specs/<slug>/spec.md`, and `plans/<slug>/` if present. After changes: keep specs, plans, tasks, and implementation aligned.
