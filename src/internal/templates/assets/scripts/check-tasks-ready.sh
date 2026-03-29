#!/bin/sh

set -eu

SPEC_SLUG="${1:-}"
PLAN_DIR=".draftspec/plans/$SPEC_SLUG"
SPEC_FILE=".draftspec/specs/$SPEC_SLUG.md"
PLAN_FILE="$PLAN_DIR/plan.md"
errors=0

if [ -z "$SPEC_SLUG" ]; then
  echo "usage: check-tasks-ready.sh <spec-slug>" >&2
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
check_file "$PLAN_FILE"
check_file "$PLAN_DIR/data-model.md"
check_file ".draftspec/templates/tasks.md"
check_file ".draftspec/templates/prompts/tasks.md"

if [ -d "$PLAN_DIR/contracts" ]; then
  echo "OK: $PLAN_DIR/contracts"
else
  echo "ERROR: missing $PLAN_DIR/contracts"
  errors=$((errors + 1))
fi

if [ -f "$SPEC_FILE" ]; then
  check_contains "$SPEC_FILE" 'AC-[0-9][0-9][0-9]' "spec has stable acceptance IDs"
fi

if [ -f "$PLAN_FILE" ]; then
  check_contains "$PLAN_FILE" 'DEC-[0-9][0-9][0-9]' "plan has stable decision IDs"
fi

if [ -x ".draftspec/scripts/inspect-spec.sh" ] && [ -f "$SPEC_FILE" ]; then
  ./.draftspec/scripts/inspect-spec.sh "$SPEC_FILE"
fi

if [ "$errors" -ne 0 ]; then
  exit 1
fi
