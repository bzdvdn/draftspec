# Draftspec Constitution Prompt

You are creating or updating `.draftspec/constitution.md` for this project.

## Goal

Produce a strict project constitution that is authoritative for both humans and development agents.

For an existing codebase, formalize the project's observable reality first, then separately codify any new mandatory rules explicitly requested by the user.

## Brownfield mode

When the project already exists, work in two layers:

- `Observed Reality`
  - record only what is supported by repository structure, configuration, key entrypoints, dependencies, and existing documentation
- `Declared Law`
  - record new mandatory development rules only when they are explicitly requested by the user or already strongly grounded in the project

Do not redesign an existing project into an idealized architecture. Describe current reality first, then formalize how future changes must be governed.

## Load First

- current user request and conversation
- `.draftspec/constitution.md`
- top-level directory structure
- only the smallest amount of repository context needed to make the constitution concrete

## Load If Present

- `README.md`
- `AGENTS.md`
- project manifests and configuration files that quickly explain language, runtime, architectural boundaries, or integrations

## Repository Evidence

Treat these as strong signals:

- explicit directory boundaries such as `api/`, `workers/`, `cmd/`, `internal/`, `contracts/`, `migrations/`
- dependencies and config that clearly reveal transports, storage systems, or runtime shape
- existing documents that already define workflow or architectural boundaries
- key entrypoint files that show the composition root, process model, or role separation

Treat these as weak signals:

- isolated files not supported by broader structure
- naming that is not confirmed by configuration or component relationships
- general best-practice expectations not grounded in the repository

Do not derive strict constitutional rules from weak signals alone.

## Do Not Read By Default

- large code areas that do not affect the constitution
- old feature artifacts unless they are required to resolve a constitutional conflict
- the whole repository by default

## Stop Conditions

Stop and ask a minimal follow-up question if:

- the project purpose cannot be stated concretely
- the development workflow rules would be guessed rather than grounded
- a constitutional conflict is visible but cannot be resolved from available context
- architecture boundaries, ownership, or workflow would have to be declared as mandatory without enough evidence

Do not broaden repository reading unless one of these conditions is met.

If the constitution is already current and does not conflict with the request, say so and do not modify the file.

## Required behavior

- Work by patching the existing `.draftspec/constitution.md` file.
- Preserve these mandatory sections exactly:
  - `## Purpose`
  - `## Core Principles`
  - `## Constraints`
  - `## Decision Priorities`
  - `## Key Quality Dimensions`
  - `## Language Policy`
  - `## Development Workflow`
  - `## Governance`
  - `## Exceptions Protocol`
  - `## Last Updated`
- Ensure there are at least 5 principle subsections under `## Core Principles` using `### Principle Name` headings.
- You may add extra sections when they materially improve project governance.
- Replace placeholder tokens like `[PROJECT_NAME]` or `[PRINCIPLE_1_NAME]` with concrete text.
- For a brownfield project, codify what the codebase already demonstrates before introducing new mandatory norms.
- If the user explicitly requests new development rules, encode them in `## Development Workflow` and `## Governance` as mandatory rules for future work.
- The `## Development Workflow` section MUST define how feature branches, specs, inspect, plans, tasks, and implementation relate to constitutional compliance.
- The `## Development Workflow` section MUST explicitly state the conditions under which a spec, inspect, plan, tasks, or implementation violates the constitution and cannot proceed.
- The `## Decision Priorities` section MUST capture 3-5 short, rule-like priorities for resolving trade-offs such as simplicity vs extensibility, correctness vs delivery speed, or maintainability vs cleverness.
- The `## Key Quality Dimensions` section MUST include only project-relevant quality dimensions. Do not write a generic quality essay; keep it to 3-5 short, testable bullets.
- The `## Exceptions Protocol` section MUST explain how acceptable deviations from the constitution are recorded and when downstream phases should treat a conflict as a blocker.
- Do not declare DDD boundaries, event-contract ownership, release policy, or branch strategy as mandatory unless they are repository-grounded or explicitly requested by the user.
- If critical information is missing, ask only the minimum necessary follow-up questions.
- Use strict, testable language. Avoid vague wording. Each principle must make it possible to answer concretely: "does this decision conform to the constitution?"
- Do not turn `## Decision Priorities`, `## Key Quality Dimensions`, or `## Exceptions Protocol` into a long handbook. Prefer compact bullets that are useful for downstream phase checks.
- The constitution is the highest-priority project document. Specs, inspection reports, plans, tasks, and implementation must conform to it.

## Update rules

- Keep existing good principles unless they conflict with new requirements.
- Prefer patching and refinement over full rewrites. When refining a principle, preserve its testability; do not replace concrete rules with vague generalizations.
- If a rule is inferred from the repository, phrase it as an observed stable norm of the project rather than an abstract best practice.
- If a rule is introduced by user intent, phrase it as law for downstream phases.
- Update `## Last Updated` with today's date in `YYYY-MM-DD` format whenever the constitution changes.

## Summary artifact

After writing or patching `constitution.md`, also write or update `.draftspec/constitution.summary.md`.

The summary MUST contain only:

- `## Purpose` — 1-2 sentences
- `## Key Constraints` — 3-5 bullets, hard non-negotiable limits only
- `## Language Policy` — 3 lines: docs, agent, comments
- `## Development Workflow` — 3-5 key rules most relevant to spec, plan, and implement phases
- `## Decision Priorities` — 3-5 bullets

Keep the summary under 60 lines. It is loaded by `implement`, `tasks`, and `verify` instead of the full constitution to reduce context overhead. The summary is not a substitute for the full constitution in phases that require constitutional consistency checks (spec, inspect, plan).

## Output expectations

- Write the updated `.draftspec/constitution.md` and `.draftspec/constitution.summary.md`
- Briefly summarize what changed and what remains unresolved
- Note what was inferred from the codebase and what was added as new mandatory law
- Mark unresolved constitutional questions as **BLOCKER** for downstream phases
