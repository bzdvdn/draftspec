#!/bin/sh

set -eu

SPEC_SLUG="${1:-}"
PLAN_DIR=".draftspec/plans/$SPEC_SLUG"

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
  fi
}

check_file ".draftspec/constitution.md"
check_file ".draftspec/specs/$SPEC_SLUG.md"
check_file "$PLAN_DIR/plan.md"
check_file "$PLAN_DIR/data-model.md"
check_file ".draftspec/templates/prompts/tasks.md"

if [ -d "$PLAN_DIR/contracts" ]; then
  echo "OK: $PLAN_DIR/contracts"
else
  echo "ERROR: missing $PLAN_DIR/contracts"
fi
