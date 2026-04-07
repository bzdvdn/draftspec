# Draftspec Verify Prompt

You are verifying one implemented feature package after task execution.

## Goal

Confirm whether the implemented work is aligned enough with tasks and project rules to proceed safely.

## Flags

`--deep`: full implementation validation mode — read all plan artifacts and inspect actual code for every completed task and acceptance criterion, not just structural checks. Produces a comprehensive report with per-AC evidence. Without this flag, verification stays structural and cheap by default.

`--persist`: write the verification report to `.draftspec/plans/<slug>/verify.md` in addition to conversation output. Without this flag, the report is output to the conversation only. When `--persist` is present, use `.draftspec/templates/verify-report.md` as the canonical template and include the machine-readable metadata block.

## Phase Contract

Inputs: `.draftspec/constitution.md`, `.draftspec/plans/<slug>/tasks.md`; spec, plan, code only to confirm concrete claims (or all artifacts in `--deep` mode).
Outputs: verdict report (`pass`, `concerns`, or `blocked`) in conversation; persisted to `.draftspec/plans/<slug>/verify.md` on request.
Stop if: slug ambiguous, tasks.md missing, or verdict would require inventing implementation facts.

## Load First

Always read these first:

- `.draftspec/constitution.summary.md` if present; otherwise `.draftspec/constitution.md`
- `.draftspec/plans/<slug>/tasks.md`

## Load If Present

Read when a specific check references content in these files (e.g., a task claims to satisfy an `AC-*`, or a `DEC-*` constrains implementation shape):

- `.draftspec/specs/<slug>/summary.md` (or `spec.md`) — when verifying acceptance coverage or task-to-AC alignment
- `.draftspec/plans/<slug>/plan.md` — when a task references a `DEC-*` or architectural decision that must be confirmed
- `.draftspec/plans/<slug>/data-model.md` — when a task touches persisted state or entity shape
- `.draftspec/plans/<slug>/contracts/` — when a task touches API or event boundaries
- `.draftspec/plans/<slug>/research.md` — only when the check depends on a documented trade-off or external dependency finding
- code files — only the specific files named in a task's `Touches:` that are needed to confirm the task was actually implemented

## Do Not Read By Default

- unrelated code areas
- broad repository history
- archives unless the current verification explicitly depends on them

## Stop Conditions

Stop and ask for clarification only if:

- the slug is ambiguous
- the tasks file is missing
- the verification would otherwise invent implementation facts
- the requested conclusion would require a broad repository sweep instead of focused evidence for this feature package
- the implementation claim cannot be confirmed from the current tasks, plan artifacts, and targeted code inspection

## Rules

- Start from `tasks.md` as the verification entrypoint.
- If `/.draftspec/scripts/check-verify-ready.*` is available, prefer it as the cheap first pass before reading deeper artifacts.
- Use `/.draftspec/scripts/verify-task-state.*` only as a fallback when the phase-readiness wrapper is unavailable.
- Prefer helper script output over reading helper script source.
- Do not read `/.draftspec/scripts/*` by default unless you are debugging the script, working on Draftspec itself, or the user explicitly asks to inspect script logic.
- Prefer confirming concrete implementation claims over broad subjective review.
- Treat verify as an evidence log, not a reassurance ritual.
- Verify that completed tasks are consistent with the current state of the feature package.
- **Traceability Evidence**: Use `/.draftspec/scripts/trace.* <slug>` to scan for `@ds-task` and `@ds-test` annotations in the code. Include these findings in the `## Checks` section as concrete implementation evidence.
- **Legacy Fallback**: If `trace` returns no findings (e.g., for older features without annotations), you MUST proceed with manual inspection of the implementation files listed in `Touches:` and run relevant tests to confirm the implementation claims. Note the lack of annotations as a minor warning in the report.
- Verify that open tasks do not contradict any claim that the feature is fully complete.
- Verify acceptance-to-task coverage consistency when `tasks.md` includes an `Acceptance Coverage` section.
- When `tasks.md` uses task IDs such as `T1.1`, reference those IDs directly in checks, findings, and conclusions.
- Prefer `concerns` over `pass` when the evidence is partial but no contradiction has been found.
- Keep default verification structural and cheap by default.
- When `--deep` is present in `$ARGUMENTS`, switch to full validation mode:
  - Read all plan artifacts (`plan.md`, `data-model.md`, `contracts/`, `research.md`).
  - For every completed task, read the actual implementation files listed in `Touches:` and confirm the work matches the task description.
  - For every `AC-*`, trace through the code to confirm the acceptance criterion is satisfied with concrete evidence.
  - The `## Scope` section must state `mode: deep` and list all surfaces inspected.
  - The `## Not Verified` section should be minimal or `none` — deep mode is expected to be thorough.
- Without `--deep`, only deepen into broader implementation validation when a concrete contradiction cannot be resolved from tasks, plan artifacts, and focused evidence.
- Use a simple verdict: `pass`, `concerns`, or `blocked`.
- Use `pass` when no blocking problems are present and only minor or no warnings remain.
- Use `concerns` when the feature can move forward, but warnings or open questions should be resolved soon.
- Use `blocked` when missing task completion or contradictory implementation state would make archive or completion claims unsafe.
- Do not use `pass` unless the completed task state is confirmed, no blocking contradiction remains, and every acceptance or implementation claim you mention is backed by inspected evidence.
- Keep the verification output in the project's configured documentation language when writing it to disk.
- Use `.draftspec/templates/verify-report.md` as the canonical template when writing the report to disk.
- When writing the report to disk, include a machine-readable metadata block at the top with `report_type`, `slug`, `status`, `docs_language`, and `generated_at`.
- Use this report structure:
  - YAML-style metadata block at the top
  - `# Verify Report: <slug>`
  - `## Scope`
  - `## Verdict`
  - `## Checks`
  - `## Errors`
  - `## Warnings`
  - `## Questions`
  - `## Not Verified`
  - `## Next Step`
- In `## Scope`, record the actual verification mode and the surfaces you really inspected.
- In `## Verdict`, include `archive_readiness` and a one-line summary that explains why the verdict is justified.
- In `## Checks`, explicitly cover:
  - `task_state` with completed/open counts
  - `acceptance_evidence` for the `AC-*` items you actually confirmed
  - `implementation_alignment` with the concrete surface inspected
- In `## Not Verified`, list any material claims or surfaces you intentionally did not check. Use `none` only when no material gaps remain inside the chosen verification scope.
- Keep claims scoped. If you only checked task state plus one endpoint or file path, say that directly instead of implying full feature validation.
- If verification discovers a workflow gap, send the feature back to the narrowest earlier phase that can honestly fix it:
  - `implement` for missing or contradictory implementation
  - `tasks` for missing, misleading, or incomplete task decomposition
  - `plan` when the implementation cannot be judged honestly because the design intent is underspecified
- For `pass`, name the exact archive command.
- For `concerns`, say whether the workflow may continue; if it may not, use an explicit return command for the earlier phase.
- For `blocked`, do not suggest archive; end with `Return to: /draftspec.<phase> <slug>` for the narrowest honest recovery phase.

## Output expectations

- Output the report to the conversation by default; persist to `.draftspec/plans/<slug>/verify.md` when `--persist` is present or the user explicitly asks
- Summarize the verdict, completed checks, remaining concerns, and whether the feature is safe to archive
- End with a summary block: `Slug`, `Status`, `Artifacts`, `Blockers`, and either `Next command` or `Return to`
- When safe to archive: `Next command: /draftspec.archive <slug>`; when returning to an earlier phase, name it explicitly with its slash command

## Self-Check

- Is every verdict claim backed by inspected evidence, not just checkbox state?
- Is the `Not Verified` section honest about what I did not check?
- Is the next step or return phase appropriate for the verdict?
