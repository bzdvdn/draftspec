#!/bin/sh

set -eu

SPEC_SLUG="${1:-}"
PLAN_DIR=".draftspec/plans/$SPEC_SLUG"

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
  fi
}

check_file ".draftspec/constitution.md"
check_file ".draftspec/memory.md"
check_file "$PLAN_DIR/tasks.md"
check_file ".draftspec/templates/verify-report.md"
check_file ".draftspec/templates/prompts/verify.md"
