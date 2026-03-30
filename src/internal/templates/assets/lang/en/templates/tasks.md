# <Spec Title> Tasks

## Phase Contract

Inputs: plan and minimal supporting artifacts for this feature.
Outputs: ordered executable tasks with coverage mapping.
Stop if: tasks would be vague or acceptance coverage cannot be mapped.

## Phase 1: Foundation

- [ ] T1.1 Establish the basic feature scaffold — the implementation entrypoint exists and matches the plan scope

## Phase 2: Core Implementation

- [ ] T2.1 Implement the primary feature behavior — the main acceptance path works end to end

## Phase 3: Validation

- [ ] T3.1 Validate the feature behavior — tests, checks, or review steps confirm the intended result

## Acceptance Coverage

- AC-001 -> T1.1, T2.1
- AC-002 -> T3.1

## Notes

- Keep task ordering aligned with the plan
- Use phase-scoped task IDs in the form `T<phase>.<index>`
- Keep tasks concrete, measurable, and executable in order
- Reference 1-2 stable IDs per task when possible (`AC-*`, `RQ-*`, `DEC-*`)
- Mark tasks complete as implementation progresses
- Do not leave acceptance criteria without task coverage
