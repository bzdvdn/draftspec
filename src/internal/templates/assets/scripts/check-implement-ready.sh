#!/bin/sh

set -eu

SPEC_SLUG="${1:-}"
PLAN_DIR=".draftspec/plans/$SPEC_SLUG"
SPEC_FILE=".draftspec/specs/$SPEC_SLUG.md"
TASKS_FILE="$PLAN_DIR/tasks.md"
errors=0

if [ -z "$SPEC_SLUG" ]; then
  echo "usage: check-implement-ready.sh <spec-slug>" >&2
  exit 1
fi

check_file() {
  path="$1"
  if [ -f "$path" ]; then
    echo "OK: $path"
  else
    echo "ERROR: missing $path"
    errors=$((errors + 1))
  fi
}

check_contains() {
  path="$1"
  pattern="$2"
  label="$3"
  if grep -Eq "$pattern" "$path"; then
    echo "OK: $label"
  else
    echo "ERROR: $label"
    errors=$((errors + 1))
  fi
}

check_file ".draftspec/constitution.md"
check_file "$SPEC_FILE"
check_file "$PLAN_DIR/plan.md"
check_file "$TASKS_FILE"
check_file "$PLAN_DIR/data-model.md"
check_file ".draftspec/templates/prompts/implement.md"

if [ -d "$PLAN_DIR/contracts" ]; then
  echo "OK: $PLAN_DIR/contracts"
else
  echo "ERROR: missing $PLAN_DIR/contracts"
  errors=$((errors + 1))
fi

if [ -f "$TASKS_FILE" ]; then
  check_contains "$TASKS_FILE" '^## (Покрытие критериев приемки|Acceptance Coverage)$' "tasks include acceptance coverage section"
  check_contains "$TASKS_FILE" 'T[0-9]+\.[0-9]+' "tasks include phase-scoped task IDs"
  check_contains "$TASKS_FILE" 'AC-[0-9][0-9][0-9].*->.*T[0-9]+\.[0-9]+' "tasks include AC-to-task coverage lines"
fi

if [ -x ".draftspec/scripts/inspect-spec.sh" ] && [ -f "$SPEC_FILE" ] && [ -f "$TASKS_FILE" ]; then
  ./.draftspec/scripts/inspect-spec.sh "$SPEC_FILE" "$TASKS_FILE"
fi

if [ -x ".draftspec/scripts/check-constitution.sh" ]; then
  ./.draftspec/scripts/check-constitution.sh .draftspec/constitution.md
fi

if [ "$errors" -ne 0 ]; then
  exit 1
fi
