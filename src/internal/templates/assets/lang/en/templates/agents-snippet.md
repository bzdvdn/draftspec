## Draftspec

Primary project context lives in `.draftspec/`.

Preferred language settings:
- Documentation: [DOCS_LANGUAGE]
- Agent interaction: [AGENT_LANGUAGE]
- Code comments: [COMMENTS_LANGUAGE]

Workflow commands:
- `/draftspec.constitution`: patch `.draftspec/constitution.md`
- `/draftspec.spec`: create or refine one file in `.draftspec/specs/<slug>.md` and work from `feature/<slug>`
- `/draftspec.inspect`: inspect one feature for consistency and quality before or after planning
- `/draftspec.plan`: create or patch `.draftspec/plans/<slug>/plan.md`, `data-model.md`, and `contracts/`
- `/draftspec.tasks`: create or patch `.draftspec/plans/<slug>/tasks.md`
- `/draftspec.implement`: execute unfinished tasks
- `/draftspec.archive`: archive one feature package under `.draftspec/archive/`

Read discipline:
- Follow `constitution -> spec -> inspect -> plan -> tasks -> implement -> archive`
- Do not skip prerequisites
- Load only the current feature slug by default
- Prefer the readiness scripts for each phase before reading deeper artifacts
- `/draftspec.spec` supports `--name`, optional `--slug`, and optional `--branch`; for chat-based input, the feature description may arrive in the next message
- For file-based `/draftspec.spec` input, prefer a top-of-file `name:` and optional `slug:` before falling back to the filename
- Allow an explicit `--branch <name>` override for repository-specific branch naming conventions such as Jira keys
- During `tasks`, start with `plan.md` and read deeper artifacts only if required
- During `implement`, start with `tasks.md` and read deeper artifacts only if required

Never load by default:
- unrelated specs or plan packages
- broad repository scans
- script source files (use readiness scripts instead)

Implementation language discipline:
- Treat the configured code comment language as the default for new or edited code comments
- Preserve an established local file convention when changing comments in existing files
- Avoid mixed-language comments in the same local code area unless there is a strong project reason

Before making meaningful changes:
- Review `.draftspec/constitution.md`
- Inspect the relevant `.draftspec/specs/<slug>.md`
- Inspect the relevant feature package in `.draftspec/plans/<slug>/` when present

After meaningful decisions or changes:
- Keep specs, plans, tasks, archive state, and implementation aligned
