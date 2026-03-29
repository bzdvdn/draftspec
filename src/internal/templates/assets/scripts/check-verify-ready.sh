#!/bin/sh

set -eu

SPEC_SLUG="${1:-}"
PLAN_DIR=".draftspec/plans/$SPEC_SLUG"
SPEC_FILE=".draftspec/specs/$SPEC_SLUG.md"
TASKS_FILE="$PLAN_DIR/tasks.md"
errors=0

if [ -z "$SPEC_SLUG" ]; then
  echo "usage: check-verify-ready.sh <spec-slug>" >&2
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

check_file ".draftspec/constitution.md"
check_file "$SPEC_FILE"
check_file "$TASKS_FILE"
check_file ".draftspec/templates/verify-report.md"
check_file ".draftspec/templates/prompts/verify.md"

if [ -x ".draftspec/scripts/inspect-spec.sh" ] && [ -f "$SPEC_FILE" ] && [ -f "$TASKS_FILE" ]; then
  ./.draftspec/scripts/inspect-spec.sh "$SPEC_FILE" "$TASKS_FILE"
fi

if [ -x ".draftspec/scripts/verify-task-state.sh" ] && [ -f "$TASKS_FILE" ]; then
  ./.draftspec/scripts/verify-task-state.sh "$SPEC_SLUG"
fi

if [ "$errors" -ne 0 ]; then
  exit 1
fi
